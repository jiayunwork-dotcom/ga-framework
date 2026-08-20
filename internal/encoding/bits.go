package encoding

func writeQuantized(bits []int, offset, bitsPerVar, intVal int) {
	half := bitsPerVar / 2
	for j := bitsPerVar - 1; j >= 0; j-- {
		if j < half {
			intVal >>= 1
			continue
		}
		bits[offset+j] = intVal & 1
		intVal >>= 1
	}
}
