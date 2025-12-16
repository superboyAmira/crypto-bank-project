# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git protobuf-dev ca-certificates

# Install protoc-gen-go (compatible with Go 1.24)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# Copy go mod files
COPY go.mod go.sum* ./

# Copy source code
COPY . .

# Generate proto files
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/exchange.proto || true

# Download dependencies and tidy
RUN go mod download && go mod tidy

# Build for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates and create non-root user
RUN apk --no-cache add ca-certificates tzdata || \
    (echo "nameserver 8.8.8.8" > /etc/resolv.conf && apk --no-cache add ca-certificates tzdata)

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser/app

# Copy the binary from builder
COPY --from=builder --chown=appuser:appuser /app/main .

# Switch to non-root user
USER appuser

# Expose ports (gRPC and HTTP)
EXPOSE 9090 8085

# Run the binary
CMD ["./main"]

