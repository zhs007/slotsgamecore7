package lowcode

import (
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
)

func TestParseWinResultLimiterJSON(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w1",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "w1",
			"configuration": {
				"type": "maxonline",
				"srcComponents": ["src"]
			}
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseWinResultLimiter(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "w1", name)

	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	cfg := cfgIface.(*WinResultLimiterConfig)
	assert.Equal(t, "maxonline", cfg.StrType)
}

func TestOnMaxOnLineKeepsMax(t *testing.T) {
	// reuse setup to get gp and game params
	gp, gparam, pr := setupGameForPlay(t)

	// create two results on same line (index 0 and 1)
	pr.Results = []*sgc7game.Result{
		{CoinWin: 100, CashWin: 100, Mul: 1, Pos: []int{0, 0}, LineIndex: 0},
		{CoinWin: 150, CashWin: 150, Mul: 1, Pos: []int{0, 0}, LineIndex: 0},
	}

	// ensure component data for "src" tracks both results
	ccd := gp.GetComponentDataWithName("src").(*BasicComponentData)
	ccd.UsedResults = []int{0, 1}

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	cfg := &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "maxonline", SrcComponents: []string{"src"}}
	err := wrl.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	var icd IComponentData = &WinResultLimiterData{}
	nc, err := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Nil(t, err)
	assert.Equal(t, "", nc)

	// first result should be zeroed, second kept
	assert.Equal(t, 0, pr.Results[0].CoinWin)
	assert.Equal(t, 150, pr.Results[1].CoinWin)

	// component data wins should equal kept win
	cd := icd.(*WinResultLimiterData)
	assert.Equal(t, 150, cd.Wins)

	// PB build
	pb := cd.BuildPBComponentData()
	pbcd := pb.(*sgc7pb.WinResultLimiterData)
	assert.Equal(t, int32(150), pbcd.Wins)
}

func TestWinResultLimiterOnPlayGameNoHistory(t *testing.T) {
	gp, gparam, pr := setupGameForPlay(t)

	// prepare two results but clear history so component is skipped
	pr.Results = []*sgc7game.Result{{CoinWin: 100, LineIndex: 0}}
	ccd := gp.GetComponentDataWithName("src").(*BasicComponentData)
	ccd.UsedResults = []int{0}
	gparam.HistoryComponents = []string{}

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	cfg := &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "maxonline", SrcComponents: []string{"src"}}
	err := wrl.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	var icd IComponentData = &WinResultLimiterData{}
	nc, err := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)
}

func TestInitReadsYAML(t *testing.T) {
	content := "type: maxonline\nsrcComponents: [\"src\"]\n"
	fn := t.TempDir() + "/wrl.yaml"
	err := os.WriteFile(fn, []byte(content), 0644)
	assert.Nil(t, err)

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	err = wrl.Init(fn, makePoolWithPaytables())
	assert.Nil(t, err)
	assert.Equal(t, WinResultLimiterTypeName, wrl.Config.ComponentType)
}

func TestOnAsciiGamePrints(t *testing.T) {
	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	cd := &WinResultLimiterData{Wins: 42}
	// OnAsciiGame prints to stdout; invoking should not error
	err := wrl.OnAsciiGame(nil, nil, nil, nil, cd)
	assert.Nil(t, err)
}

func TestNewComponentDataAndOnNewGameAndGetValExFalse(t *testing.T) {
	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	icd := wrl.NewComponentData()
	assert.IsType(t, &WinResultLimiterData{}, icd)

	d := icd.(*WinResultLimiterData)
	// OnNewGame should initialize underlying maps
	d.OnNewGame(nil, nil)
	if d.MapConfigIntVals == nil {
		t.Fatalf("MapConfigIntVals should be initialized")
	}

	// GetValEx for unknown key
	v, ok := d.GetValEx("unknown", 0)
	assert.False(t, ok)
	assert.Equal(t, 0, v)
}

func TestOnPlayGameInvalidTypeBranch(t *testing.T) {
	gp, gparam, pr := setupGameForPlay(t)

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	// manually set invalid type to hit error path
	wrl.Config = &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "invalid", Type: WinResultLimiterType(999)}

	var icd IComponentData = &WinResultLimiterData{}
	nc, err := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrInvalidComponentConfig, err)
	assert.Equal(t, "", nc)
}

