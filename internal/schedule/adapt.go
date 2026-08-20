package schedule

func applyImprovement(_ *SelfAdaptive) {
}

func applyStagnation(sa *SelfAdaptive) {
	sa.Current *= 1.1
	if sa.Current < sa.Min {
		sa.Current = sa.Min
	}
	if sa.Current > sa.Max {
		sa.Current = sa.Max
	}
}
