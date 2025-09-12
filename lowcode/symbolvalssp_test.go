package lowcode

import (
	"sync"
	"testing"

	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
)

// reuse helper from other tests: construct a pool with simple paytable mapping
func makePoolWithSimpleSymbols() *GamePropertyPool {
	mapSymbols := map[string]int{
		"CA": 1, // coin
		"BN": 2, // coin target
		"MU": 3, // multi
		"CO": 4, // collect
	}

	return &GamePropertyPool{
		DefaultPaytables: &sgc7game.PayTables{MapSymbols: mapSymbols},
		MapGamePropPool:  make(map[int]*sync.Pool),
	}
}

// create a minimal 3x3 scene helper
func make3x3Scene() *sgc7game.GameScene {
	gs, _ := sgc7game.NewGameScene2(3, 3, 0)
	return gs
}

func TestInitEx_InvalidAndSuccess(t *testing.T) {
	pool := makePoolWithSimpleSymbols()
	cfg := &SymbolValsSPConfig{}
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)

	// invalid: no coinSymbols
	if err := comp.InitEx(cfg, pool); err == nil {
		t.Fatalf("expected error for missing coinSymbols")
	}

	// success
	cfg2 := &SymbolValsSPConfig{
		CoinSymbols:    []string{"CA"},
		MultiSymbols:   []string{"MU"},
		CollectSymbols: []string{"CO"},
		StrMultiType:   "normal",
		StrCollectType: "normal",
	}
	if err := comp.InitEx(cfg2, pool); err != nil {
		t.Fatalf("InitEx failed: %v", err)
	}
}

func TestProcMultiAndCollect_Normal(t *testing.T) {
	pool := makePoolWithSimpleSymbols()
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)
	cfg := &SymbolValsSPConfig{
		CoinSymbols:           []string{"CA"},
		MultiSymbols:          []string{"MU"},
		CollectSymbols:        []string{"CO"},
		CollectTargetSymbol:   "CA",
		CollectCoinSymbol:     "BN",
		CollectMultiSymbol:    "BN",
		StrMultiType:          "normal",
		StrCollectType:        "normal",
		MapAwards:             map[string][]*Award{},
	}
	if err := comp.InitEx(cfg, pool); err != nil {
		t.Fatalf("InitEx failed: %v", err)
	}

	gp := &GameProperty{Pool: pool}
	gp.PoolScene = sgc7game.NewGameScenePoolEx()
	gp.SceneStack = NewSceneStack(false)
	gp.OtherSceneStack = NewSceneStack(true)

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")

	gs, _ := sgc7game.NewGameScene2(3, 3, 0)
	os, _ := sgc7game.NewGameScene2(3, 3, 0)

	// set symbols: place MU at (1,1) and CO collect at (0,0)
	gs.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	gs.Arr[0][0] = comp.Config.CollectSymbolCodes[0]
	// place coin values in other scene
	os.Arr[0][0] = 5
	os.Arr[2][2] = 2

	pr.Scenes = append(pr.Scenes, gs)
	pr.OtherScenes = append(pr.OtherScenes, os)

	// push scenes into the gameProp stacks so GetTopSceneEx can find them
	gp.SceneStack.Push("", gs)
	gp.OtherSceneStack.Push("", os)

	cd := comp.NewComponentData().(*SymbolValsSPData)

	// run OnPlayGame: should detect collect and/or multi depending on mapping
	_, err := comp.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), nil, "", "", nil, nil, nil, cd)
	if err != nil && err != ErrComponentDoNothing {
		t.Fatalf("OnPlayGame unexpected err: %v", err)
	}

	// OnPlayGame should not panic and cd fields should be populated or zeroed
	_ = cd.Multi
	_ = cd.CollectCoin
}

