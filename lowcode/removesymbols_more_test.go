package lowcode

import (
    "testing"
    "github.com/fatih/color"

    "github.com/stretchr/testify/assert"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
    "github.com/zhs007/slotsgamecore7/asciigame"
    "google.golang.org/protobuf/types/known/anypb"
)

func TestParseRemoveSymbolsType(t *testing.T) {
    assert.Equal(t, RSTypeAdjacentPay, parseRemoveSymbolsType("adjacentPay"))
    assert.Equal(t, RSTypeBasic, parseRemoveSymbolsType("somethingelse"))
}

func TestRemoveSymbolsDataHelpers(t *testing.T) {
    rd := &RemoveSymbolsData{RemovedNum: 2, AvgHeight: 123}

    // GetValEx for avg height
    v, ok := rd.GetValEx(CVAvgHeight, 0)
    assert.True(t, ok)
    assert.Equal(t, 123, v)

    // onNewStep resets fields
    rd.onNewStep()
    assert.Equal(t, 0, rd.RemovedNum)
    assert.Equal(t, 0, rd.AvgHeight)

    // Clone returns same type
    rd.RemovedNum = 5
    c := rd.Clone()
    _, ok2 := c.(*RemoveSymbolsData)
    assert.True(t, ok2)

    // BuildPBComponentData returns proto message
    pb := rd.BuildPBComponentData()
    assert.NotNil(t, pb)
}

func TestConfigSetLinkAndInitExSuccess(t *testing.T) {
    pool := &GamePropertyPool{DefaultPaytables: &sgc7game.PayTables{MapSymbols: map[string]int{"WL": 99, "A": 5}}, Config: &Config{Width:1, Height:1}}

    cfg := &RemoveSymbolsConfig{AddedSymbol: "WL", IgnoreSymbols: []string{"A"}}

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    err := comp.InitEx(cfg, pool)
    assert.NoError(t, err)
    assert.Equal(t, 99, comp.Config.AddedSymbolCode)
    assert.Equal(t, []int{5}, comp.Config.IgnoreSymbolCodes)

    // SetLinkComponent
    comp.Config.SetLinkComponent("next", "n1")
    comp.Config.SetLinkComponent("jump", "j1")
    assert.Equal(t, "n1", comp.Config.DefaultNextComponent)
    assert.Equal(t, "j1", comp.Config.JumpToComponent)
}

func TestOnAsciiGameAndEachUsedResultsNoop(t *testing.T) {
    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.BasicComponent = NewBasicComponent("rs", 1)

    // prepare playresult with one scene
    mat := [][]int{{1}}
    sc, _ := sgc7game.NewGameSceneWithArr2(mat)
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    pr.Scenes = append(pr.Scenes, sc)

    cd := &RemoveSymbolsData{}
    // Use a UsedScenes entry so OnAsciiGame prints
    cd.UsedScenes = []int{0}

    // should not panic - provide a non-nil SymbolColorMap
    scm := &asciigame.SymbolColorMap{MapSymbols: make(map[int]*color.Color)}
    scm.OnGetSymbolString = func(i int) string { return "" }
    err := comp.OnAsciiGame(nil, pr, nil, scm, cd)
    assert.NoError(t, err)

    // EachUsedResults is a noop; pass nil pb component data and a func
    comp.EachUsedResults(pr, nil, func(r *sgc7game.Result) {})
}

func TestNewRemoveSymbolsAndJsonBuild(t *testing.T) {
    // new remove symbols basic checks
    comp := NewRemoveSymbols("x")
    assert.NotNil(t, comp)

    // json build
    j := &jsonRemoveSymbols{Type: "adjacentPay", AddedSymbol: "X", TargetComponents: []string{"a"}, IgnoreSymbols: []string{"b"}, SourcePositionCollection: []string{"c"}}
    cfg := j.build()
    assert.Equal(t, "adjacentPay", cfg.StrType)
    assert.Equal(t, []string{"a"}, cfg.TargetComponents)
    assert.Equal(t, []string{"b"}, cfg.IgnoreSymbols)
    assert.Equal(t, []string{"c"}, cfg.SourcePositionCollection)

    // GetAll/GetNext link components should return expected length
    rs := comp.(*RemoveSymbols)
    rs.Config = &RemoveSymbolsConfig{}
    all := rs.GetAllLinkComponents()
    next := rs.GetNextLinkComponents()
    assert.Len(t, all, 2)
    assert.Len(t, next, 2)

    // call NewComponentData
    cd := rs.NewComponentData()
    _, ok := cd.(*RemoveSymbolsData)
    assert.True(t, ok)

    // call EachUsedResults with a non-nil anypb to ensure no panic
    var any anypb.Any
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    rs.EachUsedResults(pr, &any, func(r *sgc7game.Result) {})
}

