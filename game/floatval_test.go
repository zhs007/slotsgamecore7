package sgc7game

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// almostEqual checks if two float values are equal within a small epsilon
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-6
}

func Test_FloatVal32(t *testing.T) {
	// Test NewFloatVal with float32
	val1 := NewFloatVal[float32]()
	assert.Equal(t, FloatValType, val1.Type())
	assert.Equal(t, float32(0), val1.Float32())

	// Test NewFloatValEx with float32
	val2 := NewFloatValEx[float32](3.14)
	assert.Equal(t, FloatValType, val2.Type())
	assert.Equal(t, float32(3.14), val2.Float32())

	// Test ParseString
	err := val1.ParseString("3.14")
	assert.NoError(t, err)
	assert.Equal(t, float32(3.14), val1.Float32())

	// Test invalid string parsing
	err = val1.ParseString("invalid")
	assert.Error(t, err)

	// Test type conversions
	val3 := NewFloatValEx[float32](123.45)
	assert.Equal(t, int32(123), val3.Int32())
	assert.Equal(t, int64(123), val3.Int64())
	assert.Equal(t, int(123), val3.Int())
	assert.Equal(t, float32(123.45), val3.Float32())
	// When converting float32 to float64, we need to compare with the actual float32 value
	// First convert 123.45 to float32 to get the exact value we expect
	expectedFloat32 := float32(123.45)
	// Then convert this float32 value to float64
	expectedFloat64 := float64(expectedFloat32)
	assert.True(t, almostEqual(expectedFloat64, val3.Float64()))
	assert.Equal(t, "123.45", val3.String())

	// Test array conversions
	assert.Equal(t, []int32{123}, val3.Int32Arr())
	assert.Equal(t, []int64{123}, val3.Int64Arr())
	assert.Equal(t, []int{123}, val3.IntArr())
	assert.Equal(t, []float32{expectedFloat32}, val3.Float32Arr())
	float64Arr := val3.Float64Arr()
	assert.Equal(t, 1, len(float64Arr))
	assert.True(t, almostEqual(expectedFloat64, float64Arr[0]))
	assert.Equal(t, []string{"123.45"}, val3.StringArr())

	// Test GetInt
	assert.Equal(t, 123, val3.GetInt(0))
	assert.Equal(t, 0, val3.GetInt(1)) // Out of bounds should return 0
}

func Test_FloatVal64(t *testing.T) {
	// Test NewFloatVal with float64
	val1 := NewFloatVal[float64]()
	assert.Equal(t, FloatValType, val1.Type())
	assert.Equal(t, float64(0), val1.Float64())

	// Test NewFloatValEx with float64
	val2 := NewFloatValEx[float64](3.14)
	assert.Equal(t, FloatValType, val2.Type())
	assert.Equal(t, float64(3.14), val2.Float64())

	// Test IsSame
	val3 := NewFloatValEx[float64](3.14)
	assert.True(t, val2.IsSame(val3))
	assert.True(t, val3.IsSame(val2))

	val4 := NewFloatValEx[float64](2.0)
	assert.False(t, val2.IsSame(val4))

	// Test ParseString
	err := val1.ParseString("3.14")
	assert.NoError(t, err)
	assert.Equal(t, float64(3.14), val1.Float64())

	// Test type conversions
	val5 := NewFloatValEx[float64](123.45)
	assert.Equal(t, int32(123), val5.Int32())
	assert.Equal(t, int64(123), val5.Int64())
	assert.Equal(t, int(123), val5.Int())
	assert.Equal(t, float32(123.45), val5.Float32())
	assert.Equal(t, float64(123.45), val5.Float64())
	assert.Equal(t, "123.45", val5.String())

	// Test array conversions
	assert.Equal(t, []int32{123}, val5.Int32Arr())
	assert.Equal(t, []int64{123}, val5.Int64Arr())
	assert.Equal(t, []int{123}, val5.IntArr())
	assert.Equal(t, []float32{123.45}, val5.Float32Arr())
	assert.Equal(t, []float64{123.45}, val5.Float64Arr())
	assert.Equal(t, []string{"123.45"}, val5.StringArr())

	// Test GetInt
	assert.Equal(t, 123, val5.GetInt(0))
	assert.Equal(t, 0, val5.GetInt(1)) // Out of bounds should return 0
}
