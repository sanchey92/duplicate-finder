package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

const Version = "1.0.0"

type Config struct {
	// Scan params
	Path      string
	Algorithm string
	Workers   int
	// Output params
	Format       string
	OutputFile   string
	Quiet        bool
	ShowProgress bool
	// System flags
	ShowHelp    bool
	ShowVersion bool
}

func ParseFlags() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.Path, "path", "", "Directory path to scan for duplicates (required)")
	flag.StringVar(&config.Algorithm, "algo", "md5", "Hash algorithm to use (md5, sha256)")
	flag.IntVar(&config.Workers, "workers", 0, "Number of worker goroutines (default: number of CPUs)")
	flag.StringVar(&config.Format, "format", "text", "Output format (text, json, csv)")
	flag.StringVar(&config.OutputFile, "output", "", "Output file path (default: stdout)")
	flag.BoolVar(&config.Quiet, "quiet", false, "Suppress progress output")
	flag.BoolVar(&config.ShowProgress, "progress", true, "Show progress bar")
	flag.BoolVar(&config.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&config.ShowVersion, "v", false, "Show version information (shorthand)")

	flag.Usage = printUsage

	flag.Parse()

	if config.Workers <= 0 {
		config.Workers = runtime.NumCPU()
	}

	return config, nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `File Duplicates Finder v%s

DESCRIPTION:
    A fast, concurrent file duplicate finder that uses cryptographic hashes
    to identify files with identical content.

USAGE:
    %s [OPTIONS]

OPTIONS:
    -path <directory>     Directory path to scan for duplicates (required)
    -algo <algorithm>     Hash algorithm: md5, sha256 (default: md5)
    -workers <number>     Number of worker goroutines (default: %d)
    -format <format>      Output format: text, json, csv (default: text)
    -output <file>        Output file path (default: stdout)
    -progress             Show progress bar (default: true)
    -quiet                Suppress progress output
    -help, -h             Show this help message
    -version, -v          Show version information

EXAMPLES:
    # Basic usage
    %s -path /home/user/documents

    # Use SHA256 and 8 workers
    %s -path /home/user/photos -algo sha256 -workers 8

    # Output to JSON file
    %s -path /home/user/music -format json -output duplicates.json

    # Quiet mode (no progress indicators)
    %s -path /home/user/downloads -quiet

SUPPORTED ALGORITHMS:
    md5     - Fast, good for most use cases
    sha256  - More secure, slower but recommended for important data

OUTPUT FORMATS:
    text    - Human-readable format with file groupings
    json    - Machine-readable JSON with complete metadata
    csv     - Spreadsheet-compatible format for analysis

`, Version, os.Args[0], runtime.NumCPU(), os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
