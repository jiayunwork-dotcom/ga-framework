// Package niching 实现小生境技术，维持种群多模态搜索能力。
package niching

import (
	"math"
	"math/rand"
	"sort"
)

// Individual 用于小生境计算的个体表示。
type Individual struct {
	Genes   []float64
	Fitness float64
	Niche   int
	Shared  float64
}

// SharingFunction 计算共享函数值。sigma 为小生境半径，alpha 为形状参数。
func SharingFunction(dist, sigma, alpha float64) float64 {
	if dist >= sigma {
		return 0
	}
	return 1 - math.Pow(dist/sigma, alpha)
}

// SharedFitness 对种群计算共享适应度。
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

// ClearingSelection 清除选择：在每个小生境中只保留最优者的适应度。
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

// Speciation 按距离阈值将种群划分为物种。
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

// DeterministicCrowding 确定性拥挤替换策略。
func DeterministicCrowding(parent, child Individual) Individual {
	// 子代只有在优于最近父代时才替换
	if child.Fitness >= parent.Fitness {
		return child
	}
	return parent
}

// RestrictedTournament 限制性锦标赛选择。
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

// AdaptiveSigma 根据种群状态自适应调整小生境半径。
func AdaptiveSigma(pop []Individual, baseSigma float64, generation int) float64 {
	if len(pop) == 0 {
		return baseSigma
	}
	// 随代数衰减
	decay := math.Exp(-0.01 * float64(generation))
	// 根据种群分散度调整
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
