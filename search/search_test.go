package search

import (
	"aixigo/env/grid"
	"aixigo/x"
	"math"
	"math/rand"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

var meta *Meta
var env x.Environment

func init() {
	spec := [][]int{
		{0, 0, 1},
		{1, 0, 0},
		{0, 1, 2},
	}
	meta = &Meta{
		Meta:    grid.Meta,
		Horizon: 6,
		Samples: 1000,
		UCB:     math.Sqrt2,
		Model:   grid.NewModel(spec),
		Utility: func(e x.Percept, dfr int) float64 { return float64(e.R) },
		PRN:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	env = grid.New(spec)
}

func TestSearchDeterministic(t *testing.T) {
	var a x.Action
	a = GetAction(meta)
	assert.Equal(t, x.Action(1), a)
	meta.Model.Perform(a)

	a = GetAction(meta)
	assert.Equal(t, x.Action(3), a)
	meta.Model.Perform(a)

	a = GetAction(meta)
	assert.Equal(t, x.Action(1), a)
	meta.Model.Perform(a)

	a = GetAction(meta)
	assert.Equal(t, x.Action(3), a)
	meta.Model.Perform(a)
}
