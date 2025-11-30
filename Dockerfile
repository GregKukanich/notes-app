FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -o /app/main

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/db.sql .

EXPOSE 8080
CMD ["./main"]
