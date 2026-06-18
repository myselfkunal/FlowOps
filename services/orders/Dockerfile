# Stage 1 — build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o orders ./services/orders

# Stage 2 — run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/orders .
EXPOSE 50052
CMD ["./orders"]