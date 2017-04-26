package agent

import (
	"aixigo/env"
	"aixigo/x"
)

//AImu for now
type AImu struct {
	samples uint
	model   env.Environment
}

// Update AImu
func (agent *AImu) Update(a x.Action, e x.Percept) {
	agent.model.Perform(a)
}

// GetAction is the bomb!
func (agent *AImu) GetAction() {

}
