package main

import (
	"log"
	"os"

	"github.com/denkhaus/knot/internal/app"
)

func main() {
	// Create and run the application
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
