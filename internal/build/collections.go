package build

import (
	"html"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/template"
)

type CollectionPage struct {
	Title string
	URL   string
}

func BuildCollections(
	root string,
	cfg *config.Config,
	pages []string,
) error {

	collections := map[string][]CollectionPage{}

	contentRoot := filepath.Join(
		root,
		"content",
	)

	for _, page := range pages {

		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return err
		}

		if doc.Meta.Draft {
			continue
		}

		if doc.Meta.Collection == "" {
			continue
		}

		rel, err := filepath.Rel(
			contentRoot,
			page,
		)
		if err != nil {
			return err
		}

		slug := rel[:len(rel)-3]

		collections[doc.Meta.Collection] = append(
			collections[doc.Meta.Collection],
			CollectionPage{
				Title: doc.Meta.Title,
				URL:   "/" + slug + "/",
			},
		)
	}

	for name, entries := range collections {
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Title == entries[j].Title {
				return entries[i].URL < entries[j].URL
			}

			return entries[i].Title < entries[j].Title
		})

		var content strings.Builder
		content.WriteString("<h1>")
		content.WriteString(html.EscapeString(name))
		content.WriteString("</h1>\n")

		for _, p := range entries {

			content.WriteString("<p><a href=\"")
			content.WriteString(html.EscapeString(p.URL))
			content.WriteString("\">")
			content.WriteString(html.EscapeString(p.Title))
			content.WriteString("</a></p>\n")
		}

		page, err := template.RenderPage(
			root,
			cfg.Theme,
			"page",
			template.PageData{
				Site:        cfg,
				Title:       name,
				Description: name,
				Content:     htmltemplate.HTML(content.String()),
			},
		)
		if err != nil {
			return err
		}

		out := filepath.Join(
			"dist",
			name,
			"index.html",
		)

		if err := os.MkdirAll(
			filepath.Dir(out),
			0755,
		); err != nil {
			return err
		}

		if err := os.WriteFile(
			out,
			page,
			0644,
		); err != nil {
			return err
		}
	}

	return nil
}
