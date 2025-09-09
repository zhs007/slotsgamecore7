package lowcode

import (
	"context"
	"errors"
	"sync"
	"testing"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// errPlugin returns error on Random to exercise error paths
type errPlugin struct{}

func (e *errPlugin) Random(ctx context.Context, r int) (int, error) {
	return 0, errors.New("random fail")
}
func (e *errPlugin) GetUsedRngs() []*sgc7utils.RngInfo { return nil }
func (e *errPlugin) ClearUsedRngs()                    {}
func (e *errPlugin) TagUsedRngs()                      {}
func (e *errPlugin) RollbackUsedRngs() error           { return nil }
func (e *errPlugin) SetCache(arr []int)                {}
func (e *errPlugin) ClearCache()                       {}
func (e *errPlugin) Init()                             {}
func (e *errPlugin) SetScenePool(any)                  {}
func (e *errPlugin) GetScenePool() any                 { return nil }
func (e *errPlugin) SetSeed(seed int)                  {}

func TestUtils_Basics(t *testing.T) {
	if intPow(2, 10) != 1024 {
		t.Fatalf("intPow wrong")
	}
	if intPow(5, 0) != 1 {
		t.Fatalf("intPow zero")
	}
	if intPow(2, -1) != 1 {
		t.Fatalf("intPow negative should return 1 per current impl")
	}

	if !isSameBoolSlice([]bool{true, false}, []bool{true, false}) {
		t.Fatalf("isSameBoolSlice equals")
	}
	if isSameBoolSlice([]bool{true}, []bool{}) {
		t.Fatalf("isSameBoolSlice different")
	}

	if len(InsStringSliceNonRep([]string{"a"}, "a")) != 1 {
		t.Fatalf("InsStringSliceNonRep duplicate")
	}
	if len(InsStringSliceNonRep([]string{"a"}, "b")) != 2 {
		t.Fatalf("InsStringSliceNonRep add")
	}
	if len(InsSliceNonRep([]string{"a"}, []string{"b", "a", "c"})) != 3 {
		t.Fatalf("InsSliceNonRep result")
	}

	if !CmpVal(5, ">", 3) {
		t.Fatalf("CmpVal > failed")
	}

	if HasSamePos([]int{}, []int{1, 2}) {
		t.Fatalf("HasSamePos empty src")
	}
	if !HasSamePos([]int{1, 2, 3, 4}, []int{7, 8, 1, 2}) {
		t.Fatalf("HasSamePos should find")
	}

	if !IsInitialArr([]int{0, 1, 2}) {
		t.Fatalf("IsInitialArr")
	}
	if !IsSameIntArr(GenInitialArr(3), []int{0, 1, 2}) {
		t.Fatalf("GenInitialArr/IsSameIntArr")
	}

	if AbsInt(-5) != 5 || AbsInt(5) != 5 {
		t.Fatalf("AbsInt")
	}
}

func TestUtils_ProcCheat_RNG_and_ForceOutcome(t *testing.T) {
	p := sgc7plugin.NewMockPlugin()
	fo, err := ProcCheat(p, "1,2,3")
	if err != nil {
		t.Fatalf("ProcCheat rng err: %v", err)
	}
	if fo != nil {
		t.Fatalf("ProcCheat rng should not produce FO")
	}
	if len(p.Cache) == 0 {
		t.Fatalf("ProcCheat did not call SetCache")
	}

	SetAllowForceOutcome(3)
	fo2, err := ProcCheat(nil, "totalWins == 0")
	if err != nil {
		t.Fatalf("ProcCheat force outcome err: %v", err)
	}
	if fo2 == nil {
		t.Fatalf("ProcCheat force outcome should return ForceOutcome2")
	}
}

func TestUtils_Shuffle(t *testing.T) {
	p := sgc7plugin.NewMockPlugin()
	p.SetCache([]int{2, 0, 1})
	out, err := Shuffle([]int{10, 20, 30}, p)
	if err != nil {
		t.Fatalf("Shuffle err: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("Shuffle len")
	}
	if out[0] != 30 || out[1] != 10 || out[2] != 20 {
		t.Fatalf("Shuffle order wrong: %v", out)
	}

	p2 := &errPlugin{}
	_, err = Shuffle([]int{1, 2, 3}, p2)
	if err == nil {
		t.Fatalf("Shuffle should return error on Random error")
	}
}

func TestUtils_PosAndSymbols(t *testing.T) {
	gs := &sgc7game.GameScene{Arr: [][]int{{1, 2, 3}, {4, 5, 6}}}
	ret := &sgc7game.Result{Pos: []int{0, 0, 1, 2}}
	syms := []int{1, 6}
	if !HasSymbolsInResult(gs, syms, ret) {
		t.Fatalf("HasSymbolsInResult should be true")
	}
	if CountSymbolsInResult(gs, syms, ret) != 2 {
		t.Fatalf("CountSymbolsInResult expected 2")
	}

	mapCodes := map[int]int{1: 2, 6: 3}
	if v := CalcSymbolsInResultEx(gs, mapCodes, ret, WRMTypeAddSymbolMulti); v != 5 {
		t.Fatalf("Add expected 5 got %d", v)
	}
	if v := CalcSymbolsInResultEx(gs, mapCodes, ret, WRMTypeMulSymbolMulti); v != 6 {
		t.Fatalf("Mul expected 6 got %d", v)
	}
	if v := CalcSymbolsInResultEx(gs, mapCodes, ret, WRMTypeSymbolMultiOnWays); v != 6 {
		t.Fatalf("Ways expected 6 got %d", v)
	}
}

func TestUtils_OtherSmallCases(t *testing.T) {
	if !IsValidPosWithHeight(0, 0, 2, 3, false) {
		t.Fatalf("IsValidPosWithHeight expected true")
	}
	pt := &sgc7game.PayTables{MapPay: map[int][]int{1: {10}, 2: {20}, 3: {30}}}
	if len(GetExcludeSymbols(pt, []int{2})) != 2 {
		t.Fatalf("GetExcludeSymbols size")
	}
	gs := &sgc7game.GameScene{Arr: [][]int{{1}}}
	ret := &sgc7game.Result{Pos: []int{0, 0}}
	if v := CalcSymbolsInResultEx(gs, map[int]int{}, ret, WRMTypeAddSymbolMulti); v != 1 {
		t.Fatalf("edge add")
	}
	if v := CalcSymbolsInResultEx(gs, map[int]int{}, ret, WRMTypeMulSymbolMulti); v != 1 {
		t.Fatalf("edge mul")
	}
	if HasSamePos([]int{0, 1}, []int{}) {
		t.Fatalf("HasSamePos target empty should false")
	}
	if CmpVal(1, "abc", 1) {
		t.Fatalf("unknown op should return false")
	}
	if !IsInPosArea(1, 1, []int{0, 2, 0, 2}) {
		t.Fatalf("IsInPosArea should true")
	}
}

// --- helpers for procSpin/Spin/GenDefaultScene tests ---
type fakeMod struct {
	calls int
	// sequence of PlayResult slices to return (one PlayResult per call)
	seq       [][]*sgc7game.PlayResult
	errOnPlay bool
}

func (fm *fakeMod) GetName() string { return BasicGameModName }
func (fm *fakeMod) OnPlay(game sgc7game.IGame, plugin sgc7plugin.IPlugin, cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, gameData any) (*sgc7game.PlayResult, error) {
	// return one PlayResult per call according to seq
	if fm.errOnPlay {
		return nil, errors.New("play failed")
	}
	idx := fm.calls
	fm.calls++
	if idx < len(fm.seq) && len(fm.seq[idx]) > 0 {
		// append prior prs elements to simulate sequence context
		// return the first PlayResult in that slot
		return fm.seq[idx][0], nil
	}
	// default empty finish
	pr := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr.IsFinish = true
	pr.CoinWin = 0
	return pr, nil
}

func makeSimpleGame(betMul int, mod *fakeMod) *Game {
	pool := &GamePropertyPool{
		MapGamePropPool: make(map[int]*sync.Pool),
		mapComponents:   nil,
		Config:          &Config{MapBetConfigs: map[int]*BetConfig{betMul: {}}, Bets: []int{betMul}},
	}

	// create sync.Pool that returns a minimal GameProperty
	pool.MapGamePropPool[betMul] = &sync.Pool{New: func() any {
		gp := &GameProperty{
			CurBetMul:       betMul,
			Pool:            pool,
			callStack:       NewCallStack(),
			PoolScene:       sgc7game.NewGameScenePoolEx(),
			SceneStack:      NewSceneStack(false),
			OtherSceneStack: NewSceneStack(true),
		}
		return gp
	}}

	g := &Game{
		BasicGame:    sgc7game.NewBasicGame(nil),
		Pool:         pool,
		MgrComponent: nil,
	}

	// attach fake mod to basic game
	g.MapGameMods = make(map[string]sgc7game.IGameMod)
	g.MapGameMods[BasicGameModName] = mod

	return g
}

func Test_procSpin_NewGameDataNil(t *testing.T) {
	// game with no MapGamePropPool entry -> NewGameData returns nil
	g := &Game{BasicGame: sgc7game.NewBasicGame(nil), Pool: &GamePropertyPool{MapGamePropPool: map[int]*sync.Pool{}}, MgrComponent: nil}

	ips := g.Initialize()
	stake := &sgc7game.Stake{CoinBet: 1, CashBet: 1}

	_, err := procSpin(g, ips, sgc7plugin.NewMockPlugin(), stake, "", "", false)
	if err == nil {
		t.Fatalf("procSpin should return error when NewGameData nil")
	}
}

func Test_Spin_ForceOutcome_FailAndSucceed(t *testing.T) {
	// create fake mod that first returns a winning PlayResult (CoinWin>0), then a zero-win PlayResult
	pr1 := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr1.IsFinish = true
	pr1.CoinWin = 1

	// success case with CoinWin == 0
	gs := &sgc7game.GameScene{Arr: [][]int{{1}}}
	pr2 := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr2.IsFinish = true
	pr2.CoinWin = 0
	pr2.Scenes = []*sgc7game.GameScene{gs}

	// fakeMod that always returns pr1 (non-zero win)
	fmFail := &fakeMod{seq: [][]*sgc7game.PlayResult{{pr1}}}
	gFail := makeSimpleGame(1, fmFail)

	SetAllowForceOutcome(1)
	// run Spin with FO enabled; since the PlayResult has CoinWin=1, FO script "totalWins == 0" is false
	_, err := Spin(gFail, gFail.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 1}, "", "", "totalWins == 0", false)
	if err == nil {
		t.Fatalf("Spin with force-outcome should fail when FO condition not met")
	}

	// now make a game that returns a zero-win result so FO succeeds
	fmSucc := &fakeMod{seq: [][]*sgc7game.PlayResult{{pr2}}}
	gSucc := makeSimpleGame(1, fmSucc)

	SetAllowForceOutcome(1)
	rets, err := Spin(gSucc, gSucc.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 1}, "", "", "totalWins == 0", false)
	if err != nil {
		t.Fatalf("Spin with force-outcome should succeed when FO condition met: %v", err)
	}
	if len(rets) == 0 {
		t.Fatalf("Spin expected results on success")
	}
}

