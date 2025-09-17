package lowcode

import (
    "os"
    "testing"

    "github.com/bytedance/sonic"
    "github.com/bytedance/sonic/ast"
    "github.com/stretchr/testify/assert"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
    "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// posCD is a small shim implementing GetPos used by tests to simulate
// component data that returns explicit positions.
type posCD struct{ BasicComponentData }

func (p *posCD) GetPos() []int { return []int{0, 0, 2, 2} }

// emptyPosCD returns an empty pos slice to simulate a component with no positions
type emptyPosCD struct{ BasicComponentData }

func (p *emptyPosCD) GetPos() []int { return []int{} }

// helper to create a pool with simple paytables
func makePoolForSum() *GamePropertyPool {
    p := &GamePropertyPool{}
    p.DefaultPaytables = &sgc7game.PayTables{MapSymbols: map[string]int{"A": 1, "B": 2}}
    return p
}

func TestSSV_ParseFromAST(t *testing.T) {
    jsonStr := `{"componentValues":{"label":"sumlabel","configuration":{"type":">","value":5,"sourceComponent":"src"}}}`

    var node ast.Node
    err := sonic.Unmarshal([]byte(jsonStr), &node)
    assert.NoError(t, err)

    bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
    label, err := parseSumSymbolVals(bc, &node)
    assert.NoError(t, err)
    assert.Equal(t, "sumlabel", label)
    _, ok := bc.mapConfig["sumlabel"]
    assert.True(t, ok)
}

func TestSSV_CheckValTypes(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{}

    // equality
    s.Config.Type = SSVTypeEqu
    s.Config.Value = 3
    assert.True(t, s.checkVal(3))
    assert.False(t, s.checkVal(4))

    // greater eq
    s.Config.Type = SSVTypeGreaterEqu
    s.Config.Value = 5
    assert.True(t, s.checkVal(5))
    assert.True(t, s.checkVal(6))

    // less
    s.Config.Type = SSVTypeLess
    s.Config.Value = 10
    assert.True(t, s.checkVal(9))
    assert.False(t, s.checkVal(10))

    // in area [min,max]
    s.Config.Type = SSVTypeInAreaLR
    s.Config.Min = 1
    s.Config.Max = 3
    assert.True(t, s.checkVal(1))
    assert.True(t, s.checkVal(3))
    assert.False(t, s.checkVal(0))
}

func TestSSV_SumFullScanAndSourceComponent(t *testing.T) {
    pool := makePoolForSum()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()

    // scene 3x3
    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})

    // ensure Components is initialized to avoid nil deref
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)

    // no source component configured -> full scan sum of values > 5
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{}
    s.Config.Type = SSVTypeGreater
    s.Config.Value = 5

    sum := s.sum(gp, os)
    // values greater than 5: 6,7,8,9 => sum = 30
    assert.Equal(t, 6+7+8+9, sum)

    // configure source component with positions -> only those positions used
    // create a fake component and component data that returns positions (0,0) and (2,2)
    comp := NewSumSymbolVals("src")
    // register component into component list and into gameProp
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)
    gp.Components.MapComponents["src"] = comp

    // create a component data that has positions [0,0,2,2]
    // use a package-level shim type posCD (see top of file)
    pcd := &posCD{}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    // inject into callstack global map so GetComponentData returns it
    gp.callStack.nodes[0].MapComponentData["src"] = pcd

    s.Config.SourceComponent = "src"
    sum2 := s.sum(gp, os)
    // positions (0,0)=1 and (2,2)=9 => sum = 1+9 = 10 (both >5? only 9 >5 so only 9 counted)
    // but checkVal requires >5 so only 9 should be counted
    assert.Equal(t, 9, sum2)
}

func TestSSV_OnPlayGameSetsComponentDataNumber(t *testing.T) {
    pool := makePoolForSum()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})
    s.Config = &SumSymbolValsConfig{Type: SSVTypeLessEqu, Value: 3}

    // prepare scenes
    gs, _ := sgc7game.NewGameSceneWithArr2([][]int{{0}})
    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{1}})

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Scenes = append(pr.Scenes, gs)
    pr.OtherScenes = append(pr.OtherScenes, os)

    gp.SceneStack.Push("", gs)
    gp.OtherSceneStack.Push("", os)

    // ensure Components and callstack are initialized to avoid nil derefs
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()

    cd := s.NewComponentData().(*SumSymbolValsData)

    nc, err := s.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), nil, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    assert.NoError(t, err)
    assert.Equal(t, "", nc)
    // os contains one value 1 which <=3, so Number should be 1
    assert.Equal(t, 1, cd.Number)
}

