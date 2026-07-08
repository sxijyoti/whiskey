package build

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
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
	index *SiteIndex,
) error {

	if !cfg.RSS.Enabled {
		return nil
	}

	allowed := map[string]bool{}

	for _, collection := range cfg.RSS.Collections {
		allowed[collection] = true
	}

	var items []Item

	for _, page := range index.Pages {

		if !allowed[page.Collection] {
			continue
		}

		url := strings.TrimRight(
			cfg.BaseURL,
			"/",
		) + page.URL

		items = append(
			items,
			Item{
				Title: page.Title,
				Link:  url,
				Date: page.Date.Format(
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

	output := filepath.Join(
		root,
		"dist",
		"feed.xml",
	)

	if err := os.MkdirAll(
		filepath.Dir(output),
		0755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		output,
		data,
		0644,
	)
}

type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	Location string `xml:"loc"`
	LastMod  string `xml:"lastmod,omitempty"`
}

func BuildSitemap(
	root string,
	cfg *config.Config,
	index *SiteIndex,
) error {

	var sitemap Sitemap

	sitemap.XMLNS =
		"http://www.sitemaps.org/schemas/sitemap/0.9"

	base := strings.TrimRight(
		cfg.BaseURL,
		"/",
	)

	for _, page := range index.Pages {

		url := base + page.URL

		entry := SitemapURL{
			Location: url,
		}

		if !page.Date.IsZero() {

			entry.LastMod = page.Date.Format(
				"2006-01-02",
			)
		}

		sitemap.URLs = append(
			sitemap.URLs,
			entry,
		)
	}

	data, err := xml.MarshalIndent(
		sitemap,
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

	output := filepath.Join(
		root,
		"dist",
		"sitemap.xml",
	)

	if err := os.MkdirAll(
		filepath.Dir(output),
		0755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		output,
		data,
		0644,
	)
}
