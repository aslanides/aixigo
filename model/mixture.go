package model

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

// Perform (concurrent implementation)
func (m *Mixture) Perform(a x.Action) *x.Percept {
	percept := make(chan *x.Percept, 1)
	n := sample(m)
	var wg sync.WaitGroup
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		go func(idx int, model x.Model) {
			wg.Add(1)
			defer wg.Done()
			e := model.Perform(a)
			if idx == n {
				percept <- e
			}
		}(idx, model)
	}
	wg.Wait()
	return <-percept
}

// Update (concurrent)
func (m *Mixture) Update(a x.Action, e *x.Percept) {
	var wg sync.WaitGroup
	total := make(chan float64, m.n)
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		go func(idx int, model x.Model) {
			wg.Add(1)
			defer wg.Done()
			m.weights[idx] *= model.ConditionalDistribution(e)
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
	e := m.models[n].Perform(a)

	var wg sync.WaitGroup
	total := make(chan float64, m.n)
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		go func(idx int, model x.Model) {
			wg.Add(1)
			defer wg.Done()
			if idx != n {
				model.Perform(a)
			}
			m.weights[idx] *= model.ConditionalDistribution(e)
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
func (m *Mixture) ConditionalDistribution(e *x.Percept) float64 {
	var wg sync.WaitGroup
	total := make(chan float64, m.n)
	for idx, model := range m.models {
		if m.weights[idx] == 0 {
			continue
		}
		go func(idx int, model x.Model) {
			wg.Add(1)
			defer wg.Done()
			total <- m.weights[idx] * model.ConditionalDistribution(e)
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
		go func(model x.Model) {
			wg.Add(1)
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
		go func(model x.Model) {
			wg.Add(1)
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
