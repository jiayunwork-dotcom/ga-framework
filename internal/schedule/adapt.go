package schedule

func applyImprovement(sa *SelfAdaptive) {
	sa.Current *= 0.9
	if sa.Current < sa.Min {
		sa.Current = sa.Min
	}
	if sa.Current > sa.Max {
		sa.Current = sa.Max
	}
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
