# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates curl

# Create non-root user
RUN adduser -D -u 1000 appuser

# Create directories
RUN mkdir -p /app/uploads /app/temp && chown -R appuser:appuser /app

# Copy binary
COPY --from=builder /server /app/server

# Copy .env.example as reference (actual .env should be mounted)
COPY .env.example /app/.env.example

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1

# Run server
ENTRYPOINT ["/app/server"]