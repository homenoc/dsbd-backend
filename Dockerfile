# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/backend

# User API stage
FROM alpine:latest AS user-api

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# Default command for user API
CMD ["./main", "start", "user", "--config", "/config/config.json"]

# Admin API stage
FROM alpine:latest AS admin-api

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# Default command for admin API
CMD ["./main", "start", "admin", "--config", "/config/config.json"]
