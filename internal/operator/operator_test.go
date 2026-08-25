package operator

import (
	"math/rand"
	"testing"
)

func rng() *rand.Rand { return rand.New(rand.NewSource(42)) }

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	err := r.RegisterCrossover("test_cx", UniformCrossover)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	err = r.RegisterCrossover("test_cx", UniformCrossover)
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
	f, err := r.GetCrossover("test_cx")
	if err != nil || f == nil {
		t.Fatalf("get failed: %v", err)
	}
	_, err = r.GetCrossover("nonexistent")
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestUniformCrossover_Length(t *testing.T) {
	p1 := []float64{1, 2, 3, 4, 5}
	p2 := []float64{6, 7, 8, 9, 10}
	c1, c2 := UniformCrossover(p1, p2, rng())
	if len(c1) != 5 || len(c2) != 5 {
		t.Fatalf("children length mismatch")
	}
	for i := range c1 {
		if c1[i] != p1[i] && c1[i] != p2[i] {
			t.Fatalf("child1 gene %d not from either parent: %v", i, c1[i])
		}
	}
}

func TestTwoPointCrossover_Segment(t *testing.T) {
	p1 := []float64{1, 1, 1, 1, 1}
	p2 := []float64{2, 2, 2, 2, 2}
	c1, c2 := TwoPointCrossover(p1, p2, rng())
	if len(c1) != 5 || len(c2) != 5 {
		t.Fatalf("children length mismatch")
	}
	for i := range c1 {
		if c1[i] != 1 && c1[i] != 2 {
			t.Fatalf("unexpected gene value: %v", c1[i])
		}
		if c2[i] != 1 && c2[i] != 2 {
			t.Fatalf("unexpected gene value: %v", c2[i])
		}
	}
}

func TestGaussianMutation_Applies(t *testing.T) {
	genes := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	mutated := GaussianMutation(genes, 1.0, rng())
	changed := false
	for i := range genes {
		if mutated[i] != genes[i] {
			changed = true
			break
		}
	}
	if !changed {
		t.Fatalf("expected some mutation with rate=1.0")
	}
}

func TestSwapMutation_Preserves(t *testing.T) {
	genes := []float64{0, 1, 2, 3, 4}
	mutated := SwapMutation(genes, 1.0, rng())
	sum1, sum2 := 0.0, 0.0
	for _, g := range genes {
		sum1 += g
	}
	for _, g := range mutated {
		sum2 += g
	}
	if sum1 != sum2 {
		t.Fatalf("swap mutation changed sum: %v vs %v", sum1, sum2)
	}
}

func TestRouletteSelection_Index(t *testing.T) {
	fitnesses := []float64{0, 0, 10, 0, 0}
	idx := RouletteSelection(fitnesses, rng())
	if idx != 2 {
		t.Fatalf("expected index 2 (highest fitness), got %d", idx)
	}
}

func TestRankSelection_Index(t *testing.T) {
	fitnesses := []float64{1, 2, 3, 4, 5}
	counts := make(map[int]int)
	r := rng()
	for i := 0; i < 1000; i++ {
		idx := RankSelection(fitnesses, r)
		counts[idx]++
	}
	if counts[4] < counts[0] {
		t.Fatalf("rank selection should favor higher fitness: counts=%v", counts)
	}
}

func TestDefaultRegistry_HasOperators(t *testing.T) {
	r := DefaultRegistry()
	if len(r.ListCrossovers()) < 2 {
		t.Fatalf("expected at least 2 crossovers")
	}
	if len(r.ListMutations()) < 2 {
		t.Fatalf("expected at least 2 mutations")
	}
	if len(r.ListSelections()) < 2 {
		t.Fatalf("expected at least 2 selections")
	}
}
