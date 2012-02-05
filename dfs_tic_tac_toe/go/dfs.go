package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var nCPU = flag.Int("cpu", 1, "number of CPUs to use")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*nCPU) // Set number of OS threads to use.

	// Build the game graph.
	t := time.Now().UnixNano()
	g := MakeGraph()
	t = time.Now().UnixNano() - t
	fmt.Printf("make graph: %dms\n", t/1e6)
	fmt.Printf("states/edges: %d/%d\n", g.Nodes, g.Edges)

	// Count X wins.
	t = time.Now().UnixNano()
	CountWins(g)
	t = time.Now().UnixNano() - t
	fmt.Printf("count: %dms\n", t/1e6)
	fmt.Printf("x wins: %d\n", g.Wins)
}

// CountWins count X wins and stores the result in g.Wins.
func CountWins(g *Graph) {
	// Calculate how many parallel task splits to do (2^split splits).
	split := 0
	if *nCPU > 1 {
		split = math.Ilogb(float64(*nCPU)) + 4
	}
	// Create a WaitGroup to wait for completion
	// of processing of all nodes.
	var wg sync.WaitGroup
	g.Vis[0] = 1
	visitNode(g, 0, split, &wg)
	// Wait for completion.
	wg.Wait()
}

// visitNode visits a single node n (g.Vis[n] is already set to 1).
func visitNode(g *Graph, n Index, split int, wg *sync.WaitGroup) {
	// Check as to whether the current state is a win for X or O.
	if win := g.IsWin(n); win != E {
		if win == X {
			atomic.AddUint32(&g.Wins, 1)
		}
		return
	}

	// Check as to whether we need to do serial processing...
	if split == 0 {
		// ... if so visit all not visited adjacent nodes.
		for i := Index(0); i < g.Nodes; i++ {
			if needVisit(g, n, i) {
				visitNode(g, i, split, wg)
			}
		}
		return
	}

	// If we need to do a parallel split,
	// then collect all not visited adjacent nodes...
	visit := collectUnvisited(g, n)
	// ... and split them into 2 parallel tasks.
	visitSet(g, visit, split, wg)
}

// needVisit checks whether node i must be visited from node n.
// If it returns true, it marks the node i as visited.
func needVisit(g *Graph, n, i Index) bool {
	// Check that the node is adjacent and is not yet visited.
	if !g.Adj[n][i] || g.Vis[i] != 0 {
		return false
	}
	// Try to mark the node as visited.
	// If it fails then we lose a race with a concurrent goroutine.
	if !atomic.CompareAndSwapInt32(&g.Vis[i], 0, 1) {
		return false
	}
	return true
}

// visitSet visits a set of nodes (potentially in parallel).
func visitSet(g *Graph, visit []Index, split int, wg *sync.WaitGroup) {
	if split == 0 || len(visit) <= 1 {
		// Visit all serially.
		for i := 0; i < len(visit); i++ {
			visitNode(g, visit[i], split, wg)
		}
		return
	}

	// Split into 2 equal parallel tasks.
	wg.Add(2)
	mid := len(visit) / 2
	go func() {
		visitSet(g, visit[:mid], split-1, wg)
		wg.Done()
	}()
	go func() {
		visitSet(g, visit[mid:], split-1, wg)
		wg.Done()
	}()
}

// collectUnvisited collects all not visited adjacent to n nodes
// simultaneously marking them as visited.
func collectUnvisited(g *Graph, n Index) []Index {
	visit := make([]Index, 0, g.Nodes)
	for i := Index(0); i < g.Nodes; i++ {
		if needVisit(g, n, i) {
			visit = append(visit, i)
		}
	}
	return visit
}
