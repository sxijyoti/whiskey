package helpers

import (
	"bytes"
	"os/exec"
	"testing"
)

var WhiskeyBin string

type CmdResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func RunWhiskey(t *testing.T, dir string, args ...string) CmdResult {
	t.Helper()
	cmd := exec.Command(WhiskeyBin, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			t.Fatalf("failed to run whiskey: %v", err)
		}
	}

	return CmdResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}
}

func BuildCmd(t *testing.T, dir string) CmdResult {
	// Pass "." to specify the current directory as the site root
	return RunWhiskey(t, dir, "build", ".")
}

func BuildFullCmd(t *testing.T, dir string) CmdResult {
	// Pass "." as the site root, followed by the flag
	return RunWhiskey(t, dir, "build", ".", "--full")
}

func BuildOfflineCmd(t *testing.T, dir string) CmdResult {
	// Pass "." as the site root, followed by the flag
	return RunWhiskey(t, dir, "build", ".", "--offline")
}
