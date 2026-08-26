package niching

import (
	"math"
	"math/rand"
	"sort"
)

type Individual struct {
	Genes   []float64
	Fitness float64
	Niche   int
	Shared  float64
}

func SharingFunction(dist, sigma, alpha float64) float64 {
	if dist >= sigma {
		return 0
	}
	return 1 - math.Pow(dist/sigma, alpha)
}

func SharedFitness(pop []Individual, sigma, alpha float64) []float64 {
	n := len(pop)
	shared := make([]float64, n)
	for i := range pop {
		nicheCount := 0.0
		for j := range pop {
			d := euclidean(pop[i].Genes, pop[j].Genes)
			nicheCount += SharingFunction(d, sigma, alpha)
		}
		if nicheCount == 0 {
			nicheCount = 1
		}
		shared[i] = pop[i].Fitness / nicheCount
	}
	return shared
}

func ClearingSelection(pop []Individual, sigma float64, capacity int) []Individual {
	n := len(pop)
	sorted := make([]Individual, n)
	copy(sorted, pop)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Fitness > sorted[j].Fitness
	})
	winners := make([]bool, n)
	nicheCount := make([]int, n)
	for i := range sorted {
		if winners[i] {
			continue
		}
		count := 0
		for j := range sorted {
			if i == j {
				continue
			}
			d := euclidean(sorted[i].Genes, sorted[j].Genes)
			if d < sigma {
				count++
				if count >= capacity {
					sorted[j].Fitness = 0
				}
			}
		}
		winners[i] = true
		nicheCount[i] = count
	}
	return sorted
}

func Speciation(pop []Individual, threshold float64) [][]int {
	n := len(pop)
	assigned := make([]bool, n)
	species := make([][]int, 0)
	for i := 0; i < n; i++ {
		if assigned[i] {
			continue
		}
		sp := []int{i}
		assigned[i] = true
		for j := i + 1; j < n; j++ {
			if assigned[j] {
				continue
			}
			d := euclidean(pop[i].Genes, pop[j].Genes)
			if d < threshold {
				sp = append(sp, j)
				assigned[j] = true
			}
		}
		species = append(species, sp)
	}
	return species
}

func DeterministicCrowding(parent, child Individual) Individual {
	if child.Fitness >= parent.Fitness {
		return child
	}
	return parent
}

func RestrictedTournament(pop []Individual, candidate Individual, windowSize int, rnd *rand.Rand) int {
	n := len(pop)
	if windowSize > n {
		windowSize = n
	}
	bestDist := math.Inf(1)
	bestIdx := 0
	for i := 0; i < windowSize; i++ {
		idx := rnd.Intn(n)
		d := euclidean(pop[idx].Genes, candidate.Genes)
		if d < bestDist {
			bestDist = d
			bestIdx = idx
		}
	}
	return bestIdx
}

func AdaptiveSigma(pop []Individual, baseSigma float64, generation int) float64 {
	if len(pop) == 0 {
		return baseSigma
	}
	decay := math.Exp(-0.01 * float64(generation))
	avgDist := avgPairDist(pop)
	return baseSigma * decay * (1 + avgDist/baseSigma) / 2
}

func avgPairDist(pop []Individual) float64 {
	n := len(pop)
	if n < 2 {
		return 0
	}
	total := 0.0
	pairs := 0
	step := 1
	if n > 50 {
		step = n / 25
	}
	for i := 0; i < n; i += step {
		for j := i + step; j < n; j += step {
			total += euclidean(pop[i].Genes, pop[j].Genes)
			pairs++
		}
	}
	if pairs == 0 {
		return 0
	}
	return total / float64(pairs)
}

func euclidean(a, b []float64) float64 {
	sum := 0.0
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}
