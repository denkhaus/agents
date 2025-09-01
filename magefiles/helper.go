package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func changeToConfigDirWithCleanup() (func(), error) {

	// Change to the config directory to run cue commands
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	if err := os.Chdir(configDir); err != nil {
		return nil, fmt.Errorf("failed to change directory to %s: %w", configDir, err)
	}

	return func() {
		os.Chdir(cwd)
	}, nil
}

// createReleaseFile creates a new release file based on current configurations
func createReleaseFile(version string) error {
	// Check if we're already in the config directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var releaseDir string
	if filepath.Base(currentDir) == "config" {
		// We're already in the config directory
		releaseDir = "./releases"
	} else {
		// We're in the main directory
		releaseDir = configDir + "/releases"
	}

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
	// Initialize with target version for all components
	versions := &ComponentVersions{
		Prompts:      targetVersion,
		Tools:        targetVersion,
		Settings:     targetVersion,
		Compositions: targetVersion,
	}

	// Get the last tag in the config submodule
	lastTag, err := getLastConfigTag()
	if err != nil {
		// If we can't get the last tag, use target version for all components
		return versions, nil
	}

	// If there's no previous tag, all components are at target version
	if lastTag == "" {
		return versions, nil
	}

	// Get previous component versions from the last release file
	prevVersions, err := getPreviousComponentVersions(lastTag)
	if err != nil {
		// If we can't get previous versions, use target version for all components
		return versions, nil
	}

	// Get the list of changed files between the last tag and HEAD in the config submodule
	changedFiles, err := getChangedConfigFiles(lastTag)
	if err != nil {
		// If we can't get changed files, use target version for all components
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

	// For components that didn't change, keep the previous version
	if !componentsChanged["prompts"] {
		versions.Prompts = prevVersions.Prompts
	}

	if !componentsChanged["tools"] {
		versions.Tools = prevVersions.Tools
	}

	if !componentsChanged["settings"] {
		versions.Settings = prevVersions.Settings
	}

	if !componentsChanged["compositions"] {
		versions.Compositions = prevVersions.Compositions
	}

	return versions, nil
}

// getChangedConfigFiles gets the list of changed files between two commits in the config submodule
func getChangedConfigFiles(sinceTag string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", sinceTag, "HEAD")
	cmd.Dir = "config"
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git diff in config submodule: %v", err)
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

// getPreviousComponentVersions extracts component versions from the last release file
func getPreviousComponentVersions(lastTag string) (*ComponentVersions, error) {
	// Convert tag to filename format (v1.0.0 -> v1_0_0.cue)
	tagParts := strings.Split(strings.TrimPrefix(lastTag, "v"), ".")
	if len(tagParts) != 3 {
		return nil, fmt.Errorf("invalid tag format: %s", lastTag)
	}

	filename := fmt.Sprintf("v%s_%s_%s.cue", tagParts[0], tagParts[1], tagParts[2])
	releaseFile := fmt.Sprintf("config/releases/%s", filename)

	// Read the release file
	content, err := os.ReadFile(releaseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read release file %s: %v", releaseFile, err)
	}

	// Parse component versions from the file
	versions := &ComponentVersions{}
	contentStr := string(content)

	// Extract prompts version
	if matches := regexp.MustCompile(`prompts:\s+"([^"]+)"`).FindStringSubmatch(contentStr); len(matches) > 1 {
		versions.Prompts = matches[1]
	}

	// Extract tools version
	if matches := regexp.MustCompile(`tools:\s+"([^"]+)"`).FindStringSubmatch(contentStr); len(matches) > 1 {
		versions.Tools = matches[1]
	}

	// Extract settings version
	if matches := regexp.MustCompile(`settings:\s+"([^"]+)"`).FindStringSubmatch(contentStr); len(matches) > 1 {
		versions.Settings = matches[1]
	}

	// Extract compositions version
	if matches := regexp.MustCompile(`compositions:\s+"([^"]+)"`).FindStringSubmatch(contentStr); len(matches) > 1 {
		versions.Compositions = matches[1]
	}

	return versions, nil
}

// getLastConfigTag gets the most recent Git tag in the config submodule
func getLastConfigTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = "config"
	output, err := cmd.Output()
	if err != nil {
		// No tags found
		return "", nil
	}

	return strings.TrimSpace(string(output)), nil
}

// restructureDirectories moves versioned directories to flattened structure
