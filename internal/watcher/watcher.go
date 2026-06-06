package watcher

import (
	"io/fs"
	"path/filepath"

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

			return w.Add(path)
		},
	)
}