# Go E-Commerce Project

[![Project Test Coverage](https://codecov.io/gh/sagar23sj/go-ecommerce/branch/main/graph/badge.svg?token=NKLOFENXKD)](https://codecov.io/gh/sagar23sj/go-ecommerce)

## Problem Statement

<p>
Please find the below two services and the operations that are allowed.

Product Service: provides information about the product like availability, price, category

Order service: provides information about the order like orderValue, dispatchDate, orderStatus, prodQuantity

The user should be able to get the product catalogue and using that info should be able to place an order.

Once the order is placed for a particular product, the product catalogue should be updated accordingly.
(Max quantity of a particular product that can be ordered is 10)
If the order contains 3 premium different products, order value should be discounted by 10%

The Order service should be able to update the orderStatus for a particular order.
dispatchDate should be populated only when the orderStatus is 'Dispatched'.


product category values: Premium/Regular/Budget
order status values: Placed/Dispatched/Completed/Returned/Cancelled


<b>PS: Added a minor change to have Returned state as well</b>
</p>


## Setup

This Project uses key-value store BoltDB and storm toolkit to handle database queries.
There are 10 products already seeded into database and whatever updations you make on database, it will persist even after you close the application. You can run the CleanUp command to start fresh.


Firstly, run the following command to download all dependencies
```bash
go mod download
```


1. Run following command to start e-commerce Application
```bash
make run
```

2. Run following command to run unit test cases
```bash
make test
```

3. Run following command to check test coverage
```bash
make test-cover

#you can also check code test coverage on top. Click on codeccov badge to check more about test coverage
```

4. Run following command to erase database to start fresh
```bash
make clean 
```


## APIs


1. <b>List Products API</b> : `GET http://localhost:8080/products`
2. <b>Get Products Details API</b> : `GET http://localhost:8080/products/{product_id}`
3. <b>Create order API</b> : `POST http://localhost:8080/orders`
4. <b>Get Order Details API</b> : `GET http://localhost:8080/orders/{order_id}`
5. <b>List Orders API</b> : `GET http://localhost:8080/orders`
6. <b>Updare Order Status API</b> : `PATCH http://localhost:8080/orders/{order_id}/status`

## Postman Collection


[here](Go-E-Commerce.postman_collection.json)


## Project Structure

```
.
├── Go-E-Commerce.postman_collection.json
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── coverage.out
├── go.mod
├── go.sum
└── internal
    ├── api
    │   ├── order.go
    │   ├── order_test.go
    │   ├── product.go
    │   ├── product_test.go
    │   └── router.go
    ├── app
    │   ├── dependencies.go
    │   ├── order
    │   │   ├── domain.go
    │   │   ├── domain_test.go
    │   │   ├── mocks
    │   │   │   └── Service.go
    │   │   ├── service.go
    │   │   └── service_test.go
    │   └── product
    │       ├── domain.go
    │       ├── mocks
    │       │   └── Service.go
    │       ├── service.go
    │       └── service_test.go
    ├── pkg
    │   ├── apperrors
    │   │   ├── errors.go
    │   │   ├── map_errors.go
    │   │   ├── order.go
    │   │   └── product.go
    │   ├── constants
    │   │   └── app.go
    │   ├── dto
    │   │   ├── order.go
    │   │   └── product.go
    │   ├── logger
    │   │   └── logger.go
    │   └── middleware
    │       └── response_writer.go
    └── repository
        ├── boltdb
        │   ├── base.go
        │   ├── order.go
        │   ├── order_items.go
        │   └── product.go
        ├── init.go
        ├── mocks
        │   ├── OrderItemStorer.go
        │   ├── OrderStorer.go
        │   └── ProductStorer.go
        ├── order.go
        ├── order_items.go
        ├── products.go
        └── repo.go
```