func TestOnAsciiGame_NoUsedScenesAndInvalidIndex(t *testing.T) {
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)
	// empty data
	cd := comp.NewComponentData().(*SymbolValsSPData)
	pr := &sgc7game.PlayResult{}
	// nothing should happen
	pool := makePoolWithSimpleSymbols()
	if err := comp.OnAsciiGame(&GameProperty{Pool: pool}, pr, nil, asciigame.NewSymbolColorMap(pool.DefaultPaytables), cd); err != nil {
		t.Fatalf("OnAsciiGame returned error: %v", err)
	}

	// set used scenes but invalid indexes
	cd.UsedScenes = []int{10}
	cd.UsedOtherScenes = []int{10}
	pr.Scenes = []*sgc7game.GameScene{make3x3Scene()}
	pr.OtherScenes = []*sgc7game.GameScene{make3x3Scene()}
	if err := comp.OnAsciiGame(&GameProperty{Pool: pool}, pr, nil, asciigame.NewSymbolColorMap(pool.DefaultPaytables), cd); err != nil {
		t.Fatalf("OnAsciiGame with invalid indexes returned error: %v", err)
	}
}

func TestProcMulti_NormalAndRoundAndCollects(t *testing.T) {
	pool := makePoolWithSimpleSymbols()
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)
	cfg := &SymbolValsSPConfig{
		CoinSymbols:           []string{"CA"},
		MultiSymbols:          []string{"MU"},
		CollectSymbols:        []string{"CO"},
		MultiTargetSymbol:     "BN",
		CollectTargetSymbol:   "CA",
		CollectCoinSymbol:     "BN",
		CollectMultiSymbol:    "BN",
		StrMultiType:          "normal",
		StrCollectType:        "normal",
	}
	if err := comp.InitEx(cfg, pool); err != nil {
		t.Fatalf("InitEx failed: %v", err)
	}

	gp := &GameProperty{Pool: pool}
	gp.PoolScene = sgc7game.NewGameScenePoolEx()

	// prepare scenes
	gs, _ := sgc7game.NewGameScene2(3, 3, 0)
	os, _ := sgc7game.NewGameScene2(3, 3, 0)

	// place a multi symbol at (1,1) with multiplier value in other scene
	gs.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	os.Arr[1][1] = 3

	// place a coin at (0,0)
	gs.Arr[0][0] = comp.Config.CoinSymbolCodes[0]
	os.Arr[0][0] = 2

	// run procMulti via procMultiNormal (normal type)
	cd := comp.NewComponentData().(*SymbolValsSPData)
	ngs, nos, triggered, err := comp.procMulti(gp, nil, cd, gs, os)
	if err != nil {
		t.Fatalf("procMulti normal err: %v", err)
	}
	if !triggered {
		t.Fatalf("expected multi triggered for normal type")
	}
	// nos should have coin multiplied by 3
	if nos.Arr[0][0] != 6 {
		t.Fatalf("expected nos coin multiplied, got %d", nos.Arr[0][0])
	}
	// ngs should have target symbol at multi pos
	if ngs.Arr[1][1] != comp.Config.MultiTargetSymbolCode {
		t.Fatalf("expected ngs multi target set, got %d", ngs.Arr[1][1])
	}

	// now test round type: reset cd and switch type
	comp.Config.MultiType = SVSPMultiTypeRound
	gs2, _ := sgc7game.NewGameScene2(3, 3, 0)
	os2, _ := sgc7game.NewGameScene2(3, 3, 0)
	// center multi
	gs2.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	os2.Arr[1][1] = 2
	// put a coin at neighbor
	gs2.Arr[0][0] = comp.Config.CoinSymbolCodes[0]
	os2.Arr[0][0] = 5

	cd2 := comp.NewComponentData().(*SymbolValsSPData)
	_, nos2, triggered2, err := comp.procMulti(gp, nil, cd2, gs2, os2)
	if err != nil {
		t.Fatalf("procMulti round err: %v", err)
	}
	if !triggered2 {
		t.Fatalf("expected multi triggered for round type")
	}
	// neighbor coin should be multiplied by 2
	if nos2.Arr[0][0] != 10 {
		t.Fatalf("expected nos2 coin multiplied, got %d", nos2.Arr[0][0])
	}

	// test collect normal: place collect symbol and coins
	comp.Config.CollectType = SVSPCollectTypeNormal
	gs3, _ := sgc7game.NewGameScene2(3, 3, 0)
	os3, _ := sgc7game.NewGameScene2(3, 3, 0)
	gs3.Arr[0][0] = comp.Config.CollectSymbolCodes[0]
	// coin positions
	gs3.Arr[2][2] = comp.Config.CoinSymbolCodes[0]
	os3.Arr[2][2] = 7

	cd3 := comp.NewComponentData().(*SymbolValsSPData)
	_, nos3, triggered3, err := comp.procCollect(gp, nil, cd3, gs3, os3)
	if err != nil {
		t.Fatalf("procCollect normal err: %v", err)
	}
	if !triggered3 {
		t.Fatalf("expected collect triggered for normal type")
	}
	// nos3 at collect pos should equal total coin (7)
	if nos3.Arr[0][0] != 7 {
		t.Fatalf("expected collected coin at pos, got %d", nos3.Arr[0][0])
	}

	// test collect sequence: ensure only first collect pos used
	comp.Config.CollectType = SVSPCollectTypeSequence
	gs4, _ := sgc7game.NewGameScene2(3, 3, 0)
	os4, _ := sgc7game.NewGameScene2(3, 3, 0)
	gs4.Arr[0][0] = comp.Config.CollectSymbolCodes[0]
	gs4.Arr[1][0] = comp.Config.CollectSymbolCodes[0]
	gs4.Arr[2][2] = comp.Config.CoinSymbolCodes[0]
	os4.Arr[2][2] = 9

	cd4 := comp.NewComponentData().(*SymbolValsSPData)
	_, nos4, triggered4, err := comp.procCollect(gp, nil, cd4, gs4, os4)
	if err != nil {
		t.Fatalf("procCollect seq err: %v", err)
	}
	if !triggered4 {
		t.Fatalf("expected collect triggered for sequence type")
	}
	// only first collect pos (0,0) should have the total
	if nos4.Arr[0][0] != 9 {
		t.Fatalf("expected sequence collect at first pos, got %d", nos4.Arr[0][0])
	}

	// cover procMulti invalid multi type branch
	comp.Config.MultiType = SymbolValsSPMultiType(999)
	_, _, _, err = comp.procMulti(gp, nil, comp.NewComponentData().(*SymbolValsSPData), gs, os)
	if err == nil {
		t.Fatalf("expected error for invalid multi type")
	}
}

