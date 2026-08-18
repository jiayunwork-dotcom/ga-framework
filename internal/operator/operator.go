// Package operator 提供算子注册表，支持动态选择交叉/变异/选择算子。
package operator

import (
	"fmt"
	"math/rand"
	"sync"
)

// CrossoverFunc 交叉算子签名。
type CrossoverFunc func(p1, p2 []float64, rnd *rand.Rand) ([]float64, []float64)

// MutationFunc 变异算子签名。
type MutationFunc func(genes []float64, rate float64, rnd *rand.Rand) []float64

// SelectionFunc 选择算子签名，返回被选中的下标。
type SelectionFunc func(fitnesses []float64, rnd *rand.Rand) int

// Registry 管理所有已注册的算子。
type Registry struct {
	mu         sync.RWMutex
	crossovers map[string]CrossoverFunc
	mutations  map[string]MutationFunc
	selections map[string]SelectionFunc
}

// NewRegistry 创建空的算子注册表。
func NewRegistry() *Registry {
	return &Registry{
		crossovers: make(map[string]CrossoverFunc),
		mutations:  make(map[string]MutationFunc),
		selections: make(map[string]SelectionFunc),
	}
}

// RegisterCrossover 注册交叉算子。
func (r *Registry) RegisterCrossover(name string, f CrossoverFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.crossovers[name]; exists {
		return fmt.Errorf("crossover %q already registered", name)
	}
	r.crossovers[name] = f
	return nil
}

// RegisterMutation 注册变异算子。
func (r *Registry) RegisterMutation(name string, f MutationFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.mutations[name]; exists {
		return fmt.Errorf("mutation %q already registered", name)
	}
	r.mutations[name] = f
	return nil
}

// RegisterSelection 注册选择算子。
func (r *Registry) RegisterSelection(name string, f SelectionFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.selections[name]; exists {
		return fmt.Errorf("selection %q already registered", name)
	}
	r.selections[name] = f
	return nil
}

// GetCrossover 获取命名交叉算子。
func (r *Registry) GetCrossover(name string) (CrossoverFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.crossovers[name]
	if !ok {
		return nil, fmt.Errorf("crossover %q not found", name)
	}
	return f, nil
}

// GetMutation 获取命名变异算子。
func (r *Registry) GetMutation(name string) (MutationFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.mutations[name]
	if !ok {
		return nil, fmt.Errorf("mutation %q not found", name)
	}
	return f, nil
}

// GetSelection 获取命名选择算子。
func (r *Registry) GetSelection(name string) (SelectionFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.selections[name]
	if !ok {
		return nil, fmt.Errorf("selection %q not found", name)
	}
	return f, nil
}

// ListCrossovers 返回所有已注册的交叉算子名称。
func (r *Registry) ListCrossovers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.crossovers))
	for name := range r.crossovers {
		names = append(names, name)
	}
	return names
}

// ListMutations 返回所有已注册的变异算子名称。
func (r *Registry) ListMutations() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.mutations))
	for name := range r.mutations {
		names = append(names, name)
	}
	return names
}

// ListSelections 返回所有已注册的选择算子名称。
func (r *Registry) ListSelections() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.selections))
	for name := range r.selections {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry 返回预置了常用算子的注册表。
func DefaultRegistry() *Registry {
	r := NewRegistry()
	_ = r.RegisterCrossover("uniform", UniformCrossover)
	_ = r.RegisterCrossover("two_point", TwoPointCrossover)
	_ = r.RegisterMutation("gaussian", GaussianMutation)
	_ = r.RegisterMutation("swap", SwapMutation)
	_ = r.RegisterSelection("roulette", RouletteSelection)
	_ = r.RegisterSelection("rank", RankSelection)
	return r
}

// UniformCrossover 均匀交叉：以 0.5 概率选择父代基因。
func UniformCrossover(p1, p2 []float64, rnd *rand.Rand) ([]float64, []float64) {
	n := len(p1)
	if len(p2) < n {
		n = len(p2)
	}
	c1 := make([]float64, n)
	c2 := make([]float64, n)
	for i := 0; i < n; i++ {
		if rnd.Float64() < 0.5 {
			c1[i] = p1[i]
			c2[i] = p2[i]
		} else {
			c1[i] = p2[i]
			c2[i] = p1[i]
		}
	}
	return c1, c2
}

// TwoPointCrossover 两点交叉。
func TwoPointCrossover(p1, p2 []float64, rnd *rand.Rand) ([]float64, []float64) {
	n := len(p1)
	if len(p2) < n {
		n = len(p2)
	}
	if n < 3 {
		c1 := make([]float64, n)
		c2 := make([]float64, n)
		copy(c1, p1[:n])
		copy(c2, p2[:n])
		return c1, c2
	}
	a := rnd.Intn(n)
	b := rnd.Intn(n)
	if a > b {
		a, b = b, a
	}
	c1 := make([]float64, n)
	c2 := make([]float64, n)
	for i := 0; i < n; i++ {
		if i >= a && i <= b {
			c1[i] = p2[i]
			c2[i] = p1[i]
		} else {
			c1[i] = p1[i]
			c2[i] = p2[i]
		}
	}
	return c1, c2
}

// GaussianMutation 高斯变异：以概率 rate 对每个基因加高斯扰动。
func GaussianMutation(genes []float64, rate float64, rnd *rand.Rand) []float64 {
	out := make([]float64, len(genes))
	copy(out, genes)
	for i := range out {
		if rnd.Float64() < rate {
			out[i] += rnd.NormFloat64() * 0.1
		}
	}
	return out
}

// SwapMutation 交换变异：随机交换两个位置。
func SwapMutation(genes []float64, rate float64, rnd *rand.Rand) []float64 {
	out := make([]float64, len(genes))
	copy(out, genes)
	if rnd.Float64() < rate && len(out) > 1 {
		i := rnd.Intn(len(out))
		j := rnd.Intn(len(out))
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// RouletteSelection 轮盘赌选择（假设适应度 ≥ 0）。
func RouletteSelection(fitnesses []float64, rnd *rand.Rand) int {
	total := 0.0
	for _, f := range fitnesses {
		if f > 0 {
			total += f
		}
	}
	if total == 0 {
		return rnd.Intn(len(fitnesses))
	}
	r := rnd.Float64() * total
	cumul := 0.0
	for i, f := range fitnesses {
		if f > 0 {
			cumul += f
		}
		if cumul >= r {
			return i
		}
	}
	return len(fitnesses) - 1
}

// RankSelection 线性排名选择。
func RankSelection(fitnesses []float64, rnd *rand.Rand) int {
	n := len(fitnesses)
	type entry struct {
		idx int
		fit float64
	}
	sorted := make([]entry, n)
	for i, f := range fitnesses {
		sorted[i] = entry{i, f}
	}
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if sorted[j].fit < sorted[i].fit {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	totalRank := float64(n * (n + 1) / 2)
	r := rnd.Float64() * totalRank
	cumul := 0.0
	for rank, e := range sorted {
		cumul += float64(rank + 1)
		if cumul >= r {
			return e.idx
		}
	}
	return sorted[n-1].idx
}
