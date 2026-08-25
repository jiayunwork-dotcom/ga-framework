package population

import (
	"math"
	"math/rand"
	"testing"
)

func rng() *rand.Rand { return rand.New(rand.NewSource(42)) }

func TestRandomBinary_Size(t *testing.T) {
	p := RandomBinary(20, 8, rng())
	if p.Size() != 20 {
		t.Fatalf("expected 20 individuals, got %d", p.Size())
	}
	for _, ind := range p.Inds {
		if len(ind.Genes) != 8 {
			t.Fatalf("expected 8 genes, got %d", len(ind.Genes))
		}
		for _, g := range ind.Genes {
			if g != 0 && g != 1 {
				t.Fatalf("binary gene must be 0 or 1, got %v", g)
			}
		}
	}
}

func TestRandomReal_Bounds(t *testing.T) {
	p := RandomReal(30, 5, -5.0, 5.0, rng())
	for _, ind := range p.Inds {
		for _, g := range ind.Genes {
			if g < -5.0 || g > 5.0 {
				t.Fatalf("gene out of bounds: %v", g)
			}
		}
	}
}

func TestRandomPermutation_Valid(t *testing.T) {
	p := RandomPermutation(10, 6, rng())
	for _, ind := range p.Inds {
		seen := make(map[float64]bool)
		for _, g := range ind.Genes {
			if g < 0 || g >= 6 {
				t.Fatalf("gene out of range: %v", g)
			}
			if seen[g] {
				t.Fatalf("duplicate gene in permutation: %v", g)
			}
			seen[g] = true
		}
	}
}

func TestBest_Worst(t *testing.T) {
	p := New(5)
	p.Add(Individual{Genes: []float64{1}, Fitness: 10})
	p.Add(Individual{Genes: []float64{2}, Fitness: 3})
	p.Add(Individual{Genes: []float64{3}, Fitness: 20})
	p.Add(Individual{Genes: []float64{4}, Fitness: 1})
	best := p.Best()
	if best.Fitness != 20 {
		t.Fatalf("expected best fitness 20, got %v", best.Fitness)
	}
	worst := p.Worst()
	if worst.Fitness != 1 {
		t.Fatalf("expected worst fitness 1, got %v", worst.Fitness)
	}
}

func TestAvgFitness(t *testing.T) {
	p := New(3)
	p.Add(Individual{Fitness: 10})
	p.Add(Individual{Fitness: 20})
	p.Add(Individual{Fitness: 30})
	avg := p.AvgFitness()
	if math.Abs(avg-20.0) > 1e-10 {
		t.Fatalf("expected avg 20, got %v", avg)
	}
}

func TestSortByFitness(t *testing.T) {
	p := New(3)
	p.Add(Individual{Fitness: 5})
	p.Add(Individual{Fitness: 15})
	p.Add(Individual{Fitness: 10})
	p.SortByFitness()
	if p.Inds[0].Fitness != 15 || p.Inds[2].Fitness != 5 {
		t.Fatalf("sort order incorrect: %v", p.Inds)
	}
}

func TestTopN(t *testing.T) {
	p := New(5)
	p.Add(Individual{Fitness: 1})
	p.Add(Individual{Fitness: 5})
	p.Add(Individual{Fitness: 3})
	p.Add(Individual{Fitness: 4})
	p.Add(Individual{Fitness: 2})
	top := p.TopN(3)
	if len(top) != 3 {
		t.Fatalf("expected 3, got %d", len(top))
	}
	if top[0].Fitness != 5 || top[1].Fitness != 4 || top[2].Fitness != 3 {
		t.Fatalf("top3 wrong: %v", top)
	}
}

func TestClone_Independence(t *testing.T) {
	p := New(2)
	p.Add(Individual{Genes: []float64{1, 2, 3}, Fitness: 10})
	p.Add(Individual{Genes: []float64{4, 5, 6}, Fitness: 20})
	cp := p.Clone()
	cp.Inds[0].Fitness = 999
	if p.Inds[0].Fitness == 999 {
		t.Fatalf("clone should be independent")
	}
}

func TestDiversity(t *testing.T) {
	p := New(3)
	p.Add(Individual{Genes: []float64{0, 0}})
	p.Add(Individual{Genes: []float64{1, 0}})
	p.Add(Individual{Genes: []float64{0, 1}})
	d := p.Diversity()
	if d <= 0 {
		t.Fatalf("diversity should be > 0, got %v", d)
	}
}
