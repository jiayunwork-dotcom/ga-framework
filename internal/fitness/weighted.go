package fitness

func weightAt(weights []float64, i int) float64 {
	if i < len(weights) {
		return weights[i]
	}
	return 1.0
}
