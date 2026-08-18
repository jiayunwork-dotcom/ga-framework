package archive

import (
	"math"
	"testing"
)

func TestAdd_NonDominated(t *testing.T) {
	a := New(10)
	ok := a.Add(Solution{Objectives: []float64{5, 3}})
	if !ok {
		t.Fatalf("should accept first solution")
	}
	ok = a.Add(Solution{Objectives: []float64{3, 5}})
	if !ok {
		t.Fatalf("should accept non-dominated solution")
	}
	if a.Size() != 2 {
		t.Fatalf("expected 2, got %d", a.Size())
	}
}

func TestAdd_Dominated(t *testing.T) {
	a := New(10)
	a.Add(Solution{Objectives: []float64{5, 5}})
	ok := a.Add(Solution{Objectives: []float64{3, 3}})
	if ok {
		t.Fatalf("dominated solution should be rejected")
	}
}

func TestAdd_RemovesDominated(t *testing.T) {
	a := New(10)
	a.Add(Solution{Objectives: []float64{3, 3}})
	a.Add(Solution{Objectives: []float64{4, 4}})
	a.Add(Solution{Objectives: []float64{5, 5}})
	if a.Size() != 1 {
		t.Fatalf("should only have dominating solution, got %d", a.Size())
	}
}

func TestBest(t *testing.T) {
	a := New(10)
	a.Add(Solution{Objectives: []float64{3, 8}})
	a.Add(Solution{Objectives: []float64{8, 3}})
	best, ok := a.Best(0)
	if !ok || best.Objectives[0] != 8 {
		t.Fatalf("expected best obj[0]=8, got %v", best)
	}
	best1, _ := a.Best(1)
	if best1.Objectives[1] != 8 {
		t.Fatalf("expected best obj[1]=8, got %v", best1)
	}
}

func TestNonDominatedSort(t *testing.T) {
	sols := []Solution{
		{Objectives: []float64{5, 5}},
		{Objectives: []float64{3, 3}},
		{Objectives: []float64{4, 2}},
		{Objectives: []float64{1, 1}},
	}
	ranks := NonDominatedSort(sols)
	// sol[0] dominates sol[1], sol[2], sol[3]
	if ranks[0] != 0 {
		t.Fatalf("sol[0] should be rank 0, got %d", ranks[0])
	}
	if ranks[3] >= ranks[0] && ranks[3] == 0 {
		t.Fatalf("sol[3] should not be rank 0")
	}
}

func TestHyperVolume(t *testing.T) {
	a := New(10)
	a.Add(Solution{Objectives: []float64{3, 3}})
	a.Add(Solution{Objectives: []float64{1, 5}})
	ref := []float64{6, 6}
	hv := a.HyperVolume(ref)
	if hv <= 0 {
		t.Fatalf("hypervolume should be positive, got %v", hv)
	}
}

func TestSpread(t *testing.T) {
	a := New(10)
	a.Add(Solution{Objectives: []float64{1, 5}})
	a.Add(Solution{Objectives: []float64{3, 3}})
	a.Add(Solution{Objectives: []float64{5, 1}})
	s := a.Spread()
	if math.IsNaN(s) {
		t.Fatalf("spread should not be NaN")
	}
}

func TestCapacity_Trim(t *testing.T) {
	a := New(3)
	for i := 0; i < 10; i++ {
		a.Add(Solution{Objectives: []float64{float64(i), float64(10 - i)}})
	}
	if a.Size() > 3 {
		t.Fatalf("archive should respect capacity, got %d", a.Size())
	}
}

func TestClear(t *testing.T) {
	a := New(10)
	a.Add(Solution{Objectives: []float64{1, 2}})
	a.Clear()
	if a.Size() != 0 {
		t.Fatalf("expected 0 after clear, got %d", a.Size())
	}
}

func TestContains(t *testing.T) {
	a := New(10)
	s := Solution{Genes: []float64{1, 2, 3}, Objectives: []float64{5, 5}}
	a.Add(s)
	if !a.Contains(s) {
		t.Fatalf("archive should contain the added solution")
	}
	if a.Contains(Solution{Genes: []float64{4, 5, 6}}) {
		t.Fatalf("archive should not contain other solution")
	}
}
