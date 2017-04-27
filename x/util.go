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
func ToInt(o Observation) int {
	s := 0
	for _, b := range o {
		if b {
			s++
		}
	}
	return s
}
