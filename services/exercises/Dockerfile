
# Build stage
FROM golang:1.24-alpine AS builder

# Install git (required for fetching some Go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Run sqlc generate
RUN /go/bin/sqlc generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exercises .

# Final stage
FROM alpine:latest

LABEL org.opencontainers.image.source https://github.com/farzan-kh/guddy-cn

RUN apk --no-cache add ca-certificates postgresql-client python3 py3-pip

# Install Python dependencies
RUN pip3 install --no-cache-dir sqlalchemy python-dotenv psycopg2-binary --break-system-packages

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/exercises .

# Copy migration files and data
COPY --from=builder /app/internal/db/migrate ./internal/db/migrate/
COPY --from=builder /app/internal/db/data ./internal/db/data/

# Expose port 8081
EXPOSE 8081

# Run the binary
CMD ["./exercises"]
