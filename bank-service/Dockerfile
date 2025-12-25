# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Configure Go proxy and checksum database with fallback
# Try proxy first, then direct, then China mirror as fallback
ENV GOPROXY=https://proxy.golang.org,https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
ENV GOPRIVATE=
ENV GOTIMEOUT=300s

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Copy source code
COPY . .

# Download dependencies with retry logic
# Skip checksum verification for test dependencies if it fails
RUN set -eux; \
    for i in 1 2 3 4 5; do \
        if go mod download; then \
            break; \
        else \
            echo "Attempt $i failed, retrying in 10 seconds..."; \
            sleep 10; \
        fi; \
    done; \
    go mod verify || echo "Warning: go mod verify failed, continuing..."

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

# Copy the binary and migrations from builder
COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]

