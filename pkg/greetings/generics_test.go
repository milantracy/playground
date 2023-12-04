package greetings

import (
"math"
	"testing"
)

type testParam[V Number] struct {
	name     string
	input    map[int]V
	expected V
}

func withinTolerance(a, b, e float64) bool {
	d := math.Abs(a - b)
	if b == 0 {
		return d < e
	}
	return (d /math.Abs(b)) < e
}

func TestSumInt(t *testing.T) {
	for _, tst := range []testParam[int64]{
		{
			name:     "SumInts",
			input:    map[int]int64{1: 100, 2: 200, 3: 300, 4: 1000},
			expected: 1600,
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			if sum := SumInsOrFloats(tst.input); sum != tst.expected {
				t.Errorf("SumInsOrFloats() = %v, expect %v", sum, tst.expected)
			}
		})
	}
}

func TestSumFloats(t *testing.T) {
	for _, tst := range []testParam[float64]{
		{
			name:     "SumInts",
			input:    map[int]float64{1: 200, 2: 250, 3: 350, 4: 1000},
			expected: 1800,
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			if sum := SumInsOrFloats(tst.input); !withinTolerance(sum, tst.expected, 1e-10) {
				t.Errorf("SumInsOrFloats() = %v, expect %v", sum, tst.expected)
			}
		})
	}
}
