// Package crossover 实现交叉算子。
package crossover

import (
	"math/rand"

	"ga-framework/internal/genome"
)

// SinglePoint 在 [1, n-1] 之间选切点，交换两个父代切点之后的基因，返回两个子代。
// 若基因长度小于 2，则原样返回副本。rnd 用于选择切点。
func SinglePoint(a, b genome.Genome, rnd *rand.Rand) (genome.Genome, genome.Genome) {
	n := len(a.Genes)
	if n < 2 {
		return clone(a), clone(b)
	}
	point := 1 + rnd.Intn(n-1)
	c1 := make([]int, n)
	c2 := make([]int, n)
	for i := 0; i < n; i++ {
		if i < point {
			c1[i] = a.Genes[i]
			c2[i] = b.Genes[i]
		} else {
			c1[i] = b.Genes[i]
			c2[i] = a.Genes[i]
		}
	}
	return genome.Genome{Genes: c1}, genome.Genome{Genes: c2}
}

func clone(g genome.Genome) genome.Genome {
	c := make([]int, len(g.Genes))
	copy(c, g.Genes)
	return genome.Genome{Genes: c}
}
