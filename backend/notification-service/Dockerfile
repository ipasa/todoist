# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy pkg module first
COPY pkg ./pkg

# Copy service files
COPY notification-service ./notification-service

# Set working directory to service
WORKDIR /build/notification-service

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
COPY --from=builder /build/notification-service/main .

EXPOSE 8004

CMD ["./main"]
