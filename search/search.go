package search

import (
	"aixigo/x"
	"math"
	"math/rand"
)

//Meta holds search metadata
type Meta struct {
	x.Meta
	Horizon int
	Samples int
	UCB     float64
	Model   x.Model
	Utility x.Utility
	PRN     *rand.Rand
}

//GetAction does the things
func GetAction(meta *Meta) x.Action {
	model := meta.Model
	root := newDecisionNode(meta)
	model.SaveCheckpoint()
	for i := 0; i < meta.Samples; i++ {
		root.sample(0)
		model.LoadCheckpoint()
	}
	action := x.Action(-1)
	max := math.Inf(-1)
	for a := x.Action(0); a < meta.NumActions; a++ {
		cn := root.getChild(a)
		if cn.mean > max {
			max = cn.mean
			action = a
		}
	}
	return action
}

type node interface {
	addChild(v interface{})
	getChild(v interface{}) (node, bool)
	sample(dfr int) x.Reward
}

type searchNode struct {
	visits float64
	mean   float64
	meta   *Meta
}

func newSearchNode(meta *Meta) searchNode {
	return searchNode{
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

func (cn *chanceNode) getKey(e x.Percept) x.Reward {
	return x.Reward(e.O.ToInt())*cn.meta.MaxReward + e.R // TODO yuckyuck
}

func (cn *chanceNode) addChild(e x.Percept) {
	key := cn.getKey(e)
	cn.children[key] = newDecisionNode(cn.meta)
}

func (cn *chanceNode) getChild(e x.Percept) (*decisionNode, bool) {
	key := cn.getKey(e)
	child, found := cn.children[key]
	return child, found
}

func (cn *chanceNode) sample(dfr int) float64 {
	r := 0.0
	if dfr == cn.meta.Horizon {
		return r
	}
	a := cn.action
	e := cn.meta.Model.Perform(a)
	cn.meta.Model.Update(a, e)
	if _, found := cn.getChild(e); !found {
		cn.addChild(e)
	}
	child, _ := cn.getChild(e)
	r = cn.meta.Utility(e, dfr) + child.sample(dfr+1)
	cn.mean = (1.0 / (cn.visits + 1.0)) * (r + cn.visits*cn.mean)
	cn.visits += 1.0
	return r
}

func rollOut(meta *Meta, dfr int) float64 {
	r := 0.0
	for i := dfr; i <= meta.Horizon; i++ {
		a := x.Action(meta.PRN.Intn(int(meta.NumActions)))
		e := meta.Model.Perform(a)
		meta.Model.Update(a, e)
		r += meta.Utility(e, i)
	}
	return r
}

type decisionNode struct {
	searchNode
	children  []*chanceNode
	U         []int
	nChildren int
}

func newDecisionNode(meta *Meta) *decisionNode {
	return &decisionNode{
		searchNode: newSearchNode(meta),
		children:   make([]*chanceNode, meta.NumActions, meta.NumActions),
		U:          meta.PRN.Perm(int(meta.NumActions)),
		nChildren:  0,
	}
}

func (dn *decisionNode) addChild(a x.Action) {
	dn.children[a] = newChanceNode(a, dn.meta)
}

func (dn *decisionNode) getChild(a x.Action) *chanceNode {
	return dn.children[a]
}

func (dn *decisionNode) selectAction() x.Action {
	var a x.Action
	if dn.nChildren != int(dn.meta.NumActions) {
		a = x.Action(dn.U[dn.nChildren])
		dn.addChild(a)
		dn.nChildren++
	} else {
		max := math.Inf(-1)
		for b := x.Action(0); b < dn.meta.NumActions; b++ {
			child := dn.getChild(b)
			normalization := x.Reward(dn.meta.Horizon) * (dn.meta.MaxReward - dn.meta.MinReward)
			value := child.mean/float64(normalization) + dn.meta.UCB*math.Sqrt(math.Log2(dn.visits)/child.visits)
			if value > max {
				max = value
				a = b
			}
		}
	}
	return a
}

func (dn *decisionNode) sample(dfr int) float64 {
	r := 0.0
	if dfr == dn.meta.Horizon {
		return r
	}
	if dn.visits == 0.0 {
		r = rollOut(dn.meta, dfr)
	} else {
		a := dn.selectAction()
		r = dn.getChild(a).sample(dfr)
	}
	dn.mean = (1.0 / (dn.visits + 1.0)) * (r + dn.visits*dn.mean)
	dn.visits += 1.0
	return r
}
