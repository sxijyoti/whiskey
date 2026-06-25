package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/parser"
)

type TagPage struct {
	Title string
	URL   string
	Date  string
	Tags  []string
}

func BuildTags(
	root string,
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
					URL: url,
					Date: doc.Meta.Date.Format(
						"2006-01-02",
					),
					Tags: doc.Meta.Tags,
				},
			)
		}
	}

	return writeTagPages(
		tagMap,
	)
}

func writeTagPages(
	tags map[string][]TagPage,
) error {

	if err := os.MkdirAll(
		"dist/tags",
		0755,
	); err != nil {
		return err
	}

	var index string

	index += "<h1>Tags</h1>"

	for tag := range tags {

		index +=
			"<p><a href=\"/tags/" +
				tag +
				"/\">" +
				tag +
				"</a></p>"
	}

	if err := os.WriteFile(
		"dist/tags/index.html",
		[]byte(index),
		0644,
	); err != nil {
		return err
	}

	for tag, pages := range tags {

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

		html := "<h1>" + tag + "</h1>"

		for _, p := range pages {

			html +=
				"<p><a href=\"" +
				p.URL +
				"\">" +
				p.Title +
				"</a></p>"
		}

		if err := os.WriteFile(
			filepath.Join(
				dir,
				"index.html",
			),
			[]byte(html),
			0644,
		); err != nil {
			return err
		}
	}

	return nil
}