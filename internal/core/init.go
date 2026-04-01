package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/rearrange/opreiadrcligo/internal/template"
)

// InitResult holds the paths created by Init so the caller can report them.
type InitResult struct {
	Dir        string
	IndexFile  string
	CreatedDir bool // true if the directory was created, false if it already existed
}

// Init creates the ADR directory and index file.
//
// Behaviour by case:
//   - Neither dir nor index exist  → create both (fresh init)
//   - Dir exists, index is missing → create the index only (recovery)
//   - Both already exist           → nothing to do, return an error
func Init() (*InitResult, error) {
	dirExists := exists(Dir)
	indexExists := exists(IndexFile)

	// Fully initialised — nothing to do.
	if dirExists && indexExists {
		return nil, fmt.Errorf(
			"Already initialised: both %q and %q exist.\n"+
				"To start fresh, remove the directory manually: rm -rf %s",
			Dir, IndexFile, Dir,
		)
	}

	// Create the directory if it is missing.
	if !dirExists {
		if err := os.MkdirAll(Dir, 0755); err != nil {
			return nil, fmt.Errorf("Failed to create directory %q: %w", Dir, err)
		}
	}

	// Create the index file if it is missing.
	if !indexExists {
		content, err := template.RenderIndex()
		if err != nil {
			return nil, fmt.Errorf("Failed to render index template: %w", err)
		}
		if err := os.WriteFile(IndexFile, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("Failed to create index file %q: %w", IndexFile, err)
		}
	}

	return &InitResult{Dir: Dir, IndexFile: IndexFile, CreatedDir: !dirExists}, nil
}

// exists reports whether path exists on the filesystem.
func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
