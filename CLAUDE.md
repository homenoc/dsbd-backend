# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HomeNOC WebSystem Backend - A Go backend service for managing network operations, user groups, and connectivity services. The system manages user registrations, group organizations, network service requests, and integrates with external services (Stripe payments, Slack notifications, JPNIC for IP address registration).

## Build & Run Commands

### Tool Management (aqua)
```bash
# Install aqua (if not installed)
brew install aquaproj/aqua/aqua

# Install development tools (Go, Air, golangci-lint)
aqua i
```

### Local Development (Recommended)
```bash
# Initial setup
make setup              # Creates tmp/ and copies config

# Database (Docker)
make db-up              # Start MySQL container
make migrate            # Run database migrations
make seed               # Create test accounts

# Run servers
make run-user           # User API on :8080
make run-admin          # Admin API on :8081

# Stop
make db-down
```

### Docker Compose (Full Stack)
```bash
make setup-docker       # Setup with Docker config
make docker-up          # Start all containers
make docker-migrate     # Run migrations in container
make docker-seed        # Create test data in container
```

### Testing & Linting
```bash
make test               # Run all tests (go test ./...)
make lint               # Run golangci-lint

# Run single test
go test -v -run TestFunctionName ./path/to/package/...
```

### Test Accounts (created by `make seed`)
- **Master**: master@example.com / password (Level 1: full access)
- **Member**: member@example.com / password (Level 2: general member)
- **Admin API**: Uses Basic auth from config (default: admin/admin)

## Architecture

### Two Separate API Servers
- **User API** (`:8080`): User-facing endpoints for registration, login, service requests, support tickets
- **Admin API** (`:8081`): Admin endpoints for managing users, groups, services, NOC infrastructure

Both APIs use Gin framework and follow the pattern `/api/v1/*` for REST endpoints and `/ws/v1/*` for WebSocket.

### Core Package Structure (`pkg/api/`)
```
pkg/api/
├── api.go              # Router setup, defines all User/Admin API routes
├── core/               # Business logic and domain models
│   ├── interface.go    # All GORM data models (User, Group, Service, Connection, etc.)
│   ├── user/v0/        # User operations
│   ├── group/v0/       # Group and organization management
│   ├── noc/v0/         # Network Operations Center entities
│   ├── support/        # Support ticket system with WebSocket chat
│   ├── payment/v0/     # Stripe integration
│   └── tool/           # Utilities (config, token, hash, slack, mail)
└── store/              # Data access layer (GORM/MySQL operations)
```

### Key Data Models (`pkg/api/core/interface.go`)
- **User**: Individual users with authentication tokens, group membership
- **Group**: Organizations containing users and services, Stripe billing
- **Service**: Network service types (L2, L3 Static, L3 BGP, Transit)
- **Connection**: Physical/tunnel connections (EtherIP, GRE, IP-IP, Cross Connect)
- **NOC/BGPRouter/TunnelEndPointRouter**: Network infrastructure
- **Ticket/Chat**: Support ticket system with real-time messaging

### Handler Pattern
- Handlers in `pkg/api/core/{feature}/v0/` handle HTTP requests
- Naming convention: `*ByAdmin` suffix for admin-only handlers
- Store layer in `pkg/api/store/{feature}/` for database operations
- Routes defined in `pkg/api/api.go` (AdminRestAPI and UserRestAPI functions)

### Configuration
Single JSON config file (`configs/config.json`) loaded via Viper. Key sections:
- `controller.user/admin`: API server settings and ports
- `db`: MySQL connection
- `stripe`: Payment integration
- `slack`: Notification channels
- `mail`: SMTP settings
- `jpnic`: IP address registration certificates
- `template`: Service types (L2, L3 Static, L3 BGP, Transit) and connection types (EtherIP, GRE, IP-IP, Cross Connect)

## Database

- **ORM**: GORM with MySQL driver
- **Auto-migration**: Handled in `store.InitDB()` in `pkg/api/store/store.go`
- **Note**: MySQL `sql_mode=''` may be needed for long text fields
- Database diagram: https://drawsql.app/y-net/diagrams/dsbd-backend/embed

## Testing

Tests use `go-sqlmock` for database mocking. Pattern for database tests:
```go
mock, cleanup := setupTestDB(t)
defer cleanup()
// Use store.SetTestDB() / store.ClearTestDB() for test DB injection
```
