// Package core contains the business logic for managing Architecture Decision
// Records. It is an internal package — only importable within this module —
// and has no dependency on any CLI framework.
package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// Dir is the canonical directory where ADR files are stored.
	Dir = "docs/adr"

	// IndexFile is the Markdown index that tracks all ADRs.
	IndexFile = "docs/adr/README.md"
)

// gitAuthor returns the git user.name for the current working directory.
//
// Resolution order:
//  1. `git config user.name` — works when inside a git repository; git
//     automatically applies matching includeIf.gitdir rules from ~/.gitconfig.
//  2. Manual includeIf.gitdir fallback — walks the entries returned by
//     `git config --list --global`, finds the first gitdir pattern that is a
//     prefix of the current directory, and reads user.name from that included
//     config file.  This handles the case where the tool is invoked outside of
//     any git repository (e.g. a brand-new project that has not yet run
//     `git init`).
//
// Returns an empty string when git is unavailable or no user.name is set.
func gitAuthor() string {
	if out, err := exec.Command("git", "config", "user.name").Output(); err == nil {
		if name := strings.TrimSpace(string(out)); name != "" {
			return name
		}
	}
	return gitAuthorFromIncludeIf()
}

// gitAuthorFromIncludeIf manually resolves includeIf.gitdir entries.
//
// `git config --list --global` emits lines such as:
//
//	includeif.gitdir:~/Codes/GitHub/.path=~/Codes/GitHub/.gitconfig-personal
//
// For each such line, the function checks whether the current directory is
// under the gitdir pattern path and, if so, reads user.name from the
// referenced config file via `git config --file <path> user.name`.
func gitAuthorFromIncludeIf() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	out, err := exec.Command("git", "config", "--list", "--global").Output()
	if err != nil {
		return ""
	}

	const prefix = "includeif.gitdir:"
	for _, line := range strings.Split(string(out), "\n") {
		lower := strings.ToLower(line)
		if !strings.HasPrefix(lower, prefix) {
			continue
		}

		// rest = "~/Codes/GitHub/.path=~/Codes/GitHub/.gitconfig-personal"
		rest := line[len(prefix):]
		dotPath := strings.Index(rest, ".path=")
		if dotPath == -1 {
			continue
		}
		gitdirPattern := expandTilde(rest[:dotPath], homeDir)
		configFile := expandTilde(rest[dotPath+len(".path="):], homeDir)

		// A trailing slash means "any subdirectory under this path".
		matchBase := strings.TrimSuffix(gitdirPattern, "/")
		if strings.HasPrefix(cwd, matchBase) {
			if name := readUserNameFromConfig(configFile); name != "" {
				return name
			}
		}
	}
	return ""
}

// expandTilde replaces a leading "~/" with the user's home directory.
func expandTilde(path, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// readUserNameFromConfig reads user.name from an arbitrary git config file.
func readUserNameFromConfig(path string) string {
	out, err := exec.Command("git", "config", "--file", path, "user.name").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
