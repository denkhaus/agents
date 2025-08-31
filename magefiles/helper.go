//go:build mage

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// createReleaseFile creates a new release file based on current configurations
func createReleaseFile(version string) error {
	releaseDir := configDir + "/releases"
	releaseFile := fmt.Sprintf("%s/%s.cue", releaseDir, strings.ReplaceAll(version, ".", "_"))

	// Check if release file already exists
	if _, err := os.Stat(releaseFile); err == nil {
		fmt.Printf("Release file %s already exists\n", releaseFile)
		return nil
	}

	// Get current date
	currentDate := time.Now().Format("2006-01-02")

	// Create release file content
	content := fmt.Sprintf(`package releases

import (
	"github.com/denkhaus/agent-config/compositions/stable"
)

%s: {
    version: "%s"
    release_date: "%s"
    description: "Release %s"

    // Component versions
    components: {
        prompts: "%s"
        tools: "%s"
        settings: "%s"
        compositions: "%s"
    }

    // Available agents in this release
    agents: {
        coder: stable.coder
        project_manager: stable.project_manager
        researcher: stable.researcher
    }

    // Compatibility matrix
    compatibility: {
        min_go_version: "1.25"
        min_cue_version: "v0.14.1"
        supported_environments: ["development", "production"]
    }
}
`, strings.ReplaceAll(version, ".", "_"), version, currentDate, version, version, version, version, version)

	// Write release file
	return os.WriteFile(releaseFile, []byte(content), 0644)
}

// restructureDirectories moves versioned directories to flattened structure
func restructureDirectories() error {
	fmt.Println("Restructuring directories for Git-based versioning...")

	// Move prompts from v1_0 to prompts root
	promptsSrc := configDir + "/prompts/v1_0"
	promptsDst := configDir + "/prompts"

	if _, err := os.Stat(promptsSrc); err == nil {
		fmt.Println("Moving prompts from v1_0 to root...")

		// Move each file
		files, err := os.ReadDir(promptsSrc)
		if err != nil {
			return err
		}

		for _, file := range files {
			src := fmt.Sprintf("%s/%s", promptsSrc, file.Name())
			dst := fmt.Sprintf("%s/%s", promptsDst, file.Name())

			// Rename existing files with version suffix before overwriting
			if _, err := os.Stat(dst); err == nil {
				versionedDst := fmt.Sprintf("%s.v1.0.0", dst)
				fmt.Printf("Renaming existing file %s to %s\n", dst, versionedDst)
				if err := os.Rename(dst, versionedDst); err != nil {
					return err
				}
			}

			fmt.Printf("Moving %s to %s\n", src, dst)
			if err := os.Rename(src, dst); err != nil {
				return err
			}
		}

		// Remove empty v1_0 directory
		fmt.Println("Removing empty v1_0 directory...")
		if err := os.Remove(promptsSrc); err != nil {
			log.Printf("Warning: Could not remove directory %s: %v", promptsSrc, err)
		}
	}

	// Move settings from default to settings root
	settingsSrc := configDir + "/settings/default"
	settingsDst := configDir + "/settings"

	if _, err := os.Stat(settingsSrc); err == nil {
		fmt.Println("Moving settings from default to root...")

		// Move each file
		files, err := os.ReadDir(settingsSrc)
		if err != nil {
			return err
		}

		for _, file := range files {
			src := fmt.Sprintf("%s/%s", settingsSrc, file.Name())
			dst := fmt.Sprintf("%s/%s", settingsDst, file.Name())

			// Rename existing files with version suffix before overwriting
			if _, err := os.Stat(dst); err == nil {
				versionedDst := fmt.Sprintf("%s.v1.0.0", dst)
				fmt.Printf("Renaming existing file %s to %s\n", dst, versionedDst)
				if err := os.Rename(dst, versionedDst); err != nil {
					return err
				}
			}

			fmt.Printf("Moving %s to %s\n", src, dst)
			if err := os.Rename(src, dst); err != nil {
				return err
			}
		}

		// Remove empty default directory
		fmt.Println("Removing empty default directory...")
		if err := os.Remove(settingsSrc); err != nil {
			log.Printf("Warning: Could not remove directory %s: %v", settingsSrc, err)
		}
	}

	return nil
}
