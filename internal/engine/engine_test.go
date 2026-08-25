package engine

import "testing"

func TestEngine_RunsToConvergence(t *testing.T) {
	res, err := Run(Config{
		Size: 40, Genes: 6, Generations: 300, TournamentK: 2,
		MutateRate: 0.05, Elite: 2, Seed: 42,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BestFit != float64(6) {
		t.Fatalf("expected best fitness 6, got %v", res.BestFit)
	}
}

func TestEngine_InvalidConfig(t *testing.T) {
	if _, err := Run(Config{Size: 0}); err == nil {
		t.Fatalf("expected error for zero size")
	}
	if _, err := Run(Config{Size: 10, Genes: 0}); err == nil {
		t.Fatalf("expected error for zero genes")
	}
}
