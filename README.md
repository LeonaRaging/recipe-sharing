# Recipe Sharing API

> [English](#english) | **Tiếng Việt**

## Giới thiệu

Recipe Sharing là dịch vụ backend viết bằng Go, cung cấp REST API để đăng,
khám phá, đánh giá và lưu công thức nấu ăn yêu thích.

## Tính năng

- REST API JSON được xây dựng với `net/http`.
- Đăng ký và đăng nhập người dùng.
- Xác thực bằng bearer token JWT.
- Tạo, liệt kê, tìm kiếm, xem và xóa công thức.
- Đánh giá công thức và xem danh sách đánh giá.
- Thêm, xóa và liệt kê các công thức yêu thích của người dùng đã đăng nhập.
- Lưu trữ dữ liệu trên PostgreSQL thông qua GORM.
- Tự động migration cơ sở dữ liệu khi ứng dụng khởi động.
- Tài liệu API bằng Swagger.
- Docker Compose cho API, PostgreSQL và kiểm thử tải bằng k6.

## Công nghệ sử dụng

- **Ngôn ngữ:** Go 1.26
- **HTTP:** `net/http`
- **Cơ sở dữ liệu:** PostgreSQL 17
- **ORM/truy cập dữ liệu:** GORM
- **Xác thực:** JWT với HS256
- **Tài liệu API:** Swagger
- **Kiểm thử tải:** Grafana k6
- **Đóng gói và chạy cục bộ:** Docker và Docker Compose

## Kiến trúc

Dịch vụ sử dụng kiến trúc phân lớp:

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
PostgreSQL thông qua GORM
```

- `internal/handler`: giải mã request HTTP, kiểm tra dữ liệu đầu vào, ánh xạ
  lỗi thành HTTP response và mã hóa JSON response.
- `internal/httpapi`: đăng ký route và kết hợp authentication middleware.
- `internal/auth`: xác thực bearer token và đưa ID người dùng vào request context.
- `internal/service`: chứa các quy tắc nghiệp vụ của ứng dụng và domain.
- `internal/repository`: triển khai các thao tác lưu trữ bằng GORM.
- `internal/domain`: chứa các model dữ liệu cốt lõi.
- `docs`: các tệp Swagger/OpenAPI được sinh tự động.

## Yêu cầu

Để chạy bằng Docker:

- Docker Engine
- Docker Compose v2

Để chạy trực tiếp:

- Go 1.26
- PostgreSQL 17 hoặc một phiên bản PostgreSQL tương thích

## Cấu hình

Ứng dụng đọc các biến môi trường sau:

| Biến | Mô tả | Ví dụ |
| --- | --- | --- |
| `HTTP_ADDRESS` | Địa chỉ lắng nghe HTTP | `0.0.0.0:8080` |
| `JWT_KEY` | Khóa bí mật dùng để ký và xác thực JWT | `replace-with-a-long-random-secret` |
| `DATABASE_URL` | Chuỗi kết nối PostgreSQL | `host=localhost user=postgres ****** dbname=recipe_sharing port=5432 sslmode=disable` |

## Chạy bằng Docker Compose

Khởi động PostgreSQL và API:

```bash
docker compose up --build
```

API có thể truy cập tại `http://localhost:8080`. Dừng các dịch vụ bằng:

```bash
docker compose down
```

Để xóa cả volume PostgreSQL đã lưu:

```bash
docker compose down -v
```

## Chạy cục bộ không dùng Docker

Thiết lập các biến môi trường cần thiết rồi khởi động ứng dụng:

```bash
export JWT_KEY='replace-with-a-long-random-secret'
export HTTP_ADDRESS=localhost:8080
export DATABASE_URL='host=localhost user=postgres ****** dbname=recipe_sharing port=5432 sslmode=disable'
go run .
```

Khi dùng Docker Compose, tạo tệp `.env` ở thư mục gốc của dự án:

```dotenv
JWT_KEY=replace-with-a-long-random-secret
```

## Tài liệu API

Swagger UI tương tác có tại `http://localhost:8080/swagger/index.html`.
Tệp Swagger cũng có trong [`docs/swagger.yaml`](docs/swagger.yaml) và
[`docs/swagger.json`](docs/swagger.json).

Các endpoint chính:

| Phương thức | Endpoint | Xác thực | Mô tả |
| --- | --- | --- | --- |
| `POST` | `/api/auth/register` | Không | Đăng ký người dùng |
| `POST` | `/api/auth/login` | Không | Đăng nhập và nhận JWT |
| `GET` | `/api/users/me` | Có | Xem hồ sơ người dùng hiện tại |
| `GET` | `/api/recipes` | Không | Liệt kê công thức; hỗ trợ `search`, `userid`, `limit` (tối đa 100) và `offset` |
| `POST` | `/api/recipes` | Có | Tạo công thức thuộc người dùng hiện tại |
| `GET` | `/api/recipes/{id}` | Không | Xem một công thức |
| `DELETE` | `/api/recipes/{id}` | Có | Xóa công thức của người dùng hiện tại |
| `POST` | `/api/recipes/{id}/ratings` | Có | Thêm đánh giá từ 1 đến 5 |
| `GET` | `/api/recipes/{id}/ratings` | Không | Liệt kê đánh giá của một công thức |
| `POST` | `/api/recipes/{id}/favorites` | Có | Thêm công thức vào danh sách yêu thích |
| `DELETE` | `/api/recipes/{id}/favorites` | Có | Xóa công thức khỏi danh sách yêu thích |
| `GET` | `/api/users/me/favorites` | Có | Liệt kê công thức yêu thích của người dùng hiện tại |

## Kiểm thử tải

Kịch bản k6 nằm tại [`scripts/api.js`](scripts/api.js), bao phủ đăng ký,
đăng nhập, tạo và tìm kiếm công thức, yêu thích, đánh giá và xóa công thức.
Chạy bằng Docker Compose:

```bash
docker compose --profile loadtest run --rm k6
```

Tải kiểm thử gồm tăng dần lên 10 người dùng ảo trong 1 phút, duy trì 50 người
dùng ảo trong 3 phút và giảm về 0 trong 1 phút. Ngưỡng chấp nhận là tỷ lệ
request lỗi dưới 1%, p95 dưới 500 ms, p99 dưới 1.000 ms và tỷ lệ check thành
công trên 99%. Kết quả benchmark Kaggle CPU hiện tại đạt 100% check thành công
và 0% request lỗi; p95 và p99 chưa đạt các ngưỡng tương ứng.

> [English](#english)

---

## English

> [Tiếng Việt](#giới-thiệu)

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
