package template

import (
	"bytes"
	htmltemplate "html/template"
	"github.com/sxijyoti/whiskey/internal/config"
)

type PageData struct {
	Site        *config.Config
	Title       string
	Description string
	Date        string
	Content     htmltemplate.HTML
}

func RenderPage(
	siteRoot string,
	layout string,
	data PageData,
) ([]byte, error) {

	tmpl, err := LoadLayout(
		siteRoot,
		layout,
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)

	return buf.Bytes(), err
}
