package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/rearrange/opreiadrcligo/internal/core"
)

func TestNew_CreatesADRFile(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	result, err := core.New("Use Go for CLI")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if result.Number != 1 {
		t.Errorf("Number = %d, want 1", result.Number)
	}
	if result.Title != "Use Go for CLI" {
		t.Errorf("Title = %q, want %q", result.Title, "Use Go for CLI")
	}
	if _, err := os.Stat(result.FilePath); err != nil {
		t.Errorf("ADR file not found at %q: %v", result.FilePath, err)
	}
}

func TestNew_SlugifiesTitle(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	result, err := core.New("Use Go for CLI")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if !strings.HasSuffix(result.FilePath, "0001-use-go-for-cli.md") {
		t.Errorf("FilePath = %q, want suffix '0001-use-go-for-cli.md'", result.FilePath)
	}
}

func TestNew_SlugifiesSpecialCharacters(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	result, err := core.New("Choose: PostgreSQL vs. MySQL!")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if !strings.HasSuffix(result.FilePath, "0001-choose-postgresql-vs-mysql.md") {
		t.Errorf("FilePath = %q, want suffix '0001-choose-postgresql-vs-mysql.md'", result.FilePath)
	}
}

func TestNew_IncrementsNumberSequentially(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if _, err := core.New("First Decision"); err != nil {
		t.Fatalf("first New() error = %v", err)
	}
	result, err := core.New("Second Decision")
	if err != nil {
		t.Fatalf("second New() error = %v", err)
	}

	if result.Number != 2 {
		t.Errorf("Number = %d, want 2", result.Number)
	}
	if !strings.HasSuffix(result.FilePath, "0002-second-decision.md") {
		t.Errorf("FilePath = %q, want suffix '0002-second-decision.md'", result.FilePath)
	}
}

func TestNew_ADRContainsTitle(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	result, err := core.New("Adopt Hexagonal Architecture")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	data, err := os.ReadFile(result.FilePath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Adopt Hexagonal Architecture") {
		t.Errorf("ADR file missing title; content:\n%s", content)
	}
	if !strings.Contains(content, "Draft") {
		t.Errorf("ADR file missing default 'Draft' status; content:\n%s", content)
	}
}

func TestNew_UpdatesIndex(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if _, err := core.New("Use Go for CLI"); err != nil {
		t.Fatalf("New() error = %v", err)
	}

	data, err := os.ReadFile(core.IndexFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Use Go for CLI") {
		t.Error("index not updated with ADR title")
	}
	if !strings.Contains(content, "0001-use-go-for-cli.md") {
		t.Error("index not updated with ADR filename link")
	}
}

func TestNew_MultipleADRsAllAppearInIndex(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	titles := []string{"First Decision", "Second Decision", "Third Decision"}
	for _, title := range titles {
		if _, err := core.New(title); err != nil {
			t.Fatalf("New(%q) error = %v", title, err)
		}
	}

	data, err := os.ReadFile(core.IndexFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	content := string(data)
	for _, title := range titles {
		if !strings.Contains(content, title) {
			t.Errorf("index missing %q", title)
		}
	}
}

func TestNew_RequiresInit(t *testing.T) {
	chdir(t, t.TempDir())

	_, err := core.New("Some Decision")
	if err == nil {
		t.Fatal("New() expected error when workspace is not initialised, got nil")
	}
	if !strings.Contains(err.Error(), "adr init") {
		t.Errorf("error message should mention 'adr init'; got: %v", err)
	}
}
