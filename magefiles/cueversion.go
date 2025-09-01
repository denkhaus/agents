package main

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// CueVersion namespace for all cue-version related commands
type CueVersion mg.Namespace

// Show displays the current version from Git
func (CueVersion) Show() error {
	fmt.Println("Getting current version from Git...")

	cleanup, err := changeToConfigDirWithCleanup()
	if err != nil {
		return fmt.Errorf("failed to change to config directory: %w", err)
	}

	defer cleanup()

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
	mg.Deps(Cue.Validate)

	if version == "" {
		return fmt.Errorf("version is required")
	}

	// Validate version format
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	fmt.Printf("Creating new release tag: %s\n", version)

	cleanup, err := changeToConfigDirWithCleanup()
	if err != nil {
		return fmt.Errorf("failed to change to config directory: %w", err)
	}

	defer cleanup()

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

	// Push changes and tags
	fmt.Printf("Pushing changes and tags for %s\n", version)
	if err := sh.RunV("git", "push", "origin", "main"); err != nil {
		fmt.Printf("Warning: Failed to push changes: %v\n", err)
	}

	if err := sh.RunV("git", "push", "origin", version); err != nil {
		fmt.Printf("Warning: Failed to push tag: %v\n", err)
	}

	fmt.Printf("Successfully created and pushed tag %s\n", version)
	return nil
}

// List displays all Git tags
func (CueVersion) List() error {
	fmt.Println("Listing all Git tags...")

	cleanup, err := changeToConfigDirWithCleanup()
	if err != nil {
		return fmt.Errorf("failed to change to config directory: %w", err)
	}

	defer cleanup()

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
