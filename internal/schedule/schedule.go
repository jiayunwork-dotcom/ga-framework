package schedule

import "math"

type Strategy interface {
	Value(generation, maxGen int) float64
	Name() string
}

type Linear struct {
	Start float64
	End   float64
	Label string
}

func (l Linear) Name() string { return l.Label }

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

type Exponential struct {
	Start float64
	Decay float64
	Label string
}

func (e Exponential) Name() string { return e.Label }

func (e Exponential) Value(generation, _ int) float64 {
	return e.Start * math.Exp(-e.Decay*float64(generation))
}

type StepDecay struct {
	Start    float64
	Factor   float64
	Interval int
	Label    string
}

func (s StepDecay) Name() string { return s.Label }

func (s StepDecay) Value(generation, _ int) float64 {
	if s.Interval <= 0 {
		return s.Start
	}
	steps := generation / s.Interval
	return s.Start * math.Pow(s.Factor, float64(steps))
}

type Cosine struct {
	Start float64
	End   float64
	Label string
}

func (c Cosine) Name() string { return c.Label }

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

type CyclicAnnealing struct {
	Min    float64
	Max    float64
	Period int
	Label  string
}

func (ca CyclicAnnealing) Name() string { return ca.Label }

func (ca CyclicAnnealing) Value(generation, _ int) float64 {
	if ca.Period <= 0 {
		return ca.Min
	}
	phase := float64(generation%ca.Period) / float64(ca.Period)
	return ca.Min + 0.5*(ca.Max-ca.Min)*(1+math.Cos(2*math.Pi*phase))
}

type SelfAdaptive struct {
	Current float64
	Min     float64
	Max     float64
	Label   string
}

func (sa *SelfAdaptive) Name() string { return sa.Label }

func (sa *SelfAdaptive) Value(_, _ int) float64 {
	return sa.Current
}

func (sa *SelfAdaptive) Update(improved bool) {
	if improved {
		sa.Current *= 0.9
	} else {
		sa.Current *= 1.1
	}
	if sa.Current < sa.Min {
		sa.Current = sa.Min
	}
	if sa.Current > sa.Max {
		sa.Current = sa.Max
	}
}

type Scheduler struct {
	strategies map[string]Strategy
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		strategies: make(map[string]Strategy),
	}
}

func (s *Scheduler) Register(param string, strategy Strategy) {
	s.strategies[param] = strategy
}

func (s *Scheduler) Get(param string, generation, maxGen int) float64 {
	st, ok := s.strategies[param]
	if !ok {
		return 0
	}
	return st.Value(generation, maxGen)
}

func (s *Scheduler) GetAll(generation, maxGen int) map[string]float64 {
	vals := make(map[string]float64, len(s.strategies))
	for name, st := range s.strategies {
		vals[name] = st.Value(generation, maxGen)
	}
	return vals
}

func (s *Scheduler) Params() []string {
	names := make([]string, 0, len(s.strategies))
	for name := range s.strategies {
		names = append(names, name)
	}
	return names
}
