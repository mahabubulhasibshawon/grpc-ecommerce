# Order Service gRPC API

This is a gRPC-based order management service built with Go, following a hexagonal architecture. It provides functionality for user registration, authentication, and order management, with data stored in PostgreSQL. The service uses bcrypt for secure password hashing and JWT for authentication, ensuring protected access to order-related endpoints.

## Features
- **User Management**:
  - **Signup**: Register new users with a username and securely hashed password.
  - **Login**: Authenticate users and issue JWT tokens.
  - **Logout**: Simulate token invalidation (placeholder for token blacklisting).
- **Order Management**:
  - **Create Order**: Create orders with dynamic fee calculations based on city and weight.
  - **List Orders**: Retrieve paginated orders for the authenticated user.
  - **Cancel Order**: Cancel pending orders for the authenticated user.
- **Security**:
  - Passwords are hashed using bcrypt.
  - JWT-based authentication protects endpoints except Signup and Login.
- **Persistence**: PostgreSQL stores users and orders.
- **Validation**: Enforces required fields and phone number format for orders.


## Project Structure
```
order-service/
├── cmd
│   └── server
│       └── main.go
├── Dockerfile
├── er-diagram.png
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── grpc
│   │   │   ├── proto
│   │   │   │   ├── order_grpc.pb.go
│   │   │   │   ├── order.pb.go
│   │   │   │   └── order.proto
│   │   │   └── server.go
│   │   └── repository
│   │       └── postgres.go
│   ├── application
│   │   ├── auth_service.go
│   │   └── order_service.go
│   ├── domain
│   │   └── models.go
│   └── ports
│       └── ports.go
├── note-task.png
├── pkg
│   └── auth
│       └── jwt.go
├── readme.md
└── Task_ Software Engineer (GoLang).md
```

## Workflow
![alt text](note-task.png)

## DB Diagram
![db diagram](er-diagram.png)

## Prerequisites
- **Go**: 1.21 or higher
- **Docker**: For running PostgreSQL and the service
- **PostgreSQL**: For data storage
- **protoc**: For generating gRPC code
- **grpcurl**: For testing gRPC endpoints (optional)

## Setup
1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   cd order-service
   ```

2. **Generate gRPC Code**:
   Ensure `protoc` is installed, then generate Go code from the proto file:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/order.proto
   ```

3. **Run PostgreSQL**:
   Start a PostgreSQL container using Docker:
   ```bash
   docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=grpc-ecommerce
   ```

