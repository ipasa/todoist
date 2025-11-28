# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy pkg module first
COPY pkg ./pkg

# Copy service files
COPY auth-service ./auth-service

# Set working directory to service
WORKDIR /build/auth-service

# Download dependencies and tidy
RUN go mod download
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /build/auth-service/main .

# Copy migrations
COPY --from=builder /build/auth-service/infrastructure/persistence/postgres/migrations ./migrations

EXPOSE 8001

CMD ["./main"]
