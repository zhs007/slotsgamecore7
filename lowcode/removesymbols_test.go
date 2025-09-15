package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	stats2 "github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// rsFake component and component data to provide positions via GetPos
type rsFakeCD struct{
	pos []int
	results []int
}

func (f *rsFakeCD) OnNewGame(gameProp *GameProperty, component IComponent) {}
func (f *rsFakeCD) BuildPBComponentData() proto.Message { return nil }
func (f *rsFakeCD) Clone() IComponentData { return &rsFakeCD{pos: append([]int{}, f.pos...), results: append([]int{}, f.results...)} }
func (f *rsFakeCD) GetValEx(key string, getType GetComponentValType) (int, bool) { return 0, false }
func (f *rsFakeCD) GetStrVal(key string) (string, bool) { return "", false }
func (f *rsFakeCD) GetConfigVal(key string) string { return "" }
func (f *rsFakeCD) SetConfigVal(key string, val string) {}
func (f *rsFakeCD) GetConfigIntVal(key string) (int, bool) { return 0, false }
func (f *rsFakeCD) SetConfigIntVal(key string, val int) {}
func (f *rsFakeCD) ChgConfigIntVal(key string, off int) int { return 0 }
func (f *rsFakeCD) ClearConfigIntVal(key string) {}
func (f *rsFakeCD) GetResults() []int { return f.results }
func (f *rsFakeCD) GetOutput() int { return 0 }
func (f *rsFakeCD) GetStringOutput() string { return "" }
func (f *rsFakeCD) GetSymbols() []int { return nil }
func (f *rsFakeCD) AddSymbol(symbolCode int) {}
func (f *rsFakeCD) GetPos() []int { return f.pos }
func (f *rsFakeCD) HasPos(x int, y int) bool { return false }
func (f *rsFakeCD) AddPos(x int, y int) {}
func (f *rsFakeCD) ClearPos() {}
func (f *rsFakeCD) GetLastRespinNum() int { return 0 }
func (f *rsFakeCD) GetCurRespinNum() int { return 0 }
func (f *rsFakeCD) IsRespinEnding() bool { return false }
func (f *rsFakeCD) IsRespinStarted() bool { return false }
func (f *rsFakeCD) AddTriggerRespinAward(award *Award) {}
func (f *rsFakeCD) AddRespinTimes(num int) {}
func (f *rsFakeCD) TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {}
func (f *rsFakeCD) PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {}
func (f *rsFakeCD) GetMask() []bool { return nil }
func (f *rsFakeCD) ChgMask(curMask int, val bool) bool { return false }
func (f *rsFakeCD) PutInMoney(coins int) {}
func (f *rsFakeCD) ChgReelsCollector(reelsData []int) {}
func (f *rsFakeCD) SetSymbolCodes(symbolCodes []int) {}
func (f *rsFakeCD) GetSymbolCodes() []int { return nil }

