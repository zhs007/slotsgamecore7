package lowcode

import (
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func TestWinResultModifierExDataBasics(t *testing.T) {
	d := &WinResultModifierExData{}
	// onNewStep sets defaults
	d.Wins = 5
	d.WinMulti = 7
	d.onNewStep()
	assert.Equal(t, 0, d.Wins)
	assert.Equal(t, 1, d.WinMulti)

	// Clone copies values
	d.Wins = 3
	d.WinMulti = 4
	d.BasicComponentData.MapConfigIntVals = map[string]int{"x": 1}
	cd := d.Clone().(*WinResultModifierExData)
	assert.Equal(t, 3, cd.Wins)
	assert.Equal(t, 4, cd.WinMulti)

	// BuildPBComponentData and GetValEx
	pb := d.BuildPBComponentData()
	assert.NotNil(t, pb)

	v, ok := d.GetValEx(CVWins, 0)
	assert.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestJsonWinResultModifierExBuild(t *testing.T) {
	j := &jsonWinResultModifierEx{
		Type:             "addsymbolmulti",
		SourceComponents: []string{"src"},
		MapTargetSymbols: [][]any{{"A", float64(2)}, {"B", "3"}},
	}

	cfg := j.build()
	assert.Equal(t, "addsymbolmulti", cfg.StrType)
	assert.Equal(t, 2, cfg.MapTargetSymbols["A"])
	assert.Equal(t, 3, cfg.MapTargetSymbols["B"])
}

func TestInitExAndOnPlayAddSymbolMulti(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	cfg := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "addsymbolmulti", SourceComponents: []string{"src"}, MapTargetSymbols: map[string]int{"A": 2}}
	err := we.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	// prepare game and scene with symbol code 1 at (0,0)
	gp, gparam, pr := setupGameForPlay(t)
	gs, _ := sgc7game.NewGameScene(1, 1)
	gs.Arr[0][0] = 1
	gp.SceneStack.Push("src", gs)

	icd := &WinResultModifierExData{}
	nc, err := we.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Nil(t, err)
	assert.Equal(t, "", nc)
	// coin win should be multiplied by map value 2
	assert.Equal(t, 200, pr.Results[0].CoinWin)
}

func TestOnPlayGameSymbolMultiOnWaysDivisionByZero(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	cfg := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "symbolmultionways", SourceComponents: []string{"src"}, MapTargetSymbols: map[string]int{"A": 2}}
	err := we.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	gp, gparam, pr := setupGameForPlay(t)
	gs, _ := sgc7game.NewGameScene(1, 1)
	gs.Arr[0][0] = 1
	gp.SceneStack.Push("src", gs)

	// set Mul to zero to trigger division-by-zero protection
	pr.Results[0].Mul = 0

	icd := &WinResultModifierExData{}
	nc, err := we.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Equal(t, ErrInvalidComponentConfig, err)
	assert.Equal(t, "", nc)
}

func TestOnPlayGameTypeMismatchAndAscii(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	cfg := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "addsymbolmulti", SourceComponents: []string{"src"}, MapTargetSymbols: map[string]int{"A": 2}}
	err := we.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	// wrong icd type for OnPlayGame should return ErrInvalidComponentData
	_, err = we.OnPlayGame(nil, nil, nil, nil, "", "", nil, nil, nil, &WinResultModifierData{})
	assert.Equal(t, ErrInvalidComponentData, err)

	// OnAsciiGame with wrong icd should return ErrInvalidComponentData
	err = we.OnAsciiGame(nil, nil, nil, nil, &WinResultModifierData{})
	assert.Equal(t, ErrInvalidComponentData, err)
}

