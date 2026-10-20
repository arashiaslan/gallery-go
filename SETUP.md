# Project Setup Guide

## Prerequisites

### Required
- **Go 1.22+** - [Download](https://golang.org/dl/)
- **Git** - For version control

### Optional
- **Docker & Docker Compose** - For containerized development
- **curl** or **Postman** - For testing API
- **VS Code** with Go extension - For IDE

## Quick Start (Local Development)

### 1. Clone or navigate to project
```bash
cd gallery
```

### 2. Create .env file (optional, for non-defaults)
```bash
cp backend/.env.example .env
# Edit .env as needed
```

### 3. Run development server
```bash
cd backend
make dev
```

Server will start on `http://localhost:8080`

### 4. Test login (with default admin credentials)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

You'll get a session token in the response and a cookie set.

## Running with Docker

### 1. Build and run
```bash
docker-compose up
```

This will:
- Build the Go binary
- Create the database
- Start the server
- Set up health checks

### 2. Test
```bash
curl http://localhost:8080/health
```

### 3. Stop
```bash
docker-compose down
```

### 4. Clean up
```bash
docker-compose down -v  # Also remove volumes
```

## Backend Development Workflow

### Build commands
```bash
cd backend

# Run development server with auto-reload (requires watchmodo or similar)
make dev

# Build optimized binary
make build
# Output: ./bin/server

# Run tests
make test

# Clean build artifacts
make clean

# Download dependencies
make deps
```

### Database Management

**Migrations run automatically** on startup. To run manually:
```bash
cd backend
go run -tags=tools cmd/migrate/main.go up
```

To reset database:
```bash
rm data/gallery.db
make dev  # Recreates database
```

### Testing

Run all tests:
```bash
cd backend
make test
```

Run specific test:
```bash
go test -v ./internal/handlers -run TestLogin
```

Run with coverage:
```bash
go test -cover ./...
```

## Environment Variables

Copy `.env.example` to `.env` and customize:

```bash
# Server
PORT=8080
DATA_DIR=./data

# Database
DATABASE_URL=./data/gallery.db

# Auth
SESSION_SECRET=change-this-in-production
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123

# Upload
MAX_UPLOAD_MB=25
```

## Directory Structure

```
gallery/
├── backend/
│   ├── cmd/server/           # Server binary
│   ├── internal/
│   │   ├── auth/             # Authentication
│   │   ├── config/           # Configuration
│   │   ├── handlers/         # HTTP handlers
│   │   ├── store/            # Database layer
│   │   ├── media/            # Image processing
│   │   └── utils/            # Utilities
│   ├── migrations/           # SQL migrations
│   ├── Makefile              # Build targets
│   ├── go.mod / go.sum       # Dependencies
│   └── .env.example          # Config template
├── Dockerfile                # Container definition
├── docker-compose.yml        # Container orchestration
├── README.md                 # Project documentation
├── API.md                    # API reference
├── SPEC.md                   # Project specification
└── M1_COMPLETION.md          # M1 delivery report
```

## IDE Setup (VS Code)

### Recommended Extensions
- Go (golang.go) - Official Go extension
- REST Client - For testing API
- SQLite (alexcvzz.vscode-sqlite) - Database inspection

### VS Code Settings
Create `.vscode/settings.json`:
```json
{
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.useLanguageServer": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### Launch Configuration
Create `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/server",
      "cwd": "${workspaceFolder}/backend"
    }
  ]
}
```

## Troubleshooting

### "command not found: go"
Install Go from https://golang.org/dl/

### "database is locked"
SQLite with multiple writers can cause this. Solution:
- Ensure only one server instance is running
- Delete `data/gallery.db*` and restart (recreates database)

### Port already in use
Change PORT in .env:
```bash
PORT=3000 make dev
```

Or kill process:
```bash
lsof -i :8080  # Find process
kill -9 <PID>  # Kill it
```

### Tests failing
Ensure you're in backend directory:
```bash
cd backend
make test
```

If in-memory database test fails, try:
```bash
go test -v -count=1 ./...
```

### Docker build fails
Try:
```bash
docker-compose build --no-cache
docker-compose up
```

## Next Steps

1. **Review SPEC.md** - Understand full project scope
2. **Read API.md** - Understand current API
3. **Start M2** - Photo upload and image processing
4. **Create frontend** - React + Vite scaffolding

## Useful Commands

```bash
# Server health
curl http://localhost:8080/health

# Check current user
curl -b cookies.txt http://localhost:8080/api/me

# Format code
go fmt ./...

# Vet code
go vet ./...

# Build for production
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server

# View database
sqlite3 data/gallery.db ".schema"

# Delete specific session
sqlite3 data/gallery.db "DELETE FROM sessions WHERE user_id='xyz';"
```

## Resources

- [Go Documentation](https://golang.org/doc/)
- [chi Router Docs](https://github.com/go-chi/chi)
- [SQLite Docs](https://www.sqlite.org/docs.html)
- [REST API Best Practices](https://restfulapi.net/)
