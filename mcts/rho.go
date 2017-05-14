package mcts

/*
* Implementation of \rhoUCT, a history-based Monte Carlo Tree Search Algorithm
* Original paper: "A Monte Carlo AIXI Approximation",
* 					(Veness, Ng, Hutter, Uther, Silver; 2010)
*
* This file contains the vanilla serial implementation, and a version with
* root parallelism (Chaslot, Winands, & Herik, 2008)
 */

import (
	"aixigo/x"
	"fmt"
	"math"
)

//GetAction does serial MCTS
func GetAction(meta *Meta) x.Action {
	root := mcts(meta)
	return bestAction(root)
}

// GetActionRootParallel does parallel MCTS (root parallelism)
func GetActionRootParallel(meta *Meta) x.Action {
	cpu := x.NumCPU()
	samplesPerCPU := meta.Samples / cpu
	ch := make(chan *decisionNode, cpu)
	for i := 0; i < cpu; i++ {
		m := &Meta{}
		*m = *meta
		m.Samples = samplesPerCPU
		m.Model = meta.Model.Copy()
		m.PRN = x.NewPRN()
		m.U = meta.PRN.Perm(int(meta.NumActions))

		go func(m *Meta) {
			root := mcts(m)
			ch <- root
		}(m)
	}

	totals := make([]float64, int(meta.NumActions), int(meta.NumActions))
	visits := make([]float64, int(meta.NumActions), int(meta.NumActions))
	for i := 0; i < cpu; i++ {
		n := <-ch
		for a := x.Action(0); a < meta.NumActions; a++ {
			child := n.getChild(a)
			totals[int(a)] += child.mean * child.visits
			visits[int(a)] += child.visits
		}
	}
	for i := range totals {
		totals[i] /= visits[i]
	}

	a, _ := x.ArgMax(totals)
	return x.Action(a)
}

// GetActionTreeParallel ...
func GetActionTreeParallel(meta *Meta) x.Action {
	root := newDecisionNode(meta)
	model := meta.Model //.Copy()
	model.Save()
	fmt.Println(meta.Samples)
	for i := 0; i < meta.Samples; i++ {
		root.sample(0)
		model.Load()
	}

	return bestAction(root)
}

func cSample(cn *chanceNode, model x.Model, dfr int) float64 {
	cn.mutex.Lock()
	defer cn.mutex.Unlock()
	a := cn.action
	o, r := model.Perform(a)
	model.Update(a, o, r)
	if _, found := cn.getChild(o, r); !found {
		cn.addChild(o, r)
	}
	child, _ := cn.getChild(o, r)
	R := cn.meta.Utility(o, r, dfr) + dSample(child, model, dfr+1)
	cn.mean = (1.0 / (cn.visits + 1.0)) * (R + cn.visits*cn.mean)
	cn.visits += 1.0
	return R
}

func dSample(dn *decisionNode, model x.Model, dfr int) float64 {
	dn.mutex.Lock()
	defer dn.mutex.Unlock()
	R := 0.0
	if dfr == dn.meta.Horizon {
		return R
	}
	if dn.visits == 0.0 {
		R = rollOut(dn.meta, dfr)
	} else {
		a := dn.selectAction(dfr)
		R = cSample(dn.getChild(a), model, dfr)
	}
	dn.mean = (1.0 / (dn.visits + 1.0)) * (R + dn.visits*dn.mean)
	dn.visits += 1.0
	return R
}

func (cn *chanceNode) sample(dfr int) float64 {
	a := cn.action
	o, r := cn.meta.Model.Perform(a)
	cn.meta.Model.Update(a, o, r)
	if _, found := cn.getChild(o, r); !found {
		cn.addChild(o, r)
	}
	child, _ := cn.getChild(o, r)
	R := cn.meta.Utility(o, r, dfr) + child.sample(dfr+1)
	cn.mean = (1.0 / (cn.visits + 1.0)) * (R + cn.visits*cn.mean)
	cn.visits += 1.0
	return R
}

func (dn *decisionNode) sample(dfr int) float64 {
	R := 0.0
	if dfr == dn.meta.Horizon {
		return R
	}
	if dn.visits == 0.0 {
		R = rollOut(dn.meta, dfr)
	} else {
		a := dn.selectAction(dfr)
		R = dn.getChild(a).sample(dfr)
	}
	dn.mean = (1.0 / (dn.visits + 1.0)) * (R + dn.visits*dn.mean)
	dn.visits += 1.0
	return R
}

func (dn *decisionNode) selectAction(dfr int) x.Action {
	var a x.Action
	if dn.nChildren != int(dn.meta.NumActions) {
		a = x.Action(dn.meta.U[dn.nChildren])
		dn.addChild(a)
		dn.nChildren++
	} else {
		max := math.Inf(-1)
		for b := x.Action(0); b < dn.meta.NumActions; b++ {
			child := dn.getChild(b)
			normalization := float64(dn.meta.Horizon-dfr) * dn.meta.RewardRange
			// Using x.Log2 is a performance hack which means we will underestimate
			// the logarithm
			value := child.mean/normalization + dn.meta.UCB*math.Sqrt(float64(x.Log2(uint(dn.visits)))/child.visits)
			if value > max {
				max = value
				a = b
			}
		}
	}
	return a
}
