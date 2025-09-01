package lowcode

import (
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/types/known/anypb"
)

// fake component used in tests
type fakeComp struct{ *BasicComponent }

func (f *fakeComp) Init(fn string, pool *GamePropertyPool) error { return nil }
func (f *fakeComp) InitEx(cfg any, pool *GamePropertyPool) error { return nil }
func (f *fakeComp) OnGameInited(components *ComponentList) error { return nil }
func (f *fakeComp) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {
	return "", nil
}

func TestParseWinResultModifierJSON(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w1",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "w1",
			"configuration": {
				"type": "Multiply",
				"sourceComponent": ["src"],
				"winMulti": 3,
				"targetSymbols": ["A"]
			}
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseWinResultModifier(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "w1", name)

	// verify mapConfig populated
	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	cfg := cfgIface.(*WinResultModifierConfig)
	assert.Equal(t, "multiply", cfg.StrType)
	assert.Equal(t, 3, cfg.WinMulti)
	assert.Equal(t, []string{"A"}, cfg.TargetSymbols)
}

func TestInitExInvalidTargetSymbol(t *testing.T) {
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cfg := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "existsymbol", SourceComponents: []string{"src"}, TargetSymbols: []string{"Z"}}

	// pool only has symbol A
	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Equal(t, ErrInvalidSymbol, err)
}

func TestParseWinResultModifier_NoComponentValues(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w_bad",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	// JSON without componentValues or data
	js := []byte(`{"foo": {"bar": 1}}`)
	node, err := sonic.Get(js)
	assert.NoError(t, err)

	_, err = parseWinResultModifier(betCfg, &node)
	assert.Error(t, err)
}

func TestParseWinResultModifier_UnmarshalError(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w_bad2",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	// configuration is a string instead of object -> sonic.Unmarshal should error
	js := []byte(`{
		"componentValues": {
			"label": "w_bad2",
			"configuration": "not-an-object"
		}
	}`)
	node, err := sonic.Get(js)
	assert.NoError(t, err)

	_, err = parseWinResultModifier(betCfg, &node)
	assert.Error(t, err)
}

func TestParseWinResultModifier_DataNode(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w2",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"data": {
			"label": "w2",
			"configuration": {
				"type": "Divide",
				"sourceComponent": ["src"],
				"winDivisor": 2
			}
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseWinResultModifier(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "w2", name)

	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	cfg := cfgIface.(*WinResultModifierConfig)
	assert.Equal(t, "divide", cfg.StrType)
	assert.Equal(t, 2, cfg.WinDivisor)
}
func (f *fakeComp) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	return nil
}
func (f *fakeComp) NewComponentData() IComponentData { return &BasicComponentData{} }
func (f *fakeComp) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}
func (f *fakeComp) ProcRespinOnStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, canRemove bool) (string, error) {
	return "", nil
}
func (f *fakeComp) GetName() string                         { return f.Name }
func (f *fakeComp) IsRespin() bool                          { return false }
func (f *fakeComp) IsForeach() bool                         { return false }
func (f *fakeComp) NewStats2(parent string) *stats2.Feature { return nil }
func (f *fakeComp) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
}
func (f *fakeComp) IsNeedOnStepEndStats2() bool      { return false }
func (f *fakeComp) GetAllLinkComponents() []string   { return nil }
func (f *fakeComp) GetNextLinkComponents() []string  { return nil }
func (f *fakeComp) GetChildLinkComponents() []string { return nil }
func (f *fakeComp) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) {
	return false, nil
}
func (f *fakeComp) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
}
func (f *fakeComp) IsMask() bool { return false }
func (f *fakeComp) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error {
	return nil
}
func (f *fakeComp) SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, mask bool) error {
	return nil
}
func (f *fakeComp) SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error {
	return nil
}
func (f *fakeComp) EachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake,
	prs []*sgc7game.PlayResult, cd IComponentData) error {
	return nil
}
func (f *fakeComp) AddPos(cd IComponentData, x int, y int) {}
func (f *fakeComp) OnPlayGameWithSet(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData, set int) (string, error) {
	return "", nil
}
func (f *fakeComp) ClearData(icd IComponentData, bForceNow bool) {}
func (f *fakeComp) InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, ps *PlayerState, betMethod int, bet int) error {
	return nil
}
func (f *fakeComp) NewPlayerState() IComponentPS { return nil }
func (f *fakeComp) OnUpdateDataWithPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, ps *PlayerState, betMethod int, bet int, cd IComponentData) {
}
func (f *fakeComp) ChgReelsCollector(icd IComponentData, ps *PlayerState, betMethod int, bet int, reelsData []int) {
}

