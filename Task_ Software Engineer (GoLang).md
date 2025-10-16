---
title: 'Task:  Software Engineer (GoLa'

---

### Task:  Software Engineer (GoLang)

Write a golang application using go or whatever is your preferred framework, but the following features should be present.

- [Login | Optional](#Feature---Login-Optional)
- [List all orders | Required](#Feature---List-all-orders)
- [Create a new order | Required](#Feature---Create-a-new-order)
- [Cancel an order | Optional](#Feature---Cancel-an-order-Optional)
- [Logout | Optional](#Feature---Logout-Optional)


## Flow of the features

- User Has sign up already using email and password
- If the user is not logged in / invalid token
    - If skipping as optional, he can login using given email and password as user input.
    - Otherwise, shows [login error](#Feature---Login-Optional), and logs in the user with provided credentials.
- If the user is logged in, he can see his order
    - Shows the [list of orders](#Feature---List-all-orders)
  
   
- From the list of order, each order [can be cancelled](#Feature---Cancel-an-order-Optional).
    
- User can  [logout of the session](#Feature---Logout-Optional).

## Flow of the features

- User logged to system
- If the user is not logged in or give invalid token
    - If skipping as optional, use email & password to generate JWT token as user input.
    - Otherwise, shows [login ](#Feature---Login-Optional) error.
- If the user is logged in / valid token, 
    - Can fetch the [list of orders](#Feature---List-all-orders)
    - Can  [create a new order](#Feature---Create-a-new-order)
    
- From the list of order, each order [can be cancelled](#Feature---Cancel-an-order-Optional).
    - If an order is cancelled, show success message.
- user can  [logout of the session](#Feature---Logout-Optional).


## Instruction

- You can **SKIP** the *OPTIONAL* features. Complete the required features first. Then if the time allows, complete the optional features as well. For the design/template, copy the design/template from wherever you want.
- Use **git** to submit your code.
- Use Docker


## Suggestions

| Key                    | Value                                 |
| ---------------------  | ------------------------------------- |
| Host/Base URL          | `http://a.xyz`  |
| Email/Username         | `01901901901@mailinator.com`          |
| Password               | `321dsaf`                             |
| BD Number Validation   | `/^(01)[3-9]{1}[0-9]{8}$/`            |

## Feature - Login [Optional]


### JWT TOKEN
```
YOUR JWT TOKEN WILL BE HERE IF YOU ARE NOT IMPLEMENTING IT.
```


### Request

|                  |                         |
| ---------------- | ----------------------- |
| URL              | `{{HOST}}/api/v1/login` |
| Method           | `POST`                  |


```bash
curl --location '{{HOST}}/api/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "01901901901@mailinator.com",
    "password": "321dsa"
}'
```

**NOTE**: Use exactly the same **username** and **password**.

### Response

_Status code_: **200** (Success)

```json
{
    "token_type": "Bearer",
    "expires_in": 432000,
    "access_token": "ACCESS_TOKEN",
    "refresh_token": "REFRESH_TOKEN"
}
```

_Status code_: **400** (Error)
```json
{
    "message": "The user credentials were incorrect.",
    "type": "error",
    "code": 400
}
```

## Feature - Create a new order

### Request

|                  |                                     |
| ---------------- | ----------------------------------- |
| URL              | `{{HOST}}/api/v1/orders`          |
| Method           | `POST`                              |


```bash
curl --location '{{HOST}}/api/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{TOKEN}}' \
--data '{
    "store_id": 131172,
    "merchant_order_id": "[OPTIONAL FIELD] TAKE INPUT FROM USER",
    "recipient_name": "[REQUIRED FIELD] TAKE INPUT FROM USER",
    "recipient_phone": "[REQUIRED FIELD] TAKE INPUT FROM USER, VALIDATE this number if possible",
    "recipient_address": "[REQUIRED FIELD] TAKE INPUT FROM USER",
    "recipient_city": 1,
    "recipient_zone": 1,
    "recipient_area": 1,
    "delivery_type": 48,
    "item_type": 2,
    "special_instruction": "[OPTIONAL FIELD] TAKE INPUT FROM USER",
    "item_quantity": 1,
    "item_weight": 0.5,
    "amount_to_collect": "[REQUIRED FIELD] TAKE INPUT FROM USER",
    "item_description": "[OPTIONAL FIELD] TAKE INPUT FROM USER"
}'
```

#### Instructions
- Fields having **[OPTIONAL FIELD]** denotes that it is not mandatory for the users to provide those values. But **[REQUIRED FIELD]** MUST HAVE to be provided by the user. Without those required values, the form should not be submitted to the Backend server.
- Fields that uses **hardcoded** values (not specified using **[REQUIRED]** or **[OPTIONAL]**) in the example curl request, should be used as it is. Using other values may fail during validation in the backend.
- When taking `recipient_address` field's value from the user input, use "**banani, gulshan 2, dhaka, bangladesh**".
- The number validation is given in the [Data & Credentials](#Data-amp-Credentials) section.
- some clarification, if city id is 1 and weight is <= .5 kg delivery fee will be 60taka, if city id is 1 and weight is > .5 kg and <= 1kg price will be 70taka, after that per kg will cost extra 15taka.
- if city id not 1 then base price will be 100 taka instead of 60taka.
- cod_fee will be 1% of amount_to_collect

### Response

_Status code_: **200** (Success)

```json
{
    "message": "Order Created Successfully",
    "type": "success",
    "code": 200,
    "data": {
        "consignment_id": "{{CONSIGNMENT_ID}}",
        "merchant_order_id": "{{PROVIDED_MERCHANT_ORDER_ID}}",
        "order_status": "Pending",
        "delivery_fee": 60
    }
}
```

_Status code_: **401** (Error)
```json
{
    "message": "Unauthorized",
    "type": "error",
    "code": 401
}
```

_Status code_: **422** (Error)
```json
{
    "message": "Please fix the given errors",
    "type": "error",
    "code": 422,
    "errors": {
        "store_id": [
            "The store field is required",
            "Wrong Store selected"
        ],
        "recipient_name": [
            "The recipient name field is required."
        ],
        "recipient_phone": [
            "The recipient phone field is required."
        ],
        "recipient_address": [
            "The recipient address field is required."
        ],
        "delivery_type": [
            "The delivery type field is required."
        ],
        "amount_to_collect": [
            "The amount to collect field is required."
        ],
        "item_quantity": [
            "The item quantity field is required."
        ],
        "item_weight": [
            "The item weight field is required."
        ],
        "item_type": [
            "The item type field is required."
        ]
    }
}
```


## Feature - List all orders

### User Interface

The following design is an example. You should show the fields/data specified in the image, the UI is not mandatory to be the same. Use any template from the internet if you want.

<!-- ![Screenshot 2024-05-23 at 9.36.34 PM](https://hackmd.io/_uploads/Skiaf1pQR.png) -->
![Screenshot 2024-05-23 at 10.11.54 PM](https://hackmd.io/_uploads/rk8loJTX0.png)

### data need to show
```
consignment id
type
merchant order id
store name
store contact phone
recipient name
recipient phone
recipient address
delivery status
COD amount
Delivery Charge
Discount
```
### Request

|                  |                                                   |
| ---------------- | ------------------------------------------------- |
| URL              | `{{HOST}}/api/v1/orders/all`                    |
| Method           | `GET`                                             |
| Query Parameters | `transfer_status=1&archive=0&limit=10&page=2`     |


```bash
curl --location '{{HOST}}/api/v1/orders/all?transfer_status=1&archive=0&limit=10&page=2' \
--header 'Authorization: Bearer {{TOKEN}}'
```

### Response

_Status code_: **200** (Success, without data in the list)

```json
{
    "message": "Orders successfully fetched.",
    "type": "success",
    "code": 200,
    "data": {
        "data": [],
        "total": 4,
        "current_page": 2,
        "per_page": 4,
        "total_in_page": 0,
        "last_page": 1
    }
}
```

_Status code_: **200** (Success, with data)

```json
{
    "message": "Orders successfully fetched.",
    "type": "success",
    "code": 200,
    "data": {
        "data": [
            {
                "order_consignment_id": "DA230524BNWWN8",
                "order_created_at": "2024-05-23 14:05:34",
                "order_description": "item description",
                "merchant_order_id": "merchant_order_id_4",
                "recipient_name": "TEST recipient name",
                "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
                "recipient_phone": "01901901901",
                "order_amount": 1200,
                "total_fee": 72,
                "instruction": "<special instruction> (optional)",
                "order_type_id": 1,
                "cod_fee": 12,
                "promo_discount": 0,
                "discount": 0,
                "delivery_fee": 60,
                "order_status": "Pending",
                "order_type": "Delivery",
                "item_type": "Parcel",
            }
        ],
        "total": 4,
        "current_page": 1,
        "per_page": 1,
        "total_in_page": 1,
        "last_page": 4
    }
}
```

_Status code_: **401** (Error)
```json
{
    "message": "Unauthorized",
    "type": "error",
    "code": 401
}
```

#### Instructions
- It would be great if you could **paginate**. There will have a **next page** if the value of **response.data.current_page** is less than or equal to **response.data.last_page**.
- In the query parameter, **always** use `transfer_status` to `1`, and `archive` to `0`.
- If you're using pagination, use `limit` greater than or equal to `1`, and `page` greater than or equal to `1`, page value changes based on the current page you're in.
- If you're not doing with pagination, omit `limit` and `page` query paramters.


## Feature - Cancel an order [Optional]

### Request 

|                  |                                                       |
| ---------------- | ----------------------------------------------------- |
| URL              | `{{HOST}}/api/v1/orders/{{CONSIGNMENT_ID}}/cancel`  |
| Method           | `PUT`                                                 |


```bash
curl --location --request PUT '{{HOST}}/api/v1/orders/{{CONSIGNMENT_ID}}/cancel' \
--header 'Authorization: Bearer {{TOKEN}}'
```

### Response

_Status code_: **200** (Success)

```json
{
    "message": "Order Cancelled Successfully",
    "type": "success",
    "code": 200
}
```

_Status code_: **400** (Error)

```json
{
    "message": "Please contact cx to cancel order",
    "type": "error",
    "code": 400
}
```

_Status code_: **401** (Error)

```json
{
    "message": "Unauthorized",
    "type": "error",
    "code": 401
}
```

#### Instructions

- When listing all orders, each order from the list of orders will have the **order_consignment_id**. You will have to use that **order_consignment_id** in the request path.


## Feature - Logout [Optional]


```bash
curl --location --request POST '{{HOST}}/api/v1/logout' \
--header 'authorization: Bearer {{TOKEN}}'
```

### Response

_Status code_: **200** (Success)

```json
{
    "message": "Successfully logged out",
    "type": "success",
    "code": 200
}
```


_Status code_: **401** (Error)

```json
{
    "message": "Unauthorized",
    "type": "error",
    "code": 401
}
```