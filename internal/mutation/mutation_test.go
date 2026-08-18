package mutation

import (
	"math/rand"
	"testing"

	"ga-framework/internal/genome"
)

func TestMutate_FlipsAllAtRate1(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))
	g := genome.Genome{Genes: []int{0, 1, 0, 1}}
	out := Mutate(g, 1.0, rnd)
	want := []int{1, 0, 1, 0}
	for i := range want {
		if out.Genes[i] != want[i] {
			t.Fatalf("rate=1 should flip all bits: %v", out.Genes)
		}
	}
}

func TestMutate_NoFlipAtRate0(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))
	g := genome.Genome{Genes: []int{0, 1, 0, 1}}
	out := Mutate(g, 0.0, rnd)
	for i := range g.Genes {
		if out.Genes[i] != g.Genes[i] {
			t.Fatalf("rate=0 must not change genes: %v", out.Genes)
		}
	}
}
