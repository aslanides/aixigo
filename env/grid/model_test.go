package grid

import (
	"aixigo/x"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func init() {
	spec = [][]int{
		{0, 0, 1, 1, 1},
		{1, 0, 0, 1, 2},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0},
	}
}

func TestModel(t *testing.T) {
	var r x.Reward
	model := NewModel(spec)
	_, r = model.Perform(x.Action(4))
	assert.Equal(t, model.pos.X(), 0)
	assert.Equal(t, model.pos.X(), 0)
	assert.Equal(t, int(r), 0)
	model.Perform(x.Action(1))
	model.Perform(x.Action(3))
	model.Perform(x.Action(1))
	model.Perform(x.Action(3))
	model.Perform(x.Action(3))
	model.Perform(x.Action(3))
	model.Perform(x.Action(1))
	model.Perform(x.Action(1))
	model.Perform(x.Action(2))
	model.Perform(x.Action(2))
	_, r = model.Perform(x.Action(2))
	assert.Equal(t, int(r), 10)
	assert.Equal(t, model.pos.X(), 4)
	assert.Equal(t, model.pos.Y(), 1)
}

func TestSaveLoad(t *testing.T) {
	model := NewModel(spec)
	model.Save()
	o, r := model.Perform(x.Action(1))
	assert.Equal(t, 1.0, model.ConditionalDistribution(o, r))
	assert.Equal(t, 1, model.pos.X())
	model.Load()
	assert.Equal(t, 0, model.pos.X())
}

func TestCopy(t *testing.T) {
	model := NewModel(spec)
	model.savedPos = model.Tiles[2][2]
	newModel := &Model{}
	*newModel = *model
	model.pos = model.Tiles[3][2]
	assert.Equal(t, 0, newModel.pos.X())
	newModel.pos = newModel.Tiles[1][1]
	assert.Equal(t, 3, model.pos.X())
	assert.Equal(t, 1, newModel.pos.Y())
	model.savedPos = model.Tiles[4][4]
	assert.Equal(t, 2, newModel.savedPos.X())
}

func TestCopy2(t *testing.T) {
	model := NewModel(spec)
	model.savedPos = model.Tiles[2][2]
	newModel := model.Copy().(*Model) // only necessary for testing
	model.pos = model.Tiles[3][2]
	assert.Equal(t, 0, newModel.pos.X())
	newModel.pos = newModel.Tiles[1][1]
	assert.Equal(t, 3, model.pos.X())
	assert.Equal(t, 1, newModel.pos.Y())
	model.savedPos = model.Tiles[4][4]
	assert.Equal(t, 2, newModel.savedPos.X())
}

func TestMixture(t *testing.T) {
	mix := NewMixture(spec)
	mix.Perform(x.Action(1))

}
