# Guddy 🏋️‍♂️

A comprehensive fitness exercise and workout program management API built with Go. Guddy provides a robust backend service for managing exercise databases, creating custom workout programs, and accessing detailed exercise information with equipment specifications and visual guides.

## 🌟 Features

- **Exercise Database**: Comprehensive collection of exercises with detailed information
- **Equipment Categorization**: Supports various equipment types (Dumbbells, Barbell, Machine, Bodyweight, etc.)
- **Muscle Group Mapping**: Exercises mapped to specific muscle groups
- **Custom Workout Programs**: Create and manage personalized workout routines
- **Visual References**: Support for exercise demonstration images and videos
- **RESTful API**: Clean and well-documented API endpoints
- **Swagger Documentation**: Interactive API documentation
- **PostgreSQL Backend**: Robust database with proper schema design

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

**Prerequisites**
- Docker and Docker Compose

**Start all services:**
```bash
git clone https://github.com/Farzan-Kh/guddy-cn.git
cd guddy
docker compose up --build
```

This will automatically:
- ✅ Start PostgreSQL database
- ✅ Run database migrations
- ✅ Import exercise data
- ✅ Start the exercises service (port 8081)
- ✅ Start the gateway service (port 8080)

**Access the application:**
- Gateway API: `http://localhost:8080`
- Exercises API: `http://localhost:8081`

### Option 2: Local Development

**Prerequisites**

- Go 1.23.2 or higher
- PostgreSQL database
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations
- [Task](https://taskfile.dev/) for build automation
- [sqlc](https://sqlc.dev/) for SQL code generation
- [Swag](https://github.com/swaggo/swag) for Swagger documentation

**Installation**

1. **Clone the repository**
   ```bash
   git clone https://github.com/Farzan-Kh/guddy-cn.git
   cd guddy
   ```

2. **Set up environment variables**
   ```bash
   cp internal/config/.env.example internal/config/.env
   # Edit the .env file with your database configuration
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Set up the database and build the project**
   ```bash
   task default
   ```

5. **Run the application**
   ```bash
   ./build/guddy
   ```

The API will be available at `http://localhost:3000`

## 📖 API Documentation

Once the application is running, you can access the interactive Swagger documentation at:
```
http://localhost:3000/swagger/index.html
```

### Available Endpoints

- `GET /api/exercises` - Retrieve all exercises with filtering options
- `GET /api/program/{uuid}` - Get a specific workout program by UUID
- `GET /api/fullProgram/{uuid}` - Get complete program details with exercise information
- `POST /api/program` - Create a new workout program

## 🏗️ Project Structure

```
guddy/
├── docker-compose.yml              # Docker Compose configuration
├── main.go                         # Application entry point (legacy)
├── go.mod                         # Go module dependencies
├── sqlc.yaml                      # sqlc configuration
├── Taskfile.yml                   # Build and development tasks
├── services/                      # Microservices architecture
│   ├── exercises/                 # Exercises service
│   │   ├── Dockerfile            # Exercises service container
│   │   ├── docker-compose.yml    # (removed - now in root)
│   │   ├── main.go              # Exercises service entry point
│   │   ├── internal/
│   │   │   ├── db/              # Database layer
│   │   │   ├── handler/         # HTTP handlers
│   │   │   ├── models/          # Data models
│   │   │   └── router/          # Router configuration
│   │   └── wait-for-db.sh       # Database initialization script
│   ├── gateway/                  # API Gateway service
│   │   ├── Dockerfile           # Gateway service container
│   │   └── main.go              # Gateway service entry point
│   ├── docs/                     # API documentation service (planned)
│   └── logger/                   # Logger service (planned)
├── build/                         # Compiled binaries
├── data/                          # Data files and import scripts
│   ├── exercises.json             # Exercise database
│   └── import_script.py           # Data import utility
├── docs/                          # Generated Swagger documentation
├── internal/
│   ├── config/                    # Configuration and logging
│   ├── db/                        # Database layer
│   │   ├── migrate/               # Database migration files
│   │   └── queries/               # SQL queries and generated code
│   ├── handler/                   # HTTP handlers and routing
│   ├── models/                    # Data models and structures
│   └── service/                   # Business logic layer
└── logs/                          # Application logs
```

## 🐳 Docker Commands

```bash
# Start all services
docker compose up --build

# Start in background
docker compose up -d --build

# View logs
docker compose logs -f

# View specific service logs
docker compose logs -f exercises-service

# Stop all services
docker compose down

# Remove volumes (will delete database data)
docker compose down -v

# Rebuild specific service
docker compose build exercises-service
```

## 🛠️ Development

### Build Tasks

The project uses [Task](https://taskfile.dev/) for build automation. Available tasks:

```bash
task default        # Complete setup: database, migrations, and build
task build          # Build the application with docs generation
task setup_db       # Set up database with migrations and data import
task migrate        # Run database migrations
task import_data    # Import exercise data into database
task build-docs     # Generate Swagger documentation
```

### Database Schema

The application uses PostgreSQL with the following main entities:

- **exercises**: Core exercise information with equipment types
- **exercise_names**: Multiple names/aliases for exercises
- **muscles**: Muscle group definitions
- **exercise_muscle**: Many-to-many relationship between exercises and muscles
- **programs**: Workout program definitions with sets and reps
- **visuals**: Exercise demonstration media

### Equipment Types

Supported equipment categories:
- Dumbbells
- Barbell
- Machine
- Bodyweight
- Medicine Ball
- Kettlebells
- Stretches
- Cables
- Band
- Plate
- TRX
- Bosu Ball
- Foam Roll
- Exercise Ball
- Other

## 🔧 Configuration

Create an `.env` file in `internal/config/` with the following variables:

```env
DATABASE_URL=postgres://username:password@localhost:5432/guddy?sslmode=disable
PORT=3000
LOG_LEVEL=info
```

## 📊 Database Migrations

Database migrations are located in `internal/db/migrate/`. The current schema includes:

1. **000001_init_schema**: Initial database schema setup
2. **000002_add_exercises_instructions_field**: Added instructions field to exercises
3. **000003_alter_programs_table**: Modified programs table structure

## 🧪 Testing

Run tests with:
```bash
go test ./...
```

## 📝 API Examples

### Get All Exercises
```bash
curl -X GET "http://localhost:3000/api/exercises"
```

### Create a New Program
```bash
curl -X POST "http://localhost:3000/api/program" \
  -H "Content-Type: application/json" \
  -d '{
    "exercises": [
      {
        "exerciseId": 1,
        "idx": 1,
        "sets": 3,
        "reps": 10
      }
    ]
  }'
```

### Get Program by UUID
```bash
curl -X GET "http://localhost:3000/api/program/{uuid}"
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

All rights reserved. You may not use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software without prior written permission. Any unauthorized use is strictly prohibited and may be prosecuted under applicable law.
 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Exercise data sourced from community fitness databases
- Built with [Chi router](https://github.com/go-chi/chi) for HTTP routing
- Database operations powered by [pgx](https://github.com/jackc/pgx)
- Documentation generated with [Swag](https://github.com/swaggo/swag)

## 📞 Support

If you have any questions or need help, please open an issue on GitHub.

---

**Happy exercising with Guddy! 💪**