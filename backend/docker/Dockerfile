# Build stage
FROM --platform=linux/amd64 golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates build-base make

# Set working directory
WORKDIR /app

# # Set GOPROXY argument with default value
# ARG GOPROXY=https://proxy.golang.org,direct
# ENV GOPROXY=$GOPROXY

# Accept build version as argument
ARG BUILD_VERSION=dev

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init -g cmd/server/main.go -o ./docs/swagger --parseVendor --parseDependency --exclude internal/shared --parseInternal

# Build the application with version information using make
RUN make build VER=$BUILD_VERSION

# Use a small image for the final stage
FROM --platform=linux/amd64 alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/build/context-space-backend .

# Copy configuration files
COPY --from=builder /app/configs ./configs

# Copy swagger docs
COPY --from=builder /app/docs/swagger ./docs/swagger

# Expose port
EXPOSE 8080

# Run the application
CMD ["./context-space-backend"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1 