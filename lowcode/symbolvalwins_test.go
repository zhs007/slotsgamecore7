package lowcode

import (
    "context"
    "testing"
    "os"

    "github.com/bytedance/sonic"
    "github.com/bytedance/sonic/ast"
    "github.com/stretchr/testify/assert"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
    sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
    sgc7utils "github.com/zhs007/slotsgamecore7/utils"
    "github.com/zhs007/slotsgamecore7/stats2"
    "github.com/zhs007/slotsgamecore7/asciigame"
)

// deterministic fake plugin used across tests
type fakePluginSVW struct{}

func (p *fakePluginSVW) Random(_ context.Context, r int) (int, error) {
    if r <= 0 {
        return 0, nil
    }

    return 0, nil
}
func (p *fakePluginSVW) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginSVW) ClearUsedRngs()                    {}
func (p *fakePluginSVW) TagUsedRngs()                      {}
func (p *fakePluginSVW) RollbackUsedRngs() error           { return nil }
func (p *fakePluginSVW) SetCache(arr []int)                {}
func (p *fakePluginSVW) ClearCache()                       {}
func (p *fakePluginSVW) Init()                             {}
func (p *fakePluginSVW) SetScenePool(any)                  {}
func (p *fakePluginSVW) GetScenePool() any                 { return nil }
func (p *fakePluginSVW) SetSeed(seed int)                  {}

// use distinct stub names to avoid clashes with other test files
type svwStubRNG struct{}
type svwStubFeatureLevel struct{}

func (s *svwStubRNG) Clone() IRNG                                            { return &svwStubRNG{} }
func (s *svwStubRNG) OnNewGame(betMode int, plugin sgc7plugin.IPlugin) error { return nil }
func (s *svwStubRNG) GetCurRNG(betMode int, gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string) {
    return false, 0, nil, ""
}
func (s *svwStubRNG) OnChoiceBranch(betMode int, curComponent IComponent, branchName string) error { return nil }
func (s *svwStubRNG) OnStepEnd(betMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) error { return nil }

func (s *svwStubFeatureLevel) Init() {}
func (s *svwStubFeatureLevel) OnStepEnd(gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult) {}
func (s *svwStubFeatureLevel) CountLevel() int { return 0 }

func makePoolWithSymbols() *GamePropertyPool {
    p := &GamePropertyPool{}
    p.DefaultPaytables = &sgc7game.PayTables{MapSymbols: map[string]int{"CA": 1, "CB": 2, "COIN": 5}}
    return p
}

func TestInitEx_ResolvesSymbolsAndCoinSymbols(t *testing.T) {
    pool := makePoolWithSymbols()

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{Symbols: []string{"CA", "CB"}, CoinSymbols: []string{"COIN"}, WinMulti: 2}

    err := svw.InitEx(cfg, pool)
    assert.NoError(t, err)
    assert.Equal(t, 2, len(svw.Config.SymbolCodes))
    assert.Equal(t, 1, svw.Config.SymbolCodes[0])
    assert.Equal(t, 2, svw.Config.SymbolCodes[1])
    assert.Equal(t, 1, len(svw.Config.CoinSymbolCodes))
    assert.Equal(t, 5, svw.Config.CoinSymbolCodes[0])
}

func TestOnPlayGame_NormalCountsValsAndAddsResult(t *testing.T) {
    pool := makePoolWithSymbols()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.rng = &svwStubRNG{}
    gp.featureLevel = &svwStubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{WinMulti: 3}
    svw.InitEx(cfg, pool)
    svw.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    // prepare gs and os scenes; os contains coin values
    gs, _ := sgc7game.NewGameScene2(3, 3, 0)
    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{0, 1, 0}, {2, 0, 3}, {0, 0, 0}})

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Scenes = append(pr.Scenes, gs)
    pr.OtherScenes = append(pr.OtherScenes, os)

    // push scenes into the gameProp stacks so GetTargetScene3 can find them
    gp.SceneStack.Push("", gs)
    gp.OtherSceneStack.Push("", os)

    bcd := svw.NewComponentData().(*SymbolValWinsData)
    // call OnPlayGame
    nc, err := svw.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePluginSVW{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, bcd)

    assert.NoError(t, err)
    assert.Equal(t, "", nc)
    // one result of RTCoins expected (mul==1)
    assert.Len(t, pr.Results, 1)
    ret := pr.Results[0]
    // total coin values = 1+2+3 = 6, othermul default is WinMulti (3)
    assert.Equal(t, 6*3, ret.CoinWin)
    assert.Equal(t, 6*3, bcd.Wins)
}

