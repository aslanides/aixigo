package main

import (
	"aixigo/agent/aixi"
	"aixigo/env/grid"
	"aixigo/search"
	"aixigo/x"
	"fmt"
	"math"
	"math/rand"
	"time"
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

	meta := &search.Meta{
		Meta:    grid.Meta,
		Horizon: 7,
		Samples: 10000,
		UCB:     math.Sqrt2,
		Model:   grid.NewModel(spec),
		Utility: func(e x.Percept, dfr int) float64 { return float64(e.R) },
		PRN:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	agent := &aixi.AImu{Meta: meta}
	cycles := 100
	fmt.Printf("Running for %d cycles with %d samples, using horizon %d\n",
		cycles, meta.Samples, meta.Horizon)
	log := run(agent, env, 100)
	fmt.Printf("Agent's avg reward per cycle: %f\n", averageReward(log))
	fmt.Printf("Optimal avg reward per cycle: %f\n",
		float64(meta.MaxReward)*(float64(cycles)-12.0)/float64(cycles))
}

type trace struct {
	Action  x.Action
	Percept x.Percept
}

func run(agent x.Agent, env x.Environment, cycles int) []trace {
	log := make([]trace, cycles, cycles)
	var a x.Action
	var e x.Percept
	for iter := 0; iter < cycles; iter++ {
		a = agent.GetAction()
		e = env.Perform(a)
		agent.Update(a, e)
		log[iter] = trace{a, e}
	}

	return log
}

func averageReward(log []trace) float64 {
	s := 0.0
	for _, t := range log {
		s += float64(t.Percept.R)
	}

	return s / float64(len(log))
}
