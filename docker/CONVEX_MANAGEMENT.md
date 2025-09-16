# Convex Backend Management

This guide explains how to manage the Convex self-hosted backend using Mage targets.

## Quick Start

### Generate Admin Key
```bash
# Generate an admin key for API authentication
mage convex:generateadminkey
```

### Check Backend Status
```bash
# Verify backend is running and healthy
mage convex:status
```

### View Instance Credentials
```bash
# Show instance name and partial secret
mage convex:showcredentials
```

## Available Commands

### Core Management
```bash
# Generate admin key for API access
mage convex:generateadminkey

# Show current instance credentials
mage convex:showcredentials

# Check backend health and status
mage convex:status

# View backend logs in real-time
mage convex:logs

# Restart backend service
mage convex:restart
```

### Data Management
```bash
# Reset all backend data (WARNING: destructive)
mage convex:reset

# Get help and examples
mage convex:help
```

## Admin Key Usage

The admin key is used to authenticate with the Convex backend API:

```bash
# Example admin key format
convex-self-hosted|01f1a0a85af28c4d08a2da7a0449b1e961b1cdc3760dab8b26c7201e415586aa33f04805c2

# Use in HTTP requests
curl -H "Authorization: Bearer YOUR_ADMIN_KEY" http://localhost:3210/api/endpoint
```

### Key Properties
- **Persistence**: Admin keys persist across backend restarts
- **Regeneration**: Generate new keys anytime with `mage convex:generateadminkey`
- **Instance-tied**: Keys are tied to the instance name and secret
- **Security**: Store securely, never share publicly

## Service URLs

### Backend API
- **URL**: http://localhost:3210
- **Purpose**: Main Convex backend API
- **Authentication**: Admin key required for most endpoints

### Dashboard
- **URL**: http://localhost:6791
- **Purpose**: Convex development dashboard
- **Authentication**: Uses backend API internally

## Instance Management

### Instance Credentials
The backend uses persistent instance credentials:
- **Instance Name**: `convex-self-hosted` (default)
- **Instance Secret**: Randomly generated 64-character hex string
- **Storage**: Saved in `/convex/data/credentials/` inside container

### Viewing Credentials
```bash
# Show instance name and partial secret (secure)
mage convex:showcredentials

# Full credentials (inside container)
mage docker:exec backend cat /convex/data/credentials/instance_name
mage docker:exec backend cat /convex/data/credentials/instance_secret
```

## Database Integration

### PostgreSQL Connection
The Convex backend connects to PostgreSQL:
- **Database**: `convex_self_hosted`
- **Connection**: Automatic, configured via `POSTGRES_URL`
- **Fallback**: SQLite if PostgreSQL unavailable

### Database Management
```bash
# Initialize PostgreSQL databases
mage docker:initdb

# Reset PostgreSQL (affects Convex data)
mage docker:resetdb

# Check database connection status
mage convex:status
```

## Troubleshooting

### Backend Not Running
```bash
# Start all services
mage docker:up

# Check service status
mage docker:status

# View backend logs
mage convex:logs
```

### Admin Key Generation Fails
```bash
# Ensure backend is healthy
mage convex:status

# Check backend logs for errors
mage convex:logs

# Restart backend if needed
mage convex:restart
```

### Database Connection Issues
```bash
# Check if PostgreSQL is healthy
mage docker:status

# Initialize databases if needed
mage docker:initdb

# View backend connection logs
mage convex:logs | grep -i postgres
```

### Reset Everything
```bash
# Reset Convex data only
mage convex:reset

# Reset databases only
mage docker:resetdb

# Reset everything
mage docker:clean
mage docker:up
```

## Development Workflow

### Daily Development
1. **Start services**: `mage docker:up`
2. **Get admin key**: `mage convex:generateadminkey`
3. **Check status**: `mage convex:status`
4. **View logs**: `mage convex:logs` (as needed)

### Fresh Development Setup
1. **Clean slate**: `mage docker:clean`
2. **Start fresh**: `mage docker:up`
3. **Initialize DBs**: `mage docker:initdb`
4. **Generate key**: `mage convex:generateadminkey`

### Debugging Issues
1. **Check status**: `mage convex:status`
2. **View logs**: `mage convex:logs`
3. **Restart backend**: `mage convex:restart`
4. **Reset if needed**: `mage convex:reset`

## Integration Examples

### Application Configuration
```javascript
// Example Node.js configuration
const CONVEX_BACKEND_URL = 'http://localhost:3210';
const CONVEX_ADMIN_KEY = 'convex-self-hosted|your-admin-key-here';

// Use in API calls
const response = await fetch(`${CONVEX_BACKEND_URL}/api/endpoint`, {
  headers: {
    'Authorization': `Bearer ${CONVEX_ADMIN_KEY}`,
    'Content-Type': 'application/json'
  }
});
```

### Environment Variables
```bash
# .env.docker file for your application
CONVEX_BACKEND_URL=http://localhost:3210
CONVEX_ADMIN_KEY=convex-self-hosted|your-admin-key-here
```

## Production Considerations

- **Security**: Use strong instance secrets in production
- **Networking**: Configure proper network access controls
- **Monitoring**: Set up monitoring for backend health
- **Backup**: Implement backup strategies for both PostgreSQL and Convex data
- **Scaling**: Consider load balancing and high availability setup

## Summary

The Convex management namespace provides:
- ✅ **Easy admin key generation** for API authentication
- ✅ **Health monitoring** and status checking
- ✅ **Credential management** with secure display
- ✅ **Data reset capabilities** for development
- ✅ **Comprehensive logging** and debugging tools
- ✅ **Integration with PostgreSQL** database setup

Use `mage convex:help` for quick reference of all available commands.