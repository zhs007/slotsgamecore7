package lowcode

import (
	"context"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// deterministic plugin for tests
type fakePluginTC struct{}

func (p *fakePluginTC) Random(_ context.Context, r int) (int, error) {
	if r <= 0 {
		return 0, nil
	}
	// always pick the first bucket
	return 0, nil
}
func (p *fakePluginTC) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginTC) ClearUsedRngs()                    {}
func (p *fakePluginTC) TagUsedRngs()                      {}
func (p *fakePluginTC) RollbackUsedRngs() error           { return nil }
func (p *fakePluginTC) SetCache(arr []int)                {}
func (p *fakePluginTC) ClearCache()                       {}
func (p *fakePluginTC) Init()                             {}
func (p *fakePluginTC) SetScenePool(any)                  {}
func (p *fakePluginTC) GetScenePool() any                 { return nil }
func (p *fakePluginTC) SetSeed(seed int)                  {}

func TestTreasureChestDataBasics(t *testing.T) {
	cd := &TreasureChestData{}
	// OnNewGame should initialize base fields
	cd.OnNewGame(&GameProperty{}, nil)
	assert.Nil(t, cd.Selected)
	assert.Equal(t, 0, cd.Output)

	cd.Selected = []int{1, 2}
	cd.Output = 5
	cd.onNewStep()
	assert.Empty(t, cd.Selected)
	assert.Equal(t, 0, cd.Output)

	// Clone
	cd.Selected = []int{3, 4}
	cd.Output = 7
	clone := cd.Clone().(*TreasureChestData)
	assert.Equal(t, cd.Output, clone.Output)
	assert.Equal(t, cd.Selected, clone.Selected)

	// BuildPBComponentData and GetValEx
	pb := cd.BuildPBComponentData()
	assert.NotNil(t, pb)
	v, ok := cd.GetValEx(CVNumber, 0)
	assert.True(t, ok)
	assert.Equal(t, cd.Output, v)

	v2, ok2 := cd.GetValEx("noexists", 0)
	assert.False(t, ok2)
	assert.Equal(t, 0, v2)
}

func TestTreasureChestConfigSetLink(t *testing.T) {
	cfg := &TreasureChestConfig{}
	cfg.SetLinkComponent("next", "n1")
	assert.Equal(t, "n1", cfg.DefaultNextComponent)

	cfg.SetLinkComponent("5", "b5")
	// MapBranchs should exist and contain 5
	if assert.NotNil(t, cfg.MapBranchs) {
		assert.Equal(t, "b5", cfg.MapBranchs[5])
	}
}

func TestInitExValidationAndWeightLoading(t *testing.T) {
	// prepare a valweights with one int val
	iv := sgc7game.NewIntValEx[int](10)
	vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv}, []int{1})
	assert.NoError(t, err)

	pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"w": vw}}

	// FragmentCollection with invalid FragmentNum
	comp := NewTreasureChest("tc1").(*TreasureChest)
	cfg := &TreasureChestConfig{StrType: "fragmentcollection", FragmentNum: 0, StrWeight: "w"}
	err = comp.InitEx(cfg, pool)
	assert.Error(t, err)

	// SumValue with invalid OpenNum > TotalNum
	cfg2 := &TreasureChestConfig{StrType: "sumvalue", OpenNum: 5, TotalNum: 3, StrWeight: "w"}
	err = comp.InitEx(cfg2, pool)
	assert.Error(t, err)

	// valid config should load weight from pool
	cfg3 := &TreasureChestConfig{StrType: "fragmentcollection", FragmentNum: 2, StrWeight: "w"}
	err = comp.InitEx(cfg3, pool)
	assert.NoError(t, err)
	assert.Equal(t, vw, comp.Config.Weight)
}

func TestGetWeightOverrideAndGetSymbolNum(t *testing.T) {
	iv1 := sgc7game.NewIntValEx[int](1)
	iv2 := sgc7game.NewIntValEx[int](2)
	vwDefault, _ := sgc7game.NewValWeights2([]sgc7game.IVal{iv1}, []int{1})
	vwOverride, _ := sgc7game.NewValWeights2([]sgc7game.IVal{iv2}, []int{1})

	comp := NewTreasureChest("tc2").(*TreasureChest)
	comp.Config = &TreasureChestConfig{Weight: vwDefault, OpenNum: 9}

	// basicCD without CCVWeight returns default
	bcd := &BasicComponentData{MapConfigVals: map[string]string{}}
	w, err := comp.getWeight(&GameProperty{Pool: &GamePropertyPool{}}, bcd)
	assert.NoError(t, err)
	assert.Equal(t, vwDefault, w)

	// override via basicCD CCVWeight
	pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"ov": vwOverride}}
	gp := &GameProperty{Pool: pool}
	bcd2 := &BasicComponentData{MapConfigVals: map[string]string{CCVWeight: "ov"}}
	w2, err2 := comp.getWeight(gp, bcd2)
	assert.NoError(t, err2)
	assert.Equal(t, vwOverride, w2)

	// getSymbolNum: prefer basicCD config int if set
	bcd3 := &BasicComponentData{MapConfigIntVals: map[string]int{CCVOpenNum: 5}}
	sn := comp.getSymbolNum(nil, bcd3)
	assert.Equal(t, 5, sn)
	// otherwise from config
	sn2 := comp.getSymbolNum(nil, &BasicComponentData{MapConfigIntVals: map[string]int{}})
	assert.Equal(t, 9, sn2)
}

