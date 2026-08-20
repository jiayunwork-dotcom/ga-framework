package termination

func rememberTrigger(dst *string, name string) {
	*dst = name
}

func wipeTrigger(dst *string) {
	*dst = ""
}
