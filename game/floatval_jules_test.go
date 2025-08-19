package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FloatVal_Float32(t *testing.T) {
	// NewFloatVal
	v1 := NewFloatVal[float32]()
	assert.NotNil(t, v1, "Test_FloatVal_Float32: NewFloatVal")
	assert.Equal(t, v1.Float32(), float32(0), "Test_FloatVal_Float32: NewFloatVal value")

	// NewFloatValEx
	v2 := NewFloatValEx[float32](123.45)
	assert.NotNil(t, v2, "Test_FloatVal_Float32: NewFloatValEx")
	assert.Equal(t, v2.Float32(), float32(123.45), "Test_FloatVal_Float32: NewFloatValEx value")

	// Type
	assert.Equal(t, v1.Type(), FloatValType, "Test_FloatVal_Float32: Type")

	// IsSame
	v3 := NewFloatValEx[float32](123.45)
	assert.True(t, v2.IsSame(v3), "Test_FloatVal_Float32: IsSame true")
	v4 := NewFloatValEx[float32](543.21)
	assert.False(t, v2.IsSame(v4), "Test_FloatVal_Float32: IsSame false")
	v5 := NewIntValEx(123)
	assert.False(t, v2.IsSame(v5), "Test_FloatVal_Float32: IsSame different type")

	// ParseString
	err := v1.ParseString("67.89")
	assert.NoError(t, err, "Test_FloatVal_Float32: ParseString no error")
	assert.Equal(t, v1.Float32(), float32(67.89), "Test_FloatVal_Float32: ParseString value")
	err = v1.ParseString("abc")
	assert.Error(t, err, "Test_FloatVal_Float32: ParseString error")

	// Conversions
	fv := NewFloatValEx[float32](-98.7)
	assert.Equal(t, fv.Int32(), int32(-98), "Test_FloatVal_Float32: Int32")
	assert.Equal(t, fv.Int64(), int64(-98), "Test_FloatVal_Float32: Int64")
	assert.Equal(t, fv.Int(), int(-98), "Test_FloatVal_Float32: Int")
	assert.Equal(t, fv.Float32(), float32(-98.7), "Test_FloatVal_Float32: Float32")
	assert.Equal(t, fv.Float64(), float64(float32(-98.7)), "Test_FloatVal_Float32: Float64") // Be careful with precision
	assert.Equal(t, fv.String(), "-98.7", "Test_FloatVal_Float32: String")

	// Array Conversions
	assert.Equal(t, fv.Int32Arr(), []int32{-98}, "Test_FloatVal_Float32: Int32Arr")
	assert.Equal(t, fv.Int64Arr(), []int64{-98}, "Test_FloatVal_Float32: Int64Arr")
	assert.Equal(t, fv.IntArr(), []int{-98}, "Test_FloatVal_Float32: IntArr")
	assert.Equal(t, fv.Float32Arr(), []float32{-98.7}, "Test_FloatVal_Float32: Float32Arr")
	assert.Equal(t, fv.Float64Arr(), []float64{float64(float32(-98.7))}, "Test_FloatVal_Float32: Float64Arr")
	assert.Equal(t, fv.StringArr(), []string{"-98.7"}, "Test_FloatVal_Float32: StringArr")

	// GetInt
	assert.Equal(t, fv.GetInt(0), -98, "Test_FloatVal_Float32: GetInt(0)")
	assert.Equal(t, fv.GetInt(1), 0, "Test_FloatVal_Float32: GetInt(1)")
	assert.Equal(t, fv.GetInt(-1), 0, "Test_FloatVal_Float32: GetInt(-1)")
}

func Test_FloatVal_Float64(t *testing.T) {
	// NewFloatVal
	v1 := NewFloatVal[float64]()
	assert.NotNil(t, v1, "Test_FloatVal_Float64: NewFloatVal")
	assert.Equal(t, v1.Float64(), float64(0), "Test_FloatVal_Float64: NewFloatVal value")

	// NewFloatValEx
	v2 := NewFloatValEx[float64](123.456789)
	assert.NotNil(t, v2, "Test_FloatVal_Float64: NewFloatValEx")
	assert.Equal(t, v2.Float64(), 123.456789, "Test_FloatVal_Float64: NewFloatValEx value")

	// Type
	assert.Equal(t, v1.Type(), FloatValType, "Test_FloatVal_Float64: Type")

	// IsSame
	v3 := NewFloatValEx[float64](123.456789)
	assert.True(t, v2.IsSame(v3), "Test_FloatVal_Float64: IsSame true")
	v4 := NewFloatValEx[float64](987.654321)
	assert.False(t, v2.IsSame(v4), "Test_FloatVal_Float64: IsSame false")
	v5 := NewIntValEx(123)
	assert.False(t, v2.IsSame(v5), "Test_FloatVal_Float64: IsSame different type")

	// ParseString
	err := v1.ParseString("67.890123")
	assert.NoError(t, err, "Test_FloatVal_Float64: ParseString no error")
	assert.Equal(t, v1.Float64(), 67.890123, "Test_FloatVal_Float64: ParseString value")
	err = v1.ParseString("xyz")
	assert.Error(t, err, "Test_FloatVal_Float64: ParseString error")

	// Conversions
	fv := NewFloatValEx[float64](-98.765)
	assert.Equal(t, fv.Int32(), int32(-98), "Test_FloatVal_Float64: Int32")
	assert.Equal(t, fv.Int64(), int64(-98), "Test_FloatVal_Float64: Int64")
	assert.Equal(t, fv.Int(), int(-98), "Test_FloatVal_Float64: Int")
	assert.Equal(t, fv.Float32(), float32(-98.765), "Test_FloatVal_Float64: Float32")
	assert.Equal(t, fv.Float64(), float64(-98.765), "Test_FloatVal_Float64: Float64")
	assert.Equal(t, fv.String(), "-98.765", "Test_FloatVal_Float64: String")

	// Array Conversions
	assert.Equal(t, fv.Int32Arr(), []int32{-98}, "Test_FloatVal_Float64: Int32Arr")
	assert.Equal(t, fv.Int64Arr(), []int64{-98}, "Test_FloatVal_Float64: Int64Arr")
	assert.Equal(t, fv.IntArr(), []int{-98}, "Test_FloatVal_Float64: IntArr")
	assert.Equal(t, fv.Float32Arr(), []float32{-98.765}, "Test_FloatVal_Float64: Float32Arr")
	assert.Equal(t, fv.Float64Arr(), []float64{-98.765}, "Test_FloatVal_Float64: Float64Arr")
	assert.Equal(t, fv.StringArr(), []string{"-98.765"}, "Test_FloatVal_Float64: StringArr")

	// GetInt
	assert.Equal(t, fv.GetInt(0), -98, "Test_FloatVal_Float64: GetInt(0)")
	assert.Equal(t, fv.GetInt(1), 0, "Test_FloatVal_Float64: GetInt(1)")
	assert.Equal(t, fv.GetInt(-1), 0, "Test_FloatVal_Float64: GetInt(-1)")
}