func TestWinResultModifierDataBasics(t *testing.T) {
	d := &WinResultModifierData{}
	// onNewStep sets defaults
	d.Wins = 5
	d.WinMulti = 7
	d.onNewStep()
	assert.Equal(t, 0, d.Wins)
	assert.Equal(t, 1, d.WinMulti)

	// Clone copies values
	d.Wins = 3
	d.WinMulti = 4
	d.BasicComponentData.MapConfigIntVals = map[string]int{"x": 1}
	cd := d.Clone().(*WinResultModifierData)
	assert.Equal(t, 3, cd.Wins)
	assert.Equal(t, 4, cd.WinMulti)

	// BuildPBComponentData and GetValEx
	pb := d.BuildPBComponentData()
	assert.NotNil(t, pb)

	v, ok := d.GetValEx(CVWins, 0)
	assert.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestWinResultModifierConfigSetLink(t *testing.T) {
	cfg := &WinResultModifierConfig{}
	cfg.SetLinkComponent("next", "comp1")
	assert.Equal(t, "comp1", cfg.DefaultNextComponent)
}

func makePoolWithPaytables() *GamePropertyPool {
	return &GamePropertyPool{
		DefaultPaytables: &sgc7game.PayTables{MapSymbols: map[string]int{"A": 1}},
	}
}

func TestInitExNormalizeAndTargetSymbols(t *testing.T) {
	wrm := NewWinResultModifier("w1").(*WinResultModifier)

	cfg := &WinResultModifierConfig{
		BasicComponentConfig: BasicComponentConfig{},
		StrType:              "multiply",
		WinMulti:             0,
		WinDivisor:           0,
		TargetSymbols:        []string{"A"},
	}

	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	// normalized
	assert.Equal(t, 1, wrm.Config.WinMulti)
	assert.Equal(t, 1, wrm.Config.WinDivisor)
	// target symbol codes filled
	assert.Equal(t, []int{1}, wrm.Config.TargetSymbolCodes)
}

func TestNewComponentDataAndAsciiAndNew(t *testing.T) {
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cd := wrm.NewComponentData()
	assert.IsType(t, &WinResultModifierData{}, cd)

	// OnAsciiGame should not error
	d := cd.(*WinResultModifierData)
	d.WinMulti = 2
	d.Wins = 10
	err := wrm.OnAsciiGame(nil, nil, nil, nil, d)
	assert.Nil(t, err)

	// NewWinResultModifier basic
	nc := NewWinResultModifier("xx")
	assert.NotNil(t, nc)
}

func TestOnNewGameDataAndJsonBuild(t *testing.T) {
	// OnNewGame should initialize maps
	d := &WinResultModifierData{}
	d.MapConfigIntVals = nil
	d.MapConfigVals = nil
	d.OnNewGame(nil, nil)
	if d.MapConfigIntVals == nil {
		t.Fatalf("MapConfigIntVals should be initialized")
	}

	// json builder
	jwt := &jsonWinResultModifier{Type: "DIVIDE", SourceComponents: []string{"src"}, WinMulti: 5, TargetSymbols: []string{"A"}, WinDivisor: 10}
	cfg := jwt.build()
	assert.Equal(t, "divide", cfg.StrType)
	assert.Equal(t, 5, cfg.WinMulti)
	assert.Equal(t, []string{"A"}, cfg.TargetSymbols)
	assert.Equal(t, 10, cfg.WinDivisor)
}

func TestOnPlayGameSymbolBranches(t *testing.T) {
	// prepare common gp/gparam/pr
	gp, gparam, pr := setupGameForPlay(t)

	// create a scene with symbol code 1 at (0,0)
	gs, _ := sgc7game.NewGameScene(1, 1)
	gs.Arr[0][0] = 1

	// push scene into stack so GetTargetScene3 finds it
	gp.SceneStack.Push("src", gs)

	// ExistSymbol
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cfg := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "existsymbol", SourceComponents: []string{"src"}, WinMulti: 2, TargetSymbols: []string{"A"}}
	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	icd := &WinResultModifierData{}
	_, err = wrm.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Nil(t, err)
	assert.Equal(t, 200, pr.Results[0].CoinWin)

	// reset coinwin
	pr.Results[0].CoinWin = 100

	// AddSymbolMulti (winMulti * num)
	wrm2 := NewWinResultModifier("w2").(*WinResultModifier)
	cfg2 := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "addsymbolmulti", SourceComponents: []string{"src"}, WinMulti: 3, TargetSymbols: []string{"A"}}
	err = wrm2.InitEx(cfg2, makePoolWithPaytables())
	assert.Nil(t, err)
	icd2 := &WinResultModifierData{}
	_, err = wrm2.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd2)
	assert.Nil(t, err)
	assert.Equal(t, 300, pr.Results[0].CoinWin)

	// reset coinwin
	pr.Results[0].CoinWin = 100

	// MulSymbolMulti (winMulti^num)
	wrm3 := NewWinResultModifier("w3").(*WinResultModifier)
	cfg3 := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "mulsymbolmulti", SourceComponents: []string{"src"}, WinMulti: 4, TargetSymbols: []string{"A"}}
	err = wrm3.InitEx(cfg3, makePoolWithPaytables())
	assert.Nil(t, err)
	icd3 := &WinResultModifierData{}
	_, err = wrm3.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd3)
	assert.Nil(t, err)
	assert.Equal(t, 400, pr.Results[0].CoinWin)

	// reset coinwin
	pr.Results[0].CoinWin = 100

	// SymbolMultiOnWays
	wrm4 := NewWinResultModifier("w4").(*WinResultModifier)
	cfg4 := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "symbolmultionways", SourceComponents: []string{"src"}, WinMulti: 5, TargetSymbols: []string{"A"}}
	err = wrm4.InitEx(cfg4, makePoolWithPaytables())
	assert.Nil(t, err)
	icd4 := &WinResultModifierData{}
	_, err = wrm4.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd4)
	assert.Nil(t, err)
	// mul should be winMulti (5)
	assert.Equal(t, 500, pr.Results[0].CoinWin)
}

