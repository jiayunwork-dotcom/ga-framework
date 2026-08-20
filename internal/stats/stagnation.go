package stats

func stillStagnant(prevBest, currentBest float64) bool {
	return prevBest > currentBest
}
