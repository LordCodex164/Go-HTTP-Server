# HTTP Server - Production-Grade Go Web Server

A production-ready HTTP server built from scratch using Go's `net/http` package. Demonstrates best practices for building scalable, observable, and resilient web services.

## Features

### Core Functionality
- ✅ RESTful routing without frameworks
- ✅ Composable middleware architecture
- ✅ Graceful shutdown with signal handling
- ✅ Context-aware request handling
- ✅ Production-ready timeout configuration

### Security & Reliability
- ✅ Panic recovery middleware
- ✅ Rate limiting (per-IP, token bucket algorithm) -(still needs to be fixed)
- ✅ Request timeout enforcement
- ✅ CORS support (ready to add)

### Observability
- ✅ Structured JSON logging
- ✅ Request ID tracking
- ✅ Metrics collection (requests, latency, errors)
- ✅ Prometheus-compatible metrics endpoint
- ✅ Health check endpoint

## Quick Start

### Using Docker (Recommended)
```bash
# Build and run
docker-compose up

# Server available at http://localhost:8080
```

### Local Development
```bash
# Install dependencies
go mod download

# Run server
go run cmd/server/main.go

# With custom config
export SERVER_PORT=3000
export RATE_LIMIT_RPS=50
go run cmd/server/main.go
```

### Build Binary
```bash
go build -o httpserver ./cmd/server
./httpserver
```

## API Endpoints

### Core Endpoints
```bash
# Home
GET /
Response: "Welcome to the HTTP server"

# Health Check
GET /health
Response: {"status":"ok","service":"httpserver"}

# Metrics (Prometheus-style)
GET /metrics
Response: {
  "requests": {"total": 100, "success": 95, ...},
  "latency": {"avg_ms": 12, "min_ms": 1, "max_ms": 450},
  "status_codes": {"200": 95, "500": 5},
  "endpoints": {"/health": 50, "/api/v1/users": 30}
}
```

### Example API Endpoints
```bash
# Get users
GET /api/v1/users
Response: [{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]

# Create user
POST /api/v1/users
Content-Type: application/json
Body: {"name":"Charlie"}
Response: {"id":3,"name":"Charlie"}
```

## Configuration

All configuration via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_READ_TIMEOUT` | `5s` | Max time to read request |
| `SERVER_WRITE_TIMEOUT` | `10s` | Max time to write response |
| `SERVER_IDLE_TIMEOUT` | `120s` | Keep-alive timeout |
| `SERVER_REQUEST_TIMEOUT` | `5s` | Max request processing time |
| `RATE_LIMIT_RPS` | `10` | Requests per second per IP |
| `RATE_LIMIT_BURST` | `20` | Burst capacity |

### Example Configuration
```bash
# .env file
SERVER_PORT=3000
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200
SERVER_REQUEST_TIMEOUT=10s
```

## Architecture

### Project Structure
```
httpserver/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handlers/
│   │   ├── handlers.go          # HTTP handlers
│   │   └── metrics.go           # Metrics endpoint
│   ├── logger/
│   │   └── logger.go            # Structured logging
│   ├── metrics/
│   │   └── metrics.go           # Metrics collection
│   └── middleware/
│       ├── logger.go            # Request logging
│       ├── recovery.go          # Panic recovery
│       ├── requestid.go         # Request ID generation
│       ├── timeout.go           # Request timeout
│       ├── ratelimit.go         # Rate limiting
│       └── responsewriter.go    # Response wrapper
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### Middleware Chain
```
Request
  ↓
Recovery         (catch panics)
  ↓
RequestID        (generate unique ID)
  ↓
RateLimit        (enforce rate limits)
  ↓
Timeout          (enforce timeouts)
  ↓
StructuredLogger (log with metrics)
  ↓
Handler          (business logic)
  ↓
Response
```

## Deployment

### Docker
```bash
# Build image
docker build -t httpserver:latest .

# Run container
docker run -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e RATE_LIMIT_RPS=50 \
  httpserver:latest
```

### Kubernetes (Example)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpserver
spec:
  replicas: 3
  selector:
    matchLabels:
      app: httpserver
  template:
    metadata:
      labels:
        app: httpserver
    spec:
      containers:
      - name: httpserver
        image: httpserver:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_PORT
          value: "8080"
        - name: RATE_LIMIT_RPS
          value: "100"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Cloud Platforms

#### Railway
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and deploy
railway login
railway init
railway up
```

#### Render
1. Connect GitHub repo
2. Select "Web Service"
3. Build command: `go build -o httpserver ./cmd/server`
4. Start command: `./httpserver`

#### Fly.io
```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Launch app
fly launch
fly deploy
```

## Monitoring

### Structured Logs

All logs are JSON formatted for easy parsing:
```json
{
  "timestamp": "2025-01-19T10:30:00Z",
  "level": "INFO",
  "message": "HTTP Request",
  "request_id": "abc-123-def-456",
  "method": "GET",
  "path": "/api/v1/users",
  "status_code": 200,
  "latency": "12.345ms"
}
```

### Metrics Endpoint

Query `/metrics` for real-time application metrics:
```bash
curl http://localhost:8080/metrics | jq
```

### Health Checks
```bash
curl http://localhost:8080/health
```

Use in:
- Load balancers
- Kubernetes probes
- Monitoring systems (Datadog, New Relic)

## Development

### Run Tests
```bash
go test ./... -v
```

### Build
```bash
# Local platform
go build -o httpserver ./cmd/server

# Linux (for deployment)
GOOS=linux GOARCH=amd64 go build -o httpserver ./cmd/server
```

### Load Testing
```bash
# Install hey
go install github.com/rakyll/hey@latest

# Run load test
hey -n 1000 -c 50 http://localhost:8080/health
```

## Production Checklist

- [x] Structured logging
- [x] Metrics collection
- [x] Health checks
- [x] Graceful shutdown
- [x] Rate limiting
- [x] Panic recovery
- [x] Request timeouts
- [x] Configuration via env vars
- [x] Docker support
- [x] Documentation

## What's Next

This server demonstrates HTTP fundamentals. Next steps:

1. **Add Database** - PostgreSQL with CRUD operations
2. **Add Authentication** - JWT tokens, bcrypt passwords
3. **Add Caching** - Redis integration
4. **Add Testing** - Unit tests, integration tests
5. **Add Tracing** - OpenTelemetry integration

## License

MIT License - see LICENSE file

## Author

Built by [Your Name](https://github.com/LordCodex164)

Part of a series learning production-grade Go development.

- Bugs
    - rate limit not working
    - unable to access metrics of few endpoints