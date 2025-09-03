package lowcode

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/stats2"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// fakePlugin implements sgc7plugin.IPlugin for deterministic tests
type fakePluginWB struct{}

func (p *fakePluginWB) Random(_ context.Context, r int) (int, error) {
	if r <= 0 {
		return 0, nil
	}
	return 0, nil
}
func (p *fakePluginWB) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginWB) ClearUsedRngs()                    {}
func (p *fakePluginWB) TagUsedRngs()                      {}
func (p *fakePluginWB) RollbackUsedRngs() error           { return nil }
func (p *fakePluginWB) SetCache(arr []int)                {}
func (p *fakePluginWB) ClearCache()                       {}
func (p *fakePluginWB) Init()                             {}
func (p *fakePluginWB) SetScenePool(any)                  {}
func (p *fakePluginWB) GetScenePool() any                 { return nil }
func (p *fakePluginWB) SetSeed(seed int)                  {}

// plugin that returns error from Random
type fakePluginErrWB struct{}

func (p *fakePluginErrWB) Random(_ context.Context, r int) (int, error) {
	return 0, fmt.Errorf("rand err")
}
func (p *fakePluginErrWB) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginErrWB) ClearUsedRngs()                    {}
func (p *fakePluginErrWB) TagUsedRngs()                      {}
func (p *fakePluginErrWB) RollbackUsedRngs() error           { return nil }
func (p *fakePluginErrWB) SetCache(arr []int)                {}
func (p *fakePluginErrWB) ClearCache()                       {}
func (p *fakePluginErrWB) Init()                             {}
func (p *fakePluginErrWB) SetScenePool(any)                  {}
func (p *fakePluginErrWB) GetScenePool() any                 { return nil }
func (p *fakePluginErrWB) SetSeed(seed int)                  {}

// minimal stub RNG/featureLevel used by GameProperty in tests
type stubRNGWB struct{}

func (s *stubRNGWB) Clone() IRNG                                            { return &stubRNGWB{} }
func (s *stubRNGWB) OnNewGame(betMode int, plugin sgc7plugin.IPlugin) error { return nil }
func (s *stubRNGWB) GetCurRNG(betMode int, gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string) {
	return false, 0, nil, ""
}
func (s *stubRNGWB) OnChoiceBranch(betMode int, curComponent IComponent, branchName string) error {
	return nil
}
func (s *stubRNGWB) OnStepEnd(betMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) error {
	return nil
}

type stubFeatureLevelWB struct{}

func (s *stubFeatureLevelWB) Init() {}
func (s *stubFeatureLevelWB) OnStepEnd(gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult) {
}
func (s *stubFeatureLevelWB) CountLevel() int { return 0 }

func prepareVW2(name string) *sgc7game.ValWeights2 {
	sv := sgc7game.NewStrValEx(name)
	vw2, _ := sgc7game.NewValWeights2([]sgc7game.IVal{sv}, []int{1})
	return vw2
}

func prepareVW2Multi(names ...string) *sgc7game.ValWeights2 {
	vals := make([]sgc7game.IVal, 0, len(names))
	weights := make([]int, 0, len(names))
	for _, n := range names {
		vals = append(vals, sgc7game.NewStrValEx(n))
		weights = append(weights, 1)
	}
	vw2, _ := sgc7game.NewValWeights2(vals, weights)
	return vw2
}

