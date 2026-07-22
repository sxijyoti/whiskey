package build

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func DiscoverPages(root string) ([]string, error) {
	var pages []string

	err := filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if strings.HasPrefix(
				filepath.Base(path),
				".",
			) {
				return nil
			}

			if filepath.Ext(path) != ".md" {
				return nil
			}

			pages = append(
				pages,
				path,
			)

			return nil
		},
	)

	return pages, err
}
