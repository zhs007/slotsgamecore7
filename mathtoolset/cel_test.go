package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ScriptCore(t *testing.T) {
	mgrGenMath := NewGamMathMgr()
	assert.NotNil(t, mgrGenMath)

	script, err := NewScriptCore(mgrGenMath)
	assert.NoError(t, err)
	assert.NotNil(t, script)

	err = script.Compile(`calcLineRTP("../unittestdata/paytables.xlsx","../unittestdata/rss.xlsx",[1,2,3,4,5,6,7,8,9],[0],10,10)`)
	assert.NoError(t, err)

	out, err := script.Eval(mgrGenMath)
	assert.NoError(t, err)
	assert.True(t, out.Value().(float64) > 0)

	err = script.Compile(`calcLineRTP("../unittestdata/paytables.xlsx","../unittestdata/rss.xlsx",[1,2,3,4,5,6,7,8,9],[0],10,10) + calcScatterRTP("../unittestdata/paytables.xlsx","../unittestdata/rss.xlsx",[10],3)`)
	assert.NoError(t, err)

	out1, err := script.Eval(mgrGenMath)
	assert.NoError(t, err)
	assert.True(t, out1.Value().(float64) > out.Value().(float64))

	t.Logf("Test_ScriptCore OK")
}
