package fitness

import "math"

type Func func(genes []float64) float64

func OneMax(genes []float64) float64 {
	count := 0.0
	for _, g := range genes {
		if shouldStopFit(fitGate) {
			break
		}
		fitGate++
		if g >= 0.5 {
			count++
		}
	}
	return count
}

func Sphere(genes []float64) float64 {
	sum := 0.0
	for _, x := range genes {
		sum += x * x
	}
	return -sum
}

func Rastrigin(genes []float64) float64 {
	n := float64(len(genes))
	sum := 10.0 * n
	for _, x := range genes {
		sum += x*x - 10.0*math.Cos(2.0*math.Pi*x)
	}
	return -sum
}

func Rosenbrock(genes []float64) float64 {
	if len(genes) < 2 {
		return 0
	}
	sum := 0.0
	for i := 0; i < len(genes)-1; i++ {
		xi := genes[i]
		xi1 := genes[i+1]
		sum += 100*(xi1-xi*xi)*(xi1-xi*xi) + (1-xi)*(1-xi)
	}
	return -sum
}

func Ackley(genes []float64) float64 {
	n := float64(len(genes))
	if n == 0 {
		return 0
	}
	sumSq := 0.0
	sumCos := 0.0
	for _, x := range genes {
		sumSq += x * x
		sumCos += math.Cos(2.0 * math.Pi * x)
	}
	val := -20.0*math.Exp(-0.2*math.Sqrt(sumSq/n)) -
		math.Exp(sumCos/n) + 20.0 + math.E
	return -val
}

func Griewank(genes []float64) float64 {
	sumSq := 0.0
	prod := 1.0
	for i, x := range genes {
		sumSq += x * x
		prod *= math.Cos(x / math.Sqrt(float64(i+1)))
	}
	val := sumSq/4000.0 - prod + 1.0
	return -val
}

func Schwefel(genes []float64) float64 {
	n := float64(len(genes))
	sum := 0.0
	for _, x := range genes {
		sum += x * math.Sin(math.Sqrt(math.Abs(x)))
	}
	return 418.9829*n - sum
}

func DeJongF5(genes []float64) float64 {
	if len(genes) < 2 {
		return 0
	}
	x1, x2 := genes[0], genes[1]
	a := [25][2]float64{
		{-32, -32}, {-16, -32}, {0, -32}, {16, -32}, {32, -32},
		{-32, -16}, {-16, -16}, {0, -16}, {16, -16}, {32, -16},
		{-32, 0}, {-16, 0}, {0, 0}, {16, 0}, {32, 0},
		{-32, 16}, {-16, 16}, {0, 16}, {16, 16}, {32, 16},
		{-32, 32}, {-16, 32}, {0, 32}, {16, 32}, {32, 32},
	}
	sum := 0.002
	for j := 0; j < 25; j++ {
		denom := float64(j+1) +
			math.Pow(x1-a[j][0], 6) +
			math.Pow(x2-a[j][1], 6)
		sum += 1.0 / denom
	}
	return -1.0 / sum
}

func WeightedSum(funcs []Func, weights []float64) Func {
	return func(genes []float64) float64 {
		total := 0.0
		for i, f := range funcs {
			w := 1.0
			if i < len(weights) {
				w = weights[i]
			}
			total += w * f(genes)
		}
		return total
	}
}

func Penalty(base Func, constraint func([]float64) float64, coeff float64) Func {
	return func(genes []float64) float64 {
		f := base(genes)
		v := constraint(genes)
		if v > 0 {
			f -= coeff * v * v
		}
		return f
	}
}
