// Package selection 实现个体选择算子。
package selection

import (
	"math/rand"

	"ga-framework/internal/genome"
)

// Tournament 在种群中随机抽取 k 个个体，返回适应度最高的下标。
// k 小于 1 时按 1 处理；rnd 必须非 nil。
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
