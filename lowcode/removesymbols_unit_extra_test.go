package lowcode

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
    "github.com/zhs007/slotsgamecore7/asciigame"
)

// focused unit tests to cover small helpers in removesymbols.go
func TestRemoveSymbols_InitFromFileAndInitExSuccess(t *testing.T) {
    // create a temporary YAML config
    data := `type: "adjacentPay"
addedSymbol: "WL"
ignoreSymbols: ["A"]
targetComponents: ["bg-pay"]
emptySymbolVal: -2
outputToComp: "out"
`
    tmp, err := os.CreateTemp("", "rs-*.yaml")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmp.Name())

    _, err = tmp.WriteString(data)
    if err != nil {
        t.Fatal(err)
    }
    tmp.Close()

    pool := &GamePropertyPool{DefaultPaytables: &sgc7game.PayTables{MapSymbols: map[string]int{"WL": 99, "A": 5}}, Config: &Config{Width: 3, Height: 3}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    err = comp.Init(tmp.Name(), pool)
    assert.NoError(t, err)
    assert.NotNil(t, comp.Config)
    assert.Equal(t, 99, comp.Config.AddedSymbolCode)
    // ignore symbol resolved
    assert.Equal(t, 1, len(comp.Config.IgnoreSymbolCodes))
    assert.Equal(t, RSTypeAdjacentPay, comp.Config.Type)
}

func TestRemoveSymbols_OnAsciiGame_And_ComponentDataHelpers(t *testing.T) {
    // prepare a trivial scene and playresult
    mat := [][]int{{1}}
    sc, _ := sgc7game.NewGameSceneWithArr2(mat)
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr.Scenes = append(pr.Scenes, sc)

    pool := &GamePropertyPool{DefaultPaytables: &sgc7game.PayTables{MapSymbols: map[string]int{"A": 1}}, Config: &Config{Width: 1, Height: 1}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

    gp := pool.newGameProp(1)

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    // ensure NewComponentData returns expected type
    cd := comp.NewComponentData()
    _, ok := cd.(*RemoveSymbolsData)
    assert.True(t, ok)

    // exercise Clone/BuildPBComponentData
    rsd := &RemoveSymbolsData{RemovedNum: 2}
    cloned := rsd.Clone()
    _, ok = cloned.(*RemoveSymbolsData)
    assert.True(t, ok)
    pb := rsd.BuildPBComponentData()
    assert.NotNil(t, pb)

    // exercise OnAsciiGame when UsedScenes populated
    rsd.UsedScenes = []int{0}
    scm := asciigame.NewSymbolColorMap(pool.DefaultPaytables)
    // OnGetSymbolString is already set in NewSymbolColorMap
    err := comp.OnAsciiGame(gp, pr, nil, scm, rsd)
    assert.NoError(t, err)
}

func TestRemoveSymbols_SetLink_GetLinks_EachUsedResults_JsonBuild(t *testing.T) {
    // SetLinkComponent usage
    cfg := &RemoveSymbolsConfig{}
    cfg.SetLinkComponent("next", "nxt")
    cfg.SetLinkComponent("jump", "jmp")
    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.Config = cfg
    // GetAllLinkComponents and GetNextLinkComponents
    all := comp.GetAllLinkComponents()
    next := comp.GetNextLinkComponents()
    assert.Contains(t, all, "nxt")
    assert.Contains(t, next, "jmp")

    // EachUsedResults is a no-op; ensure calling it doesn't panic
    comp.EachUsedResults(nil, nil, nil)

    // json helper build
    j := &jsonRemoveSymbols{Type: "adjacentPay", AddedSymbol: "WL", TargetComponents: []string{"a"}, IgnoreSymbols: []string{"A"}, IsNeedProcSymbolVals: true, EmptySymbolVal: -1, OutputToComponent: "out", SourcePositionCollection: []string{"src"}}
    full := j.build()
    assert.Equal(t, "WL", full.AddedSymbol)
    assert.Equal(t, -1, full.EmptySymbolVal)
}

func TestRemoveSymbols_ProcControllers_NoPanic(t *testing.T) {
    pool := &GamePropertyPool{DefaultPaytables: &sgc7game.PayTables{MapSymbols: map[string]int{"A": 1}}, Config: &Config{Width: 1, Height: 1}}
    pool.newRNG = func() IRNG { return &stubRNG{} }
    pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }
    gp := pool.newGameProp(1)

    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.Config = &RemoveSymbolsConfig{Awards: []*Award{{AwardType: "respinTimes", Vals: []int{1}, StrParams: []string{"a"}}}}

    // should not panic
    comp.ProcControllers(gp, nil, nil, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), 0, "")
}

