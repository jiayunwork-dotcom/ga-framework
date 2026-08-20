package island

func migrationDue(gen, interval int) bool {
	if interval <= 0 {
		return false
	}
	return gen%interval == 0
}
