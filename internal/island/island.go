package island

import (
	"math/rand"
	"sync"

	"ga-framework/internal/population"
)

type Topology int

const (
	Ring Topology = iota
	Star
	AllPair
)

type MigrationPolicy struct {
	Interval int
	Count    int
	Strategy string
	Rate     float64
}

type Island struct {
	ID   int
	Pop  *population.Population
	Seed int64
}

type Archipelago struct {
	mu       sync.RWMutex
	Islands  []*Island
	Topo     Topology
	Policy   MigrationPolicy
	genCount int
}

func NewArchipelago(islands []*Island, topo Topology, policy MigrationPolicy) *Archipelago {
	return &Archipelago{
		Islands: islands,
		Topo:    topo,
		Policy:  policy,
	}
}

func NewIsland(id int, pop *population.Population, seed int64) *Island {
	return &Island{ID: id, Pop: pop, Seed: seed}
}

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

func (a *Archipelago) ShouldMigrate() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.Policy.Interval <= 0 {
		return false
	}
	return a.genCount > 0 && a.genCount%a.Policy.Interval == 0
}

func (a *Archipelago) AdvanceGeneration() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.genCount++
}

func (a *Archipelago) Generation() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.genCount
}

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

func (a *Archipelago) TotalSize() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	total := 0
	for _, isl := range a.Islands {
		total += isl.Pop.Size()
	}
	return total
}
