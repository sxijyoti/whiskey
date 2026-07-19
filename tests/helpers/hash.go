package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

// HashDir walks a directory and creates a deterministic hash of all file names and contents.
func HashDir(t *testing.T, dir string) string {
	t.Helper()
	h := sha256.New()
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		
		// Hash relative path
		rel, _ := filepath.Rel(dir, path)
		h.Write([]byte(rel))
		
		// Hash content
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		
		return nil
	})
	
	if err != nil {
		t.Fatalf("Failed to hash directory %s: %v", dir, err)
	}
	
	return hex.EncodeToString(h.Sum(nil))
}