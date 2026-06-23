package planner

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sxijyoti/whiskey/internal/fingerprint"
)

func LocalDirtyPages(
	contentRoot string,
	store fingerprint.Store,
) ([]string, error) {

	var dirty []string

	err := filepath.WalkDir(
		contentRoot,
		func(
			path string,
			d fs.DirEntry,
			err error,
		) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".md" {
				return nil
			}

			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			hash := fingerprint.SHA256(
				raw,
			)

			old := store[path]

			if old.Hash != hash {

				dirty = append(
					dirty,
					path,
				)

				store[path] =
					fingerprint.Entry{
						Hash: hash,
					}
			}

			return nil
		},
	)

	return dirty, err
}