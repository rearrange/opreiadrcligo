package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rearrange/opreiadrcligo/internal/core"
)

func TestList_RequiresInit(t *testing.T) {
	chdir(t, t.TempDir())

	_, err := core.List()
	if err == nil {
		t.Fatal("List() expected error when workspace is not initialised, got nil")
	}
	if !strings.Contains(err.Error(), "adr init") {
		t.Errorf("error should mention 'adr init'; got: %v", err)
	}
}

func TestList_EmptyWhenNoADRs(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(adrs) != 0 {
		t.Errorf("List() = %d entries, want 0", len(adrs))
	}
}

func TestList_ReturnsAllADRs(t *testing.T) {
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

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(adrs) != len(titles) {
		t.Fatalf("List() = %d entries, want %d", len(adrs), len(titles))
	}
}

func TestList_NumbersAreSequential(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	for _, title := range []string{"Alpha", "Beta", "Gamma"} {
		if _, err := core.New(title); err != nil {
			t.Fatalf("New(%q) error = %v", title, err)
		}
	}

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	for i, e := range adrs {
		if e.Number != i+1 {
			t.Errorf("adrs[%d].Number = %d, want %d", i, e.Number, i+1)
		}
	}
}

func TestList_ParsesTitlesCorrectly(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	want := "Use PostgreSQL as the primary database"
	if _, err := core.New(want); err != nil {
		t.Fatalf("New() error = %v", err)
	}

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if adrs[0].Title != want {
		t.Errorf("Title = %q, want %q", adrs[0].Title, want)
	}
}

func TestList_ReadsActualStatusFromFile(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	result, err := core.New("Adopt hexagonal architecture")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Simulate a user manually changing the status from Draft to Accepted.
	data, err := os.ReadFile(result.FilePath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	updated := strings.ReplaceAll(string(data), "| Draft |", "| Accepted |")
	if err := os.WriteFile(result.FilePath, []byte(updated), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if adrs[0].Status != "Accepted" {
		t.Errorf("Status = %q, want %q", adrs[0].Status, "Accepted")
	}
}

func TestList_StatusDefaultsToInitialDraft(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if _, err := core.New("Some decision"); err != nil {
		t.Fatalf("New() error = %v", err)
	}

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if adrs[0].Status != "Draft" {
		t.Errorf("Status = %q, want %q", adrs[0].Status, "Draft")
	}
}

func TestList_IgnoresNonADRFiles(t *testing.T) {
	chdir(t, t.TempDir())
	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if _, err := core.New("Real Decision"); err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Place a stray file that should be ignored.
	stray := filepath.Join(core.Dir, "notes.md")
	if err := os.WriteFile(stray, []byte("# Not an ADR"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	adrs, err := core.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(adrs) != 1 {
		t.Errorf("List() = %d entries, want 1 (stray file should be ignored)", len(adrs))
	}
}
