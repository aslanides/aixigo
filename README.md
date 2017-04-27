# AIXIgo
Fast, scalable parallel MC-AIXI implementation in Golang.

Here, I'm implementing algorithmic and implementation-specific optimizations to see how fast & scalable we can get AIXI with Monte Carlo Tree Search :)

## Performance
We implement MCTS with root parallelism (Chaslot, Winands, & Herik, 2008b), and we get close to linear speedup over the serial implementation. Here's a benchmark on a small deterministic Gridworld, running AI$\mu$ on an i7-3770 (8 virtual cores):

```
BenchmarkHorizon10Samples1k-8            	     300	   5184511 ns/op
BenchmarkHorizon20Samples1k-8            	     200	   7593272 ns/op
BenchmarkHorizon10Samples10k-8           	      20	  61568848 ns/op
BenchmarkParallelHorizon10Samples10K-8   	     200	   7805558 ns/op
PASS
ok  	aixigo/search	8.065s
```