func TestProcFragmentCollectionAndOnPlayGame(t *testing.T) {
	// prepare weights where the first value will be selected repeatedly
	v1 := sgc7game.NewIntValEx[int](100)
	v2 := sgc7game.NewIntValEx[int](200)
	vw, _ := sgc7game.NewValWeights2([]sgc7game.IVal{v1, v2}, []int{10, 1})

	comp := NewTreasureChest("tc3").(*TreasureChest)
	comp.Config = &TreasureChestConfig{
		Type:        TreasureChestTypeFragmentCollection,
		FragmentNum: 3,
		TotalNum:    9,
		Weight:      vw,
		MapBranchs:  map[int]string{100: "nextcomp"},
		MapControllers: map[int][]*Award{
			100: {{AwardType: ""}},
		},
	}

	gp := &GameProperty{Pool: &GamePropertyPool{}}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	plugin := &fakePluginTC{}
	gp.Pool = &GamePropertyPool{}
	gp.Pool = &GamePropertyPool{}

	// ensure BasicComponent config is initialized to avoid nil access in onStepEnd
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})
	cd := comp.NewComponentData().(*TreasureChestData)
	next, err := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
	assert.NoError(t, err)
	// next should be from MapBranchs
	assert.Equal(t, "nextcomp", next)
	// Selected length should equal FragmentNum
	assert.Equal(t, comp.Config.FragmentNum, len(cd.Selected))
	// Output should be the int value of the first val
	assert.Equal(t, 100, cd.Output)
}

func TestProcControllersBranches(t *testing.T) {
	comp := NewTreasureChest("tc4").(*TreasureChest)
	comp.Config = &TreasureChestConfig{Controllers: []*Award{{AwardType: ""}}, MapControllers: map[int][]*Award{1: {{AwardType: ""}}}}
	gp := &GameProperty{Pool: &GamePropertyPool{}}
	plugin := &fakePluginTC{}
	pr := &sgc7game.PlayResult{}
	params := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	// val == -1 && strVal == "sumValue" -> should call Controllers
	comp.ProcControllers(gp, plugin, pr, params, -1, "sumValue")
	// map controllers path
	comp.ProcControllers(gp, plugin, pr, params, 1, "")
}

func TestOnAsciiGameNewStats2(t *testing.T) {
	comp := NewTreasureChest("tc5").(*TreasureChest)
	// OnAsciiGame returns nil
	err := comp.OnAsciiGame(&GameProperty{}, sgc7game.NewPlayResult("m", 0, 0, "t"), nil, nil, comp.NewComponentData())
	assert.NoError(t, err)

	// OnStats2 & NewStats2
	comp.BasicComponent = NewBasicComponent("tc5", 0)
	s2 := stats2.NewCache(1)
	cd := comp.NewComponentData().(*TreasureChestData)
	comp.OnStats2(cd, s2, &GameProperty{}, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), sgc7game.NewPlayResult("m", 0, 0, "t"), false)
	f := comp.NewStats2("")
	assert.NotNil(t, f)
}

func TestJsonBuildAndParseInit(t *testing.T) {
	j := &jsonTreasureChest{Type: "fragmentCollection", FragmentNum: 3, OpenNum: 4, TotalNum: 5, Weight: "w"}
	cfg := j.build()
	assert.Equal(t, "fragmentcollection", cfg.StrType)
	assert.Equal(t, 3, cfg.FragmentNum)

	// test Init reads yaml file via Init
	tmp := "test_tc_init.yaml"
	_ = os.WriteFile(tmp, []byte("type: fragmentCollection\nfragmentNum: 2\nweight: w\n"), 0644)
	defer os.Remove(tmp)
	pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"w": sgc7game.NewValWeights2Ex()}}
	// make sure LoadIntWeights returns something (we keep empty ValWeights2 to avoid nil panics)
	comp := NewTreasureChest("tc6").(*TreasureChest)
	// Init should not panic; underlying yaml may not have all fields and InitEx will validate
	_ = comp.Init(tmp, pool)
}

