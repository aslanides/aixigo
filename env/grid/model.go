package grid

import (
	"aixigo/model"
	"aixigo/x"
)

// Model is just an environment wrapper (for usage with AImu)
type Model struct {
	Gridworld
	savedPos tile
}

// NewModel is our constructor
func NewModel(spec [][]int) *Model {
	m := New(spec)
	return &Model{
		Gridworld: *m,
		savedPos:  m.Tiles[0][0],
	}
}

// Update ...
func (model *Model) Update(a x.Action, o x.Observation, r x.Reward) {}

// SaveCheckpoint ...
func (model *Model) SaveCheckpoint() {
	model.savedPos = model.pos
}

// LoadCheckpoint ...
func (model *Model) LoadCheckpoint() {
	model.pos = model.savedPos
}

// ConditionalDistribution ...
func (model *Model) ConditionalDistribution(o x.Observation, r x.Reward) float64 {
	op, rp := model.Perform(x.Action(4)) // hmm, this seems dodgy, maybe use genPercept instead?
	if o == op && r == rp {
		return 1.0
	}
	return 0.0

}

// Copy ...
func (model *Model) Copy() x.Model {
	newModel := &Model{}
	*newModel = *model
	return newModel
}

// NewMixture does what you expect
func NewMixture(spec [][]int) x.Model {
	n := len(spec)
	models := make([]x.Model, 0)
	cpy := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if spec[i][j] != 0 {
				continue
			}
			copy(cpy, spec)
			cpy[i][j] = 2
			models = append(models, NewModel(cpy))
		}
	}
	return model.NewMixture(models)
}
