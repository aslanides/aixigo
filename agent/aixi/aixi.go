package aixi

import (
	"aixigo/search"
	"aixigo/x"
)

//AImu for now
type AImu struct {
	Meta *search.Meta
}

// Update AImu
func (agent *AImu) Update(a x.Action, e x.Percept) {
	agent.Meta.Model.Perform(a)
}

// GetAction ...
func (agent *AImu) GetAction() x.Action {
	return search.GetAction(agent.Meta)
}