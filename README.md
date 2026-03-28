# Go Microservices E-Commerce

A distributed e-commerce platform built with Go, demonstrating modern microservices architecture using gRPC for inter-service communication and GraphQL as a unified API gateway.

## Overview

This project implements a scalable e-commerce backend system following microservices principles. Each service is independently deployable and maintains its own database, ensuring loose coupling and high cohesion. The system uses gRPC for efficient service-to-service communication and exposes a GraphQL API for client applications.

## Features

- **Account Management**: Create and manage user accounts with unique identifiers
- **Product Catalog**: Full-text search capabilities for product discovery with Elasticsearch
- **Order Processing**: Complete order lifecycle management with transaction support
- **Unified API Gateway**: Single GraphQL endpoint for all client operations
- **Service Independence**: Each microservice operates autonomously with its own database
- **Connection Resilience**: Automatic retry logic for database connections
- **Pagination Support**: Built-in pagination for all list operations
- **GraphQL Playground**: Interactive API explorer for testing and development

## Architecture

### Microservices

The system consists of four main services:

#### 1. Account Service
- Manages user account information
- PostgreSQL database for persistent storage
- Exposes gRPC endpoints for account operations
- Supports account creation and retrieval with pagination

#### 2. Catalog Service
- Handles product catalog management
- Elasticsearch for full-text search capabilities
- Multi-match search across product names and descriptions
- Bulk product retrieval for optimized queries

#### 3. Orders Service
- Processes and manages customer orders
- PostgreSQL database with relational integrity
- Validates orders against Account and Catalog services
- Automatic price calculation and transaction management

#### 4. GraphQL API Gateway
- Unified entry point for all client requests
- Aggregates data from backend microservices
- Provides flexible querying capabilities
- Interactive playground at `/playground`

### Technology Stack

**Core Technologies**
- Go 1.26.1
- gRPC 1.79.3
- Protocol Buffers 1.36.11
- GraphQL (gqlgen 0.17.88)

**Databases**
- PostgreSQL 18.3 (Account and Orders services)
- Elasticsearch 8.17.1 (Catalog service)

**Key Libraries**
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/olivere/elastic/v7` - Elasticsearch client
- `github.com/kelseyhightower/envconfig` - Environment configuration
- `github.com/segmentio/ksuid` - Unique identifier generation
- `github.com/tinrab/retry` - Connection retry logic

**Infrastructure**
- Docker and Docker Compose for containerization
- Multi-stage Docker builds for optimized images

### Communication Patterns

- **Synchronous Communication**: gRPC for service-to-service calls
- **API Gateway Pattern**: GraphQL aggregates backend services
- **Database Per Service**: Each service maintains its own data store
- **Repository Pattern**: Clean separation of data access logic

## Prerequisites

- Docker (version 20.10 or higher)
- Docker Compose (version 2.0 or higher)
- Go 1.26.1 (only required for local development without Docker)

## Getting Started

### Running with Docker Compose

Since the vendor folder is excluded from version control, you need to download dependencies before the first build:

```bash
# Download Go dependencies (creates vendor directory)
go mod vendor

# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up --build -d
```

The GraphQL API will be available at `http://localhost:8000/graphql` and the interactive playground at `http://localhost:8000/playground`.

### Stopping the Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (removes all data)
docker-compose down -v
```

### Viewing Logs

```bash
# View logs from all services
docker-compose logs -f

# View logs from a specific service
docker-compose logs -f graphql
docker-compose logs -f account
docker-compose logs -f catalog
docker-compose logs -f orders
```

## API Usage

### GraphQL Playground

Navigate to `http://localhost:8000/playground` in your browser to access the interactive GraphQL playground. This provides an intuitive interface for exploring the API schema and executing queries.

### Example Queries

**Create an Account**
```graphql
mutation {
  createAccount(account: {name: "John Doe"}) {
    id
    name
  }
}
```

**Create a Product**
```graphql
mutation {
  createProduct(product: {name: "Laptop", description: "High-performance laptop", price: 999.99}) {
    id
    name
    description
    price
  }
}
```

**Create an Order**
```graphql
mutation {
  createOrder(order: {
    accountId: "account-id-here"
    products: [
      {id: "product-id-here", quantity: 1}
    ]
  }) {
    id
    createdAt
    totalPrice
  }
}
```

**Get Account with Orders**
```graphql
query {
  account(id: "account-id-here") {
    id
    name
    orders {
      id
      totalPrice
      createdAt
      products {
        name
        quantity
      }
    }
  }
}
```

