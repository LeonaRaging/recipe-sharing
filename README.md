### Prerequisites
- Docker installed
### How to build
```
docker compose up --build
```
### How to do performance testing
```
docker compose --profile loadtest run --rm k6
```