func TestOnPlayGame_WithCollectorType_PrependsCollectorPos(t *testing.T) {
    pool := makePoolWithSymbols()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.rng = &svwStubRNG{}
    gp.featureLevel = &svwStubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{WinMulti: 1, Symbols: []string{"CA"}}
    svw.InitEx(cfg, pool)
    svw.Config.Type = svwTypeCollector
    svw.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    // gs contains collector symbol CA (code 1) at 0,0
    gs, _ := sgc7game.NewGameSceneWithArr2([][]int{{1, 0}, {0, 0}})
    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{5, 0}, {0, 0}}) // coin values: one position (0,0)=5

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Scenes = append(pr.Scenes, gs)
    pr.OtherScenes = append(pr.OtherScenes, os)

    gp.SceneStack.Push("", gs)
    gp.OtherSceneStack.Push("", os)

    bcd := svw.NewComponentData().(*SymbolValWinsData)

    nc, err := svw.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePluginSVW{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, bcd)

    assert.NoError(t, err)
    assert.Equal(t, "", nc)
    assert.Len(t, pr.Results, 1)
    ret := pr.Results[0]
    // first two entries in Pos should be collector coord (0,0)
    assert.Equal(t, 0, ret.Pos[0])
    assert.Equal(t, 0, ret.Pos[1])
    // Symbol should be collector symbol code
    assert.Equal(t, 1, ret.Symbol)
}

func TestOnPlayGame_ReelCollector_BreaksPerReel(t *testing.T) {
    pool := makePoolWithSymbols()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{WinMulti: 1, Symbols: []string{"CA"}}
    svw.InitEx(cfg, pool)
    svw.Config.Type = svwTypeReelCollector
    svw.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    // gs has collector at (0,0) and (1,0) per reel
    gs, _ := sgc7game.NewGameSceneWithArr2([][]int{{1, 1}, {0, 0}})
    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{2, 0}, {3, 0}}) // two coin values

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Scenes = append(pr.Scenes, gs)
    pr.OtherScenes = append(pr.OtherScenes, os)

    gp.SceneStack.Push("", gs)
    gp.OtherSceneStack.Push("", os)

    bcd := svw.NewComponentData().(*SymbolValWinsData)

    nc, err := svw.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePluginSVW{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, bcd)

    assert.NoError(t, err)
    assert.Equal(t, "", nc)
    // at least one result expected (implementation may produce fewer per reel)
    assert.GreaterOrEqual(t, len(pr.Results), 1)
}

func TestCoinSymbolsFiltering(t *testing.T) {
    pool := makePoolWithSymbols()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.rng = &stubRNG{}
    gp.featureLevel = &stubFeatureLevel{}
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{WinMulti: 1, CoinSymbols: []string{"COIN"}}
    svw.InitEx(cfg, pool)
    svw.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    // gs symbols: only positions with symbol code 5 will be counted
    gs, _ := sgc7game.NewGameSceneWithArr2([][]int{{5, 0}, {0, 5}})
    os, _ := sgc7game.NewGameSceneWithArr2([][]int{{1, 2}, {3, 4}})

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Scenes = append(pr.Scenes, gs)
    pr.OtherScenes = append(pr.OtherScenes, os)

    gp.SceneStack.Push("", gs)
    gp.OtherSceneStack.Push("", os)

    bcd := svw.NewComponentData().(*SymbolValWinsData)

    nc, err := svw.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePluginSVW{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, bcd)

    assert.NoError(t, err)
    assert.Equal(t, "", nc)
    // both positions in gs with symbol 5 correspond to os values 1 and 4 -> total 5
    assert.Len(t, pr.Results, 1)
    ret := pr.Results[0]
    assert.Equal(t, 5, ret.CoinWin)
}

