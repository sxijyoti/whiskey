package planner

import (
	"io/fs"
	"path/filepath"
)

func LocalDirtyPages(
	contentRoot string,
	state *State,
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

			info, err := d.Info()
			if err != nil {
				return err
			}

			if state.LastBuild.IsZero() {

				dirty = append(
					dirty,
					path,
				)

				return nil
			}

			if info.ModTime().After(
				state.LastBuild,
			) {

				dirty = append(
					dirty,
					path,
				)
			}

			return nil
		},
	)

	return dirty, err
}