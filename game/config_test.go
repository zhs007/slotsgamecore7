package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	cfg := NewConfig()

	err := cfg.LoadLine5("../unittestdata/linedata.json")
	assert.Nil(t, err, "Test_Config LoadLine5")

	err = cfg.LoadLine5("../unittestdata/linedata1.json")
	assert.NotNil(t, err, "Test_Config LoadLine5")

	err = cfg.LoadPayTables5("../unittestdata/paytables.json")
	assert.Nil(t, err, "Test_Config LoadPayTables5")

	err = cfg.LoadPayTables5("../unittestdata/paytables1.json")
	assert.NotNil(t, err, "Test_Config LoadPayTables5")

	err = cfg.LoadReels5("bg", "../unittestdata/reels.json")
	assert.Nil(t, err, "Test_Config LoadReels5")

	err = cfg.LoadReels5("fg", "../unittestdata/reels2.json")
	assert.Nil(t, err, "Test_Config LoadReels5")

	err = cfg.LoadReels5("fg2", "../unittestdata/reels1.json")
	assert.NotNil(t, err, "Test_Config LoadReels5")

	t.Logf("Test_Config OK")
}
