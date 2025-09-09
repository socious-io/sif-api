# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Quick Start
```bash
# Setup dependencies and database
cp .tmp.config.yml config.yml
sudo docker-compose up -d
go get
go run cmd/migrate/main.go up
go run cmd/app/main.go
```

### Running the Application
- **Start server**: `go run cmd/app/main.go` (runs on port 3000 by default)
- **Build**: `go build -C cmd/app -o ../../build`

### Database Operations
- **Run migrations**: `go run cmd/migrate/main.go up`
- **Create migration**: `go run cmd/migrate/main.go new {migration_name}`
- **Force migration version**: `go run cmd/migrate/main.go force {version}`

### Testing
- **Run tests**: `go test ./tests/... -c test.config.yml`
- **Run with focus**: `go test ./tests/... -focus "specific test" -c test.config.yml`
- Tests use Ginkgo/Gomega framework
- Test config: `test.config.yml` (copy from `.tmp.config.yml`)

## Architecture Overview

### Project Structure
- **`cmd/`**: Application entry points
  - `cmd/app/main.go`: Main API server
  - `cmd/migrate/main.go`: Database migration tool
- **`src/apps/`**: Core application logic
  - `models/`: Database models and business logic
  - `views/`: HTTP handlers and routing
  - `auth/`: Authentication and JWT middleware
  - `utils/`: Shared utilities (GCS upload, security, etc.)
- **`src/sql/`**: SQL queries organized by entity
  - `migrations/`: Database schema migrations
  - Entity folders: `users/`, `projects/`, `donations/`, etc.
- **`tests/`**: Integration tests using Ginkgo

### Key Technologies
- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL with sqlx
- **Authentication**: JWT tokens via `goaccount` library
- **Payments**: Integration with Stripe and blockchain via `gopay` library
- **Storage**: Google Cloud Storage for file uploads
- **Testing**: Ginkgo/Gomega

### Configuration
- Main config: `config.yml` (copy from `.tmp.config.yml`)
- Environment variables can override config values
- Database URL format: `postgresql://user:pass@host:port/db?sslmode=disable`
- Supports multiple payment providers (Stripe, blockchain)

### Database Query System
- SQL files organized in `src/sql/{entity}/` directories
- Each entity has CRUD operations as separate `.sql` files
- Migrations in `src/sql/migrations/` with timestamp prefixes
- Schema defined in `src/sql/schema.sql`

### Security Features
- XSS protection via request sanitization
- CORS configuration for allowed origins
- JWT-based authentication
- Input validation and secure headers

### Key Models
- **Users**: Identity management with KYC/KYB verification
- **Projects**: Crowdfunding projects with quadratic funding
- **Donations**: Payment tracking with blockchain support  
- **Organizations**: Multi-user organizations
- **Rounds**: Funding rounds for projects
- **Votes**: Community voting on projects

## Development Notes
- Uses Go modules (`go.mod`) for dependency management
- Docker Compose provides PostgreSQL database (port 5432)
- Application runs on port 3000 by default
- Debug mode enables permissive CORS for development
- File uploads handled via Google Cloud Storage with CDN