package build

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sxijyoti/whiskey/internal/parser"
)

type Page struct {
	Title       string
	Description string

	Slug string
	URL  string

	Layout     string
	Collection string

	Tags []string

	Date      time.Time
	DateLabel string

	Draft bool
}

type SiteIndex struct {
	Pages []Page

	Collections map[string][]Page
	Tags        map[string][]Page
}

func BuildIndex(
	root string,
	pages []string,
) (*SiteIndex, error) {

	contentRoot := filepath.Join(
		root,
		"content",
	)

	index := &SiteIndex{
		Collections: make(map[string][]Page),
		Tags:         make(map[string][]Page),
	}

	for _, page := range pages {

		raw, err := os.ReadFile(page)
		if err != nil {
			return nil, err
		}

		doc, err := parser.ParseFrontmatter(raw)
		if err != nil {
			return nil, err
		}

		if doc.Meta.Draft {
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

		url := "/" + slug + "/"

		if slug == "index" {
			url = "/"
		}

		dateLabel := ""

		if !doc.Meta.Date.IsZero() {
			dateLabel = doc.Meta.Date.Format(
				"2006-01-02",
			)
		}

		p := Page{
			Title:       doc.Meta.Title,
			Description: doc.Meta.Description,

			Slug: slug,
			URL:  url,

			Layout:     doc.Meta.Layout,
			Collection: doc.Meta.Collection,

			Tags: doc.Meta.Tags,

			Date:      doc.Meta.Date,
			DateLabel: dateLabel,

			Draft: doc.Meta.Draft,
		}

		index.Pages = append(
			index.Pages,
			p,
		)

		if p.Collection != "" {
			index.Collections[p.Collection] = append(
				index.Collections[p.Collection],
				p,
			)
		}

		for _, tag := range p.Tags {
			index.Tags[tag] = append(
				index.Tags[tag],
				p,
			)
		}
	}

	sort.Slice(
		index.Pages,
		func(i, j int) bool {
			return index.Pages[i].Date.After(
				index.Pages[j].Date,
			)
		},
	)

	for collection := range index.Collections {

		sort.Slice(
			index.Collections[collection],
			func(i, j int) bool {
				return index.Collections[collection][i].Date.After(
					index.Collections[collection][j].Date,
				)
			},
		)
	}

	for tag := range index.Tags {

		sort.Slice(
			index.Tags[tag],
			func(i, j int) bool {
				return index.Tags[tag][i].Date.After(
					index.Tags[tag][j].Date,
				)
			},
		)
	}

	return index, nil
}