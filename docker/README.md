# Docker Environment for Agents Project

This directory contains the Docker Compose configuration for running the Agents project with all its dependencies.

## Services

### PostgreSQL Database
- **Port**: 6888 (mapped from container port 5432)
- **Database**: agents
- **User**: agents
- **Password**: agents (configurable via environment variables)

### ADK Web Frontend
- **Port**: 4200
- **Backend URL**: Configurable via `ADK_BACKEND_URL` environment variable
- **Default Backend**: http://host.docker.internal:6999

## Quick Start

1. **Start all services** (automatically creates .env from .env.example):
   ```bash
   mage docker:up
   ```

2. **View logs**:
   ```bash
   # All services
   mage docker:logs
   
   # Specific service
   mage docker:logs adk-web
   mage docker:logs postgres
   ```

3. **Stop services**:
   ```bash
   mage docker:down
   ```

4. **Get help**:
   ```bash
   mage docker:help
   ```

## Configuration

### Environment Variables

The `.env` file is automatically created from `.env.example` when you run `mage docker:up`. You can customize it as needed:

```bash
# Database Configuration
AGENTS_DB_USER=agents
AGENTS_DB_PASSWORD=your_secure_password
AGENTS_DB_NAME=agents

# ADK Web Frontend Configuration
ADK_BACKEND_URL=http://host.docker.internal:6999
NODE_ENV=development
```

### Mage Commands

This project uses [Mage](https://magefile.org/) for task automation. All Docker operations are managed through Mage targets:

```bash
mage docker:up          # Start all services
mage docker:down        # Stop all services
mage docker:restart     # Restart all services
mage docker:logs        # Show logs for all services
mage docker:logs <svc>  # Show logs for specific service
mage docker:build       # Build all services
mage docker:build <svc> # Build specific service
mage docker:status      # Show service status
mage docker:clean       # Remove all containers, networks, and volumes
mage docker:exec        # Execute command in service container
mage docker:pull        # Pull latest images
mage docker:config      # Validate and show configuration
mage docker:help        # Show detailed help
```

### Backend Integration

The ADK web frontend is configured to connect to your Go backend running on the host machine at port 6999. If your backend runs on a different port or host, update the `ADK_BACKEND_URL` environment variable.

For Docker Desktop users:
- Use `host.docker.internal:6999` to connect to services running on the host
- For Linux users, you might need to use `172.17.0.1:6999` or the actual host IP

## Service URLs

Once running, you can access:

- **ADK Web Frontend**: http://localhost:4200
- **PostgreSQL**: localhost:6888

## Troubleshooting

### ADK Web Service Issues

1. **Check if the service is healthy**:
   ```bash
   mage docker:status
   ```

2. **View detailed logs**:
   ```bash
   mage docker:logs adk-web
   ```

3. **Rebuild the service**:
   ```bash
   mage docker:build adk-web
   mage docker:up
   ```

### Database Connection Issues

1. **Check PostgreSQL logs**:
   ```bash
   mage docker:logs postgres
   ```

2. **Connect to database directly**:
   ```bash
   mage docker:exec postgres psql -U agents -d agents
   ```

### Network Issues

If the ADK web frontend can't connect to your backend:

1. **Check backend URL configuration**:
   ```bash
   mage docker:exec adk-web env | grep BACKEND_URL
   ```

2. **Test connectivity from container**:
   ```bash
   mage docker:exec adk-web wget -qO- http://host.docker.internal:6999/health
   ```

## Development

### Rebuilding Services

To rebuild services after changes:

```bash
# Rebuild specific service
mage docker:build adk-web

# Rebuild all services
mage docker:build

# Force rebuild without cache (use clean)
mage docker:clean
```

### Updating ADK Web

The ADK web frontend is cloned from the official repository during the Docker build. To update to the latest version:

```bash
mage docker:clean  # This rebuilds without cache
```

### Common Mage Commands

```bash
# Quick development cycle
mage docker:up                    # Start services
mage docker:logs adk-web          # Monitor ADK web logs
mage docker:restart               # Restart after changes

# Debugging
mage docker:status                # Check service health
mage docker:config                # Validate configuration
mage docker:exec adk-web sh       # Shell into ADK web container

# Maintenance
mage docker:pull                  # Update base images
mage docker:clean                 # Full cleanup and rebuild
```