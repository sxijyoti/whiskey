package build

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/parser"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Date  string `xml:"pubDate"`
}

func BuildRSS(
	root string,
	cfg *config.Config,
	pages []string,
) error {

	contentRoot := filepath.Join(
		root,
		"content",
	)

	allowed := map[string]bool{}

	for _, collection := range cfg.RSS.Collections {

		allowed[collection] = true
	}

	var items []Item

	for _, page := range pages {

		raw, err := os.ReadFile(
			page,
		)

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

		if !allowed[
			doc.Meta.Collection,
		] {
			continue
		}

		rel, err := filepath.Rel(
			contentRoot,
			page,
		)

		if err != nil {
			return err
		}

		slug := strings.TrimSuffix(
			rel,
			filepath.Ext(rel),
		)

		url := strings.TrimRight(
			cfg.BaseURL,
			"/",
		) + "/" + slug + "/"

		if slug == "index" {

			url = strings.TrimRight(
				cfg.BaseURL,
				"/",
			) + "/"
		}

		items = append(
			items,
			Item{
				Title: doc.Meta.Title,
				Link:  url,
				Date: doc.Meta.Date.Format(
					"Mon, 02 Jan 2006 15:04:05 MST",
				),
			},
		)
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       cfg.Title,
			Description: cfg.Description,
			Link:        cfg.BaseURL,
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(
		rss,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	data = append(
		[]byte(xml.Header),
		data...,
	)

	return os.WriteFile(
		"dist/feed.xml",
		data,
		0644,
	)
}