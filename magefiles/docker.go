package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Docker contains targets for managing Docker services
type Docker mg.Namespace

const (
	dockerDir = "./docker"
)

// Up starts all Docker services
func (Docker) Up() error {
	printStatus("Starting Agents Docker environment...")

	if err := ensureEnvFile(); err != nil {
		return err
	}

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "up", "-d"); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	printSuccess("Services started successfully!")
	fmt.Println()
	printStatus("Service URLs:")
	fmt.Println("  - Agents Backend: http://localhost:9451")
	fmt.Println("  - Convex Backend: http://localhost:3210")
	fmt.Println("  - Convex Dashboard: http://localhost:6791")
	fmt.Println("  - PostgreSQL Database: localhost:6888")
	fmt.Println()
	printStatus("Use 'mage docker:logs' to view logs")
	printStatus("Use 'mage docker:down' to stop services")

	return nil
}

// Down stops all Docker services
func (Docker) Down() error {
	printStatus("Stopping Agents Docker environment...")

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "down"); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	printSuccess("Services stopped successfully!")
	return nil
}

// Restart restarts all Docker services
func (Docker) Restart() error {
	printStatus("Restarting Agents Docker environment...")

	mg.Deps(Docker.Down)
	mg.Deps(Docker.Up)

	printSuccess("Services restarted successfully!")
	return nil
}

// Logs shows logs for all services or a specific service
// Usage: mage docker:logs [service]
func (Docker) Logs() error {
	args := os.Args
	var service string

	// Check if a service name was provided
	for i, arg := range args {
		if strings.Contains(arg, "docker:logs") && i+1 < len(args) {
			service = args[i+1]
			break
		}
	}

	cmd := []string{"docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "logs", "-f"}

	if service != "" {
		printStatus(fmt.Sprintf("Showing logs for service: %s", service))
		cmd = append(cmd, service)
	} else {
		printStatus("Showing logs for all services...")
	}

	return sh.RunV(cmd[0], cmd[1:]...)
}

// Build builds all services or a specific service
// Usage: mage docker:build [service]
func (Docker) Build() error {
	args := os.Args
	var service string

	// Check if a service name was provided
	for i, arg := range args {
		if strings.Contains(arg, "docker:build") && i+1 < len(args) {
			service = args[i+1]
			break
		}
	}

	cmd := []string{"docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "build"}

	if service != "" {
		printStatus(fmt.Sprintf("Building service: %s", service))
		cmd = append(cmd, service)
	} else {
		printStatus("Building all services...")
	}

	if err := sh.RunV(cmd[0], cmd[1:]...); err != nil {
		return fmt.Errorf("failed to build services: %w", err)
	}

	printSuccess("Build completed!")
	return nil
}

// Status shows the status of all services
func (Docker) Status() error {
	printStatus("Service status:")
	return sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "ps")
}

// Clean removes all containers, networks, and volumes (with confirmation)
func (Docker) Clean() error {
	printWarning("This will remove all containers, networks, and volumes!")

	// In a real scenario, you might want to add interactive confirmation
	// For now, we'll just print the warning and proceed
	fmt.Print("Are you sure? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		printStatus("Cleanup cancelled.")
		return nil
	}

	printStatus("Cleaning up Docker environment...")

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "down", "-v", "--remove-orphans"); err != nil {
		return fmt.Errorf("failed to clean up: %w", err)
	}

	printSuccess("Cleanup completed!")
	return nil
}

// Exec executes a command in a running service container
// Usage: mage docker:exec service command
func (Docker) Exec() error {
	args := os.Args

	if len(args) < 4 {
		return fmt.Errorf("usage: mage docker:exec <service> <command>")
	}

	// Find the service and command arguments
	var service, command string
	var cmdArgs []string

	for i, arg := range args {
		if strings.Contains(arg, "docker:exec") {
			if i+1 < len(args) {
				service = args[i+1]
			}
			if i+2 < len(args) {
				command = args[i+2]
			}
			if i+3 < len(args) {
				cmdArgs = args[i+3:]
			}
			break
		}
	}

	if service == "" || command == "" {
		return fmt.Errorf("usage: mage docker:exec <service> <command>")
	}

	cmd := []string{"docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", service, command}
	cmd = append(cmd, cmdArgs...)

	printStatus(fmt.Sprintf("Executing '%s' in service '%s'", command, service))
	return sh.RunV(cmd[0], cmd[1:]...)
}

