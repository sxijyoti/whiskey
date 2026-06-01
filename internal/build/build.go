package build

import (
	htmltemplate "html/template"
	"path/filepath"
	"os"
	"strings"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/template"
	
)

func BuildPage(input, output string) error {
	raw,err := os.ReadFile(input)

	if err != nil {
		return err
	}

	doc,err := parser.ParseFrontmatter(raw)
	if err != nil {
		return err
	}
	
	html,err := parser.MdToHTML(doc.Body)
	if err != nil {
		return err
	}

	page,err := template.RenderPage(template.PageData{
		Title: doc.Meta.Title,
		Content: htmltemplate.HTML(html),
	})
	if err != nil {
		return err
	}	

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}

	return os.WriteFile(output, page, 0644)

}

func BuildSite() error {
	pages, err := DiscoverPages("site")
	if err != nil {
		return err
	}

	for _, page := range pages {
		rel, err := filepath.Rel("site", page)
		if err != nil {
			return err
		}

		out := strings.TrimSuffix(
			rel,
			filepath.Ext(rel),
		) + ".html"

		output := filepath.Join(
			"dist",
			out,
		)

		if err := BuildPage(page, output); err != nil {
			return err
		}
	}

	return nil
}