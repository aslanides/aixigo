package x

import (
	"math"
	"math/rand"
	"time"
)

// NewPRN creates a new rand.Rand object with its own seed
func NewPRN() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

//ArgMax returns the argmax (index) for an array of floats
func ArgMax(A []float64) int {
	max := math.Inf(-1)
	idx := -1
	for i, v := range A {
		if v > max {
			max = v
			idx = i
		}
	}
	return idx
}

//ToInt helper for Observation objects
func ToInt(o ObservationBits) Observation {
	n := 0
	for _, b := range o {
		n <<= 1
		if b {
			n++
		}
	}
	return Observation(n)
}

// Equals checks equality of percepts
func Equals(e, p *Percept) bool {
	return p.R == e.R && p.O == e.O
}

// RLUtility is the utility function for normal reward-based reinforcement learners
func RLUtility(o Observation, r Reward, dfr int) float64 {
	return float64(r)
}

// Log2 for integers to hackishly improve performance in MCTS
func Log2(v uint) int {
	n := -1
	for ; v > 0; n++ {
		v >>= 1
	}
	return n
}