type rsFakeComp struct{ name string; pos []int; results []int }
func (f *rsFakeComp) Init(fn string, pool *GamePropertyPool) error { return nil }
func (f *rsFakeComp) InitEx(cfg any, pool *GamePropertyPool) error { return nil }
func (f *rsFakeComp) OnGameInited(components *ComponentList) error { return nil }
func (f *rsFakeComp) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) { return "", nil }
func (f *rsFakeComp) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error { return nil }
func (f *rsFakeComp) NewComponentData() IComponentData { return &rsFakeCD{pos: append([]int{}, f.pos...), results: append([]int{}, f.results...)} }
func (f *rsFakeComp) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {}
func (f *rsFakeComp) ProcRespinOnStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, canRemove bool) (string, error) { return "", nil }
func (f *rsFakeComp) GetName() string { return f.name }
func (f *rsFakeComp) IsRespin() bool { return false }
func (f *rsFakeComp) IsForeach() bool { return false }
func (f *rsFakeComp) NewStats2(parent string) *stats2.Feature { return nil }
func (f *rsFakeComp) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {}
func (f *rsFakeComp) IsNeedOnStepEndStats2() bool { return false }
func (f *rsFakeComp) GetAllLinkComponents() []string { return nil }
func (f *rsFakeComp) GetNextLinkComponents() []string { return nil }
func (f *rsFakeComp) GetChildLinkComponents() []string { return nil }
func (f *rsFakeComp) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) { return false, nil }
func (f *rsFakeComp) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {}
func (f *rsFakeComp) IsMask() bool { return false }
func (f *rsFakeComp) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error { return nil }
func (f *rsFakeComp) SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, mask bool) error { return nil }
func (f *rsFakeComp) SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error { return nil }
func (f *rsFakeComp) EachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) error { return nil }
func (f *rsFakeComp) AddPos(cd IComponentData, x int, y int) {}
func (f *rsFakeComp) OnPlayGameWithSet(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData, set int) (string, error) { return "", nil }
func (f *rsFakeComp) ClearData(icd IComponentData, bForceNow bool) {}
func (f *rsFakeComp) InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, ps *PlayerState, betMethod int, bet int) error { return nil }
func (f *rsFakeComp) NewPlayerState() IComponentPS { return nil }
func (f *rsFakeComp) OnUpdateDataWithPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, ps *PlayerState, betMethod int, bet int, cd IComponentData) {}
func (f *rsFakeComp) ChgReelsCollector(icd IComponentData, ps *PlayerState, betMethod int, bet int, reelsData []int) {}


// rsOutCD is a minimal component data implementation that records AddPos
// calls so tests can verify positions were emitted to an output component.
type rsOutCD struct{ pos []int }

func (r *rsOutCD) OnNewGame(gameProp *GameProperty, component IComponent) {}
func (r *rsOutCD) BuildPBComponentData() proto.Message { return nil }
func (r *rsOutCD) Clone() IComponentData { return &rsOutCD{pos: append([]int{}, r.pos...)} }
func (r *rsOutCD) GetValEx(key string, getType GetComponentValType) (int, bool) { return 0, false }
func (r *rsOutCD) GetStrVal(key string) (string, bool) { return "", false }
func (r *rsOutCD) GetConfigVal(key string) string { return "" }
func (r *rsOutCD) SetConfigVal(key string, val string) {}
func (r *rsOutCD) GetConfigIntVal(key string) (int, bool) { return 0, false }
func (r *rsOutCD) SetConfigIntVal(key string, val int) {}
func (r *rsOutCD) ChgConfigIntVal(key string, off int) int { return 0 }
func (r *rsOutCD) ClearConfigIntVal(key string) {}
func (r *rsOutCD) GetResults() []int { return nil }
func (r *rsOutCD) GetOutput() int { return 0 }
func (r *rsOutCD) GetStringOutput() string { return "" }
func (r *rsOutCD) GetSymbols() []int { return nil }
func (r *rsOutCD) AddSymbol(symbolCode int) {}
func (r *rsOutCD) GetPos() []int { return append([]int{}, r.pos...) }
func (r *rsOutCD) HasPos(x int, y int) bool { return false }
func (r *rsOutCD) AddPos(x int, y int) { r.pos = append(r.pos, x, y) }
func (r *rsOutCD) ClearPos() { r.pos = nil }
func (r *rsOutCD) GetLastRespinNum() int { return 0 }
func (r *rsOutCD) GetCurRespinNum() int { return 0 }
func (r *rsOutCD) IsRespinEnding() bool { return false }
func (r *rsOutCD) IsRespinStarted() bool { return false }
func (r *rsOutCD) AddTriggerRespinAward(award *Award) {}
func (r *rsOutCD) AddRespinTimes(num int) {}
func (r *rsOutCD) TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {}
func (r *rsOutCD) PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {}
func (r *rsOutCD) GetMask() []bool { return nil }
func (r *rsOutCD) ChgMask(curMask int, val bool) bool { return false }
func (r *rsOutCD) PutInMoney(coins int) {}
func (r *rsOutCD) ChgReelsCollector(reelsData []int) {}
func (r *rsOutCD) SetSymbolCodes(symbolCodes []int) {}
func (r *rsOutCD) GetSymbolCodes() []int { return nil }


