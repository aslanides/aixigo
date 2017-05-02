package ctw

import "aixigo/x"

// CTW ...
type CTW struct {
	// TODO: Wrap the context tree and implement the x.Model interface
	// interface methods are stubbed out right now
}

// NewCTW ...
func NewCTW() x.Model {
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

// SaveCheckpoint ...
func (m *CTW) SaveCheckpoint() {}

// LoadCheckpoint ...
func (m *CTW) LoadCheckpoint() {}
