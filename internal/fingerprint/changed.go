package fingerprint

import (
	"os"
	"strconv"
	"strings"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/source"
)

// ChangedNodes fingerprints every node in the dependency graph and returns
// the IDs of nodes whose content has changed since the last build.
//
// Node types and how their bytes are obtained:
//
//   - PageNode   — read directly from the filesystem (absolute path = node ID)
//   - SourceNode — local: prefix → read from filesystem;
//     everything else → read from the materialized workspace
//   - LayoutNode, PartialNode, AssetNode — read directly from the filesystem
//
// After this function returns, store contains up-to-date hashes for all nodes
// that were successfully read. Nodes that could not be read are skipped
// (they are not marked as changed unless they previously had a hash in the store).
func ChangedNodes(
	root string,
	g *graph.Graph,
	store Store,
) ([]string, error) {

	var changed []string

	for _, node := range g.Nodes {

		data, err := readNodeBytes(root, node)
		if err != nil {
			// If we previously had a hash for this node and now cannot read it,
			// treat it as changed so dependents are invalidated.
			if old := store[node.ID]; old.Hash != "" {
				changed = append(changed, node.ID)
				delete(store, node.ID)
			}
			continue
		}

		hash := SHA256(data)
		old := store[node.ID]

		if old.Hash != hash {
			changed = append(changed, node.ID)
			store[node.ID] = Entry{Hash: hash}
		}
	}

	return changed, nil
}

// readNodeBytes returns the raw bytes for a graph node.
func readNodeBytes(root string, node *graph.Node) ([]byte, error) {

	switch node.Type {

	case graph.PageNode,
		graph.LayoutNode,
		graph.PartialNode,
		graph.AssetNode:

		// IDs for these node types are absolute filesystem paths.
		return os.ReadFile(node.ID)

	case graph.SourceNode:

		if strings.HasPrefix(node.ID, "local:") {
			path := strings.TrimPrefix(node.ID, "local:")
			return os.ReadFile(path)
		}

		// Remote sources are read from the materialized workspace.
		return source.ReadWorkspace(root, node.ID)
	}

	return nil, nil
}

// ConfigHash produces a stable hash of the configuration fields that, when
// changed, require a full rebuild (theme, URLs, RSS settings, etc.).
func ConfigHash(cfg *config.Config) string {

	input :=
		cfg.Theme + "|" +
			cfg.BaseURL + "|" +
			cfg.Title + "|" +
			cfg.Description + "|" +
			cfg.Favicon + "|" +
			strconv.FormatBool(cfg.RSS.Enabled)

	for _, c := range cfg.RSS.Collections {
		input += "|" + c
	}

	for _, item := range cfg.ExplicitNav {
		input += "|" + item.Title + ":" + item.URL
	}

	return SHA256([]byte(input))
}
