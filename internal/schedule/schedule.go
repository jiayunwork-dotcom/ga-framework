// Package schedule 实现参数自适应调度策略。
package schedule

import "math"

// Strategy 参数调度策略接口。
type Strategy interface {
	Value(generation, maxGen int) float64
	Name() string
}

// Linear 线性衰减策略。
type Linear struct {
	Start float64
	End   float64
	Label string
}

// Name 返回策略名。
func (l Linear) Name() string { return l.Label }

// Value 返回给定代数的参数值。
func (l Linear) Value(generation, maxGen int) float64 {
	if maxGen <= 0 {
		return l.Start
	}
	t := float64(generation) / float64(maxGen)
	if t > 1 {
		t = 1
	}
	return l.Start + (l.End-l.Start)*t
}

// Exponential 指数衰减策略。
type Exponential struct {
	Start float64
	Decay float64
	Label string
}

// Name 返回策略名。
func (e Exponential) Name() string { return e.Label }

// Value 返回值。
func (e Exponential) Value(generation, _ int) float64 {
	return e.Start * math.Exp(-e.Decay*float64(generation))
}

// StepDecay 阶梯衰减。
type StepDecay struct {
	Start    float64
	Factor   float64
	Interval int
	Label    string
}

// Name 返回策略名。
func (s StepDecay) Name() string { return s.Label }

// Value 返回值。
func (s StepDecay) Value(generation, _ int) float64 {
	if s.Interval <= 0 {
		return s.Start
	}
	steps := generation / s.Interval
	return s.Start * math.Pow(s.Factor, float64(steps))
}

// Cosine 余弦退火策略。
type Cosine struct {
	Start float64
	End   float64
	Label string
}

// Name 返回策略名。
func (c Cosine) Name() string { return c.Label }

// Value 使用余弦退火公式。
func (c Cosine) Value(generation, maxGen int) float64 {
	if maxGen <= 0 {
		return c.Start
	}
	t := float64(generation) / float64(maxGen)
	if t > 1 {
		t = 1
	}
	return c.End + 0.5*(c.Start-c.End)*(1+math.Cos(math.Pi*t))
}

// CyclicAnnealing 循环退火策略。
type CyclicAnnealing struct {
	Min    float64
	Max    float64
	Period int
	Label  string
}

// Name 返回策略名。
func (ca CyclicAnnealing) Name() string { return ca.Label }

// Value 循环退火值。
func (ca CyclicAnnealing) Value(generation, _ int) float64 {
	if ca.Period <= 0 {
		return ca.Min
	}
	phase := float64(generation%ca.Period) / float64(ca.Period)
	return ca.Min + 0.5*(ca.Max-ca.Min)*(1+math.Cos(2*math.Pi*phase))
}

// SelfAdaptive 自适应策略：基于种群改进率动态调整。
type SelfAdaptive struct {
	Current float64
	Min     float64
	Max     float64
	Label   string
}

// Name 返回策略名。
func (sa *SelfAdaptive) Name() string { return sa.Label }

// Value 当前值。
func (sa *SelfAdaptive) Value(_, _ int) float64 {
	return sa.Current
}

// Update 根据改进情况更新参数。
func (sa *SelfAdaptive) Update(improved bool) {
	if improved {
		applyImprovement(sa)
		return
	}
	applyStagnation(sa)
}

// Scheduler 管理多个参数的调度。
type Scheduler struct {
	strategies map[string]Strategy
}

// NewScheduler 创建调度器。
func NewScheduler() *Scheduler {
	return &Scheduler{
		strategies: make(map[string]Strategy),
	}
}

// Register 注册参数调度策略。
func (s *Scheduler) Register(param string, strategy Strategy) {
	s.strategies[param] = strategy
}

// Get 获取参数当前值。
func (s *Scheduler) Get(param string, generation, maxGen int) float64 {
	st, ok := s.strategies[param]
	if !ok {
		return 0
	}
	return st.Value(generation, maxGen)
}

// GetAll 获取所有参数的当前值。
func (s *Scheduler) GetAll(generation, maxGen int) map[string]float64 {
	vals := make(map[string]float64, len(s.strategies))
	for name, st := range s.strategies {
		vals[name] = st.Value(generation, maxGen)
	}
	return vals
}

// Params 返回所有已注册参数名。
func (s *Scheduler) Params() []string {
	names := make([]string, 0, len(s.strategies))
	for name := range s.strategies {
		names = append(names, name)
	}
	return names
}
