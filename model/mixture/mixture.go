package mixture

import (
	"aixigo/x"
	"math/rand"
	"sync"
)

// Mixture is Bayes mixture
type Mixture struct {
	models  []x.Model
	weights []float64
	n       int
	prn     *rand.Rand
}

// NewMixture ...
func NewMixture(models []x.Model) x.Model {
	// only use a uniform distribution for now
	n := len(models)
	weights := make([]float64, n, n)
	for i := 0; i < n; i++ {
		weights[i] = 1.0 / float64(n)
	}

	return &Mixture{
		models:  models,
		weights: weights,
		n:       n,
		prn:     x.NewPRN(),
	}
}

// Not sure if all this concurrency will really help, since all cores will be
// maxed most of the time in search anyway...
// TODO: profile serial vs concurrent

// Occasionally we will need to perform actions & update separately

// Perform (serial implementation)
func (m *Mixture) Perform(a x.Action) (x.Observation, x.Reward) {
	n := sample(m)
	o, r := m.models[n].Perform(a)
	for idx, model := range m.models {
		if m.weights[idx] == 0 || idx == n {
			continue
		}
		model.Perform(a)
	}
	return o, r
}

// ParPerform (concurrent implementation)
func (m *Mixture) ParPerform(a x.Action) (x.Observation, x.Reward) {
	n := sample(m)
	var wg sync.WaitGroup
	o, r := m.models[n].Perform(a)
	for idx, model := range m.models {
		if m.weights[idx] == 0 || idx == n {
			continue
		}
		wg.Add(1)
		go func(idx int, model x.Model) {
			defer wg.Done()
			model.Perform(a)
		}(idx, model)
	}
	wg.Wait()
	return o, r
}

// Update (serial)
func (m *Mixture) Update(a x.Action, o x.Observation, r x.Reward) {
	xi := 0.0
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		m.weights[idx] *= model.ConditionalDistribution(o, r)
		xi += m.weights[idx]
	}

	for i := 0; i < m.n; i++ {
		if m.weights[i] == 0 {
			continue
		}
		m.weights[i] /= xi
	}
}

// ParUpdate (concurrent)
func (m *Mixture) ParUpdate(a x.Action, o x.Observation, r x.Reward) {
	var wg sync.WaitGroup
	total := make(chan float64, m.n)
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		wg.Add(1)
		go func(idx int, model x.Model) {
			defer wg.Done()
			m.weights[idx] *= model.ConditionalDistribution(o, r)
			total <- m.weights[idx]
		}(idx, model)
	}
	wg.Wait()
	close(total)
	xi := 0.0
	for x := range total {
		xi += x
	}

	// unfortunately i think this has to be serial code ...
	for i := 0; i < m.n; i++ {
		m.weights[i] /= xi // some wasted iterations but c'est la vie (i think)
	}
}

// GeneratePerceptAndUpdate is a streamlined version of:
//
// e := m.Perform(a)
// m.Update(a, e)
//
// for use in MCTS (where performance is critical)
//
// Yes, this is a *lot* of code duplication for a possibly minor performance gain
func (m *Mixture) GeneratePerceptAndUpdate(a x.Action) {
	n := sample(m)
	o, r := m.models[n].Perform(a)

	var wg sync.WaitGroup
	total := make(chan float64, m.n)
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		wg.Add(1)
		go func(idx int, model x.Model) {
			defer wg.Done()
			if idx != n {
				model.Perform(a)
			}
			m.weights[idx] *= model.ConditionalDistribution(o, r)
			total <- m.weights[idx]
		}(idx, model)
	}
	wg.Wait()
	close(total)
	xi := 0.0
	for x := range total {
		xi += x
	}

	// unfortunately i think this has to be serial code ...
	for i := 0; i < m.n; i++ {
		m.weights[i] /= xi // some wasted iterations but c'est la vie (i think)
	}
}

// ConditionalDistribution (won't get used much)
func (m *Mixture) ConditionalDistribution(o x.Observation, r x.Reward) float64 {
	var wg sync.WaitGroup
	total := make(chan float64, m.n)
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		wg.Add(1)
		go func(idx int, model x.Model) {
			defer wg.Done()
			total <- m.weights[idx] * model.ConditionalDistribution(o, r)
		}(idx, model)
	}
	wg.Wait()
	close(total)
	xi := 0.0
	for x := range total {
		xi += x
	}
	return xi
}

// SaveCheckpoint ...
func (m *Mixture) SaveCheckpoint() {
	var wg sync.WaitGroup
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		wg.Add(1)
		go func(model x.Model) {
			defer wg.Done()
			model.SaveCheckpoint()
		}(model)
	}
	wg.Wait()
}

// LoadCheckpoint ...
func (m *Mixture) LoadCheckpoint() {
	var wg sync.WaitGroup
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		wg.Add(1)
		go func(model x.Model) {
			defer wg.Done()
			model.LoadCheckpoint()
		}(model)
	}
	wg.Wait()
}

// Copy (shouldn't be called often, so is allowed to be slow)
func (m *Mixture) Copy() x.Model {
	models := make([]x.Model, m.n, m.n)
	weights := make([]float64, m.n, m.n)
	copy(weights, m.weights)
	for idx, model := range m.models {
		models[idx] = model.Copy()
	}

	return &Mixture{
		models:  models,
		weights: weights,
		n:       m.n,
		prn:     x.NewPRN(),
	}
}

// this function has to go here because it needs access to its own m.prn
// since the global rand.Rand object locks
func sample(m *Mixture) int {
	s := m.prn.Float64()
	p := 0.0
	for i, w := range m.weights {
		if s <= p {
			return i - 1
		}
		p += w
	}
	return len(m.weights) - 1
}
