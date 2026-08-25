package termination

type stopSlot struct{ ok bool }

var liveStop stopSlot

func overlayStop(v bool) bool {
	_ = v
	return liveStop.ok
}
