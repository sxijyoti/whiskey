package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/sxijyoti/whiskey/tests/helpers"
)

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "whiskey-test-")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	helpers.WhiskeyBin = filepath.Join(tmpDir, "whiskey")
	cmd := exec.Command("go", "build", "-o", helpers.WhiskeyBin, "../../cmd/whiskey")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build whiskey binary: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}