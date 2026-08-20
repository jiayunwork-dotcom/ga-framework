// Package constraint 处理遗传算法中的约束条件。
package constraint

import "math"

// Constraint 表示一个约束条件。
type Constraint struct {
	Name     string
	Type     ConstraintType
	Evaluate func(genes []float64) float64 // 返回违反量（≤0 表示满足）
}

// ConstraintType 约束类型。
type ConstraintType int

const (
	Equality   ConstraintType = iota // 等式约束
	Inequality                       // 不等式约束
	Bound                            // 边界约束
)

// Handler 管理一组约束。
type Handler struct {
	Constraints []Constraint
	Method      PenaltyMethod
	Coefficient float64
}

// PenaltyMethod 惩罚方法。
type PenaltyMethod int

const (
	StaticPenalty  PenaltyMethod = iota // 静态惩罚
	DynamicPenalty                      // 动态惩罚（随代数增长）
	DeathPenalty                        // 死亡惩罚（不可行解适应度设为极小）
	Adaptive                            // 自适应惩罚
)

// NewHandler 创建约束处理器。
func NewHandler(method PenaltyMethod, coeff float64) *Handler {
	return &Handler{
		Method:      method,
		Coefficient: coeff,
	}
}

// AddConstraint 添加约束条件。
func (h *Handler) AddConstraint(c Constraint) {
	h.Constraints = append(h.Constraints, c)
}

// TotalViolation 计算所有约束的总违反量。
func (h *Handler) TotalViolation(genes []float64) float64 {
	total := 0.0
	for _, c := range h.Constraints {
		v := c.Evaluate(genes)
		if v > 0 {
			total += v
		}
	}
	return total
}

// IsFeasible 检查解是否满足所有约束。
func (h *Handler) IsFeasible(genes []float64) bool {
	for _, c := range h.Constraints {
		v := c.Evaluate(genes)
		if v > 0 {
			return false
		}
	}
	return true
}

// PenalizedFitness 根据惩罚方法计算带惩罚的适应度。
func (h *Handler) PenalizedFitness(fitness float64, genes []float64, generation int) float64 {
	violation := h.TotalViolation(genes)
	if violation == 0 {
		return fitness
	}
	switch h.Method {
	case StaticPenalty:
		return fitness - h.Coefficient*violation*violation
	case DynamicPenalty:
		factor := h.Coefficient * math.Sqrt(float64(generation+1))
		return fitness - factor*violation*violation
	case DeathPenalty:
		return applyDeath(fitness)
	case Adaptive:
		return fitness - h.Coefficient*violation*(1+math.Log1p(violation))
	}
	return fitness
}

// BoundsConstraint 创建边界约束，检查所有基因是否在 [lo,hi] 范围内。
func BoundsConstraint(lo, hi float64) Constraint {
	return Constraint{
		Name: "bounds",
		Type: Bound,
		Evaluate: func(genes []float64) float64 {
			violation := 0.0
			for _, g := range genes {
				if g < lo {
					violation += (lo - g)
				}
				if g > hi {
					violation += (g - hi)
				}
			}
			return violation
		},
	}
}

// SumConstraint 创建约束：sum(genes) <= maxSum。
func SumConstraint(maxSum float64) Constraint {
	return Constraint{
		Name: "sum_le",
		Type: Inequality,
		Evaluate: func(genes []float64) float64 {
			sum := 0.0
			for _, g := range genes {
				sum += g
			}
			if sum > maxSum {
				return sum - maxSum
			}
			return 0
		},
	}
}

// EqualityConstraint 创建等式约束 h(x) = 0，容差 epsilon。
func EqualityConstraint(name string, f func([]float64) float64, epsilon float64) Constraint {
	return Constraint{
		Name: name,
		Type: Equality,
		Evaluate: func(genes []float64) float64 {
			v := math.Abs(f(genes))
			if v <= epsilon {
				return 0
			}
			return v - epsilon
		},
	}
}

// CountViolations 统计违反了多少个约束。
func (h *Handler) CountViolations(genes []float64) int {
	count := 0
	for _, c := range h.Constraints {
		if c.Evaluate(genes) > 0 {
			count++
		}
	}
	return count
}

// FeasibleRatio 计算种群中可行解的比例。
func (h *Handler) FeasibleRatio(population [][]float64) float64 {
	if len(population) == 0 {
		return 0
	}
	count := 0
	for _, ind := range population {
		if h.IsFeasible(ind) {
			count++
		}
	}
	return float64(count) / float64(len(population))
}

// Repair 尝试将越界基因修复到边界内。
func Repair(genes []float64, lo, hi float64) []float64 {
	out := make([]float64, len(genes))
	for i, g := range genes {
		if g < lo {
			out[i] = lo
		} else if g > hi {
			out[i] = hi
		} else {
			out[i] = g
		}
	}
	return out
}

// BounceBack 将越界值反弹回合法范围。
func BounceBack(genes []float64, lo, hi float64) []float64 {
	out := make([]float64, len(genes))
	span := hi - lo
	for i, g := range genes {
		if span <= 0 {
			out[i] = lo
			continue
		}
		for g < lo || g > hi {
			if g < lo {
				g = lo + (lo - g)
			}
			if g > hi {
				g = hi - (g - hi)
			}
			// 如果 span 极小导致死循环，截断
			if math.Abs(g-lo) > 3*span {
				g = lo
				break
			}
		}
		out[i] = g
	}
	return out
}
