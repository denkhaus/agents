//go:build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

	// Get component versions based on actual changes
	componentVersions, err := getComponentVersions(version)
	if err != nil {
		return fmt.Errorf("failed to determine component versions: %v", err)
	}

	// Get current date
	currentDate := time.Now().Format("2006-01-02")

	// Create release file content with independent component versions
	content := fmt.Sprintf(`package releases

import (
	"github.com/denkhaus/agent-config/compositions/stable"
)

%s: {
    version: "%s"
    release_date: "%s"
    description: "Release %s"

    // Component versions - set independently based on actual changes
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
`, strings.ReplaceAll(version, ".", "_"), version, currentDate, version,
   componentVersions.Prompts, componentVersions.Tools,
   componentVersions.Settings, componentVersions.Compositions)

	// Write release file
	return os.WriteFile(releaseFile, []byte(content), 0644)
}

// ComponentVersions holds version information for each component
type ComponentVersions struct {
	Prompts      string
	Tools        string
	Settings     string
	Compositions string
}

// getComponentVersions determines component versions based on actual changes
func getComponentVersions(targetVersion string) (*ComponentVersions, error) {
	versions := &ComponentVersions{
		Prompts:      "v1.0.0",
		Tools:        "v1.0.0",
		Settings:     "v1.0.0",
		Compositions: "v1.0.0",
	}

	// Get the last tag
	lastTag, err := getLastTag()
	if err != nil {
		// If we can't get the last tag, use target version for all components
		versions.Prompts = targetVersion
		versions.Tools = targetVersion
		versions.Settings = targetVersion
		versions.Compositions = targetVersion
		return versions, nil
	}

	// If there's no previous tag, all components are at target version
	if lastTag == "" {
		versions.Prompts = targetVersion
		versions.Tools = targetVersion
		versions.Settings = targetVersion
		versions.Compositions = targetVersion
		return versions, nil
	}

	// Get the list of changed files between the last tag and HEAD
	changedFiles, err := getChangedFiles(lastTag)
	if err != nil {
		// If we can't get changed files, use target version for all components
		versions.Prompts = targetVersion
		versions.Tools = targetVersion
		versions.Settings = targetVersion
		versions.Compositions = targetVersion
		return versions, nil
	}

	// Track which components have changed
	componentsChanged := map[string]bool{
		"prompts":      false,
		"tools":        false,
		"settings":     false,
		"compositions": false,
	}

	// Analyze which components have changed
	for _, file := range changedFiles {
		if file == "" {
			continue
		}
		
		// Determine which component this file belongs to
		switch {
		case strings.HasPrefix(file, "prompts/") || strings.Contains(file, "prompts"):
			componentsChanged["prompts"] = true
		case strings.HasPrefix(file, "tools/") || strings.Contains(file, "tools"):
			componentsChanged["tools"] = true
		case strings.HasPrefix(file, "settings/") || strings.Contains(file, "settings"):
			componentsChanged["settings"] = true
		case strings.HasPrefix(file, "compositions/") || strings.Contains(file, "compositions"):
			componentsChanged["compositions"] = true
		}
	}

	// For components that changed, use the target version
	// For components that didn't change, keep the previous version (v1.0.0 for now)
	if componentsChanged["prompts"] {
		versions.Prompts = targetVersion
	} else {
		versions.Prompts = "v1.0.0"
	}
	
	if componentsChanged["tools"] {
		versions.Tools = targetVersion
	} else {
		versions.Tools = "v1.0.0"
	}
	
	if componentsChanged["settings"] {
		versions.Settings = targetVersion
	} else {
		versions.Settings = "v1.0.0"
	}
	
	if componentsChanged["compositions"] {
		versions.Compositions = targetVersion
	} else {
		versions.Compositions = "v1.0.0"
	}

	return versions, nil
}

// getChangedFiles gets the list of changed files between two commits
func getChangedFiles(sinceTag string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", sinceTag, "HEAD")
	cmd.Dir = "config"
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git diff: %v", err)
	}

	// Split output into lines
	lines := strings.Split(string(output), "\n")
	
	// Filter out empty lines
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}
	
	return files, nil
}

// getLastTag gets the most recent Git tag
func getLastTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = "."
	output, err := cmd.Output()
	if err != nil {
		// No tags found
		return "", nil
	}
	
	return strings.TrimSpace(string(output)), nil
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