// Pull pulls the latest images for all services
func (Docker) Pull() error {
	printStatus("Pulling latest images...")

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "pull"); err != nil {
		return fmt.Errorf("failed to pull images: %w", err)
	}

	printSuccess("Images pulled successfully!")
	return nil
}

// Config validates and shows the Docker Compose configuration
func (Docker) Config() error {
	printStatus("Docker Compose configuration:")
	return sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "config")
}

// Core starts only the core services (agents and postgres)
func (Docker) Core() error {
	printStatus("Starting core services (agents and postgres)...")

	if err := ensureEnvFile(); err != nil {
		return err
	}

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "up", "-d", "agents", "postgres"); err != nil {
		return fmt.Errorf("failed to start core services: %w", err)
	}

	printSuccess("Core services started successfully!")
	fmt.Println()
	printStatus("Service URLs:")
	fmt.Println("  - Agents Backend: http://localhost:9451")
	fmt.Println("  - PostgreSQL Database: localhost:6888")
	fmt.Println()
	printStatus("Use 'mage docker:logs agents' to view agents logs")
	fmt.Println("Use 'mage docker:down' to stop services")

	return nil
}

// InitDB initializes or reinitializes the PostgreSQL databases
func (Docker) InitDB() error {
	printStatus("Initializing PostgreSQL databases...")
	
	// Check if postgres service is running
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "ps", "postgres"); err != nil {
		printWarning("PostgreSQL service not running. Starting it first...")
		if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "up", "-d", "postgres"); err != nil {
			return fmt.Errorf("failed to start PostgreSQL service: %w", err)
		}
		printStatus("Waiting for PostgreSQL to be ready...")
		if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "postgres", "pg_isready", "-U", "agents"); err != nil {
			return fmt.Errorf("PostgreSQL not ready: %w", err)
		}
	}

	printStatus("Creating required databases...")
	
	// Create agents database if it doesn't exist
	printStatus("Creating 'agents' database...")
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "-T", "postgres", "createdb", "-U", "agents", "agents"); err != nil {
		printWarning("Database 'agents' may already exist (this is OK)")
	}
	
	// Create convex_self_hosted database if it doesn't exist
	printStatus("Creating 'convex_self_hosted' database...")
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "-T", "postgres", "createdb", "-U", "agents", "convex_self_hosted"); err != nil {
		printWarning("Database 'convex_self_hosted' may already exist (this is OK)")
	}
	
	// Verify databases were created
	printStatus("Verifying database creation...")
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "-T", "postgres", "psql", "-U", "agents", "-l"); err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}
	
	printSuccess("Database initialization completed!")
	fmt.Println()
	printStatus("Created databases:")
	fmt.Println("  - agents: Used by the agents service")
	fmt.Println("  - convex_self_hosted: Used by the Convex backend")
	fmt.Println()
	printStatus("Connection strings:")
	fmt.Println("  - Agents: postgres://agents:agents@localhost:6888/agents?sslmode=disable")
	fmt.Println("  - Convex: postgres://agents:agents@localhost:6888?sslmode=disable")
	
	return nil
}