func TestDataCloneAndGetValExAndSetLink(t *testing.T) {
	d := &WinResultLimiterData{}
	d.Wins = 7
	d.MapConfigIntVals = map[string]int{"x": 1}

	c := d.Clone().(*WinResultLimiterData)
	assert.Equal(t, 7, c.Wins)

	// GetValEx for CVWins
	v, ok := c.GetValEx(CVWins, 0)
	assert.True(t, ok)
	assert.Equal(t, 7, v)

	// SetLinkComponent
	cfg := &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}}
	cfg.SetLinkComponent("next", "comp1")
	assert.Equal(t, "comp1", cfg.DefaultNextComponent)
}

func TestInitExInvalidTypeAndParseType(t *testing.T) {
	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)

	// passing wrong type should return ErrInvalidComponentConfig
	err := wrl.InitEx(&BasicComponentConfig{}, nil)
	assert.Equal(t, ErrInvalidComponentConfig, err)

	// parseWinResultLimiterType should be case-insensitive
	ty := parseWinResultLimiterType("MAXONLINE")
	assert.Equal(t, WRLTypeMaxOnLine, ty)
}

func TestOnMaxOnLine_NoUsedResultsAndZeroWins(t *testing.T) {
	gp, gparam, pr := setupGameForPlay(t)

	// ensure src component exists but has no UsedResults -> mapLinesWin empty
	ccd := gp.GetComponentDataWithName("src").(*BasicComponentData)
	ccd.UsedResults = nil

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	cfg := &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "maxonline", SrcComponents: []string{"src"}}
	err := wrl.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	var icd IComponentData = &WinResultLimiterData{}
	nc, err := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)

	// Now set UsedResults but all CoinWin are zero -> cd.Wins == 0 -> ErrComponentDoNothing
	pr.Results = []*sgc7game.Result{{CoinWin: 0, LineIndex: 0}}
	ccd.UsedResults = []int{0}

	nc2, err2 := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err2)
	assert.Equal(t, "", nc2)
}

func TestParseWinResultLimiter_NoComponentValues(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w_bad",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{"foo": {"bar": 1}}`)
	node, err := sonic.Get(js)
	assert.NoError(t, err)

	_, err = parseWinResultLimiter(betCfg, &node)
	assert.Error(t, err)
}

func TestParseWinResultLimiter_UnmarshalError(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w_bad2",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "w_bad2",
			"configuration": "not-an-object"
		}
	}`)
	node, err := sonic.Get(js)
	assert.NoError(t, err)

	_, err = parseWinResultLimiter(betCfg, &node)
	assert.Error(t, err)
}

func TestParseWinResultLimiter_DataNode(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "w2",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"data": {
			"label": "w2",
			"configuration": {
				"type": "maxonline",
				"srcComponents": ["src"]
			}
		}
	}`)
	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseWinResultLimiter(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "w2", name)

	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	cfg := cfgIface.(*WinResultLimiterConfig)
	assert.Equal(t, "maxonline", cfg.StrType)
}

func TestInit_ReadFileError(t *testing.T) {
	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	// point to a non-existent file
	err := wrl.Init("/no/such/file/doesnotexist.yaml", makePoolWithPaytables())
	assert.Error(t, err)
}

func TestOnMaxOnLine_SingleResultPerLine(t *testing.T) {
	gp, gparam, pr := setupGameForPlay(t)

	// two results on different lines -> each list len == 1 -> should do nothing
	pr.Results = []*sgc7game.Result{
		{CoinWin: 100, CashWin: 100, Mul: 1, Pos: []int{0, 0}, LineIndex: 0},
		{CoinWin: 200, CashWin: 200, Mul: 1, Pos: []int{0, 0}, LineIndex: 1},
	}

	ccd := gp.GetComponentDataWithName("src").(*BasicComponentData)
	ccd.UsedResults = []int{0, 1}

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	cfg := &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "maxonline", SrcComponents: []string{"src"}}
	err := wrl.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	var icd IComponentData = &WinResultLimiterData{}
	nc, err := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)
}

func TestOnMaxOnLine_ResultIndexOutOfRange(t *testing.T) {
	gp, gparam, pr := setupGameForPlay(t)

	// pr.Results length is 1 but UsedResults references index 5 -> out of range
	pr.Results = []*sgc7game.Result{{CoinWin: 100, LineIndex: 0}}
	ccd := gp.GetComponentDataWithName("src").(*BasicComponentData)
	ccd.UsedResults = []int{5}

	wrl := NewWinResultLimiter("w1").(*WinResultLimiter)
	cfg := &WinResultLimiterConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "maxonline", SrcComponents: []string{"src"}}
	err := wrl.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	var icd IComponentData = &WinResultLimiterData{}
	nc, err := wrl.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	// out of range entries are skipped and result set becomes empty -> do nothing
	assert.Equal(t, ErrComponentDoNothing, err)
	assert.Equal(t, "", nc)
}
