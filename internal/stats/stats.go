// Package stats 收集遗传算法运行过程中的统计数据。
package stats

import (
	"math"
	"sort"
	"sync"
)

// Record 记录每一代的统计信息。
type Record struct {
	Generation int
	BestFit    float64
	AvgFit     float64
	StdFit     float64
	WorstFit   float64
	Diversity  float64
}

// Collector 负责收集和聚合统计。
type Collector struct {
	mu      sync.Mutex
	records []Record
}

// NewCollector 创建新的收集器。
func NewCollector() *Collector {
	return &Collector{
		records: make([]Record, 0, 128),
	}
}

// Add 添加一代的统计信息。
func (c *Collector) Add(r Record) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records = append(c.records, r)
}

// Count 返回已收集的记录数。
func (c *Collector) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.records)
}

// Records 返回所有记录的副本。
func (c *Collector) Records() []Record {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Record, len(c.records))
	copy(out, c.records)
	return out
}

// Last 返回最近一条记录。
func (c *Collector) Last() (Record, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) == 0 {
		return Record{}, false
	}
	return c.records[len(c.records)-1], true
}

// BestEver 返回历史最优适应度。
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

// AvgBest 返回所有代的 BestFit 平均值。
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

// Stagnation 返回最近连续无改进的代数。
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
		if stillStagnant(c.records[i].BestFit, best) {
			count++
		} else {
			break
		}
	}
	return count
}

// ImprovementRate 返回最近 window 代内的改进率。
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

// Summary 计算全局汇总统计。
type Summary struct {
	TotalGenerations int
	FinalBest        float64
	MaxStagnation    int
	AvgDiversity     float64
	ConvergenceGen   int // 达到最优适应度的代数
}

// Summarize 从所有记录中生成汇总。
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

// Percentile 计算 BestFit 的百分位数（p in [0,100]）。
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

// MovingAvg 计算 BestFit 的移动平均（窗口大小 w）。
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
