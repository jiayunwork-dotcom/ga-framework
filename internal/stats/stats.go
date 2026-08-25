package stats

import (
	"math"
	"sort"
	"sync"
)

type Record struct {
	Generation int
	BestFit    float64
	AvgFit     float64
	StdFit     float64
	WorstFit   float64
	Diversity  float64
}

type Collector struct {
	mu      sync.Mutex
	records []Record
}

func NewCollector() *Collector {
	return &Collector{
		records: make([]Record, 0, 128),
	}
}

func (c *Collector) Add(r Record) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records = append(c.records, r)
}

func (c *Collector) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.records)
}

func (c *Collector) Records() []Record {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Record, len(c.records))
	copy(out, c.records)
	return out
}

func (c *Collector) Last() (Record, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) == 0 {
		return Record{}, false
	}
	return c.records[len(c.records)-1], true
}

func (c *Collector) BestEver() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) == 0 {
		return math.Inf(-1)
	}
	best := c.records[0].BestFit
	for _, r := range c.records[1:] {
		if r.BestFit > best {
			best = r.BestFit
		}
	}
	return best
}

func (c *Collector) AvgBest() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range c.records {
		sum += r.BestFit
	}
	return sum / float64(len(c.records))
}

func (c *Collector) Stagnation() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := len(c.records)
	if n < 2 {
		return 0
	}
	best := c.records[n-1].BestFit
	count := 0
	for i := n - 2; i >= 0; i-- {
		if c.records[i].BestFit >= best {
			count++
		} else {
			break
		}
	}
	return count
}

func (c *Collector) ImprovementRate(window int) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := len(c.records)
	if n < 2 || window < 2 {
		return 0
	}
	if window > n {
		window = n
	}
	start := c.records[n-window].BestFit
	end := c.records[n-1].BestFit
	if start == 0 {
		return 0
	}
	return (end - start) / math.Abs(start)
}

type Summary struct {
	TotalGenerations int
	FinalBest        float64
	MaxStagnation    int
	AvgDiversity     float64
	ConvergenceGen   int
}

func (c *Collector) Summarize() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := Summary{TotalGenerations: len(c.records)}
	if len(c.records) == 0 {
		return s
	}
	bestFit := c.records[0].BestFit
	convGen := 0
	maxStag := 0
	curStag := 0
	divSum := 0.0
	for i, r := range c.records {
		divSum += r.Diversity
		if r.BestFit > bestFit {
			bestFit = r.BestFit
			convGen = i
			curStag = 0
		} else {
			curStag++
			if curStag > maxStag {
				maxStag = curStag
			}
		}
	}
	s.FinalBest = bestFit
	s.MaxStagnation = maxStag
	s.ConvergenceGen = convGen
	s.AvgDiversity = divSum / float64(len(c.records))
	return s
}

func (c *Collector) Percentile(p float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) == 0 {
		return 0
	}
	vals := make([]float64, len(c.records))
	for i, r := range c.records {
		vals[i] = r.BestFit
	}
	sort.Float64s(vals)
	idx := int(math.Floor(p / 100.0 * float64(len(vals)-1)))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(vals) {
		idx = len(vals) - 1
	}
	return vals[idx]
}

func (c *Collector) MovingAvg(w int) []float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := len(c.records)
	if n == 0 || w < 1 {
		return nil
	}
	result := make([]float64, n)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += c.records[i].BestFit
		if i >= w {
			sum -= c.records[i-w].BestFit
			result[i] = sum / float64(w)
		} else {
			result[i] = sum / float64(i+1)
		}
	}
	return result
}
