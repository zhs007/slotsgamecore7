package lowcode

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

// Test that parseGenMask correctly loads config from a JSON cell (componentValues)
func Test_ParseGenMask_JSON(t *testing.T) {
	// minimal BetConfig context with initialized maps
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "gm_json",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	// a single cell JSON mimicking the designer output for genMask
	js := []byte(`{
        "componentValues": {
            "label": "gm_json",
            "configuration": {
                "type": "Random",
                "maskLen": 3,
                "outputMask": "out_json",
                "weightValue": 10000,
                "initMask": [1, 0, 1]
            }
        }
    }`)

	// parse into an ast.Node
	node, err := sonic.Get(js)
	assert.NoError(t, err)

	// invoke component JSON loader
	name, err := parseGenMask(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "gm_json", name)

	// verify structures are populated
	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	gmCfg, ok := cfgIface.(*GenMaskConfig)
	assert.True(t, ok)

	// field assertions
	assert.Equal(t, "random", gmCfg.StrType)
	assert.Equal(t, 3, gmCfg.MaskLen)
	assert.Equal(t, "out_json", gmCfg.OutputMask)
	assert.Equal(t, 10000, gmCfg.WeightValue)
	assert.Equal(t, []bool{true, false, true}, gmCfg.InitMask)

	// component list entry exists
	assert.Len(t, betCfg.Components, 1)
	assert.Equal(t, "gm_json", betCfg.Components[0].Name)
	assert.Equal(t, GenMaskTypeName, betCfg.Components[0].Type)
}

// Test that parseGenMask also handles srcMask and controller array
func Test_ParseGenMask_JSON_WithSrcAndControllers(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "gm_json2",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "gm_json2",
			"configuration": {
				"type": "Or",
				"maskLen": 2,
				"outputMask": "out_json2",
				"srcMask": ["m1", "m2"]
			},
			"controller": [
				{ "type": "addRespinTimes", "target": "fg-start", "times": 2 }
			]
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseGenMask(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "gm_json2", name)

	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	gmCfg, ok := cfgIface.(*GenMaskConfig)
	assert.True(t, ok)

	assert.Equal(t, "or", gmCfg.StrType)
	assert.Equal(t, 2, gmCfg.MaskLen)
	assert.Equal(t, "out_json2", gmCfg.OutputMask)
	assert.ElementsMatch(t, []string{"m1", "m2"}, gmCfg.SrcMask)

	// controllers parsed
	if assert.Len(t, gmCfg.Controllers, 1) {
		aw := gmCfg.Controllers[0]
		assert.Equal(t, "respinTimes", aw.AwardType)
		assert.Equal(t, []int{2}, aw.Vals)
		assert.Equal(t, []string{"fg-start"}, aw.StrParams)
	}

	// component list entry exists
	assert.Len(t, betCfg.Components, 1)
	assert.Equal(t, GenMaskTypeName, betCfg.Components[0].Type)
}

// Test that parseGenMask returns error on unsupported controller type
func Test_ParseGenMask_JSON_InvalidController(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "gm_badctrl",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "gm_badctrl",
			"configuration": {
				"type": "Or",
				"maskLen": 2,
				"outputMask": "out_badctrl",
				"srcMask": ["m1", "m2"]
			},
			"controller": [
				{ "type": "__unsupported__", "times": 1 }
			]
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	_, err = parseGenMask(betCfg, &node)
	assert.Error(t, err)
}

// Test that a single-mask config with both srcMask and initMask parses,
// but will be rejected by InitEx later (ensuring pipeline catches it)
func Test_ParseGenMask_JSON_InitExRejectsConflict(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "gm_conflict",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "gm_conflict",
			"configuration": {
				"type": "Set",
				"maskLen": 3,
				"outputMask": "out_conflict",
				"srcMask": ["onlyone"],
				"initMask": [1,0,1]
			}
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseGenMask(betCfg, &node)
	assert.NoError(t, err)
	cfgIface := betCfg.mapConfig[name]
	gmCfg := cfgIface.(*GenMaskConfig)

	// Now feed this config into InitEx which should reject it
	gm := &GenMask{BasicComponent: NewBasicComponent("gm_conflict", 0)}
	err = gm.InitEx(gmCfg, nil)
	assert.Error(t, err)
}

// Test that missing required fields (e.g., outputMask) parse but fail InitEx
func Test_ParseGenMask_JSON_InitEx_MissingOutputMask(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "gm_missing",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "gm_missing",
			"configuration": {
				"type": "Random",
				"maskLen": 3,
				"weightValue": 5000
			}
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseGenMask(betCfg, &node)
	assert.NoError(t, err)
	cfgIface := betCfg.mapConfig[name]
	gmCfg := cfgIface.(*GenMaskConfig)

	gm := &GenMask{BasicComponent: NewBasicComponent("gm_missing", 0)}
	err = gm.InitEx(gmCfg, nil)
	assert.Error(t, err)
}
