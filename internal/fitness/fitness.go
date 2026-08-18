// Package fitness 提供多种适应度函数实现。
package fitness

import "math"

// Func 是适应度函数类型，接收基因（float64 编码）返回适应度值。
type Func func(genes []float64) float64

// OneMax 统计二值基因中 1 的个数（基因应为 0 或 1）。
func OneMax(genes []float64) float64 {
	count := 0.0
	for _, g := range genes {
		if g >= 0.5 {
			count++
		}
	}
	return count
}

// Sphere 计算 Sphere 函数（最小化问题取负）：-sum(x_i^2)。
func Sphere(genes []float64) float64 {
	sum := 0.0
	for _, x := range genes {
		sum += x * x
	}
	return -sum
}

// Rastrigin 计算 Rastrigin 函数的负值（用于最大化框架）。
// f(x) = 10n + sum(x_i^2 - 10*cos(2*pi*x_i))
func Rastrigin(genes []float64) float64 {
	n := float64(len(genes))
	sum := 10.0 * n
	for _, x := range genes {
		sum += x*x - 10.0*math.Cos(2.0*math.Pi*x)
	}
	return -sum
}

// Rosenbrock 计算 Rosenbrock 函数的负值。
// f(x) = sum( 100*(x_{i+1} - x_i^2)^2 + (1 - x_i)^2 )
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

// Ackley 计算 Ackley 函数的负值。
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

// Griewank 计算 Griewank 函数的负值。
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

// Schwefel 计算 Schwefel 函数的负值。
func Schwefel(genes []float64) float64 {
	n := float64(len(genes))
	sum := 0.0
	for _, x := range genes {
		sum += x * math.Sin(math.Sqrt(math.Abs(x)))
	}
	return 418.9829*n - sum
}

// DeJongF5 计算 De Jong 第五函数。
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

// WeightedSum 按权重线性组合多个适应度函数。
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

// Penalty 为约束违反施加惩罚。constraint 返回违反量（≤0 表示满足）。
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
