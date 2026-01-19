# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Run go mod tidy and download
RUN go mod tidy
RUN go mod download

# Copy source code
COPY . .

# Build binary - correct path
RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver ./cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

#change the working directory to root
# this is the default working directory for the container
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/httpserver .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run
CMD ["./httpserver"]