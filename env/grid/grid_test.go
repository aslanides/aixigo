package grid

import (
	"aixigo/x"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

var spec [][]int
var grid *Gridworld

func init() {
	spec = [][]int{
		{0, 0, 1, 1, 1},
		{1, 0, 0, 1, 2},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0},
	}
	grid = New(spec)
}

func TestConnection(t *testing.T) {
	m := grid.Tiles[0][0]
	n := m.GetNeighbor(4)
	assert.True(t, n != nil)
	assert.Equal(t, m, n)
	n = m.GetNeighbor(1)
	m = grid.Tiles[1][0]
	assert.True(t, n != nil)
	assert.Equal(t, m, n)
}

func TestMovement(t *testing.T) {
	var r x.Reward
	_, r = grid.Perform(0) // left fails
	assert.Equal(t, wallPenalty, r)
	_, r = grid.Perform(1) // right succeeds
	assert.Equal(t, emptyReward, r)
	_, r = grid.Perform(2) // up fails
	assert.Equal(t, wallPenalty, r)
	assert.Equal(t, grid.pos.X(), 1)
	assert.Equal(t, grid.pos.Y(), 0)
	_, r = grid.Perform(4) // stay succeeds
	assert.Equal(t, emptyReward, r)
}
