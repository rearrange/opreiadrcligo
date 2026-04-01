package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/rearrange/opreiadrcligo/internal/core"
)

// chdir changes into dir and restores the original working directory on
// test cleanup. Go 1.24+ provides t.Chdir, but we implement it manually here
// to be explicit about the lifecycle.
func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("chdir: getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir cleanup: %v", err)
		}
	})
}

func TestInit_FreshInit_CreatesBothDirAndIndex(t *testing.T) {
	chdir(t, t.TempDir())

	result, err := core.Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if result.Dir != core.Dir {
		t.Errorf("Dir = %q, want %q", result.Dir, core.Dir)
	}
	if result.IndexFile != core.IndexFile {
		t.Errorf("IndexFile = %q, want %q", result.IndexFile, core.IndexFile)
	}
	if !result.CreatedDir {
		t.Error("CreatedDir = false, want true")
	}

	if _, err := os.Stat(core.Dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
	if _, err := os.Stat(core.IndexFile); err != nil {
		t.Errorf("index file not created: %v", err)
	}
}

func TestInit_FreshInit_IndexContainsExpectedContent(t *testing.T) {
	chdir(t, t.TempDir())

	if _, err := core.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	data, err := os.ReadFile(core.IndexFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	content := string(data)
	for _, want := range []string{"Architecture Decision Records", "## Index"} {
		if !strings.Contains(content, want) {
			t.Errorf("index missing %q", want)
		}
	}
}

func TestInit_AlreadyInitialised_ReturnsError(t *testing.T) {
	chdir(t, t.TempDir())

	if _, err := core.Init(); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}

	_, err := core.Init()
	if err == nil {
		t.Fatal("second Init() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Already initialised") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestInit_RecoveryMode_CreatesIndexWhenDirExists(t *testing.T) {
	chdir(t, t.TempDir())

	// Pre-create only the directory, not the index file.
	if err := os.MkdirAll(core.Dir, 0755); err != nil {
		t.Fatalf("setup MkdirAll error = %v", err)
	}

	result, err := core.Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result.CreatedDir {
		t.Error("CreatedDir = true, want false (dir was pre-existing)")
	}
	if _, err := os.Stat(core.IndexFile); err != nil {
		t.Errorf("index file not created in recovery mode: %v", err)
	}
}
