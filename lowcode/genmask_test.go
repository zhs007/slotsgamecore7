package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// helper to build minimal mask component and game context
func newTestMaskComponent(name string, num int, ignoreFalse bool) (*Mask, *GameProperty, sgc7plugin.IPlugin, *GameParams, *sgc7game.PlayResult) {
	// build pool/config
	cfg := &Config{Width: num, Height: 1, Bets: []int{1}, MapBetConfigs: map[int]*BetConfig{1: {Bet: 1, Start: name}}}
	pool := &GamePropertyPool{Config: cfg, mapComponents: make(map[int]*ComponentList)}
	pool.DefaultPaytables = &sgc7game.PayTables{MapSymbols: map[string]int{"A": 1}}

	// register components list
	lst := NewComponentList()
	pool.onAddComponentList(1, lst)

	// create mask component and add
	mask := &Mask{BasicComponent: NewBasicComponent(name, 1), Config: &MaskConfig{Num: num, IgnoreFalse: ignoreFalse}}
	lst.AddComponent(name, mask)

	// init game property
	pool.newRNG = NewBasicRNG
	pool.newFeatureLevel = NewEmptyFeatureLevel
	gp := pool.newGameProp(1)
	gp.Components = lst
	_ = pool.InitStats(1)

	plugin := sgc7plugin.NewFastPlugin()
	pr := &sgc7game.PlayResult{}
	stake := &sgc7game.Stake{CashBet: 1, CoinBet: 1}
	_ = gp.OnNewGame(stake, plugin)
	_ = gp.OnNewStep()
	gparams := NewGameParam(stake, NewPlayerState())

	return mask, gp, plugin, gparams, pr
}

func Test_GenMask_SetAndNot(t *testing.T) {
	// prepare a source mask component
	srcMask, gp, plugin, gpms, pr := newTestMaskComponent("src", 6, false)
	// set some mask values
	_ = gp.Pool.SetMask(plugin, gp, pr, gpms, srcMask.GetName(), []bool{false, true, true, false, true, false}, false)

	// genMask set
	gm := &GenMask{BasicComponent: NewBasicComponent("gm", 0), Config: &GenMaskConfig{StrType: "set", MaskLen: 6, OutputMask: "out", SrcMask: []string{"src"}}}
	gm.InitEx(gm.Config, gp.Pool)
	gp.Components.AddComponent("gm", gm)

	// add output mask to component list
	outMask := &Mask{BasicComponent: NewBasicComponent("out", 1), Config: &MaskConfig{Num: 6}}
	gp.Components.AddComponent("out", outMask)

	// run set
	next, err := gm.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CashBet: 1, CoinBet: 1}, NewPlayerState()), plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	_ = next

	// verify copied
	mv, _ := gp.GetMask("out")
	assert.Equal(t, []bool{false, true, true, false, true, false}, mv)

	// not
	gm.Config.StrType = "not"
	gm.Config.Type = parseGenMaskType("not")
	_, err = gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv2, _ := gp.GetMask("out")
	assert.Equal(t, []bool{true, false, false, true, false, true}, mv2)
}

func Test_GenMask_Random_NoSrc(t *testing.T) {
	_, gp, plugin, gpms, pr := newTestMaskComponent("src2", 5, false)

	gm := &GenMask{BasicComponent: NewBasicComponent("gm2", 0), Config: &GenMaskConfig{StrType: "random", MaskLen: 5, OutputMask: "out2", WeightValue: 10000}}
	gm.InitEx(gm.Config, gp.Pool)
	gp.Components.AddComponent("gm2", gm)

	outMask := &Mask{BasicComponent: NewBasicComponent("out2", 1), Config: &MaskConfig{Num: 5}}
	gp.Components.AddComponent("out2", outMask)

	_, err := gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ := gp.GetMask("out2")
	// weight=10000 should produce all true
	assert.Equal(t, []bool{true, true, true, true, true}, mv)
}

