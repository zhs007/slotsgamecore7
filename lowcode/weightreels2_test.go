package lowcode

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/stretchr/testify/assert"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// fakePlugin implements sgc7plugin.IPlugin for deterministic tests
type fakePlugin struct{}

func (p *fakePlugin) Random(_ context.Context, r int) (int, error) {
	if r <= 0 {
		return 0, nil
	}
	return 0, nil
}
func (p *fakePlugin) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePlugin) ClearUsedRngs()                    {}
func (p *fakePlugin) TagUsedRngs()                      {}
func (p *fakePlugin) RollbackUsedRngs() error           { return nil }
func (p *fakePlugin) SetCache(arr []int)                {}
func (p *fakePlugin) ClearCache()                       {}
func (p *fakePlugin) Init()                             {}
func (p *fakePlugin) SetScenePool(any)                  {}
func (p *fakePlugin) GetScenePool() any                 { return nil }
func (p *fakePlugin) SetSeed(seed int)                  {}

func TestWeightReels2DataBasics(t *testing.T) {
	wd := &WeightReels2Data{}

	// OnNewStep sets ReelSetIndex = -1 and clears UsedScenes
	wd.UsedScenes = []int{1, 2}
	wd.ReelSetIndex = 5
	wd.onNewStep()

	assert.Equal(t, -1, wd.ReelSetIndex)
	assert.Empty(t, wd.UsedScenes)

	// GetValEx for selected index
	wd.ReelSetIndex = 3
	v, ok := wd.GetValEx(CVSelectedIndex, GCVTypeNormal)
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	// Clone
	c := wd.Clone().(*WeightReels2Data)
	assert.Equal(t, wd.ReelSetIndex, c.ReelSetIndex)

	// BuildPBComponentData should return protobuf message
	pb := wd.BuildPBComponentData()
	assert.NotNil(t, pb)
}

// stub implementations for interfaces used by GameProperty
type stubRNG struct{}

func (s *stubRNG) Clone() IRNG                                            { return &stubRNG{} }
func (s *stubRNG) OnNewGame(betMode int, plugin sgc7plugin.IPlugin) error { return nil }
func (s *stubRNG) GetCurRNG(betMode int, gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string) {
	return false, 0, nil, ""
}
func (s *stubRNG) OnChoiceBranch(betMode int, curComponent IComponent, branchName string) error {
	return nil
}
func (s *stubRNG) OnStepEnd(betMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) error {
	return nil
}

type stubFeatureLevel struct{}

func (s *stubFeatureLevel) Init() {}
func (s *stubFeatureLevel) OnStepEnd(gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult) {
}
func (s *stubFeatureLevel) CountLevel() int { return 0 }

func TestWeightReels2_GetReelSetWeightAndOnPlayGame(t *testing.T) {
	// prepare ValWeights2 with one StrVal so RandValEx returns it deterministically
	sv := sgc7game.NewStrValEx("rname")
	vw2, err := sgc7game.NewValWeights2([]sgc7game.IVal{sv}, []int{1})
	assert.NoError(t, err)

	// prepare pool and config
	pool := &GamePropertyPool{
		mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw2},
	}

	cfg := &Config{
		MapReels: map[string]*sgc7game.ReelsData{"rname": {Reels: [][]int{{1, 2, 3}}}},
	}

	// create a minimal GameProperty and initialize PoolScene and RNG/featureLevel stubs
	gameProp := &GameProperty{Pool: pool}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	// stub rng and featureLevel to avoid nil calls
	gameProp.rng = &stubRNG{}
	gameProp.featureLevel = &stubFeatureLevel{}
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)
	pool.Config = cfg
	pool.MapSymbolColor = nil

	// component and basicCD
	comp := NewWeightReels2("wr2").(*WeightReels2)
	comp.Config = &WeightReels2Config{IsExpandReel: false}
	// ensure embedded BasicComponent has a config to avoid nil deref in onStepEnd
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	bcd := &BasicComponentData{MapConfigVals: map[string]string{}}

	// when basicCD has no CCVReelSetWeight, should return comp.Config.ReelSetsWeightVW (nil)
	assert.Nil(t, comp.getReelSetWeight(gameProp, bcd))

	// when basicCD has CCVReelSetWeight set, pool.LoadStrWeights should return vw2
	bcd.MapConfigVals = map[string]string{CCVReelSetWeight: "w1"}
	// attach Config.ReelSetsWeightVW nil
	comp.Config.ReelSetsWeightVW = nil

	// call OnPlayGame with component data that contains CCVReelSetWeight
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)
	plugin := &fakePlugin{}

	icd := comp.NewComponentData().(*WeightReels2Data)
	icd.MapConfigVals = map[string]string{CCVReelSetWeight: "w1"}

	nc, err := comp.OnPlayGame(gameProp, pr, gp, plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)
	// because default next component is empty, onStepEnd returns DefaultNextComponent which is empty -> nc == ""
	assert.Equal(t, "", nc)

	// after play, gameProp.CurReels should be set
	assert.NotNil(t, gameProp.CurReels)
}

