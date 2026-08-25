package stats

import (
	"math"
	"testing"
)

func TestCollector_AddAndCount(t *testing.T) {
	c := NewCollector()
	c.Add(Record{Generation: 0, BestFit: 10})
	c.Add(Record{Generation: 1, BestFit: 15})
	if c.Count() != 2 {
		t.Fatalf("expected 2, got %d", c.Count())
	}
}

func TestBestEver(t *testing.T) {
	c := NewCollector()
	c.Add(Record{BestFit: 5})
	c.Add(Record{BestFit: 20})
	c.Add(Record{BestFit: 12})
	if c.BestEver() != 20 {
		t.Fatalf("expected 20, got %v", c.BestEver())
	}
}

func TestStagnation(t *testing.T) {
	c := NewCollector()
	c.Add(Record{BestFit: 10})
	c.Add(Record{BestFit: 15})
	c.Add(Record{BestFit: 15})
	c.Add(Record{BestFit: 15})
	if c.Stagnation() != 2 {
		t.Fatalf("expected stagnation 2, got %d", c.Stagnation())
	}
}

func TestImprovementRate(t *testing.T) {
	c := NewCollector()
	c.Add(Record{BestFit: 10})
	c.Add(Record{BestFit: 12})
	c.Add(Record{BestFit: 15})
	rate := c.ImprovementRate(3)
	expected := (15 - 10) / 10.0
	if math.Abs(rate-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, rate)
	}
}

func TestSummarize(t *testing.T) {
	c := NewCollector()
	c.Add(Record{BestFit: 5, Diversity: 1.0})
	c.Add(Record{BestFit: 10, Diversity: 0.8})
	c.Add(Record{BestFit: 10, Diversity: 0.6})
	s := c.Summarize()
	if s.FinalBest != 10 {
		t.Fatalf("expected final best 10, got %v", s.FinalBest)
	}
	if s.ConvergenceGen != 1 {
		t.Fatalf("expected convergence gen 1, got %d", s.ConvergenceGen)
	}
	if s.TotalGenerations != 3 {
		t.Fatalf("expected 3 generations, got %d", s.TotalGenerations)
	}
}

func TestPercentile(t *testing.T) {
	c := NewCollector()
	for i := 1; i <= 100; i++ {
		c.Add(Record{BestFit: float64(i)})
	}
	p50 := c.Percentile(50)
	if p50 < 49 || p50 > 51 {
		t.Fatalf("P50 should be around 50, got %v", p50)
	}
}

func TestMovingAvg(t *testing.T) {
	c := NewCollector()
	for i := 0; i < 10; i++ {
		c.Add(Record{BestFit: float64(i)})
	}
	ma := c.MovingAvg(3)
	if len(ma) != 10 {
		t.Fatalf("expected 10 values, got %d", len(ma))
	}
	if math.Abs(ma[5]-4.0) > 1e-10 {
		t.Fatalf("expected ma[5]=4, got %v", ma[5])
	}
}

func TestLast(t *testing.T) {
	c := NewCollector()
	_, ok := c.Last()
	if ok {
		t.Fatalf("expected no last record")
	}
	c.Add(Record{Generation: 5, BestFit: 99})
	r, ok := c.Last()
	if !ok || r.BestFit != 99 {
		t.Fatalf("expected last record with BestFit=99")
	}
}