func TestGetValExAndGetWinMultiOverride(t *testing.T) {
	d := &WinResultModifierData{}
	// non CVWins should return false
	v, ok := d.GetValEx("something", 0)
	assert.False(t, ok)
	assert.Equal(t, 0, v)

	// GetWinMulti with BasicComponentData override
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	basic := &BasicComponentData{MapConfigIntVals: map[string]int{CCVWinMulti: 7}}
	val := wrm.GetWinMulti(basic)
	assert.Equal(t, 7, val)
}

func TestInitReadsYAMLFile(t *testing.T) {
	// create temp yaml
	content := "type: divide\nsourceComponents: [\"src\"]\nwinDivisor: 5\ntargetSymbols: [\"A\"]\n"
	fn := t.TempDir() + "/wrm.yaml"
	err := os.WriteFile(fn, []byte(content), 0644)
	assert.Nil(t, err)

	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	err = wrm.Init(fn, makePoolWithPaytables())
	assert.Nil(t, err)
	assert.Equal(t, WRMTypeDivide, wrm.Config.Type)
	assert.Equal(t, 5, wrm.Config.WinDivisor)
}

func TestInitErrorPaths(t *testing.T) {
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	// non-existent file
	err := wrm.Init("/non/existent/file.yaml", makePoolWithPaytables())
	if err == nil {
		t.Fatalf("expected error for non-existent file")
	}

	// invalid yaml (malformed)
	fn := t.TempDir() + "/bad.yaml"
	// missing closing bracket -> should cause unmarshal error
	err = os.WriteFile(fn, []byte("type: [1,2"), 0644)
	assert.Nil(t, err)
	err = wrm.Init(fn, makePoolWithPaytables())
	if err == nil {
		t.Fatalf("expected unmarshal error for bad yaml")
	}
}

