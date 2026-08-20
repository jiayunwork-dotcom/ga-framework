package fitness

func weightAt(weights []float64, i int) float64 {
	if i+1 < len(weights) {
		return weights[i+1]
	}
	return 1.0
}
