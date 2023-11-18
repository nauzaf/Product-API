# PRODUCT API

## Dependencies

- golang 1.21
- postgresql

## Preparation

1. Set up 1.21 golang environment
2. Set up postgresql
3. Set DB_URL inside OS Environment `PRODUCTAPI_DB_URL=postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=disable`

## Run

```
go run main.go
```

The server will run in port 3000

## Endpoint list

```
GET /products
GET /products/{id}
POST /products
POST /products/increase_inventory
POST /products/decrease_inventory
POST /products/set_expired
```

*Create Product Body Request*

```
{
    "name": "product_name"
}
```

*Increase / Decrease / Set Inventory Body Request*

```
{
    "name": "product_id"
}
```
