package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Convex contains targets for managing the Convex backend
type Convex mg.Namespace

// GenerateAdminKey generates an admin key for the Convex backend
func (Convex) GenerateAdminKey() error {
	printStatus("Generating Convex admin key...")

	// Check if backend service is running
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "ps", "backend"); err != nil {
		return fmt.Errorf("backend service not running. Start it with: mage docker:up")
	}

	// Generate admin key
	output, err := sh.Output("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "-T", "backend", "/convex/generate_admin_key.sh")
	if err != nil {
		return fmt.Errorf("failed to generate admin key: %w", err)
	}

	adminKey := strings.TrimSpace(output)
	if adminKey == "" {
		return fmt.Errorf("admin key generation returned empty result")
	}

	printSuccess("Admin key generated successfully!")
	fmt.Println()
	fmt.Printf("🔑 Admin Key: %s\n", adminKey)
	fmt.Println()
	printStatus("Usage:")
	fmt.Println("  - Use this key to authenticate with the Convex backend")
	fmt.Println("  - Store it securely in your application configuration")
	fmt.Println("  - Do not share this key publicly")
	fmt.Println()
	printStatus("Backend URL: http://localhost:3210")

	return nil
}

// ShowCredentials displays the current instance credentials
func (Convex) ShowCredentials() error {
	printStatus("Retrieving Convex instance credentials...")

	// Check if backend service is running
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "ps", "backend"); err != nil {
		return fmt.Errorf("backend service not running. Start it with: mage docker:up")
	}

	// Get instance name
	instanceName, err := sh.Output("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "-T", "backend", "cat", "/convex/data/credentials/instance_name")
	if err != nil {
		return fmt.Errorf("failed to read instance name: %w", err)
	}

	// Get instance secret (first 8 characters for security)
	instanceSecret, err := sh.Output("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "exec", "-T", "backend", "cat", "/convex/data/credentials/instance_secret")
	if err != nil {
		return fmt.Errorf("failed to read instance secret: %w", err)
	}

	instanceName = strings.TrimSpace(instanceName)
	instanceSecret = strings.TrimSpace(instanceSecret)

	fmt.Println("Convex Instance Credentials")
	fmt.Println("===========================")
	fmt.Printf("Instance Name: %s\n", instanceName)
	fmt.Printf("Instance Secret: %s...\n", instanceSecret[:8]) // Show only first 8 chars for security
	fmt.Println()
	printStatus("Backend URL: http://localhost:3210")
	printStatus("Dashboard URL: http://localhost:6791")

	return nil
}

// Status shows the status of the Convex backend
func (Convex) Status() error {
	printStatus("Checking Convex backend status...")

	// Check if backend service is running
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "ps", "backend"); err != nil {
		printWarning("Backend service not found or not running")
		return nil
	}

	// Check backend health
	healthOutput, err := sh.Output("curl", "-s", "-f", "http://localhost:3210/version")
	if err != nil {
		printWarning("Backend is not responding to health checks")
		printStatus("Try: mage docker:logs backend")
		return nil
	}

	// Show database connection info
	printStatus("Checking database connection...")
	dbLogs, err := sh.Output("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "logs", "--tail=50", "backend")
	if err == nil {
		if strings.Contains(dbLogs, "Connected to Postgres") {
			printSuccess("✅ Connected to PostgreSQL database")
		} else if strings.Contains(dbLogs, "Connected to SQLite") {
			printWarning("⚠️  Using SQLite database (PostgreSQL connection failed)")
		} else {
			printStatus("Database connection status unclear")
		}
	}

	printSuccess("✅ Convex backend is running and healthy")
	fmt.Println()
	fmt.Printf("Version info: %s\n", strings.TrimSpace(healthOutput))
	fmt.Println()
	printStatus("Service URLs:")
	fmt.Println("  - Backend API: http://localhost:3210")
	fmt.Println("  - Dashboard: http://localhost:6791")

	return nil
}

// Logs shows logs for the Convex backend
func (Convex) Logs() error {
	printStatus("Showing Convex backend logs...")
	return sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "logs", "-f", "backend")
}

