package build

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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

			if strings.HasPrefix(d.Name(), ".") {
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

func CopyStatic(
	siteRoot,
	theme string,
) error {

	themeStatic := filepath.Join(
		"themes",
		theme,
		"static",
	)

	siteStatic := filepath.Join(
		siteRoot,
		"static",
	)

	dist := filepath.Join(
		siteRoot,
		"dist",
	)

	// Copy theme assets first.
	if err := copyDir(
		themeStatic,
		dist,
	); err != nil {
		return err
	}

	// Site assets override theme assets.
	if err := copyDir(
		siteStatic,
		dist,
	); err != nil {
		return err
	}

	return nil
}

func DisplayName(
	name string,
) string {

	name = strings.ReplaceAll(
		name,
		"-",
		" ",
	)

	name = strings.ReplaceAll(
		name,
		"_",
		" ",
	)

	words := strings.Fields(name)

	for i, word := range words {

		runes := []rune(strings.ToLower(word))

		if len(runes) == 0 {
			continue
		}

		runes[0] = unicode.ToUpper(
			runes[0],
		)

		words[i] = string(runes)
	}

	return strings.Join(
		words,
		" ",
	)
}
