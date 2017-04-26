package search

import (
	"aixigo/model"
	"aixigo/x"
)

type searchMeta struct {
	x.Meta
	Horizon int
	Model   model.Model
}

//Node seems useful
type Node interface {
	addChild(v interface{})
	getChild(v interface{}) (Node, bool)
	Sample(dfr int) x.Reward
}

type searchNode struct {
	visits float64
	mean   float64
	meta   *searchMeta
}

func newSearchNode(meta *searchMeta) searchNode {
	return searchNode{
		visits: 0.0,
		mean:   0.0,
		meta:   meta,
	}
}

//ChanceNode like the rapper
type ChanceNode struct {
	searchNode
	children map[x.Reward]*DecisionNode // TODO: yuck
	action   x.Action
}

func newChanceNode(a x.Action, meta *searchMeta) *ChanceNode {
	return &ChanceNode{
		searchNode: newSearchNode(meta),
		children:   make(map[x.Reward]*DecisionNode),
		action:     a,
	}
}

func (cn *ChanceNode) getKey(e x.Percept) x.Reward {
	return x.Reward(e.O.ToInt())*cn.meta.MaxReward + e.R // TODO yuckyuck
}

func (cn *ChanceNode) addChild(e x.Percept) {
	key := cn.getKey(e)
	cn.children[key] = newDecisionNode(e, cn.meta)
}

func (cn *ChanceNode) getChild(e x.Percept) (*DecisionNode, bool) {
	key := cn.getKey(e)
	child, found := cn.children[key]
	return child, found
}

//Sample dat (all on the floor)
func (cn *ChanceNode) Sample(dfr int) x.Reward {
	r := x.Reward(0)
	if dfr > cn.meta.Horizon {
		return r
	}
	a := cn.action
	e := cn.meta.Model.Perform(a)
	cn.meta.Model.Update(a, e)
	if _, found := cn.getChild(e); !found {
		cn.addChild(e)
	}
	child, _ := cn.getChild(e)
	r = e.R + child.Sample(dfr+1)
	cn.mean = (1.0/cn.visits + 1.0) * (float64(r) + cn.visits*cn.mean)
	cn.visits += 1.0
	return r
}

func rollOut(meta *searchMeta, dfr int) {
	return // TODO
}

//DecisionNode like dis
type DecisionNode struct {
	searchNode
	children []*ChanceNode
}

func (dn *DecisionNode) addChild(a x.Action) {
	dn.children[a] = newChanceNode(a, dn.meta)
}

func (dn *DecisionNode) getChild(a x.Action) *ChanceNode {
	return dn.children[a]
}

func (dn *DecisionNode) selectAction() x.Action {
	return x.Action(0)
}

// Sample tHat
func (dn *DecisionNode) Sample(dfr int) x.Reward {
	r := x.Reward(0)
	if dfr > dn.meta.Horizon {
		return r
	}
	if dn.visits == 0.0 {
		rollOut(dn.meta, dfr)
	} else {
		a := dn.selectAction()
		r = dn.getChild(a).Sample(dfr)
	}
	dn.mean = (1.0/dn.visits + 1.0) * (float64(r) + dn.visits*dn.mean)
	dn.visits += 1.0
	return r
}

func newDecisionNode(e x.Percept, meta *searchMeta) *DecisionNode {
	return &DecisionNode{
		searchNode: newSearchNode(meta),
		children:   make([]*ChanceNode, meta.NumActions, meta.NumActions),
	}
}
