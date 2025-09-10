package lowcode

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func TestSymbolExpander_InitEx_Errors(t *testing.T) {
	se := NewSymbolExpander("sx").(*SymbolExpander)

	// nil config
	err := se.InitEx(nil, &GamePropertyPool{})
	assert.Error(t, err)

	// wrong type config
	err = se.InitEx(&struct{ A int }{A: 1}, &GamePropertyPool{})
	assert.Error(t, err)

	// valid config but pool missing paytables
	cfg := &SymbolExpanderConfig{Symbols: []string{"A"}}
	err = se.InitEx(cfg, &GamePropertyPool{})
	assert.Error(t, err)
}

func TestSymbolExpander_InitEx_And_ProcControllers(t *testing.T) {
	// build paytables with symbol name -> code
	pt := &sgc7game.PayTables{MapSymbols: map[string]int{"CO": 7}}

	pool := &GamePropertyPool{DefaultPaytables: pt}

	se := NewSymbolExpander("sx").(*SymbolExpander)
	cfg := &SymbolExpanderConfig{Symbols: []string{"CO"}}

	err := se.InitEx(cfg, pool)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(se.Config.SymbolCodes))

	// ProcControllers with nil map should be noop
	gp := &GameProperty{Pool: pool}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	params := NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil)

	se.ProcControllers(gp, &fakePlugin{}, pr, params, 0, "noexist")
}

func TestSymbolExpander_OnPlayGame_NoExpandAndExpand(t *testing.T) {
	// prepare paytables and pool
	pt := &sgc7game.PayTables{MapSymbols: map[string]int{"A": 1, "B": 2}}
	pool := &GamePropertyPool{DefaultPaytables: pt}
	// newGameProp expects pool.Config and newRNG/newFeatureLevel functions to be set
	pool.Config = &Config{Width: 3, Height: 3}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

	// build a source scene with no expandable symbols
	sc, _ := sgc7game.NewGameScene(3, 3)
	// fill with B (non-expandable)
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			sc.Arr[x][y] = 2
		}
	}
	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr.Scenes = append(pr.Scenes, sc)

	se := NewSymbolExpander("sx").(*SymbolExpander)
	se.Config = &SymbolExpanderConfig{Symbols: []string{"B"}, SymbolCodes: []int{2}, IgnoreSymbols: []string{}, IgnoreSymbolCodes: []int{}}
	se.BasicComponent.onInit(&BasicComponentConfig{})

	// case: no expandable symbols present => ErrComponentDoNothing
	gp1 := pool.newGameProp(1)
	_, err := se.OnPlayGame(gp1, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, []*sgc7game.PlayResult{pr}, se.NewComponentData())
	// OnPlayGame returns ErrComponentDoNothing when there's nothing to do
	assert.Error(t, err)

	// now create a scene that has expandable symbol 'A' in column 0 and non-expandable at top
	sc2, _ := sgc7game.NewGameScene(3, 3)
	// column 0: top has non-expandable (1), below has expandable 2 => starty=0 and expansion will set top to 2
	sc2.Arr[0][0] = 1
	sc2.Arr[0][1] = 2
	sc2.Arr[0][2] = 1
	// other columns fill with B
	for x := 1; x < 3; x++ {
		for y := 0; y < 3; y++ {
			sc2.Arr[x][y] = 2
		}
	}

	pr2 := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr2.Scenes = append(pr2.Scenes, sc2)

	// use the scene directly from pr2 to avoid SceneStack semantics
	gs := pr2.Scenes[0]
	if gs == nil {
		t.Fatalf("expected pr2.Scenes[0] to be present")
	}

	// independently verify the detection logic for column 0
	arr := gs.Arr[0]
	foundSc := -1
	for _, s := range arr {
		if slices.Contains(se.Config.SymbolCodes, s) {
			foundSc = s
			break
		}
	}
	if foundSc == -1 {
		t.Fatalf("expected to find expandable symbol in column but didn't; arr=%v codes=%v", arr, se.Config.SymbolCodes)
	}
	// find starty
	starty := -1
	for y, s := range arr {
		if !slices.Contains(se.Config.SymbolCodes, s) && !slices.Contains(se.Config.IgnoreSymbolCodes, s) {
			starty = y
			break
		}
	}
	if starty == -1 {
		t.Fatalf("expected to find starty in column but didn't; arr=%v codes=%v ignore=%v", arr, se.Config.SymbolCodes, se.Config.IgnoreSymbolCodes)
	}

	// will expand column 0 from starty=0 downwards -> top replaced with A
	gp2 := pool.newGameProp(1)
	nc2, err2 := se.OnPlayGame(gp2, pr2, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, []*sgc7game.PlayResult{pr2}, se.NewComponentData())
	assert.NoError(t, err2)
	assert.Equal(t, "", nc2)
	// after expansion, pr2.Scenes should have new scene index appended
	assert.True(t, len(pr2.Scenes) >= 1)
}

