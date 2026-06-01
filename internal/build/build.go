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
		Description: doc.Meta.Description,
		Date: doc.Meta.Date.Format("2006-01-02"),
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

func BuildSite(root string) error {
	pages, err := DiscoverPages(root)
	if err != nil {
		return err
	}

	for _, page := range pages {
	rel, err := filepath.Rel(root, page)
	if err != nil {
		return err
	}

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
	}}

	return nil
}