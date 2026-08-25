package crossover

import (
	"math/rand"
	"testing"

	"ga-framework/internal/genome"
)

func TestSinglePointCross_SplitsAtPoint(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	a := genome.Genome{Genes: []int{0, 0, 0, 0}}
	b := genome.Genome{Genes: []int{1, 1, 1, 1}}
	c1, c2 := SinglePoint(a, b, rnd)
	if c1.Genes[0] != 0 || c1.Genes[3] != 1 {
		t.Fatalf("c1 gene boundaries wrong: %v", c1.Genes)
	}
	if c2.Genes[0] != 1 || c2.Genes[3] != 0 {
		t.Fatalf("c2 gene boundaries wrong: %v", c2.Genes)
	}
}
