package fitness

import (
	"math"
	"testing"
)

func TestOneMax(t *testing.T) {
	genes := []float64{1, 0, 1, 1, 0, 1}
	got := OneMax(genes)
	if got != 4 {
		t.Fatalf("OneMax: expected 4, got %v", got)
	}
}

func TestSphere_Origin(t *testing.T) {
	genes := []float64{0, 0, 0}
	got := Sphere(genes)
	if got != 0 {
		t.Fatalf("Sphere at origin should be 0, got %v", got)
	}
}

func TestRastrigin_Origin(t *testing.T) {
	genes := []float64{0, 0, 0}
	got := Rastrigin(genes)
	// f(0,0,0) = -(10*3 + sum(0 - 10*1)) = -(30-30) = 0
	if math.Abs(got) > 1e-6 {
		t.Fatalf("Rastrigin at origin expected 0, got %v", got)
	}
}

func TestRosenbrock_Minimum(t *testing.T) {
	genes := []float64{1, 1, 1}
	got := Rosenbrock(genes)
	if math.Abs(got) > 1e-10 {
		t.Fatalf("Rosenbrock at (1,1,1) should be 0, got %v", got)
	}
}

func TestAckley_Origin(t *testing.T) {
	genes := []float64{0, 0}
	got := Ackley(genes)
	if math.Abs(got) > 1e-10 {
		t.Fatalf("Ackley at origin should be 0, got %v", got)
	}
}

func TestGriewank_Origin(t *testing.T) {
	genes := []float64{0, 0, 0}
	got := Griewank(genes)
	if math.Abs(got) > 1e-10 {
		t.Fatalf("Griewank at origin should be 0, got %v", got)
	}
}

func TestWeightedSum(t *testing.T) {
	f := WeightedSum([]Func{Sphere, OneMax}, []float64{0.5, 2.0})
	genes := []float64{1, 0, 1}
	got := f(genes)
	expected := 0.5*Sphere(genes) + 2.0*OneMax(genes)
	if math.Abs(got-expected) > 1e-10 {
		t.Fatalf("WeightedSum: expected %v, got %v", expected, got)
	}
}

func TestPenalty(t *testing.T) {
	constraint := func(g []float64) float64 {
		sum := 0.0
		for _, x := range g {
			sum += x
		}
		return sum - 3.0 // violated if sum > 3
	}
	f := Penalty(OneMax, constraint, 10.0)
	genes := []float64{1, 1, 1, 1} // sum=4, violation=1
	got := f(genes)
	base := OneMax(genes)
	expected := base - 10.0*1.0*1.0
	if math.Abs(got-expected) > 1e-10 {
		t.Fatalf("Penalty: expected %v, got %v", expected, got)
	}
}
