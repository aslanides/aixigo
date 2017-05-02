package aixi

import (
	"aixigo/mcts/rho"
	"aixigo/x"
)

//AImu only for now
type AImu struct {
	Meta *rho.Meta
}

// Update ...
func (agent *AImu) Update(a x.Action, o x.Observation, r x.Reward) {
	agent.Meta.Model.Perform(a)
}

// GetAction uses the Parallel implementation
func (agent *AImu) GetAction() x.Action {
	return rho.GetActionParallel(agent.Meta)
}