func TestOnBasic_TargetComponents_NoOtherScene(t *testing.T) {
    // produce a scenario where TargetComponents cause removal and os == nil
    mat := [][]int{
        {1,1,1},
        {1,1,1},
        {1,1,1},
    }
    sc, _ := sgc7game.NewGameSceneWithArr2(mat)
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    pr.Scenes = append(pr.Scenes, sc)

    pool := &GamePropertyPool{Config: &Config{Width:3, Height:3}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }
    gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

    gameProp := pool.newGameProp(1)
    gameProp.Components = NewComponentList()
    // register a fake history component so GetCurComponentDataWithName can resolve it
    gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": &rsFakeComp{name: "bg-pay", results: []int{0}}}
    _ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
    gameProp.SceneStack.Push("rs", sc)

    // prepare a result in pr
    r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex:0, Symbol:0, Pos: []int{0,1,1,1}}
    pr.Results = append(pr.Results, r)

    // fake prev comp data
    fcd := &rsFakeCD{results: []int{0}}
    gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": fcd}
    gp.HistoryComponents = []string{"bg-pay"}

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.BasicComponent = NewBasicComponent("rs", 1)
    comp.BasicComponent.Config = &BasicComponentConfig{}
    comp.Config = &RemoveSymbolsConfig{TargetComponents: []string{"bg-pay"}}

    cd := comp.NewComponentData()
    _, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
    assert.NoError(t, err)
    rsd := cd.(*RemoveSymbolsData)
    assert.Equal(t, 2, rsd.RemovedNum)
}

func TestOnAdjacentPay_NoOtherScene(t *testing.T) {
    mat := [][]int{{1,1,1},{1,1,1},{1,1,1}}
    sc, _ := sgc7game.NewGameSceneWithArr2(mat)
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    pr.Scenes = append(pr.Scenes, sc)

    pool := &GamePropertyPool{Config: &Config{Width:3, Height:3}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }
    gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

    gameProp := pool.newGameProp(1)
    gameProp.Components = NewComponentList()
    // register a fake history component so GetCurComponentDataWithName can resolve it
    gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": &rsFakeComp{name: "bg-pay", results: []int{0}}}
    _ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
    gameProp.SceneStack.Push("rs", sc)

    r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex:0, Symbol:0, Pos: []int{0,1,1,1,2,1}}
    pr.Results = append(pr.Results, r)

    fcd := &rsFakeCD{results: []int{0}}
    gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": fcd}
    gp.HistoryComponents = []string{"bg-pay"}

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.BasicComponent = NewBasicComponent("rs", 1)
    comp.BasicComponent.Config = &BasicComponentConfig{}
    comp.Config = &RemoveSymbolsConfig{Type: RSTypeAdjacentPay, AddedSymbolCode: 77, TargetComponents: []string{"bg-pay"}}

    cd := comp.NewComponentData()
    _, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
    assert.NoError(t, err)
    rsd := cd.(*RemoveSymbolsData)
    assert.Equal(t, 3, rsd.RemovedNum)
}

func TestProcControllersWithAwards(t *testing.T) {
    // basic smoke to ensure ProcControllers calls procAwards without panic
    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.Config = &RemoveSymbolsConfig{Awards: []*Award{{AwardType: "cash", Vals: []int{10}}}}

    pool := &GamePropertyPool{Config: &Config{Width:1, Height:1}, DefaultPaytables: &sgc7game.PayTables{MapSymbols: map[string]int{"A":1}}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }
    gameProp := pool.newGameProp(1)

    // should not panic
    comp.ProcControllers(gameProp, nil, sgc7game.NewPlayResult("m",0,0,"t"), NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil), -1, "")
}

func TestBuildPBAndEachUsedResults_GetLinks(t *testing.T) {
    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    cd := &RemoveSymbolsData{RemovedNum: 4}

    // BuildPBComponentData should return a proto message
    pb := cd.BuildPBComponentData()
    if pb == nil {
        t.Fatalf("BuildPBComponentData returned nil")
    }

    // EachUsedResults is a noop and should not panic when passed nil
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    comp.EachUsedResults(pr, nil, func(r *sgc7game.Result) { t.Log("ok") })

    // GetAll/GetNext link components should return a slice of length 2 even if cfg nil
    comp.Config = &RemoveSymbolsConfig{}
    all := comp.GetAllLinkComponents()
    next := comp.GetNextLinkComponents()
    assert.Len(t, all, 2)
    assert.Len(t, next, 2)
}

func TestOnAsciiGame_NoUsedScenesAndInitExErrors(t *testing.T) {
    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    // call OnAsciiGame with component data that has no UsedScenes - should noop
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    sc := &asciigame.SymbolColorMap{MapSymbols: map[int]*color.Color{}}
    sc.OnGetSymbolString = func(i int) string { return "" }
    err := comp.OnAsciiGame(nil, pr, nil, sc, &RemoveSymbolsData{})
    assert.NoError(t, err)

    // InitEx with wrong type should error (already tested elsewhere but exercise here)
    err = comp.InitEx(nil, &GamePropertyPool{})
    assert.Error(t, err)
}

