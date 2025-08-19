package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IntArrVal_ParseString(t *testing.T) {
	val := NewIntArrVal[int]()
	err := val.ParseString("1,2, 3")
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, val.IntArr())

	t.Logf("Test_IntArrVal_ParseString OK")
}

func Test_IntArrVal_IsSame(t *testing.T) {
	val1 := NewIntArrVal[int]()
	val1.ParseString("1,2,3")

	val2 := NewIntArrVal[int]()
	val2.ParseString("1,2,3")

	assert.True(t, val1.IsSame(val2))

	val3 := NewIntArrVal[int]()
	val3.ParseString("1,2,4")
	assert.False(t, val1.IsSame(val3))

	t.Logf("Test_IntArrVal_IsSame OK")
}

func Test_IntArrVal_Conversions(t *testing.T) {
	val := NewIntArrVal[int]()
	val.ParseString("1,2,3")

	assert.Equal(t, int32(1), val.Int32())
	assert.Equal(t, int64(1), val.Int64())
	assert.Equal(t, int(1), val.Int())
	assert.Equal(t, float32(1), val.Float32())
	assert.Equal(t, float64(1), val.Float64())
	assert.Equal(t, "1,2,3", val.String())

	assert.Equal(t, []int32{1, 2, 3}, val.Int32Arr())
	assert.Equal(t, []int64{1, 2, 3}, val.Int64Arr())
	assert.Equal(t, []int{1, 2, 3}, val.IntArr())
	assert.Equal(t, []float32{1, 2, 3}, val.Float32Arr())
	assert.Equal(t, []float64{1, 2, 3}, val.Float64Arr())
	assert.Equal(t, []string{"1", "2", "3"}, val.StringArr())

	assert.Equal(t, 1, val.GetInt(0))
	assert.Equal(t, 2, val.GetInt(1))
	assert.Equal(t, 3, val.GetInt(2))

	t.Logf("Test_IntArrVal_Conversions OK")
}
