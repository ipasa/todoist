# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy pkg module first
COPY pkg ./pkg

# Copy service files
COPY task-service ./task-service

# Copy infrastructure files
COPY task-service/infrastructure ./task-service/infrastructure

# Set working directory to service
WORKDIR /build/task-service

# Download dependencies and tidy
RUN go mod download
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /build/task-service/main .

EXPOSE 8002

CMD ["./main"]
