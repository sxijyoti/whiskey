package graph

import (
	"os"

	"github.com/sxijyoti/whiskey/internal/dependency"
)

func BuildPageGraph(
	path string,
) (*Graph, error) {

	g := New()

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	g.AddNode(
		path,
		PageNode,
	)

	directives := dependency.Extract(
		string(raw),
	)

	for _, d := range directives {

		g.AddNode(
			d.Ref,
			SourceNode,
		)

		g.AddEdge(
			path,
			d.Ref,
		)
	}

	return g, nil
}