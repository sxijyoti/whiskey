package planner

import "github.com/sxijyoti/whiskey/internal/graph"

func DirtyPages(
	g *graph.Graph,
	changed []string,
) []string {

	seen := map[string]struct{}{}

	for _, node := range changed {

		for _, dep := range g.ReachableFrom(
			node,
		) {

			n := g.Nodes[dep]

			if n == nil {
				continue
			}

			if n.Type != graph.PageNode {
				continue
			}

			seen[dep] = struct{}{}
		}
	}

	var pages []string

	for page := range seen {

		pages = append(
			pages,
			page,
		)
	}

	return pages
}