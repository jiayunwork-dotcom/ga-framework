package constraint

import (
	"math"
	"testing"
)

func TestIsFeasible(t *testing.T) {
	h := NewHandler(StaticPenalty, 1.0)
	h.AddConstraint(BoundsConstraint(-5, 5))
	if !h.IsFeasible([]float64{0, 1, -3}) {
		t.Fatalf("should be feasible")
	}
	if h.IsFeasible([]float64{0, 1, 6}) {
		t.Fatalf("should be infeasible")
	}
}

func TestTotalViolation(t *testing.T) {
	h := NewHandler(StaticPenalty, 1.0)
	h.AddConstraint(BoundsConstraint(0, 10))
	v := h.TotalViolation([]float64{-2, 12, 5})
	// violations: 2 (from -2) + 2 (from 12)
	if math.Abs(v-4.0) > 1e-10 {
		t.Fatalf("expected violation 4, got %v", v)
	}
}

func TestStaticPenalty(t *testing.T) {
	h := NewHandler(StaticPenalty, 10.0)
	h.AddConstraint(SumConstraint(5.0))
	fit := h.PenalizedFitness(100, []float64{3, 3}, 0)
	violation := 1.0 // sum=6, excess=1
	expected := 100 - 10.0*violation*violation
	if math.Abs(fit-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, fit)
	}
}

func TestDynamicPenalty(t *testing.T) {
	h := NewHandler(DynamicPenalty, 1.0)
	h.AddConstraint(SumConstraint(5.0))
	fit0 := h.PenalizedFitness(100, []float64{3, 3}, 0)
	fit10 := h.PenalizedFitness(100, []float64{3, 3}, 10)
	if fit10 >= fit0 {
		t.Fatalf("dynamic penalty should increase with generation: gen0=%v, gen10=%v", fit0, fit10)
	}
}

func TestDeathPenalty(t *testing.T) {
	h := NewHandler(DeathPenalty, 1.0)
	h.AddConstraint(BoundsConstraint(0, 1))
	fit := h.PenalizedFitness(100, []float64{2}, 0)
	if !math.IsInf(fit, -1) {
		t.Fatalf("death penalty should give -Inf, got %v", fit)
	}
}

func TestEqualityConstraint(t *testing.T) {
	eq := EqualityConstraint("sum_eq_10", func(g []float64) float64 {
		s := 0.0
		for _, x := range g {
			s += x
		}
		return s - 10
	}, 0.01)
	if eq.Evaluate([]float64{5, 5}) > 0 {
		t.Fatalf("sum=10 should satisfy equality within epsilon")
	}
	if eq.Evaluate([]float64{5, 6}) <= 0 {
		t.Fatalf("sum=11 should violate equality")
	}
}

func TestRepair(t *testing.T) {
	genes := []float64{-2, 3, 12, 5}
	repaired := Repair(genes, 0, 10)
	for i, g := range repaired {
		if g < 0 || g > 10 {
			t.Fatalf("gene %d out of bounds after repair: %v", i, g)
		}
	}
	if repaired[0] != 0 || repaired[2] != 10 {
		t.Fatalf("repair incorrect: %v", repaired)
	}
}

func TestBounceBack(t *testing.T) {
	genes := []float64{-1, 11, 5}
	bounced := BounceBack(genes, 0, 10)
	for i, g := range bounced {
		if g < 0 || g > 10 {
			t.Fatalf("gene %d out of bounds after bounce: %v", i, g)
		}
	}
	if bounced[0] != 1 { // -1 bounces to 1
		t.Fatalf("expected bounce to 1, got %v", bounced[0])
	}
}

func TestFeasibleRatio(t *testing.T) {
	h := NewHandler(StaticPenalty, 1.0)
	h.AddConstraint(BoundsConstraint(0, 10))
	pop := [][]float64{
		{1, 2, 3},
		{-1, 2, 3},
		{5, 5, 5},
		{0, 11, 5},
	}
	ratio := h.FeasibleRatio(pop)
	if math.Abs(ratio-0.5) > 1e-10 {
		t.Fatalf("expected ratio 0.5, got %v", ratio)
	}
}

func TestCountViolations(t *testing.T) {
	h := NewHandler(StaticPenalty, 1.0)
	h.AddConstraint(BoundsConstraint(0, 10))
	h.AddConstraint(SumConstraint(5))
	count := h.CountViolations([]float64{3, 4})
	if count != 1 {
		t.Fatalf("expected 1 violation (sum=7>5), got %d", count)
	}
}