func TestProcControllers_MapAwards(t *testing.T) {
	vw := NewWeightReels2("wr2").(*WeightReels2)
	// build a gameProp with procAwards invoked — we just ensure no panic when mapAwards contains key
	gp := &GameProperty{Pool: &GamePropertyPool{}}
	bp := &sgc7game.PlayResult{}
	params := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	// mapAwards with one award (empty awards array)
	vw.Config = &WeightReels2Config{MapAwards: map[string][]*Award{"k": {&Award{AwardType: "respinTimes", Vals: []int{1}, StrParams: []string{"x"}}}}}

	// should not panic
	vw.ProcControllers(gp, &fakePlugin{}, bp, params, -1, "k")
}

func TestOnAsciiGame_NoPanic(t *testing.T) {
	comp := NewWeightReels2("wr2").(*WeightReels2)
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	lst := []*sgc7game.PlayResult{}
	scm := asciigame.NewSymbolColorMap(&sgc7game.PayTables{})

	// OnAsciiGame with an empty WeightReels2Data should not panic
	err := comp.OnAsciiGame(&GameProperty{}, pr, lst, scm, comp.NewComponentData())
	assert.NoError(t, err)
}

func TestWeightReels2_InitExAndJSONBuild(t *testing.T) {
	// test json build
	j := &jsonWeightReels2{ReelSetWeight: "rw", IsExpandReel: true}
	cfg := j.build()
	assert.Equal(t, "rw", cfg.ReelSetsWeight)
	assert.True(t, cfg.IsExpandReel)

	// test SetLinkComponent
	cfg.SetLinkComponent("next", "nextcomp")
	assert.Equal(t, "nextcomp", cfg.DefaultNextComponent)

	// prepare ValWeights2 and pool
	sv := sgc7game.NewStrValEx("rname")
	vw2, err := sgc7game.NewValWeights2([]sgc7game.IVal{sv}, []int{1})
	assert.NoError(t, err)

	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw2}, mapUsedWeights: make(map[string]string)}

	// prepare WeightReels2 and call InitEx with ReelSetsWeight set
	comp := NewWeightReels2("wrx").(*WeightReels2)
	wcfg := &WeightReels2Config{ReelSetsWeight: "w1", MapAwards: map[string][]*Award{"k": {{AwardType: "respinTimes"}}}}

	err = comp.InitEx(wcfg, pool)
	assert.NoError(t, err)
	// after InitEx, ReelSetsWeightVW should be set from pool
	assert.Equal(t, vw2, comp.Config.ReelSetsWeightVW)
}

