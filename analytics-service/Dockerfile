# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum* ./

# Copy source code
COPY . .

# Configure Go proxy and checksum database with fallback
ENV GOPROXY=https://proxy.golang.org,https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
ENV GOPRIVATE=
ENV GOTIMEOUT=300s

# Download dependencies with retry logic
RUN set -eux; \
    for i in 1 2 3 4 5; do \
        if go mod download; then \
            break; \
        else \
            echo "Download attempt $i failed, retrying in 10 seconds..."; \
            sleep 10; \
        fi; \
    done; \
    echo "Running go mod tidy..."; \
    for i in 1 2 3 4 5; do \
        if go mod tidy; then \
            break; \
        else \
            echo "Tidy attempt $i failed, retrying in 10 seconds..."; \
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

# Copy the binary from builder
COPY --from=builder --chown=appuser:appuser /app/main .

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8082

# Run the binary
CMD ["./main"]

