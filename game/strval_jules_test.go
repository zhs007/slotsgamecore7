package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StrVal_jules(t *testing.T) {
	// Test NewStrVal
	sv1 := NewStrVal()
	assert.NotNil(t, sv1)
	assert.Equal(t, "", sv1.String())
	assert.Equal(t, StrValType, sv1.Type())

	// Test NewStrValEx
	sv2 := NewStrValEx("123")
	assert.NotNil(t, sv2)
	assert.Equal(t, "123", sv2.String())

	// Test ParseString
	err := sv1.ParseString("456")
	assert.NoError(t, err)
	assert.Equal(t, "456", sv1.String())

	// Test IsSame
	assert.True(t, sv1.IsSame(NewStrValEx("456")))
	assert.False(t, sv1.IsSame(NewStrValEx("123")))
	assert.False(t, sv1.IsSame(NewIntValEx(456)))

	// Test Int32
	assert.Equal(t, int32(456), sv1.Int32())
	sv1.ParseString("abc")
	assert.Equal(t, int32(0), sv1.Int32())

	// Test Int64
	sv1.ParseString("789")
	assert.Equal(t, int64(789), sv1.Int64())
	sv1.ParseString("abc")
	assert.Equal(t, int64(0), sv1.Int64())

	// Test Int
	sv1.ParseString("101112")
	assert.Equal(t, int(101112), sv1.Int())
	sv1.ParseString("abc")
	assert.Equal(t, int(0), sv1.Int())

	// Test Float32
	sv1.ParseString("1.23")
	assert.InEpsilon(t, float32(1.23), sv1.Float32(), 1e-6)
	sv1.ParseString("abc")
	assert.Equal(t, float32(0), sv1.Float32())

	// Test Float64
	sv1.ParseString("4.56")
	assert.InEpsilon(t, float64(4.56), sv1.Float64(), 1e-6)
	sv1.ParseString("abc")
	assert.Equal(t, float64(0), sv1.Float64())

	// Test String
	sv1.ParseString("hello")
	assert.Equal(t, "hello", sv1.String())

	// Test Int32Arr
	sv1.ParseString("123")
	assert.Equal(t, []int32{123}, sv1.Int32Arr())

	// Test Int64Arr
	assert.Equal(t, []int64{123}, sv1.Int64Arr())

	// Test IntArr
	assert.Equal(t, []int{123}, sv1.IntArr())

	// Test Float32Arr
	sv1.ParseString("1.23")
	assert.Equal(t, []float32{1.23}, sv1.Float32Arr())

	// Test Float64Arr
	assert.Equal(t, []float64{1.23}, sv1.Float64Arr())

	// Test StringArr
	assert.Equal(t, []string{"1.23"}, sv1.StringArr())

	// Test GetInt
	sv1.ParseString("42")
	assert.Equal(t, 42, sv1.GetInt(0))
	assert.Equal(t, 0, sv1.GetInt(1))
}
