# Build stage
FROM golang:1.25-alpine AS builder

# Install git for version info
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
ARG VERSION=1.1.0
RUN CGO_ENABLED=0 GOOS=linux go build -a \
  -ldflags "-s -w -X main.version=${VERSION} -X main.commit=$(git rev-parse --short HEAD 2>/dev/null || echo docker)" \
  -o amazon-vl ./cmd

# Runtime stage
FROM alpine:3.19

# Security: run as non-root user
RUN addgroup -g 1000 appgroup && \
  adduser -D -u 1000 -G appgroup appuser

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/amazon-vl /usr/local/bin/

# Set ownership
RUN chown appuser:appgroup /usr/local/bin/amazon-vl

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-9000}/healthz || exit 1

# Default environment variables
ENV AUTH_USER=joaquim
ENV AUTH_REALM=amazon-server-logs.com

EXPOSE 9000

ENTRYPOINT ["amazon-vl"]
CMD ["/logs", "9000"]