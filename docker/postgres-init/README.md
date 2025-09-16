# PostgreSQL Database Initialization

This directory contains initialization scripts for the PostgreSQL database used by the Agents project.

## Overview

The PostgreSQL container is configured to automatically run initialization scripts when a new database cluster is created. This happens when the `postgres_data` volume is empty (first startup or after volume deletion).

## Database Structure

The initialization creates the following databases:

### 1. `agents` Database
- **Purpose**: Used by the agents service for storing agent data, tasks, projects, etc.
- **Owner**: `agents` user
- **Connection**: `postgres://agents:agents@postgres:5432/agents?sslmode=disable`

### 2. `convex_self_hosted` Database  
- **Purpose**: Used by the Convex self-hosted backend for internal data storage
- **Owner**: `agents` user
- **Connection**: `postgres://agents:agents@postgres:5432?sslmode=disable`
- **Note**: Convex manages database selection internally, so no database name in URL

## Initialization Process

### Automatic Initialization (Recommended)
When you start the services for the first time:
```bash
cd docker
docker-compose up -d
```

The PostgreSQL container will automatically:
1. Create the database cluster
2. Run all `.sql` files in `/docker-entrypoint-initdb.d/` (alphabetically)
3. Create the required databases and users
4. Set up proper permissions

### Manual Database Management (Using Mage)

**Initialize databases (recommended):**
```bash
mage docker:initdb
```

**Reset databases completely (WARNING: Deletes all data):**
```bash
mage docker:resetdb
```

### Manual Re-initialization (Low-level)
If you need to reset manually without Mage:

```bash
# Stop all services
docker-compose down

# Remove the postgres data volume (WARNING: This deletes all data!)
docker volume rm docker_postgres_data

# Start services again (will trigger re-initialization)
docker-compose up -d
```

### Verify Initialization
Check that databases were created correctly:
```bash
# List all databases
docker-compose exec postgres psql -U agents -l

# Connect to agents database
docker-compose exec postgres psql -U agents -d agents

# Connect to convex database
docker-compose exec postgres psql -U agents -d convex_self_hosted
```

## Files

- `01-init-databases.sql`: Main initialization script that creates databases and sets permissions
- `README.md`: This documentation file

## Adding Custom Initialization

To add custom initialization logic:

1. Create new `.sql` files in this directory
2. Use numerical prefixes to control execution order (e.g., `02-custom-setup.sql`)
3. Restart with fresh volume to trigger re-initialization

## Environment Variables

The initialization uses these environment variables from `.env.docker`:
- `POSTGRES_USER=agents`: Database superuser
- `POSTGRES_PASSWORD=agents`: Password for the superuser
- `POSTGRES_DB=agents`: Default database (not used since we create specific ones)

## Troubleshooting

### Initialization Not Running
- Check that `postgres_data` volume is empty
- Initialization only runs on first container startup with empty data directory

### Permission Issues
- Ensure the `agents` user has proper privileges
- Check the initialization logs: `docker-compose logs postgres`

### Connection Issues
- Verify the databases exist: `docker-compose exec postgres psql -U agents -l`
- Check connection strings match the format documented above
- Ensure healthcheck passes: `docker-compose ps postgres`