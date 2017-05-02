package ctw

type tree struct {
	history     []symbol
	historySize int
	maxDepth    int
	root        *node
}

func newTree(depth int) *tree {
	return &tree{
		history:  make([]symbol, 0),
		maxDepth: depth,
		root:     newNode(nil),
	}
}

func (t *tree) reset() {
	t.root = newNode(nil) // runtime.GC handles the rest
}

func (t *tree) updateList(symList []symbol) {
	for _, sym := range symList {
		t.update(sym)
	}
}

func (t *tree) update(sym symbol) {
	if len(t.history) >= t.maxDepth {
		t.doUpdate(sym, 0, t.root, update)
	}
	t.history = append(t.history, sym)
	t.historySize++
}

func (t *tree) doUpdate(sym symbol, depth int, n *node, action updateType) {
	if depth == t.maxDepth {
		n.updateLeaf(sym, action)
		return
	}

	//
	// childSym := t.history[t.historySize-depth-1]
	// TODO: up to github.com/surajx/mc-aixi-ctw/CTW/ContextTree.cpp:85
}
