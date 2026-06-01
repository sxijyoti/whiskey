package template

import (
	"html/template"
	"os"
	"path/filepath"
)

func LoadLayout(
	siteRoot string,
	layout string,
) (*template.Template, error) {

	siteLayout := filepath.Join(
		siteRoot,
		"layouts",
		layout+".html",
	)

	if _, err := os.Stat(siteLayout); err == nil {
		return template.ParseFiles(siteLayout)
	}

	defaultLayout := filepath.Join(
		"themes",
		"default",
		"layouts",
		layout+".html",
	)

	return template.ParseFiles(defaultLayout)
}
