package termination

import (
	"testing"
	"time"
)

func TestMaxGenerations(t *testing.T) {
	c := MaxGenerations{Max: 100}
	if c.ShouldStop(State{Generation: 50}) {
		t.Fatalf("should not stop at gen 50")
	}
	if !c.ShouldStop(State{Generation: 100}) {
		t.Fatalf("should stop at gen 100")
	}
}

func TestFitnessTarget(t *testing.T) {
	c := FitnessTarget{Target: 95.0}
	if c.ShouldStop(State{BestFitness: 80}) {
		t.Fatalf("should not stop below target")
	}
	if !c.ShouldStop(State{BestFitness: 95}) {
		t.Fatalf("should stop at target")
	}
}

func TestMaxStagnation(t *testing.T) {
	c := MaxStagnation{Max: 20}
	if c.ShouldStop(State{Stagnation: 10}) {
		t.Fatalf("should not stop with stagnation=10")
	}
	if !c.ShouldStop(State{Stagnation: 20}) {
		t.Fatalf("should stop with stagnation=20")
	}
}

func TestTimeBudget(t *testing.T) {
	c := TimeBudget{Budget: 5 * time.Minute}
	if c.ShouldStop(State{Elapsed: 2 * time.Minute}) {
		t.Fatalf("should not stop within budget")
	}
	if !c.ShouldStop(State{Elapsed: 6 * time.Minute}) {
		t.Fatalf("should stop when budget exceeded")
	}
}

func TestDiversityLow(t *testing.T) {
	c := DiversityLow{Threshold: 0.01}
	if c.ShouldStop(State{Diversity: 0.005, Generation: 5}) {
		t.Fatalf("should not stop early even if diversity low")
	}
	if !c.ShouldStop(State{Diversity: 0.005, Generation: 20}) {
		t.Fatalf("should stop when diversity drops after initial generations")
	}
}

func TestCombined(t *testing.T) {
	c := Combined{Conditions: []Condition{
		MaxGenerations{Max: 100},
		FitnessTarget{Target: 50},
	}}
	if c.ShouldStop(State{Generation: 50, BestFitness: 30}) {
		t.Fatalf("should not stop when neither condition met")
	}
	if !c.ShouldStop(State{Generation: 50, BestFitness: 50}) {
		t.Fatalf("should stop when fitness target met")
	}
}

func TestAllOf(t *testing.T) {
	c := AllOf{Conditions: []Condition{
		MaxGenerations{Max: 100},
		FitnessTarget{Target: 50},
	}}
	// Only one met -> don't stop
	if c.ShouldStop(State{Generation: 100, BestFitness: 30}) {
		t.Fatalf("should not stop when only one condition met")
	}
	if !c.ShouldStop(State{Generation: 100, BestFitness: 50}) {
		t.Fatalf("should stop when all conditions met")
	}
}

func TestChecker(t *testing.T) {
	ch := NewChecker(
		MaxGenerations{Max: 100},
		FitnessTarget{Target: 95},
	)
	if ch.Check(State{Generation: 50, BestFitness: 80}) {
		t.Fatalf("should not trigger")
	}
	if ch.Triggered() != "" {
		t.Fatalf("no trigger expected")
	}
	if !ch.Check(State{Generation: 50, BestFitness: 95}) {
		t.Fatalf("should trigger on fitness target")
	}
	if ch.Triggered() != "fitness_target" {
		t.Fatalf("expected fitness_target, got %q", ch.Triggered())
	}
	ch.Reset()
	if ch.Triggered() != "" {
		t.Fatalf("should be empty after reset")
	}
}

func TestMaxEvals(t *testing.T) {
	c := MaxEvals{Max: 10000}
	if c.ShouldStop(State{FitnessEvals: 5000}) {
		t.Fatalf("should not stop below max evals")
	}
	if !c.ShouldStop(State{FitnessEvals: 10000}) {
		t.Fatalf("should stop at max evals")
	}
}

func TestConditionCount(t *testing.T) {
	ch := NewChecker(MaxGenerations{Max: 10}, FitnessTarget{Target: 5})
	if ch.ConditionCount() != 2 {
		t.Fatalf("expected 2, got %d", ch.ConditionCount())
	}
}
