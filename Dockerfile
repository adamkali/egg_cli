# Multi-stage Dockerfile for egg_cli
# Build stage
FROM golang:1.24.2-alpine AS builder

# Install git for version info and make for build system
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version info using make
RUN make build

# Runtime stage - minimal Alpine image
FROM alpine:latest

# Install ca-certificates for HTTPS requests and git for project generation
RUN apk --no-cache add ca-certificates git

# Create non-root user for security
RUN addgroup -g 1001 -S egg && \
    adduser -u 1001 -S egg -G egg

# Set working directory where projects will be generated
WORKDIR /workspace

# Copy the binary from builder stage
COPY --from=builder /app/egg_cli /usr/local/bin/egg_cli

# Make sure binary is executable
RUN chmod +x /usr/local/bin/egg_cli

# Change ownership of workspace to egg user
RUN chown -R egg:egg /workspace

# Switch to non-root user
USER egg

# Set environment variables for better UX
ENV EGG_CLI_VERSION=container

# Set the binary as entrypoint
ENTRYPOINT ["egg_cli"]

# Default command shows help
CMD ["--help"]

# Usage examples:
# docker build -t egg_cli .
# docker run --rm -v $(pwd):/workspace egg_cli init
# docker run --rm -v $(pwd):/workspace egg_cli generate -c config.yaml