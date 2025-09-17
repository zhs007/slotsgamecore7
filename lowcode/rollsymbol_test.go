package lowcode

import (
    "context"
    "fmt"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
    sgc7utils "github.com/zhs007/slotsgamecore7/utils"
    "github.com/zhs007/slotsgamecore7/plugin"
    "github.com/zhs007/slotsgamecore7/asciigame"
    "github.com/zhs007/slotsgamecore7/stats2"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/anypb"
    "github.com/bytedance/sonic"
    "github.com/bytedance/sonic/ast"
)

// fakePlugin deterministic plugin for RandVal
type fakePluginRS struct{}

func (p *fakePluginRS) Random(_ context.Context, r int) (int, error) {
    if r <= 0 {
        return 0, nil
    }

    // always return zero so RandVal picks first index deterministically
    return 0, nil
}
func (p *fakePluginRS) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginRS) ClearUsedRngs()                    {}
func (p *fakePluginRS) TagUsedRngs()                      {}
func (p *fakePluginRS) RollbackUsedRngs() error           { return nil }
func (p *fakePluginRS) SetCache(arr []int)                {}
func (p *fakePluginRS) ClearCache()                       {}
func (p *fakePluginRS) Init()                             {}
func (p *fakePluginRS) SetScenePool(any)                  {}
func (p *fakePluginRS) GetScenePool() any                 { return nil }
func (p *fakePluginRS) SetSeed(seed int)                  {}

func TestRollSymbol_InitEx_ErrNoWeight(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    // empty pool
    pool := &GamePropertyPool{}

    // cfg without Weight should return ErrInvalidComponentConfig
    cfg := &RollSymbolConfig{SymbolNum: 1}

    err := comp.InitEx(cfg, pool)
    assert.Error(t, err)
}

func TestRollSymbol_InitEx_SuccessAndOnPlay(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)

    // prepare ValWeights2 with two int vals [10,20]
    iv1 := sgc7game.NewIntValEx(10)
    iv2 := sgc7game.NewIntValEx(20)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv1, iv2}, []int{1, 1})
    assert.NoError(t, err)

    pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"w1": vw}, DefaultPaytables: &sgc7game.PayTables{}}

    cfg := &RollSymbolConfig{Weight: "w1", SymbolNum: 2}

    err = comp.InitEx(cfg, pool)
    assert.NoError(t, err)

    // prepare gameProp manually to avoid using pool.newGameProp which requires initializers
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    // minimal stubs
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    // component data
    cd := comp.NewComponentData().(*RollSymbolData)

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    gpar := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

    // run OnPlayGame with fake plugin
    nc, err := comp.OnPlayGame(gp, pr, gpar, &fakePluginRS{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    assert.NoError(t, err)
    // since DefaultNextComponent empty, nc == ""
    assert.Equal(t, "", nc)

    // symbol codes should have been generated (2 entries)
    assert.Equal(t, 2, len(cd.SymbolCodes))

    // BuildPBComponentData and Clone cover methods
    pb := cd.BuildPBComponentData()
    assert.NotNil(t, pb)
    clone := cd.Clone().(*RollSymbolData)
    assert.Equal(t, cd.SymbolCodes, clone.SymbolCodes)
}

// reuse stubRNG and stubFeatureLevel defined in other test files

func TestRollSymbol_GetSymbolNum_Overrides(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)

    // minimal config
    comp.Config = &RollSymbolConfig{Weight: "", SymbolNum: 5}

    // basic component data with CCVSymbolNum override
    bcd := &BasicComponentData{}
    bcd.MapConfigIntVals = map[string]int{CCVSymbolNum: 3}

    gp := &GameProperty{Pool: &GamePropertyPool{}}

    // should return 3 from config val
    got := comp.getSymbolNum(gp, bcd)
    assert.Equal(t, 3, got)

    // test default SymbolNum when no overrides (clear previous config int)
    bcd.MapConfigIntVals = map[string]int{}
    comp.Config.SymbolNum = 7
    got2 := comp.getSymbolNum(gp, bcd)
    assert.Equal(t, 7, got2)
}

func TestRollSymbol_OnAsciiGame_NoPanic(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    // Create simple pool and gameProp so paytables lookup won't panic
    pool := &GamePropertyPool{DefaultPaytables: &sgc7game.PayTables{}}
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    cd := comp.NewComponentData().(*RollSymbolData)
    cd.SymbolCodes = []int{1, 2}

    err := comp.OnAsciiGame(gp, sgc7game.NewPlayResult("m", 0, 0, "t"), nil, nil, cd)
    assert.NoError(t, err)
}