func TestParseWinResultModifierExJSON(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "wex",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "wex",
			"configuration": {
				"type": "addSymbolMulti",
				"sourceComponent": ["src"],
				"mapTargetSymbols": [["A", 2]]
			}
		}
	}`)

	node, err := sonic.Get(js)
	assert.NoError(t, err)

	name, err := parseWinResultModifierEx(betCfg, &node)
	assert.NoError(t, err)
	assert.Equal(t, "wex", name)

	cfgIface, ok := betCfg.mapConfig[name]
	assert.True(t, ok)
	cfg := cfgIface.(*WinResultModifierExConfig)
	assert.Equal(t, 2, cfg.MapTargetSymbols["A"])
}

func TestInitReadsYAMLFileEx(t *testing.T) {
	content := "type: addsymbolmulti\nsourceComponents: [\"src\"]\nmapTargetSymbols: {A: 2}\n"
	fn := t.TempDir() + "/wre.yaml"
	err := os.WriteFile(fn, []byte(content), 0644)
	assert.Nil(t, err)

	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	err = we.Init(fn, makePoolWithPaytables())
	assert.Nil(t, err)
	// MapTargetSymbolCodes should be filled based on paytable
	_, ok := we.Config.MapTargetSymbolCodes[1]
	assert.True(t, ok)
}

func TestAdditionalWinResultModifierExPaths(t *testing.T) {
	// OnNewGame should initialize maps in BasicComponentData
	d := &WinResultModifierExData{}
	d.MapConfigIntVals = nil
	d.MapConfigVals = nil
	d.OnNewGame(nil, nil)
	if d.MapConfigIntVals == nil || d.MapConfigVals == nil {
		t.Fatalf("OnNewGame did not initialize maps")
	}

	// SetLinkComponent should set DefaultNextComponent
	cfg := &WinResultModifierExConfig{}
	cfg.SetLinkComponent("next", "comp1")
	if cfg.DefaultNextComponent != "comp1" {
		t.Fatalf("SetLinkComponent failed to set default next component")
	}

	// NewComponentData returns correct type
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	cd := we.NewComponentData()
	if _, ok := cd.(*WinResultModifierExData); !ok {
		t.Fatalf("NewComponentData returned wrong type")
	}

	// OnAsciiGame success path
	icd := &WinResultModifierExData{Wins: 10, WinMulti: 2}
	err := we.OnAsciiGame(nil, nil, nil, nil, icd)
	if err != nil {
		t.Fatalf("OnAsciiGame failed on valid icd: %v", err)
	}

	// GetValEx negative case
	v, ok := icd.GetValEx("notwins", 0)
	if ok || v != 0 {
		t.Fatalf("GetValEx returned unexpected value")
	}

	// InitEx invalid type should return ErrInvalidComponentConfig
	cfg2 := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "existsymbol"}
	err = we.InitEx(cfg2, makePoolWithPaytables())
	if err != ErrInvalidComponentConfig {
		t.Fatalf("InitEx did not reject invalid type for Ex: %v", err)
	}
}

func TestInitErrorPathsEx(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	// non-existent file should error
	err := we.Init("/no/such/file.yaml", makePoolWithPaytables())
	if err == nil {
		t.Fatalf("expected error for non-existent file")
	}
}

func TestInitExInvalidTargetSymbolEx(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	cfg := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "addsymbolmulti", SourceComponents: []string{"src"}, MapTargetSymbols: map[string]int{"Z": 2}}
	err := we.InitEx(cfg, makePoolWithPaytables())
	if err != ErrInvalidSymbol {
		t.Fatalf("expected ErrInvalidSymbol, got %v", err)
	}
}

func TestOnPlayGameMulOneDoesNothing(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	// map target symbols does not include A, so mul should be 1
	// use empty mapping so no matched symbols -> mul == 1
	cfg := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "addsymbolmulti", SourceComponents: []string{"src"}, MapTargetSymbols: map[string]int{}}
	err := we.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	gp, gparam, pr := setupGameForPlay(t)
	// push a scene so GetTargetScene3 returns non-nil
	gs, _ := sgc7game.NewGameScene(1, 1)
	gs.Arr[0][0] = 0
	gp.SceneStack.Push("src", gs)
	icd := &WinResultModifierExData{}
	nc, err := we.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	if err != ErrComponentDoNothing {
		t.Fatalf("expected ErrComponentDoNothing, got %v", err)
	}
	if nc != "" {
		t.Fatalf("expected empty next component, got %v", nc)
	}
}

func TestOnPlayGameSymbolMultiOnWaysSuccess(t *testing.T) {
	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	cfg := &WinResultModifierExConfig{BasicComponentConfig: BasicComponentConfig{}, StrType: "symbolmultionways", SourceComponents: []string{"src"}, MapTargetSymbols: map[string]int{"A": 5}}
	err := we.InitEx(cfg, makePoolWithPaytables())
	assert.Nil(t, err)

	gp, gparam, pr := setupGameForPlay(t)
	gs, _ := sgc7game.NewGameScene(1, 1)
	gs.Arr[0][0] = 1
	gp.SceneStack.Push("src", gs)

	// ensure Mul is positive
	pr.Results[0].Mul = 1

	icd := &WinResultModifierExData{}
	nc, err := we.OnPlayGame(gp, pr, gparam, nil, "", "", nil, nil, nil, icd)
	assert.Nil(t, err)
	assert.Equal(t, "", nc)
	// coin win should be scaled via symbol multi on ways logic
	// original coinwin 100 -> divided by Mul(1) then * mul (5) => 500
	assert.Equal(t, 500, pr.Results[0].CoinWin)
}

func TestJsonWinResultModifierExBuildEdgeCases(t *testing.T) {
	j := &jsonWinResultModifierEx{
		Type:             "addsymbolmulti",
		MapTargetSymbols: [][]any{{"A"}, {123, 2}},
	}
	cfg := j.build()
	// nothing should be added for invalid rows
	if len(cfg.MapTargetSymbols) != 0 {
		t.Fatalf("expected empty MapTargetSymbols, got %v", cfg.MapTargetSymbols)
	}
}

func TestParseWinResultModifierEx_NoComponentValues(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "wex_bad",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{"foo": {"bar": 1}}`)
	node, err := sonic.Get(js)
	if err != nil {
		t.Fatalf("sonic.Get error: %v", err)
	}

	_, err = parseWinResultModifierEx(betCfg, &node)
	if err == nil {
		t.Fatalf("expected error when no componentValues present")
	}
}

