# Hackathon Backend API

A clean, scalable Go HTTP API skeleton using Echo framework with PostgreSQL connection via pgxpool.

## Architecture

This project follows clean architecture principles with clear separation of concerns:

```
├── config/          # Configuration management
├── database/        # Database connection and utilities
├── handlers/        # HTTP handlers (presentation layer)
├── models/          # Data models and DTOs
├── repositories/    # Data access layer
├── routes/          # Route definitions
├── services/        # Business logic layer
└── scripts/         # Database initialization scripts
```

## Features

- ✅ Echo HTTP framework
- ✅ PostgreSQL with pgxpool connection
- ✅ Clean architecture (handlers → services → repositories)
- ✅ Environment configuration
- ✅ Docker & Docker Compose setup
- ✅ Graceful shutdown
- ✅ CORS middleware
- ✅ Request validation
- ✅ Error handling
- ✅ Health check endpoint

## Quick Start

### Using Docker (Recommended)

1. **Clone and setup:**
   ```bash
   git clone <your-repo>
   cd backend
   ```

2. **Start services:**
   ```bash
   make docker-up
   ```

3. **API will be available at:** `http://localhost:8080`

### Manual Setup

1. **Install dependencies:**
   ```bash
   make deps
   ```

2. **Setup environment:**
   ```bash
   cp env.example .env
   # Edit .env with your database credentials
   ```

3. **Start PostgreSQL:**
   ```bash
   docker-compose up -d postgres
   ```

4. **Run the application:**
   ```bash
   make run
   ```

## API Endpoints

### Health Check
- `GET /health` - Server health status

### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users (with pagination)
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## Example Usage

### Create a user:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User"
  }'
```

### Get all users:
```bash
curl http://localhost:8080/api/v1/users
```

### Get user by ID:
```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

## Development

### Available Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make clean          # Clean build artifacts
make docker-up      # Start services with Docker
make docker-down    # Stop services
make docker-build   # Build and start services
make dev-setup      # Setup development environment
```

### Project Structure

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and orchestration
- **Repositories**: Data access and database operations
- **Models**: Data structures and validation
- **Config**: Environment and configuration management
- **Database**: Connection pooling and database utilities

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | password | Database password |
| `DB_NAME` | hackathon_db | Database name |
| `DB_SSL_MODE` | disable | SSL mode |
| `SERVER_PORT` | 8080 | Server port |
| `SERVER_HOST` | localhost | Server host |
| `ENV` | development | Environment |

## Database Schema

The application includes a sample `users` table with the following structure:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Adding New Features

1. **Create model** in `models/` directory
2. **Create repository** in `repositories/` directory
3. **Create service** in `services/` directory
4. **Create handler** in `handlers/` directory
5. **Add routes** in `routes/routes.go`

## Production Considerations

- Add proper logging (structured logging)
- Implement authentication/authorization
- Add rate limiting
- Implement caching
- Add monitoring and metrics
- Set up proper error tracking
- Add API documentation (Swagger)
- Implement database migrations
- Add comprehensive testing

## License

MIT License - feel free to use this skeleton for your hackathon project!
