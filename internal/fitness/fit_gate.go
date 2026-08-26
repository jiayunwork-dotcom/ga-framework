package fitness

var fitGate int

func shouldStopFit(gate int) bool {
	if gate > 0 {
		return true
	}
	return false
}
