# Recipe Sharing API

Recipe Sharing is a Go backend service for publishing, discovering, rating, and
favoriting recipes.

## Features

- JSON REST API implemented with Go's `net/http`.
- User registration and login.
- JWT bearer-token authentication.
- Create, list, search, view, and delete recipes.
- Rate recipes and view recipe ratings.
- Add and remove favorite recipes and list the authenticated user's favorites.
- PostgreSQL persistence through GORM.
- Automatic database migration on application startup.
- Swagger documentation.
- Docker Compose setup for the API, PostgreSQL, and k6 load testing.

## Technology Stack

- **Language:** Go 1.26
- **HTTP:** `net/http`
- **Database:** PostgreSQL 17
- **ORM/data access:** GORM
- **Authentication:** JWT using HS256
- **API documentation:** Swagger
- **Load testing:** Grafana k6
- **Packaging and local deployment:** Docker and Docker Compose

## Architecture

The service follows a layered architecture:

```text
HTTP request
    |
    v
API / HTTP handlers + authentication middleware
    |
    v
Business services
    |
    v
Repositories / data-access layer
    |
    v
PostgreSQL through GORM
```

The layers are organized as follows:

- `internal/handler`: decodes HTTP requests, validates request-level input,
  maps errors to HTTP responses, and encodes JSON responses.
- `internal/httpapi`: registers routes and composes authentication middleware.
- `internal/auth`: validates bearer tokens and places the authenticated user ID
  in the request context.
- `internal/service`: contains application and domain business rules. Services
  depend on repository interfaces and do not import web-framework packages or
  GORM.
- `internal/repository`: implements persistence operations using GORM.
- `internal/domain`: contains the core data models.
- `docs`: generated Swagger/OpenAPI artifacts.

## Project Structure

```text
.
├── main.go
├── internal/
│   ├── auth/          # JWT authentication middleware
│   ├── config/        # Environment-based configuration
│   ├── domain/        # User, recipe, rating, and favorite models
│   ├── handler/       # HTTP handlers
│   ├── httpapi/       # Router and route composition
│   ├── repository/    # GORM data-access implementations
│   └── service/       # Business logic
├── docs/              # Swagger output
├── scripts/api.js     # k6 load-test scenario
├── compose.yaml       # API, database, and load-test services
└── dockerfile         # Multi-stage API image build
```

## Requirements

For local Docker deployment:

- Docker Engine
- Docker Compose v2

For running the application directly:

- Go 1.26
- PostgreSQL 17 or a compatible PostgreSQL instance

## Configuration

The application reads these environment variables:

| Variable | Description | Example |
| --- | --- | --- |
| `HTTP_ADDRESS` | HTTP listen address | `0.0.0.0:8080` |
| `JWT_KEY` | Secret key used to sign and verify JWTs | `replace-with-a-long-random-secret` |
| `DATABASE_URL` | PostgreSQL connection string | `host=localhost user=postgres password=postgres dbname=recipe_sharing port=5432 sslmode=disable` |

## Run with Docker Compose

Start PostgreSQL and the API:

```bash
docker compose up --build
```

The API is available at:

```text
http://localhost:8080
```

Stop the services:

```bash
docker compose down
```

To also remove the persisted PostgreSQL volume:

```bash
docker compose down -v
```

## Run Locally without Docker

Set the required environment variables and start the application:

```bash
export JWT_KEY='replace-with-a-long-random-secret'
export HTTP_ADDRESS=localhost:8080
export DATABASE_URL='host=localhost user=postgres password=postgres dbname=recipe_sharing port=5432 sslmode=disable'
go run .
```

When using Docker Compose, create a `.env` file in the project root before
starting the services:

```dotenv
JWT_KEY=replace-with-a-long-random-secret
```

## API Documentation

Interactive Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

The Swagger documents are also available in:

- [`docs/swagger.yaml`](docs/swagger.yaml)
- [`docs/swagger.json`](docs/swagger.json)

### Endpoint Summary

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/api/auth/register` | No | Register a user |
| `POST` | `/api/auth/login` | No | Log in and receive a JWT |
| `GET` | `/api/users/me` | Yes | Get the authenticated user's profile |
| `GET` | `/api/recipes` | No | List up to 50 recipes; supports `search`, `userid`, `limit` (maximum 100), and `offset` query parameters |
| `POST` | `/api/recipes` | Yes | Create a recipe owned by the authenticated user |
| `GET` | `/api/recipes/{id}` | No | Get one recipe |
| `DELETE` | `/api/recipes/{id}` | Yes | Delete an owned recipe |
| `POST` | `/api/recipes/{id}/ratings` | Yes | Add a rating from 1 to 5 |
| `GET` | `/api/recipes/{id}/ratings` | No | List ratings for a recipe |
| `POST` | `/api/recipes/{id}/favorites` | Yes | Add a recipe to the user's favorites |
| `DELETE` | `/api/recipes/{id}/favorites` | Yes | Remove a recipe from favorites |
| `GET` | `/api/users/me/favorites` | Yes | List the authenticated user's favorites |

Example recipe creation:

```http
POST /api/recipes
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "title": "Tomato Pasta",
  "description": "A simple pasta recipe"
}
```

Example rating:

```http
POST /api/recipes/1/ratings
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "score": 5,
  "comment": "Easy and delicious"
}
```

## Load Testing on Kaggle CPU

The repository includes a k6 scenario at
[`scripts/api.js`](scripts/api.js). It exercises registration, login, protected
recipe creation, recipe listing/search, recipe lookup, favorites, ratings, and
deletion.

Run it through Docker Compose:

```bash
docker compose --profile loadtest run --rm k6
```

The configured workload is:

- 1 minute ramp-up to 10 virtual users
- 3 minutes at 50 virtual users
- 1 minute ramp-down to 0 virtual users

The scenario defines these acceptance thresholds:

- HTTP request failure rate below 1%
- 95th percentile latency below 500 ms
- 99th percentile latency below 1,000 ms
- Check success rate above 99%

### Baseline Kaggle CPU Benchmark

The following result was collected by running the configured k6 workload on a
Kaggle CPU notebook.

| Metric | Result | Threshold | Status |
| --- | ---: | ---: | --- |
| Checks | 100.00% (21,844 / 21,844) | > 99% | Passed |
| HTTP request failure rate | 0.00% (0 / 16,383) | < 1% | Passed |
| HTTP request rate | 54.59 requests/s | Informational | — |
| Average HTTP request duration | 437.88 ms | Informational | — |
| Median HTTP request duration | 349.73 ms | Informational | — |
| p90 HTTP request duration | 942.25 ms | Informational | — |
| p95 HTTP request duration | 1.14 s | < 500 ms | **Failed** |
| p99 HTTP request duration | 1.51 s | < 1,000 ms | **Failed** |
| Maximum HTTP request duration | 2.31 s | Informational | — |
| Average iteration duration | 1.31 s | Informational | — |
| Iterations | 5,461 | Informational | — |
| Peak virtual users | 50 | Informational | — |

Additional execution results:

- `checks_total`: 21,844
- `http_reqs`: 16,383
- `iterations`: 5,461
- `checks_succeeded`: 100.00%
- `data_received`: 4.3 MB
- `data_sent`: 4.6 MB
