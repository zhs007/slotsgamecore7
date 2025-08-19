package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IntVal_ParseString(t *testing.T) {
	val := NewIntVal[int]()
	err := val.ParseString("123")
	assert.NoError(t, err)
	assert.Equal(t, 123, val.Int())

	t.Logf("Test_IntVal_ParseString OK")
}

func Test_IntVal_IsSame(t *testing.T) {
	val1 := NewIntValEx[int](123)
	val2 := NewIntValEx[int](123)

	assert.True(t, val1.IsSame(val2))

	val3 := NewIntValEx[int](456)
	assert.False(t, val1.IsSame(val3))

	t.Logf("Test_IntVal_IsSame OK")
}

func Test_IntVal_Conversions(t *testing.T) {
	val := NewIntValEx[int](123)

	assert.Equal(t, int32(123), val.Int32())
	assert.Equal(t, int64(123), val.Int64())
	assert.Equal(t, int(123), val.Int())
	assert.Equal(t, float32(123), val.Float32())
	assert.Equal(t, float64(123), val.Float64())
	assert.Equal(t, "123", val.String())

	assert.Equal(t, []int32{123}, val.Int32Arr())
	assert.Equal(t, []int64{123}, val.Int64Arr())
	assert.Equal(t, []int{123}, val.IntArr())
	assert.Equal(t, []float32{123}, val.Float32Arr())
	assert.Equal(t, []float64{123}, val.Float64Arr())
	assert.Equal(t, []string{"123"}, val.StringArr())

	assert.Equal(t, 123, val.GetInt(0))
	assert.Equal(t, 0, val.GetInt(1))

	t.Logf("Test_IntVal_Conversions OK")
}
