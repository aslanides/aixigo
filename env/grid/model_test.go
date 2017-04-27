package grid

import (
	"aixigo/x"
	"fmt"
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
	var e x.Percept
	model := NewModel(spec)
	e = model.Perform(x.Action(4))
	assert.Equal(t, model.pos.X(), 0)
	assert.Equal(t, model.pos.X(), 0)
	assert.Equal(t, int(e.R), 0)
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
	e = model.Perform(x.Action(2))
	assert.Equal(t, int(e.R), 10)
	assert.Equal(t, model.pos.X(), 4)
	assert.Equal(t, model.pos.Y(), 1)
	fmt.Println(model.pos)
	fmt.Println(e)

}