func setupGameForPlay(t *testing.T) (*GameProperty, *GameParams, *sgc7game.PlayResult) {
	// minimal gameProp
	gp := &GameProperty{}
	gp.Pool = makePoolWithPaytables()
	gp.Components = NewComponentList()
	gp.SceneStack = NewSceneStack(false)
	gp.OtherSceneStack = NewSceneStack(true)
	// add a fake component named "src"
	fc := &fakeComp{BasicComponent: NewBasicComponent("src", 0)}
	gp.Components.AddComponent("src", fc)

	// init callstack and onNewGame so component data can be created
	gp.callStack = NewCallStack()
	gp.callStack.OnNewGame()

	// gameParam with history contains src
	gparam := NewGameParam(&sgc7game.Stake{CashBet: 1, CoinBet: 1}, nil)
	gparam.HistoryComponents = []string{"src"}

	// playresult with one result
	pr := &sgc7game.PlayResult{}
	pr.Results = []*sgc7game.Result{{CoinWin: 100, CashWin: 100, Mul: 1, Pos: []int{0, 0}}}

	// ensure component data exists and mark used result index
	ccd := gp.GetComponentDataWithName("src").(*BasicComponentData)
	ccd.UsedResults = []int{0}

	return gp, gparam, pr
}

func TestOnPlayGameEarlyDoNothing(t *testing.T) {
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cfg := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "addsymbolmulti", SourceComponents: []string{"src"}, WinMulti: 1}
	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	gp, gparam, pr := setupGameForPlay(t)

	icd := &WinResultModifierData{}

	// winMulti == 1 and type needs multiply -> should return ErrComponentDoNothing
	nc, err := wrm.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)
}

func TestOnPlayGameDivideAndMultiply(t *testing.T) {
	// DIVIDE
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cfg := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "divide", SourceComponents: []string{"src"}, WinDivisor: 10}
	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	gp, gparam, pr := setupGameForPlay(t)

	icd := &WinResultModifierData{}

	nc, err := wrm.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Nil(t, err)
	assert.Equal(t, "", nc)
	assert.Equal(t, 10, pr.Results[0].CoinWin)

	// MULTIPLY
	wrm2 := NewWinResultModifier("w2").(*WinResultModifier)
	cfg2 := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "multiply", SourceComponents: []string{"src"}, WinMulti: 3}
	err2 := wrm2.InitEx(cfg2, makePoolWithPaytables())
	assert.Nil(t, err2)

	gp2, gparam2, pr2 := setupGameForPlay(t)
	icd2 := &WinResultModifierData{}

	nc2, err2 := wrm2.OnPlayGame(gp2, pr2, gparam2, nil, "", "", nil, nil, nil, icd2)
	assert.Nil(t, err2)
	assert.Equal(t, "", nc2)
	assert.Equal(t, 300, pr2.Results[0].CoinWin)
}

func TestOnPlayGameNoScene(t *testing.T) {
	// When type needs game scene and no scene exists, should return ErrComponentDoNothing
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cfg := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "existsymbol", SourceComponents: []string{"src"}, WinMulti: 2, TargetSymbols: []string{"A"}}
	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	// setup game but do NOT push any scene
	gp, gparam, pr := setupGameForPlay(t)

	icd := &WinResultModifierData{}
	nc, err := wrm.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)
}

func TestOnPlayGameNoHistory(t *testing.T) {
	// When source components are not in history, component should do nothing
	wrm := NewWinResultModifier("w1").(*WinResultModifier)
	cfg := &WinResultModifierConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "multiply", SourceComponents: []string{"src"}, WinMulti: 2}
	err := wrm.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	gp, gparam, pr := setupGameForPlay(t)
	// clear history so source component is not marked as executed
	gparam.HistoryComponents = []string{}

	icd := &WinResultModifierData{}
	nc, err := wrm.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)
}
