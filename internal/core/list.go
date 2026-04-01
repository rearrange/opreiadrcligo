package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ADREntry holds the parsed metadata for a single ADR file.
type ADREntry struct {
	Number int
	Title  string
	Date   string
	Status string
}

// List scans the ADR directory and returns one ADREntry per ADR file, sorted
// by ADR number ascending.  The Status field reflects the actual value inside
// each file, so a manually updated "Accepted" is returned correctly instead
// of the original "Draft".
//
// List returns an error if the workspace has not been initialised.
func List() ([]ADREntry, error) {
	if !exists(Dir) {
		return nil, fmt.Errorf("ADR workspace not initialised; run 'adr init' first")
	}

	entries, err := os.ReadDir(Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read ADR directory: %w", err)
	}

	var adrs []ADREntry
	for _, e := range entries {
		if e.IsDir() || e.Name() == "README.md" {
			continue
		}
		var n int
		if _, err := fmt.Sscanf(e.Name(), "%d-", &n); err != nil {
			continue // skip non-ADR files (e.g. stray docs)
		}
		entry, err := parseADR(filepath.Join(Dir, e.Name()), n)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", e.Name(), err)
		}
		adrs = append(adrs, entry)
	}
	return adrs, nil
}

// parseADR reads a single ADR Markdown file and extracts its number, title,
// date, and current status.
//
// Expected file structure (produced by adr new):
//
//	# ADR 0001: Title
//	...
//	| Date       | Author | Status  |
//	|------------|--------|---------|
//	| 2 Apr 2026 | Alice  | Draft   |
func parseADR(path string, number int) (ADREntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ADREntry{}, err
	}

	var title, date, status string
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		// Title is on the first "# ADR NNNN: ..." heading.
		if title == "" && strings.HasPrefix(line, "# ADR ") {
			if idx := strings.Index(line, ": "); idx != -1 {
				title = strings.TrimSpace(line[idx+2:])
			}
			continue
		}

		// The data row immediately follows the |---| separator.
		if strings.HasPrefix(line, "|---") && i+1 < len(lines) {
			cells := strings.Split(lines[i+1], "|")
			// cells: ["", " date ", " author ", " status ", ""]
			if len(cells) >= 4 {
				date = strings.TrimSpace(cells[1])
				status = strings.TrimSpace(cells[3])
			}
			break
		}
	}

	return ADREntry{Number: number, Title: title, Date: date, Status: status}, nil
}