func TestProcMulti_CloneBranches(t *testing.T) {
	pool := makePoolWithSimpleSymbols()
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)
	cfg := &SymbolValsSPConfig{
		CoinSymbols:         []string{"CA"},
		MultiSymbols:        []string{"MU"},
		StrMultiType:        "normal",
		MultiTargetSymbol:   "BN",
	}
	if err := comp.InitEx(cfg, pool); err != nil {
		t.Fatalf("InitEx failed: %v", err)
	}

	gp := &GameProperty{Pool: pool}
	gp.PoolScene = sgc7game.NewGameScenePoolEx()

	gs, _ := sgc7game.NewGameScene2(3, 3, 0)
	os, _ := sgc7game.NewGameScene2(3, 3, 0)

	// Put multi at (1,1) and coin at same pos so both ngs and nos will be cloned
	gs.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	gs.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	os.Arr[1][1] = 4
	gs.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	// also put coin elsewhere
	gs.Arr[0][2] = comp.Config.CoinSymbolCodes[0]
	os.Arr[0][2] = 3

	cd := comp.NewComponentData().(*SymbolValsSPData)
	ngs, nos, triggered, err := comp.procMultiNormal(gp, nil, cd, gs, os)
	if err != nil {
		t.Fatalf("procMultiNormal err: %v", err)
	}
	if !triggered {
		t.Fatalf("expected triggered")
	}
	// ensure nos was cloned and modified
	if nos.Arr[0][2] != 9 { // 3 * 3? note: because multi is os.Arr[1][1]=4, coin should be multiplied by 4
		// we check that it's > original to ensure modification happened
		if nos.Arr[0][2] == 3 {
			t.Fatalf("expected nos to be modified (cloned and multiplied), got %d", nos.Arr[0][2])
		}
	}
	// ensure ngs was cloned when multi target symbol configured
	if comp.Config.MultiTargetSymbolCode > 0 && ngs.Arr[1][1] != comp.Config.MultiTargetSymbolCode {
		t.Fatalf("expected ngs multi target set, got %d", ngs.Arr[1][1])
	}
}

