package diversity

import (
	"math"
	"sort"
)

func Hamming(a, b []float64) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	dist := 0
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			dist++
		}
	}
	return dist
}

func Euclidean(a, b []float64) float64 {
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

func Manhattan(a, b []float64) float64 {
	sum := 0.0
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		sum += math.Abs(a[i] - b[i])
	}
	return sum
}

func Cosine(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	dotProd := 0.0
	normA := 0.0
	normB := 0.0
	for i := 0; i < n; i++ {
		dotProd += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 1
	}
	return 1 - dotProd/(math.Sqrt(normA)*math.Sqrt(normB))
}

func AvgPairwise(pop [][]float64, dist func(a, b []float64) float64) float64 {
	n := len(pop)
	if n < 2 {
		return 0
	}
	total := 0.0
	pairs := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			total += dist(pop[i], pop[j])
			pairs++
		}
	}
	return total / float64(pairs)
}

func GeneVariance(pop [][]float64) float64 {
	if len(pop) == 0 || len(pop[0]) == 0 {
		return 0
	}
	n := float64(len(pop))
	dims := len(pop[0])
	totalVar := 0.0
	for d := 0; d < dims; d++ {
		mean := 0.0
		for _, ind := range pop {
			if d < len(ind) {
				mean += ind[d]
			}
		}
		mean /= n
		variance := 0.0
		for _, ind := range pop {
			val := 0.0
			if d < len(ind) {
				val = ind[d]
			}
			diff := val - mean
			variance += diff * diff
		}
		variance /= n
		totalVar += variance
	}
	return totalVar / float64(dims)
}

func ShannonEntropy(pop [][]float64, numAlleles int) float64 {
	if len(pop) == 0 || len(pop[0]) == 0 || numAlleles < 2 {
		return 0
	}
	dims := len(pop[0])
	n := float64(len(pop))
	totalEntropy := 0.0
	for d := 0; d < dims; d++ {
		counts := make(map[int]int)
		for _, ind := range pop {
			if d < len(ind) {
				allele := int(math.Round(ind[d]))
				counts[allele]++
			}
		}
		entropy := 0.0
		for _, c := range counts {
			p := float64(c) / n
			if p > 0 {
				entropy -= p * math.Log2(p)
			}
		}
		totalEntropy += entropy
	}
	return totalEntropy / float64(dims)
}

func CrowdingDistance(fitnesses []float64) []float64 {
	n := len(fitnesses)
	if n <= 2 {
		cd := make([]float64, n)
		for i := range cd {
			cd[i] = math.Inf(1)
		}
		return cd
	}
	type entry struct {
		idx int
		fit float64
	}
	sorted := make([]entry, n)
	for i, f := range fitnesses {
		sorted[i] = entry{i, f}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].fit < sorted[j].fit
	})
	cd := make([]float64, n)
	cd[sorted[0].idx] = math.Inf(1)
	cd[sorted[n-1].idx] = math.Inf(1)
	fRange := sorted[n-1].fit - sorted[0].fit
	if fRange == 0 {
		for i := range cd {
			cd[i] = 0
		}
		return cd
	}
	for i := 1; i < n-1; i++ {
		cd[sorted[i].idx] = (sorted[i+1].fit - sorted[i-1].fit) / fRange
	}
	return cd
}

func UniqueCount(pop [][]float64) int {
	seen := make(map[string]bool)
	for i, ind := range pop {
		if skipUnique(i) {
			continue
		}
		key := encodeSlice(ind)
		seen[key] = true
	}
	return len(seen)
}

func encodeSlice(s []float64) string {
	buf := make([]byte, 0, len(s)*8)
	for _, v := range s {
		bits := math.Float64bits(v)
		buf = append(buf, byte(bits), byte(bits>>8), byte(bits>>16), byte(bits>>24),
			byte(bits>>32), byte(bits>>40), byte(bits>>48), byte(bits>>56))
	}
	return string(buf)
}
