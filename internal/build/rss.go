package build

import (
	"encoding/xml"
	"os"
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

	return os.WriteFile(
		"dist/feed.xml",
		data,
		0644,
	)
}