func TestWeightReels2Data_OnNewGameAndInit(t *testing.T) {
	// OnNewGame should initialize MapConfigVals
	cd := &WeightReels2Data{}
	gp := &GameProperty{Pool: &GamePropertyPool{}}
	comp := NewWeightReels2("wtest")

	cd.OnNewGame(gp, comp)
	// underlying BasicComponentData.OnNewGame should set maps
	bcd := &cd.BasicComponentData
	if bcd.MapConfigVals == nil {
		t.Fatalf("MapConfigVals should be initialized")
	}

	// test Init reads YAML file and calls InitEx
	tmpf := "test_weightreels2_init.yaml"
	yaml := "reelSetWeight: w1\n"
	err := os.WriteFile(tmpf, []byte(yaml), 0644)
	if err != nil {
		t.Fatalf("write temp file err=%v", err)
	}
	defer os.Remove(tmpf)

	// prepare pool with weight named w1
	sv := sgc7game.NewStrValEx("rname")
	vw2, err := sgc7game.NewValWeights2([]sgc7game.IVal{sv}, []int{1})
	assert.NoError(t, err)

	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw2}, mapUsedWeights: make(map[string]string)}

	wc := NewWeightReels2("wx").(*WeightReels2)
	err = wc.Init(tmpf, pool)
	assert.NoError(t, err)
	assert.Equal(t, vw2, wc.Config.ReelSetsWeightVW)
}

func TestWeightReels2_OnPlayGame_VW2Nil(t *testing.T) {
	// When Config.ReelSetsWeightVW is nil but gameProp.CurReels is provided
	comp := NewWeightReels2("wrnil").(*WeightReels2)
	comp.Config = &WeightReels2Config{IsExpandReel: false, ReelSetsWeightVW: nil}
	comp.BasicComponent.onInit(&BasicComponentConfig{DefaultNextComponent: ""})

	pool := &GamePropertyPool{}
	gp := &GameProperty{Pool: pool}
	gp.PoolScene = sgc7game.NewGameScenePoolEx()
	gp.SceneStack = NewSceneStack(false)
	gp.OtherSceneStack = NewSceneStack(true)
	gp.rng = &stubRNG{}
	gp.featureLevel = &stubFeatureLevel{}

	// provide CurReels so RandReelsWithReelData can run
	gp.CurReels = &sgc7game.ReelsData{Reels: [][]int{{1, 2, 3}}}

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp2 := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	plugin := &fakePlugin{}

	_, err := comp.OnPlayGame(gp, pr, gp2, plugin, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, comp.NewComponentData())
	assert.NoError(t, err)
}

func TestOnAsciiGame_WithUsedScenes(t *testing.T) {
	comp := NewWeightReels2("wrascii").(*WeightReels2)
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")

	sc, _ := sgc7game.NewGameScene(3, 3)
	pr.Scenes = append(pr.Scenes, sc)

	icd := comp.NewComponentData().(*WeightReels2Data)
	icd.UsedScenes = []int{0}

	scm := asciigame.NewSymbolColorMap(&sgc7game.PayTables{})

	err := comp.OnAsciiGame(&GameProperty{}, pr, nil, scm, icd)
	assert.NoError(t, err)
}

func TestWeightReels2Data_OnNewGameMethod(t *testing.T) {
	wrd := &WeightReels2Data{}
	gp := &GameProperty{Pool: &GamePropertyPool{}}
	comp := NewWeightReels2("wn")

	wrd.OnNewGame(gp, comp)
	if wrd.MapConfigVals == nil {
		t.Fatalf("MapConfigVals should be initialized by OnNewGame")
	}
}

// plugin that returns error from Random
type fakePluginErr struct{}

func (p *fakePluginErr) Random(_ context.Context, r int) (int, error) {
	return 0, fmt.Errorf("rand err")
}
func (p *fakePluginErr) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (p *fakePluginErr) ClearUsedRngs()                    {}
func (p *fakePluginErr) TagUsedRngs()                      {}
func (p *fakePluginErr) RollbackUsedRngs() error           { return nil }
func (p *fakePluginErr) SetCache(arr []int)                {}
func (p *fakePluginErr) ClearCache()                       {}
func (p *fakePluginErr) Init()                             {}
func (p *fakePluginErr) SetScenePool(any)                  {}
func (p *fakePluginErr) GetScenePool() any                 { return nil }
func (p *fakePluginErr) SetSeed(seed int)                  {}

