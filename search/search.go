package search

import (
	"aixigo/x"
	"math"
	"math/rand"
	"runtime"
)

// Meta holds search metadata
type Meta struct {
	x.Meta             // Base Meta struct, holding environment metadata
	Horizon int        // MCTS planning horizon
	Samples int        // MCTS samples
	UCB     float64    // MCTS exploration constant
	Model   x.Model    // Environment model
	Utility x.Utility  // Agent utility function
	PRN     *rand.Rand // Pseudorandom number generator for sampling
}

//GetAction does serial MCTS
func GetAction(meta *Meta) x.Action {
	root := mcts(meta)
	return bestAction(meta, root)
}

// GetActionParallel does parallel MCTS (root parallelism)
func GetActionParallel(meta *Meta) x.Action {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	samplesPerCPU := meta.Samples / cpu
	ch := make(chan *decisionNode, cpu)
	for i := 0; i < cpu; i++ {
		m := &Meta{}
		*m = *meta
		m.Samples = samplesPerCPU
		m.Model = meta.Model.Copy()
		m.PRN = x.NewPRN()

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

	return x.Action(x.ArgMax(totals))
}

func mcts(meta *Meta) *decisionNode {
	model := meta.Model
	root := newDecisionNode(meta)
	model.SaveCheckpoint()
	for i := 0; i < meta.Samples; i++ {
		root.sample(0)
		model.LoadCheckpoint()
	}
	return root
}

func bestAction(meta *Meta, root *decisionNode) x.Action {
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

func (cn *chanceNode) sample(dfr int) float64 {
	R := 0.0
	if dfr == cn.meta.Horizon {
		return R
	}
	a := cn.action
	o, r := cn.meta.Model.Perform(a)
	cn.meta.Model.Update(a, o, r)
	if _, found := cn.getChild(o, r); !found {
		cn.addChild(o, r)
	}
	child, _ := cn.getChild(o, r)
	R = cn.meta.Utility(o, r, dfr) + child.sample(dfr+1)
	cn.mean = (1.0 / (cn.visits + 1.0)) * (R + cn.visits*cn.mean)
	cn.visits += 1.0
	return R
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

func (dn *decisionNode) selectAction(dfr int) x.Action {
	var a x.Action
	if dn.nChildren != int(dn.meta.NumActions) {
		a = x.Action(dn.U[dn.nChildren])
		dn.addChild(a)
		dn.nChildren++
	} else {
		max := math.Inf(-1)
		for b := x.Action(0); b < dn.meta.NumActions; b++ {
			child := dn.getChild(b)
			normalization := float64(dn.meta.Horizon-dfr) * dn.meta.RewardRange
			value := child.mean/normalization + dn.meta.UCB*math.Sqrt(math.Log2(dn.visits)/child.visits)
			if value > max {
				max = value
				a = b
			}
		}
	}
	return a
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
