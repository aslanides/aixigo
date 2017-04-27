package x

import (
	"math"
	"math/rand"
	"time"
)

//NewPRN ...
func NewPRN() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

//ArgMax ...
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