// helper component/data to provide symbol lists
type cbSymbolsCD struct{
    symbols []int
}
func (c *cbSymbolsCD) OnNewGame(gameProp *GameProperty, component IComponent) {}
func (c *cbSymbolsCD) BuildPBComponentData() proto.Message { return nil }
func (c *cbSymbolsCD) Clone() IComponentData { return &cbSymbolsCD{symbols: append([]int{}, c.symbols...)} }
func (c *cbSymbolsCD) GetValEx(key string, getType GetComponentValType) (int, bool) { return 0, false }
func (c *cbSymbolsCD) GetStrVal(key string) (string, bool) { return "", false }
func (c *cbSymbolsCD) GetConfigVal(key string) string { return "" }
func (c *cbSymbolsCD) SetConfigVal(key string, val string) {}
func (c *cbSymbolsCD) GetConfigIntVal(key string) (int, bool) { return 0, false }
func (c *cbSymbolsCD) SetConfigIntVal(key string, val int) {}
func (c *cbSymbolsCD) ChgConfigIntVal(key string, off int) int { return 0 }
func (c *cbSymbolsCD) ClearConfigIntVal(key string) {}
func (c *cbSymbolsCD) GetResults() []int { return nil }
func (c *cbSymbolsCD) GetOutput() int { return 0 }
func (c *cbSymbolsCD) GetStringOutput() string { return "" }
func (c *cbSymbolsCD) GetSymbols() []int { return append([]int{}, c.symbols...) }
func (c *cbSymbolsCD) AddSymbol(symbolCode int) { c.symbols = append(c.symbols, symbolCode) }
func (c *cbSymbolsCD) GetPos() []int { return nil }
func (c *cbSymbolsCD) HasPos(x int, y int) bool { return false }
func (c *cbSymbolsCD) AddPos(x int, y int) {}
func (c *cbSymbolsCD) ClearPos() {}
func (c *cbSymbolsCD) GetLastRespinNum() int { return 0 }
func (c *cbSymbolsCD) GetCurRespinNum() int { return 0 }
func (c *cbSymbolsCD) IsRespinEnding() bool { return false }
func (c *cbSymbolsCD) IsRespinStarted() bool { return false }
func (c *cbSymbolsCD) AddTriggerRespinAward(award *Award) {}
func (c *cbSymbolsCD) AddRespinTimes(num int) {}
func (c *cbSymbolsCD) TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {}
func (c *cbSymbolsCD) PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {}
func (c *cbSymbolsCD) GetMask() []bool { return nil }
func (c *cbSymbolsCD) ChgMask(curMask int, val bool) bool { return false }
func (c *cbSymbolsCD) PutInMoney(coins int) {}
func (c *cbSymbolsCD) ChgReelsCollector(reelsData []int) {}
func (c *cbSymbolsCD) SetSymbolCodes(symbolCodes []int) {}
func (c *cbSymbolsCD) GetSymbolCodes() []int { return nil }

type cbSymbolsComp struct{ name string; syms []int }
func (c *cbSymbolsComp) Init(fn string, pool *GamePropertyPool) error { return nil }
func (c *cbSymbolsComp) InitEx(cfg any, pool *GamePropertyPool) error { return nil }
func (c *cbSymbolsComp) OnGameInited(components *ComponentList) error { return nil }
func (c *cbSymbolsComp) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
    cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) { return "", nil }
func (c *cbSymbolsComp) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error { return nil }
func (c *cbSymbolsComp) NewComponentData() IComponentData { return &cbSymbolsCD{symbols: append([]int{}, c.syms...)} }
func (c *cbSymbolsComp) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {}
func (c *cbSymbolsComp) ProcRespinOnStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, canRemove bool) (string, error) { return "", nil }
func (c *cbSymbolsComp) GetName() string { return c.name }
func (c *cbSymbolsComp) IsRespin() bool { return false }
func (c *cbSymbolsComp) IsForeach() bool { return false }
func (c *cbSymbolsComp) NewStats2(parent string) *stats2.Feature { return nil }
func (c *cbSymbolsComp) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {}
func (c *cbSymbolsComp) IsNeedOnStepEndStats2() bool { return false }
func (c *cbSymbolsComp) GetAllLinkComponents() []string { return nil }
func (c *cbSymbolsComp) GetNextLinkComponents() []string { return nil }
func (c *cbSymbolsComp) GetChildLinkComponents() []string { return nil }
func (c *cbSymbolsComp) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) { return false, nil }
func (c *cbSymbolsComp) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {}
func (c *cbSymbolsComp) IsMask() bool { return false }
func (c *cbSymbolsComp) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error { return nil }
func (c *cbSymbolsComp) SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, mask bool) error { return nil }
func (c *cbSymbolsComp) SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error { return nil }
func (c *cbSymbolsComp) EachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) error { return nil }
func (c *cbSymbolsComp) AddPos(cd IComponentData, x int, y int) {}
func (c *cbSymbolsComp) OnPlayGameWithSet(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
    cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData, set int) (string, error) { return "", nil }
