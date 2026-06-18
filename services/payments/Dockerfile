# Stage 1 — build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payments ./services/payments

# Stage 2 — run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payments .
EXPOSE 50053
CMD ["./payments"]