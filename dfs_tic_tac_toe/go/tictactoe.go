package tictactoe

const N = 3            // Board size is 3x3.
type State int         // Packed game state.
type Index int         // Node index in the graph.
type Board [N * N]byte // Game board.
const X = byte('x')    // Board cell is occupied by X.
const O = byte('o')    // Board cell is occupied by O.
const E = byte(0)      // Board cell is empty.

// Graph represents 
type Graph struct {
	Nodes      Index           // Number of nodes.
	Edges      int             // Number of edges.
	Wins       uint32          // Number of X wins.
	Adj        [][]bool        // Adjacency matrix.
	Vis        []int32         // Is a node visited?
	stateToIdx map[State]Index // Maps packed game state to node index.
	idxToState []State         // Maps node index to packed game state.
	indexSeq   Index           // Last assigned node index.
}

// MakeGraph creates a graph with nodes representing all possible game states
// and edges representing legal transitions (legal player moves).
func MakeGraph() *Graph {
	var g Graph
	g.Nodes = countGameStates(N*N, 0, 0)
	g.Vis = make([]int32, g.Nodes)
	g.Adj = make([][]bool, g.Nodes)
	g.stateToIdx = make(map[State]Index)
	g.idxToState = make([]State, 1)
	for i := Index(0); i < g.Nodes; i++ {
		g.Adj[i] = make([]bool, g.Nodes)
	}
	g.generateStates(0, 0, Board{})
	return &g
}

// countGameStates counts total number of possible game states.
func countGameStates(n, x, o int) Index {
	if n == 0 {
		if x == o || x == o+1 {
			return 1
		}
		return 0
	}
	c := countGameStates(n-1, x, o)
	c += countGameStates(n-1, x+1, o)
	c += countGameStates(n-1, x, o+1)
	return c
}

// generateStates assigns indices to game states and adds edges to the graph.
func (g *Graph) generateStates(idx0 Index, played int, b Board) {
	t := O
	if played%2 == 0 {
		t = X
	}
	for i := 0; i < N*N; i++ {
		if b[i] != E {
			continue
		}
		b[i] = t
		idx := g.stateIdx(b)
		g.Adj[idx0][idx] = true
		g.Edges++
		if played != N*N-1 {
			g.generateStates(idx, played+1, b)
		}
		b[i] = E
	}
}

// stateIdx returns node index for the game state b.
func (g *Graph) stateIdx(b Board) Index {
	state := State(0)
	for i, token := range b {
		if token == 0 {
			continue
		}
		v := 1 << (uint32(i) * 2)
		if token == 'x' {
			v <<= 1
		}
		state += State(v)
	}
	idx := g.stateToIdx[state]
	if idx != 0 {
		return idx
	}
	g.indexSeq++
	g.stateToIdx[state] = g.indexSeq
	g.idxToState = append(g.idxToState, state)
	return g.indexSeq
}

// IsWin checks whether node idx represents a win for X or O,
// otherwise returns E.
func (g *Graph) IsWin(idx Index) byte {
	state := g.idxToState[idx]
	b := Board{}
	for i := 0; i < N*N; i++ {
		t := state & 3
		if t == 1 {
			b[i] = O
		} else if t == 2 {
			b[i] = X
		}
		state >>= 2
	}
	if isWinFor(b, X) {
		return X
	}
	if isWinFor(b, O) {
		return O
	}
	return E
}

// isWinFor checks whether game state b is a win for t (either X or O).
func isWinFor(b Board, t byte) bool {
	// three across top row
	if b[0] == t && b[1] == t && b[2] == t {
		return true
	}
	// three across middle row
	if b[3] == t && b[4] == t && b[5] == t {
		return true
	}
	// three across bottom row
	if b[6] == t && b[7] == t && b[8] == t {
		return true
	}
	// three vertical on left
	if b[0] == t && b[3] == t && b[6] == t {
		return true
	}
	// three vertical in middle
	if b[1] == t && b[4] == t && b[7] == t {
		return true
	}
	// three vertical on right
	if b[2] == t && b[5] == t && b[8] == t {
		return true
	}
	// three on main diagonal
	if b[0] == t && b[4] == t && b[8] == t {
		return true
	}
	// three on off diagonal
	if b[2] == t && b[5] == t && b[6] == t {
		return true
	}
	return false
}
