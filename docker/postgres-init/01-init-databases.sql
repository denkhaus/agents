-- PostgreSQL initialization script for Agents project
-- This script creates the necessary databases for the application
-- Create database for Convex self-hosted backend
CREATE DATABASE convex_self_hosted;

-- Grant all privileges to the agents user (already created via environment variables)
GRANT ALL PRIVILEGES ON DATABASE agents TO agents;

GRANT ALL PRIVILEGES ON DATABASE convex_self_hosted TO agents;

-- Optional: Create additional users or configure specific permissions here
-- Example:
-- CREATE USER readonly_user WITH PASSWORD 'readonly_password';
-- GRANT CONNECT ON DATABASE agents TO readonly_user;
-- GRANT USAGE ON SCHEMA public TO readonly_user;
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;

-- Log successful initialization
\echo 'Database initialization completed successfully'
\echo 'Created databases: convex_self_hosted'
\echo 'Granted privileges to user: agents'
