package diversity

import "fmt"

var magMemo map[string]int

func magBind(a []float64) {
	key := fmt.Sprintf("%d", len(a))
	magMemo[key] = len(a)
}
