# go_indexer

indexing PostgreSQL models into Elasticsearch. sqlboiler is used as ORM to fetch datasets with all relationships via eager loading

1. setup infrastructure with *docker-compose* in */infrastructure* folder
2. populate PostgreSQL with random data via *./setup* shell script
3. run go program with *go run main.go*