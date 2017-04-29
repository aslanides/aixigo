package search

import (
	"aixigo/env/grid"
	"aixigo/x"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

var meta *Meta
var env x.Environment
var spec [][]int

func init() {
	spec = [][]int{
		{0, 0, 1},
		{1, 0, 0},
		{0, 1, 2},
	}
	meta = NewMeta(grid.Meta, grid.NewModel(spec), 10000)
	env = grid.New(spec)
}

func helperSearchDeterministic(t *testing.T, planner func(*Meta) x.Action) {
	meta.Model = grid.NewModel(spec)
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

// Benchmarks
//
//

func BenchmarkHorizon10Samples1k(b *testing.B) {
	meta.Horizon = 10
	meta.Samples = 1000
	for n := 0; n < b.N; n++ {
		GetAction(meta)
	}
}

func BenchmarkHorizon20Samples1k(b *testing.B) {
	meta.Horizon = 20
	meta.Samples = 1000
	for n := 0; n < b.N; n++ {
		GetAction(meta)
	}
}

func BenchmarkHorizon10Samples10k(b *testing.B) {
	meta.Horizon = 10
	meta.Samples = 10000
	for n := 0; n < b.N; n++ {
		GetAction(meta)
	}
}

func BenchmarkParallelHorizon10Samples10K(b *testing.B) {
	meta.Horizon = 10
	meta.Samples = 10000
	for n := 0; n < b.N; n++ {
		GetActionParallel(meta)
	}
}