// Restart restarts the Convex backend service
func (Convex) Restart() error {
	printStatus("Restarting Convex backend...")

	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "restart", "backend"); err != nil {
		return fmt.Errorf("failed to restart backend: %w", err)
	}

	printSuccess("Convex backend restarted successfully!")
	fmt.Println()
	printStatus("Use 'mage convex:status' to check the status")
	printStatus("Use 'mage convex:logs' to view logs")

	return nil
}

// Reset resets the Convex backend data (with confirmation)
func (Convex) Reset() error {
	printWarning("This will reset all Convex backend data including:")
	fmt.Println("  - All databases and tables")
	fmt.Println("  - All uploaded functions")
	fmt.Println("  - All stored files")
	fmt.Println("  - Instance credentials")
	fmt.Print("Are you sure? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		printStatus("Convex reset cancelled.")
		return nil
	}

	printStatus("Resetting Convex backend data...")

	// Stop backend service
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "stop", "backend"); err != nil {
		return fmt.Errorf("failed to stop backend: %w", err)
	}

	// Remove convex data volume
	if err := sh.RunV("docker", "volume", "rm", "docker_convex_data"); err != nil {
		printWarning("Failed to remove convex data volume (may not exist)")
	}

	// Start backend service with fresh data
	if err := sh.RunV("docker-compose", "-f", filepath.Join(dockerDir, "docker-compose.yml"), "up", "-d", "backend"); err != nil {
		return fmt.Errorf("failed to start backend: %w", err)
	}

	printSuccess("Convex backend reset completed!")
	fmt.Println()
	printStatus("New instance credentials have been generated")
	printStatus("Use 'mage convex:generateadminkey' to get a new admin key")
	printStatus("Use 'mage convex:status' to check the status")

	return nil
}

// ClearData clears all data from the Convex database
func (Convex) ClearData() error {
	printStatus("Clearing Convex database...")
	
	projectCanvasDir := filepath.Join("contrib", "project-canvas")
	
	// Check if package.json exists
	if err := checkFileExists(filepath.Join(projectCanvasDir, "package.json")); err != nil {
		return fmt.Errorf("package.json not found in %s: %w", projectCanvasDir, err)
	}
	
	// Change to project-canvas directory
	cleanup, err := changeToDirWithCleanup(projectCanvasDir)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", projectCanvasDir, err)
	}
	defer cleanup()
	
	// Ensure dependencies are installed
	printStatus("Checking npm dependencies...")
	if err := sh.RunV("npm", "install"); err != nil {
		return fmt.Errorf("failed to install npm dependencies: %w", err)
	}
	
	// Run the clear command
	err = sh.RunV("npx", "convex", "dev", "--once", "--run", "clearDatabase:clearDatabase")
	if err != nil {
		return fmt.Errorf("failed to clear database: %w", err)
	}

	printSuccess("✅ Convex database cleared successfully!")
	return nil
}

// SeedData seeds the Convex database with dummy data
func (Convex) SeedData() error {
	printStatus("Seeding Convex database with dummy data...")
	
	projectCanvasDir := filepath.Join("contrib", "project-canvas")
	
	// Check if package.json exists
	if err := checkFileExists(filepath.Join(projectCanvasDir, "package.json")); err != nil {
		return fmt.Errorf("package.json not found in %s: %w", projectCanvasDir, err)
	}
	
	// Change to project-canvas directory
	cleanup, err := changeToDirWithCleanup(projectCanvasDir)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", projectCanvasDir, err)
	}
	defer cleanup()
	
	// Ensure dependencies are installed
	printStatus("Checking npm dependencies...")
	if err := sh.RunV("npm", "install"); err != nil {
		return fmt.Errorf("failed to install npm dependencies: %w", err)
	}
	
	// Run the seed command
	err = sh.RunV("npx", "convex", "dev", "--once", "--run", "seed:seedDatabase")
	if err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	printSuccess("✅ Convex database seeded successfully!")
	return nil
}