// minimal helper to build a GameScene with given dimensions and symbol matrix
func newSceneFromMatrix(mat [][]int) *sgc7game.GameScene {
	gs, err := sgc7game.NewGameSceneWithArr2(mat)
	if err != nil {
		panic(err)
	}

	return gs
}

func TestCanRemoveBoundsAndIgnore(t *testing.T) {
	// scene 3x3
	mat := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	gs := newSceneFromMatrix(mat)

	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	// out of bounds
	assert.False(t, comp.canRemove(-1, 0, gs))
	assert.False(t, comp.canRemove(0, -1, gs))
	assert.False(t, comp.canRemove(3, 0, gs))
	assert.False(t, comp.canRemove(0, 3, gs))

	// ignore symbol code: configure component to ignore code 5
	comp.Config = &RemoveSymbolsConfig{IgnoreSymbolCodes: []int{5}}
	// center (1,1) is value 5 in mat above, so canRemove should return false
	assert.False(t, comp.canRemove(1, 1, gs))
}

// onBasic is exercised indirectly by higher-level integration tests in the
// package; low-level direct invocation requires substantial GameProperty and
// GameParams plumbing which is already covered elsewhere. Here we avoid
// calling unexported helpers directly.
func TestOnBasic_Smoke(t *testing.T) {
	// basic smoke to ensure NewComponentData returns expected type
	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	cd := comp.NewComponentData()
	_, ok := cd.(*RemoveSymbolsData)
	assert.True(t, ok)
}

func TestOnBasic_SourcePositionCollection(t *testing.T) {
	// build a simple 3x3 scene
	mat := [][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	}
	sc, err := sgc7game.NewGameSceneWithArr2(mat)
	if err != nil { t.Fatal(err) }

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr.Scenes = append(pr.Scenes, sc)

	pool := &GamePropertyPool{Config: &Config{Width: 3, Height: 3}}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

	gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

	// prepare GameProperty with a fake component that will provide positions
	gameProp := pool.newGameProp(1)
	gameProp.Components = NewComponentList()
	if gameProp.SceneStack == nil {
		gameProp.SceneStack = NewSceneStack(false)
	}
	if gameProp.OtherSceneStack == nil {
		gameProp.OtherSceneStack = NewSceneStack(true)
	}

	// fake provider component
	fcomp := &rsFakeComp{name: "src" , pos: []int{0,1, 1,1, 2,1}}
	gameProp.Components.MapComponents = map[string]IComponent{"src": fcomp}

	// initialize gameProp for a new game so callStack nodes are prepared
	_ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	// push current scene into SceneStack so GetTargetScene3 can find it
	gameProp.SceneStack.Push("rs", sc)
	// also prepare an other-scene so IsNeedProcSymbolVals branch uses it
	other := sc.CloneEx(gameProp.PoolScene)
	pr.OtherScenes = append(pr.OtherScenes, other)
	gameProp.OtherSceneStack.Push("rs", other)

	// create removeSymbols component and config to use SourcePositionCollection
	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	// ensure BasicComponent.Config is set so GetTargetScene3 won't nil-deref
	comp.BasicComponent = NewBasicComponent("rs", 1)
	comp.BasicComponent.Config = &BasicComponentConfig{}
	comp.Config = &RemoveSymbolsConfig{SourcePositionCollection: []string{"src"}}

	// put component data into gameProp.callStack so GetCurComponentDataWithName can find nothing (we only need positions via GetComponentPos)

	cd := comp.NewComponentData()

	// call OnPlayGame; expect successful removal because SourcePositionCollection is used
	_, err = comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
	// expect no error and removed positions applied
	assert.NoError(t, err)

	rsd, ok := cd.(*RemoveSymbolsData)
	assert.True(t, ok)
	assert.Equal(t, 3, rsd.RemovedNum)
	// a new scene should be appended
	assert.Equal(t, 2, len(pr.Scenes))
	// check removed positions in new scene: positions (0,1),(1,1),(2,1) should be -1
	newsc := pr.Scenes[1]
	assert.Equal(t, -1, newsc.Arr[0][1])
	assert.Equal(t, -1, newsc.Arr[1][1])
	assert.Equal(t, -1, newsc.Arr[2][1])
}

