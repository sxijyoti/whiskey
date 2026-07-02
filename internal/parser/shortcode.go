package parser

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

type Shortcode struct {
	Name string

	Positional []string
	Args       map[string]string
}

type ShortcodeRenderer func(Shortcode) string

var shortcodePattern = regexp.MustCompile(
	`(?s)\{\{<\s*(.*?)\s*>\}\}`,
)
var namedArgPattern = regexp.MustCompile(
	`^[A-Za-z_][A-Za-z0-9_-]*=`,
)

var renderers = map[string]ShortcodeRenderer{}

func init() {

	RegisterShortcode(
		"image",
		renderImage,
	)

	RegisterShortcode(
		"youtube",
		renderYoutube,
	)
}

func RegisterShortcode(
	name string,
	renderer ShortcodeRenderer,
) {

	renderers[name] = renderer
}

func ExpandShortcodes(
	markdown string,
) (string, error) {

	out := shortcodePattern.ReplaceAllStringFunc(
		markdown,
		func(match string) string {

			shortcode := parseShortcode(
				match,
			)

			return renderShortcode(
				shortcode,
				match,
			)
		},
	)

	return out, nil
}

func renderShortcode(
	shortcode Shortcode,
	original string,
) string {

	renderer, ok := renderers[
		shortcode.Name,
	]

	if !ok {
		return original
	}

	return renderer(
		shortcode,
	)
}

func parseShortcode(
	input string,
) Shortcode {

	input = strings.TrimPrefix(
		input,
		"{{<",
	)

	input = strings.TrimSuffix(
		input,
		">}}",
	)

	fields := tokenize(
		strings.TrimSpace(input),
	)

	shortcode := Shortcode{
		Args: make(map[string]string),
	}

	if len(fields) == 0 {
		return shortcode
	}

	shortcode.Name = fields[0]

	for _, field := range fields[1:] {

		if namedArgPattern.MatchString(field) {

			parts := strings.SplitN(
				field,
				"=",
				2,
			)

			key := strings.TrimSpace(parts[0])

			value := strings.TrimSpace(parts[1])

			value = strings.Trim(
				value,
				`"'`,
			)

			shortcode.Args[key] = value

			continue
		}

		shortcode.Positional = append(
			shortcode.Positional,
			strings.Trim(
				field,
				`"'`,
			),
		)
	}

	return shortcode
}

func tokenize(
	input string,
) []string {

	var (
		tokens []string

		current strings.Builder

		inQuote bool

		quote rune
	)

	for _, r := range input {

		switch {

		case inQuote:

			if r == quote {

				inQuote = false

				continue
			}

			current.WriteRune(
				r,
			)

		case r == '"' || r == '\'':

			inQuote = true

			quote = r

		case r == ' ' ||
			r == '\n' ||
			r == '\t' ||
			r == '\r':

			if current.Len() > 0 {

				tokens = append(
					tokens,
					current.String(),
				)

				current.Reset()
			}

		default:

			current.WriteRune(
				r,
			)
		}
	}

	if current.Len() > 0 {

		tokens = append(
			tokens,
			current.String(),
		)
	}

	return tokens
}

func normalizeYoutubeID(
	value string,
) string {

	value = strings.TrimSpace(
		value,
	)

	if value == "" {
		return ""
	}

	if strings.Contains(
		value,
		"youtu.be/",
	) {

		return value[
			strings.LastIndex(
				value,
				"/",
			)+1:
		]
	}

	if strings.Contains(
		value,
		"watch?v=",
	) {

		u, err := url.Parse(
			value,
		)

		if err == nil {

			id := u.Query().Get(
				"v",
			)

			if id != "" {
				return id
			}
		}
	}

	return value
}

func renderImage(
	shortcode Shortcode,
) string {

	src := ""

	if len(shortcode.Positional) > 0 {
		src = shortcode.Positional[0]
	}

	if value, ok := shortcode.Args["src"]; ok {
		src = value
	}

	if src == "" {
		return ""
	}

	width := shortcode.Args["width"]
	caption := shortcode.Args["caption"]

	alt := shortcode.Args["alt"]
	if alt == "" {
		alt = caption
	}

	var out strings.Builder

	class := shortcode.Args["class"]

	out.WriteString(`<figure class="image`)

	if class != "" {

		out.WriteString(` `)
		out.WriteString(html.EscapeString(class))
	}

	out.WriteString(`"`)

	if width != "" {

		out.WriteString(` style="--image-width:`)
		out.WriteString(html.EscapeString(width))
		out.WriteString(`px;"`)
	}

	out.WriteString(`>`)

	out.WriteString(`<img src="`)
	out.WriteString(html.EscapeString(src))
	out.WriteString(`"`)

	out.WriteString(` alt="`)
	out.WriteString(html.EscapeString(alt))
	out.WriteString(`"`)

	out.WriteString(`>`)

	if caption != "" {

		out.WriteString(`<figcaption>`)
		out.WriteString(html.EscapeString(caption))
		out.WriteString(`</figcaption>`)
	}

	out.WriteString(`</figure>`)

	return out.String()
}

func renderYoutube(
	shortcode Shortcode,
) string {

	video := ""

	if len(shortcode.Positional) > 0 {
		video = shortcode.Positional[0]
	}

	if value, ok := shortcode.Args["id"]; ok {
		video = value
	}

	video = normalizeYoutubeID(
		video,
	)

	if video == "" {
		return ""
	}

	return fmt.Sprintf(
		`<div class="youtube"><iframe src="https://www.youtube.com/embed/%s" title="YouTube video" loading="lazy" allowfullscreen></iframe></div>`,
		html.EscapeString(video),
	)
}