func TestJsonWinResultModifierExBuildNumbers(t *testing.T) {
	j := &jsonWinResultModifierEx{
		Type:             "addsymbolmulti",
		MapTargetSymbols: [][]any{{"I", 7}, {"I64", int64(9)}, {"F", float64(2.0)}},
	}

	cfg := j.build()
	if cfg.MapTargetSymbols["I"] != 7 {
		t.Fatalf("expected I=7 got %v", cfg.MapTargetSymbols["I"])
	}
	if cfg.MapTargetSymbols["I64"] != 9 {
		t.Fatalf("expected I64=9 got %v", cfg.MapTargetSymbols["I64"])
	}
	if cfg.MapTargetSymbols["F"] != 2 {
		t.Fatalf("expected F=2 got %v", cfg.MapTargetSymbols["F"])
	}
}

func TestJsonWinResultModifierExBuildNonNumericString(t *testing.T) {
	j := &jsonWinResultModifierEx{
		Type:             "addsymbolmulti",
		MapTargetSymbols: [][]any{{"A", "not-a-number"}, {"B"}},
	}

	cfg := j.build()
	// neither entry should be added
	if len(cfg.MapTargetSymbols) != 0 {
		t.Fatalf("expected empty MapTargetSymbols, got %v", cfg.MapTargetSymbols)
	}
}

func TestParseWinResultModifierEx_UnmarshalError(t *testing.T) {
	betCfg := &BetConfig{
		Bet:            1,
		Start:          "wex",
		Components:     []*ComponentConfig{},
		mapConfig:      make(map[string]IComponentConfig),
		mapBasicConfig: make(map[string]*BasicComponentConfig),
	}

	js := []byte(`{
		"componentValues": {
			"label": "wex_bad",
			"configuration": {
				"type": "addSymbolMulti",
				"sourceComponent": ["src"],
				"mapTargetSymbols": "not-an-object"
			}
		}
	}`)

	node, err := sonic.Get(js)
	if err != nil {
		t.Fatalf("sonic.Get error: %v", err)
	}

	_, err = parseWinResultModifierEx(betCfg, &node)
	if err == nil {
		t.Fatalf("expected Unmarshal error from parseWinResultModifierEx")
	}
}

func TestInitReadsYAMLFileExUnmarshalError(t *testing.T) {
	// write invalid YAML to temp file
	content := "type: [unclosed"
	fn := t.TempDir() + "/wre_bad.yaml"
	err := os.WriteFile(fn, []byte(content), 0644)
	if err != nil {
		t.Fatalf("write tmp file err: %v", err)
	}

	we := NewWinResultModifierEx("wex").(*WinResultModifierEx)
	err = we.Init(fn, makePoolWithPaytables())
	if err == nil {
		t.Fatalf("expected unmarshal error for invalid yaml")
	}
}