func TestWeightBranchData_And_Init(t *testing.T) {
	// test WeightBranchData methods
	wbd := &WeightBranchData{}
	wbd.WeightVW = prepareVW2("x")
	wbd.IgnoreBranches = []string{"x"}

	gp := &GameProperty{Pool: &GamePropertyPool{}}
	comp := NewWeightBranch("twb")
	wbd.OnNewGame(gp, comp)
	if wbd.WeightVW != nil || wbd.IgnoreBranches != nil {
		t.Fatalf("OnNewGame should clear WeightVW and IgnoreBranches")
	}

	// Clone should copy Value
	wbd.Value = "v"
	c := wbd.Clone().(*WeightBranchData)
	if c.Value != "v" {
		t.Fatalf("Clone should copy Value")
	}

	// BuildPBComponentData should return non-nil proto
	pb := wbd.BuildPBComponentData()
	if pb == nil {
		t.Fatalf("BuildPBComponentData nil")
	}

	// test Init with yaml file
	tmpf := "test_weightbranch_init.yaml"
	yaml := "weight: w1\n"
	_ = os.WriteFile(tmpf, []byte(yaml), 0644)
	defer os.Remove(tmpf)

	vw := prepareVW2("w1v")
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw}, mapUsedWeights: make(map[string]string)}

	comp2 := NewWeightBranch("wbinit").(*WeightBranch)
	err := comp2.Init(tmpf, pool)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if comp2.Config.WeightVW == nil {
		t.Fatalf("Init should set WeightVW")
	}
}

func TestOnBranch_SuccessPaths(t *testing.T) {
	// case: wbd.WeightVW == nil path
	vw := prepareVW2Multi("a", "b", "c")
	comp := NewWeightBranch("obr").(*WeightBranch)
	comp.Config = &WeightBranchConfig{ForceTriggerOnce: []string{"a"}}

	wbd := &WeightBranchData{}
	err := comp.onBranch("a", wbd, vw)
	if err != nil {
		t.Fatalf("onBranch should succeed for first-time a: %v", err)
	}
	if wbd.WeightVW == nil {
		t.Fatalf("onBranch should set WeightVW")
	}

	// case: wbd.WeightVW != nil and CloneExcludeVal path
	// prepare vw with two vals
	vw2 := prepareVW2Multi("x", "y")
	wbd2 := &WeightBranchData{WeightVW: vw2}
	comp.Config.ForceTriggerOnce = []string{"y"}
	err = comp.onBranch("y", wbd2, vw2)
	if err != nil {
		t.Fatalf("onBranch CloneExcludeVal failed: %v", err)
	}
	// after excluding y, GetWeight(y) should be 0
	if wbd2.WeightVW.GetWeight(sgc7game.NewStrValEx("y")) != 0 {
		t.Fatalf("WeightVW should not contain y")
	}
}

func TestOnPlayGame_ForceBranch(t *testing.T) {
	vw := prepareVW2("fb")
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw}, mapUsedWeights: make(map[string]string)}

	gameProp := &GameProperty{Pool: pool}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)
	gameProp.rng = &stubRNGWB{}
	gameProp.featureLevel = &stubFeatureLevelWB{}

	comp := NewWeightBranch("wbf").(*WeightBranch)
	comp.Config = &WeightBranchConfig{ForceBranch: "fb", MapBranchs: map[string]*BranchNode{"fb": {JumpToComponent: "jx"}}}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	icd := comp.NewComponentData().(*WeightBranchData)
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
	plugin := &fakePluginWB{}

	nc, err := comp.OnPlayGame(gameProp, pr, gp, plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	if err != nil {
		t.Fatalf("OnPlayGame forceBranch failed: %v", err)
	}
	if nc != "jx" {
		t.Fatalf("expected next jx, got %s", nc)
	}
}

func TestAscii_And_LinkComponents_Stats(t *testing.T) {
	comp := NewWeightBranch("wascii").(*WeightBranch)
	comp.Config = &WeightBranchConfig{MapBranchs: map[string]*BranchNode{"k": {JumpToComponent: "n1"}}}

	// valid ascii
	icd := comp.NewComponentData().(*WeightBranchData)
	icd.Value = "v1"
	err := comp.OnAsciiGame(&GameProperty{}, nil, nil, nil, icd)
	if err != nil {
		t.Fatalf("OnAsciiGame valid failed: %v", err)
	}

	// link components
	lst := comp.GetAllLinkComponents()
	if len(lst) == 0 {
		t.Fatalf("GetAllLinkComponents empty")
	}
	lst2 := comp.GetNextLinkComponents()
	if len(lst2) == 0 {
		t.Fatalf("GetNextLinkComponents empty")
	}

	// stats
	s2 := stats2.NewCache(1)
	f := stats2.NewFeature(comp.GetName(), []stats2.Option{stats2.OptStrVal})
	s2.AddFeature(comp.GetName(), f, false)
	comp.OnStats2(icd, s2, &GameProperty{}, nil, nil, false)
}

