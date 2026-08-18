// Package mutation 实现变异算子。
package mutation

import (
	"math/rand"

	"ga-framework/internal/genome"
)

// Mutate 以概率 rate 翻转每个基因（0<->1）。返回新个体，不修改入参。
func Mutate(g genome.Genome, rate float64, rnd *rand.Rand) genome.Genome {
	out := make([]int, len(g.Genes))
	copy(out, g.Genes)
	for i := range out {
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
