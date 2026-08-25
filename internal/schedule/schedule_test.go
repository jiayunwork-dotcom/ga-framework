package schedule

import (
	"math"
	"testing"
)

func TestLinear(t *testing.T) {
	l := Linear{Start: 1.0, End: 0.0, Label: "mutation_rate"}
	if l.Value(0, 100) != 1.0 {
		t.Fatalf("expected 1.0 at gen 0")
	}
	if math.Abs(l.Value(50, 100)-0.5) > 1e-10 {
		t.Fatalf("expected 0.5 at gen 50, got %v", l.Value(50, 100))
	}
	if l.Value(100, 100) != 0.0 {
		t.Fatalf("expected 0.0 at gen 100")
	}
}

func TestExponential(t *testing.T) {
	e := Exponential{Start: 1.0, Decay: 0.01, Label: "temp"}
	v0 := e.Value(0, 100)
	v50 := e.Value(50, 100)
	if v0 != 1.0 {
		t.Fatalf("expected 1.0 at gen 0")
	}
	if v50 >= v0 {
		t.Fatalf("should decrease: %v >= %v", v50, v0)
	}
}

func TestStepDecay(t *testing.T) {
	s := StepDecay{Start: 1.0, Factor: 0.5, Interval: 10, Label: "lr"}
	if s.Value(0, 100) != 1.0 {
		t.Fatalf("expected 1.0 at gen 0")
	}
	if s.Value(10, 100) != 0.5 {
		t.Fatalf("expected 0.5 at gen 10, got %v", s.Value(10, 100))
	}
	if s.Value(20, 100) != 0.25 {
		t.Fatalf("expected 0.25 at gen 20, got %v", s.Value(20, 100))
	}
}

func TestCosine(t *testing.T) {
	c := Cosine{Start: 1.0, End: 0.0, Label: "rate"}
	v0 := c.Value(0, 100)
	v100 := c.Value(100, 100)
	if math.Abs(v0-1.0) > 1e-10 {
		t.Fatalf("expected 1.0 at gen 0, got %v", v0)
	}
	if math.Abs(v100) > 1e-10 {
		t.Fatalf("expected 0.0 at gen 100, got %v", v100)
	}
}

func TestCyclicAnnealing(t *testing.T) {
	ca := CyclicAnnealing{Min: 0.01, Max: 0.1, Period: 20, Label: "temp"}
	v0 := ca.Value(0, 100)
	v10 := ca.Value(10, 100)
	if math.Abs(v0-0.1) > 1e-10 {
		t.Fatalf("expected 0.1 at gen 0, got %v", v0)
	}
	if math.Abs(v10-0.01) > 1e-10 {
		t.Fatalf("expected 0.01 at gen 10, got %v", v10)
	}
}

func TestSelfAdaptive(t *testing.T) {
	sa := &SelfAdaptive{Current: 0.1, Min: 0.01, Max: 0.5, Label: "rate"}
	sa.Update(false)
	if sa.Current <= 0.1 {
		t.Fatalf("should increase on no improvement: %v", sa.Current)
	}
	sa.Update(true)
	prev := sa.Current
	sa.Update(true)
	if sa.Current >= prev {
		t.Fatalf("should decrease on improvement: %v >= %v", sa.Current, prev)
	}
}

func TestScheduler_GetAll(t *testing.T) {
	s := NewScheduler()
	s.Register("mutation", Linear{Start: 0.5, End: 0.01})
	s.Register("crossover", Linear{Start: 0.9, End: 0.6})
	vals := s.GetAll(50, 100)
	if len(vals) != 2 {
		t.Fatalf("expected 2 params, got %d", len(vals))
	}
	if vals["mutation"] > 0.5 || vals["mutation"] < 0 {
		t.Fatalf("mutation out of range: %v", vals["mutation"])
	}
}

func TestScheduler_Params(t *testing.T) {
	s := NewScheduler()
	s.Register("a", Linear{Start: 1, End: 0})
	s.Register("b", Linear{Start: 1, End: 0})
	params := s.Params()
	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}
}