func TestWeightBranch_Getters_And_ConfigMods(t *testing.T) {
	// GetValEx and GetStrVal
	wbd := &WeightBranchData{Value: "valx"}
	v, ok := wbd.GetValEx("any", 0)
	if ok || v != 0 {
		t.Fatalf("GetValEx unexpected")
	}
	sv, sok := wbd.GetStrVal(CSVValue)
	if !sok || sv != "valx" {
		t.Fatalf("GetStrVal unexpected: %v %v", sv, sok)
	}

	// SetConfigVal CCVWeight clears WeightVW
	wbd.MapConfigVals = make(map[string]string)
	wbd.WeightVW = prepareVW2("zzz")
	wbd.SetConfigVal(CCVWeight, "abc")
	if wbd.WeightVW != nil {
		t.Fatalf("SetConfigVal should clear WeightVW")
	}

	// SetConfigIntVal CCVClearForceTriggerOnceCache clears caches
	wbd.WeightVW = prepareVW2("zzz")
	wbd.IgnoreBranches = []string{"a"}
	wbd.SetConfigIntVal(CCVClearForceTriggerOnceCache, 1)
	if wbd.WeightVW != nil || wbd.IgnoreBranches != nil {
		t.Fatalf("SetConfigIntVal should clear caches")
	}

	// ChgConfigIntVal special handling
	wbd.IgnoreBranches = []string{"a"}
	v2 := wbd.ChgConfigIntVal(CCVClearForceTriggerOnceCache, 1)
	if v2 != 0 || wbd.IgnoreBranches != nil {
		t.Fatalf("ChgConfigIntVal did not clear as expected")
	}

	// SetLinkComponent non-next
	cfg := &WeightBranchConfig{}
	cfg.SetLinkComponent("abc", "cmpname")
	if cfg.MapBranchs["abc"].JumpToComponent != "cmpname" {
		t.Fatalf("SetLinkComponent failed")
	}

	// getWeight reading from component data config
	vw := prepareVW2("g1")
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"g1": vw}}
	gp := &GameProperty{Pool: pool}
	comp := NewWeightBranch("gw").(*WeightBranch)
	comp.Config = &WeightBranchConfig{}
	cd := &WeightBranchData{}
	cd.MapConfigVals = map[string]string{CCVWeight: "g1"}
	got := comp.getWeight(gp, cd)
	if got != vw {
		t.Fatalf("getWeight did not return expected vw")
	}

	// when cd.WeightVW is set, getWeight should return it directly
	cd2 := &WeightBranchData{WeightVW: prepareVW2("direct")}
	got2 := comp.getWeight(gp, cd2)
	if got2 != cd2.WeightVW {
		t.Fatalf("getWeight did not return cd.WeightVW")
	}

	// NewStats2 returns non-nil
	if comp.NewStats2("p") == nil {
		t.Fatalf("NewStats2 nil")
	}
}

func TestWeightBranch_ForceBranch_FromComponentData(t *testing.T) {
	comp := NewWeightBranch("fbcomp").(*WeightBranch)
	comp.Config = &WeightBranchConfig{ForceBranch: "cfgfb", MapBranchs: map[string]*BranchNode{"cfgfb": {JumpToComponent: "cj"}}}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	icd := comp.NewComponentData().(*WeightBranchData)
	icd.MapConfigVals = map[string]string{CCVForceBranch: "datfb"}
	// when component data has CCVForceBranch it should override config ForceBranch
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
	gameProp := &GameProperty{Pool: &GamePropertyPool{}}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)
	gameProp.rng = &stubRNGWB{}
	gameProp.featureLevel = &stubFeatureLevelWB{}

	// without MapBranch for datfb, next should be default
	nc, err := comp.OnPlayGame(gameProp, pr, gp, &fakePluginWB{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)
	// since datfb not in MapBranchs, next is empty
	assert.Equal(t, "", nc)

	// now add MapBranch for datfb
	comp.Config.MapBranchs["datfb"] = &BranchNode{JumpToComponent: "tj"}
	nc2, err := comp.OnPlayGame(gameProp, pr, gp, &fakePluginWB{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)
	assert.Equal(t, "tj", nc2)
}

func TestProcControllers_WithAwards_NoPanic(t *testing.T) {
	comp := NewWeightBranch("pc").(*WeightBranch)
	comp.Config = &WeightBranchConfig{MapBranchs: map[string]*BranchNode{"a": {Awards: []*Award{{AwardType: "respinTimes"}}}}}
	// gameProp minimal
	gp := &GameProperty{Pool: &GamePropertyPool{}}
	pr := &sgc7game.PlayResult{}
	// should not panic
	comp.ProcControllers(gp, &fakePluginWB{}, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), -1, "a")
}

