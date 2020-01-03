# indexer-testing

```sql
INSERT INTO simple.customer(
	customer_id, description)
	VALUES (generate_series(1,10), md5(random()::text));

```

```sql
INSERT INTO simple.order(
	order_id, order_description, customer_id)
	VALUES (generate_series(1,30), md5(random()::text), (random() * 	10 + 1)::int);
```

```sql
INSERT INTO simple.item(
	item_id, item_description, order_id)
	VALUES (generate_series(1,90), md5(random()::text), (random() * 	30 + 1)::int);
```