// ResetData clears and re-seeds the Convex database
func (Convex) ResetData() error {
	printStatus("Resetting Convex database (clear + seed)...")
	
	// Clear first
	if err := (Convex{}).ClearData(); err != nil {
		return err
	}
	
	// Then seed
	if err := (Convex{}).SeedData(); err != nil {
		return err
	}
	
	printSuccess("✅ Convex database reset completed!")
	return nil
}

// Dev starts Convex development server
func (Convex) Dev() error {
	printStatus("Starting Convex development server...")
	
	projectCanvasDir := filepath.Join("contrib", "project-canvas")
	
	// Check if package.json exists
	if err := checkFileExists(filepath.Join(projectCanvasDir, "package.json")); err != nil {
		return fmt.Errorf("package.json not found in %s: %w", projectCanvasDir, err)
	}
	
	// Change to project-canvas directory
	cleanup, err := changeToDirWithCleanup(projectCanvasDir)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", projectCanvasDir, err)
	}
	defer cleanup()
	
	// Ensure dependencies are installed
	printStatus("Checking npm dependencies...")
	if err := sh.RunV("npm", "install"); err != nil {
		return fmt.Errorf("failed to install npm dependencies: %w", err)
	}
	
	// Start development server
	return sh.RunV("npx", "convex", "dev")
}

// Help shows available Convex commands and usage examples
func (Convex) Help() {
	fmt.Println("Convex Backend Management Commands")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("Available targets:")
	fmt.Println("  mage convex:generateadminkey    Generate admin key for authentication")
	fmt.Println("  mage convex:showcredentials     Show instance credentials")
	fmt.Println("  mage convex:status              Check backend status and health")
	fmt.Println("  mage convex:logs                Show backend logs (follow mode)")
	fmt.Println("  mage convex:restart             Restart backend service")
	fmt.Println("  mage convex:reset               Reset all backend data (WARNING: destructive)")
	fmt.Println("  mage convex:cleardata           Clear Convex database")
	fmt.Println("  mage convex:seeddata            Seed database with dummy data")
	fmt.Println("  mage convex:resetdata           Clear and re-seed database")
	fmt.Println("  mage convex:dev                 Start Convex development server")
	fmt.Println("  mage convex:help                Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mage convex:generateadminkey    # Generate admin key for API access")
	fmt.Println("  mage convex:status              # Check if backend is healthy")
	fmt.Println("  mage convex:logs                # View real-time logs")
	fmt.Println("  mage convex:showcredentials     # View instance name and partial secret")
	fmt.Println("  mage convex:resetdata           # Quick database reset with fresh data")
	fmt.Println()
	fmt.Println("Database Management:")
	fmt.Println("  mage convex:cleardata           # Remove all data from database")
	fmt.Println("  mage convex:seeddata            # Add dummy projects, tasks, agents")
	fmt.Println("  mage convex:resetdata           # Clear + seed in one command")
	fmt.Println()
	fmt.Println("Service URLs:")
	fmt.Println("  - Backend API: http://localhost:3210")
	fmt.Println("  - Dashboard: http://localhost:6791")
	fmt.Println()
	fmt.Println("Authentication:")
	fmt.Println("  Use the admin key from 'generateadminkey' to authenticate API requests")
	fmt.Println("  Admin keys are tied to the instance credentials and persist across restarts")
	fmt.Println()
	fmt.Println("Data Management:")
	fmt.Println("  - Backend data is stored in 'docker_convex_data' volume")
	fmt.Println("  - Database data is stored in PostgreSQL 'convex_self_hosted' database")
	fmt.Println("  - Use 'mage convex:reset' to completely start fresh")
	fmt.Println()
	fmt.Println("Related Commands:")
	fmt.Println("  mage docker:up          Start all services including Convex")
	fmt.Println("  mage docker:initdb      Initialize PostgreSQL databases")
	fmt.Println("  mage docker:status      Show status of all Docker services")
}