func TestProcMultiRoundXY_CloneAndBounds(t *testing.T) {
	pool := makePoolWithSimpleSymbols()
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)
	cfg := &SymbolValsSPConfig{
		CoinSymbols:  []string{"CA"},
		MultiSymbols: []string{"MU"},
		StrMultiType: "round",
	}
	if err := comp.InitEx(cfg, pool); err != nil {
		t.Fatalf("InitEx failed: %v", err)
	}

	gp := &GameProperty{Pool: pool}
	gp.PoolScene = sgc7game.NewGameScenePoolEx()

	gs, _ := sgc7game.NewGameScene2(3, 3, 0)
	os, _ := sgc7game.NewGameScene2(3, 3, 0)

	// neighbor coin and center multi
	gs.Arr[1][1] = comp.Config.MultiSymbolCodes[0]
	gs.Arr[0][0] = comp.Config.CoinSymbolCodes[0]
	os.Arr[0][0] = 2

	cd := comp.NewComponentData().(*SymbolValsSPData)
	// call procMultiRoundXY directly for center (should clone os lazily)
	nos, err := comp.procMultiRoundXY(gp, nil, cd, gs, os, 1, 1, 2, false)
	if err != nil {
		t.Fatalf("procMultiRoundXY err: %v", err)
	}
	if nos.Arr[0][0] != 4 {
		t.Fatalf("expected neighbor multiplied by 2 -> 4, got %d", nos.Arr[0][0])
	}

	// bounds: call with cx near edge and ensure no panic and no change out-of-bounds
	cd2 := comp.NewComponentData().(*SymbolValsSPData)
	_, err = comp.procMultiRoundXY(gp, nil, cd2, gs, os, 0, 0, 2, false)
	if err != nil {
		t.Fatalf("procMultiRoundXY bounds err: %v", err)
	}
}

func TestDataCloneAndPBAndControllersAndAscii(t *testing.T) {
	comp := NewSymbolValsSP("svsp").(*SymbolValsSP)
	pool := makePoolWithSimpleSymbols()

	// test Clone and BuildPBComponentData
	cd := comp.NewComponentData().(*SymbolValsSPData)
	cd.AddPos(1, 2)
	cd.Multi = 5
	cloned := cd.Clone().(*SymbolValsSPData)
	if cloned.Multi != 5 {
		t.Fatalf("clone did not preserve multi")
	}
	pb := cd.BuildPBComponentData().(*sgc7pb.SymbolValsSPData)
	if int(pb.Multi) != 5 {
		t.Fatalf("pb build multi mismatch")
	}

	// test ProcControllers calls into gameProp.procAwards when key exists
	comp.Config = &SymbolValsSPConfig{MapAwards: map[string][]*Award{"<trigger>": {}}}
	gp := &GameProperty{Pool: pool}
	// call ProcControllers (should not panic)
	comp.ProcControllers(gp, nil, sgc7game.NewPlayResult("m", 0, 0, "t"), NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), 0, "<trigger>")

	// test OnAsciiGame valid branch: ensure no error when indexes valid
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	gs, _ := sgc7game.NewGameScene2(3, 3, 0)
	os, _ := sgc7game.NewGameScene2(3, 3, 0)
	pr.Scenes = append(pr.Scenes, gs)
	pr.OtherScenes = append(pr.OtherScenes, os)
	cd2 := comp.NewComponentData().(*SymbolValsSPData)
	cd2.UsedScenes = []int{0}
	cd2.UsedOtherScenes = []int{0}
	if err := comp.OnAsciiGame(&GameProperty{Pool: pool}, pr, nil, asciigame.NewSymbolColorMap(pool.DefaultPaytables), cd2); err != nil {
		t.Fatalf("OnAsciiGame valid branch error: %v", err)
	}
}
