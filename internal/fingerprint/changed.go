package fingerprint

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/parser"
	"github.com/sxijyoti/whiskey/internal/source"
)

func ChangedSources(
	root string,
	g *graph.Graph,
	store Store,
) ([]string, error) {

	var changed []string

	for _, node := range g.Nodes {

		switch node.Type {

		case graph.SourceNode:

			data, err := source.ReadWorkspace(
				root,
				node.ID,
			)
			if err != nil {
				continue
			}

			hash := SHA256(data)

			old := store[node.ID]

			if old.Hash != hash {

				changed = append(
					changed,
					node.ID,
				)

				store[node.ID] = Entry{
					Hash: hash,
				}
			}

		case graph.LayoutNode,
			graph.PartialNode,
			graph.AssetNode:

			data, err := os.ReadFile(
				node.ID,
			)
			if err != nil {
				continue
			}

			hash := SHA256(data)

			old := store[node.ID]

			if old.Hash != hash {

				changed = append(
					changed,
					node.ID,
				)

				store[node.ID] = Entry{
					Hash: hash,
				}
			}
		}
	}

	return changed, nil
}

// incase the config changes and rss description
func ConfigHash(
	cfg *config.Config,
) string {

	input :=
		cfg.Theme + "|" +
			cfg.BaseURL + "|" +
			cfg.Title + "|" +
			cfg.Description + "|" +
			strconv.FormatBool(
				cfg.RSS.Enabled,
			)

	for _, c := range cfg.RSS.Collections {
		input += "|" + c
	}

	return SHA256([]byte(input))
}

func UpdateLocalPages(
	contentRoot string,
	store Store,
) error {

	return filepath.WalkDir(
		contentRoot,
		func(
			path string,
			d fs.DirEntry,
			err error,
		) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if strings.HasPrefix(
				filepath.Base(path),
				".",
			) {
				return nil
			}

			if filepath.Ext(path) != ".md" {
				return nil
			}

			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			doc, err := parser.ParseFrontmatter(raw)
			if err == nil && doc.Meta.Draft {
				return nil
			}

			store[path] = Entry{
				Hash: SHA256(raw),
			}

			return nil
		},
	)
}
