package template

import (
	"bytes"
	"github.com/sxijyoti/whiskey/internal/config"
	htmltemplate "html/template"
)

type PageData struct {
	Site        *config.Config
	Title       string
	Description string
	Date        string
	Tags        []string
	Content     htmltemplate.HTML
}

func RenderPage(
	siteRoot string,
	theme string,
	layout string,
	data PageData,
) ([]byte, error) {

	tmpl, err := LoadLayout(
		siteRoot,
		theme,
		layout,
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)

	return buf.Bytes(), err
}
