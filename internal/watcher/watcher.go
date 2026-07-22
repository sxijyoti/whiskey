package watcher

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func WatchRecursive(
	w *fsnotify.Watcher,
	root string,
) error {

	return filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			if strings.HasPrefix(
				filepath.Base(path),
				".",
			) {
				return nil
			}

			return w.Add(path)
		},
	)
}

func IgnoreFile(path string) bool {

	name := filepath.Base(path)

	switch {

	case name == ".DS_Store":
		return true

	case strings.HasSuffix(name, "~"):
		return true

	case strings.HasSuffix(name, ".swp"):
		return true

	case strings.HasSuffix(name, ".tmp"):
		return true

	case strings.HasPrefix(name, ".#"):
		return true

	case strings.Contains(
		path,
		string(filepath.Separator)+"dist"+string(filepath.Separator),
	):
		return true

	case strings.Contains(
		path,
		string(filepath.Separator)+".whiskey"+string(filepath.Separator),
	):
		return true

	default:
		return false
	}
}

func IsStatic(path string) bool {

	return strings.Contains(
		path,
		string(filepath.Separator)+"static"+string(filepath.Separator),
	)
}

func IsLayout(path string) bool {

	return strings.Contains(
		path,
		string(filepath.Separator)+"layouts"+string(filepath.Separator),
	)
}

func IsConfig(path string) bool {

	return filepath.Base(path) == "whiskey.toml"
}
