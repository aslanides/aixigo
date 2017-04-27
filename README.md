# AIXIgo
Fast, scalable parallel MC-AIXI implementation in Golang.

Here, I'm implementing algorithmic and implementation-specific optimizations to see how fast & scalable we can get AIXI with Monte Carlo Tree Search :)

## API
We opt for a simpler API than the one in [AIXIjs](https://github.com/aslanides/aixijs). The  agent-environment loop looks like:

```go
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
```

## Performance
We implement MCTS with root parallelism (Chaslot, Winands, & Herik, 2008b), and we get close to linear speedup over the serial implementation. Here's a benchmark on a small deterministic Gridworld, running AImu on an i7-3770 (8 virtual cores):

```
BenchmarkHorizon10Samples1k-8            	     300	   5184511 ns/op
BenchmarkHorizon20Samples1k-8            	     200	   7593272 ns/op
BenchmarkHorizon10Samples10k-8           	      20	  61568848 ns/op
BenchmarkParallelHorizon10Samples10K-8   	     200	   7805558 ns/op
PASS
ok  	aixigo/search	8.065s
```

From `BenchmarkParallelHorizon10Samples10K-8`, we see that we can run 10,000 samples with a horizon of 10 in _7 milliseconds_. This is definitely a lower bound for the performance we can realistically expect for AIXI, given that AImu does no model updates. Introducing a Bayes mixture will slow things down by a factor of M, where M is the cardinality of the model class. 
