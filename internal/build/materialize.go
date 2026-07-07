package build

import (
	"fmt"
	"os"
	"strings"

	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/source"
)

func EnsureWorkspace(
	root string,
) error {

	if _, err := source.LoadManifest(
		root,
	); err != nil {
		return err
	}

	return os.MkdirAll(
		source.WorkspaceDir(root),
		0755,
	)
}

func MaterializeSources(
	root string,
	g *graph.Graph,
	manifest *source.Manifest,
) (*MaterializationResult, error) {

	failed := make(
		map[string]error,
	)
	dirty := false

	offlineCached := make(
		map[string]bool,
	)

	for _, node := range g.Nodes {

		if node.Type != graph.SourceNode {
			continue
		}

		src, err := source.Resolve(
			node.ID,
		)

		if err != nil {
			failed[node.ID] = err
			continue
		}

		if strings.HasPrefix(
			src.ID(),
			"local:",
		) {
			continue
		}

		if source.Offline {
			if source.WorkspaceExists(
				root,
				src.ID(),
			) {
				offlineCached[src.ID()] = true

				fmt.Printf(
					"[source] %s (offline)\n",
					src.ID(),
				)

				continue
			}

			failed[src.ID()] = fmt.Errorf(
				"workspace missing",
			)

			continue
		}

		entry, exists := manifest.Sources[src.ID()]

		var oldMeta *source.Metadata

		if exists {

			oldMeta = &source.Metadata{
				ETag:         entry.State["etag"],
				LastModified: entry.State["last_modified"],
			}
		}

		var meta *source.Metadata

		if conditional, ok := src.(source.ConditionalSource); ok {

			meta, err = conditional.ConditionalMetadata(
				oldMeta,
			)

		} else {

			meta, err = src.Metadata()
		}

		if err != nil {

			if source.WorkspaceExists(
				root,
				src.ID(),
			) {

				offlineCached[src.ID()] = true

				fmt.Printf(
					"[source] %s (cached)\n",
					src.ID(),
				)

				continue
			}

			failed[src.ID()] = err

			fmt.Printf(
				"[source] %s (failed): %v\n",
				src.ID(),
				err,
			)

			continue
		}

		if meta.NotModified {

			fmt.Printf(
				"[source] %s (cached)\n",
				src.ID(),
			)

			continue
		}

		result, err := source.Materialize(
			root,
			src,
			meta,
		)
		if err != nil {

			if source.WorkspaceExists(
				root,
				src.ID(),
			) {

				offlineCached[src.ID()] = true

				fmt.Printf(
					"[source] %s (cached)\n",
					src.ID(),
				)

				continue
			}

			failed[src.ID()] = err

			fmt.Printf(
				"[source] %s (failed): %v\n",
				src.ID(),
				err,
			)

			continue
		}

		if exists &&
			entry.ContentHash == result.ContentHash {
			fmt.Printf(
				"[source] %s (unchanged)\n",
				src.ID(),
			)

			continue
		}

		if err := source.WriteWorkspace(
			root,
			src.ID(),
			result.Content,
		); err != nil {
			failed[src.ID()] = err
			fmt.Printf(
				"[source] %s (failed): %v\n",
				src.ID(),
				err,
			)
			continue
		}

		manifest.Sources[src.ID()] = source.ManifestEntry{
			Workspace:   result.Workspace,
			ContentHash: result.ContentHash,
		}
		dirty = true

		fmt.Printf(
			"[source] %s (updated)\n",
			src.ID(),
		)
	}

	if dirty {
		if err := source.SaveManifest(
			root,
			manifest,
		); err != nil {
			return nil, err
		}
	}

	return &MaterializationResult{
		Failed:        failed,
		OfflineCached: offlineCached,
	}, nil
}

type MaterializationResult struct {
	Failed map[string]error

	// Source could not be refreshed,
	// but an existing workspace copy was reused.
	OfflineCached map[string]bool
}
