package main

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Cue.Validate

const (
	configDir             = "./config"
	projectCanvasDir      = "./contrib/project-canvas"
	dockerDir             = "./docker"
	convexGolangClientDir = "./convex_client"
)
