package formatter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/sanchey92/duplicate-finder/internal/types"
)

type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatCSV  OutputFormat = "csv"
)

type Formatter struct {
	format OutputFormat
	writer io.Writer
}

func New(f OutputFormat, w io.Writer) (*Formatter, error) {
	switch f {
	case FormatText, FormatJSON, FormatCSV:
	default:
		return nil, fmt.Errorf("unsupported output format: %s", f)
	}

	if w == nil {
		w = os.Stdout
	}

	return &Formatter{
		format: f,
		writer: w,
	}, nil
}

func (f *Formatter) PrintResults(duplicates types.DuplicateGroup, stats *types.ScanStats) error {
	switch f.format {
	case FormatText:
		return f.printText(duplicates, stats)
	case FormatJSON:
		return f.printJSON(duplicates, stats)
	case FormatCSV:
		return f.printCSV(duplicates)
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

func (f *Formatter) printText(duplicates types.DuplicateGroup, stats *types.ScanStats) error {
	if len(duplicates) == 0 {
		_, _ = fmt.Fprintln(f.writer, "No duplicates found!")
		return nil
	}

	_, _ = fmt.Fprintf(f.writer, "\n=== Found %d groups of duplicate files ===\n\n\n", len(duplicates))

	hashes := make([]string, 0, len(duplicates))

	for hash := range duplicates {
		hashes = append(hashes, hash)
	}
	sort.Strings(hashes)

	for i, hash := range hashes {
		group := duplicates[hash]
		if len(group) == 0 {
			continue
		}

		_, _ = fmt.Fprintf(f.writer, "Group %d:\n", i+1)

		for j, file := range group {
			_, _ = fmt.Fprintf(f.writer, "  %d. %s (%s)\n", j+1, file.Path, formatSize(file.Size))
		}

		wastedSize := int64(len(group)-1) * group[0].Size
		_, _ = fmt.Fprintf(f.writer, "  → Wasted space: %s\n\n", formatSize(wastedSize))
	}

	if stats != nil {
		f.printSummary(stats)
	}
	return nil
}

func (f *Formatter) printJSON(duplicates types.DuplicateGroup, stats *types.ScanStats) error {
	result := struct {
		Duplicates types.DuplicateGroup `json:"duplicates"`
		Stats      *types.ScanStats     `json:"stats"`
	}{
		Duplicates: duplicates,
		Stats:      stats,
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", " ")
	return encoder.Encode(result)
}

func (f *Formatter) printCSV(duplicates types.DuplicateGroup) error {
	writer := csv.NewWriter(f.writer)
	defer writer.Flush()

	if err := writer.Write([]string{"Hash", "FilePath", "FileSize", "GroupSize"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for h, g := range duplicates {
		for _, file := range g {
			record := []string{
				h,
				file.Path,
				strconv.FormatInt(file.Size, 10),
				strconv.Itoa(len(g)),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write CSV record: %w", err)
			}
		}
	}

	return nil
}

func (f *Formatter) printSummary(stats *types.ScanStats) {
	if stats == nil {
		return
	}
	_, _ = fmt.Fprintln(f.writer, "=== Summary ===")
	_, _ = fmt.Fprintf(f.writer, "Total files scanned: %d\n", stats.TotalFiles)
	_, _ = fmt.Fprintf(f.writer, "Files processed: %d\n", stats.ProcessedFiles)
	_, _ = fmt.Fprintf(f.writer, "Duplicate groups: %d\n", stats.DuplicateGroup)
	_, _ = fmt.Fprintf(f.writer, "Total wasted space: %s\n", formatSize(stats.TotalWastedSpace))
}

func (f *Formatter) PrintHeader(rootPath, hashAlgo string, numWorkers int) {
	_, _ = fmt.Fprintln(f.writer, "═══════════════════════════════════════")
	_, _ = fmt.Fprintln(f.writer, "     File Duplicates Finder")
	_, _ = fmt.Fprintln(f.writer, "═══════════════════════════════════════")
	_, _ = fmt.Fprintf(f.writer, "Path: %s\n", rootPath)
	_, _ = fmt.Fprintf(f.writer, "Algorithm: %s\n", hashAlgo)
	_, _ = fmt.Fprintf(f.writer, "Workers: %d\n", numWorkers)
	_, _ = fmt.Fprintln(f.writer, "───────────────────────────────────────")
}

func formatSize(bytes int64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0

	for n := bytes / unit; n > unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (f *Formatter) PrintProgressBar(current, total int) {
	if total == 0 {
		return
	}

	percentage := float64(current) / float64(total) * 100
	_, _ = fmt.Fprintf(f.writer, "\rProgress: %.1f%% (%d/%d)", percentage, current, total)

	if current == total {
		_, _ = fmt.Fprintln(f.writer)
	}
}

func IsValidFormat(format string) bool {
	switch OutputFormat(format) {
	case FormatText, FormatJSON, FormatCSV:
		return true
	default:
		return false
	}
}

func GetSupportedFormats() []string {
	return []string{
		string(FormatText),
		string(FormatJSON),
		string(FormatCSV),
	}
}
