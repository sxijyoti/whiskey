package graph

type NodeType string

const (
	PageNode   NodeType = "page"
	SourceNode NodeType = "source"
	TagNode    NodeType = "tag"
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
}

func New() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
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
}

func (g *Graph) Dependencies(
	id string,
) []string {

	var deps []string

	for _, edge := range g.Edges {

		if edge.From == id {
			deps = append(
				deps,
				edge.To,
			)
		}
	}

	return deps
}

func (g *Graph) Dependents(
	id string,
) []string {

	var deps []string

	for _, edge := range g.Edges {

		if edge.To == id {
			deps = append(
				deps,
				edge.From,
			)
		}
	}

	return deps
}