func Test_GenDefaultScene_Basic(t *testing.T) {
	// GenDefaultScene calls procSpin in a loop until it finds a single PlayResult with CoinWin==0
	// build a fakeMod that returns pr1 (coinwin>0) then pr2 (coinwin==0)
	pr1 := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr1.IsFinish = true
	pr1.CoinWin = 5

	gs := &sgc7game.GameScene{Arr: [][]int{{9}}}
	pr2 := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr2.IsFinish = true
	pr2.CoinWin = 0
	pr2.Scenes = []*sgc7game.GameScene{gs}

	fm := &fakeMod{seq: [][]*sgc7game.PlayResult{{pr1}, {pr2}}}
	g := makeSimpleGame(1, fm)

	scene, err := GenDefaultScene(g, 1)
	if err != nil {
		t.Fatalf("GenDefaultScene failed: %v", err)
	}
	if scene == nil || scene.Arr[0][0] != 9 {
		t.Fatalf("GenDefaultScene returned wrong scene: %+v", scene)
	}
}

func Test_procSpin_PlayError(t *testing.T) {
	// fake mod that returns an error
	fm := &fakeMod{seq: nil, errOnPlay: true}
	g := makeSimpleGame(1, fm)

	_, err := procSpin(g, g.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 1}, "", "", false)
	if err == nil {
		t.Fatalf("procSpin should return error when Play errors")
	}
}