func TestOnPlayGame_RandVal_ErrorPath(t *testing.T) {
	vw := prepareVW2Multi("e1", "e2")
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"we": vw}}

	gameProp := &GameProperty{Pool: pool}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)
	gameProp.rng = &stubRNGWB{}
	gameProp.featureLevel = &stubFeatureLevelWB{}

	comp := NewWeightBranch("we").(*WeightBranch)
	comp.Config = &WeightBranchConfig{Weight: "we", MapBranchs: map[string]*BranchNode{"e1": {}, "e2": {}}, WeightVW: vw}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	icd := comp.NewComponentData().(*WeightBranchData)
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	// plugin that errors on Random
	perr := &fakePluginErrWB{}

	_, err := comp.OnPlayGame(gameProp, pr, gp, perr, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	if err == nil {
		t.Fatalf("expected error from RandVal path when plugin.Random fails")
	}
}

func TestSetLinkComponent_Next(t *testing.T) {
	cfg := &WeightBranchConfig{}
	cfg.SetLinkComponent("next", "ncomp")
	if cfg.DefaultNextComponent != "ncomp" {
		t.Fatalf("SetLinkComponent next failed")
	}
}

func TestInitEx_NoWeight_Error(t *testing.T) {
	comp := NewWeightBranch("initerr").(*WeightBranch)
	// config without weight should cause InitEx to return ErrInvalidComponentConfig
	cfg := &WeightBranchConfig{}
	_, _ = comp, cfg
	_, pool := comp, &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{}}
	err := comp.InitEx(cfg, pool)
	if err == nil {
		t.Fatalf("InitEx should return error when weight not set")
	}
}

func TestSetLinkComponent_UpdateExisting(t *testing.T) {
	cfg := &WeightBranchConfig{MapBranchs: map[string]*BranchNode{"x": {JumpToComponent: "old"}}}
	cfg.SetLinkComponent("x", "new")
	if cfg.MapBranchs["x"].JumpToComponent != "new" {
		t.Fatalf("SetLinkComponent did not update existing entry")
	}
}

func TestInitEx_WithAwardsCallsInit(t *testing.T) {
	vw := prepareVW2("w1v")
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw}, mapUsedWeights: make(map[string]string)}

	comp := NewWeightBranch("initaw").(*WeightBranch)
	cfg := &WeightBranchConfig{Weight: "w1", MapBranchs: map[string]*BranchNode{"a": {Awards: []*Award{{AwardType: "respinTimes"}}}}}

	err := comp.InitEx(cfg, pool)
	if err != nil {
		t.Fatalf("InitEx with awards failed: %v", err)
	}
	// ensure awards got their Type initialized
	if comp.Config.MapBranchs["a"].Awards[0].Type != AwardRespinTimes {
		t.Fatalf("award.Init did not set Type")
	}
}

