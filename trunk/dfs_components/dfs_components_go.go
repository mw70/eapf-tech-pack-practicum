// Parallel marking of connected components in a graph
// represented by an adjacency matrix implemented by means of depth-first traversal.

package main

import (
	"flag"
	"fmt"
	"math"
	"rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	nVertex = flag.Int("v", 1000, "number of vertices")
	nEdge   = flag.Int("e", 100, "mean number of edges per vertex")
	nCPU    = flag.Int("cpu", 1, "number of CPUs to use")
)

type Graph struct {
	adj  [][]bool // Adjacency matrix.
	comp []uint32 // Component index (0 means not marked).
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*nCPU) // Set number of OS threads to use.

	// Build random graph.
	t := time.Nanoseconds()
	g := MakeGraph(*nVertex, *nEdge)
	t = time.Nanoseconds() - t
	fmt.Printf("make graph: %dms\n", t/1e6)

	// Parallel mark the graph.
	t = time.Nanoseconds()
	g.Mark()
	t = time.Nanoseconds() - t
	fmt.Printf("mark graph: %dms\n", t/1e6)
}

// MakeGraph creates a random graph with v vertices
// and e edges per vertex (on the average).
func MakeGraph(v, e int) *Graph {
	var g Graph
	g.comp = make([]uint32, v)
	g.adj = make([][]bool, v)
	for i := 0; i < v; i++ {
		g.adj[i] = make([]bool, v)
	}
	var wg sync.WaitGroup
	wg.Add(v)
	for i := 0; i < v; i++ {
		go func(i int) {
			r := rand.New(rand.NewSource(int64(i)))
			c := float64(v) / float64(e)
			for j := 0; j < i; j++ {
				if r.Float64()*c < 1 {
					g.adj[i][j] = true
					g.adj[j][i] = true
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return &g
}

// Mark marks connected components in g.
func (g *Graph) Mark() {
	comp := 0 // Component index sequence.
	splitThreshold := 0
	if *nCPU > 1 {
		splitThreshold = math.Ilogb(float64(*nCPU)) + 2
	}
	for i := 0; i < len(g.adj); i++ {
		if g.comp[i] == 0 {
			comp++
			g.comp[i] = uint32(comp)
			var wg sync.WaitGroup
			g.visit(i, comp, splitThreshold, &wg)
			wg.Wait()
		}
	}
}

// visit visits a single vertex n.
func (g *Graph) visit(n, comp, splitThreshold int, wg *sync.WaitGroup) {
	v := len(g.adj)
	if splitThreshold == 0 {
		for i := 0; i < v; i++ {
			if !g.adj[n][i] || g.comp[i] != 0 {
				continue
			}
			if !atomic.CompareAndSwapUint32(&g.comp[i], 0, uint32(comp)) {
				continue
			}
			g.visit(i, comp, splitThreshold, wg)
		}
		return
	}

	visit := g.collectUnvisited(n, comp)
	g.visitSet(visit, comp, splitThreshold, wg)
}

func (g *Graph) visitSet(visit []int, comp, splitThreshold int, wg *sync.WaitGroup) {
	if len(visit) == 0 {
		return
	} else if len(visit) == 1 {
		g.visit(visit[0], comp, splitThreshold, wg)
		return
	} else if splitThreshold == 0 {
		for i := 0; i < len(visit); i++ {
			g.visit(visit[i], comp, splitThreshold, wg)
		}
		return
	}

	wg.Add(2)
	mid := len(visit) / 2
	go func() {
		g.visitSet(visit[:mid], comp, splitThreshold-1, wg)
		wg.Done()
	}()
	go func() {
		g.visitSet(visit[mid:], comp, splitThreshold-1, wg)
		wg.Done()
	}()
}

func (g *Graph) collectUnvisited(n, comp int) []int {
	v := len(g.adj)
	visit := make([]int, 0, v)
	for i := 0; i < v; i++ {
		// Check that the vertex is adjacent and is not yet marked.
		if !g.adj[n][i] || g.comp[i] != 0 {
			continue
		}
		// Try to mark the vertex.
		// If it fails then we lose a race with a concurrent goroutine.
		if !atomic.CompareAndSwapUint32(&g.comp[i], 0, uint32(comp)) {
			continue
		}
		visit = append(visit, i)
	}
	return visit
}