func Test_procSpin_NextCmds_Params(t *testing.T) {
	// first pr returns NextCmds and IsFinish=false, second pr finishes
	pr1 := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr1.IsFinish = false
	pr1.NextCmds = []string{"CMDX"}
	pr1.NextCmdParams = []string{"P1"}

	pr2 := sgc7game.NewPlayResult(BasicGameModName, 1, 0, "")
	pr2.IsFinish = true

	fm := &fakeMod{seq: [][]*sgc7game.PlayResult{{pr1}, {pr2}}}
	g := makeSimpleGame(1, fm)

	rets, err := procSpin(g, g.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 1}, "", "", false)
	if err != nil {
		t.Fatalf("procSpin NextCmds err: %v", err)
	}
	if len(rets) != 2 {
		t.Fatalf("expected 2 results got %d", len(rets))
	}
}

func Test_procSpin_IsWait_isNotAutoSelect(t *testing.T) {
	pr := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr.IsFinish = false
	pr.IsWait = true

	fm := &fakeMod{seq: [][]*sgc7game.PlayResult{{pr}}}
	g := makeSimpleGame(1, fm)

	rets, err := procSpin(g, g.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 1}, "", "", true)
	if err != nil {
		t.Fatalf("procSpin IsWait err: %v", err)
	}
	if len(rets) != 1 {
		t.Fatalf("expected 1 result on IsWait break got %d", len(rets))
	}
}

