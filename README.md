
# ⭐️ [Scalable E-commerce Backend (Microservices)]

## Tech Stack

- Golang
- Nest.js
- gRPC
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
- **Communication**:

  + GraphQL is used for communication between client and server, providing a flexible and efficient API for frontend.

  + gRPC is used for internal communication between microservices, ensuring high performance and type-safe message exchange.
- **Authentication & Authorization**: `JWT` is used for authentication.

## Pre-requisites

- Docker & Docker Compose should be installed.

  ```bash
  docker --version
  docker compose version
  ```

- Create and update all `.env` files with the required values for each service.

## How to run the project using docker (Recommended)

```bash
docker compose up --build
```

Here `--build` is used to build the image again if there are any changes in the code.

## Github Actions (CI/CD) Requirements

- Add `DOCKER_USERNAME` & `DOCKER_PASSWORD` to github secrets to push the image to docker hub.
