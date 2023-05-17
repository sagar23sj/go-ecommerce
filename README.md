# Go E-Commerce Project

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
order status values: Placed/Dispatched/Completed/Cancelled

PS: Also added a Returned order status
</p>


## Setup

This Project uses sqlite database and gorm to handle database queries.
There are 9 products already seeded into database and whatever updations you make to database will persist until cleanup


1. Run following command to start e-commerce Application
```bash
make run
```


2. Run following command to erase database to start fresh
```bash
make clean 
```


## APIs


1. <b>List Products API</b> : `GET http://localhost:8080/products`
2. <b>Get Products Details API</b> : `GET http://localhost:8080/product/1`
3. <b>Create order API</b> : `POST http://localhost:8080/order`
4. <b>Get Order Details API</b> : `GET http://localhost:8080/order/2`
5. <b>List Orders API</b> : `GET http://localhost:8080/orders`
6. <b>Updare Order Status API</b> : `PATCH http://localhost:8080/order/5/status`

## Postman Collection


[here](Go_E-Commerce.postman_collection.json)