FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GGOS=linux go build -o recipe-sharing

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/recipe-sharing .
EXPOSE 8080
CMD ["./recipe-sharing"]