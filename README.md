# Clothing Shop API

## Overview
The Clothing Shop API is a backend service for an online clothing shop application. It provides endpoints for user authentication, product management, and order processing. This API is built using Go and follows a clean architecture pattern.

## Project Structure
```
clothing-shop-api
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── api
│   │   ├── handlers
│   │   │   ├── auth.go
│   │   │   ├── products.go
│   │   │   ├── orders.go
│   │   │   └── users.go
│   │   ├── middleware
│   │   │   ├── auth.go
│   │   │   └── logging.go
│   │   └── routes.go
│   ├── config
│   │   └── config.go
│   ├── domain
│   │   ├── models
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   └── user.go
│   │   └── services
│   │       ├── product_service.go
│   │       ├── order_service.go
│   │       └── user_service.go
│   └── repository
│       ├── product_repository.go
│       ├── order_repository.go
│       └── user_repository.go
├── pkg
│   ├── database
│   │   └── db.go
│   └── utils
│       ├── jwt.go
│       └── password.go
├── migrations
│   └── schema.sql
├── go.mod
├── go.sum
├── .env.example
├── Dockerfile
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.16 or later
- MySQL or any compatible database
- Docker (optional, for containerization)

### Installation
1. Clone the repository:
   ```
   git clone https://github.com/yourusername/clothing-shop-api.git
   cd clothing-shop-api
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up your environment variables. Create a `.env` file based on the `.env.example` provided.

### Database Setup
1. Run the SQL migration to set up the database schema:
   ```
   mysql -u username -p < migrations/schema.sql
   ```

### Running the Application
1. Start the server:
   ```
   go run cmd/server/main.go
   ```

2. The API will be available at `http://localhost:8080`.

### API Endpoints
- **Authentication**
  - POST `/api/auth/login`: Login a user
  - POST `/api/auth/register`: Register a new user

- **Products**
  - GET `/api/products`: List all products
  - GET `/api/products/{id}`: Get product details
  - POST `/api/products`: Add a new product (admin only)

- **Orders**
  - GET `/api/orders`: List all orders (admin only)
  - POST `/api/orders`: Create a new order

### Docker
To build and run the application in a Docker container, use the following command:
```
docker build -t clothing-shop-api .
docker run -p 8080:8080 clothing-shop-api
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.