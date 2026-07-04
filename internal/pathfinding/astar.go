package pathfinding

import (
	"container/heap"
	"fmt"

	"svitlonav/internal/geo"
)

type item struct {
	id       int64
	priority float64
	index    int
}

type priorityQueue []*item

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}
func (pq *priorityQueue) Push(x any) {
	it := x.(*item)
	it.index = len(*pq)
	*pq = append(*pq, it)
}
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.index = -1
	*pq = old[:n-1]
	return it
}

func AStar(g *Graph, start, goal int64) ([]int64, error) {
	startNode, ok := g.Node(start)
	if !ok {
		return nil, fmt.Errorf("astar: start lamp %d not in graph", start)
	}
	goalNode, ok := g.Node(goal)
	if !ok {
		return nil, fmt.Errorf("astar: goal lamp %d not in graph", goal)
	}
	if start == goal {
		return []int64{start}, nil
	}

	heuristic := func(id int64) float64 {
		n, _ := g.Node(id)
		return geo.HaversineMeters(n.Lat, n.Lon, goalNode.Lat, goalNode.Lon)
	}

	gScore := map[int64]float64{start: 0}
	cameFrom := map[int64]int64{}
	visited := map[int64]bool{}

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{id: start, priority: heuristic(start)})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*item)

		if visited[current.id] {
			continue
		}
		visited[current.id] = true

		if current.id == goal {
			return reconstructPath(cameFrom, start, goal), nil
		}

		for _, edge := range g.Neighbors(current.id) {
			if visited[edge.To] {
				continue
			}

			tentativeG := gScore[current.id] + edge.Weight

			if existing, ok := gScore[edge.To]; !ok || tentativeG < existing {
				gScore[edge.To] = tentativeG
				cameFrom[edge.To] = current.id
				heap.Push(pq, &item{id: edge.To, priority: tentativeG + heuristic(edge.To)})
			}
		}
	}

	return nil, fmt.Errorf(
		"astar: no path between lamp %d (%.5f,%.5f) and lamp %d (%.5f,%.5f) — the lamp graph is likely disconnected in this area, try a larger LAMP_RADIUS_METERS in cmd/build-graph",
		start, startNode.Lat, startNode.Lon, goal, goalNode.Lat, goalNode.Lon,
	)
}

func reconstructPath(cameFrom map[int64]int64, start, goal int64) []int64 {
	path := []int64{goal}
	for current := goal; current != start; {
		current = cameFrom[current]
		path = append([]int64{current}, path...)
	}
	return path
}
