package main

import (
	"fmt"
	"os"

	"github.com/sanchey92/duplicate-finder/internal/app"
	"github.com/sanchey92/duplicate-finder/internal/config"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if cfg.ShowHelp {
		os.Exit(0)
	}

	if cfg.ShowVersion {
		fmt.Printf("Duplicate Finder v%s\n", config.Version)
		os.Exit(0)
	}

	if cfg.Path == "" {
		fmt.Fprintln(os.Stderr, "Error: -path flag is required")
		fmt.Fprintln(os.Stderr, "Use -help for usage information")
		os.Exit(1)
	}

	application, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
