package fingerprint

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sort"

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
			var (
				data []byte
				err  error
			)

			if strings.HasPrefix(
				node.ID,
				"local:",
			) {

				path := filepath.Join(
					root,
					"content",
					strings.TrimPrefix(
						node.ID,
						"local:",
					),
				)

				data, err = os.ReadFile(
					path,
				)

			} else {

				data, err = source.ReadWorkspace(
					root,
					node.ID,
				)
			}

			if err != nil {
				continue
			}

			hash := SHA256(
				data,
			)

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

func PageSetHash(
	pages []string,
) string {

	pages = append([]string(nil), pages...)

	sort.Strings(
		pages,
	)

	return SHA256(
		[]byte(
			strings.Join(
				pages,
				"\n",
			),
		),
	)
}

func GraphHash(
	g *graph.Graph,
) string {

	var items []string

	for id := range g.Nodes {
		items = append(
			items,
			"N:"+id,
		)
	}

	for _, edge := range g.Edges {

		items = append(
			items,
			"E:"+edge.From+"->"+edge.To,
		)
	}

	sort.Strings(
		items,
	)

	return SHA256(
		[]byte(
			strings.Join(
				items,
				"\n",
			),
		),
	)
}