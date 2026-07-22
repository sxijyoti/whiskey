package build

import (
	"os"
	"path/filepath"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/template"
)

func RenderList(
	root string,
	cfg *config.Config,
	title string,
	pages []Page,
	output string,
) error {

	rendered, err := template.RenderPage(
		root,
		cfg.Theme,
		"list",
		template.PageData{
			Site:  cfg,
			Title: DisplayName(title),
			Pages: pages,
		},
	)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(output, rendered, 0644); err != nil {
		return err
	}

	return nil
}