4. **Set Environment Variables**:
   Configure database connection details:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=grpc-ecommerce
   ```

5. **Build and Run**:
   Build the run the service:
   ```bash
   go mod tidy
   go run cmd/server/main.go
   ```

## API Endpoints
The service exposes the following gRPC endpoints under the `order.OrderService` service, accessible at `localhost:50051`. Use `grpcurl` or a gRPC client to interact with them.

### 1. Signup
- **Purpose**: Register a new user with a username and password (hashed with bcrypt).
- **Request**: `SignupRequest { username, password }`
- **Response**: `SignupResponse { message, type, code }`
- **Authentication**: None (public endpoint)
- **Example**:
  ```bash
  grpcurl -plaintext -d '{"username":"user@example.com","password":"securepass"}' localhost:50051 order.OrderService/Signup
  ```
  **Expected Output**:
  ```json
  {
    "message": "User registered successfully",
    "type": "success",
    "code": 200
  }
  ```
  **Error Cases**:
  - Username already exists: `{ "message": "username already exists", "type": "error", "code": 400 }`
  - Missing fields: `{ "message": "username and password are required", "type": "error", "code": 400 }`

### 2. Login
- **Purpose**: Authenticate a user and return a JWT token.
- **Request**: `LoginRequest { username, password }`
- **Response**: `LoginResponse { token_type, expires_in, access_token, refresh_token, message, type, code }`
- **Authentication**: None (public endpoint)
- **Example**:
  ```bash
  grpcurl -plaintext -d '{"username":"user@example.com","password":"securepass"}' localhost:50051 order.OrderService/Login
  ```
  **Expected Output**:
  ```json
  {
    "tokenType": "Bearer",
    "expiresIn": 432000,
    "accessToken": "<jwt-token>",
    "refreshToken": "dummy-refresh",
    "message": "Logged in",
    "type": "success",
    "code": 200
  }
  ```
  **Error Cases**:
  - Invalid credentials: `{ "message": "invalid credentials", "type": "error", "code": 400 }`

### 3. Create Order
- **Purpose**: Create a new order with recipient details and calculate fees.
- **Request**: `CreateOrderRequest { store_id, merchant_order_id, recipient_name, recipient_phone, recipient_address, recipient_city, recipient_zone, recipient_area, delivery_type, item_type, special_instruction, item_quantity, item_weight, amount_to_collect, item_description }`
- **Response**: `CreateOrderResponse { message, type, code, data }`
- **Authentication**: Requires JWT token in `authorization: Bearer <token>` header
- **Example**:
  ```bash
  grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{
    "store_id": 1,
    "recipient_name": "John Doe",
    "recipient_phone": "01712345678",
    "recipient_address": "123 Main St",
    "recipient_city": 1,
    "recipient_zone": 1,
    "recipient_area": 1,
    "delivery_type": 48,
    "item_type": 2,
    "item_quantity": 5,
    "item_weight": 1.5,
    "amount_to_collect": 1000.0
  }' localhost:50051 order.OrderService/CreateOrder
  ```
  **Expected Output**:
  ```json
  {
    "message": "Order Created Successfully",
    "type": "success",
    "code": 200,
    "data": {
      "consignmentId": "DA251021BNWWN123",
      "merchantOrderId": "",
      "orderStatus": "Pending",
      "deliveryFee": 85.0
    }
  }
  ```
  **Error Cases**:
  - Missing required fields: `{ "message": "missing required fields", "type": "error", "code": 422 }`
  - Invalid phone number: `{ "message": "invalid phone number", "type": "error", "code": 422 }`
  - Unauthorized: `{ "code": 16, "message": "Unauthorized" }` (gRPC status)

### 4. List Orders
- **Purpose**: Retrieve paginated orders for the authenticated user.
- **Request**: `ListOrdersRequest { transfer_status, archive, limit, page }`
- **Response**: `ListOrdersResponse { message, type, code, data }`
- **Authentication**: Requires JWT token
- **Example**:
  ```bash
  grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{"transfer_status":1,"archive":0,"limit":10,"page":1}' localhost:50051 order.OrderService/ListOrders
  ```
  **Expected Output**:
  ```json
  {
    "message": "Orders successfully fetched.",
    "type": "success",
    "code": 200,
    "data": {
      "orders": [
        {
          "orderConsignmentId": "DA251021BNWWN123",
          "orderCreatedAt": "2025-10-21T16:58:00Z",
          "orderDescription": "",
          "merchantOrderId": "",
          "recipientName": "John Doe",
          "recipientAddress": "123 Main St",
          "recipientPhone": "01712345678",
          "orderAmount": 1000.0,
          "totalFee": 85.1,
          "instruction": "",
          "orderTypeId": 1,
          "codFee": 10.0,
          "promoDiscount": 0.0,
          "discount": 0.0,
          "deliveryFee": 75.0,
          "orderStatus": "Pending",
          "orderType": "Delivery",
          "itemType": 2,
          "storeName": "Default Store",
          "storeContactPhone": "123456789",
          "codAmount": 1000.0,
          "deliveryCharge": 75.0,
          "storeId": 1,
          "recipientCity": 1,
          "recipientZone": 1,
          "recipientArea": 1,
          "deliveryType": 48,
          "itemQuantity": 5,
          "itemWeight": 1.5,
          "amountToCollect": 1000.0
        }
      ],
      "total": 1,
      "currentPage": 1,
      "perPage": 10,
      "totalInPage": 1,
      "lastPage": 1
    }
  }
  ```
  **Error Cases**:
  - Unauthorized: `{ "code": 16, "message": "Unauthorized" }`

### 5. Cancel Order
- **Purpose**: Cancel a pending order by consignment ID.
- **Request**: `CancelOrderRequest { consignment_id }`
- **Response**: `CancelOrderResponse { message, type, code }`
- **Authentication**: Requires JWT token
- **Example**:
  ```bash
  grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{"consignment_id":"DA251021BNWWN123"}' localhost:50051 order.OrderService/CancelOrder
  ```
  **Expected Output**:
  ```json
  {
    "message": "Order Cancelled Successfully",
    "type": "success",
    "code": 200
  }
  ```
  **Error Cases**:
  - Order not found or not pending: `{ "message": "order not found, unauthorized, or cannot cancel", "type": "error", "code": 400 }`
  - Unauthorized: `{ "code": 16, "message": "Unauthorized" }`

### 6. Logout
- **Purpose**: Simulate logout (placeholder for token invalidation).
- **Request**: `LogoutRequest {}`
- **Response**: `LogoutResponse { message, type, code }`
- **Authentication**: Requires JWT token
- **Example**:
  ```bash
  grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{}' localhost:50051 order.OrderService/Logout
  ```
  **Expected Output**:
  ```json
  {
    "message": "Successfully logged out",
    "type": "success",
    "code": 200
  }
  ```
  **Error Cases**:
  - Unauthorized: `{ "code": 16, "message": "Unauthorized" }`

## Testing Workflow
1. **Register a User**:
   ```bash
   grpcurl -plaintext -d '{"username":"testuser@example.com","password":"securepass"}' localhost:50051 order.OrderService/Signup
   ```

2. **Login to Get JWT Token**:
   ```bash
   grpcurl -plaintext -d '{"username":"testuser@example.com","password":"securepass"}' localhost:50051 order.OrderService/Login
   ```
   Copy the `accessToken` from the response.

3. **Create an Order**:
   Use the JWT token in the `authorization` header:
   ```bash
   grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{
     "store_id": 1,
     "recipient_name": "John Doe",
     "recipient_phone": "01712345678",
     "recipient_address": "123 Main St",
     "recipient_city": 1,
     "recipient_zone": 1,
     "recipient_area": 1,
     "delivery_type": 48,
     "item_type": 2,
     "item_quantity": 5,
     "item_weight": 1.5,
     "amount_to_collect": 1000.0
   }' localhost:50051 order.OrderService/CreateOrder
   ```
   Note the `consignmentId` from the response.

4. **List Orders**:
   ```bash
   grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{"transfer_status":1,"archive":0,"limit":10,"page":1}' localhost:50051 order.OrderService/ListOrders
   ```

5. **Cancel an Order**:
   Use the `consignmentId` from the create order response:
   ```bash
   grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{"consignment_id":"DA251021BNWWN123"}' localhost:50051 order.OrderService/CancelOrder
   ```

6. **Logout**:
   ```bash
   grpcurl -plaintext -H "authorization: Bearer <jwt-token>" -d '{}' localhost:50051 order.OrderService/Logout
   ```

## Dependencies
- `github.com/golang-jwt/jwt/v5`: JWT token handling
- `golang.org/x/crypto`: bcrypt for password hashing
- `google.golang.org/grpc`: gRPC framework
- `google.golang.org/protobuf`: Protocol Buffers
- `github.com/lib/pq`: PostgreSQL driver

## Security Notes
- **Password Hashing**: Passwords are securely hashed using bcrypt with the default cost factor.
- **JWT Secret**: Currently hardcoded (`your-secret-key`). In production, configure via environment variables.
- **Logout**: Currently a placeholder; implement token blacklisting in production for proper session invalidation.
- **Database**: Uses PostgreSQL with SSL disabled (`sslmode=disable`). Enable SSL in production.

## Notes
- A default user (`01901901901@mailinator.com` / `321dsaf`) is inserted on startup with a hashed password for testing.
- The service listens on port `50051`.
- Fee calculations for orders are based on city (60 for city 1, 100 otherwise) and weight (extra charges for >0.5kg).
- Phone numbers are validated with the regex `^(01)[3-9]{1}[0-9]{8}$`.