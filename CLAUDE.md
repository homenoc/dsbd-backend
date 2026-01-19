# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HomeNOC WebSystem Backend - A Go backend service for managing network operations, user groups, and connectivity services. The system manages user registrations, group organizations, network service requests, and integrates with external services (Stripe payments, Slack notifications, JPNIC for IP address registration).

## Build & Run Commands

```bash
# Build
go build -o tmp/main cmd/backend/main.go

# Initialize database (required on first run)
go run cmd/backend/main.go init database --config config.json

# Start User API (port 8080)
go run cmd/backend/main.go start user --config config.json

# Start Admin API (port 8081)
go run cmd/backend/main.go start admin --config config.json

# Run tests
go test ./...

# Local development with Docker Compose (includes MySQL + hot reload via Air)
mkdir tmp
cp configs/config.json tmp/config.json
docker compose up -d
```

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

### Configuration
Single JSON config file (`configs/config.json`) loaded via Viper. Key sections:
- `controller.user/admin`: API server settings and ports
- `db`: MySQL connection
- `stripe`: Payment integration
- `slack`: Notification channels
- `mail`: SMTP settings
- `jpnic`: IP address registration certificates

## Database

- **ORM**: GORM with MySQL driver
- **Auto-migration**: Handled in `InitDB()`
- **Note**: MySQL `sql_mode=''` may be needed for long text fields
- Database diagram: https://drawsql.app/y-net/diagrams/dsbd-backend/embed
