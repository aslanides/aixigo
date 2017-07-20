package mcts

import (
	"aixigo/x"
	"math"
	"math/rand"
	"sync"
)

// Meta holds search metadata
type Meta struct {
	x.Meta             // Base Meta struct, holding environment metadata
	Horizon int        // Planning horizon
	Samples int        // Samples
	UCB     float64    // Exploration constant
	Model   x.Model    // Environment model
	Utility x.Utility  // Agent utility function
	PRN     *rand.Rand // Pseudorandom number generator for sampling
	U       []int      // Order in which to explore actions
	// (NOTE: putting this here is a performance hack)
	// search for commit message "hackish" to undo this
}

// NewMeta gives sensible defaults
func NewMeta(meta x.Meta, model x.Model, samples int) *Meta {
	// sensible defaults

	m := Meta{
		Meta:    meta,
		Horizon: 10,
		Samples: samples,
		UCB:     math.Sqrt2,
		Model:   model,
		Utility: x.RLUtility,
		PRN:     x.NewPRN(),
		U:       make([]int, meta.NumActions, meta.NumActions),
	}
	m.U = m.PRN.Perm(int(meta.NumActions))

	return &m
}

type node interface {
	addChild(v interface{})
	getChild(v interface{}) (node, bool)
	sample(dfr int) x.Reward
}

func mcts(meta *Meta) *decisionNode {
	model := meta.Model
	root := newDecisionNode(meta)
	model.Save()
	for i := 0; i < meta.Samples; i++ {
		root.sample(0)
		model.Load()
	}
	return root
}

func bestAction(root *decisionNode) x.Action {
	action := x.Action(-1)
	max := math.Inf(-1)
	for a := x.Action(0); a < root.meta.NumActions; a++ {
		cn := root.getChild(a)
		if cn.mean > max {
			max = cn.mean
			action = a
		}
	}
	return action
}

type searchNode struct {
	mutex  sync.Mutex
	visits float64
	mean   float64
	meta   *Meta
}

func newSearchNode(meta *Meta) searchNode {
	return searchNode{
		mutex:  sync.Mutex{},
		visits: 0.0,
		mean:   0.0,
		meta:   meta,
	}
}

type chanceNode struct {
	searchNode
	children map[x.Reward]*decisionNode // TODO: yuck
	action   x.Action
}

func newChanceNode(a x.Action, meta *Meta) *chanceNode {
	return &chanceNode{
		searchNode: newSearchNode(meta),
		children:   make(map[x.Reward]*decisionNode),
		action:     a,
	}
}

func (cn *chanceNode) getKey(o x.Observation, r x.Reward) x.Reward {
	return x.Reward(o)*cn.meta.MaxReward + r // TODO yuckyuck
}

func (cn *chanceNode) addChild(o x.Observation, r x.Reward) {
	key := cn.getKey(o, r)
	cn.children[key] = newDecisionNode(cn.meta)
}

func (cn *chanceNode) getChild(o x.Observation, r x.Reward) (*decisionNode, bool) {
	key := cn.getKey(o, r)
	child, found := cn.children[key]
	return child, found
}

type decisionNode struct {
	searchNode
	children  []*chanceNode
	nChildren int
}

func newDecisionNode(meta *Meta) *decisionNode {
	return &decisionNode{
		searchNode: newSearchNode(meta),
		children:   make([]*chanceNode, meta.NumActions, meta.NumActions),
		nChildren:  0,
	}
}

func (dn *decisionNode) addChild(a x.Action) {
	dn.children[a] = newChanceNode(a, dn.meta)
}

func (dn *decisionNode) getChild(a x.Action) *chanceNode {
	return dn.children[a]
}

func rollOut(meta *Meta, dfr int) float64 {
	R := 0.0
	for i := dfr; i <= meta.Horizon; i++ {
		a := x.Action(meta.PRN.Intn(int(meta.NumActions)))
		o, r := meta.Model.Perform(a)
		meta.Model.Update(a, o, r)
		R += meta.Utility(o, r, i)
	}
	return R
}
