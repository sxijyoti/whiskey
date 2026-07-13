package planner

import "github.com/sxijyoti/whiskey/internal/graph"

// DirtyPages returns the set of publishable pages that are transitively
// reachable from any of the changed nodes via reverse edges in the dependency
// graph.
//
// "Publishable" means the node is a PageNode that is not marked as a draft.
// Draft pages are still PageNodes in the graph (for dependency tracking) but
// they are never rendered, so they must not appear in the dirty-page set.
func DirtyPages(
	g *graph.Graph,
	changed []string,
) []string {

	seen := map[string]struct{}{}

	for _, node := range changed {

		// If the changed node itself is a PageNode, it must be rebuilt.
		if n := g.Nodes[node]; n != nil && n.Type == graph.PageNode && !n.Draft {
			seen[node] = struct{}{}
		}

		for _, dep := range g.ReachableFrom(node) {

			n := g.Nodes[dep]

			if n == nil {
				continue
			}

			if n.Type != graph.PageNode {
				continue
			}

			// Draft pages participate in dependency tracking but are never
			// rendered; exclude them from the rebuild set.
			if n.Draft {
				continue
			}

			seen[dep] = struct{}{}
		}
	}

	dirty := make([]string, 0, len(seen))

	for page := range seen {
		dirty = append(dirty, page)
	}

	return dirty
}