package main

import (
	"aixigo/agent/aixi"
	"aixigo/env/grid"
	"aixigo/mcts"
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

	meta := mcts.NewMeta(grid.Meta, grid.NewModel(spec), 10000)
	agent := &aixi.AImu{Meta: meta}
	cycles := 100
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(agent, env, cycles)
	}
}
