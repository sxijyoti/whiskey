package build

import (
	"fmt"
	"os"

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

		meta, err := src.Metadata()

		if err != nil {

			if source.WorkspaceExists(
				root,
				src.ID(),
			) {

				offlineCached[src.ID()] = true

				fmt.Printf(
					"[source] %s (offline cached)\n",
					src.ID(),
				)

				continue
			}

			failed[src.ID()] = err

			fmt.Printf(
				"[source] %s (metadata failed): %v\n",
				src.ID(),
				err,
			)

			continue
		}

		entry, exists := manifest.Sources[src.ID()]

		if exists &&
			source.WorkspaceExists(
				root,
				src.ID(),
			) &&
			entry.State["etag"] == meta.ETag &&
			entry.State["last_modified"] == meta.LastModified {

			continue
		}

		result, err := source.Materialize(
			root,
			src,
		)

		if err != nil {

			if source.WorkspaceExists(
				root,
				src.ID(),
			) {

				offlineCached[src.ID()] = true

				fmt.Printf(
					"[source] %s (offline cached)\n",
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

		manifest.Sources[src.ID()] = source.ManifestEntry{
			Workspace: result.Workspace,
			ContentHash: result.ContentHash,
			State: map[string]string{
				"etag":          result.Metadata.ETag,
				"last_modified": result.Metadata.LastModified,
			},
		}

		fmt.Printf(
			"[source] %s (updated)\n",
			src.ID(),
		)
	}

	if err := source.SaveManifest(
		root,
		manifest,
	); err != nil {
		return nil, err
	}

	return &MaterializationResult{
		Failed:         failed,
		OfflineCached:  offlineCached,
	}, nil
}

type MaterializationResult struct {
	Failed map[string]error

	// Source could not be refreshed,
	// but an existing workspace copy was reused.
	OfflineCached map[string]bool
}