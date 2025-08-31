//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Cue is a namespace for cue related commands
type Cue mg.Namespace

// Update updates all cue definitions and dependencies.
func (Cue) Update() error {
	fmt.Println("Updating CUE definitions and dependencies...")

	configDir := "config"

	// Find all Go packages to generate CUE definitions from
	packages, err := findCUEGeneratedPackages(configDir)
	if err != nil {
		return fmt.Errorf("failed to find CUE generated packages: %w", err)
	}

	// Change to the config directory to run cue commands
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	if err := os.Chdir(configDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", configDir, err)
	}
	defer os.Chdir(cwd)

	// Update CUE dependencies
	fmt.Println("Running cue mod tidy...")
	if err := sh.RunV("cue", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run cue mod tidy: %w", err)
	}

	// Update generated CUE files from Go packages
	fmt.Println("Updating generated CUE files from Go packages...")
	for _, pkg := range packages {
		fmt.Printf("  - cue get go %s\n", pkg)
		if err := sh.RunV("cue", "get", "go", pkg); err != nil {
			fmt.Printf("Warning: failed to run cue get go %s: %v\n", pkg, err)
		}
	}

	fmt.Println("CUE definitions and dependencies updated successfully.")
	return nil
}

func findCUEGeneratedPackages(configDir string) ([]string, error) {
	genDir := filepath.Join(configDir, "cue.mod", "gen")
	packages := make(map[string]struct{})

	generateCommentRegex := regexp.MustCompile(`//cue:generate cue get go (.*)`)

	err := filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_go_gen.cue") {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			matches := generateCommentRegex.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				if len(match) > 1 {
					pkg := strings.TrimSpace(match[1])
					packages[pkg] = struct{}{}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var result []string
	for pkg := range packages {
		result = append(result, pkg)
	}
	sort.Strings(result)
	return result, nil
}