func TestOnPlayGame_Basic_And_AdjacentPay_CoverBranches(t *testing.T) {
    // helper to setup common pool and gameProp
    makeEnv := func(mat [][]int) (*GameProperty, *sgc7game.PlayResult, *GameParams) {
        sc, _ := sgc7game.NewGameSceneWithArr2(mat)
        pr := sgc7game.NewPlayResult("m", 0, 0, "t")
        pr.Scenes = append(pr.Scenes, sc)

        pool := &GamePropertyPool{Config: &Config{Width: len(mat), Height: len(mat[0])}}
        pool.newRNG = func() IRNG { return &stubRNG{} }
        pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }
        gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

        gameProp := pool.newGameProp(1)
        gameProp.Components = NewComponentList()
        _ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

        // ensure scene stack has the base scene
        gameProp.SceneStack.Push("rs", sc)

        return gameProp, pr, gp
    }

    // positions for a 3-length result in middle row
    mat := [][]int{{1,1,1},{1,1,1},{1,1,1}}
    gameProp, pr, gp := makeEnv(mat)

    // create a fake result (3 in a row horizontally at y=1)
    r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex: 0, Symbol: 0, Pos: []int{0,1,1,1,2,1}}
    pr.Results = append(pr.Results, r)

    // fake component data that reports this result index
    frcd := &rsFakeCD{results: []int{0}}
    // set callStack global map
    gameProp.callStack = NewCallStack()
    gameProp.callStack.OnNewGame()
    gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": frcd}

    // register fake component
    fakeComp := &rsFakeComp{name: "bg-pay"}
    gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": fakeComp}

    // mark history
    gp.HistoryComponents = []string{"bg-pay"}

    // Test RSTypeBasic without other scene
    comp := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp.BasicComponent = NewBasicComponent("rs", 1)
    comp.BasicComponent.Config = &BasicComponentConfig{}
    comp.Config = &RemoveSymbolsConfig{Type: RSTypeBasic, TargetComponents: []string{"bg-pay"}}
    cd := comp.NewComponentData()
    _, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
    // expect success and removed positions appended
    assert.NoError(t, err)
    rsd := cd.(*RemoveSymbolsData)
    assert.Greater(t, rsd.RemovedNum, 0)

    // Test RSTypeBasic with other scene (IsNeedProcSymbolVals true)
    // prepare other scene
    other := pr.Scenes[0].CloneEx(gameProp.PoolScene)
    pr.OtherScenes = append(pr.OtherScenes, other)
    gameProp.OtherSceneStack.Push("rs", other)

    comp2 := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp2.BasicComponent = NewBasicComponent("rs", 1)
    comp2.BasicComponent.Config = &BasicComponentConfig{}
    comp2.Config = &RemoveSymbolsConfig{Type: RSTypeBasic, TargetComponents: []string{"bg-pay"}, IsNeedProcSymbolVals: true, EmptySymbolVal: -5}
    cd2 := comp2.NewComponentData()
    _, err = comp2.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd2)
    assert.NoError(t, err)
    rsd2 := cd2.(*RemoveSymbolsData)
    assert.GreaterOrEqual(t, rsd2.RemovedNum, 0)

    // Test RSTypeAdjacentPay without other scene
    // reset playresult scenes to original
    mat2 := [][]int{{1,1,1},{1,1,1},{1,1,1}}
    sc2, _ := sgc7game.NewGameSceneWithArr2(mat2)
    pr2 := sgc7game.NewPlayResult("m", 0, 0, "t")
    pr2.Scenes = append(pr2.Scenes, sc2)
    pr2.Results = append(pr2.Results, &sgc7game.Result{Type: sgc7game.RTAdjacentPay, Pos: []int{0,1,1,1,2,1}})

    gameProp2, _, gp2 := makeEnv(mat2)
    gameProp2.callStack = NewCallStack()
    gameProp2.callStack.OnNewGame()
    gameProp2.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": &rsFakeCD{results: []int{0}}}
    gameProp2.Components.MapComponents = map[string]IComponent{"bg-pay": &rsFakeComp{name: "bg-pay"}}
    gp2.HistoryComponents = []string{"bg-pay"}

    comp3 := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp3.BasicComponent = NewBasicComponent("rs", 1)
    comp3.BasicComponent.Config = &BasicComponentConfig{}
    comp3.Config = &RemoveSymbolsConfig{Type: RSTypeAdjacentPay, TargetComponents: []string{"bg-pay"}, AddedSymbolCode: 77}
    cd3 := comp3.NewComponentData()
    _, err = comp3.OnPlayGame(gameProp2, pr2, gp2, nil, "", "", nil, nil, nil, cd3)
    assert.NoError(t, err)
    rsd3 := cd3.(*RemoveSymbolsData)
    assert.Equal(t, 3, rsd3.RemovedNum)

    // Test RSTypeAdjacentPay with other scene
    other2 := pr2.Scenes[0].CloneEx(gameProp2.PoolScene)
    pr2.OtherScenes = append(pr2.OtherScenes, other2)
    gameProp2.OtherSceneStack.Push("rs", other2)
    comp4 := NewRemoveSymbols("rs").(*RemoveSymbols)
    comp4.BasicComponent = NewBasicComponent("rs", 1)
    comp4.BasicComponent.Config = &BasicComponentConfig{}
    comp4.Config = &RemoveSymbolsConfig{Type: RSTypeAdjacentPay, TargetComponents: []string{"bg-pay"}, AddedSymbolCode: 88, IsNeedProcSymbolVals: true, EmptySymbolVal: -9}
    cd4 := comp4.NewComponentData()
    _, err = comp4.OnPlayGame(gameProp2, pr2, gp2, nil, "", "", nil, nil, nil, cd4)
    assert.NoError(t, err)
    rsd4 := cd4.(*RemoveSymbolsData)
    assert.GreaterOrEqual(t, rsd4.RemovedNum, 0)
}
