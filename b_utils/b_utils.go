package b_utils

import (
	"math"
)

func MaxInt64(a int64, b int64) int64 {

	return int64(math.Max(float64(a), float64(b)))
}

func MaxInt(a int, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func MinInt64(a int64, b int64) int64 {
	return int64(math.Min(float64(a), float64(b)))
}

func MinInt(a int, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
