# Multi-stage Dockerfile for PicoClaw
# Produces a minimal runtime image (<15MB) with security hardening

# ============================================================================
# Stage 1: Builder - Compile Go binary
# ============================================================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    make \
    ca-certificates \
    tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first (for better layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build binary with optimizations
# - CGO_ENABLED=0: Static binary (no libc dependency)
# - -ldflags: Strip debug info, set version
# - -trimpath: Remove absolute paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'docker') -X main.buildTime=$(date +%FT%T%z)" \
    -trimpath \
    -o picoclaw \
    ./cmd/picoclaw

# Verify binary size
RUN ls -lh picoclaw && \
    SIZE=$(stat -c %s picoclaw) && \
    echo "Binary size: $((SIZE / 1024 / 1024))MB" && \
    if [ $SIZE -gt 15000000 ]; then \
        echo "WARNING: Binary exceeds 15MB"; \
    fi

# ============================================================================
# Stage 2: Runtime - Minimal Alpine image
# ============================================================================
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && addgroup -g 1000 picoclaw \
    && adduser -D -u 1000 -G picoclaw picoclaw

# Create required directories
RUN mkdir -p /app/workspace /app/skills /app/config && \
    chown -R picoclaw:picoclaw /app

# Copy binary from builder
COPY --from=builder --chown=picoclaw:picoclaw /build/picoclaw /app/picoclaw

# Copy skills directory (optional)
COPY --from=builder --chown=picoclaw:picoclaw /build/skills /app/skills

# Set working directory
WORKDIR /app

# Switch to non-root user
USER picoclaw

# Environment variables
ENV PICOCLAW_HOME=/app \
    PICOCLAW_WORKSPACE=/app/workspace \
    PICOCLAW_PORT=8080

# Expose port (configurable via -p flag)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PICOCLAW_PORT}/healthz || exit 1

# Volume for persistent data
VOLUME ["/app/workspace"]

# Run picoclaw
ENTRYPOINT ["/app/picoclaw"]
CMD ["--config", "/app/config/config.json"]