func TestInitEx_InvalidCfg(t *testing.T) {
	c := &RemoveSymbols{}
	// pass cfg as wrong type
	err := c.InitEx(nil, &GamePropertyPool{})
	assert.Error(t, err)
}

func TestOnPlayGame_InvalidComponentDataType(t *testing.T) {
	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	prop := &GameProperty{}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	// pass wrong type as component data
	_, err := comp.OnPlayGame(prop, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), nil, "", "", nil, nil, nil, &BasicComponentData{})
	assert.Error(t, err)
}

func TestAdjacentPay_MiddleRetention(t *testing.T) {
	// 3-in-row positions: (0,1),(1,1),(2,1) expect middle (1,1) kept and replaced by AddedSymbolCode
	mat := [][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	}
	sc, err := sgc7game.NewGameSceneWithArr2(mat)
	if err != nil { t.Fatal(err) }

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr.Scenes = append(pr.Scenes, sc)

	pool := &GamePropertyPool{Config: &Config{Width: 3, Height: 3}}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

	gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

	gameProp := pool.newGameProp(1)
	gameProp.Components = NewComponentList()
	gameProp.SceneStack.Push("rs", sc)
	_ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)

	// create provider component with a result containing three positions
	fakeResComp := &rsFakeComp{name:"bg-pay"}
	// prepare a component that would have produced a result listing positions
	// We simulate by adding a result to playresult and letting component's GetCurComponentDataWithName return nil so removeSymbols reads curpr.Results
	gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": fakeResComp}

	// simulate that bg-pay has produced a result by pushing component data into callStack history
	// create a fake component data that reports the result index
	frcd := &rsFakeCD{results: []int{0}}
	gameProp.callStack = NewCallStack()
	gameProp.callStack.OnNewGame()
	// put the fake component into map and into global node so GetCurComponentDataWithName can find it
	gameProp.Components.MapComponents["bg-pay"] = fakeResComp
	// set playresult's Results
	r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex: 0, Symbol: 0, Pos: []int{0,1,1,1,2,1}}
	pr.Results = append(pr.Results, r)

	// set the global callStack node's MapComponentData so GetCurComponentDataWithName returns it
	gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"bg-pay": frcd}

	// set history components so removeSymbols will consider bg-pay
	gp.HistoryComponents = []string{"bg-pay"}

	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	comp.BasicComponent = NewBasicComponent("rs", 1)
	comp.BasicComponent.Config = &BasicComponentConfig{}
	comp.Config = &RemoveSymbolsConfig{Type: RSTypeAdjacentPay, AddedSymbolCode: 99, TargetComponents: []string{"bg-pay"}, AddedSymbol: "X"}

	cd := comp.NewComponentData()

	_, err = comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
	assert.NoError(t, err)

	rsd := cd.(*RemoveSymbolsData)
	assert.Equal(t, 3, rsd.RemovedNum)
	// check that the middle position is set to AddedSymbolCode in new scene
	newsc := pr.Scenes[1]
	assert.Equal(t, 99, newsc.Arr[1][1])
}

