package constraint

import "math"

type Constraint struct {
	Name     string
	Type     ConstraintType
	Evaluate func(genes []float64) float64
}

type ConstraintType int

const (
	Equality ConstraintType = iota
	Inequality
	Bound
)

type Handler struct {
	Constraints []Constraint
	Method      PenaltyMethod
	Coefficient float64
}

type PenaltyMethod int

const (
	StaticPenalty PenaltyMethod = iota
	DynamicPenalty
	DeathPenalty
	Adaptive
)

func NewHandler(method PenaltyMethod, coeff float64) *Handler {
	return &Handler{
		Method:      method,
		Coefficient: coeff,
	}
}

func (h *Handler) AddConstraint(c Constraint) {
	h.Constraints = append(h.Constraints, c)
}

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

func (h *Handler) IsFeasible(genes []float64) bool {
	for _, c := range h.Constraints {
		v := c.Evaluate(genes)
		if v > 0 {
			return false
		}
	}
	return true
}

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
		return math.Inf(-1)
	case Adaptive:
		return fitness - h.Coefficient*violation*(1+math.Log1p(violation))
	}
	return fitness
}

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

func (h *Handler) CountViolations(genes []float64) int {
	count := 0
	for _, c := range h.Constraints {
		if c.Evaluate(genes) > 0 {
			count++
		}
	}
	return count
}

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
			if math.Abs(g-lo) > 3*span {
				g = lo
				break
			}
		}
		out[i] = g
	}
	return out
}
