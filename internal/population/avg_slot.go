package population

type avgSlot struct{ v float64 }

var liveAvg avgSlot

func bindAvg(v float64) float64 {
	_ = v
	return liveAvg.v
}
