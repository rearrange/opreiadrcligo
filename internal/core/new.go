package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rearrange/opreiadrcligo/internal/template"
)

// NewResult holds the ADR created by New.
type NewResult struct {
	Number   int
	FilePath string
	Title    string
}

// New creates a new ADR file for the given title, then updates the index.
//
// The file is placed in docs/adr/ with a sequential 4-digit prefix and the
// title slugified as the basename, e.g. 0003-use-postgresql.md.
// The docs/adr/README.md index is updated with a new table row.
//
// New returns an error if the directory or index file has not been initialised (adr init).
func New(title string) (*NewResult, error) {
	if !exists(Dir) || !exists(IndexFile) {
		return nil, fmt.Errorf("ADR directory not initialized; run 'adr init' first")
	}

	number, err := nextNumber(Dir)
	if err != nil {
		return nil, fmt.Errorf("Failed to determine next ADR number: %w", err)
	}

	date := time.Now().Format("2 Jan 2006")
	author := gitAuthor()
	slug := slugify(title)
	filename := fmt.Sprintf("%04d-%s.md", number, slug)
	filePath := filepath.Join(Dir, filename)

	content, err := template.RenderADR(template.ADRData{
		Number: number,
		Title:  title,
		Date:   date,
		Author: author,
		Status: "Draft",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render ADR template: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write ADR file %q: %w", filePath, err)
	}

	if err := updateIndex(IndexFile, number, title, filename, date, author); err != nil {
		return nil, fmt.Errorf("failed to update index: %w", err)
	}

	return &NewResult{Number: number, FilePath: filePath, Title: title}, nil
}

// nextNumber scans dir for existing ADR files (NNNN-*.md) and returns max+1.
func nextNumber(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	max := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		var n int
		if _, err := fmt.Sscanf(e.Name(), "%d-", &n); err == nil && n > max {
			max = n
		}
	}
	return max + 1, nil
}

// slugify converts a title to a URL-safe lowercase hyphen-separated string.
func slugify(title string) string {
	s := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// updateIndex appends a new row to the Markdown table inside indexFile.
//
// It locates the table separator line (|---|...) and inserts the new row
// immediately after any existing data rows, preserving the rest of the file.
func updateIndex(indexFile string, number int, title, filename, date, author string) error {
	data, err := os.ReadFile(indexFile)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	// Find the table header separator line.
	sepIdx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "|---") {
			sepIdx = i
			break
		}
	}
	if sepIdx == -1 {
		return fmt.Errorf("index file is missing a Markdown table; cannot update")
	}

	// Advance past any existing data rows.
	insertIdx := sepIdx + 1
	for insertIdx < len(lines) && strings.HasPrefix(lines[insertIdx], "|") {
		insertIdx++
	}

	row := fmt.Sprintf("| [%04d](%s) | %s | %s | %s |", number, filename, title, date, author)

	updated := make([]string, 0, len(lines)+1)
	updated = append(updated, lines[:insertIdx]...)
	updated = append(updated, row)
	updated = append(updated, lines[insertIdx:]...)

	return os.WriteFile(indexFile, []byte(strings.Join(updated, "\n")), 0644)
}
