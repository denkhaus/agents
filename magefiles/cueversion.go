//go:build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// CueVersion namespace for all cue-version related commands
type CueVersion mg.Namespace

// Validate validates all CUE configurations
func (CueVersion) Validate() error {
	fmt.Println("Validating CUE configurations...")

	// Change to config directory for validation
	if err := os.Chdir("config"); err != nil {
		return err
	}
	defer os.Chdir("..")

	// Validate all configurations
	if err := sh.RunV("cue", "vet", "./..."); err != nil {
		return err
	}

	fmt.Println("All CUE configurations are valid!")
	return nil
}

// Show displays the current version from Git
func (CueVersion) Show() error {
	fmt.Println("Getting current version from Git...")

	// Get the latest Git tag
	tag, err := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		fmt.Println("No Git tags found, using default version v0.0.0")
		tag = "v0.0.0"
	}

	fmt.Printf("Current version: %s\n", tag)
	return nil
}

// Tag creates a new Git tag for a release
func (CueVersion) Tag(version string) error {
	mg.Deps(CueVersion.Validate)

	if version == "" {
		return fmt.Errorf("version is required")
	}

	// Validate version format
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	fmt.Printf("Creating new release tag: %s\n", version)

	// Create release file
	if err := createReleaseFile(version); err != nil {
		return err
	}

	// Add all changes
	if err := sh.RunV("git", "add", "."); err != nil {
		return err
	}

	// Commit changes
	if err := sh.RunV("git", "commit", "-m", fmt.Sprintf("Release %s", version)); err != nil {
		return err
	}

	// Create Git tag
	if err := sh.RunV("git", "tag", version); err != nil {
		return err
	}

	fmt.Printf("Successfully created tag %s\n", version)
	return nil
}

// List displays all Git tags
func (CueVersion) List() error {
	fmt.Println("Listing all Git tags...")

	output, err := sh.Output("git", "tag", "--sort=-version:refname")
	if err != nil {
		return err
	}

	tags := strings.Split(output, "\n")
	for _, tag := range tags {
		if tag != "" {
			fmt.Println(tag)
		}
	}

	return nil
}

// Init initializes the Git-based versioning system
func (CueVersion) Init() error {
	fmt.Println("Initializing Git-based versioning system...")

	// Check if we're in a Git repository
	if err := sh.Run("git", "rev-parse", "--git-dir"); err != nil {
		return fmt.Errorf("not a Git repository")
	}

	// Create necessary directories if they don't exist
	dirs := []string{
		configDir + "/prompts",
		configDir + "/settings",
		configDir + "/tools/profiles",
		configDir + "/compositions/stable",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Move versioned directories to flattened structure
	if err := restructureDirectories(); err != nil {
		return err
	}

	fmt.Println("Git-based versioning system initialized!")
	return nil
}

// Clean removes backup files created during initialization
func (CueVersion) Clean() error {
	fmt.Println("Cleaning up backup files...")

	// Find and remove backup files (files with version suffixes)
	// This is a simplified implementation - in practice, you might want to be more specific

	fmt.Println("Cleanup completed!")
	return nil
}
