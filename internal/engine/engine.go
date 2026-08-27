package engine

import (
	"errors"
	"math/rand"

	"ga-framework/internal/crossover"
	"ga-framework/internal/genome"
	"ga-framework/internal/mutation"
	"ga-framework/internal/selection"
)

type Config struct {
	Size        int
	Genes       int
	Generations int
	TournamentK int
	MutateRate  float64
	Elite       int
	Seed        int64
}

type Result struct {
	Best        genome.Genome
	BestFit     float64
	Generations int
}

func Run(cfg Config) (Result, error) {
	if cfg.Size < 1 || cfg.Genes < 1 || cfg.Generations < 1 {
		return Result{}, errors.New("invalid config: size/genes/generations must be >= 1")
	}
	k := cfg.TournamentK
	if k < 1 {
		k = 2
	}
	elite := cfg.Elite
	if elite < 0 {
		elite = 0
	}
	if elite > cfg.Size {
		elite = cfg.Size
	}
	rnd := rand.New(rand.NewSource(cfg.Seed))
	pop := make([]genome.Genome, cfg.Size)
	for i := range pop {
		g := genome.NewRandom(cfg.Genes, rnd)
		g.Fitness = float64(g.Ones())
		pop[i] = g
	}
	best := pop[0]
	for gen := 1; gen <= cfg.Generations; gen++ {
		for i := range pop {
			pop[i].Fitness = float64(pop[i].Ones())
			if pop[i].Fitness > best.Fitness {
				best = pop[i]
			}
		}
		next := make([]genome.Genome, 0, cfg.Size)
		for _, e := range pickElite(pop, elite) {
			next = append(next, cloneGenome(e))
		}
		for len(next) < cfg.Size {
			ia := selection.Tournament(pop, k, rnd)
			ib := selection.Tournament(pop, k, rnd)
			c1, c2 := crossover.SinglePoint(pop[ia], pop[ib], rnd)
			c1 = mutation.Mutate(c1, cfg.MutateRate, rnd)
			c2 = mutation.Mutate(c2, cfg.MutateRate, rnd)
			next = append(next, c1, c2)
		}
		pop = next[:cfg.Size]
	}
	return Result{Best: best, BestFit: best.Fitness, Generations: cfg.Generations}, nil
}

func pickElite(pop []genome.Genome, n int) []genome.Genome {
	idx := make([]int, len(pop))
	for i := range idx {
		idx[i] = i
	}
	for i := 0; i < n && i < len(pop); i++ {
		max := i
		for j := i + 1; j < len(pop); j++ {
			if pop[idx[j]].Fitness > pop[idx[max]].Fitness {
				max = j
			}
		}
		idx[i], idx[max] = idx[max], idx[i]
	}
	out := make([]genome.Genome, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, pop[idx[i]])
	}
	return out
}

func cloneGenome(g genome.Genome) genome.Genome {
	c := make([]int, len(g.Genes))
	copy(c, g.Genes)
	return genome.Genome{Genes: c, Fitness: g.Fitness}
}
