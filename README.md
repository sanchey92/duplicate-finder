# Duplicate Finder

A fast, concurrent file duplicate detector written in Go that identifies files with identical content using cryptographic hashing.

## Features

- **Fast Concurrent Processing** - Multi-worker goroutine pool for parallel file hashing
- **Multiple Hash Algorithms** - Support for MD5 (fast) and SHA256 (secure)
- **Flexible Output Formats** - Text, JSON, and CSV output formats
- **Progress Tracking** - Real-time progress bar during scanning
- **Comprehensive Statistics** - Total files scanned, duplicate groups, and wasted space calculation
- **Smart Filtering** - Automatically excludes empty files, symlinks, and non-regular files
- **Graceful Shutdown** - Handles SIGINT and SIGTERM for clean cancellation

## Installation

```bash
go install github.com/sanchey92/duplicate-finder/cmd@latest
```

Or build from source:

```bash
git clone https://github.com/sanchey92/duplicate-finder.git
cd duplicate-finder
go build -o duplicate-finder ./cmd
```

## Usage

### Basic Usage

```bash
# Scan a directory for duplicates
duplicate-finder -path /home/user/documents
```

### Advanced Usage

```bash
# Use SHA256 algorithm with 8 workers
duplicate-finder -path /home/user/photos -algo sha256 -workers 8

# Output results to JSON file
duplicate-finder -path /home/user/music -format json -output duplicates.json

# Quiet mode (no progress bar)
duplicate-finder -path /home/user/downloads -quiet

# Output to CSV for analysis
duplicate-finder -path /home/user/data -format csv -output report.csv
```

## Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-path` | **Required.** Directory to scan for duplicates | - |
| `-algo` | Hash algorithm to use (`md5`, `sha256`) | `md5` |
| `-workers` | Number of concurrent workers | CPU count |
| `-format` | Output format (`text`, `json`, `csv`) | `text` |
| `-output` | Output file path (use `-` for stdout) | stdout |
| `-progress` | Show progress bar | `true` |
| `-quiet` | Suppress progress output | `false` |
| `-help`, `-h` | Show help message | - |
| `-version`, `-v` | Show version information | - |

## Output Formats

### Text Format (Default)

```
Duplicate Finder v1.0.0
Path: /home/user/documents
Algorithm: md5
Workers: 8

Scanning... ████████████████████████████████ 100% (1234/1234)

Duplicate group #1 (3 files, 45.2 MB wasted):
  - /home/user/documents/file1.txt (15.1 MB)
  - /home/user/documents/backup/file1.txt (15.1 MB)
  - /home/user/documents/old/file1.txt (15.1 MB)

Total wasted space: 45.2 MB
```

### JSON Format

```json
{
  "stats": {
    "total_files": 1234,
    "processed_files": 1234,
    "duplicate_groups": 5,
    "wasted_space": 47448064
  },
  "duplicates": {
    "5d41402abc4b2a76b9719d911017c592": [
      {
        "path": "/home/user/documents/file1.txt",
        "size": 15728640,
        "hash": "5d41402abc4b2a76b9719d911017c592"
      }
    ]
  }
}
```

### CSV Format

```csv
Hash,FilePath,FileSize,GroupSize
5d41402abc4b2a76b9719d911017c592,/home/user/documents/file1.txt,15728640,3
5d41402abc4b2a76b9719d911017c592,/home/user/documents/backup/file1.txt,15728640,3
```

## Architecture

The project follows Go's standard directory layout with clean separation of concerns:

```
duplicate-finder/
├── cmd/                    # Application entry point
│   └── main.go
├── internal/              # Private application packages
│   ├── app/              # Application orchestration
│   ├── config/           # CLI configuration
│   ├── finder/           # Duplicate detection logic
│   ├── hash/             # File hashing
│   ├── types/            # Data type definitions
│   ├── walker/           # File system traversal
│   └── wp/               # Worker pool implementation
└── pkg/                   # Public packages
    └── formatter/        # Output formatting
```

### Key Components

- **Walker** - Recursively traverses directory tree and filters files
- **Hasher** - Computes cryptographic hashes using MD5 or SHA256
- **Worker Pool** - Manages concurrent file processing with configurable workers
- **Finder** - Orchestrates scanning, groups files by hash, identifies duplicates
- **Formatter** - Outputs results in text, JSON, or CSV formats
- **App** - High-level orchestration with signal handling

## How It Works

1. **File Collection** - Recursively scans the target directory, filtering out empty files, symlinks, and non-regular files
2. **Concurrent Hashing** - Distributes files across worker goroutines for parallel hash computation
3. **Grouping** - Groups files by their hash values
4. **Duplicate Detection** - Identifies groups with 2+ files (duplicates)
5. **Statistics** - Calculates total wasted space and scan metrics
6. **Output** - Formats and displays results in the chosen format

## Performance

- **Concurrent Processing**: Automatically scales to available CPU cores
- **Efficient I/O**: Worker pool pattern prevents resource exhaustion
- **Memory Efficient**: Streams files through processing pipeline
- **Fast Hashing**: MD5 algorithm optimized for speed; SHA256 for security

## Requirements

- Go 1.25.1 or higher
- No external dependencies (uses only Go standard library)

## Development

### Building

```bash
go build -o duplicate-finder ./cmd
```

## Version

Current version: **1.0.0**

## License

This project is open source. Please check the repository for license details.
