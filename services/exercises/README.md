# Exercises Service

A Go service for managing exercise data with PostgreSQL database.

## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

## Quick Start with Docker Compose

1. **Start all services from the project root:**
   ```bash
   cd /path/to/guddy-cn
   docker compose up --build
   ```

   This will:
   - Start a PostgreSQL database
   - Run database migrations
   - Import exercise data from `internal/db/data/exercises.json`
   - Start the exercises service on port 8081
   - Start the gateway service on port 8080

2. **Access the services:**
   - Gateway API: `http://localhost:8080`
   - Direct exercises API: `http://localhost:8081`
   - Available endpoints:
     - `GET /api/exercises` - Get all exercises
     - `GET /api/program/{uuid}` - Get a program by UUID
     - `GET /api/completeProgram/{uuid}` - Get complete program details
     - `POST /api/program` - Create a new program

## Local Development

1. **Start PostgreSQL:**
   ```bash
   docker run --name postgres-exercises -e POSTGRES_DB=exercises -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15-alpine
   ```

2. **Run migrations:**
   ```bash
   # Using psql
   psql -h localhost -p 5432 -U postgres -d exercises -f internal/db/migrate/000001_init_schema.up.sql
   psql -h localhost -p 5432 -U postgres -d exercises -f internal/db/migrate/000002_add_exercises_instructions_field.up.sql
   psql -h localhost -p 5432 -U postgres -d exercises -f internal/db/migrate/000003_alter_programs_table.up.sql
   ```

3. **Import data:**
   ```bash
   cd internal/db/data
   pip install sqlalchemy python-dotenv psycopg2-binary
   python import_script.py
   ```

4. **Run the service:**
   ```bash
   go run main.go
   ```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/exercises?sslmode=disable`)

## Database Schema

The service uses the following main tables:
- `exercises`: Exercise definitions with equipment type and instructions
- `exercise_names`: Alternative names for exercises
- `muscles`: Muscle groups
- `exercise_muscle`: Many-to-many relationship between exercises and muscles
- `programs`: Workout programs containing multiple exercises
- `visuals`: Exercise images/videos (currently unused)
