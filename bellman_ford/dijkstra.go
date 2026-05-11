package main

import (
	"container/heap"
	"math"
)

// 定义一个由指向边的指针组成的优先队列
type PriorityQueue []*Edge

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].val < pq[j].val
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Edge))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	edge := old[n-1]
	*pq = old[0 : n-1]
	return edge
}

func dijkstra(n, m int, edges [][]int, start, end int) {
	//初始化边矩阵
	grid := make([][]Edge, n+1)
	for _, edge := range edges {
		p1, p2, val := edge[0], edge[1], edge[2]
		grid[p1] = append(grid[p1], Edge{to: p2, val: val})
	}
	//初始化最小距离
	minDist := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt64
	}
	minDist[start] = 0
	//初始化标记数组（节点是否入列）
	visited := make([]bool, n+1)
	//初始化优先队列，使用heap的方法把边放进去
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Edge{to: start, val: 0})
	//
	for pq.Len() > 0 {

		edge := heap.Pop(pq).(*Edge)

		if visited[edge.to] {
			continue
		}
		visited[edge.to] = true
		for _, new := range grid[edge.to] {
			if !visited[new.to] && new.val+minDist[edge.to] < minDist[new.to] {
				minDist[new.to] = new.val + minDist[edge.to]
				heap.Push(pq, &Edge{to: new.to, val: minDist[new.to]})
			}
		}
	}

}