func Test_GenMask_AndOrXor(t *testing.T) {
	// two source masks
	src1, gp, plugin, gpms, pr := newTestMaskComponent("s1", 4, false)
	_ = gp.Pool.SetMask(plugin, gp, pr, gpms, src1.GetName(), []bool{true, false, true, false}, false)
	src2 := &Mask{BasicComponent: NewBasicComponent("s2", 1), Config: &MaskConfig{Num: 4}}
	gp.Components.AddComponent("s2", src2)
	_ = gp.Pool.SetMask(plugin, gp, pr, gpms, src2.GetName(), []bool{true, true, false, false}, false)

	out := &Mask{BasicComponent: NewBasicComponent("out3", 1), Config: &MaskConfig{Num: 4}}
	gp.Components.AddComponent("out3", out)

	gm := &GenMask{BasicComponent: NewBasicComponent("gm3", 0), Config: &GenMaskConfig{StrType: "and", MaskLen: 4, OutputMask: "out3", SrcMask: []string{"s1", "s2"}}}
	gm.InitEx(gm.Config, gp.Pool)
	gp.Components.AddComponent("gm3", gm)

	// AND
	_, err := gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ := gp.GetMask("out3")
	assert.Equal(t, []bool{true, false, false, false}, mv)

	// OR
	gm.Config.StrType = "or"
	gm.Config.Type = parseGenMaskType("or")
	_, err = gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ = gp.GetMask("out3")
	assert.Equal(t, []bool{true, true, true, false}, mv)

	// XOR
	gm.Config.StrType = "xor"
	gm.Config.Type = parseGenMaskType("xor")
	_, err = gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ = gp.GetMask("out3")
	assert.Equal(t, []bool{false, true, true, false}, mv)
}

func Test_GenMask_LengthMismatch(t *testing.T) {
	src, gp, plugin, gpms, pr := newTestMaskComponent("slen", 3, false)
	_ = gp.Pool.SetMask(plugin, gp, pr, gpms, src.GetName(), []bool{true, false, true}, false)

	out := &Mask{BasicComponent: NewBasicComponent("outlen", 1), Config: &MaskConfig{Num: 4}}
	gp.Components.AddComponent("outlen", out)

	gm := &GenMask{BasicComponent: NewBasicComponent("gmlen", 0), Config: &GenMaskConfig{StrType: "and", MaskLen: 4, OutputMask: "outlen", SrcMask: []string{"slen"}}}
	gm.InitEx(gm.Config, gp.Pool)

	// should error due to length mismatch in getAllMask
	_, err := gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.Error(t, err)
}

func Test_GenMask_Random_WithSrc(t *testing.T) {
	// src has some false; random should only affect true positions
	src, gp, plugin, gpms, pr := newTestMaskComponent("srcrand", 6, false)
	_ = gp.Pool.SetMask(plugin, gp, pr, gpms, src.GetName(), []bool{false, true, false, true, false, true}, false)

	out := &Mask{BasicComponent: NewBasicComponent("outrand", 1), Config: &MaskConfig{Num: 6}}
	gp.Components.AddComponent("outrand", out)

	gm := &GenMask{BasicComponent: NewBasicComponent("gmr", 0), Config: &GenMaskConfig{StrType: "random", MaskLen: 6, OutputMask: "outrand", SrcMask: []string{"srcrand"}, WeightValue: 0}}
	gm.InitEx(gm.Config, gp.Pool)
	gp.Components.AddComponent("gmr", gm)

	// WeightValue=0 should make all false, but only where src is true it applies; others stay false
	_, err := gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ := gp.GetMask("outrand")
	assert.Equal(t, []bool{false, false, false, false, false, false}, mv)

	// now set WeightValue=10000 to force all true on true positions; false positions remain false
	gm.Config.WeightValue = 10000
	_, err = gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ = gp.GetMask("outrand")
	assert.Equal(t, []bool{false, true, false, true, false, true}, mv)
}

