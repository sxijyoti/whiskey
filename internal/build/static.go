package build

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(
		filepath.Dir(dst),
		0755,
	); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

func copyDir(srcRoot, dstRoot string) error {
	if _, err := os.Stat(srcRoot); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(
		srcRoot,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(srcRoot, path)
			if err != nil {
				return err
			}

			dst := filepath.Join(
				dstRoot,
				rel,
			)

			return copyFile(path, dst)
		},
	)
}

func CopyStatic(siteRoot string) error {
	themeStatic := filepath.Join(
		"themes",
		"default",
		"static",
	)

	siteStatic := filepath.Join(
		siteRoot,
		"static",
	)

	// copy theme first
	if err := copyDir(
		themeStatic,
		"dist",
	); err != nil {
		return err
	}

	// then site overrides
	if err := copyDir(
		siteStatic,
		"dist",
	); err != nil {
		return err
	}

	return nil
}