package planner

import "github.com/sxijyoti/whiskey/internal/graph"

func DirtyPages(
	g *graph.Graph,
	changed []string,
) []string {

	seen := map[string]struct{}{}

	for _, src := range changed {

		for _, page := range g.Dependents(
			src,
		) {

			seen[page] = struct{}{}
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