func Test_GenMask_YamlInit(t *testing.T) {
	// simulate YAML config by using InitEx (we avoid file IO in unit tests)
	_, gp, plugin, gpms, pr := newTestMaskComponent("srcyaml", 3, false)

	out := &Mask{BasicComponent: NewBasicComponent("outyaml", 1), Config: &MaskConfig{Num: 3}}
	gp.Components.AddComponent("outyaml", out)

	cfg := &GenMaskConfig{StrType: "set", MaskLen: 3, OutputMask: "outyaml", InitMask: []bool{true, false, true}}
	gm := &GenMask{BasicComponent: NewBasicComponent("gmyaml", 0)}
	err := gm.InitEx(cfg, gp.Pool)
	assert.NoError(t, err)
	gp.Components.AddComponent("gmyaml", gm)

	_, err = gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ := gp.GetMask("outyaml")
	assert.Equal(t, []bool{true, false, true}, mv)
}

func Test_GenMask_InitError_Set_NoSrcNoInit(t *testing.T) {
	// set 类型：没有 srcMask 且没有 initMask，应在 InitEx 阶段报错
	_, gp, _, _, _ := newTestMaskComponent("dummy", 3, false)
	gm := &GenMask{BasicComponent: NewBasicComponent("gmerr1", 0)}
	cfg := &GenMaskConfig{StrType: "set", MaskLen: 3, OutputMask: "out"}
	err := gm.InitEx(cfg, gp.Pool)
	assert.Error(t, err)
}

func Test_GenMask_InitError_Random_BadWeight(t *testing.T) {
	// random 类型：weightValue 越界应在 InitEx 阶段报错
	_, gp, _, _, _ := newTestMaskComponent("dummy2", 3, false)
	gm := &GenMask{BasicComponent: NewBasicComponent("gmerr2", 0)}
	cfg := &GenMaskConfig{StrType: "random", MaskLen: 3, OutputMask: "out", WeightValue: -1}
	err := gm.InitEx(cfg, gp.Pool)
	assert.Error(t, err)

	cfg.WeightValue = 10001
	err = gm.InitEx(cfg, gp.Pool)
	assert.Error(t, err)
}

func Test_GenMask_InitError_And_TooFewMasks(t *testing.T) {
	// and 类型：没有 initMask 且 srcMask 少于 2，应在 InitEx 报错
	_, gp, _, _, _ := newTestMaskComponent("dummy3", 3, false)
	gm := &GenMask{BasicComponent: NewBasicComponent("gmerr3", 0)}
	cfg := &GenMaskConfig{StrType: "and", MaskLen: 3, OutputMask: "out", SrcMask: []string{"onlyone"}}
	err := gm.InitEx(cfg, gp.Pool)
	assert.Error(t, err)
}

func Test_GenMask_Or_WithInitMask(t *testing.T) {
	// or 类型：有 initMask 且至少 1 个 srcMask，结果应为两者 OR
	src, gp, plugin, gpms, pr := newTestMaskComponent("s_or_init", 3, false)
	_ = gp.Pool.SetMask(plugin, gp, pr, gpms, src.GetName(), []bool{false, true, false}, false)

	out := &Mask{BasicComponent: NewBasicComponent("out_or_init", 1), Config: &MaskConfig{Num: 3}}
	gp.Components.AddComponent("out_or_init", out)

	cfg := &GenMaskConfig{StrType: "or", MaskLen: 3, OutputMask: "out_or_init", SrcMask: []string{"s_or_init"}, InitMask: []bool{true, false, true}}
	gm := &GenMask{BasicComponent: NewBasicComponent("gm_or_init", 0)}
	err := gm.InitEx(cfg, gp.Pool)
	assert.NoError(t, err)
	gp.Components.AddComponent("gm_or_init", gm)

	_, err = gm.OnPlayGame(gp, pr, gpms, plugin, "", "", nil, nil, nil, gm.BasicComponent.NewComponentData())
	assert.NoError(t, err)
	mv, _ := gp.GetMask("out_or_init")
	assert.Equal(t, []bool{true, true, true}, mv)
}
