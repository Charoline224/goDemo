// 核心思想是 对所有边进行松弛n-1次操作（n为节点数量），从而求得目标最短路
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	from, to, val int
}

type Grid struct {
	n     int
	edges []Edge
}

type ShortestPathAlgo interface {
	Run(grid Grid, start int) int
}
type Bellmanford struct{}

func (b Bellmanford) Run(grid Grid, start int) int {
	n := grid.n
	minDist := make([]int, n+1)
	for i := range minDist {
		minDist[i] = math.MaxInt32
	}
	minDist[1] = 0
	for i := 1; i < n; i++ {
		updated := false
		for _, edge := range grid.edges {
			if minDist[edge.from] != math.MaxInt32 &&
				minDist[edge.to] > minDist[edge.from]+edge.val {
				minDist[edge.to] = minDist[edge.from] + edge.val
				updated = true
			}
		}
		if !updated {
			break
		}
	}
	return minDist[n]
}
func main() {
	in := bufio.NewReader(os.Stdin)

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([]Edge, m)
	grid := Grid{
		n:     n,
		edges: edges,
	}
	for i := 0; i < m; i++ {
		var p1, p2, val int
		fmt.Fscan(in, &p1, &p2, &val)
		grid.edges[i] = Edge{from: p1, to: p2, val: val}
	}

	start := 1

	var algo ShortestPathAlgo = Bellmanford{}
	dist := algo.Run(grid, start)

	if dist == math.MaxInt32 {
		fmt.Println("unconnected")
	} else {
		fmt.Println(dist)
	}
}