// ResetDB stops services, removes database volume, and reinitializes databases
func (Docker) ResetDB() error {
	printWarning("This will completely reset all databases and data!")
	fmt.Print("Are you sure? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		printStatus("Database reset cancelled.")
		return nil
	}

	printStatus("Resetting PostgreSQL databases...")
	
	// Stop all services
	printStatus("Stopping services...")
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "down"); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}
	
	// Remove postgres data volume
	printStatus("Removing postgres data volume...")
	if err := sh.RunV("docker", "volume", "rm", "docker_postgres_data"); err != nil {
		printWarning("Failed to remove postgres volume (may not exist)")
	}
	
	// Start postgres service (will trigger initialization)
	printStatus("Starting PostgreSQL with fresh volume...")
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "up", "-d", "postgres"); err != nil {
		return fmt.Errorf("failed to start PostgreSQL: %w", err)
	}
	
	// Wait for postgres to be ready
	printStatus("Waiting for PostgreSQL to initialize...")
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "postgres", "pg_isready", "-U", "agents"); err != nil {
		return fmt.Errorf("PostgreSQL not ready after reset: %w", err)
	}
	
	printSuccess("Database reset completed!")
	printStatus("The initialization scripts in docker/postgres-init/ have been executed.")
	
	return nil
}

// Help shows available Docker commands and usage examples
func (Docker) Help() {
	fmt.Println("Docker Management Commands")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Println("Available targets:")
	fmt.Println("  mage docker:up          Start all services")
	fmt.Println("  mage docker:down        Stop all services")
	fmt.Println("  mage docker:restart     Restart all services")
	fmt.Println("  mage docker:logs        Show logs for all services")
	fmt.Println("  mage docker:logs <svc>  Show logs for specific service")
	fmt.Println("  mage docker:build       Build all services")
	fmt.Println("  mage docker:build <svc> Build specific service")
	fmt.Println("  mage docker:status      Show service status")
	fmt.Println("  mage docker:clean       Remove all containers, networks, and volumes")
	fmt.Println("  mage docker:exec        Execute command in service container")
	fmt.Println("  mage docker:pull        Pull latest images")
	fmt.Println("  mage docker:config      Validate and show configuration")
	fmt.Println("  mage docker:core        Start only core services (agents + postgres)")
	fmt.Println("  mage docker:initdb      Initialize PostgreSQL databases")
	fmt.Println("  mage docker:resetdb     Reset all databases (removes all data)")
	fmt.Println("  mage docker:help        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mage docker:up                    # Start all services")
	fmt.Println("  mage docker:logs agents           # Show agents backend logs")
	fmt.Println("  mage docker:logs backend          # Show Convex backend logs")
	fmt.Println("  mage docker:build agents          # Rebuild agents service")
	fmt.Println("  mage docker:exec postgres psql    # Connect to PostgreSQL")
	fmt.Println("  mage docker:initdb                 # Create required databases")
	fmt.Println("  mage docker:resetdb                # Reset databases (WARNING: deletes all data)")
	fmt.Println()
	fmt.Println("Service URLs:")
	fmt.Println("  - Agents Backend: http://localhost:9451")
	fmt.Println("  - Convex Backend: http://localhost:3210")
	fmt.Println("  - Convex Dashboard: http://localhost:6791")
	fmt.Println("  - PostgreSQL Database: localhost:6888")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Environment variables can be customized in .env.docker (project root)")
	fmt.Println("  All services use the same .env.docker file for consistent configuration")
	fmt.Println()
	fmt.Println("Related Commands:")
	fmt.Println("  mage convex:help        Convex backend management (admin keys, status, etc.)")
}

// ensureEnvFile creates .env.docker file from .env.example if it doesn't exist
func ensureEnvFile() error {
	// We use the root .env.docker file, so check if it exists
	rootEnvPath := ".env.docker"

	// Check if root .env.docker file exists
	if _, err := os.Stat(rootEnvPath); os.IsNotExist(err) {
		return fmt.Errorf("root .env.docker file not found. Please create .env.docker in the project root directory")
	}

	printStatus("Using existing .env.docker file from project root")
	return nil
}

// Helper functions for colored output (reusing from existing helper.go style)
func printStatus(msg string) {
	fmt.Printf("\033[0;34m[INFO]\033[0m %s\n", msg)
}

func printSuccess(msg string) {
	fmt.Printf("\033[0;32m[SUCCESS]\033[0m %s\n", msg)
}

func printWarning(msg string) {
	fmt.Printf("\033[1;33m[WARNING]\033[0m %s\n", msg)
}

func printError(msg string) {
	fmt.Printf("\033[0;31m[ERROR]\033[0m %s\n", msg)
}