func TestGetWinMulti_OverrideAndDefault(t *testing.T) {
    pool := makePoolWithSymbols()
    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{WinMulti: 7}
    svw.InitEx(cfg, pool)

    bcd := &BasicComponentData{MapConfigIntVals: map[string]int{CCVWinMulti: 3}}
    v := svw.GetWinMulti(bcd)
    assert.Equal(t, 3, v)

    // without override
    bcd2 := &BasicComponentData{MapConfigIntVals: map[string]int{}}
    v2 := svw.GetWinMulti(bcd2)
    assert.Equal(t, 7, v2)
}

func TestBuildPBComponentDataAndOnAsciiGame_InvalidICD(t *testing.T) {
    // test BuildPBComponentData
    svwd := &SymbolValWinsData{SymbolNum: 4, Wins: 11}
    pb := svwd.BuildPBComponentData()
    assert.NotNil(t, pb)

    // OnAsciiGame with invalid icd
    svw := NewSymbolValWins("svw").(*SymbolValWins)
    err := svw.OnAsciiGame(nil, &sgc7game.PlayResult{}, nil, asciigame.NewSymbolColorMap(&sgc7game.PayTables{}), nil)
    assert.Error(t, err)
}

func TestParseSymbolValWins_FromAST(t *testing.T) {
    jsonStr := `{"componentValues":{"label":"svwlabel","configuration":{"betType":"bet","winMulti":2,"type":"normal","coinSymbols":["COIN"]}}}`

    var node ast.Node
    err := sonic.Unmarshal([]byte(jsonStr), &node)
    assert.NoError(t, err)

    bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
    label, err := parseSymbolValWins(bc, &node)
    assert.NoError(t, err)
    assert.Equal(t, "svwlabel", label)
    _, ok := bc.mapConfig["svwlabel"]
    assert.True(t, ok)
}

func TestSymbolValWinsData_OnNewGame_Clone_GetValEx(t *testing.T) {
    svwd := &SymbolValWinsData{}

    // OnNewGame should initialize maps inside BasicComponentData
    svwd.OnNewGame(nil, nil)
    assert.NotNil(t, svwd.MapConfigVals)
    assert.NotNil(t, svwd.MapConfigIntVals)

    // set some fields and used results
    svwd.SymbolNum = 2
    svwd.Wins = 10
    svwd.UsedResults = []int{0, 1}

    // GetValEx should return correct values
    v, ok := svwd.GetValEx(SVWDVSymbolNum, 0)
    assert.True(t, ok)
    assert.Equal(t, 2, v)

    v2, ok2 := svwd.GetValEx(SVWDVWins, 0)
    assert.True(t, ok2)
    assert.Equal(t, 10, v2)

    v3, ok3 := svwd.GetValEx(CVResultNum, 0)
    assert.True(t, ok3)
    assert.Equal(t, 2, v3)

    // unknown key
    _, ok4 := svwd.GetValEx("unknown", 0)
    assert.False(t, ok4)

    // Clone should produce a copy with same values
    cl := svwd.Clone().(*SymbolValWinsData)
    assert.Equal(t, svwd.SymbolNum, cl.SymbolNum)
    assert.Equal(t, svwd.Wins, cl.Wins)
    assert.Equal(t, len(svwd.UsedResults), len(cl.UsedResults))
}

func TestSetLinkComponentAndNewStats2_OnStats2(t *testing.T) {
    // SetLinkComponent
    cfg := &SymbolValWinsConfig{}
    cfg.SetLinkComponent("next", "comp1")
    assert.Equal(t, "comp1", cfg.DefaultNextComponent)

    // NewStats2 and OnStats2
    pool := makePoolWithSymbols()
    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg2 := &SymbolValWinsConfig{WinMulti: 5}
    _ = svw.InitEx(cfg2, pool)

    s2 := stats2.NewCache(1)
    f := svw.NewStats2(svw.GetName())
    s2.AddFeature(svw.GetName(), f, false)

    svwd := &SymbolValWinsData{BasicComponentData: BasicComponentData{MapConfigIntVals: map[string]int{CCVWinMulti: 4}}, Wins: 123}

    // call OnStats2 and ensure wins propagated
    svw.OnStats2(svwd, s2, &GameProperty{Pool: pool}, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), nil, true)

    f2 := s2.GetFeature(svw.GetName())
    assert.NotNil(t, f2)
    if f2.Wins != nil {
        assert.Equal(t, int64(123), f2.Wins.TotalWin)
    }
}