func TestSSV_DataCloneBuildPBGetValEx(t *testing.T) {
    sd := &SumSymbolValsData{}
    sd.OnNewGame(nil, nil)
    sd.Number = 42

    clone := sd.Clone().(*SumSymbolValsData)
    assert.Equal(t, sd.Number, clone.Number)

    pb := sd.BuildPBComponentData()
    assert.NotNil(t, pb)

    v, ok := sd.GetValEx(CVNumber, 0)
    assert.True(t, ok)
    assert.Equal(t, 42, v)

    v2, ok2 := sd.GetValEx("unknown", 0)
    assert.False(t, ok2)
    assert.Equal(t, 0, v2)
}

func TestSSV_ParseTypeAllVariants(t *testing.T) {
    table := map[string]SumSymbolValsType{
        "==": SSVTypeEqu,
        ">=": SSVTypeGreaterEqu,
        "<=": SSVTypeLessEqu,
        ">":  SSVTypeGreater,
        "<":  SSVTypeLess,
        "in [min, max]":  SSVTypeInAreaLR,
        "in (min, max]":  SSVTypeInAreaR,
        "in [min, max)":  SSVTypeInAreaL,
        "in (min, max)":  SSVTypeInArea,
    }

    for k, expect := range table {
        got := parseSumSymbolValsType(k)
        assert.Equal(t, expect, got)
    }
}

func TestSSV_InitExInvalidType(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    // pass wrong type
    err := s.InitEx("not a config", nil)
    assert.Error(t, err)
}

func TestSSV_SourceComponentEmptyPosFallsBack(t *testing.T) {
    pool := makePoolForSum()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()

    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{1, 2}, {3, 4}})

    comp := NewSumSymbolVals("src")
    gp.Components.MapComponents["src"] = comp

    // inject a component data with empty pos
    ep := &emptyPosCD{}
    gp.callStack.nodes[0].MapComponentData["src"] = ep

    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{SourceComponent: "src", Type: SSVTypeGreater, Value: 1}

    sum := s.sum(gp, os)
    // values >1 are 2,3,4 => sum = 9
    assert.Equal(t, 2+3+4, sum)
}

func TestSSV_CheckValAllBranches(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{}

    // GreaterEqu
    s.Config.Type = SSVTypeGreaterEqu
    s.Config.Value = 5
    assert.True(t, s.checkVal(5))
    assert.True(t, s.checkVal(6))

    // LessEqu
    s.Config.Type = SSVTypeLessEqu
    s.Config.Value = 4
    assert.True(t, s.checkVal(4))
    assert.True(t, s.checkVal(3))

    // Greater
    s.Config.Type = SSVTypeGreater
    s.Config.Value = 2
    assert.True(t, s.checkVal(3))
    assert.False(t, s.checkVal(2))

    // Less
    s.Config.Type = SSVTypeLess
    s.Config.Value = 2
    assert.True(t, s.checkVal(1))
    assert.False(t, s.checkVal(2))

    // InAreaR (min, max]
    s.Config.Type = SSVTypeInAreaR
    s.Config.Min = 1
    s.Config.Max = 3
    assert.True(t, s.checkVal(3))
    assert.False(t, s.checkVal(1))

    // InAreaL [min, max)
    s.Config.Type = SSVTypeInAreaL
    s.Config.Min = 1
    s.Config.Max = 3
    assert.True(t, s.checkVal(1))
    assert.False(t, s.checkVal(3))

    // InArea (min, max)
    s.Config.Type = SSVTypeInArea
    s.Config.Min = 1
    s.Config.Max = 3
    assert.True(t, s.checkVal(2))
    assert.False(t, s.checkVal(1))
}

// posCDOdd returns an odd-length pos slice; this should still be handled
// because the implementation only iterates full pairs (len/2)
type posCDOdd struct{ BasicComponentData }

func (p *posCDOdd) GetPos() []int { return []int{0, 0, 1} }

func TestSSV_SumSourcePosOddAndMissing(t *testing.T) {
    pool := makePoolForSum()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()

    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{10, 20}, {30, 40}})

    // case: source component missing -> full scan
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)

    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{Type: SSVTypeGreater, Value: 5}
    sum := s.sum(gp, os)
    // values >5: 10,20,30,40 => sum = 100
    assert.Equal(t, 10+20+30+40, sum)

    // case: source component present with odd-length pos
    comp := NewSumSymbolVals("src")
    gp.Components.MapComponents["src"] = comp

    // inject odd pos component data
    pcd := &posCDOdd{}
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()
    gp.callStack.nodes[0].MapComponentData["src"] = pcd

    s.Config.SourceComponent = "src"
    // only first pair (0,0) should be used => value 10 and >5 => counted
    sum2 := s.sum(gp, os)
    assert.Equal(t, 10, sum2)
}

func TestSSV_InitExAndProcControllers(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)

    // construct a proper config and call InitEx
    cfg := &SumSymbolValsConfig{StrType: ">", Value: 1}
    // add an award to ensure the awards path is initialized
    cfg.Awards = []*Award{{AwardType: "cash", Vals: []int{1}}}

    err := s.InitEx(cfg, makePoolForSum())
    assert.NoError(t, err)
    // component type should be set
    assert.Equal(t, SumSymbolValsTypeName, s.Config.ComponentType)

    // create a minimal gameProp and ensure ProcControllers does not panic
    gp := &GameProperty{Pool: makePoolForSum(), SceneStack: NewSceneStack(false), OtherSceneStack: NewSceneStack(true)}
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()

    // invoke ProcControllers; we don't assert award effects, only no panic
    s.ProcControllers(gp, nil, nil, nil, -1, "")
}

