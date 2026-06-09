package graph

type NodeType string

const (
	PageNode   NodeType = "page"
	SourceNode NodeType = "source"
	TagNode    NodeType = "tag"
)

type Node struct {
	ID   string
	Type NodeType
}

type Edge struct {
	From string
	To   string
}

type Graph struct {
	Nodes map[string]*Node
	Edges []Edge
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