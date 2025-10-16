package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	var err error
	
	// Get log level from environment variable
	logLevel := os.Getenv("PM_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error" // Default to error level for clean CLI output
	}

	// Only enable logging in debug mode for CLI usage
	if logLevel == "debug" {
		// Development logger for debug mode
		Log, err = zap.NewDevelopment()
		if err != nil {
			Log = zap.NewNop()
		}
	} else {
		// For non-debug modes, use minimal logging to avoid JSON output
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel) // Only log errors
		config.OutputPaths = []string{"stderr"} // Send to stderr to not interfere with CLI output
		
		// Set log level if specified
		switch logLevel {
		case "warn":
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		case "off":
			Log = zap.NewNop() // No logging at all
			return
		}

		Log, err = config.Build()
		if err != nil {
			Log = zap.NewNop() // Fallback to no-op logger
		}
	}
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return Log
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		Log.Sync()
	}
}