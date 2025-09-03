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
	envFile   = ".env"
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
	fmt.Println("  - React Chat UI: http://localhost:3000")
	fmt.Println("  - ADK Web Frontend: http://localhost:4200")
	fmt.Println("  - Backend API: http://localhost:6999")
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

// ReactChat starts only the React Chat interface and backend
func (Docker) ReactChat() error {
	printStatus("Starting React Chat interface with backend...")

	if err := ensureEnvFile(); err != nil {
		return err
	}

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "up", "-d", "react-chat", "agents", "postgres"); err != nil {
		return fmt.Errorf("failed to start React Chat services: %w", err)
	}

	printSuccess("React Chat services started successfully!")
	fmt.Println()
	printStatus("Service URLs:")
	fmt.Println("  - React Chat UI: http://localhost:3000")
	fmt.Println("  - Backend API: http://localhost:6999")
	fmt.Println("  - PostgreSQL Database: localhost:6888")
	fmt.Println()
	printStatus("Use 'mage docker:logs react-chat' to view React Chat logs")
	printStatus("Use 'mage docker:down' to stop services")

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
	fmt.Println("  mage docker:reactchat   Start only React Chat interface")
	fmt.Println("  mage docker:help        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mage docker:up                    # Start all services")
	fmt.Println("  mage docker:logs react-chat       # Show React Chat logs")
	fmt.Println("  mage docker:logs web               # Show ADK web logs")
	fmt.Println("  mage docker:build react-chat      # Rebuild React Chat service")
	fmt.Println("  mage docker:build web              # Rebuild ADK web service")
	fmt.Println("  mage docker:exec postgres psql    # Connect to PostgreSQL")
	fmt.Println()
	fmt.Println("Service URLs:")
	fmt.Println("  - React Chat UI: http://localhost:3000")
	fmt.Println("  - ADK Web Frontend: http://localhost:4200")
	fmt.Println("  - Backend API: http://localhost:6999")
	fmt.Println("  - PostgreSQL Database: localhost:6888")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Environment variables can be customized in docker/.env")
	fmt.Println("  Backend URL: ADK_BACKEND_URL (default: http://host.docker.internal:6999)")
}

// ensureEnvFile creates .env file from .env.example if it doesn't exist
func ensureEnvFile() error {
	envPath := filepath.Join(dockerDir, envFile)
	examplePath := filepath.Join(dockerDir, envFile+".example")

	// Check if .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		printWarning(".env file not found. Creating from .env.example...")

		// Copy .env.example to .env
		if err := sh.Copy(envPath, examplePath); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}

		printSuccess("Created .env file. You can customize it if needed.")
	}

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