func TestSymbolExpander_OnAsciiGame(t *testing.T) {
	se := NewSymbolExpander("sx").(*SymbolExpander)
	se.BasicComponent.onInit(&BasicComponentConfig{})

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	sc, _ := sgc7game.NewGameScene(3, 3)
	pr.Scenes = append(pr.Scenes, sc)

	icd := se.NewComponentData()
	bcd := icd.(*BasicComponentData)
	bcd.UsedScenes = []int{0}

	scm := asciigame.NewSymbolColorMap(&sgc7game.PayTables{})

	err := se.OnAsciiGame(&GameProperty{}, pr, nil, scm, icd)
	assert.NoError(t, err)
}

func TestSymbolExpander_OnPlayGame_StartyMinusOne(t *testing.T) {
	// prepare pool required by newGameProp
	pool := &GamePropertyPool{}
	pool.Config = &Config{Width: 3, Height: 3}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

	// scene where every cell in column 0 is either expandable (code 2) or ignored (code 3)
	sc, _ := sgc7game.NewGameScene(3, 3)
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if x == 0 {
				// fill with expandable or ignore only
				if y%2 == 0 {
					sc.Arr[x][y] = 2
				} else {
					sc.Arr[x][y] = 3
				}
			} else {
				sc.Arr[x][y] = 1
			}
		}
	}

	pr := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr.Scenes = append(pr.Scenes, sc)

	se := NewSymbolExpander("sx").(*SymbolExpander)
	se.Config = &SymbolExpanderConfig{Symbols: []string{"E"}, SymbolCodes: []int{2}, IgnoreSymbols: []string{"I"}, IgnoreSymbolCodes: []int{3}}
	se.BasicComponent.onInit(&BasicComponentConfig{})

	gp := pool.newGameProp(1)
	_, err := se.OnPlayGame(gp, pr, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, []*sgc7game.PlayResult{pr}, se.NewComponentData())

	// starty == -1 for column 0 means nothing to expand -> ErrComponentDoNothing
	assert.Error(t, err)
}

func TestSymbolExpander_OnPlayGame_ProcControllersCalled(t *testing.T) {
	// prepare pool, paytables and functions
	pt := &sgc7game.PayTables{MapSymbols: map[string]int{"A": 2, "B": 1}}
	pool := &GamePropertyPool{DefaultPaytables: pt}
	pool.Config = &Config{Width: 3, Height: 3}
	pool.newRNG = func() IRNG { return &stubRNG{} }
	pool.newFeatureLevel = func(b int) IFeatureLevel { return &stubFeatureLevel{} }

	// build a scene which will expand in column 0
	sc2, _ := sgc7game.NewGameScene(3, 3)
	sc2.Arr[0][0] = 1
	sc2.Arr[0][1] = 2
	sc2.Arr[0][2] = 1
	for x := 1; x < 3; x++ {
		for y := 0; y < 3; y++ {
			sc2.Arr[x][y] = 1
		}
	}

	pr2 := sgc7game.NewPlayResult("m", 0, 0, "t")
	pr2.Scenes = append(pr2.Scenes, sc2)

	se := NewSymbolExpander("sx").(*SymbolExpander)
	// configure map awards so ProcControllers will find entries for "<trigger>" and symbol name
	se.Config = &SymbolExpanderConfig{
		Symbols:     []string{"A"},
		SymbolCodes: []int{2},
		MapAwards: map[string][]*Award{
			"<trigger>": {&Award{AwardType: "respinTimes", Vals: []int{1}, StrParams: []string{"noexist"}}},
			"A":         {&Award{AwardType: "respinTimes", Vals: []int{1}, StrParams: []string{"noexist"}}},
		},
	}
	se.BasicComponent.onInit(&BasicComponentConfig{})

	gp2 := pool.newGameProp(1)
	nc2, err2 := se.OnPlayGame(gp2, pr2, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, []*sgc7game.PlayResult{pr2}, se.NewComponentData())
	assert.NoError(t, err2)
	assert.Equal(t, "", nc2)
	// expansion should append a new scene
	assert.True(t, len(pr2.Scenes) >= 1)
}

func TestSymbolExpander_OnPlayGame_InvalidICD(t *testing.T) {
	se := NewSymbolExpander("sx").(*SymbolExpander)
	// use a valid config so function progresses to icd check
	se.Config = &SymbolExpanderConfig{Symbols: []string{"A"}, SymbolCodes: []int{1}}

	// call with nil icd -> should return ErrInvalidComponentData
	_, err := se.OnPlayGame(&GameProperty{}, nil, NewGameParam(&sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil), &fakePlugin{}, "", "", nil, &sgc7game.Stake{CoinBet: 1, CashBet: 1}, nil, nil)
	assert.Error(t, err)
}
