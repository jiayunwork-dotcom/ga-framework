// Package encoding 提供不同编码方案之间的转换工具。
package encoding

import (
	"math"
	"math/rand"
)

// BinaryToReal 将二进制串按 bitsPerVar 切分，映射到 [lo, hi]。
func BinaryToReal(bits []int, bitsPerVar int, lo, hi float64) []float64 {
	if bitsPerVar < 1 {
		return nil
	}
	n := len(bits) / bitsPerVar
	result := make([]float64, n)
	maxVal := math.Pow(2, float64(bitsPerVar)) - 1
	span := hi - lo
	for i := 0; i < n; i++ {
		val := 0.0
		for j := 0; j < bitsPerVar; j++ {
			idx := i*bitsPerVar + j
			if idx < len(bits) && bits[idx] == 1 {
				val += math.Pow(2, float64(bitsPerVar-1-j))
			}
		}
		result[i] = lo + (val/maxVal)*span
	}
	return result
}

// RealToBinary 将实数序列编码为二进制串。
func RealToBinary(reals []float64, bitsPerVar int, lo, hi float64) []int {
	if bitsPerVar < 1 || hi <= lo {
		return nil
	}
	bits := make([]int, len(reals)*bitsPerVar)
	maxVal := math.Pow(2, float64(bitsPerVar)) - 1
	span := hi - lo
	for i, r := range reals {
		if r < lo {
			r = lo
		}
		if r > hi {
			r = hi
		}
		intVal := int(math.Round((r - lo) / span * maxVal))
		for j := bitsPerVar - 1; j >= 0; j-- {
			bits[i*bitsPerVar+j] = intVal & 1
			intVal >>= 1
		}
	}
	return bits
}

// GrayEncode 将整数转换为格雷码。
func GrayEncode(n int) int {
	return n ^ (n >> 1)
}

// GrayDecode 将格雷码转换回整数。
func GrayDecode(g int) int {
	n := g
	for g >>= 1; g != 0; g >>= 1 {
		n ^= g
	}
	return n
}

// GrayToBinary 将格雷码位串转为标准二进制位串。
func GrayToBinary(gray []int) []int {
	if len(gray) == 0 {
		return nil
	}
	bin := make([]int, len(gray))
	bin[0] = gray[0]
	for i := 1; i < len(gray); i++ {
		bin[i] = bin[i-1] ^ gray[i]
	}
	return bin
}

// BinaryToGray 将标准二进制位串转为格雷码位串。
func BinaryToGray(bin []int) []int {
	if len(bin) == 0 {
		return nil
	}
	gray := make([]int, len(bin))
	gray[0] = bin[0]
	for i := 1; i < len(bin); i++ {
		gray[i] = bin[i-1] ^ bin[i]
	}
	return gray
}

// PermutationToEdges 将排列编码转为邻接表示（用于 TSP）。
func PermutationToEdges(perm []int) map[int]int {
	edges := make(map[int]int, len(perm))
	for i := 0; i < len(perm); i++ {
		next := (i + 1) % len(perm)
		edges[perm[i]] = perm[next]
	}
	return edges
}

// EdgesToPermutation 将邻接表示转回排列。
func EdgesToPermutation(edges map[int]int, start int) []int {
	perm := make([]int, 0, len(edges))
	visited := make(map[int]bool)
	cur := start
	for len(perm) < len(edges) {
		if visited[cur] {
			break
		}
		visited[cur] = true
		perm = append(perm, cur)
		cur = edges[cur]
	}
	return perm
}

// RandomPerm 生成 [0,n) 的随机排列。
func RandomPerm(n int, rnd *rand.Rand) []int {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	rnd.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}

// IntegerEncode 将 [lo,hi] 范围的整数值编码为规范化 [0,1] 浮点数。
func IntegerEncode(values []int, lo, hi int) []float64 {
	span := float64(hi - lo)
	if span == 0 {
		return make([]float64, len(values))
	}
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = float64(v-lo) / span
	}
	return result
}

// IntegerDecode 将 [0,1] 浮点数解码为 [lo,hi] 整数。
func IntegerDecode(encoded []float64, lo, hi int) []int {
	span := float64(hi - lo)
	result := make([]int, len(encoded))
	for i, e := range encoded {
		result[i] = lo + int(math.Round(e*span))
		if result[i] < lo {
			result[i] = lo
		}
		if result[i] > hi {
			result[i] = hi
		}
	}
	return result
}