func (c *cbSymbolsComp) ClearData(icd IComponentData, bForceNow bool) {}
func (c *cbSymbolsComp) InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, ps *PlayerState, betMethod int, bet int) error { return nil }
func (c *cbSymbolsComp) NewPlayerState() IComponentPS { return nil }
func (c *cbSymbolsComp) OnUpdateDataWithPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, ps *PlayerState, betMethod int, bet int, cd IComponentData) {}
func (c *cbSymbolsComp) ChgReelsCollector(icd IComponentData, ps *PlayerState, betMethod int, bet int, reelsData []int) {}

func TestRollSymbol_GetValWeight_FilterAndTargetCollection(t *testing.T) {
    // prepare ValWeights2 with int vals [1,2,3]
    iv1 := sgc7game.NewIntValEx(1)
    iv2 := sgc7game.NewIntValEx(2)
    iv3 := sgc7game.NewIntValEx(3)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv1, iv2, iv3}, []int{1, 1, 1})
    assert.NoError(t, err)

    pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"w1": vw}, DefaultPaytables: &sgc7game.PayTables{}}

    comp := NewRollSymbol("rs").(*RollSymbol)
    cfg := &RollSymbolConfig{Weight: "w1", SymbolNum: 3}
    err = comp.InitEx(cfg, pool)
    assert.NoError(t, err)

    // gameProp with a source component that reports only symbol 2
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.callStack = NewCallStack()
    gp.callStack = NewCallStack()

    // register source component
    gp.Components = NewComponentList()
    gp.Components.MapComponents = map[string]IComponent{"src": &cbSymbolsComp{name: "src", syms: []int{2}}}
    gp.callStack.OnNewGame()
    // inject component data into global callstack node so GetComponentData can find it
    gp.callStack.nodes[0].MapComponentData["src"] = &cbSymbolsCD{symbols: []int{2}}

    comp.Config.SrcSymbolCollection = "src"
    vw2 := comp.getValWeight(gp)
    // after filtering, only one value should remain and it should be 2
    assert.NotNil(t, vw2)
    assert.Equal(t, 1, len(vw2.Vals))
    assert.Equal(t, 2, vw2.Vals[0].Int())

    // test ignore collection removes a symbol
    gp.Components.MapComponents["ig"] = &cbSymbolsComp{name: "ig", syms: []int{2}}
    gp.callStack.nodes[0].MapComponentData["ig"] = &cbSymbolsCD{symbols: []int{2}}
    comp.Config.SrcSymbolCollection = ""
    comp.Config.IgnoreSymbolCollection = "ig"
    vw3 := comp.getValWeight(gp)
    assert.NotNil(t, vw3)
    // now values should be 1 and 3
    vals := []int{vw3.Vals[0].Int(), vw3.Vals[1].Int()}
    assert.Contains(t, vals, 1)
    assert.Contains(t, vals, 3)

    // test TargetSymbolCollection: create target component to collect symbols
    gp.Components.MapComponents["tgt"] = &cbSymbolsComp{name: "tgt", syms: []int{}}
    gp.callStack.nodes[0].MapComponentData["tgt"] = &cbSymbolsCD{symbols: []int{}}
    comp.Config.TargetSymbolCollection = "tgt"

    cd := comp.NewComponentData().(*RollSymbolData)
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    gpar := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
    // run OnPlayGame; it should append symbols to target component's data
    _, err = comp.OnPlayGame(gp, pr, gpar, &fakePluginRS{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    assert.NoError(t, err)
    syms := gp.GetComponentSymbols("tgt")
    assert.NotNil(t, syms)
    assert.GreaterOrEqual(t, len(syms), 1)
}

func TestRollSymbol_OnPlayGame_NoWeight_DoNothing(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    // set config with weight pointing to missing weight in pool
    comp.Config = &RollSymbolConfig{Weight: "missing", SymbolNum: 2}

    // ensure embedded BasicComponent has a config to avoid nil deref in onStepEnd
    comp.BasicComponent.Config = &BasicComponentConfig{DefaultNextComponent: ""}

    gp := &GameProperty{Pool: &GamePropertyPool{}}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.callStack = NewCallStack()

    cd := comp.NewComponentData().(*RollSymbolData)
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    gpar := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

    nc, err := comp.OnPlayGame(gp, pr, gpar, &fakePluginRS{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    // when no weights available, component should report ErrComponentDoNothing
    assert.Equal(t, "", nc)
    assert.Error(t, err)
}

// test Init reading from a YAML file to cover Init path
func TestRollSymbol_Init_File(t *testing.T) {
    tmpf := "test_rollsymbol_init.yaml"
    yaml := "weight: w1\nsymbolNum: 2\n"
    err := os.WriteFile(tmpf, []byte(yaml), 0644)
    if err != nil { t.Fatalf("write temp file err=%v", err) }
    defer os.Remove(tmpf)

    // prepare pool with weight named w1
    iv := sgc7game.NewIntValEx(1)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv}, []int{1})
    assert.NoError(t, err)
    pool := &GamePropertyPool{mapIntValWeights: map[string]*sgc7game.ValWeights2{"w1": vw}, DefaultPaytables: &sgc7game.PayTables{}}

    comp := NewRollSymbol("rfi").(*RollSymbol)
    err = comp.Init(tmpf, pool)
    assert.NoError(t, err)
}

// extra tests to increase coverage on small branches
func TestRollSymbol_SetLinkComponentAndGetValWeightNoFilter(t *testing.T) {
    // SetLinkComponent
    cfg := &RollSymbolConfig{}
    cfg.SetLinkComponent("next", "nxt")
    if cfg.DefaultNextComponent != "nxt" {
        t.Fatalf("SetLinkComponent failed")
    }

    // getValWeight without filters should return the original WeightVW pointer
    comp := NewRollSymbol("rs").(*RollSymbol)
    iv := sgc7game.NewIntValEx(5)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv}, []int{1})
    assert.NoError(t, err)

    comp.Config = &RollSymbolConfig{WeightVW: vw}

    got := comp.getValWeight(nil)
    if got != vw {
        t.Fatalf("getValWeight expected same pointer")
    }
}

