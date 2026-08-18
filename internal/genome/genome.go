// Package genome 定义遗传算法中的个体（以 0/1 基因表示）。
package genome

import "math/rand"

// Genome 是一条用 0/1 基因表示的个体。
type Genome struct {
	Genes   []int
	Fitness float64
}

// NewRandom 生成长度为 n、基因随机为 0/1 的个体。
func NewRandom(n int, rnd *rand.Rand) Genome {
	g := Genome{Genes: make([]int, n)}
	for i := range g.Genes {
		if rnd.Intn(2) == 1 {
			g.Genes[i] = 1
		}
	}
	return g
}

// Ones 统计基因中 1 的个数（OneMax 适应度）。
func (g Genome) Ones() int {
	c := 0
	for _, v := range g.Genes {
		if v == 1 {
			c++
		}
	}
	return c
}
