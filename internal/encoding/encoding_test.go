package encoding

import (
	"math"
	"math/rand"
	"testing"
)

func TestBinaryToReal_Basic(t *testing.T) {
	bits := []int{1, 1, 1, 1}
	result := BinaryToReal(bits, 4, 0, 10)
	if len(result) != 1 {
		t.Fatalf("expected 1 variable, got %d", len(result))
	}
	if math.Abs(result[0]-10.0) > 1e-10 {
		t.Fatalf("expected 10, got %v", result[0])
	}
}

func TestRealToBinary_Roundtrip(t *testing.T) {
	reals := []float64{2.5, 7.5}
	bits := RealToBinary(reals, 8, 0, 10)
	decoded := BinaryToReal(bits, 8, 0, 10)
	for i := range reals {
		if math.Abs(decoded[i]-reals[i]) > 0.05 {
			t.Fatalf("roundtrip mismatch at %d: %v vs %v", i, reals[i], decoded[i])
		}
	}
}

func TestGrayEncode_Decode(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := GrayEncode(i)
		back := GrayDecode(g)
		if back != i {
			t.Fatalf("gray roundtrip failed for %d: encoded=%d, decoded=%d", i, g, back)
		}
	}
}

func TestGrayToBinary_BinaryToGray(t *testing.T) {
	bin := []int{1, 0, 1, 1}
	gray := BinaryToGray(bin)
	back := GrayToBinary(gray)
	for i := range bin {
		if back[i] != bin[i] {
			t.Fatalf("bit %d mismatch: %d vs %d", i, bin[i], back[i])
		}
	}
}

func TestPermutationToEdges_Roundtrip(t *testing.T) {
	perm := []int{3, 1, 4, 0, 2}
	edges := PermutationToEdges(perm)
	back := EdgesToPermutation(edges, 3)
	if len(back) != len(perm) {
		t.Fatalf("length mismatch: %d vs %d", len(back), len(perm))
	}
	for i := range perm {
		if back[i] != perm[i] {
			t.Fatalf("position %d mismatch: %d vs %d", i, back[i], perm[i])
		}
	}
}

func TestRandomPerm_Valid(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	perm := RandomPerm(10, rnd)
	seen := make(map[int]bool)
	for _, v := range perm {
		if v < 0 || v >= 10 {
			t.Fatalf("value out of range: %d", v)
		}
		if seen[v] {
			t.Fatalf("duplicate value: %d", v)
		}
		seen[v] = true
	}
}

func TestIntegerEncode_Decode(t *testing.T) {
	values := []int{3, 7, 10}
	encoded := IntegerEncode(values, 0, 10)
	decoded := IntegerDecode(encoded, 0, 10)
	for i := range values {
		if decoded[i] != values[i] {
			t.Fatalf("position %d: expected %d, got %d", i, values[i], decoded[i])
		}
	}
}
