package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NextGameModParams(t *testing.T) {
	type TestParams struct {
		A int     `json:"a"`
		B string  `json:"b"`
		C float32 `json:"c"`
		D []int   `json:"d"`
	}

	tp := &TestParams{A: 123, B: "456", C: 7.89, D: []int{10, 11, 12}}
	pr := &PlayResult{
		CurGameModParams: tp,
	}

	buf, err := PlayResult2JSON(pr)
	assert.NoError(t, err)
	t.Logf("%s", string(buf))

	pr2 := &PlayResult{
		CurGameModParams: &TestParams{},
	}
	pr1, err := JSON2PlayResult(buf, pr2)
	assert.NoError(t, err)

	tp1, isok := pr1.CurGameModParams.(*TestParams)
	assert.Equal(t, isok, true)
	assert.NotNil(t, tp1)

	assert.Equal(t, tp1.A, 123)
	assert.Equal(t, tp1.B, "456")
	assert.Equal(t, tp1.C, float32(7.89))
	assert.Equal(t, tp1.D[0], 10)
	assert.Equal(t, tp1.D[1], 11)
	assert.Equal(t, tp1.D[2], 12)

	t.Logf("Test_NextGameModParams OK")
}
