
# ⭐️ [Scalable E-commerce Backend (Microservices)]

## Tech Stack

- Golang
- Nest.js
- gRPC
- Redis
- GraphQL
- PostgresQL
- JWT
- Docker


## Services & Features

- [x] User Service: Responsible for handling user-related business logic and interactions.
- [x] Product Service: Manages product listings, categories, details, and searching.
- [x] Content Service: Manages webpage information.
- [x] Channel Service: Manages channel listing, information.
- [x] Account Service: Manages User account, including authentication & authorization (JWT), chaging password, Signing up.
- [x] Shopping Cart Service: Manages users’ shopping carts, including adding/removing items and updating quantities , Still working.
- [x] Order Service: Still working.
- [x] Payment Service: Still working.
- [x] Notification Service: Still working.

## Architecture

- **Microservices Architecture**: Each service is a separate codebase, with its own database.
- **Containerization**: `Docker` is used to containerize each service, making it easy to deploy and scale the services.
- **Database**: `PostgresQL` is used as the database for all services, providing a flexible schema and scalability.
- **Cache**: `Redis` is used for caching
- **Communication**:

  + GraphQL is used for communication between client and server, providing a flexible and efficient API for frontend.

  + gRPC is used for internal communication between microservices, ensuring high performance and type-safe message exchange.
- **Authentication & Authorization**: `JWT` is used for authentication.

## Future Plans/Goals

- Finished all services
- Connected to Saleor storefront
- Get product deployed

