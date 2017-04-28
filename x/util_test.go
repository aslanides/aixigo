package x

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	var o Observation
	o = Observation{false, true, true, false, true}
	assert.Equal(t, 13, ToInt(o))
	o = Observation{false, false, false}
	assert.Equal(t, 0, ToInt(o))
	o = Observation{true, false, false, false, false}
	assert.Equal(t, 16, ToInt(o))
	o = Observation{false, false, false, false, true}
	assert.Equal(t, 1, ToInt(o))
}

func TestArgMax(t *testing.T) {
	var A []float64
	A = []float64{1.1, 3.3, -123.2, 40000.1, 1e6, -1e7, 22.3, -3.14}
	assert.Equal(t, 4, ArgMax(A))
}
