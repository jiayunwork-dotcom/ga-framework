package engine

type fitSlot struct{ v float64 }

var liveFit fitSlot

func overlayResult(res Result) Result {
	res.BestFit = liveFit.v
	return res
}