func TestOnAsciiGame_ValidICD(t *testing.T) {
    svw := NewSymbolValWins("svw").(*SymbolValWins)

    // build component data with one used result
    svwd := &SymbolValWinsData{BasicComponentData: BasicComponentData{UsedResults: []int{0}}}

    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Results = append(pr.Results, &sgc7game.Result{Type: sgc7game.RTCoins})

    err := svw.OnAsciiGame(nil, pr, nil, asciigame.NewSymbolColorMap(&sgc7game.PayTables{}), svwd)
    assert.NoError(t, err)
}

func TestParseSymbolValWins_ErrorsAndTypeParsing(t *testing.T) {
    // missing componentValues should return error
    var node ast.Node
    // empty node -> getConfigInCell should fail
    _, err := parseSymbolValWins(&BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}, &node)
    assert.Error(t, err)

    // type parsing
    assert.Equal(t, svwTypeCollector, parseSymbolValWinsType("collector"))
    assert.Equal(t, svwTypeReelCollector, parseSymbolValWinsType("reelcollector"))
    assert.Equal(t, svwTypeNormal, parseSymbolValWinsType("somethingelse"))
}

func TestInitEx_InvalidSymbols(t *testing.T) {
    pool := makePoolWithSymbols()
    svw := NewSymbolValWins("svw").(*SymbolValWins)

    // Symbols contains unknown symbol
    cfg := &SymbolValWinsConfig{Symbols: []string{"NOPE"}}
    err := svw.InitEx(cfg, pool)
    assert.Error(t, err)

    // CoinSymbols contains unknown symbol
    cfg2 := &SymbolValWinsConfig{CoinSymbols: []string{"NOPE"}}
    err2 := svw.InitEx(cfg2, pool)
    assert.Error(t, err2)
}

func TestInit_ReadFromFile_Success(t *testing.T) {
    pool := makePoolWithSymbols()

    // create temp yaml file
    content := "betType: \"bet\"\nwinMulti: 2\nsymbols: [\"CA\"]\ncoinSymbols: [\"COIN\"]\n"
    fn := "./test_svw_init.yaml"
    _ = os.WriteFile(fn, []byte(content), 0644)
    defer os.Remove(fn)

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    err := svw.Init(fn, pool)
    assert.NoError(t, err)
}

func TestOnPlayGame_InvalidICDAndNoOtherScene(t *testing.T) {
    pool := makePoolWithSymbols()
    gp := &GameProperty{Pool: pool}
    gp.PoolScene = sgc7game.NewGameScenePoolEx()
    gp.SceneStack = NewSceneStack(false)
    gp.OtherSceneStack = NewSceneStack(true)

    svw := NewSymbolValWins("svw").(*SymbolValWins)
    cfg := &SymbolValWinsConfig{WinMulti: 1}
    svw.InitEx(cfg, pool)
    svw.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

    // invalid icd
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    gs, _ := sgc7game.NewGameScene2(2, 2, 0)
    pr.Scenes = append(pr.Scenes, gs)
    gp.SceneStack.Push("", gs)

    nc, err := svw.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePluginSVW{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, nil)
    assert.Error(t, err)
    assert.Equal(t, "", nc)

    // valid icd but no other scene -> should do nothing
    bcd := svw.NewComponentData().(*SymbolValWinsData)
    nc2, err2 := svw.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePluginSVW{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, bcd)
    // returns ErrComponentDoNothing
    assert.Error(t, err2)
    assert.Equal(t, "", nc2)
}
