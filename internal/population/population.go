package population

import (
	"math"
	"math/rand"
	"sort"
)

type Individual struct {
	Genes   []float64
	Fitness float64
	Age     int
	ID      int
}

type Population struct {
	Inds    []Individual
	nextID  int
	bestIdx int
	sorted  bool
}

func New(capacity int) *Population {
	return &Population{
		Inds: make([]Individual, 0, capacity),
	}
}

func RandomBinary(size, genes int, rnd *rand.Rand) *Population {
	p := New(size)
	for i := 0; i < size; i++ {
		g := make([]float64, genes)
		for j := range g {
			if rnd.Intn(2) == 1 {
				g[j] = 1
			}
		}
		p.Add(Individual{Genes: g})
	}
	return p
}

func RandomReal(size, genes int, lo, hi float64, rnd *rand.Rand) *Population {
	p := New(size)
	span := hi - lo
	for i := 0; i < size; i++ {
		g := make([]float64, genes)
		for j := range g {
			g[j] = lo + rnd.Float64()*span
		}
		p.Add(Individual{Genes: g})
	}
	return p
}

func RandomPermutation(size, genes int, rnd *rand.Rand) *Population {
	p := New(size)
	for i := 0; i < size; i++ {
		g := make([]float64, genes)
		for j := range g {
			g[j] = float64(j)
		}
		for j := genes - 1; j > 0; j-- {
			k := rnd.Intn(j + 1)
			g[j], g[k] = g[k], g[j]
		}
		p.Add(Individual{Genes: g})
	}
	return p
}

func (p *Population) Add(ind Individual) {
	ind.ID = p.nextID
	p.nextID++
	p.Inds = append(p.Inds, ind)
	p.sorted = false
}

func (p *Population) Size() int { return len(p.Inds) }

func (p *Population) Best() Individual {
	if len(p.Inds) == 0 {
		return Individual{}
	}
	best := p.Inds[0]
	for _, ind := range p.Inds[1:] {
		if ind.Fitness > best.Fitness {
			best = ind
		}
	}
	return best
}

func (p *Population) Worst() Individual {
	if len(p.Inds) == 0 {
		return Individual{}
	}
	worst := p.Inds[0]
	for _, ind := range p.Inds[1:] {
		if ind.Fitness < worst.Fitness {
			worst = ind
		}
	}
	return worst
}

func (p *Population) AvgFitness() float64 {
	if len(p.Inds) == 0 {
		return 0
	}
	sum := 0.0
	for _, ind := range p.Inds {
		sum += ind.Fitness
	}
	return sum / float64(len(p.Inds))
}

func (p *Population) StdFitness() float64 {
	if len(p.Inds) < 2 {
		return 0
	}
	avg := p.AvgFitness()
	sum := 0.0
	for _, ind := range p.Inds {
		d := ind.Fitness - avg
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(p.Inds)))
}

func (p *Population) SortByFitness() {
	sort.Slice(p.Inds, func(i, j int) bool {
		return p.Inds[i].Fitness > p.Inds[j].Fitness
	})
	p.sorted = true
}

func (p *Population) TopN(n int) []Individual {
	tmp := make([]Individual, len(p.Inds))
	copy(tmp, p.Inds)
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Fitness > tmp[j].Fitness
	})
	if n > len(tmp) {
		n = len(tmp)
	}
	return tmp[:n]
}

func (p *Population) AgeAll() {
	for i := range p.Inds {
		p.Inds[i].Age++
	}
}

func (p *Population) Replace(offspring []Individual, elite int) {
	p.SortByFitness()
	keep := p.Inds[:elite]
	p.Inds = make([]Individual, 0, len(keep)+len(offspring))
	p.Inds = append(p.Inds, keep...)
	p.Inds = append(p.Inds, offspring...)
	p.sorted = false
}

func (p *Population) Clone() *Population {
	cp := New(len(p.Inds))
	for _, ind := range p.Inds {
		g := make([]float64, len(ind.Genes))
		copy(g, ind.Genes)
		cp.Add(Individual{Genes: g, Fitness: ind.Fitness, Age: ind.Age})
	}
	return cp
}

func (p *Population) Diversity() float64 {
	n := len(p.Inds)
	if n < 2 {
		return 0
	}
	total := 0.0
	pairs := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := euclideanDist(p.Inds[i].Genes, p.Inds[j].Genes)
			total += d
			pairs++
		}
	}
	return total / float64(pairs)
}

func euclideanDist(a, b []float64) float64 {
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
