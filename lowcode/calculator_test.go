package lowcode

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/bytedance/sonic"
    "github.com/bytedance/sonic/ast"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func TestNewCalculatorAndNewComponentData(t *testing.T) {
    comp := NewCalculator("calc")
    assert.NotNil(t, comp)

    cd := comp.NewComponentData()
    assert.IsType(t, &BasicComponentData{}, cd)

    // OnPlayGame should proceed to default next
    comp.(*Calculator).BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: "next"})

    gp := &GameProperty{Pool: &GamePropertyPool{}}
    pr := sgc7game.NewPlayResult("m", 0, 0, "t")
    plugin := &fakePluginTC{}

    nc, err := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, cd)
    assert.NoError(t, err)
    assert.Equal(t, "next", nc)
}

func TestParseCalculator(t *testing.T) {
    jsonStr := `{"componentValues": {"label": "calc1", "configuration": {}}}`

    var node ast.Node
    err := sonic.Unmarshal([]byte(jsonStr), &node)
    if err != nil {
        t.Fatalf("sonic.Unmarshal err=%v", err)
    }

    bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
    label, err := parseCalculator(bc, &node)
    assert.NoError(t, err)
    assert.Equal(t, "calc1", label)
    _, ok := bc.mapConfig["calc1"]
    assert.True(t, ok)
}
