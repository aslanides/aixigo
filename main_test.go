package main

import (
	"aixigo/agent/aixi"
	"aixigo/env/grid"
	"aixigo/search"
	"aixigo/x"
	"math"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	spec := [][]int{
		{0, 0, 1, 1, 1},
		{1, 0, 0, 1, 2},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0},
	}
	env := grid.New(spec)

	meta := &search.Meta{
		Meta:    grid.Meta,
		Horizon: 10,
		Samples: 10000,
		UCB:     math.Sqrt2,
		Model:   grid.NewModel(spec),
		Utility: func(e x.Percept, dfr int) float64 { return float64(e.R) },
		PRN:     x.NewPRN(),
	}
	agent := &aixi.AImu{Meta: meta}
	cycles := 100
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(agent, env, cycles)
	}
}
