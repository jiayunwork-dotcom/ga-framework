package niching

import (
	"math"
	"math/rand"
	"testing"
)

func TestSharingFunction_InRange(t *testing.T) {
	s := SharingFunction(1.0, 5.0, 1.0)
	expected := 1 - 1.0/5.0
	if math.Abs(s-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, s)
	}
}

func TestSharingFunction_OutOfRange(t *testing.T) {
	s := SharingFunction(10.0, 5.0, 1.0)
	if s != 0 {
		t.Fatalf("expected 0 for dist >= sigma, got %v", s)
	}
}

func TestSharedFitness_Reduces(t *testing.T) {
	pop := []Individual{
		{Genes: []float64{0, 0}, Fitness: 10},
		{Genes: []float64{0.1, 0.1}, Fitness: 10},
		{Genes: []float64{100, 100}, Fitness: 10},
	}
	shared := SharedFitness(pop, 5.0, 1.0)
	if shared[0] >= shared[2] {
		t.Fatalf("close individuals should have lower shared fitness: %v vs %v", shared[0], shared[2])
	}
}

func TestSpeciation(t *testing.T) {
	pop := []Individual{
		{Genes: []float64{0, 0}},
		{Genes: []float64{0.5, 0}},
		{Genes: []float64{10, 10}},
		{Genes: []float64{10.5, 10}},
	}
	species := Speciation(pop, 2.0)
	if len(species) != 2 {
		t.Fatalf("expected 2 species, got %d", len(species))
	}
}

func TestClearingSelection(t *testing.T) {
	pop := []Individual{
		{Genes: []float64{0, 0}, Fitness: 10},
		{Genes: []float64{0.1, 0}, Fitness: 8},
		{Genes: []float64{0.2, 0}, Fitness: 6},
		{Genes: []float64{10, 10}, Fitness: 9},
	}
	result := ClearingSelection(pop, 2.0, 1)
	clearedCount := 0
	for _, ind := range result {
		if ind.Fitness == 0 {
			clearedCount++
		}
	}
	if clearedCount == 0 {
		t.Fatalf("clearing should zero some fitnesses")
	}
}

func TestDeterministicCrowding(t *testing.T) {
	parent := Individual{Fitness: 10}
	better := Individual{Fitness: 15}
	worse := Individual{Fitness: 5}
	if DeterministicCrowding(parent, better).Fitness != 15 {
		t.Fatalf("better child should replace parent")
	}
	if DeterministicCrowding(parent, worse).Fitness != 10 {
		t.Fatalf("worse child should not replace parent")
	}
}

func TestRestrictedTournament(t *testing.T) {
	pop := []Individual{
		{Genes: []float64{0, 0}},
		{Genes: []float64{1, 1}},
		{Genes: []float64{1000, 1000}},
	}
	candidate := Individual{Genes: []float64{0.5, 0.5}}
	rnd := rand.New(rand.NewSource(99))
	found2 := 0
	trials := 100
	for i := 0; i < trials; i++ {
		idx := RestrictedTournament(pop, candidate, 3, rnd)
		if idx == 2 {
			found2++
		}
	}
	if found2 > 5 {
		t.Fatalf("should rarely pick the farthest individual, got %d/%d", found2, trials)
	}
}

func TestAdaptiveSigma(t *testing.T) {
	pop := []Individual{
		{Genes: []float64{0, 0}},
		{Genes: []float64{5, 5}},
		{Genes: []float64{10, 10}},
	}
	s0 := AdaptiveSigma(pop, 3.0, 0)
	s100 := AdaptiveSigma(pop, 3.0, 100)
	if s100 >= s0 {
		t.Fatalf("sigma should decrease with generation: s0=%v, s100=%v", s0, s100)
	}
}