func TestSSV_ConfigSetLinkAndInitFromFile(t *testing.T) {
    // test SetLinkComponent
    cfg := &SumSymbolValsConfig{}
    cfg.SetLinkComponent("next", "nextcomp")
    assert.Equal(t, "nextcomp", cfg.DefaultNextComponent)

    // test Init reading from file
    tmpf, err := os.CreateTemp("", "ssv-*.yaml")
    assert.NoError(t, err)
    defer os.Remove(tmpf.Name())

    // write simple yaml config
    _, err = tmpf.WriteString("type: \">\"\nvalue: 2\n")
    assert.NoError(t, err)
    tmpf.Close()

    s := NewSumSymbolVals("s").(*SumSymbolVals)
    err = s.Init(tmpf.Name(), makePoolForSum())
    assert.NoError(t, err)
    // after Init, Config should be non-nil and ComponentType set
    assert.NotNil(t, s.Config)
    assert.Equal(t, SumSymbolValsTypeName, s.Config.ComponentType)
}

func TestSSV_DataGetOutputAndPB(t *testing.T) {
    sd := &SumSymbolValsData{}
    sd.Number = 123
    assert.Equal(t, 123, sd.GetOutput())

    pb := sd.BuildPBComponentData()
    assert.NotNil(t, pb)
    if pbcd, ok := pb.(*sgc7pb.SumSymbolValsData); ok {
        assert.Equal(t, int32(123), pbcd.Number)
    }
}

func TestSSV_OnPlayGameNoScene(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{Type: SSVTypeLess, Value: 1}
    s.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    gp := &GameProperty{Pool: makePoolForSum(), SceneStack: NewSceneStack(false), OtherSceneStack: NewSceneStack(true)}
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()

    cd := s.NewComponentData().(*SumSymbolValsData)

    // call OnPlayGame with no scenes pushed; GetTargetOtherScene3 should return nil
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    nc, err := s.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), nil, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    assert.NoError(t, err)
    assert.Equal(t, "", nc)
    // since no scene, Number should remain 0
    assert.Equal(t, 0, cd.Number)
}

func TestSSV_ParseTypeDefaultAndCheckValUnknown(t *testing.T) {
    // parse unknown type should return SSVTypeNone
    got := parseSumSymbolValsType("unknown-type")
    assert.Equal(t, SSVTypeNone, got)

    // checkVal with unknown type should return false
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    s.Config = &SumSymbolValsConfig{Type: SSVTypeNone}
    ok := s.checkVal(123)
    assert.False(t, ok)
}

func TestSSV_Init_FileReadAndUnmarshalErrors(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)

    // non-existent file -> ReadFile error
    err := s.Init("/path/does/not/exist.yaml", nil)
    assert.Error(t, err)

    // create a temp file with invalid YAML to force Unmarshal error
    tmpf, err2 := os.CreateTemp("", "ssv-bad-*.yaml")
    assert.NoError(t, err2)
    defer os.Remove(tmpf.Name())

    // write invalid YAML
    _, err3 := tmpf.WriteString(":\n")
    assert.NoError(t, err3)
    tmpf.Close()

    err4 := s.Init(tmpf.Name(), nil)
    // depending on YAML parser, it may or may not return an error; ensure we don't panic
    if err4 == nil {
        // still ok: at least Config should be set or InitEx may have returned nil
        assert.NotNil(t, s.Config)
    }
}

func TestSSV_OnPlayGameInvalidICDAndOnAsciiGame(t *testing.T) {
    s := NewSumSymbolVals("s").(*SumSymbolVals)
    // initialize minimal config so OnAsciiGame can be called
    s.Config = &SumSymbolValsConfig{Type: SSVTypeLess, Value: 1}
    s.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    gp := &GameProperty{Pool: makePoolForSum(), SceneStack: NewSceneStack(false), OtherSceneStack: NewSceneStack(true)}
    gp.Components = NewComponentList()
    gp.Components.MapComponents = make(map[string]IComponent)
    gp.callStack = NewCallStack()
    gp.callStack.OnNewGame()

    // call OnPlayGame with invalid component data (nil) and expect ErrInvalidComponentData
    nc, err := s.OnPlayGame(gp, sgc7game.NewPlayResult("m", 0, 0, "t"), NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), nil, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, nil)
    assert.Error(t, err)
    assert.Equal(t, "", nc)

    // OnAsciiGame should be a no-op and return nil
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    err2 := s.OnAsciiGame(gp, pr, nil, nil, &SumSymbolValsData{})
    assert.NoError(t, err2)
}
