package sgc7utils

import "math"

var FloatPrecision float64 = 0.00000001

func IsFloatEquals(a, b float64) bool {
	return math.Abs(a-b) < FloatPrecision
}
