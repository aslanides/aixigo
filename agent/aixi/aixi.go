package aixi

import (
	"aixigo/search"
	"aixigo/x"
)

//AImu only for now
type AImu struct {
	Meta *search.Meta
}

// Update ...
func (agent *AImu) Update(a x.Action, e x.Percept) {
	agent.Meta.Model.Perform(a)
}

// GetAction uses the Parallel implementation
func (agent *AImu) GetAction() x.Action {
	return search.GetActionParallel(agent.Meta)
}
