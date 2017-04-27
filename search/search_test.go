package search

import (
	"aixigo/env/grid"
	"aixigo/x"
	"math"
	"testing"

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
		Samples: 10000,
		UCB:     math.Sqrt2,
		Model:   grid.NewModel(spec),
		Utility: func(e x.Percept, dfr int) float64 { return float64(e.R) },
		PRN:     x.NewPRN(),
	}
	env = grid.New(spec)
}

func helperSearchDeterministic(t *testing.T, planner func(*Meta) x.Action) {
	var a x.Action
	a = planner(meta)
	assert.Equal(t, x.Action(1), a)
	meta.Model.Perform(a)

	a = planner(meta)
	assert.Equal(t, x.Action(3), a)
	meta.Model.Perform(a)

	a = planner(meta)
	assert.Equal(t, x.Action(1), a)
	meta.Model.Perform(a)

	a = planner(meta)
	assert.Equal(t, x.Action(3), a)
	meta.Model.Perform(a)

	a = planner(meta)
	assert.Equal(t, x.Action(4), a)
}

func TestSearchDeterministicSerial(t *testing.T) {
	helperSearchDeterministic(t, GetAction)
}

func TestSearchDeterministicParallel(t *testing.T) {
	helperSearchDeterministic(t, GetActionParallel)
}

func BenchmarkSearchHorizon10(b *testing.B) {
	meta.Horizon = 10
	for n := 0; n < b.N; n++ {
		GetAction(meta)
	}
}

func BenchmarkSearchHorizon20(b *testing.B) {
	meta.Horizon = 20
	for n := 0; n < b.N; n++ {
		GetAction(meta)
	}
}

func BenchmarkSearchSamples10k(b *testing.B) {
	meta.Horizon = 10
	meta.Samples = 10000
	for n := 0; n < b.N; n++ {
		GetAction(meta)
	}
}

func BenchmarkSearchSamplesParallel10k(b *testing.B) {
	meta.Horizon = 10
	meta.Samples = 10000
	for n := 0; n < b.N; n++ {
		GetActionParallel(meta)
	}
}
