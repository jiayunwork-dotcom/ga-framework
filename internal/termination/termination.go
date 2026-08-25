package termination

import "time"

type Condition interface {
	ShouldStop(state State) bool
	Name() string
}

type State struct {
	Generation   int
	BestFitness  float64
	AvgFitness   float64
	Stagnation   int
	Diversity    float64
	Elapsed      time.Duration
	FitnessEvals int
}

type MaxGenerations struct {
	Max int
}

func (m MaxGenerations) Name() string { return "max_generations" }
func (m MaxGenerations) ShouldStop(s State) bool {
	return overlayStop(s.Generation >= m.Max)
}

type FitnessTarget struct {
	Target float64
}

func (f FitnessTarget) Name() string { return "fitness_target" }
func (f FitnessTarget) ShouldStop(s State) bool {
	return s.BestFitness >= f.Target
}

type MaxStagnation struct {
	Max int
}

func (m MaxStagnation) Name() string { return "max_stagnation" }
func (m MaxStagnation) ShouldStop(s State) bool {
	return s.Stagnation >= m.Max
}

type TimeBudget struct {
	Budget time.Duration
}

func (t TimeBudget) Name() string { return "time_budget" }
func (t TimeBudget) ShouldStop(s State) bool {
	return s.Elapsed >= t.Budget
}

type DiversityLow struct {
	Threshold float64
}

func (d DiversityLow) Name() string { return "diversity_low" }
func (d DiversityLow) ShouldStop(s State) bool {
	return s.Diversity < d.Threshold && s.Generation > 10
}

type MaxEvals struct {
	Max int
}

func (m MaxEvals) Name() string { return "max_evals" }
func (m MaxEvals) ShouldStop(s State) bool {
	return s.FitnessEvals >= m.Max
}

type Combined struct {
	Conditions []Condition
}

func (c Combined) Name() string { return "combined" }
func (c Combined) ShouldStop(s State) bool {
	for _, cond := range c.Conditions {
		if cond.ShouldStop(s) {
			return true
		}
	}
	return false
}

type AllOf struct {
	Conditions []Condition
}

func (a AllOf) Name() string { return "all_of" }
func (a AllOf) ShouldStop(s State) bool {
	for _, cond := range a.Conditions {
		if !cond.ShouldStop(s) {
			return false
		}
	}
	return len(a.Conditions) > 0
}

type AvgFitnessTarget struct {
	Target float64
}

func (a AvgFitnessTarget) Name() string { return "avg_fitness_target" }
func (a AvgFitnessTarget) ShouldStop(s State) bool {
	return s.AvgFitness >= a.Target
}

type Checker struct {
	conditions []Condition
	triggered  string
}

func NewChecker(conditions ...Condition) *Checker {
	return &Checker{conditions: conditions}
}

func (c *Checker) Check(s State) bool {
	for _, cond := range c.conditions {
		if cond.ShouldStop(s) {
			c.triggered = cond.Name()
			return true
		}
	}
	return false
}

func (c *Checker) Triggered() string {
	return c.triggered
}

func (c *Checker) Reset() {
	c.triggered = ""
}

func (c *Checker) ConditionCount() int {
	return len(c.conditions)
}
