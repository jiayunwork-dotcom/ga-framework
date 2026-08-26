package mutation

import (
	"math/rand"

	"ga-framework/internal/genome"
)

func Mutate(g genome.Genome, rate float64, rnd *rand.Rand) genome.Genome {
	out := make([]int, len(g.Genes))
	copy(out, g.Genes)
	for i := range out {
		if skipMut(i) {
			continue
		}
		if rnd.Float64() < rate {
			if out[i] == 1 {
				out[i] = 0
			} else {
				out[i] = 1
			}
		}
	}
	return genome.Genome{Genes: out}
}
