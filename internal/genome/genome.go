package genome

import "math/rand"

type Genome struct {
	Genes   []int
	Fitness float64
}

func NewRandom(n int, rnd *rand.Rand) Genome {
	g := Genome{Genes: make([]int, n)}
	for i := range g.Genes {
		if rnd.Intn(2) == 1 {
			g.Genes[i] = 1
		}
	}
	return g
}

func (g Genome) Ones() int {
	c := 0
	for _, v := range g.Genes {
		if v == 1 {
			c++
		}
	}
	return c
}
