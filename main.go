package main

import (
	"aixigo/agent/aixi"
	"aixigo/env/grid"
	"aixigo/search"
	"aixigo/x"
	"fmt"
)

func main() {
	spec := [][]int{
		{0, 0, 1, 1, 1},
		{1, 0, 0, 1, 2},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0},
	}
	env := grid.New(spec)
	meta := search.NewMeta(grid.Meta, grid.NewModel(spec), 10000)

	agent := &aixi.AImu{Meta: meta}
	cycles := 100
	fmt.Printf("Running for %d cycles with %d samples, using horizon %d\n",
		cycles, meta.Samples, meta.Horizon)
	log := run(agent, env, 100)
	fmt.Printf("Agent's avg reward per cycle: %f\n", averageReward(log))
	fmt.Printf("Optimal avg reward per cycle: %f\n",
		float64(meta.MaxReward)*(float64(cycles)-10.0)/float64(cycles))
}

type trace struct {
	Action      x.Action
	Observation x.Observation
	Reward      x.Reward
}

func run(agent x.Agent, env x.Environment, cycles int) []trace {
	log := make([]trace, cycles, cycles)
	var a x.Action
	var o x.Observation
	var r x.Reward
	for iter := 0; iter < cycles; iter++ {
		a = agent.GetAction()
		o, r = env.Perform(a)
		agent.Update(a, o, r)
		log[iter] = trace{a, o, r}
	}

	return log
}

func averageReward(log []trace) float64 {
	s := 0.0
	for _, t := range log {
		s += float64(t.Reward)
	}

	return s / float64(len(log))
}
