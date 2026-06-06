package template

import (
	"html/template"
	"os"
	"path/filepath"
)

func layoutFile(siteRoot, layout string) string {
	siteLayout := filepath.Join(
		siteRoot,
		"layouts",
		layout+".html",
	)

	if _, err := os.Stat(siteLayout); err == nil {
		return siteLayout
	}

	return filepath.Join(
		"themes",
		"default",
		"layouts",
		layout+".html",
	)
}

func LoadLayout(
	siteRoot string,
	layout string,
) (*template.Template, error) {

	base := filepath.Join(
		"themes",
		"default",
		"layouts",
		"base.html",
	)

	layoutPath := layoutFile(
		siteRoot,
		layout,
	)

	partials, err := filepath.Glob(
		filepath.Join(
			"themes",
			"default",
			"layouts",
			"partials",
			"*.html",
		),
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