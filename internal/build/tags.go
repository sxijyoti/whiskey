package build

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/template"
)

func BuildTags(
	root string,
	cfg *config.Config,
	index *SiteIndex,
) error {

	if err := os.MkdirAll(
		"dist/tags",
		0755,
	); err != nil {
		return err
	}

	tagNames := make(
		[]string,
		0,
		len(index.Tags),
	)

	for tag := range index.Tags {
		tagNames = append(
			tagNames,
			tag,
		)
	}

	sort.Strings(tagNames)

	indexPage, err := template.RenderPage(
		root,
		cfg.Theme,
		"tags",
		template.PageData{
			Site:        cfg,
			Title:       "Tags",
			Description: "Tags",
			Tags:        tagNames,
		},
	)

	if err != nil {
		return err
	}

	if err := os.WriteFile(
		"dist/tags/index.html",
		indexPage,
		0644,
	); err != nil {
		return err
	}

	for tag, pages := range index.Tags {

		if err := RenderList(
			root,
			cfg,
			tag,
			pages,
			filepath.Join(
				"dist",
				"tags",
				tag,
				"index.html",
			),
		); err != nil {
			return err
		}
	}

	return nil
}

func BuildCollections(
	root string,
	cfg *config.Config,
	index *SiteIndex,
) error {

	for name, pages := range index.Collections {

		if err := RenderList(
			root,
			cfg,
			name,
			pages,
			filepath.Join(
				"dist",
				name,
				"index.html",
			),
		); err != nil {
			return err
		}
	}

	return nil
}