FROM golang:1.21-alpine3.19 AS builder

WORKDIR /app

# Install git and security updates
RUN apk update && apk upgrade && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with security optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o main cmd/main.go

# Final stage - use minimal and secure image
FROM alpine:3.19

# Install security updates and required packages
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations/

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /root && \
    chmod +x /root/main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 4200

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4200/health || exit 1

# Run the application
CMD ["./main"]