func TestOnPlayGame_MapReelsMissingAndExpandAndRandErr(t *testing.T) {
	// MapReels missing -> ErrInvalidReels
	sv := sgc7game.NewStrValEx("nonexist")
	vw2, _ := sgc7game.NewValWeights2([]sgc7game.IVal{sv}, []int{1})

	pool := &GamePropertyPool{mapStrValWeights: map[string]*sgc7game.ValWeights2{"w1": vw2}}
	cfg := &Config{MapReels: map[string]*sgc7game.ReelsData{}}
	gameProp := &GameProperty{Pool: pool}
	gameProp.PoolScene = sgc7game.NewGameScenePoolEx()
	gameProp.SceneStack = NewSceneStack(false)
	gameProp.OtherSceneStack = NewSceneStack(true)
	gameProp.rng = &stubRNG{}
	gameProp.featureLevel = &stubFeatureLevel{}
	pool.Config = cfg

	comp := NewWeightReels2("wrx").(*WeightReels2)
	comp.Config = &WeightReels2Config{IsExpandReel: false}
	comp.BasicComponent.onInit(&BasicComponentConfig{})

	icd := comp.NewComponentData().(*WeightReels2Data)
	icd.MapConfigVals = map[string]string{CCVReelSetWeight: "w1"}

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gp := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	_, err := comp.OnPlayGame(gameProp, pr, gp, &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	if err == nil {
		t.Fatalf("expected ErrInvalidReels when MapReels missing")
	}

	// Expand reel branch: provide valid reels and set IsExpandReel true
	cfg.MapReels = map[string]*sgc7game.ReelsData{"nonexist": {Reels: [][]int{{1, 2, 3}}}}
	sv2a := sgc7game.NewStrValEx("nonexist")
	sv2b := sgc7game.NewStrValEx("other")
	vw3, _ := sgc7game.NewValWeights2([]sgc7game.IVal{sv2a, sv2b}, []int{1, 1})
	pool.mapStrValWeights["w2"] = vw3

	comp.Config.IsExpandReel = true
	icd.MapConfigVals = map[string]string{CCVReelSetWeight: "w2"}
	_, err = comp.OnPlayGame(gameProp, pr, gp, &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	assert.NoError(t, err)

	// RandValEx error path: plugin returns error
	icd.MapConfigVals = map[string]string{CCVReelSetWeight: "w2"}
	_, err = comp.OnPlayGame(gameProp, pr, gp, &fakePluginErr{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, icd)
	if err == nil {
		t.Fatalf("expected error from RandValEx when plugin.Random returns error")
	}
}

func TestParseWeightReels2_Errors(t *testing.T) {
	// call parseWeightReels2 with nil cell to hit error path
	_, err := parseWeightReels2(&BetConfig{mapConfig: make(map[string]IComponentConfig)}, nil)
	if err == nil {
		t.Fatalf("expected error when parsing nil cell")
	}

	// call Init with non-existent file to hit ReadFile error path
	comp := NewWeightReels2("wmissing").(*WeightReels2)
	err = comp.Init("/non/existent/file.yaml", &GamePropertyPool{})
	if err == nil {
		t.Fatalf("expected error when Init with missing file")
	}
}

func TestParseWeightReels2_Success(t *testing.T) {
	// build JSON that matches expected structure for getConfigInCell
	jsonStr := `{
        "componentValues": {
            "label": "wlabel",
            "configuration": {
                "reelSetWeight": "rw",
                "isExpandReel": true
            },
            "controller": [
                {"type": "addRespinTimes", "target": "tg", "times": 1}
            ]
        }
    }`

	var node ast.Node
	err := sonic.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		t.Fatalf("sonic.Unmarshal err=%v", err)
	}

	bc := &BetConfig{mapConfig: make(map[string]IComponentConfig), mapBasicConfig: make(map[string]*BasicComponentConfig)}
	label, err := parseWeightReels2(bc, &node)
	assert.NoError(t, err)
	assert.Equal(t, "wlabel", label)
	// ensure config added
	_, ok := bc.mapConfig["wlabel"]
	assert.True(t, ok)
}
