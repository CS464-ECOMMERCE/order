# Order Service

The Order Service is responsible for managing orders in the e-commerce application.

## Features

- Place orders from a user's cart
- View order details
- Track orders by user or merchant
- Update order status
- Cancel orders (with inventory management)
- Delete orders

## Endpoints

- `POST /api/orders/checkout` - Place an order
- `GET /api/orders/:id` - Get order details
- `GET /api/orders/user/:userId` - Get orders by user
- `GET /api/orders/merchant/:merchantId` - Get orders by merchant
- `PUT /api/orders/:id/status` - Update order status
- `POST /api/orders/:id/cancel` - Cancel an order
- `DELETE /api/orders/:id` - Delete an order

## Development

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run the service: `go run main.go`

## Docker

Build the Docker image:
```
docker build -f docker/Dockerfile -t order-service .
```

Run the Docker container:
```
docker run -p 8080:8080 order-service
```

## Configuration

The service can be configured using environment variables:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres) 
- `DB_NAME` - Database name (default: order_service)
- `REDIS_ADDR` - Redis address (default: localhost:6379)
- `REDIS_PASSWORD` - Redis password (default: empty)
- `REDIS_DB` - Redis database (default: 0)
- `REDIS_DEFAULT_TTL` - Redis default TTL in seconds (default: 3600)
- `ORDER_PORT` - Order service port (default: 8080)
- `PRODUCT_SERVICE_HOST` - Product service host (default: localhost)
- `PRODUCT_SERVICE_PORT` - Product service port (default: 50051) 