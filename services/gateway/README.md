# Gateway Microservice

This is the API Gateway microservice for the Guddy application. It acts as a single entry point for all client requests and routes them to the appropriate backend services.

## Features

- **Request Routing**: Routes incoming requests to the correct microservice based on the URL path
- **Header Forwarding**: Preserves all HTTP headers from the client request
- **Query Parameter Forwarding**: Maintains query parameters in the forwarded requests
- **Response Proxying**: Forwards responses back to the client with all headers and status codes intact
- **Logging**: Comprehensive logging using Zap for monitoring and debugging
- **Chi Router**: Uses Chi router for better performance and middleware support

## Architecture

The gateway uses a path-based routing strategy:

- `/api/*` - Routes to the exercises service (handles both exercises and programs)
- `/exercises/*` - Routes to the exercises service
- `/programs/*` - Routes to the exercises service (programs are handled by the exercises service)
- `/docs/*` - Routes to the documentation service
- `/logger/*` - Routes to the logging service

## Service Configuration

Services are configured in the `services` map in `main.go`:

```go
var services = map[string]ServiceConfig{
    "exercises": {Name: "exercises", Host: "localhost", Port: "8081"},
    "programs":  {Name: "exercises", Host: "localhost", Port: "8081"},
    "docs":      {Name: "docs", Host: "localhost", Port: "8082"},
    "logger":    {Name: "logger", Host: "localhost", Port: "8083"},
}
```

## Running the Gateway

1. Ensure all backend services are running on their configured ports
2. Start the gateway:

```bash
go run main.go
```

The gateway will start on port 8080.

## Usage Examples

### Get all exercises
```bash
curl http://localhost:8080/api/exercises
```

### Get exercises with filters
```bash
curl "http://localhost:8080/api/exercises?name=push&limit=10"
```

### Get a program by UUID
```bash
curl http://localhost:8080/api/program/123e4567-e89b-12d3-a456-426614174000
```

### Create a new program
```bash
curl -X POST http://localhost:8080/api/program \
  -H "Content-Type: application/json" \
  -d '[{"exerciseId": 1, "idx": 1, "sets": 3, "reps": 10}]'
```

### Service-specific routing
```bash
# These also work
curl http://localhost:8080/exercises/
curl http://localhost:8080/programs/
```

## Development

To add a new service:

1. Add the service configuration to the `services` map
2. Add a new route handler if needed
3. Update the service to run on the configured port

## Dependencies

- `github.com/go-chi/chi/v5` - Lightweight, idiomatic and composable router for Go
- `go.uber.org/zap` - Structured logging
