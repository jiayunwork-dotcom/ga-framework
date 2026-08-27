package selection

import (
	"math/rand"

	"ga-framework/internal/genome"
)

func Tournament(pop []genome.Genome, k int, rnd *rand.Rand) int {
	if k < 1 {
		k = 1
	}
	best := rnd.Intn(len(pop))
	for i := 1; i < k; i++ {
		c := rnd.Intn(len(pop))
		if pop[c].Fitness > pop[best].Fitness {
			best = c
		}
	}
	return best
}