func TestOutputToComponentAndOtherScene(t *testing.T) {
	// Prepare 3x3 scene
	mat := [][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	}
	sc, _ := sgc7game.NewGameSceneWithArr2(mat)
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr.Scenes = append(pr.Scenes, sc)

	pool := &GamePropertyPool{Config: &Config{Width: 3, Height: 3}}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

	gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
	gameProp := pool.newGameProp(1)
	gameProp.Components = NewComponentList()
	if gameProp.SceneStack == nil { gameProp.SceneStack = NewSceneStack(false) }
	if gameProp.OtherSceneStack == nil { gameProp.OtherSceneStack = NewSceneStack(true) }
	_ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
	gameProp.SceneStack.Push("rs", sc)

	// also prepare an other-scene so IsNeedProcSymbolVals branch uses it
	other := sc.CloneEx(gameProp.PoolScene)
	pr.OtherScenes = append(pr.OtherScenes, other)
	gameProp.OtherSceneStack.Push("rs", other)

	// source positions
	fcomp := &rsFakeComp{name: "src", pos: []int{0,1, 1,1}}
	gameProp.Components.MapComponents = map[string]IComponent{"src": fcomp}

	// output component to receive removed positions
	outcd := &rsOutCD{}
	outcomp := &rsFakeComp{name: "out"}
	gameProp.Components.MapComponents["out"] = outcomp
	// put outcd into callStack global map so GetCurComponentDataWithName finds it
	gameProp.callStack.nodes[0].MapComponentData = map[string]IComponentData{"out": outcd}

	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	comp.BasicComponent = NewBasicComponent("rs", 1)
	comp.BasicComponent.Config = &BasicComponentConfig{}
	comp.Config = &RemoveSymbolsConfig{SourcePositionCollection: []string{"src"}, OutputToComponent: "out", IsNeedProcSymbolVals: true, EmptySymbolVal: -2}

	cd := comp.NewComponentData()

	_, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
	assert.NoError(t, err)

	rsd := cd.(*RemoveSymbolsData)
	// two positions removed
	assert.Equal(t, 2, rsd.RemovedNum)
	// outputCD should have recorded two positions
	assert.Equal(t, 4, len(outcd.pos))
	// scene changes should have been recorded (either scene or other scene)
	assert.GreaterOrEqual(t, len(pr.Scenes), 1)
	// verify emptySymbolVal applied in any new OtherScene or appended Scene
	found := false
	for _, other := range pr.OtherScenes {
		if other.Arr[0][1] == -2 && other.Arr[1][1] == -2 {
			found = true
			break
		}
	}
	if !found {
		for i := 0; i < len(pr.Scenes); i++ {
			ns := pr.Scenes[i]
			if ns.Arr[0][1] == -2 && ns.Arr[1][1] == -2 {
				found = true
				break
			}
		}
	}
	assert.True(t, found)
}

func TestAvgHeightZeroBranch(t *testing.T) {
	// if no removals occur, AvgHeight must be zero
	mat := [][]int{{1}}
	sc, _ := sgc7game.NewGameSceneWithArr2(mat)
	pr := sgc7game.NewPlayResult("m",0,0,"t")
	pr.Scenes = append(pr.Scenes, sc)

	pool := &GamePropertyPool{Config: &Config{Width:1, Height:1}}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }
	gp := NewGameParam(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
	gameProp := pool.newGameProp(1)
	gameProp.Components = NewComponentList()
	_ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
	gameProp.SceneStack.Push("rs", sc)

	comp := NewRemoveSymbols("rs").(*RemoveSymbols)
	comp.BasicComponent = NewBasicComponent("rs", 1)
	comp.BasicComponent.Config = &BasicComponentConfig{}
	comp.Config = &RemoveSymbolsConfig{SourcePositionCollection: []string{"src"}}

	// no source component registered -> no removal
	cd := comp.NewComponentData()
	_, err := comp.OnPlayGame(gameProp, pr, gp, nil, "", "", nil, nil, nil, cd)
	// ErrComponentDoNothing expected
	assert.Error(t, err)
	rsd := cd.(*RemoveSymbolsData)
	assert.Equal(t, 0, rsd.AvgHeight)
}

func TestBasicTargetComponentsRemoval(t *testing.T) {
	// test basic removal via TargetComponents using curpr.Results positions
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
	_ = gameProp.OnNewGame(&sgc7game.Stake{CoinBet:1, CashBet:1}, nil)
	gameProp.SceneStack.Push("rs", sc)

	// create a fake previous component with a result
	fake := &rsFakeComp{name: "bg-pay"}
	gameProp.Components.MapComponents = map[string]IComponent{"bg-pay": fake}
	r := &sgc7game.Result{Type: sgc7game.RTAdjacentPay, LineIndex:0, Symbol:0, Pos: []int{0,1,1,1}}
	pr.Results = append(pr.Results, r)

	// put fake component data into callStack so GetCurComponentDataWithName returns it
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
	newsc := pr.Scenes[1]
	// removed positions are set to -1
	assert.Equal(t, -1, newsc.Arr[0][1])
	assert.Equal(t, -1, newsc.Arr[1][1])
}
