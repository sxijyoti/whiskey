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

type TagPage struct {
	Title string
	URL   string
	Date  string
	Tags  []string
}

func BuildTags(
	root string,
	cfg *config.Config,
	pages []string,
) error {

	tagMap := map[string][]TagPage{}

	contentRoot := filepath.Join(
		root,
		"content",
	)

	for _, page := range pages {

		raw, err := os.ReadFile(page)
		if err != nil {
			return err
		}

		doc, err := parser.ParseFrontmatter(
			raw,
		)

		if err != nil {
			return err
		}

		if doc.Meta.Draft {
			continue
		}

		rel, err := filepath.Rel(
			contentRoot,
			page,
		)

		if err != nil {
			return err
		}

		for _, tag := range doc.Meta.Tags {

			slug := strings.TrimSuffix(
				rel,
				filepath.Ext(rel),
			)

			url := "/" + slug + "/"

			if slug == "index" {
				url = "/"
			}

			tagMap[tag] = append(
				tagMap[tag],
				TagPage{
					Title: doc.Meta.Title,
					URL:   url,
					Date: doc.Meta.Date.Format(
						"2006-01-02",
					),
					Tags: doc.Meta.Tags,
				},
			)
		}
	}

	return writeTagPages(
		root,
		cfg,
		tagMap,
	)
}

func writeTagPages(
	root string,
	cfg *config.Config,
	tags map[string][]TagPage,
) error {

	if err := os.MkdirAll(
		"dist/tags",
		0755,
	); err != nil {
		return err
	}

	tagNames := make([]string, 0, len(tags))
	for tag := range tags {
		tagNames = append(tagNames, tag)
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

	for tag, pages := range tags {
		sort.Slice(pages, func(i, j int) bool {
			if pages[i].Date == pages[j].Date {
				return pages[i].Title < pages[j].Title
			}

			return pages[i].Date > pages[j].Date
		})

		dir := filepath.Join(
			"dist",
			"tags",
			tag,
		)

		if err := os.MkdirAll(
			dir,
			0755,
		); err != nil {
			return err
		}

		var content strings.Builder
		content.WriteString("<h1>")
		content.WriteString(html.EscapeString(tag))
		content.WriteString("</h1>")

		for _, p := range pages {

			content.WriteString("<p><a href=\"")
			content.WriteString(html.EscapeString(p.URL))
			content.WriteString("\">")
			content.WriteString(html.EscapeString(p.Title))
			content.WriteString("</a></p>")
		}

		page, err := template.RenderPage(
			root,
			cfg.Theme,
			"page",
			template.PageData{
				Site:        cfg,
				Title:       tag,
				Description: tag,
				Content:     htmltemplate.HTML(content.String()),
			},
		)
		if err != nil {
			return err
		}

		if err := os.WriteFile(
			filepath.Join(
				dir,
				"index.html",
			),
			page,
			0644,
		); err != nil {
			return err
		}
	}

	return nil
}
