# Build stage
FROM docker.io/library/golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/api-server cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/migrate cmd/migrate/main.go

# Runtime stage
FROM docker.io/library/alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Copy binaries from builder
COPY --from=builder /app/bin/api-server /app/api-server
COPY --from=builder /app/bin/migrate /app/migrate

# Copy migrations
COPY migrations /app/migrations

# Create non-root user
RUN adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

# Run migrations then start server
CMD ["/app/api-server"]
