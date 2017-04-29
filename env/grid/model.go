package grid

import "aixigo/x"

//Model currently is just an environment wrapper (for usage with AImu)
// TODO generalize to mixtures
type Model struct {
	Gridworld
	savedPos tile
}

//NewModel is our constructor
func NewModel(spec [][]int) *Model {
	m := New(spec)
	return &Model{
		Gridworld: *m,
		savedPos:  m.Tiles[0][0],
	}
}

//Update ...
func (model *Model) Update(a x.Action, o x.Observation, r x.Reward) {}

//SaveCheckpoint ...
func (model *Model) SaveCheckpoint() {
	model.savedPos = model.pos
}

//LoadCheckpoint ...
func (model *Model) LoadCheckpoint() {
	model.pos = model.savedPos
}

//ConditionalDistribution ...
func (model *Model) ConditionalDistribution(o x.Observation, r x.Reward) float64 {
	op, rp := model.Perform(x.Action(4)) // hmm, this seems dodgy, maybe use genPercept instead?
	if o == op && r == rp {
		return 1.0
	}
	return 0.0

}

//Copy ...
func (model *Model) Copy() x.Model {
	newModel := &Model{}
	*newModel = *model
	return newModel
}