**Search Products**
```graphql
query {
  products(query: "laptop", pagination: {skip: 0, take: 10}) {
    id
    name
    description
    price
  }
}
```

## Project Structure

```
.
├── account/                 # Account microservice
│   ├── app.dockerfile      # Service container definition
│   ├── db.dockerfile       # PostgreSQL container with migrations
│   ├── up.sql              # Database schema
│   ├── client.go           # gRPC client
│   ├── server.go           # gRPC server implementation
│   └── pb/                 # Generated protobuf code
├── catalog/                 # Catalog microservice
│   ├── app.dockerfile      # Service container definition
│   ├── client.go           # gRPC client
│   ├── server.go           # gRPC server implementation
│   └── pb/                 # Generated protobuf code
├── orders/                  # Orders microservice
│   ├── app.dockerfile      # Service container definition
│   ├── db.dockerfile       # PostgreSQL container with migrations
│   ├── up.sql              # Database schema
│   ├── client.go           # gRPC client
│   ├── server.go           # gRPC server implementation
│   └── pb/                 # Generated protobuf code
├── graphql/                 # GraphQL API gateway
│   ├── app.dockerfile      # Service container definition
│   ├── graph.go            # Resolver implementations
│   ├── schema.graphql      # GraphQL schema definition
│   └── gqlgen.yml          # Code generation config
├── docker-compose.yaml      # Multi-service orchestration
├── go.mod                   # Go module definition
└── go.sum                   # Dependency checksums
```

## Database Schemas

### Account Service

```sql
CREATE TABLE accounts (
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(24) NOT NULL
);
```

### Orders Service

```sql
CREATE TABLE orders (
    id CHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    account_id CHAR(27) NOT NULL,
    total_price MONEY NOT NULL
);

CREATE TABLE order_products (
    order_id CHAR(27) REFERENCES orders(id) ON DELETE CASCADE,
    product_id CHAR(27),
    quantity INT NOT NULL,
    PRIMARY KEY (product_id, order_id)
);
```

### Catalog Service

Uses Elasticsearch with the following document structure:
- `name` (text): Product name
- `description` (text): Product description
- `price` (float): Product price

## Service Ports

- **GraphQL Gateway**: 8000 (exposed to host)
- **Account Service**: 8080 (internal)
- **Catalog Service**: 8080 (internal)
- **Orders Service**: 8080 (internal)
- **Account Database**: 5432 (internal)
- **Orders Database**: 5432 (internal)
- **Catalog Database**: 9200 (internal)

## Development

### Regenerating Protocol Buffers

If you modify any `.proto` files:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       account/account.proto

protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       catalog/catalog.proto

protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       orders/order.proto
```

### Regenerating GraphQL Code

If you modify `schema.graphql`:

```bash
cd graphql
go run github.com/99designs/gqlgen generate
```

### Adding Dependencies

When adding new Go dependencies:

```bash
# Add dependency to go.mod
go get github.com/package/name

# Update vendor directory
go mod vendor

# Rebuild services
docker-compose up --build
```

## Configuration

All services are configured via environment variables defined in `docker-compose.yaml`:

**Account Service**
- `DATABASE_URL`: PostgreSQL connection string

**Catalog Service**
- `DATABASE_URL`: Elasticsearch HTTP endpoint

**Orders Service**
- `DATABASE_URL`: PostgreSQL connection string
- `ACCOUNT_SERVICE_URL`: Account service gRPC address
- `CATALOG_SERVICE_URL`: Catalog service gRPC address

**GraphQL Service**
- `ACCOUNT_SERVICE_URL`: Account service gRPC address
- `CATALOG_SERVICE_URL`: Catalog service gRPC address
- `ORDER_SERVICE_URL`: Orders service gRPC address

## Limitations and Considerations

- **Authentication**: This project does not implement authentication or authorization. All endpoints are publicly accessible. This is suitable for development and learning purposes but should not be deployed to production without proper security measures.

- **Testing**: No automated tests are currently included. Consider adding unit tests, integration tests, and end-to-end tests for production use.

- **Monitoring**: The system lacks observability tools such as metrics collection, distributed tracing, or centralized logging. For production deployments, consider integrating tools like Prometheus, Jaeger, or ELK stack.

- **Configuration Management**: Database credentials are hardcoded in docker-compose.yaml. Use environment variables or secret management tools in production.

- **Service Discovery**: Service URLs are statically configured. For dynamic environments, consider using service mesh solutions like Istio or Consul.

- **Error Handling**: Error messages are basic and may need enhancement for better debugging and user feedback in production scenarios.

## License

This project is provided as-is for educational and demonstration purposes.

## Author

Shivank Sharma
