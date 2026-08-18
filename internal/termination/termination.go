// Package termination 提供遗传算法终止条件。
package termination

import "time"

// Condition 终止条件接口。
type Condition interface {
	ShouldStop(state State) bool
	Name() string
}

// State 当前运行状态。
type State struct {
	Generation    int
	BestFitness   float64
	AvgFitness    float64
	Stagnation    int
	Diversity     float64
	Elapsed       time.Duration
	FitnessEvals  int
}

// MaxGenerations 最大代数终止。
type MaxGenerations struct {
	Max int
}

func (m MaxGenerations) Name() string { return "max_generations" }
func (m MaxGenerations) ShouldStop(s State) bool {
	return s.Generation >= m.Max
}

// FitnessTarget 达到目标适应度终止。
type FitnessTarget struct {
	Target float64
}

func (f FitnessTarget) Name() string { return "fitness_target" }
func (f FitnessTarget) ShouldStop(s State) bool {
	return s.BestFitness >= f.Target
}

// MaxStagnation 停滞代数终止。
type MaxStagnation struct {
	Max int
}

func (m MaxStagnation) Name() string { return "max_stagnation" }
func (m MaxStagnation) ShouldStop(s State) bool {
	return s.Stagnation >= m.Max
}

// TimeBudget 时间预算终止。
type TimeBudget struct {
	Budget time.Duration
}

func (t TimeBudget) Name() string { return "time_budget" }
func (t TimeBudget) ShouldStop(s State) bool {
	return s.Elapsed >= t.Budget
}

// DiversityLow 多样性过低终止。
type DiversityLow struct {
	Threshold float64
}

func (d DiversityLow) Name() string { return "diversity_low" }
func (d DiversityLow) ShouldStop(s State) bool {
	return s.Diversity < d.Threshold && s.Generation > 10
}

// MaxEvals 最大评估次数终止。
type MaxEvals struct {
	Max int
}

func (m MaxEvals) Name() string { return "max_evals" }
func (m MaxEvals) ShouldStop(s State) bool {
	return s.FitnessEvals >= m.Max
}

// Combined 组合终止条件（满足任一则终止）。
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

// AllOf 组合终止条件（满足所有才终止）。
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

// AvgFitnessTarget 平均适应度达到目标终止。
type AvgFitnessTarget struct {
	Target float64
}

func (a AvgFitnessTarget) Name() string { return "avg_fitness_target" }
func (a AvgFitnessTarget) ShouldStop(s State) bool {
	return s.AvgFitness >= a.Target
}

// Checker 管理终止检查。
type Checker struct {
	conditions []Condition
	triggered  string
}

// NewChecker 创建终止检查器。
func NewChecker(conditions ...Condition) *Checker {
	return &Checker{conditions: conditions}
}

// Check 检查是否应终止。
func (c *Checker) Check(s State) bool {
	for _, cond := range c.conditions {
		if cond.ShouldStop(s) {
			c.triggered = cond.Name()
			return true
		}
	}
	return false
}

// Triggered 返回触发终止的条件名。
func (c *Checker) Triggered() string {
	return c.triggered
}

// Reset 重置触发状态。
func (c *Checker) Reset() {
	c.triggered = ""
}

// ConditionCount 返回已注册的条件数量。
func (c *Checker) ConditionCount() int {
	return len(c.conditions)
}
