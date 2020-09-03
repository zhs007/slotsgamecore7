package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	cfg := NewConfig()

	err := cfg.LoadLine5("../unittestdata/linedata.json")
	assert.NoError(t, err, "Test_Config LoadLine5")

	err = cfg.LoadLine5("../unittestdata/linedata1.json")
	assert.Error(t, err, "Test_Config LoadLine5")

	err = cfg.LoadPayTables5("../unittestdata/paytables.json")
	assert.NoError(t, err, "Test_Config LoadPayTables5")

	err = cfg.LoadPayTables5("../unittestdata/paytables1.json")
	assert.NotNil(t, err, "Test_Config LoadPayTables5")

	err = cfg.LoadReels5("bg", "../unittestdata/reels.json")
	assert.NoError(t, err, "Test_Config LoadReels5")

	err = cfg.LoadReels5("fg", "../unittestdata/reels2.json")
	assert.NoError(t, err, "Test_Config LoadReels5")

	err = cfg.LoadReels5("fg2", "../unittestdata/reels1.json")
	assert.Error(t, err, "Test_Config LoadReels5")

	err = cfg.SetDefaultSceneString("[[1,9,2],[10,7,1],[10,8,2],[4,6,1],[10,6,4]]")
	assert.NoError(t, err)

	err = cfg.SetDefaultSceneString("")
	assert.Error(t, err)

	t.Logf("Test_Config OK")
}