func TestParseTreasureChestTypeAndPB(t *testing.T) {
	tt := parseTreasureChestType("sumvalue")
	assert.Equal(t, TreasureChestTypeSumValue, tt)
	tt2 := parseTreasureChestType("other")
	assert.Equal(t, TreasureChestTypeFragmentCollection, tt2)

	// BuildPBComponentData
	cd := &TreasureChestData{}
	cd.Selected = []int{1, 2, 3}
	cd.Output = 99
	pb := cd.BuildPBComponentData().(*sgc7pb.TreasureChestData)
	assert.Equal(t, int32(99), pb.BasicComponentData.Output)
	assert.Equal(t, int32(1), pb.Selected[0])
}

func TestProcSumValueAndOnPlayGameInvalidType(t *testing.T) {
	// prepare weights
	v1 := sgc7game.NewIntValEx[int](7)
	vw, _ := sgc7game.NewValWeights2([]sgc7game.IVal{v1}, []int{1})

	comp := NewTreasureChest("tc_sum").(*TreasureChest)
	comp.Config = &TreasureChestConfig{Type: TreasureChestTypeSumValue, OpenNum: 2, Weight: vw}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	gp := &GameProperty{Pool: &GamePropertyPool{}}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	plugin := &fakePluginTC{}

	cd := comp.NewComponentData().(*TreasureChestData)
	nc, err := comp.procSumValue(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), plugin, nil, nil, nil, cd)
	assert.NoError(t, err)
	assert.Equal(t, "", nc)
	assert.Equal(t, 14, cd.Output)

	// OnPlayGame with invalid type
	comp.Config.Type = TreasureChestType(999)
	_, err2 := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), plugin, "", "", nil, nil, nil, cd)
	assert.Error(t, err2)
}

func TestNewTreasureChestAndNewComponentData(t *testing.T) {
	comp := NewTreasureChest("ntc")
	assert.NotNil(t, comp)
	cd := comp.NewComponentData()
	assert.IsType(t, &TreasureChestData{}, cd)
}

func TestInitExWithControllers(t *testing.T) {
	iv := sgc7game.NewIntValEx[int](1)
	vw, _ := sgc7game.NewValWeights2([]sgc7game.IVal{iv}, []int{1})

	pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"w": vw}}

	comp := NewTreasureChest("tc_initc").(*TreasureChest)
	cfg := &TreasureChestConfig{
		StrType:        "fragmentCollection",
		FragmentNum:    1,
		StrWeight:      "w",
		MapControllers: map[int][]*Award{1: {{AwardType: "cash", Vals: []int{10}}}},
		Controllers:    []*Award{{AwardType: "respinTimes", Vals: []int{2}}},
	}

	err := comp.InitEx(cfg, pool)
	assert.NoError(t, err)
	assert.Equal(t, vw, comp.Config.Weight)
	// awards initialized
	assert.Equal(t, AwardCash, comp.Config.MapControllers[1][0].Type)
	assert.Equal(t, AwardRespinTimes, comp.Config.Controllers[0].Type)
}

func TestProcSumValue_RandErr(t *testing.T) {
	v1 := sgc7game.NewIntValEx[int](3)
	v2 := sgc7game.NewIntValEx[int](4)
	vw, _ := sgc7game.NewValWeights2([]sgc7game.IVal{v1, v2}, []int{1, 1})

	comp := NewTreasureChest("tc_sum_err").(*TreasureChest)
	comp.Config = &TreasureChestConfig{Type: TreasureChestTypeSumValue, OpenNum: 2, Weight: vw, Controllers: []*Award{{AwardType: "respinTimes", Vals: []int{1}}}}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	gp := &GameProperty{Pool: &GamePropertyPool{}}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	pluginErr := &fakePluginErrTC{}

	cd := comp.NewComponentData().(*TreasureChestData)
	_, err := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), pluginErr, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
	assert.Error(t, err)
}

func TestParseTreasureChest_NoControllersAndUnmarshalError(t *testing.T) {
	// no controller field
	jsonStr := `{
		"componentValues": {
			"label": "tcparse2",
			"configuration": {
				"type": "fragmentCollection",
				"fragmentNum": 1
			}
		}
	}`

	var node ast.Node
	err := sonic.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		t.Fatalf("sonic.Unmarshal err=%v", err)
	}

	bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
	label, err := parseTreasureChest(bc, &node)
	assert.NoError(t, err)
	assert.Equal(t, "tcparse2", label)
	_, ok := bc.mapConfig["tcparse2"]
	assert.True(t, ok)

	// configuration is not an object -> Unmarshal should error
	jsonBad := `{
		"componentValues": {
			"label": "bad",
			"configuration": "not-an-object"
		}
	}`
	var node2 ast.Node
	_ = sonic.Unmarshal([]byte(jsonBad), &node2)
	_, err2 := parseTreasureChest(&BetConfig{mapConfig: make(map[string]IComponentConfig)}, &node2)
	assert.Error(t, err2)
}

