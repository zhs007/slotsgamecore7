package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config_NewConfig(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg, "Test_Config_NewConfig: NewConfig should not return nil")
	assert.NotNil(t, cfg.Reels, "Test_Config_NewConfig: Reels map should be initialized")
	assert.NotNil(t, cfg.SWReels, "Test_Config_NewConfig: SWReels map should be initialized")
	assert.Nil(t, cfg.Lines, "Test_Config_NewConfig: Lines should be nil initially")
	assert.Nil(t, cfg.PayTables, "Test_Config_NewConfig: PayTables should be nil initially")
}

func Test_Config_SetDefaultSceneString(t *testing.T) {
	cfg := NewConfig()

	// Valid case
	validSceneStr := "[[1,2,3],[4,5,6]]"
	err := cfg.SetDefaultSceneString(validSceneStr)
	assert.NoError(t, err, "Test_Config_SetDefaultSceneString: Valid scene string should not produce an error")
	assert.NotNil(t, cfg.DefaultScene, "Test_Config_SetDefaultSceneString: DefaultScene should be set")
	assert.Equal(t, 3, cfg.DefaultScene.Height, "Test_Config_SetDefaultSceneString: Scene height should be 3")
	assert.Equal(t, 2, cfg.DefaultScene.Width, "Test_Config_SetDefaultSceneString: Scene width should be 2")

	// Invalid JSON
	invalidJSONStr := "[[1,2,3],[4,5,6]"
	err = cfg.SetDefaultSceneString(invalidJSONStr)
	assert.Error(t, err, "Test_Config_SetDefaultSceneString: Invalid JSON should produce an error")

	// Valid JSON, but wrong structure
	wrongStructureStr := `{"a": 1}`
	err = cfg.SetDefaultSceneString(wrongStructureStr)
	assert.Error(t, err, "Test_Config_SetDefaultSceneString: Wrong JSON structure should produce an error")
}

func Test_Config_AddDefaultSceneString2(t *testing.T) {
	cfg := NewConfig()

	// Valid case
	validSceneStr1 := "[[1,2],[3,4]]"
	err := cfg.AddDefaultSceneString2(validSceneStr1)
	assert.NoError(t, err, "Test_Config_AddDefaultSceneString2: First valid scene should not error")
	assert.Len(t, cfg.DefaultScene2, 1, "Test_Config_AddDefaultSceneString2: DefaultScene2 should have 1 element")
	assert.Equal(t, 2, cfg.DefaultScene2[0].Height, "Test_Config_AddDefaultSceneString2: First scene height")

	validSceneStr2 := "[[5,6,7],[8,9,10],[11,12,13]]"
	err = cfg.AddDefaultSceneString2(validSceneStr2)
	assert.NoError(t, err, "Test_Config_AddDefaultSceneString2: Second valid scene should not error")
	assert.Len(t, cfg.DefaultScene2, 2, "Test_Config_AddDefaultSceneString2: DefaultScene2 should have 2 elements")
	assert.Equal(t, 3, cfg.DefaultScene2[1].Height, "Test_Config_AddDefaultSceneString2: Second scene height")

	// Invalid JSON
	invalidJSONStr := "[[1,2]"
	err = cfg.AddDefaultSceneString2(invalidJSONStr)
	assert.Error(t, err, "Test_Config_AddDefaultSceneString2: Invalid JSON should produce an error")
	assert.Len(t, cfg.DefaultScene2, 2, "Test_Config_AddDefaultSceneString2: Length should not change after error")
}

func Test_Config_LoadErrors(t *testing.T) {
	cfg := NewConfig()

	// Test LoadLine with an invalid reel number
	err := cfg.LoadLine("dummy.json", 99)
	assert.Error(t, err, "Test_Config_LoadErrors: LoadLine with invalid reels should error")
	assert.Equal(t, ErrInvalidReels, err, "Test_Config_LoadErrors: LoadLine should return ErrInvalidReels")

	// Test LoadPayTables with an invalid reel number
	err = cfg.LoadPayTables("dummy.json", 99)
	assert.Error(t, err, "Test_Config_LoadErrors: LoadPayTables with invalid reels should error")
	assert.Equal(t, ErrInvalidReels, err, "Test_Config_LoadErrors: LoadPayTables should return ErrInvalidReels")

	// Test LoadReels with an invalid reel number
	err = cfg.LoadReels("dummyname", "dummy.json", 99)
	assert.Error(t, err, "Test_Config_LoadErrors: LoadReels with invalid reels should error")
	assert.Equal(t, ErrInvalidReels, err, "Test_Config_LoadErrors: LoadReels should return ErrInvalidReels")
}