func TestOnBasic_WithOtherSceneAndOutput(t *testing.T) {
    // 3x3 scene
    mat := [][]int{{1,1,1},{1,1,1},{1,1,1}}
    sc, _ := sgc7game.NewGameSceneWithArr2(mat)

    pr := sgc7game.NewPlayResult("m",0,0,"t")
    pr.Scenes = append(pr.Scenes, sc)

    pool := &GamePropertyPool{Config: &Config{Width:3, Height:3}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

    gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

    gameProp := pool.newGameProp(1)
    gameProp.Components = NewComponentList()
    // register components
    gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": &rsFakeComp{name: "bg-pay", results: []int{0}}, "out": &rsFakeComp{name: "out"}}
    _ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
    gameProp.SceneStack.Push("rs", sc)
    // prepare other scene
    other := sc.CloneEx(gameProp.PoolScene)
    pr.OtherScenes = append(pr.OtherScenes, other)
    gameProp.OtherSceneStack.Push("rs", other)

    // prepare result referring to two positions
    r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex:0, Symbol:0, Pos: []int{0,1,1,1}}
    pr.Results = append(pr.Results, r)

    // put fake prev comp data into global map so GetCurComponentDataWithName returns it
    fcd := &rsFakeCD{results: []int{0}}
    gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": fcd, "out": &rsOutCD{}}
    gp.HistoryComponents = []string{"bg-pay"}

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.BasicComponent = NewBasicComponent("rs", 1)
    comp.BasicComponent.Config = &BasicComponentConfig{}
    comp.Config = &RemoveSymbolsConfig{TargetComponents: []string{"bg-pay"}, IsNeedProcSymbolVals: true, OutputToComponent: "out", EmptySymbolVal: -9}

    cd := comp.NewComponentData()
    _, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
    assert.NoError(t, err)

    rsd := cd.(*RemoveSymbolsData)
    // two positions removed
    assert.Equal(t, 2, rsd.RemovedNum)

    // output component recorded positions
    outcd, _ := gameProp.callStack.nodes[0].MapComponentData["out"].(*rsOutCD)
    assert.Len(t, outcd.pos, 4)

    // new scenes appended
    assert.GreaterOrEqual(t, len(pr.Scenes), 2)
    assert.GreaterOrEqual(t, len(pr.OtherScenes), 2)
}

func TestAdjacentPay_WithOtherScene_RetainMiddleAndAddSymbol(t *testing.T) {
    mat := [][]int{{1,1,1},{1,1,1},{1,1,1}}
    sc, _ := sgc7game.NewGameSceneWithArr2(mat)
    pr := sgc7game.NewPlayResult("m",0,0,"t")
    pr.Scenes = append(pr.Scenes, sc)

    pool := &GamePropertyPool{Config: &Config{Width:3, Height:3}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

    gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

    gameProp := pool.newGameProp(1)
    gameProp.Components = NewComponentList()
    gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": &rsFakeComp{name: "bg-pay", results: []int{0}}}
    _ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
    gameProp.SceneStack.Push("rs", sc)
    other := sc.CloneEx(gameProp.PoolScene)
    pr.OtherScenes = append(pr.OtherScenes, other)
    gameProp.OtherSceneStack.Push("rs", other)

    // 3-in-row positions -> middle should be retained (not removed) and later set to AddedSymbolCode
    r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex:0, Symbol:0, Pos: []int{0,1,1,1,2,1}}
    pr.Results = append(pr.Results, r)

    fcd := &rsFakeCD{results: []int{0}}
    gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": fcd}
    gp.HistoryComponents = []string{"bg-pay"}

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.BasicComponent = NewBasicComponent("rs", 1)
    comp.BasicComponent.Config = &BasicComponentConfig{}
    comp.Config = &RemoveSymbolsConfig{Type: RSTypeAdjacentPay, AddedSymbolCode: 77, TargetComponents: []string{"bg-pay"}, IsNeedProcSymbolVals: true, EmptySymbolVal: -1}

    cd := comp.NewComponentData()
    _, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
    assert.NoError(t, err)

    rsd := cd.(*RemoveSymbolsData)
    // should have removed three positions (logic increments RemovedNum for all checked positions)
    assert.Equal(t, 3, rsd.RemovedNum)

    // the ngs (new scene) should have the middle position set to AddedSymbolCode (77)
    newsc := pr.Scenes[len(pr.Scenes)-1]
    // middle is (1,1)
    assert.Equal(t, 77, newsc.Arr[1][1])
}
