package genome

import (
	"math/rand"
	"testing"
)

func newRand() *rand.Rand { return rand.New(rand.NewSource(1)) }

func TestOnes(t *testing.T) {
	g := Genome{Genes: []int{1, 0, 1, 1, 0}}
	if g.Ones() != 3 {
		t.Fatalf("expected 3 ones, got %d", g.Ones())
	}
}

func TestNewRandom_Size(t *testing.T) {
	g := NewRandom(10, newRand())
	if len(g.Genes) != 10 {
		t.Fatalf("expected length 10, got %d", len(g.Genes))
	}
	for _, v := range g.Genes {
		if v != 0 && v != 1 {
			t.Fatalf("gene must be 0 or 1, got %d", v)
		}
	}
}
