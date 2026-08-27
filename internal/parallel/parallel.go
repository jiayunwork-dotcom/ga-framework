package parallel

import (
	"runtime"
	"sync"
)

type EvalFunc func(genes []float64) float64

type EvalResult struct {
	Index   int
	Fitness float64
}

func EvalAll(population [][]float64, f EvalFunc, workers int) []float64 {
	n := len(population)
	if n == 0 {
		return nil
	}
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	results := make([]float64, n)
	var wg sync.WaitGroup
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range ch {
				results[idx] = f(population[idx])
			}
		}()
	}
	wg.Wait()
	return results
}

func BatchEval(population [][]float64, f EvalFunc, batchSize int) <-chan EvalResult {
	out := make(chan EvalResult, len(population))
	go func() {
		defer close(out)
		for i := 0; i < len(population); i += batchSize {
			end := i + batchSize
			if end > len(population) {
				end = len(population)
			}
			var wg sync.WaitGroup
			for j := i; j < end; j++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					fit := f(population[idx])
					out <- EvalResult{Index: idx, Fitness: fit}
				}(j)
			}
			wg.Wait()
		}
	}()
	return out
}

type Pool struct {
	workers int
	tasks   chan task
	wg      sync.WaitGroup
	running bool
	mu      sync.Mutex
}

type task struct {
	fn   func()
	done chan struct{}
}

func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	p := &Pool{
		workers: workers,
		tasks:   make(chan task, workers*2),
		running: true,
	}
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	for t := range p.tasks {
		t.fn()
		if t.done != nil {
			close(t.done)
		}
	}
}

func (p *Pool) Submit(fn func()) {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
	done := make(chan struct{})
	p.tasks <- task{fn: fn, done: done}
	<-done
}

func (p *Pool) SubmitAsync(fn func()) <-chan struct{} {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		ch := make(chan struct{})
		close(ch)
		return ch
	}
	p.mu.Unlock()
	done := make(chan struct{})
	p.tasks <- task{fn: fn, done: done}
	return done
}

func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return
	}
	p.running = false
	close(p.tasks)
}

func (p *Pool) MapEval(population [][]float64, f EvalFunc) []float64 {
	n := len(population)
	results := make([]float64, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		idx := i
		p.mu.Lock()
		if !p.running {
			p.mu.Unlock()
			wg.Done()
			continue
		}
		p.mu.Unlock()
		p.tasks <- task{
			fn: func() {
				results[idx] = f(population[idx])
				wg.Done()
			},
		}
	}
	wg.Wait()
	return results
}
