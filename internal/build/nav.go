package build

import (
	"sort"

	"github.com/sxijyoti/whiskey/internal/config"
)

func BuildNav(
	index *SiteIndex,
) []config.NavItem {

	contentPages := map[string]config.NavItem{}

	collections := map[string]config.NavItem{}

	seenTags := false

	for _, page := range index.Pages {

		if len(page.Tags) > 0 {
			seenTags = true
		}

		if page.Collection != "" {

			if _, ok := collections[
				page.Collection,
			]; !ok {

				collections[
					page.Collection,
				] = config.NavItem{

					Title: DisplayName(
						page.Collection,
					),

					URL: "/" +
						page.Collection +
						"/",
				}
			}

			continue
		}

		if page.Slug == "index" {
			continue
		}

		contentPages[
			page.Slug,
		] = config.NavItem{

			Title: page.Title,

			URL: page.URL,
		}
	}

	items := []config.NavItem{
		{
			Title: "Home",
			URL: "/",
		},
	}

	var pageNames []string

	for slug := range contentPages {

		pageNames = append(
			pageNames,
			slug,
		)
	}

	sort.Strings(
		pageNames,
	)

	for _, slug := range pageNames {

		items = append(
			items,
			contentPages[slug],
		)
	}

	var collectionNames []string

	for name := range collections {

		collectionNames = append(
			collectionNames,
			name,
		)
	}

	sort.Strings(
		collectionNames,
	)

	for _, name := range collectionNames {

		items = append(
			items,
			collections[name],
		)
	}

	if seenTags {

		items = append(
			items,
			config.NavItem{
				Title: "Tags",
				URL: "/tags/",
			},
		)
	}

	if len(index.Pages) > 0 {

		items = append(
			items,
			config.NavItem{
				Title: "RSS",
				URL: "/feed.xml",
			},
		)
	}

	return items
}