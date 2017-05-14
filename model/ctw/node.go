package ctw

import (
	"log"
	"math"
)

var nodeCount uint // ??

type symbol bool

type updateType bool

const update = updateType(true)
const revert = updateType(false)

type node struct {
	Count           uint
	LogProbWeighted float64
	logProbKT       float64
	symCount        [2]uint
	totalSymCount   uint
	parent          *node
	children        [2]*node
}

func bool2int(sym symbol) int {
	if sym {
		return 1
	}
	return 0
}

func newNode(parent *node) *node {
	return &node{
		Count:           0,
		LogProbWeighted: 0.,
		logProbKT:       0.,
		symCount:        [2]uint{0, 0},
		parent:          parent,
		children:        [2]*node{nil, nil},
	}
}

func (n *node) addChild(sym symbol) {
	n.children[bool2int(sym)] = newNode(n)
}

func (n *node) size() uint {
	var size uint
	for _, child := range n.children {
		if child == nil {
			continue
		}
		size += child.size()
	}

	return size
}

func (n *node) logKTMul(sym symbol) float64 {
	return math.Log(float64(n.symCount[bool2int(sym)])+0.5) - math.Log2(float64(n.totalSymCount)+1.0)
}

func (n *node) updateKT(sym symbol, action updateType) {
	if action == update {
		n.logProbKT += n.logKTMul(sym)
		n.symCount[bool2int(sym)]++
		n.totalSymCount++
	} else { // revert
		if n.symCount[bool2int(sym)] == 0 {
			log.Fatal("Bad news!")
		}
		n.symCount[bool2int(sym)]--
		n.logProbKT -= n.logKTMul(sym)
	}
}

func (n *node) update(sym symbol, action updateType) {
	n.updateKT(sym, action)
	logProbW0, logProbW1 := 0., 0.
	if n.children[0] != nil {
		logProbW0 = n.children[0].LogProbWeighted
	}

	if n.children[1] != nil {
		logProbW0 = n.children[1].LogProbWeighted
	}

	KTRatio := math.Exp(logProbW0 + logProbW1 - n.logProbKT)
	if KTRatio > 1.0 {
		KTRatio = math.Exp(n.logProbKT - logProbW0 - logProbW1)
		n.LogProbWeighted = logProbW0 + logProbW1
	} else {
		n.LogProbWeighted = n.logProbKT
	}

	if math.IsNaN(KTRatio) {
		KTRatio = 0.
	}
	n.LogProbWeighted += math.Log1p(KTRatio) - math.Log(2)
}

func (n *node) updateLeaf(sym symbol, action updateType) {
	n.updateKT(sym, action)
	n.LogProbWeighted = n.logProbKT
}
