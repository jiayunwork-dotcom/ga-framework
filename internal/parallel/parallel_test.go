package parallel

import (
	"math"
	"sync/atomic"
	"testing"
)

func sumFitness(genes []float64) float64 {
	s := 0.0
	for _, g := range genes {
		s += g
	}
	return s
}

func TestEvalAll(t *testing.T) {
	pop := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	results := EvalAll(pop, sumFitness, 2)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if math.Abs(results[0]-6) > 1e-10 {
		t.Fatalf("expected 6, got %v", results[0])
	}
	if math.Abs(results[2]-24) > 1e-10 {
		t.Fatalf("expected 24, got %v", results[2])
	}
}

func TestBatchEval(t *testing.T) {
	pop := make([][]float64, 10)
	for i := range pop {
		pop[i] = []float64{float64(i)}
	}
	ch := BatchEval(pop, sumFitness, 3)
	count := 0
	for range ch {
		count++
	}
	if count != 10 {
		t.Fatalf("expected 10 results, got %d", count)
	}
}

func TestPool_Submit(t *testing.T) {
	p := NewPool(4)
	defer p.Close()
	var counter int64
	for i := 0; i < 100; i++ {
		p.Submit(func() {
			atomic.AddInt64(&counter, 1)
		})
	}
	if counter != 100 {
		t.Fatalf("expected 100, got %d", counter)
	}
}

func TestPool_SubmitAsync(t *testing.T) {
	p := NewPool(2)
	defer p.Close()
	var counter int64
	dones := make([]<-chan struct{}, 50)
	for i := range dones {
		dones[i] = p.SubmitAsync(func() {
			atomic.AddInt64(&counter, 1)
		})
	}
	for _, done := range dones {
		<-done
	}
	if counter != 50 {
		t.Fatalf("expected 50, got %d", counter)
	}
}

func TestPool_MapEval(t *testing.T) {
	p := NewPool(2)
	defer p.Close()
	pop := [][]float64{{1, 1}, {2, 2}, {3, 3}}
	results := p.MapEval(pop, sumFitness)
	if math.Abs(results[0]-2) > 1e-10 || math.Abs(results[1]-4) > 1e-10 || math.Abs(results[2]-6) > 1e-10 {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestEvalAll_Empty(t *testing.T) {
	results := EvalAll(nil, sumFitness, 2)
	if results != nil {
		t.Fatalf("expected nil for empty input")
	}
}
