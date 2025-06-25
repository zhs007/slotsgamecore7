package sgc7utils

import "math"

// IsEqualFloat32 - checks if two float values are equal within a small epsilon
func IsEqualFloat32(a, b float32) bool {
	return math.Abs(float64(a-b)) <= 1e-6
}

// IsEqualFloat64 - checks if two float values are equal within a small epsilon
func IsEqualFloat64(a, b float64) bool {
	return math.Abs(a-b) <= 1e-12
}
