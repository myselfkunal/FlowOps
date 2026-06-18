# Stage 1 — build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-gateway ./services/api-gateway

# Stage 2 — run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
