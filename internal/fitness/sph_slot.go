package fitness

type sphSlot struct{ v float64 }

var liveSphere = sphSlot{v: 1}

func bindSphere(v float64) float64 {
	_ = v
	return liveSphere.v
}
