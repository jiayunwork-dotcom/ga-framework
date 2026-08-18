package island

import (
	"math/rand"
	"testing"

	"ga-framework/internal/population"
)

func makePop(size int, rnd *rand.Rand) *population.Population {
	p := population.RandomReal(size, 5, 0, 10, rnd)
	for i := range p.Inds {
		p.Inds[i].Fitness = rnd.Float64() * 100
	}
	return p
}

func TestNeighbors_Ring(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	islands := make([]*Island, 4)
	for i := range islands {
		islands[i] = NewIsland(i, makePop(5, rnd), int64(i))
	}
	a := NewArchipelago(islands, Ring, MigrationPolicy{Interval: 5, Count: 1, Strategy: "best", Rate: 1.0})
	n := a.Neighbors(0)
	if len(n) != 2 || n[0] != 3 || n[1] != 1 {
		t.Fatalf("ring neighbors of 0: expected [3,1], got %v", n)
	}
}

func TestNeighbors_Star(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	islands := make([]*Island, 4)
	for i := range islands {
		islands[i] = NewIsland(i, makePop(5, rnd), int64(i))
	}
	a := NewArchipelago(islands, Star, MigrationPolicy{})
	n := a.Neighbors(0)
	if len(n) != 3 {
		t.Fatalf("star center should have 3 neighbors, got %d", len(n))
	}
	n2 := a.Neighbors(2)
	if len(n2) != 1 || n2[0] != 0 {
		t.Fatalf("star leaf should connect to center: %v", n2)
	}
}

func TestMigrate_IncreasesPopulation(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))
	islands := make([]*Island, 3)
	for i := range islands {
		islands[i] = NewIsland(i, makePop(10, rnd), int64(i))
	}
	a := NewArchipelago(islands, AllPair, MigrationPolicy{
		Interval: 1, Count: 2, Strategy: "best", Rate: 1.0,
	})
	totalBefore := a.TotalSize()
	a.Migrate(rnd)
	totalAfter := a.TotalSize()
	if totalAfter <= totalBefore {
		t.Fatalf("migration should add individuals: before=%d, after=%d", totalBefore, totalAfter)
	}
}

func TestGlobalBest(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	islands := make([]*Island, 3)
	for i := range islands {
		islands[i] = NewIsland(i, makePop(5, rnd), int64(i))
	}
	islands[1].Pop.Inds[0].Fitness = 9999
	a := NewArchipelago(islands, Ring, MigrationPolicy{})
	best := a.GlobalBest()
	if best.Fitness != 9999 {
		t.Fatalf("expected global best 9999, got %v", best.Fitness)
	}
}

func TestShouldMigrate(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	islands := []*Island{NewIsland(0, makePop(5, rnd), 0)}
	a := NewArchipelago(islands, Ring, MigrationPolicy{Interval: 3})
	if a.ShouldMigrate() {
		t.Fatalf("gen 0 should not trigger migration")
	}
	a.AdvanceGeneration()
	a.AdvanceGeneration()
	a.AdvanceGeneration()
	if !a.ShouldMigrate() {
		t.Fatalf("gen 3 should trigger migration with interval=3")
	}
}
