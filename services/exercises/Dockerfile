# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install sqlc and compile SQL code
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
	&& /go/bin/sqlc generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exercises .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client python3 py3-pip

# Install Python dependencies
RUN pip3 install --no-cache-dir sqlalchemy python-dotenv psycopg2-binary --break-system-packages

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/exercises .

# Copy migration files and data
COPY --from=builder /app/internal/db/migrate ./internal/db/migrate/
COPY --from=builder /app/internal/db/data ./internal/db/data/

# Copy the .env file into the Docker image
COPY .env .env

# Expose port 8081
EXPOSE 8081

# Run the binary
CMD ["./exercises"]
