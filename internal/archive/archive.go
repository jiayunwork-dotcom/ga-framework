// Package archive 实现解归档，支持帕累托前沿存储与管理。
package archive

import (
	"math"
	"sort"
	"sync"
)

// Solution 归档中的解。
type Solution struct {
	Genes      []float64
	Objectives []float64
	Rank       int
}

// Archive 管理非支配解集。
type Archive struct {
	mu       sync.Mutex
	capacity int
	items    []Solution
}

// New 创建指定容量的归档。
func New(capacity int) *Archive {
	return &Archive{
		capacity: capacity,
		items:    make([]Solution, 0, capacity),
	}
}

// Size 返回归档中的解数量。
func (a *Archive) Size() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.items)
}

// Items 返回归档中所有解的副本。
func (a *Archive) Items() []Solution {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Solution, len(a.items))
	copy(out, a.items)
	return out
}

// Add 尝试将解加入归档。如果被归档中某解支配则拒绝。
func (a *Archive) Add(sol Solution) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	// 检查是否被已有解支配
	for _, existing := range a.items {
		if dominates(existing.Objectives, sol.Objectives) {
			return false
		}
	}
	// 移除被新解支配的解
	kept := a.items[:0]
	for _, existing := range a.items {
		if !dominates(sol.Objectives, existing.Objectives) {
			kept = append(kept, existing)
		}
	}
	a.items = append(kept, sol)
	// 超容量时裁剪（移除拥挤距离最小者）
	if len(a.items) > a.capacity {
		a.trim()
	}
	return true
}

// Clear 清空归档。
func (a *Archive) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.items = a.items[:0]
}

// Contains 检查解是否已在归档中（基于基因相等）。
func (a *Archive) Contains(sol Solution) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, existing := range a.items {
		if genesEqual(existing.Genes, sol.Genes) {
			return true
		}
	}
	return false
}

// Best 返回第一个目标值最优的解。
func (a *Archive) Best(objectiveIdx int) (Solution, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.items) == 0 {
		return Solution{}, false
	}
	best := a.items[0]
	for _, s := range a.items[1:] {
		if objectiveIdx < len(s.Objectives) && objectiveIdx < len(best.Objectives) {
			if s.Objectives[objectiveIdx] > best.Objectives[objectiveIdx] {
				best = s
			}
		}
	}
	return best, true
}

// HyperVolume 计算归档中解集的超体积指标（2D 情况）。
func (a *Archive) HyperVolume(reference []float64) float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.items) == 0 || len(reference) < 2 {
		return 0
	}
	// 按第一个目标排序
	sorted := make([]Solution, len(a.items))
	copy(sorted, a.items)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Objectives[0] > sorted[j].Objectives[0]
	})
	hv := 0.0
	prevY := reference[1]
	for _, s := range sorted {
		if len(s.Objectives) < 2 {
			continue
		}
		x := reference[0] - s.Objectives[0]
		y := prevY - s.Objectives[1]
		if x > 0 && y > 0 {
			hv += x * y
		}
		if s.Objectives[1] < prevY {
			prevY = s.Objectives[1]
		}
	}
	return hv
}

// Spread 计算解集的分散性指标。
func (a *Archive) Spread() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	n := len(a.items)
	if n < 2 {
		return 0
	}
	distances := make([]float64, 0, n-1)
	sorted := make([]Solution, n)
	copy(sorted, a.items)
	sort.Slice(sorted, func(i, j int) bool {
		if len(sorted[i].Objectives) == 0 {
			return false
		}
		if len(sorted[j].Objectives) == 0 {
			return true
		}
		return sorted[i].Objectives[0] < sorted[j].Objectives[0]
	})
	for i := 1; i < n; i++ {
		d := objDist(sorted[i-1].Objectives, sorted[i].Objectives)
		distances = append(distances, d)
	}
	avg := 0.0
	for _, d := range distances {
		avg += d
	}
	avg /= float64(len(distances))
	sum := 0.0
	for _, d := range distances {
		diff := d - avg
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(distances)))
}

func (a *Archive) trim() {
	// 简单策略：移除离最近邻最近的解
	minDist := math.Inf(1)
	minIdx := 0
	for i := range a.items {
		nearest := math.Inf(1)
		for j := range a.items {
			if i == j {
				continue
			}
			d := objDist(a.items[i].Objectives, a.items[j].Objectives)
			if d < nearest {
				nearest = d
			}
		}
		if nearest < minDist {
			minDist = nearest
			minIdx = i
		}
	}
	a.items[minIdx] = a.items[len(a.items)-1]
	a.items = a.items[:len(a.items)-1]
}

// dominates 判断 a 是否支配 b（所有目标 ≥ 且至少一个 >）。
func dominates(a, b []float64) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	allGE := true
	anyGT := false
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			allGE = false
			break
		}
		if a[i] > b[i] {
			anyGT = true
		}
	}
	return allGE && anyGT
}

func genesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func objDist(a, b []float64) float64 {
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

// NonDominatedSort 执行非支配排序，返回每个解的前沿等级。
func NonDominatedSort(solutions []Solution) []int {
	n := len(solutions)
	ranks := make([]int, n)
	dominated := make([][]int, n)
	domCount := make([]int, n)
	for i := range dominated {
		dominated[i] = make([]int, 0)
	}
	var front []int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dominates(solutions[i].Objectives, solutions[j].Objectives) {
				dominated[i] = append(dominated[i], j)
				domCount[j]++
			} else if dominates(solutions[j].Objectives, solutions[i].Objectives) {
				dominated[j] = append(dominated[j], i)
				domCount[i]++
			}
		}
		if domCount[i] == 0 {
			ranks[i] = 0
			front = append(front, i)
		}
	}
	rank := 0
	for len(front) > 0 {
		var nextFront []int
		for _, i := range front {
			ranks[i] = rank
			for _, j := range dominated[i] {
				domCount[j]--
				if domCount[j] == 0 {
					nextFront = append(nextFront, j)
				}
			}
		}
		rank++
		front = nextFront
	}
	return ranks
}
