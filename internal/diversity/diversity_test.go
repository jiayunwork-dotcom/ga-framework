package diversity

import (
	"math"
	"testing"
)

func TestHamming(t *testing.T) {
	a := []float64{1, 0, 1, 1, 0}
	b := []float64{0, 0, 1, 0, 1}
	d := Hamming(a, b)
	if d != 3 {
		t.Fatalf("expected hamming 3, got %d", d)
	}
}

func TestEuclidean(t *testing.T) {
	a := []float64{0, 0}
	b := []float64{3, 4}
	d := Euclidean(a, b)
	if math.Abs(d-5.0) > 1e-10 {
		t.Fatalf("expected 5, got %v", d)
	}
}

func TestManhattan(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 1, 0}
	d := Manhattan(a, b)
	if math.Abs(d-7.0) > 1e-10 {
		t.Fatalf("expected 7, got %v", d)
	}
}

func TestCosine_Identical(t *testing.T) {
	a := []float64{1, 2, 3}
	d := Cosine(a, a)
	if math.Abs(d) > 1e-10 {
		t.Fatalf("identical vectors should have cosine distance 0, got %v", d)
	}
}

func TestAvgPairwise(t *testing.T) {
	pop := [][]float64{{0, 0}, {1, 0}, {0, 1}}
	avg := AvgPairwise(pop, Euclidean)
	// distances: 1, 1, sqrt(2) => avg = (1+1+1.414)/3 ≈ 1.138
	if avg < 1.0 || avg > 1.5 {
		t.Fatalf("unexpected avg pairwise: %v", avg)
	}
}

func TestGeneVariance(t *testing.T) {
	pop := [][]float64{{0, 10}, {10, 0}, {5, 5}}
	v := GeneVariance(pop)
	if v < 10 {
		t.Fatalf("expected high variance, got %v", v)
	}
}

func TestShannonEntropy(t *testing.T) {
	pop := [][]float64{{1, 1}, {1, 1}, {1, 1}}
	e := ShannonEntropy(pop, 2)
	if e != 0 {
		t.Fatalf("identical population should have entropy 0, got %v", e)
	}
	pop2 := [][]float64{{0, 0}, {1, 1}, {0, 1}, {1, 0}}
	e2 := ShannonEntropy(pop2, 2)
	if e2 <= 0 {
		t.Fatalf("diverse population should have positive entropy, got %v", e2)
	}
}

func TestCrowdingDistance(t *testing.T) {
	fitnesses := []float64{1, 3, 5, 7, 9}
	cd := CrowdingDistance(fitnesses)
	if !math.IsInf(cd[0], 1) || !math.IsInf(cd[4], 1) {
		t.Fatalf("boundary individuals should have infinite CD")
	}
	for i := 1; i < 4; i++ {
		if cd[i] <= 0 {
			t.Fatalf("interior CD should be positive, got %v at %d", cd[i], i)
		}
	}
}

func TestUniqueCount(t *testing.T) {
	pop := [][]float64{{1, 2}, {1, 2}, {3, 4}, {1, 2}}
	u := UniqueCount(pop)
	if u != 2 {
		t.Fatalf("expected 2 unique, got %d", u)
	}
}