func TestGetWeight_ReturnsNilWhenMissing(t *testing.T) {
	comp := NewTreasureChest("tcgw").(*TreasureChest)
	comp.Config = &TreasureChestConfig{Weight: nil}

	gp := &GameProperty{Pool: &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{}}}
	bcd := &BasicComponentData{MapConfigVals: map[string]string{CCVWeight: "noexists"}}

	vw, err := comp.getWeight(gp, bcd)
	assert.NoError(t, err)
	assert.Nil(t, vw)
}

func TestInit_ErrorPaths(t *testing.T) {
	comp := NewTreasureChest("tinit").(*TreasureChest)
	// non-existent file
	err := comp.Init("/no/such/file.yaml", &GamePropertyPool{})
	assert.Error(t, err)

	// bad yaml file
	tmp := "test_tc_bad.yaml"
	_ = os.WriteFile(tmp, []byte("type: [bad"), 0644)
	defer os.Remove(tmp)
	err2 := comp.Init(tmp, &GamePropertyPool{})
	assert.Error(t, err2)
}

func TestParseTreasureChest_ControllerParseError(t *testing.T) {
	// controller with stringVal that isn't an int to force String2Int64 error
	jsonStr := `{
		"componentValues": {
			"label": "tcparse3",
			"configuration": {
				"type": "fragmentCollection",
				"fragmentNum": 1
			},
			"controller": [ { "stringVal": "notint", "type": "addRespinTimes" } ]
		}
	}`

	var node ast.Node
	err := sonic.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		t.Fatalf("sonic.Unmarshal err=%v", err)
	}

	_, err2 := parseTreasureChest(&BetConfig{mapConfig: make(map[string]IComponentConfig)}, &node)
	assert.Error(t, err2)
}

// plugin that returns error from Random
type fakePluginErrTC struct{}

func (p *fakePluginErrTC) Random(_ context.Context, r int) (int, error) {
	return 0, sgc7game.ErrInvalidValWeights
}
func (p *fakePluginErrTC) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginErrTC) ClearUsedRngs()                    {}
func (p *fakePluginErrTC) TagUsedRngs()                      {}
func (p *fakePluginErrTC) RollbackUsedRngs() error           { return nil }
func (p *fakePluginErrTC) SetCache(arr []int)                {}
func (p *fakePluginErrTC) ClearCache()                       {}
func (p *fakePluginErrTC) Init()                             {}
func (p *fakePluginErrTC) SetScenePool(any)                  {}
func (p *fakePluginErrTC) GetScenePool() any                 { return nil }
func (p *fakePluginErrTC) SetSeed(seed int)                  {}

func TestProcFragmentCollection_EmptyWeightsAndRandErr(t *testing.T) {
	// empty weight should return ErrInvalidComponentConfig
	comp := NewTreasureChest("tc_empty").(*TreasureChest)
	comp.Config = &TreasureChestConfig{Type: TreasureChestTypeFragmentCollection, FragmentNum: 1, Weight: sgc7game.NewValWeights2Ex()}

	gp := &GameProperty{Pool: &GamePropertyPool{}}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	plugin := &fakePluginTC{}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	cd := comp.NewComponentData().(*TreasureChestData)
	_, err := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
	t.Logf("err1=%v", err)
	assert.Error(t, err)

	// RandVal error path
	v1 := sgc7game.NewIntValEx[int](9)
	v2 := sgc7game.NewIntValEx[int](8)
	vw, _ := sgc7game.NewValWeights2([]sgc7game.IVal{v1, v2}, []int{1, 1})
	comp.Config.Weight = vw
	pluginErr := &fakePluginErrTC{}
	cd2 := comp.NewComponentData().(*TreasureChestData)
	_, err2 := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), pluginErr, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd2)
	t.Logf("err2=%v", err2)
	assert.Error(t, err2)
}

func TestParseTreasureChest_JSON(t *testing.T) {
	jsonStr := `{
	    "componentValues": {
	        "label": "tcparse",
	        "configuration": {
	            "type": "fragmentCollection",
	            "fragmentNum": 2,
	            "weight": "w"
	        },
	        "controller": [
	            {"type": "addRespinTimes", "target": "tg", "times": 1}
	        ]
	    }
	}`

	var node ast.Node
	err := sonic.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		t.Fatalf("sonic.Unmarshal err=%v", err)
	}

	bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
	label, err := parseTreasureChest(bc, &node)
	assert.NoError(t, err)
	assert.Equal(t, "tcparse", label)
	_, ok := bc.mapConfig["tcparse"]
	assert.True(t, ok)

	// parseTreasureChest with nil cell should error
	_, err2 := parseTreasureChest(&BetConfig{mapConfig: make(map[string]IComponentConfig)}, nil)
	assert.Error(t, err2)
}
