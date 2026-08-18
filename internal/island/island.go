// Package island 实现岛屿模型并行遗传算法。
package island

import (
	"math/rand"
	"sync"

	"ga-framework/internal/population"
)

// Topology 定义岛屿间的迁移拓扑。
type Topology int

const (
	Ring    Topology = iota // 环形拓扑
	Star                   // 星形拓扑
	AllPair                // 全连接
)

// MigrationPolicy 迁移策略。
type MigrationPolicy struct {
	Interval int     // 每隔多少代执行迁移
	Count    int     // 每次迁移个体数
	Strategy string  // "best" 或 "random"
	Rate     float64 // 迁移概率
}

// Island 表示一座岛。
type Island struct {
	ID   int
	Pop  *population.Population
	Seed int64
}

// Archipelago 多岛管理器。
type Archipelago struct {
	mu       sync.RWMutex
	Islands  []*Island
	Topo     Topology
	Policy   MigrationPolicy
	genCount int
}

// NewArchipelago 创建群岛。
func NewArchipelago(islands []*Island, topo Topology, policy MigrationPolicy) *Archipelago {
	return &Archipelago{
		Islands: islands,
		Topo:    topo,
		Policy:  policy,
	}
}

// NewIsland 创建单座岛。
func NewIsland(id int, pop *population.Population, seed int64) *Island {
	return &Island{ID: id, Pop: pop, Seed: seed}
}

// Neighbors 根据拓扑返回邻居岛 ID。
func (a *Archipelago) Neighbors(islandID int) []int {
	n := len(a.Islands)
	switch a.Topo {
	case Ring:
		prev := (islandID - 1 + n) % n
		next := (islandID + 1) % n
		return []int{prev, next}
	case Star:
		if islandID == 0 {
			ids := make([]int, n-1)
			for i := 1; i < n; i++ {
				ids[i-1] = i
			}
			return ids
		}
		return []int{0}
	case AllPair:
		ids := make([]int, 0, n-1)
		for i := 0; i < n; i++ {
			if i != islandID {
				ids = append(ids, i)
			}
		}
		return ids
	}
	return nil
}

// ShouldMigrate 判断当前代是否执行迁移。
func (a *Archipelago) ShouldMigrate() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.Policy.Interval <= 0 {
		return false
	}
	return a.genCount > 0 && a.genCount%a.Policy.Interval == 0
}

// AdvanceGeneration 推进一代。
func (a *Archipelago) AdvanceGeneration() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.genCount++
}

// Generation 返回当前代数。
func (a *Archipelago) Generation() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.genCount
}

// Migrate 执行一次迁移操作。
func (a *Archipelago) Migrate(rnd *rand.Rand) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, isl := range a.Islands {
		neighbors := a.neighborsLocked(isl.ID)
		if len(neighbors) == 0 || isl.Pop.Size() == 0 {
			continue
		}
		for c := 0; c < a.Policy.Count; c++ {
			if rnd.Float64() > a.Policy.Rate {
				continue
			}
			destID := neighbors[rnd.Intn(len(neighbors))]
			dest := a.Islands[destID]
			emigrant := selectEmigrant(isl.Pop, a.Policy.Strategy, rnd)
			dest.Pop.Add(emigrant)
		}
	}
}

func (a *Archipelago) neighborsLocked(islandID int) []int {
	n := len(a.Islands)
	switch a.Topo {
	case Ring:
		prev := (islandID - 1 + n) % n
		next := (islandID + 1) % n
		return []int{prev, next}
	case Star:
		if islandID == 0 {
			ids := make([]int, n-1)
			for i := 1; i < n; i++ {
				ids[i-1] = i
			}
			return ids
		}
		return []int{0}
	case AllPair:
		ids := make([]int, 0, n-1)
		for i := 0; i < n; i++ {
			if i != islandID {
				ids = append(ids, i)
			}
		}
		return ids
	}
	return nil
}

func selectEmigrant(pop *population.Population, strategy string, rnd *rand.Rand) population.Individual {
	switch strategy {
	case "best":
		return pop.Best()
	default:
		idx := rnd.Intn(pop.Size())
		ind := pop.Inds[idx]
		g := make([]float64, len(ind.Genes))
		copy(g, ind.Genes)
		return population.Individual{Genes: g, Fitness: ind.Fitness}
	}
}

// GlobalBest 返回所有岛中最优个体。
func (a *Archipelago) GlobalBest() population.Individual {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var best population.Individual
	first := true
	for _, isl := range a.Islands {
		if isl.Pop.Size() == 0 {
			continue
		}
		b := isl.Pop.Best()
		if first || b.Fitness > best.Fitness {
			best = b
			first = false
		}
	}
	return best
}

// TotalSize 返回所有岛的总个体数。
func (a *Archipelago) TotalSize() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	total := 0
	for _, isl := range a.Islands {
		total += isl.Pop.Size()
	}
	return total
}
