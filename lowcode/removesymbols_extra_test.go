package lowcode

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
    "github.com/zhs007/slotsgamecore7/asciigame"
)

// Test exercising multiple small helpers and code paths in removesymbols.go
func TestRemoveSymbols_ExerciseHelpers(t *testing.T) {
    // NewRemoveSymbols, NewComponentData
    rs := NewRemoveSymbols("rs1").(*RemoveSymbols)
    assert.NotNil(t, rs)

    cd := rs.NewComponentData()
    assert.NotNil(t, cd)

    // RemoveSymbolsData helpers
    rsd := &RemoveSymbolsData{}
    // onNewStep should initialize fields without panic
    rsd.onNewStep()
    assert.Equal(t, 0, rsd.RemovedNum)

    // GetValEx CVAvgHeight
    v, ok := rsd.GetValEx(CVAvgHeight, 0)
    assert.True(t, ok)
    assert.Equal(t, 0, v)

    // Clone and BuildPBComponentData
    cl := rsd.Clone()
    assert.NotNil(t, cl)

    pb := rsd.BuildPBComponentData()
    assert.NotNil(t, pb)

    // SetLinkComponent variants
    cfg := &RemoveSymbolsConfig{}
    cfg.SetLinkComponent("next", "n1")
    cfg.SetLinkComponent("jump", "j1")
    assert.Equal(t, "n1", cfg.DefaultNextComponent)
    assert.Equal(t, "j1", cfg.JumpToComponent)

    // NewRemoveSymbols default next/jump getters
    rs.Config = &RemoveSymbolsConfig{}
    all := rs.GetAllLinkComponents()
    next := rs.GetNextLinkComponents()
    assert.Len(t, all, 2)
    assert.Len(t, next, 2)

    // EachUsedResults is a no-op; ensure calling it doesn't panic
    pr := &sgc7game.PlayResult{}
    rs.EachUsedResults(pr, nil, nil)

    // OnAsciiGame: when UsedScenes empty should no-op
    gp := &GameProperty{}
    // supply a simple symbol color map
    scm := asciigame.NewSymbolColorMap(nil)
    err := rs.OnAsciiGame(gp, pr, nil, scm, &RemoveSymbolsData{})
    assert.Nil(t, err)

    // Init with invalid file should return error
    err2 := rs.Init("/no/such/file.yaml", nil)
    assert.Error(t, err2)

    // InitEx invalid cfg should return error
    err3 := rs.InitEx(nil, nil)
    assert.Error(t, err3)

    // json builder build function
    j := &jsonRemoveSymbols{Type: "basic", AddedSymbol: "A", TargetComponents: []string{"x"}}
    cfg2 := j.build()
    assert.NotNil(t, cfg2)

    // parseRemoveSymbols with bad AST should error - use empty gamecfg
    // We cannot easily construct an AST node here; skip calling parseRemoveSymbols

    // cleanup: ensure temporary variables don't leak
    _ = os.TempDir()
}