func Test_Spin_CheckStakeError(t *testing.T) {
	fm := &fakeMod{seq: nil}
	g := makeSimpleGame(1, fm)

	// stake with CashBet not in cfg.Bets
	_, err := Spin(g, g.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 999}, "", "", "", false)
	if err == nil {
		t.Fatalf("Spin should return error on invalid stake")
	}
}

// Note: not testing pr==nil branch because sgc7game.BasicGame.Play dereferences
// the returned PlayResult; using nil PlayResult causes panic in BasicGame.Play.

func Test_Spin_NoForceOutcome(t *testing.T) {
	pr := sgc7game.NewPlayResult(BasicGameModName, 0, 0, "")
	pr.IsFinish = true
	pr.CoinWin = 0

	fm := &fakeMod{seq: [][]*sgc7game.PlayResult{{pr}}}
	g := makeSimpleGame(1, fm)

	rets, err := Spin(g, g.Initialize(), sgc7plugin.NewMockPlugin(), &sgc7game.Stake{CoinBet: 1, CashBet: 1}, "", "", "", false)
	if err != nil {
		t.Fatalf("Spin no FO should succeed: %v", err)
	}
	if len(rets) != 1 {
		t.Fatalf("expected 1 result got %d", len(rets))
	}
}

func Test_SmallHelpers_Extra(t *testing.T) {
	if IsInitialArr([]int{1, 2, 3}) {
		t.Fatalf("IsInitialArr false case")
	}
	if IsSameIntArr([]int{1, 2}, []int{1, 3}) {
		t.Fatalf("IsSameIntArr false case")
	}

	if !CmpVal(5, "==", 5) {
		t.Fatalf("CmpVal == failed")
	}
	if !CmpVal(5, ">=", 5) {
		t.Fatalf(">= failed")
	}
	if !CmpVal(5, "<=", 5) {
		t.Fatalf("<= failed")
	}
	if !CmpVal(4, "<", 5) {
		t.Fatalf("< failed")
	}

	// IsValidPosWithHeight reversal case: expects false for these params
	if IsValidPosWithHeight(0, 1, 1, 3, true) {
		t.Fatalf("IsValidPosWithHeight reversal expected false for these params")
	}
	if !IsValidPosWithHeight(0, 0, 3, 3, false) {
		t.Fatalf("IsValidPosWithHeight curHeight>=height expected true")
	}
}