// outCD provides GetOutput for testing SymbolNumComponent override
type outCD struct{ BasicComponentData }
func (o *outCD) GetOutput() int { return 9 }

func TestRollSymbol_GetSymbolNum_ComponentOverride(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    comp.Config = &RollSymbolConfig{SymbolNum: 4, SymbolNumComponent: "other"}

    gp := &GameProperty{Pool: &GamePropertyPool{}}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    gp.Components = NewComponentList()
    gp.Components.MapComponents = map[string]IComponent{"other": &cbSymbolsComp{name: "other"}}
    gp.callStack.nodes[0].MapComponentData["other"] = &outCD{}

    bcd := &BasicComponentData{}

    got := comp.getSymbolNum(gp, bcd)
    assert.Equal(t, 9, got)
}

// plugin that returns error from Random to exercise RandVal error branch
type fakePluginErrLocal struct{}
func (p *fakePluginErrLocal) Random(_ context.Context, r int) (int, error) { return 0, fmt.Errorf("rnderr") }
func (p *fakePluginErrLocal) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginErrLocal) ClearUsedRngs()                    {}
func (p *fakePluginErrLocal) TagUsedRngs()                      {}
func (p *fakePluginErrLocal) RollbackUsedRngs() error           { return nil }
func (p *fakePluginErrLocal) SetCache(arr []int)                {}
func (p *fakePluginErrLocal) ClearCache()                       {}
func (p *fakePluginErrLocal) Init()                             {}
func (p *fakePluginErrLocal) SetScenePool(any)                  {}
func (p *fakePluginErrLocal) GetScenePool() any                 { return nil }
func (p *fakePluginErrLocal) SetSeed(seed int)                  {}

