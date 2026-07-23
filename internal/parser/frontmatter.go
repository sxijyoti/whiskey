package parser

import (
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

type Frontmatter struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Date        time.Time `yaml:"date"`
	Tags        []string  `yaml:"tags"`
	Collection  string    `yaml:"collection"`
	Layout      string    `yaml:"layout"`
	Draft       bool      `yaml:"draft"`
	Unlisted    bool      `yaml:"unlisted"`
	Hidden      bool      `yaml:"hidden"`
}

func (f Frontmatter) IsUnlisted() bool {
	return f.Unlisted || f.Hidden
}

type Document struct {
	Meta Frontmatter
	Body string
}

func ParseFrontmatter(content []byte) (*Document, error) {
	txt := string(content)

	if !strings.HasPrefix(txt, "---\n") {
		return &Document{
			Body: txt,
		}, nil
	}

	parts := strings.SplitN(txt, "---", 3)

	if len(parts) < 3 {
		return &Document{
			Body: txt,
		}, nil
	}

	var meta Frontmatter

	if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
		return nil, err
	}

	return &Document{
		Meta: meta,
		Body: strings.TrimSpace(parts[2]),
	}, nil
}
