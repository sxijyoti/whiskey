package build

import (
	"os"
	"path/filepath"

	"github.com/sxijyoti/whiskey/internal/parser"
)

type CollectionPage struct {
	Title string
	URL   string
}

func BuildCollections(
	root string,
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

		collections[
			doc.Meta.Collection,
		] = append(
			collections[
				doc.Meta.Collection,
			],
			CollectionPage{
				Title: doc.Meta.Title,
				URL: "/" + slug + "/",
			},
		)
	}

	for name, entries := range collections {

		html := "<h1>" + name + "</h1>\n"

		for _, p := range entries {

			html +=
				"<p><a href=\"" +
					p.URL +
					"\">" +
					p.Title +
					"</a></p>\n"
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
			[]byte(html),
			0644,
		); err != nil {
			return err
		}
	}

	return nil
}