func TestRollSymbol_OnPlayGame_RandValError(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)

    // prepare ValWeights2 with two values so RandVal will call plugin.Random
    iv1 := sgc7game.NewIntValEx(1)
    iv2 := sgc7game.NewIntValEx(2)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv1, iv2}, []int{1, 1})
    assert.NoError(t, err)

    comp.Config = &RollSymbolConfig{WeightVW: vw, SymbolNum: 1}
    comp.BasicComponent.Config = &BasicComponentConfig{DefaultNextComponent: ""}

    gp := &GameProperty{Pool: &GamePropertyPool{}}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}

    cd := comp.NewComponentData().(*RollSymbolData)
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    gpar := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

    _, err = comp.OnPlayGame(gp, pr, gpar, &fakePluginErrLocal{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    assert.Error(t, err)
}

func TestParseRollSymbol_FromAST(t *testing.T) {
    jsonStr := `{"componentValues":{"label":"rsl","configuration":{"weight":"w1","symbolNum":2}}}`

    var node ast.Node
    err := sonic.Unmarshal([]byte(jsonStr), &node)
    assert.NoError(t, err)

    bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
    label, err := parseRollSymbol(bc, &node)
    assert.NoError(t, err)
    assert.Equal(t, "rsl", label)
}

func TestRollSymbolData_OnNewGame_Cover(t *testing.T) {
    d := &RollSymbolData{}
    // should not panic
    d.OnNewGame(&GameProperty{}, NewRollSymbol("r"))
}

func TestRollSymbol_GetValWeight_FiltersEmptyResult(t *testing.T) {
    // weight vals [1,2], but src collection only has 3 -> CloneWithIntArray returns nil
    iv1 := sgc7game.NewIntValEx(1)
    iv2 := sgc7game.NewIntValEx(2)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv1, iv2}, []int{1, 1})
    assert.NoError(t, err)

    comp := NewRollSymbol("rs").(*RollSymbol)
    comp.Config = &RollSymbolConfig{WeightVW: vw, SrcSymbolCollection: "src"}

    gp := &GameProperty{Pool: &GamePropertyPool{}}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    gp.Components = NewComponentList()
    gp.Components.MapComponents = map[string]IComponent{"src": &cbSymbolsComp{name: "src", syms: []int{3}}}
    gp.callStack.nodes[0].MapComponentData["src"] = &cbSymbolsCD{symbols: []int{3}}

    got := comp.getValWeight(gp)
    assert.Nil(t, got)
}

func TestRollSymbol_GetValWeight_IgnoreRemovesAll(t *testing.T) {
    // weight vals [1], ignore collection has [1] -> CloneWithoutIntArray returns nil
    iv1 := sgc7game.NewIntValEx(1)
    vw, err := sgc7game.NewValWeights2([]sgc7game.IVal{iv1}, []int{1})
    assert.NoError(t, err)

    comp := NewRollSymbol("rs").(*RollSymbol)
    comp.Config = &RollSymbolConfig{WeightVW: vw, IgnoreSymbolCollection: "ig"}

    gp := &GameProperty{Pool: &GamePropertyPool{}}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    gp.Components = NewComponentList()
    gp.Components.MapComponents = map[string]IComponent{"ig": &cbSymbolsComp{name: "ig", syms: []int{1}}}
    gp.callStack.nodes[0].MapComponentData["ig"] = &cbSymbolsCD{symbols: []int{1}}

    got := comp.getValWeight(gp)
    assert.Nil(t, got)
}

func TestRollSymbol_InitEx_WrongType(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    pool := &GamePropertyPool{}

    // pass wrong type
    err := comp.InitEx("not-a-config", pool)
    assert.Error(t, err)
}

func TestRollSymbol_OnPlayGame_InvalidICD(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    gp := &GameProperty{Pool: &GamePropertyPool{}}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    gpar := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

    // pass wrong icd type
    nc, err := comp.OnPlayGame(gp, pr, gpar, &fakePluginRS{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, &cbSymbolsCD{})
    assert.Equal(t, "", nc)
    assert.Error(t, err)
}

func TestRollSymbol_OnAsciiGame_InvalidICD(t *testing.T) {
    comp := NewRollSymbol("rs").(*RollSymbol)
    pool := &GamePropertyPool{DefaultPaytables: &sgc7game.PayTables{}}
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    err := comp.OnAsciiGame(gp, sgc7game.NewPlayResult("m", 0, 0, "t"), nil, nil, &cbSymbolsCD{})
    assert.Error(t, err)
}

func TestJSONRollSymbol_Build(t *testing.T) {
    jr := &jsonRollSymbol{Weight: "w", SymbolNum: 3, SymbolNumComponent: "c", SrcSymbolCollection: "s", IgnoreSymbolCollection: "ig", TargetSymbolCollection: "tgt"}
    cfg := jr.build()
    assert.Equal(t, "w", cfg.Weight)
    assert.Equal(t, 3, cfg.SymbolNum)
    assert.Equal(t, "c", cfg.SymbolNumComponent)
    assert.Equal(t, "s", cfg.SrcSymbolCollection)
}
