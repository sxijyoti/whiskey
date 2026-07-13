package dependency

import (
	"os"
	"path/filepath"
	"strings"
)

func CanonicalPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved, nil
	}

	return filepath.Clean(abs), nil
}

func ResolveLocalPath(contextFile, ref string) (string, error) {
	path := strings.TrimPrefix(ref, "local:")

	if !filepath.IsAbs(path) {
		base := filepath.Dir(contextFile)
		candidate := filepath.Join(base, path)
		if _, err := os.Stat(candidate); err == nil {
			return CanonicalPath(candidate)
		}

		for dir := base; ; dir = filepath.Dir(dir) {
			fallback := filepath.Join(dir, path)
			if _, err := os.Stat(fallback); err == nil {
				return CanonicalPath(fallback)
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}

		path = candidate
	}

	return CanonicalPath(path)
}

func ResolveLocalRef(contextFile, ref string) (string, error) {
	path, err := ResolveLocalPath(contextFile, ref)
	if err != nil {
		return "", err
	}

	return "local:" + path, nil
}
