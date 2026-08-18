package selection

import (
	"math/rand"
	"testing"

	"ga-framework/internal/genome"
)

func TestTournamentSelect_ReturnsBestAmongSampled(t *testing.T) {
	pop := []genome.Genome{
		{Genes: []int{0}, Fitness: 0.1},
		{Genes: []int{1}, Fitness: 0.9},
		{Genes: []int{0}, Fitness: 0.3},
	}
	k := 3
	rep := rand.New(rand.NewSource(1))
	sampled := []int{rep.Intn(len(pop))}
	best := sampled[0]
	for i := 1; i < k; i++ {
		c := rep.Intn(len(pop))
		sampled = append(sampled, c)
		if pop[c].Fitness > pop[best].Fitness {
			best = c
		}
	}
	call := rand.New(rand.NewSource(1))
	got := Tournament(pop, k, call)
	if got != best {
		t.Fatalf("expected sampled-best index %d, got %d", best, got)
	}
	inSet, maxFit := false, pop[sampled[0]].Fitness
	for _, idx := range sampled {
		if idx == got {
			inSet = true
		}
		if pop[idx].Fitness > maxFit {
			maxFit = pop[idx].Fitness
		}
	}
	if !inSet {
		t.Fatalf("selected index %d not in sampled set %v", got, sampled)
	}
	if pop[got].Fitness != maxFit {
		t.Fatalf("selected fitness %v != sampled max %v", pop[got].Fitness, maxFit)
	}
}
