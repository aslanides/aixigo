package ctw

import "math"

type tree struct {
	history  []symbol
	maxDepth int
	root     *node
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
}

func (t *tree) doUpdate(sym symbol, depth int, n *node, action updateType) {
	if depth == t.maxDepth {
		n.updateLeaf(sym, action)
		return
	}

	childSym := t.history[len(t.history)-depth-1]

	if n.children[bool2int(childSym)] == nil {
		n.addChild(childSym)
	}

	t.doUpdate(sym, depth+1, n.children[bool2int(childSym)], action)

	n.update(sym, action)
	if action == revert && n.totalSymCount == 0 {
		n.parent.children[bool2int(childSym)] = nil
	}
}

func (t *tree) updateHistory(symList []symbol) {
	for _, sym := range symList {
		t.history = append(t.history, sym)
	}
}

func (t *tree) revert() {
	n := len(t.history)
	sym := t.history[n-1]
	t.history = t.history[:(n - 1)]
	t.doUpdate(sym, 0, t.root, revert)
}

func (t *tree) revertHistory(newSize int) {
	t.history = t.history[:newSize]
}

func (t *tree) predict(sym symbol) float64 {
	lpw := t.root.LogProbWeighted
	t.update(sym) // ? O(N) ew
	ret := t.root.LogProbWeighted - lpw
	t.revert()
	return math.Exp2(ret)
}

func (t *tree) predictList(symList []symbol) float64 {
	for _, sym := range symList {
		t.update(sym)
	}
	lpw := t.root.LogProbWeighted
	for _ = range symList {
		t.revert()
	}
	return math.Exp2(lpw - t.root.LogProbWeighted)
}

func (t *tree) predictSymbol() symbol {
	if t.predict(true) > 0.5 {
		return true
	}
	return false
}

func (t *tree) genNextSymbolsAndUpdate(symList []symbol, n int) {
	var sym symbol
	for i := 0; i < n; i++ {
		sym = t.predictSymbol()
		t.update(sym)
		symList[i] = sym
	}
}

func (t *tree) genNextSymbols(symList []symbol, n int) {
	t.genNextSymbolsAndUpdate(symList, n)
	for i := 0; i < n; i++ {
		t.revert()
	}
}
