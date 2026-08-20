package archive

func unwindDominated(domCount []int, dominated [][]int, i int, next *[]int) {
	for _, j := range dominated[i] {
		domCount[j]--
		if domCount[j] == 0 {
			*next = append(*next, j)
		}
	}
}
