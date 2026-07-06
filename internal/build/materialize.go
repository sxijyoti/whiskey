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
) error {

	for _, node := range g.Nodes {

		if node.Type != graph.SourceNode {
			continue
		}

		src, err := source.Resolve(
			node.ID,
		)
		if err != nil {
			return err
		}

		meta, err := src.Metadata()
		if err != nil {
			return err
		}

		entry, exists := manifest.Sources[src.ID()]

		if exists &&
			source.WorkspaceExists(
				root,
				src.ID(),
			) &&
			entry.State["etag"] == meta.ETag &&
			entry.State["last_modified"] == meta.LastModified {

			fmt.Printf(
				"[source] %s (cached)\n",
				src.ID(),
			)

			continue
		}

		result, err := source.Materialize(
			root,
			src,
		)
		if err != nil {
			return err
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

	return source.SaveManifest(
		root,
		manifest,
	)
}