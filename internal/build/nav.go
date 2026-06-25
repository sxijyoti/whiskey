package build

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/parser"
)

func BuildNav(root string, pages []string) ([]config.NavItem, error) {
	contentRoot := filepath.Join(
		root,
		"content",
	)

	contentPages := map[string]config.NavItem{}
	collections := map[string]config.NavItem{}
	seenTags := map[string]struct{}{}

	for _, page := range pages {
		rawBytes, err := os.ReadFile(page)
		if err != nil {
			return nil, err
		}

		raw, err := parser.ParseFrontmatter(rawBytes)
		if err != nil {
			return nil, err
		}

		if raw.Meta.Draft {
			continue
		}

		rel, err := filepath.Rel(
			contentRoot,
			page,
		)
		if err != nil {
			return nil, err
		}

		slug := strings.TrimSuffix(
			rel,
			filepath.Ext(rel),
		)

		for _, tag := range raw.Meta.Tags {
			seenTags[tag] = struct{}{}
		}

		if raw.Meta.Collection != "" {
			if _, ok := collections[raw.Meta.Collection]; !ok {
				collections[raw.Meta.Collection] = config.NavItem{
					Title: titleCase(raw.Meta.Collection),
					URL:   "/" + raw.Meta.Collection + "/",
				}
			}

			continue
		}

		if slug == "index" {
			continue
		}

		contentPages[slug] = config.NavItem{
			Title: raw.Meta.Title,
			URL:   "/" + slug + "/",
		}
	}

	items := []config.NavItem{{Title: "Home", URL: "/"}}

	contentTitles := make([]string, 0, len(contentPages))
	for slug := range contentPages {
		contentTitles = append(contentTitles, slug)
	}
	sort.Strings(contentTitles)
	for _, slug := range contentTitles {
		items = append(items, contentPages[slug])
	}

	collectionNames := make([]string, 0, len(collections))
	for name := range collections {
		collectionNames = append(collectionNames, name)
	}
	sort.Strings(collectionNames)
	for _, name := range collectionNames {
		items = append(items, collections[name])
	}

	if len(seenTags) > 0 {
		items = append(items, config.NavItem{Title: "Tags", URL: "/tags/"})
	}

	return items, nil
}

func titleCase(value string) string {
	if value == "" {
		return value
	}

	return strings.ToUpper(value[:1]) + value[1:]
}
