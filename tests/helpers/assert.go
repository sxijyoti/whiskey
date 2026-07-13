package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func AssertBuildSucceeds(t *testing.T, res CmdResult) {
	t.Helper()
	if res.ExitCode != 0 {
		t.Fatalf("Expected build to succeed, got exit code %d\nStdout: %s\nStderr: %s", res.ExitCode, res.Stdout, res.Stderr)
	}
}

func AssertBuildFails(t *testing.T, res CmdResult) {
	t.Helper()
	if res.ExitCode == 0 {
		t.Fatalf("Expected build to fail, but it succeeded\nStdout: %s", res.Stdout)
	}
}

func AssertFileExists(t *testing.T, dir, path string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, path)); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to exist", path)
	}
}

func AssertFileNotExists(t *testing.T, dir, path string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, path)); err == nil {
		t.Fatalf("Expected file %s to NOT exist", path)
	}
}

func AssertContains(t *testing.T, dir, path, expected string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, path))
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	if !strings.Contains(string(b), expected) {
		t.Fatalf("Expected %s to contain %q.\nContent:\n%s", path, expected, string(b))
	}
}

func AssertNotContains(t *testing.T, dir, path, unexpected string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, path))
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	if strings.Contains(string(b), unexpected) {
		t.Fatalf("Expected %s to NOT contain %q.\nContent:\n%s", path, unexpected, string(b))
	}
}

func AssertWorkspaceExists(t *testing.T, dir string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, ".whiskey")); os.IsNotExist(err) {
		t.Fatalf("Expected .whiskey workspace to exist")
	}
}

func AssertOutputContains(t *testing.T, res CmdResult, expected string) {
	t.Helper()
	if !strings.Contains(res.Stdout, expected) && !strings.Contains(res.Stderr, expected) {
		t.Fatalf("Expected build output to contain %q.\nStdout:\n%s\nStderr:\n%s", expected, res.Stdout, res.Stderr)
	}
}

func AssertNothingChanged(t *testing.T, res CmdResult) {
	t.Helper()
	AssertBuildSucceeds(t, res)
	AssertOutputContains(t, res, "no changes")
}

func AssertIncrementalBuild(t *testing.T, res CmdResult) {
	t.Helper()
	AssertBuildSucceeds(t, res)
	AssertOutputContains(t, res, "Incremental")
}

func AssertFullBuild(t *testing.T, res CmdResult, reason string) {
	t.Helper()
	AssertBuildSucceeds(t, res)
	AssertOutputContains(t, res, "Full rebuild")
	if reason != "" {
		AssertOutputContains(t, res, reason)
	}
}