func TestWeightBranch_OnPlayGame_RandAndForceOnce(t *testing.T) {
	// prepare ValWeights2
	vw2 := prepareVW2("b1")

	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw2}, mapUsedWeights: make(map[string]string)}

	gameProp := &GameProperty{Pool: pool}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)
	gameProp.rng = &stubRNGWB{}
	gameProp.featureLevel = &stubFeatureLevelWB{}

	// component config using weight "w1"
	comp := NewWeightBranch("wb").(*WeightBranch)
	comp.Config = &WeightBranchConfig{Weight: "w1", MapBranchs: map[string]*BranchNode{"b1": {JumpToComponent: "next"}}, ForceTriggerOnce: []string{"b1"}}
	// set preloaded ValWeights to simulate InitEx behavior
	comp.Config.WeightVW = vw2
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	// icd
	icd := comp.NewComponentData().(*WeightBranchData)

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
	plugin := &fakePluginWB{}

	// first play: should pick b1 and set IgnoreBranches
	nc, err := comp.OnPlayGame(gameProp, pr, gp, plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)
	assert.Equal(t, "next", nc)
	assert.Equal(t, "b1", icd.Value)
	// since ForceTriggerOnce contains b1, it should be in IgnoreBranches
	assert.Contains(t, icd.IgnoreBranches, "b1")

	// second play: since b1 is ignored, the weight list for this data should not contain b1
	icd2 := icd.Clone().(*WeightBranchData)
	icd2.OnNewGame(gameProp, comp) // clear per-game via OnNewGame
	// manually set IgnoreBranches to simulate previous pick
	icd2.IgnoreBranches = []string{"b1"}
	// now call onBranch with b1 should return ErrInvalidBranch when called again
	err = comp.onBranch("b1", icd2, vw2)
	assert.Error(t, err)
}

func TestWeightBranch_OnPlayGame_PlayerSelect(t *testing.T) {
	vw2 := prepareVW2("b2")
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w2": vw2}}

	gameProp := &GameProperty{Pool: pool}
	gameProp.rng = &stubRNGWB{}
	gameProp.featureLevel = &stubFeatureLevelWB{}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)

	comp := NewWeightBranch("wb2").(*WeightBranch)
	comp.Config = &WeightBranchConfig{Weight: "w2", MapBranchs: map[string]*BranchNode{"b2": {JumpToComponent: "nx"}}, IsNeedPlayerSelect: true}
	comp.Config.WeightVW = vw2
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	icd := comp.NewComponentData().(*WeightBranchData)

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
	plugin := &fakePluginWB{}

	// call with DefaultCmd should populate NextCmds and NextCmdParams and return no next component
	nc, err := comp.OnPlayGame(gameProp, pr, gp, plugin, DefaultCmd, "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)
	assert.Equal(t, "", nc)
	assert.False(t, pr.IsFinish)
	assert.True(t, pr.IsWait)

	// now simulate player selecting invalid command
	_, err = comp.OnPlayGame(gameProp, pr, gp, plugin, "WRONG", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.Error(t, err)

	// simulate correct selection
	nc2, err := comp.OnPlayGame(gameProp, pr, gp, plugin, comp.Name, "b2", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)
	// should jump to configured next
	assert.Equal(t, "nx", nc2)
}

func TestWeightBranch_GetWeight_Missing(t *testing.T) {
	// missing weight should cause getWeight to try pool and return nil -> OnPlayGame returns ErrInvalidComponentConfig
	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{}}
	gameProp := &GameProperty{Pool: pool}
	gameProp.rng = &stubRNGWB{}
	gameProp.featureLevel = &stubFeatureLevelWB{}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)

	comp := NewWeightBranch("wb3").(*WeightBranch)
	comp.Config = &WeightBranchConfig{Weight: "wmissing", MapBranchs: map[string]*BranchNode{}}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	icd := comp.NewComponentData().(*WeightBranchData)

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
	plugin := &fakePluginWB{}

	_, err := comp.OnPlayGame(gameProp, pr, gp, plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.Error(t, err)
}

func TestWeightBranch_OnAsciiGame_And_OnStats2_InvalidICD(t *testing.T) {
	comp := NewWeightBranch("aascii").(*WeightBranch)
	// OnAsciiGame with invalid icd should return ErrInvalidComponentData
	err := comp.OnAsciiGame(&GameProperty{}, nil, nil, nil, nil)
	assert.Error(t, err)

	// OnStats2 with invalid icd should not panic; pass nil icd to hit error logging path
	s2 := stats2.NewCache(1)
	comp.OnStats2(nil, s2, &GameProperty{}, nil, nil, false)
	// ensure no panic and stats map remains empty
	if len(s2.MapStats) != 0 {
		t.Fatalf("expected empty stats map")
	}
}

