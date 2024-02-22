package mathtoolset2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ScriptCore_genStackReels(t *testing.T) {
	file, err := os.Open("../unittestdata/reelsstats2.xlsx")
	assert.NoError(t, err)

	mapfd, err := NewFileDataMap("")
	assert.NoError(t, err)

	mapfd.AddReader("reelsstats2.xlsx", file)

	jsondata, err := mapfd.ToJson()
	assert.NoError(t, err)

	sc, err := NewScriptCore(jsondata)
	assert.NoError(t, err)

	err = sc.Run("genStackReels(\"output.xlsx\", \"reelsstats2.xlsx\", [2, 3], [\"SC\"])")
	assert.NoError(t, err)

	assert.Equal(t, len(sc.ErrInRun), 0)
	assert.NotNil(t, sc.MapOutputFiles.MapFiles["output.xlsx"])

	t.Logf("Test_ScriptCore_genStackReels OK")
}
