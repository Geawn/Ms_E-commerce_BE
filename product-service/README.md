# Product Service

GraphQL service for managing products in BK Commerce.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL (if running locally)

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
3. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Service

1. Start PostgreSQL:
   ```bash
   docker-compose up -d
   ```

2. Run the service:
   ```bash
   go run src/main.go
   ```

3. Access GraphQL Playground:
   - Open http://localhost:8080 in your browser
   - You can test queries and mutations in the playground

## Example Queries

### Get Product by Slug
```graphql
query {
  product(slug: "example-product", channel: "default-channel") {
    id
    name
    description
    variants {
      id
      name
      pricing {
        price {
          amount
          currency
        }
      }
    }
  }
}
```

### List Products
```graphql
query {
  products(first: 10) {
    edges {
      node {
        id
        name
        description
        pricing {
          priceRange {
            start {
              amount
              currency
            }
            stop {
              amount
              currency
            }
          }
        }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

### Get Categories
```graphql
query {
  categories(first: 10) {
    edges {
      node {
        id
        name
        description
        products {
          id
          name
        }
      }
    }
  }
} 