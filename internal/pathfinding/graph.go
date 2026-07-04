package pathfinding

type Node struct {
	ID  int64
	Lat float64
	Lon float64
}

type Edge struct {
	To     int64
	Weight float64
}

type Graph struct {
	nodes     map[int64]Node
	adjacency map[int64][]Edge
}

func NewGraph() *Graph {
	return &Graph{
		nodes:     make(map[int64]Node),
		adjacency: make(map[int64][]Edge),
	}
}

func (g *Graph) AddNode(n Node) {
	g.nodes[n.ID] = n
}

func (g *Graph) AddEdge(from, to int64, weightM float64) {
	g.adjacency[from] = append(g.adjacency[from], Edge{To: to, Weight: weightM})
}

func (g *Graph) Node(id int64) (Node, bool) {
	n, ok := g.nodes[id]
	return n, ok
}

func (g *Graph) Neighbors(id int64) []Edge {
	return g.adjacency[id]
}

func (g *Graph) Len() int {
	return len(g.nodes)
}
