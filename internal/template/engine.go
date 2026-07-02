package template

import (
	"html/template"
	"os"
	"path/filepath"
	"sort"
)

func resolveLayoutFile(siteRoot, theme, layout string) string {
	return resolveThemeFile(
		siteRoot,
		theme,
		filepath.Join(
			"layouts",
			layout+".html",
		),
	)
}

func resolveThemeFile(siteRoot, theme, rel string) string {
	siteFile := filepath.Join(
		siteRoot,
		rel,
	)

	if _, err := os.Stat(siteFile); err == nil {
		return siteFile
	}

	return filepath.Join(
		"themes",
		theme,
		rel,
	)
}

func resolvedPartials(siteRoot, theme string) ([]string, error) {
	partials := map[string]string{}

	themePartials, err := filepath.Glob(
		filepath.Join(
			"themes",
			theme,
			"layouts",
			"partials",
			"*.html",
		),
	)
	if err != nil {
		return nil, err
	}

	for _, partial := range themePartials {
		partials[filepath.Base(partial)] = partial
	}

	sitePartials, err := filepath.Glob(
		filepath.Join(
			siteRoot,
			"layouts",
			"partials",
			"*.html",
		),
	)
	if err != nil {
		return nil, err
	}

	for _, partial := range sitePartials {
		partials[filepath.Base(partial)] = partial
	}

	names := make([]string, 0, len(partials))
	for name := range partials {
		names = append(names, name)
	}

	sort.Strings(names)

	files := make([]string, 0, len(names))
	for _, name := range names {
		files = append(files, partials[name])
	}

	return files, nil
}

func LoadLayout(
	siteRoot string,
	theme string,
	layout string,
) (*template.Template, error) {

	base := resolveThemeFile(
		siteRoot,
		theme,
		filepath.Join(
			"layouts",
			"base.html",
		),
	)

	layoutPath := resolveLayoutFile(
		siteRoot,
		theme,
		layout,
	)

	partials, err := resolvedPartials(
		siteRoot,
		theme,
	)
	if err != nil {
		return nil, err
	}

	files := []string{
		base,
		layoutPath,
	}

	files = append(
		files,
		partials...,
	)

	return template.ParseFiles(files...)
}
