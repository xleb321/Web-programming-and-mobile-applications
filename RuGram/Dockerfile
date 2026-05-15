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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o main ./cmd/main.go

# Final stage
FROM alpine:3.19

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates wget && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy binary and migrations
COPY --from=builder /app/main .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations/

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /root && \
    chmod +x /root/main

USER appuser

EXPOSE 4200

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4200/health || exit 1

CMD ["./main"]