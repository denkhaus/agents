# Database Setup Guide

This guide explains how to properly set up and manage PostgreSQL databases for the Agents project.

## Quick Start

### 1. First-time Setup
```bash
# Start all services (includes automatic database initialization)
mage docker:up

# Or initialize databases explicitly
mage docker:initdb
```

### 2. Verify Setup
```bash
# Check service status
mage docker:status

# View database list
mage docker:exec postgres psql -U agents -l
```

## Database Architecture

The project uses **two separate PostgreSQL databases**:

### 1. `agents` Database
- **Purpose**: Stores agent data, projects, tasks, dependencies
- **Used by**: Agents service (Go application)
- **Connection**: `postgres://agents:agents@postgres:5432/agents?sslmode=disable`
- **External access**: `postgres://agents:agents@localhost:6888/agents?sslmode=disable`

### 2. `convex_self_hosted` Database  
- **Purpose**: Convex backend internal storage (real-time data, functions, etc.)
- **Used by**: Convex self-hosted backend
- **Connection**: `postgres://agents:agents@postgres:5432?sslmode=disable`
- **External access**: `postgres://agents:agents@localhost:6888?sslmode=disable`
- **Note**: Convex automatically selects the `convex_self_hosted` database

## Available Commands

### Database Management (Mage Targets)
```bash
# Initialize/create required databases
mage docker:initdb

# Reset all databases (WARNING: Deletes all data)
mage docker:resetdb

# Start all services
mage docker:up

# Stop all services
mage docker:down

# View database logs
mage docker:logs postgres

# Execute commands in postgres container
mage docker:exec postgres <command>
```

### Direct Database Access
```bash
# Connect to agents database
mage docker:exec postgres psql -U agents -d agents

# Connect to convex database
mage docker:exec postgres psql -U agents -d convex_self_hosted

# List all databases
mage docker:exec postgres psql -U agents -l

# Check database sizes
mage docker:exec postgres psql -U agents -c "SELECT datname, pg_size_pretty(pg_database_size(datname)) FROM pg_database;"
```

## Configuration Files

### Environment Variables (`.env.docker`)
```bash
# PostgreSQL Configuration
POSTGRES_USER=agents
POSTGRES_PASSWORD=agents
POSTGRES_DB=agents  # Default database (not used directly)

# Database URLs
AGENTS_DATABASE_URL="postgres://agents:agents@postgres:5432/agents?sslmode=disable"
POSTGRES_URL="postgres://agents:agents@postgres:5432?sslmode=disable"
DISABLE_BEACON="true"
```

### Docker Compose (`docker/docker-compose.yml`)
- PostgreSQL service with healthcheck
- Automatic initialization via mounted scripts
- Volume mounting for persistence
- Network configuration for service communication

### Initialization Scripts (`docker/postgres-init/`)
- `01-init-databases.sql`: Creates both databases and sets permissions
- `README.md`: Detailed initialization documentation

## Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running and healthy
mage docker:status

# View PostgreSQL logs
mage docker:logs postgres

# Test connection
mage docker:exec postgres pg_isready -U agents
```

### Convex Backend Issues
```bash
# Check if convex_self_hosted database exists
mage docker:exec postgres psql -U agents -l | grep convex

# View Convex backend logs
mage docker:logs backend

# Verify Convex is using PostgreSQL (not SQLite)
mage docker:logs backend | grep -E "(Connected to|db_connection)"
```

### Reset Everything
```bash
# Complete reset (removes all data)
mage docker:resetdb

# Or manual cleanup
mage docker:down
docker volume rm docker_postgres_data
mage docker:up
```

## Development Workflow

### Making Schema Changes
1. **For agents database**: Modify Go models and run migrations
2. **For convex database**: Use Convex schema management tools
3. **For both**: Consider using `mage docker:resetdb` for fresh development setup

### Backup and Restore
```bash
# Backup agents database
mage docker:exec postgres pg_dump -U agents agents > backup_agents.sql

# Backup convex database  
mage docker:exec postgres pg_dump -U agents convex_self_hosted > backup_convex.sql

# Restore (after reset)
mage docker:exec postgres psql -U agents agents < backup_agents.sql
mage docker:exec postgres psql -U agents convex_self_hosted < backup_convex.sql
```

### Monitoring
```bash
# Monitor database activity
mage docker:exec postgres psql -U agents -c "SELECT * FROM pg_stat_activity;"

# Check database sizes
mage docker:exec postgres psql -U agents -c "SELECT datname, pg_size_pretty(pg_database_size(datname)) FROM pg_database WHERE datname IN ('agents', 'convex_self_hosted');"
```

## Production Considerations

- Use strong passwords in production environment
- Configure proper backup strategies
- Monitor database performance and storage
- Consider read replicas for scaling
- Implement proper SSL/TLS for external connections
- Use connection pooling for high-traffic applications

## Summary

The database setup provides:
- ✅ **Automated initialization** via Mage targets
- ✅ **Proper separation** between agents and Convex data
- ✅ **Development-friendly** reset and management commands
- ✅ **Production-ready** configuration with health checks
- ✅ **Comprehensive documentation** and troubleshooting guides

Use `mage docker:help` to see all available commands.