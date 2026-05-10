package main
import (
    "container/heap"
    "fmt"
    "math"
)
type Edge struct{
	to,val int
}

type PriorityQueue []*Edge

func (pq PriorityQueue)Len() int{
	return len(pq)
}
func (pq PriorityQueue)Less(i,j int) bool{
	return pq[i].val < pq[j].val
}
func (pq PriorityQueue)Swap(i,j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue)Push(x interface{}){
	*pq = append(*pq,x.(*Edge))
}
func (pq *PriorityQueue)Pop() interface{}{
	old := *pq
	n := len(old)
	edge := old[n-1]
	*pq = old[0:n-1]
	return edge
}

func dijkstra(n, m int, edges [][]int, start, end int){
	grid := make([][]Edge, n+1)
	for _, edge := range edges{
		p1,p2,val = edge[0],edge[1],edge[2]
		grid[p1] = append(grid[p1], Edge{to:p2,val:val} )
	}
	minDist := make([]int, n+1)
	for i := range minDist{
		minDist[i] = math.MaxInt64
	}
	visited := make([]bool, n+1)
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(*pq,&Edge{to:start,val:0})

	minDist[start] = 0

	for pq.Len()>0{

		edge Edge = heap.Pop(*pq)

		if visted[edge.to] {
			continue
		}
		visited[edge.to] = true
		for _, new := range grid[edge.to]{
			if !visited[new.to] && new.val + minDist[edge.to] < minDist[new.to]{
				minDist[new.to] = new.val + minDist[edge.to]
				heap.Push(*pq, &Edge{to:new.to, val:minDist[new.to]})
			}
		}
	}

}