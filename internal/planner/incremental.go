package planner

import "github.com/sxijyoti/whiskey/internal/graph"

func IncrementalDirtySet(
	g *graph.Graph,
	local []string,
	changedSources []string,
) []string {

	seen := map[string]struct{}{}

	for _, page := range local {

		seen[page] = struct{}{}
	}

	for _, page := range DirtyPages(
		g,
		changedSources,
	) {

		seen[page] = struct{}{}
	}

	var dirty []string

	for page := range seen {

		dirty = append(
			dirty,
			page,
		)
	}

	return dirty
}