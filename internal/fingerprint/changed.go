package fingerprint

import (
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/source"
)

func ChangedSources(
	g *graph.Graph,
	store Store,
) ([]string, error) {

	var changed []string

	for _, node := range g.Nodes {

		if node.Type != graph.SourceNode {
			continue
		}

		src, err := source.Resolve(
			node.ID,
		)

		if err != nil {
			return nil, err
		}

		meta, err := src.Metadata()

		if err != nil {
			return nil, err
		}

		old := store[node.ID]

		needsFetch := true

		if old.ETag != "" &&
			meta.ETag != "" {

			needsFetch =
				old.ETag != meta.ETag
		}

		if old.LastModified != "" &&
			meta.LastModified != "" {

			needsFetch =
				old.LastModified != meta.LastModified
		}

		if !needsFetch {
			continue
		}

		hash, err := FingerprintSource(
			src,
		)

		if err != nil {
			return nil, err
		}

		if old.Hash != hash {

			changed = append(
				changed,
				node.ID,
			)
		}

		store[node.ID] = Entry{
			Hash: hash,
			ETag: meta.ETag,
			LastModified: meta.LastModified,
		}
	}

	return changed, nil
}