package graph

type NodeType string

const (
	PageNode    NodeType = "page"
	SourceNode  NodeType = "source"
	LayoutNode  NodeType = "layout"
	PartialNode NodeType = "partial"
	AssetNode   NodeType = "asset"
)

type Node struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Graph struct {
	Nodes map[string]*Node `json:"nodes"`
	Edges []Edge           `json:"edges"`

	Out map[string][]string `json:"-"`
	In  map[string][]string `json:"-"`
}

func New() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Out:   make(map[string][]string),
		In:    make(map[string][]string),
	}
}

func (g *Graph) AddNode(
	id string,
	typ NodeType,
) {

	if _, exists := g.Nodes[id]; exists {
		return
	}

	g.Nodes[id] = &Node{
		ID:   id,
		Type: typ,
	}
}

func (g *Graph) AddEdge(
	from string,
	to string,
) {

	g.Edges = append(
		g.Edges,
		Edge{
			From: from,
			To:   to,
		},
	)

	g.Out[from] = append(
		g.Out[from],
		to,
	)

	g.In[to] = append(
		g.In[to],
		from,
	)
}

func (g *Graph) Dependencies(
	id string,
) []string {
	return g.Out[id]
}

func (g *Graph) Dependents(
	id string,
) []string {
	return g.In[id]
}

func (g *Graph) ReachableFrom(
	start string,
) []string {

	seen := map[string]bool{}
	queue := []string{start}

	for len(queue) > 0 {

		current := queue[0]
		queue = queue[1:]

		for _, next := range g.In[current] {

			if seen[next] {
				continue
			}

			seen[next] = true

			queue = append(
				queue,
				next,
			)
		}
	}

	var nodes []string

	for node := range seen {
		nodes = append(
			nodes,
			node,
		)
	}

	return nodes
}