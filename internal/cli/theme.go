package cli

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage local themes",
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available themes",
	Args:  cobra.NoArgs,

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		entries, err := os.ReadDir("themes")
		if err != nil {
			return err
		}

		var themes []string

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			themes = append(
				themes,
				entry.Name(),
			)
		}

		themes = orderedThemes(themes)

		fmt.Println("Available themes")
		fmt.Println()

		for _, theme := range themes {
			fmt.Println(theme)
		}

		return nil
	},
}

var themeNewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new theme",
	Args:  cobra.ExactArgs(1),

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		name := args[0]

		dst := filepath.Join(
			"themes",
			name,
		)

		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf(
				"theme %q already exists",
				name,
			)
		}

		if err := copyThemeDir(
			filepath.Join("themes", "minimal"),
			dst,
		); err != nil {
			return err
		}

		fmt.Printf(
			"Created theme %s\n",
			name,
		)

		return nil
	},
}

func orderedThemes(themes []string) []string {
	seen := map[string]bool{}
	available := map[string]bool{}

	for _, theme := range themes {
		available[theme] = true
	}

	var ordered []string

	for _, theme := range []string{
		"minimal",
		"terminal",
		"paper",
	} {
		if !available[theme] {
			continue
		}

		ordered = append(
			ordered,
			theme,
		)
		seen[theme] = true
	}

	var rest []string
	for _, theme := range themes {
		if seen[theme] {
			continue
		}

		rest = append(
			rest,
			theme,
		)
	}

	sort.Strings(rest)

	return append(
		ordered,
		rest...,
	)
}

func copyThemeDir(srcRoot, dstRoot string) error {
	if _, err := os.Stat(srcRoot); err != nil {
		return err
	}

	return filepath.WalkDir(
		srcRoot,
		func(
			path string,
			d fs.DirEntry,
			err error,
		) error {

			if err != nil {
				return err
			}

			rel, err := filepath.Rel(
				srcRoot,
				path,
			)
			if err != nil {
				return err
			}

			dst := filepath.Join(
				dstRoot,
				rel,
			)

			if d.IsDir() {
				return os.MkdirAll(
					dst,
					0755,
				)
			}

			return copyThemeFile(
				path,
				dst,
			)
		},
	)
}

func copyThemeFile(src, dst string) error {
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

	_, err = io.Copy(
		out,
		in,
	)

	return err
}

func init() {
	themeCmd.AddCommand(
		themeListCmd,
		themeNewCmd,
	)

	rootCmd.AddCommand(
		themeCmd,
	)
}
