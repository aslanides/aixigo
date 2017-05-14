package ctw

import "aixigo/x"

// CTW ...
type CTW struct {
	ObsBits    []symbol
	RewBits    []symbol
	ActionBits []symbol
	// TODO: Wrap the context tree and implement the x.Model interface
	// interface methods are stubbed out right now
}

// NewCTW ...
func NewCTW(meta *x.Meta) x.Model {
	return &CTW{} // TODO
}

// Perform ...
func (m *CTW) Perform(a x.Action) (x.Observation, x.Reward) {
	return x.Observation(0), x.Reward(0)
}

// Update ...
func (m *CTW) Update(a x.Action, o x.Observation, r x.Reward) {}

// ConditionalDistribution ...
func (m *CTW) ConditionalDistribution(o x.Observation, r x.Reward) float64 {
	return 0. // TODO
}

// Copy ...
func (m *CTW) Copy() x.Model {
	return m // TODO
}

// Save ...
func (m *CTW) Save() {}

// Load ...
func (m *CTW) Load() {}
