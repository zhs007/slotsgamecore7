package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ScriptCore(t *testing.T) {
	pool, err := NewGamePropertyPool("../data/game001/rtp96.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	gameProp := pool.newGameProp(10)
	assert.NotNil(t, gameProp)

	script, err := NewScriptCore(gameProp)
	assert.NoError(t, err)
	assert.NotNil(t, script)

	ast, iss := script.Cel.Compile(`setVal("abc",123)`)
	assert.NoError(t, iss.Err())
	assert.NotNil(t, ast)

	prg, err := script.Cel.Program(ast)
	assert.NoError(t, err)

	_, _, err = prg.Eval(map[string]any{})
	assert.NoError(t, err)

	t.Logf("Test_ScriptCore OK")
}
