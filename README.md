# AIXIgo
Fast, scalable parallel MC-AIXI implementation in Golang.

Here, I'm implementing algorithmic and implementation-specific optimizations to see how fast & scalable we can get AIXI with Monte Carlo Tree Search (ρUCT, Veness et al. 2010) :)

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
We implement MCTS with root parallelism (Chaslot, Winands, & Herik, 2008b), and we get close to linear speedup over the serial implementation. Here's an initial benchmark on a small deterministic Gridworld, running AImu on an i7-3770 (8 virtual cores):

```
BenchmarkHorizon10Samples1k-8            	    1000	   2519791 ns/op
BenchmarkHorizon20Samples1k-8            	     500	   3432620 ns/op
BenchmarkHorizon10Samples10k-8           	     100	  23130465 ns/op
BenchmarkParallelHorizon10Samples10K-8   	     500	   3932721 ns/op
PASS
ok  	aixigo/search	9.710s
```

From `BenchmarkParallelHorizon10Samples10K-8`, we see that we can run 10,000 samples with a horizon of 10 around 4 milliseconds. This is obviously a very loose lower bound for the performance we can realistically expect from AIXI, given that AImu does no model updates. Introducing a Bayes mixture will slow things down by a factor proportional to the cardinality of the model class. Using a sequential density estimator like CTW will introduce a factor proportional to the depth of the context tree.

Note that using root parallelism will generally result in a less accurate value function estimate (and thus a weaker policy), since more time is spent doing rollouts as we build each of the independent trees from scratch. 

### Profiling

We profile the main `run` method in `main_test.go` using:

`go test -run=^$ -bench=. -cpuprofile=cpu.out`

followed by running the interactive profiler `pprof`

`go tool pprof aixigo.test cpu.out`

We can then generate the execution graph using `web`.

## TODO

* Port over the mixture and Dirichlet models from AIXIjs
* Generalize the agent implementation and port KSA/Thompson sampling
* Try out tree parallelism and see how this affects performance (clearly we'll need to use mutexes within the tree, and so we will lose some performance by having blocked goroutines, but we gain in having 8x less memory for the GC to keep track of. We can also maybe run > nCPU goroutines to make up for the blocking, so that the scheduler can always find a non-blocked thread for each CPU. Maybe.)
* Add more instrumentation, and implement more interesting environments (POCman, anyone?)
