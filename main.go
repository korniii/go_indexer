package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/korniii/go_indexer/models"
	_ "github.com/lib/pq" //psql db driver
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

type Customer struct {
	CustomerID          int64    `json:"customer_id"`
	CustomerDescription string   `json:"customer_description`
	Orders              []*Order `json:"orders"`
}

type Order struct {
	OrderID          int64   `json:"order_id"`
	OrderDescription string  `json:"order_description"`
	Items            []*Item `json:"items"`
}

type Item struct {
	ItemID          int64  `json:"item_id"`
	ItemDescription string `json:"item_description"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "example"
)

var (
	querySize int
	batch     int
)

func init() {
	flag.IntVar(&querySize, "querySize", 100, "Number of customers pulled in one query")
	flag.IntVar(&batch, "batch", 255, "Number of documents to send in one batch")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func main() {

	start := time.Now().UTC()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	numOfRecords, err := models.Customers().Count(ctx, db)
	if err != nil {
		log.Fatalln(err)
	}

	convNumOfRecords := int(numOfRecords)

	//ToDo: sent messages to channel
	elMessageChannel := make(chan []*Customer)

	counter := 0

	for i := 0; i < convNumOfRecords; i += querySize {
		log.Printf("Indexing records from %d to %d", i, i+querySize)
		go fetchCustomers(db, ctx, elMessageChannel, querySize, i)
		counter++
	}

	idx := 0
	for data := range elMessageChannel {
		indexCustomers(data)
		idx++
		if (idx == counter) {
			close(elMessageChannel)
		}
	}

	duration := time.Since(start)
	log.Printf("Total time taken -> %s", duration.Truncate(time.Millisecond))

}

func indexCustomers(customers []*Customer) {

	count := len(customers)

	//see https://github.com/elastic/go-elasticsearch/tree/master/_examples/bulk
	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				ID     string `json:"_id"`
				Result string `json:"result"`
				Status int    `json:"status"`
				Error  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
					Cause  struct {
						Type   string `json:"type"`
						Reason string `json:"reason"`
					} `json:"caused_by"`
				} `json:"error"`
			} `json:"index"`
		} `json:"items"`
	}

	var (
		buf bytes.Buffer
		res *esapi.Response
		err error
		raw map[string]interface{}
		blk *bulkResponse

		indexName = "customers"

		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
	)

	//connect to elastic search
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	//(re-)create index
	// if res, err = es.Indices.Delete([]string{indexName}); err != nil {
	// 	log.Fatalf("Cannot delete index: %s", err)
	// }
	// res, err = es.Indices.Create(indexName)
	// if err != nil {
	// 	log.Fatalf("Cannot create index: %s", err)
	// }
	// if res.IsError() {
	// 	log.Fatalf("Cannot create index: %s", res)
	// }

	if count%batch == 0 {
		numBatches = (count / batch)
	} else {
		numBatches = (count / batch) + 1
	}

	start := time.Now().UTC()

	// Loop over the collection
	//
	for i, a := range customers {
		numItems++

		currBatch = i / batch
		if i == count-1 {
			currBatch++
		}

		// Prepare the metadata payload
		//
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, a.CustomerID, "\n"))
		// fmt.Printf("%s", meta) // <-- Uncomment to see the payload

		// Prepare the data payload: encode article to JSON
		//
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", a.CustomerID, err)
		}

		// Append newline to the data payload
		//
		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch
		// fmt.Printf("%s", data) // <-- Uncomment to see the payload

		// // Uncomment next block to trigger indexing errors -->
		// if a.ID == 11 || a.ID == 101 {
		// 	data = []byte(`{"published" : "INCORRECT"}` + "\n")
		// }
		// // <--------------------------------------------------

		// Append payloads to the buffer (ignoring write errors)
		//
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		//
		if i > 0 && i%batch == 0 || i == count-1 {
			log.Printf("> Batch %-2d of %d", currBatch, numBatches)

			res, err = es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
			if err != nil {
				log.Fatalf("Failure indexing batch %d: %s", currBatch, err)
			}
			// If the whole request failed, print error and mark all documents as failed
			//
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					log.Printf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}
				// A successful response might still contain errors for particular documents...
				//
			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {
						// ... so for any HTTP status above 201 ...
						//
						if d.Index.Status > 201 {
							// ... increment the error counter ...
							//
							numErrors++

							// ... and print the response status and error information ...
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// ... otherwise increase the success counter.
							//
							numIndexed++
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			//
			res.Body.Close()

			// Reset the buffer and items counter
			//
			buf.Reset()
			numItems = 0
		}
	}

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	log.Println(strings.Repeat("=", 80))

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%d] documents with [%d] errors in %s (%.0f docs/sec)",
			numIndexed,
			numErrors,
			dur.Truncate(time.Millisecond),
			1000.0/float64(dur/time.Millisecond)*float64(numIndexed),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%d] documents in %s (%.0f docs/sec)",
			numIndexed,
			dur.Truncate(time.Millisecond),
			1000.0/float64(dur/time.Millisecond)*float64(numIndexed),
		)
	}

}

func fetchCustomers(db *sql.DB, ctx context.Context, ch chan<- []*Customer,_limit int, _offset int) {

	custmrs, err := models.Customers(Load("Orders.Items"), Limit(_limit), Offset(_offset)).All(ctx, db)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Data fetched from Postgres")

	var elEntries []*Customer

	for _, custmr := range custmrs {
		var tempOrders []*Order
		for _, ordr := range custmr.R.Orders {
			var tempItems []*Item
			for _, itm := range ordr.R.Items {
				tempItems = append(tempItems, &Item{
					ItemID:          itm.ItemID,
					ItemDescription: itm.ItemDescription.String,
				})
			}
			tempOrders = append(tempOrders, &Order{
				OrderID:          ordr.OrderID,
				OrderDescription: ordr.OrderDescription.String,
				Items:            tempItems,
			})
		}
		elEntries = append(elEntries, &Customer{
			CustomerID:          custmr.CustomerID,
			CustomerDescription: custmr.Description.String,
			Orders:              tempOrders,
		})
	}

	ch <- elEntries
}
