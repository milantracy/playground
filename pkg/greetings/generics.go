package greetings

type Number interface {
	int64 | float64
}

// SumInsOrFloats sums up the values of the given map m. It supports int64 and float64 as types of the map values.
func SumInsOrFloats[K comparable, V Number](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}