func TestParseWeightBranch_ParseErr(t *testing.T) {
	// nil cell should return error
	_, err := parseWeightBranch(&BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}, nil)
	assert.Error(t, err)

	// build a valid JSON node and ensure parse succeeds
	jsonStr := `{"componentValues": {"label": "wlabel","configuration": {"weight": "w1","isNeedPlayerSelect": false},"controller": []}}`
	var node ast.Node
	err = sonic.Unmarshal([]byte(jsonStr), &node)
	assert.NoError(t, err)

	bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
	label, err := parseWeightBranch(bc, &node)
	assert.NoError(t, err)
	assert.Equal(t, "wlabel", label)
	_, ok := bc.mapConfig["wlabel"]
	assert.True(t, ok)
}

func TestParseWeightBranch_WithControllersSuccess(t *testing.T) {
	// build a JSON node with controller that maps a stringVal to an award
	jsonStr := `{"componentValues": {"label": "wlabel3","configuration": {"weight": "w1","isNeedPlayerSelect": false},"controller": [{"type":"addRespinTimes","stringVal":"branchA","times":2}]}}`
	var node ast.Node
	err := sonic.Unmarshal([]byte(jsonStr), &node)
	assert.NoError(t, err)

	bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
	label, err := parseWeightBranch(bc, &node)
	assert.NoError(t, err)
	assert.Equal(t, "wlabel3", label)

	// ensure the controller populated MapBranchs with branchA awards
	ccfg, ok := bc.mapConfig[label]
	assert.True(t, ok)
	cfg := ccfg.(*WeightBranchConfig)
	bn, ok2 := cfg.MapBranchs["branchA"]
	assert.True(t, ok2)
	assert.NotEmpty(t, bn.Awards)
}

func TestParseWeightBranch_WithControllersFail(t *testing.T) {
	// controller with unsupported type should cause parse failure
	jsonStr := `{"componentValues": {"label": "wlabel4","configuration": {"weight": "w1"},"controller": [{"type":"unknownType","stringVal":"x"}]}}`
	var node ast.Node
	err := sonic.Unmarshal([]byte(jsonStr), &node)
	assert.NoError(t, err)

	bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
	_, err = parseWeightBranch(bc, &node)
	assert.Error(t, err)
}

func TestConfigIntVal_ChgConfigIntVal_OtherKeys(t *testing.T) {
	wbd := &WeightBranchData{}
	// initialize maps
	wbd.MapConfigIntVals = make(map[string]int)

	// SetConfigIntVal with other key should set the int value
	wbd.SetConfigIntVal("other", 7)
	if v, ok := wbd.MapConfigIntVals["other"]; !ok || v != 7 {
		t.Fatalf("SetConfigIntVal did not set value, got %v %v", v, ok)
	}

	// ChgConfigIntVal with other key should add offset
	v2 := wbd.ChgConfigIntVal("other", 3)
	if v2 != 10 {
		t.Fatalf("ChgConfigIntVal unexpected, got %d", v2)
	}
}

func TestGetStrVal_Negative(t *testing.T) {
	wbd := &WeightBranchData{Value: "vv"}
	s, ok := wbd.GetStrVal("nope")
	if ok || s != "" {
		t.Fatalf("GetStrVal negative unexpected: %v %v", s, ok)
	}
}

func TestInit_FileReadError(t *testing.T) {
	comp := NewWeightBranch("initbad").(*WeightBranch)
	// call Init with non-existent file
	pool := &GamePropertyPool{mapStrValWeights: make(map[string]*sgc7game.ValWeights2)}
	err := comp.Init("this_file_does_not_exist.yaml", pool)
	if err == nil {
		t.Fatalf("Init should fail for missing file")
	}
}
