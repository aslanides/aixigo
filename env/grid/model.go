package grid

import "aixigo/x"

// TODO generalize to mixtures

//Model lol
type Model struct {
	Gridworld
	savedPos tile
}

//NewModel does the things
func NewModel(spec [][]int) *Model {
	m := New(spec)
	return &Model{
		Gridworld: *m,
		savedPos:  m.Tiles[0][0],
	}
}

//Update method
func (model *Model) Update(a x.Action, e x.Percept) {}

//SaveCheckpoint ...
func (model *Model) SaveCheckpoint() {
	model.savedPos = model.pos
}

//LoadCheckpoint ...
func (model *Model) LoadCheckpoint() {
	model.pos = model.savedPos
}

//ConditionalDistribution kek
func (model *Model) ConditionalDistribution(e x.Percept) float64 {
	return 1.0
}

//Copy ...
func (model *Model) Copy() x.Model {
	newModel := &Model{}
	*newModel = *model
	return newModel
}
