//go:build mage

package main

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = CueVersion.Validate

const (
	configDir = "./config"
)
