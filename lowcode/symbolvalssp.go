package lowcode

import (
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

// SymbolValsSPTypeName is the component type name used in configuration
// and component registration for the SymbolValsSP component.
const SymbolValsSPTypeName = "symbolValsSP"

// SymbolValsSPType defines the runtime variant of the SymbolValsSP
// component. Currently only a normal type is supported. This type
// is used to switch behavior in OnPlayGame.
type SymbolValsSPType int

const (
	SVSPTypeNormal SymbolValsSPType = 0 // normal
)

// parseSymbolValsSPType parses a textual type representation into
// SymbolValsSPType. Unknown or empty values default to SVSPTypeNormal.
func parseSymbolValsSPType(_ string) SymbolValsSPType {
	return SVSPTypeNormal
}

// SymbolValsSPMultiType defines how the multi (multiplier) logic
// should be processed. Values affect procMulti.
type SymbolValsSPMultiType int

const (
	SVSPMultiTypeNormal SymbolValsSPMultiType = 0 // normal
	SVSPMultiTypeRound  SymbolValsSPMultiType = 1 // round
)

// parseSymbolValsSPMultiType converts a string to SymbolValsSPMultiType.
// Supported values: "round" -> SVSPMultiTypeRound, otherwise normal.
func parseSymbolValsSPMultiType(str string) SymbolValsSPMultiType {
	if str == "round" {
		return SVSPMultiTypeRound
	}

	return SVSPMultiTypeNormal
}

// SymbolValsSPCollectType defines the collect behavior used when
// collect symbols are present. Sequence means only the first collect
// position is used; normal will use all collect positions.
type SymbolValsSPCollectType int

const (
	SVSPCollectTypeNormal   SymbolValsSPCollectType = 0 // normal
	SVSPCollectTypeSequence SymbolValsSPCollectType = 1 // sequence
)

// parseSymbolValsSPCollectType converts a string to SymbolValsSPCollectType.
// Supported values: "sequence" -> SVSPCollectTypeSequence, otherwise normal.
func parseSymbolValsSPCollectType(str string) SymbolValsSPCollectType {
	if str == "sequence" {
		return SVSPCollectTypeSequence
	}

	return SVSPCollectTypeNormal
}

// SymbolValsSPData holds per-play state for the SymbolValsSP component.
// It tracks positions (Pos), scenes used, and counters produced by the
// multi and collect processing so later components or output can access
// the produced values.
type SymbolValsSPData struct {
	BasicComponentData
	// Pos stores positions in the form [x,y, x,y, -1, x,y,...]. -1 is a
	// separator inserted by newLine when a logical group changes (for
	// example between multi and collect results).
	Pos                  []int
	// UsedScenes and UsedOtherScenes store indices into PlayResult Scenes
	// and OtherScenes that were added during processing. They are used
	// by OnAsciiGame for display. -1 separators may appear in these slices.
	UsedScenes           []int // 使用的场景, -1分隔
	MultiSymbolNum       int   // number of multi symbols detected
	MultiCoinSymbolNum   int   // number of coin symbols affected by multi
	Multi                int   // computed multiplier value
	CollectSymbolNum     int   // number of collect symbols detected
	CollectCoinSymbolNum int   // number of coin symbols collected
	CollectCoin          int   // total coin amount collected
}

// OnNewGame -
// OnNewGame resets per-game state by delegating to BasicComponentData.OnNewGame.
func (symbolValsSPData *SymbolValsSPData) OnNewGame(gameProp *GameProperty, component IComponent) {
	symbolValsSPData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
// onNewStep clears state used within a single play step. It resets
// counters and position lists so the component can accumulate fresh data.
func (symbolValsSPData *SymbolValsSPData) onNewStep() {
	symbolValsSPData.UsedScenes = nil
	symbolValsSPData.UsedOtherScenes = nil

	symbolValsSPData.Pos = nil
	symbolValsSPData.MultiSymbolNum = 0
	symbolValsSPData.MultiCoinSymbolNum = 0
	// default multiplier is 1
	symbolValsSPData.Multi = 1
	symbolValsSPData.CollectSymbolNum = 0
	symbolValsSPData.CollectCoinSymbolNum = 0
	symbolValsSPData.CollectCoin = 0
}

// Clone
func (symbolValsSPData *SymbolValsSPData) Clone() IComponentData {
	target := &SymbolValsSPData{
		BasicComponentData:   symbolValsSPData.CloneBasicComponentData(),
		MultiSymbolNum:       symbolValsSPData.MultiSymbolNum,
		MultiCoinSymbolNum:   symbolValsSPData.MultiCoinSymbolNum,
		Multi:                symbolValsSPData.Multi,
		CollectSymbolNum:     symbolValsSPData.CollectSymbolNum,
		CollectCoinSymbolNum: symbolValsSPData.CollectCoinSymbolNum,
		CollectCoin:          symbolValsSPData.CollectCoin,
		Pos:                  slices.Clone(symbolValsSPData.Pos),
	}

	return target
}

// BuildPBComponentData
func (symbolValsSPData *SymbolValsSPData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolValsSPData{
		BasicComponentData:   symbolValsSPData.BuildPBBasicComponentData(),
		Pos:                  make([]int32, 0, len(symbolValsSPData.Pos)),
		MultiSymbolNum:       int32(symbolValsSPData.MultiSymbolNum),
		MultiCoinSymbolNum:   int32(symbolValsSPData.MultiCoinSymbolNum),
		Multi:                int32(symbolValsSPData.Multi),
		CollectSymbolNum:     int32(symbolValsSPData.CollectSymbolNum),
		CollectCoinSymbolNum: int32(symbolValsSPData.CollectCoinSymbolNum),
		CollectCoin:          int32(symbolValsSPData.CollectCoin),
	}

	for _, v := range symbolValsSPData.Pos {
		pbcd.Pos = append(pbcd.Pos, int32(v))
	}

	return pbcd
}

// GetPos -
func (symbolValsSPData *SymbolValsSPData) GetPos() []int {
	return symbolValsSPData.Pos
}

// HasPos -
func (symbolValsSPData *SymbolValsSPData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(symbolValsSPData.Pos, x, y, 0) >= 0
}

// AddPos -
func (symbolValsSPData *SymbolValsSPData) AddPos(x int, y int) {
	symbolValsSPData.Pos = append(symbolValsSPData.Pos, x, y)
}

// ClearPos -
func (symbolValsSPData *SymbolValsSPData) ClearPos() {
	symbolValsSPData.Pos = nil
}

// AddPosEx -
func (symbolValsSPData *SymbolValsSPData) AddPosEx(x int, y int) {
	if !symbolValsSPData.HasPos(x, y) {
		symbolValsSPData.AddPos(x, y)
	}
}

// newLine -
func (symbolValsSPData *SymbolValsSPData) newLine() {
	if len(symbolValsSPData.Pos) > 0 {
		symbolValsSPData.Pos = append(symbolValsSPData.Pos, -1)
	}
}

// SymbolValsSPConfig - configuration for SymbolValsSP
// SymbolValsSPConfig is the YAML/JSON-backed configuration for the
// SymbolValsSP component. Human friendly fields (strings, symbol
// names) are provided and later resolved to integer symbol codes during
// InitEx (stored in fields with Code suffix).
type SymbolValsSPConfig struct {
	BasicComponentConfig    `yaml:",inline" json:",inline"`
	StrType                 string                  `yaml:"type" json:"type"`
	Type                    SymbolValsSPType        `yaml:"-" json:"-"`
	// CoinSymbols are the symbols considered "coin" in the OtherScene.
	CoinSymbols             []string                `yaml:"coinSymbols" json:"coinSymbols"`
	CoinSymbolCodes         []int                   `yaml:"-" json:"-"`
	// MultiSymbols are symbols that produce multiplier effects using
	// values from OtherScene.
	MultiSymbols            []string                `yaml:"multiSymbols" json:"multiSymbols"`
	MultiSymbolCodes        []int                   `yaml:"-" json:"-"`
	StrMultiType            string                  `yaml:"multiType" json:"multiType"`
	MultiType               SymbolValsSPMultiType   `yaml:"-" json:"-"`
	MultiTargetSymbol       string                  `yaml:"multiTargetSymbol" json:"multiTargetSymbol"`
	MultiTargetSymbolCode   int                     `yaml:"-" json:"-"`
	// CollectSymbols are symbols that will collect coin values into
	// a target position when present.
	CollectSymbols          []string                `yaml:"collectSymbols" json:"collectSymbols"`
	CollectSymbolCodes      []int                   `yaml:"-" json:"-"`
	CollectTargetSymbol     string                  `yaml:"collectTargetSymbol" json:"collectTargetSymbol"`
	CollectTargetSymbolCode int                     `yaml:"-" json:"-"`
	CollectCoinSymbol       string                  `yaml:"collectCoinSymbol" json:"collectCoinSymbol"`
	CollectCoinSymbolCode   int                     `yaml:"-" json:"-"`
	CollectMultiSymbol      string                  `yaml:"collectMultiSymbol" json:"collectMultiSymbol"`
	CollectMultiSymbolCode  int                     `yaml:"-" json:"-"`
	StrCollectType          string                  `yaml:"collectType" json:"collectType"`
	CollectType             SymbolValsSPCollectType `yaml:"-" json:"-"`
	// MapAwards maps controller labels to award lists executed on triggers
	// like <trigger>, <multi>, <collect>.
	MapAwards               map[string][]*Award     `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
// SetLinkComponent sets named link components. Currently supports
// "next" to set the default next component name in the chain.
func (cfg *SymbolValsSPConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	}
}

// SymbolValsSP is a component that inspects the main GameScene and a
// corresponding OtherScene to implement "multi" and "collect" mechanics.
// It produces modified scenes/other-scenes and records positional
// information in SymbolValsSPData.
type SymbolValsSP struct {
	*BasicComponent `json:"-"`
	Config          *SymbolValsSPConfig `json:"config"`
}

// Init -
// Init reads a YAML file at fn and unmarshals it into the component
// configuration then delegates to InitEx.
func (symbolValsSP *SymbolValsSP) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolValsSP.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolValsSPConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolValsSP.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return symbolValsSP.InitEx(cfg, pool)
}

// InitEx -
// InitEx initializes the component from an already-unmarshaled config
// object. It resolves symbol names to symbol codes using pool.DefaultPaytables
// and validates required fields.
func (symbolValsSP *SymbolValsSP) InitEx(cfg any, pool *GamePropertyPool) error {
	cfgp, ok := cfg.(*SymbolValsSPConfig)
	if !ok || cfgp == nil {
		goutils.Error("SymbolValsSP.InitEx:invalid cfg",
			goutils.Err(ErrInvalidComponentConfig))
		return ErrInvalidComponentConfig
	}

	symbolValsSP.Config = cfgp
	symbolValsSP.Config.ComponentType = SymbolValsSPTypeName

	symbolValsSP.Config.Type = parseSymbolValsSPType(symbolValsSP.Config.StrType)

	if len(symbolValsSP.Config.MultiSymbols) > 0 {
		for _, v := range symbolValsSP.Config.MultiSymbols {
			sc, isok := pool.DefaultPaytables.MapSymbols[v]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:MultiSymbols",
					slog.String("symbol", v),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}
			symbolValsSP.Config.MultiSymbolCodes = append(symbolValsSP.Config.MultiSymbolCodes, sc)
		}

		symbolValsSP.Config.MultiType = parseSymbolValsSPMultiType(symbolValsSP.Config.StrMultiType)

		if symbolValsSP.Config.MultiTargetSymbol != "" {
			sc, isok := pool.DefaultPaytables.MapSymbols[symbolValsSP.Config.MultiTargetSymbol]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:MultiTargetSymbol",
					slog.String("symbol", symbolValsSP.Config.MultiTargetSymbol),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			symbolValsSP.Config.MultiTargetSymbolCode = sc
		}
	}

	if len(symbolValsSP.Config.CoinSymbols) > 0 {
		for _, v := range symbolValsSP.Config.CoinSymbols {
			sc, isok := pool.DefaultPaytables.MapSymbols[v]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:CoinSymbols",
					slog.String("symbol", v),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}
			symbolValsSP.Config.CoinSymbolCodes = append(symbolValsSP.Config.CoinSymbolCodes, sc)
		}
	} else {
		goutils.Error("SymbolValsSP.InitEx:CoinSymbols:no-symbols",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if len(symbolValsSP.Config.CollectSymbols) > 0 {
		for _, v := range symbolValsSP.Config.CollectSymbols {
			sc, isok := pool.DefaultPaytables.MapSymbols[v]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:CollectSymbols",
					slog.String("symbol", v),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}
			symbolValsSP.Config.CollectSymbolCodes = append(symbolValsSP.Config.CollectSymbolCodes, sc)
		}

		symbolValsSP.Config.CollectType = parseSymbolValsSPCollectType(symbolValsSP.Config.StrCollectType)

		if symbolValsSP.Config.CollectTargetSymbol != "" {
			sc, isok := pool.DefaultPaytables.MapSymbols[symbolValsSP.Config.CollectTargetSymbol]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:CollectTargetSymbol",
					slog.String("symbol", symbolValsSP.Config.CollectTargetSymbol),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}
			symbolValsSP.Config.CollectTargetSymbolCode = sc
		}

		if symbolValsSP.Config.CollectCoinSymbol != "" {
			sc, isok := pool.DefaultPaytables.MapSymbols[symbolValsSP.Config.CollectCoinSymbol]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:CollectCoinSymbol",
					slog.String("symbol", symbolValsSP.Config.CollectCoinSymbol),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}
			symbolValsSP.Config.CollectCoinSymbolCode = sc
		}

		if symbolValsSP.Config.CollectMultiSymbol != "" {
			sc, isok := pool.DefaultPaytables.MapSymbols[symbolValsSP.Config.CollectMultiSymbol]
			if !isok {
				goutils.Error("SymbolValsSP.InitEx:CollectMultiSymbol",
					slog.String("symbol", symbolValsSP.Config.CollectMultiSymbol),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}
			symbolValsSP.Config.CollectMultiSymbolCode = sc
		}
	}

	for _, awards := range symbolValsSP.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	symbolValsSP.onInit(&symbolValsSP.Config.BasicComponentConfig)

	return nil
}

// procMultiNormal applies the normal (non-round) multiplier logic.
// It finds MultiSymbols in the GameScene, sums corresponding values
// from OtherScene and multiplies coin symbols by the sum when
// multi > 1. It returns possibly cloned/modified versions of the
// game scene and other scene, plus a boolean indicating whether a
// multi trigger occurred.
func (symbolValsSP *SymbolValsSP) procMultiNormal(gameProp *GameProperty, _ sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {

	isTriggerMulti := false

	ngs := gs
	nos := os

	multi := 0

	mulpos := make([]int, 0, gs.Width*gs.Height*2)

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if slices.Contains(symbolValsSP.Config.MultiSymbolCodes, s) {
				if os.Arr[x][y] > 0 {
					multi += os.Arr[x][y]

					cd.MultiSymbolNum++

					mulpos = append(mulpos, x, y)
				}
			}
		}
	}

	cd.Multi = multi

	if multi > 1 {
		isTriggerMulti = true

		cd.newLine()

		for i := 0; i < len(mulpos); i += 2 {
			x := mulpos[i]
			y := mulpos[i+1]

			cd.AddPos(x, y)
		}

		if nos == os {
			nos = os.CloneEx(gameProp.PoolScene)
		}

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if slices.Contains(symbolValsSP.Config.CoinSymbolCodes, s) {
					if nos.Arr[x][y] > 0 {
						nos.Arr[x][y] = nos.Arr[x][y] * multi

						cd.MultiCoinSymbolNum++

						cd.AddPos(x, y)
					}
				}
			}
		}

		if symbolValsSP.Config.MultiTargetSymbolCode > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			for i := 0; i < len(mulpos); i += 2 {
				x := mulpos[i]
				y := mulpos[i+1]

				ngs.Arr[x][y] = symbolValsSP.Config.MultiTargetSymbolCode
			}
		}
	}

	return ngs, nos, isTriggerMulti, nil
}

// procMultiRoundXY applies round-type multiplier effect to the 3x3
// neighborhood centered at (cx,cy). It multiplies coins in OtherScene by
// multi and records affected positions. The function clones the other
// scene lazily (only when a modification is required) to avoid needless copies.
func (symbolValsSP *SymbolValsSP) procMultiRoundXY(gameProp *GameProperty, _ sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene, cx, cy int, multi int, isNewOS bool) (*sgc7game.GameScene, error) {

	cd.MultiSymbolNum++
	cd.newLine()
	cd.AddPos(cx, cy)

	nos := os

	for x := cx - 1; x <= cx+1; x++ {
		if x < 0 || x >= gs.Width {
			continue
		}

		for y := cy - 1; y <= cy+1; y++ {
			if y < 0 || y >= gs.Height {
				continue
			}

			if slices.Contains(symbolValsSP.Config.CoinSymbolCodes, gs.Arr[x][y]) {
				if nos.Arr[x][y] > 0 {
					if !isNewOS && nos == os {
						nos = os.CloneEx(gameProp.PoolScene)
					}

					nos.Arr[x][y] = nos.Arr[x][y] * multi

					cd.MultiCoinSymbolNum++

					cd.AddPos(x, y)
				}
			}
		}
	}

	return nos, nil
}

// procMultiRound handles the "round" multi type. For each multi symbol
// it delegates to procMultiRoundXY to apply the neighborhood multiplier
// and optionally sets target symbols in the returned game scene.
func (symbolValsSP *SymbolValsSP) procMultiRound(gameProp *GameProperty, plugin sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {

	ngs := gs
	nos := os

	cd.Multi = 1

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if slices.Contains(symbolValsSP.Config.MultiSymbolCodes, s) {
				if os.Arr[x][y] > 0 {
					cd.Multi *= os.Arr[x][y]

					// 这个函数里面还有一层检查是否需要clone os,如果clone,tos就是新的nos
					tos, err := symbolValsSP.procMultiRoundXY(gameProp, plugin, cd, gs, nos, x, y, os.Arr[x][y], os != nos)
					if err != nil {
						goutils.Error("SymbolValsSP.procMultiRound:procMultiRoundXY",
							slog.Int("x", x),
							slog.Int("y", y),
							goutils.Err(err))

						return nil, nil, false, err
					}

					nos = tos

					if symbolValsSP.Config.MultiTargetSymbolCode > 0 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						ngs.Arr[x][y] = symbolValsSP.Config.MultiTargetSymbolCode
					}
				}
			}
		}
	}

	return ngs, nos, os != nos, nil
}

// procMulti dispatches to the configured MultiType implementation
// and returns modified scenes, a trigger flag and possible error.
func (symbolValsSP *SymbolValsSP) procMulti(gameProp *GameProperty, plugin sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {

	switch symbolValsSP.Config.MultiType {
	case SVSPMultiTypeNormal:
		return symbolValsSP.procMultiNormal(gameProp, plugin, cd, gs, os)
	case SVSPMultiTypeRound:
		return symbolValsSP.procMultiRound(gameProp, plugin, cd, gs, os)
	}

	goutils.Error("SymbolValsSP.procMulti:InvalidMultiType",
		slog.String("type", symbolValsSP.Config.StrMultiType),
		goutils.Err(ErrInvalidComponentConfig))

	return nil, nil, false, ErrInvalidComponentConfig
}

// procCollectNormal implements the normal collect behavior: when any
// collect symbols appear, it sums coin values from OtherScene into each
// collect position and optionally replaces coin positions with a
// configured collect coin symbol in the returned game scene.
func (symbolValsSP *SymbolValsSP) procCollectNormal(gameProp *GameProperty, _ sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {

	isTriggerCollect := false

	ngs := gs
	nos := os

	totalcoin := 0
	coinpos := make([]int, 0, gs.Width*gs.Height*2)
	collectpos := make([]int, 0, gs.Width*gs.Height*2)
	mulpos := make([]int, 0, gs.Width*gs.Height*2)

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if slices.Contains(symbolValsSP.Config.CoinSymbolCodes, s) {
				if os.Arr[x][y] > 0 {
					totalcoin += os.Arr[x][y]

					cd.CollectCoinSymbolNum++

					coinpos = append(coinpos, x, y)
				}
			} else if slices.Contains(symbolValsSP.Config.CollectSymbolCodes, s) {
				cd.CollectSymbolNum++

				collectpos = append(collectpos, x, y)
			} else if slices.Contains(symbolValsSP.Config.MultiSymbolCodes, s) {
				mulpos = append(mulpos, x, y)
			}
		}
	}

	if len(collectpos) > 0 {
		isTriggerCollect = true
		cd.newLine()

		for i := 0; i < len(collectpos); i += 2 {
			x := collectpos[i]
			y := collectpos[i+1]

			if nos == os {
				nos = os.CloneEx(gameProp.PoolScene)
			}

			nos.Arr[x][y] = totalcoin

			cd.CollectCoin += totalcoin

			cd.AddPos(x, y)

			if symbolValsSP.Config.CollectTargetSymbolCode > 0 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = symbolValsSP.Config.CollectTargetSymbolCode
			}

		}

		if symbolValsSP.Config.CollectCoinSymbolCode > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			for i := 0; i < len(coinpos); i += 2 {
				x := coinpos[i]
				y := coinpos[i+1]

				cd.AddPos(x, y)

				ngs.Arr[x][y] = symbolValsSP.Config.CollectCoinSymbolCode
			}
		} else {
			for i := 0; i < len(coinpos); i += 2 {
				x := coinpos[i]
				y := coinpos[i+1]

				cd.AddPos(x, y)
			}
		}

		if symbolValsSP.Config.CollectMultiSymbolCode > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			for i := 0; i < len(mulpos); i += 2 {
				x := mulpos[i]
				y := mulpos[i+1]

				cd.AddPos(x, y)

				ngs.Arr[x][y] = symbolValsSP.Config.CollectMultiSymbolCode
			}
		}
	}

	return ngs, nos, isTriggerCollect, nil
}

// procCollectSequence implements collect behavior where only the first
// collect position receives the total coin amount (sequence mode).
func (symbolValsSP *SymbolValsSP) procCollectSequence(gameProp *GameProperty, _ sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {

	isTriggerCollect := false

	ngs := gs
	nos := os

	totalcoin := 0
	coinpos := make([]int, 0, gs.Width*gs.Height*2)
	collectpos := make([]int, 0, gs.Width*gs.Height*2)
	mulpos := make([]int, 0, gs.Width*gs.Height*2)

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if slices.Contains(symbolValsSP.Config.CoinSymbolCodes, s) {
				if os.Arr[x][y] > 0 {
					totalcoin += os.Arr[x][y]

					cd.CollectCoinSymbolNum++

					coinpos = append(coinpos, x, y)
				}
			} else if slices.Contains(symbolValsSP.Config.CollectSymbolCodes, s) {
				cd.CollectSymbolNum++

				collectpos = append(collectpos, x, y)
			} else if slices.Contains(symbolValsSP.Config.MultiSymbolCodes, s) {
				mulpos = append(mulpos, x, y)
			}
		}
	}

	if len(collectpos) > 0 {
		isTriggerCollect = true
		cd.newLine()

		x := collectpos[0]
		y := collectpos[1]

		if nos == os {
			nos = os.CloneEx(gameProp.PoolScene)
		}

		nos.Arr[x][y] = totalcoin

		cd.CollectCoin += totalcoin

		cd.AddPos(x, y)

		if symbolValsSP.Config.CollectTargetSymbolCode > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			ngs.Arr[x][y] = symbolValsSP.Config.CollectTargetSymbolCode
		}

		if symbolValsSP.Config.CollectCoinSymbolCode > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			for i := 0; i < len(coinpos); i += 2 {
				x := coinpos[i]
				y := coinpos[i+1]

				cd.AddPos(x, y)

				ngs.Arr[x][y] = symbolValsSP.Config.CollectCoinSymbolCode
			}
		} else {
			for i := 0; i < len(coinpos); i += 2 {
				x := coinpos[i]
				y := coinpos[i+1]

				cd.AddPos(x, y)
			}
		}

		if symbolValsSP.Config.CollectMultiSymbolCode > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			for i := 0; i < len(mulpos); i += 2 {
				x := mulpos[i]
				y := mulpos[i+1]

				cd.AddPos(x, y)

				ngs.Arr[x][y] = symbolValsSP.Config.CollectMultiSymbolCode
			}
		}
	}

	return ngs, nos, isTriggerCollect, nil
}

// procCollect dispatches to the configured CollectType implementation
// and returns modified scenes and a boolean indicating a collect trigger.
func (symbolValsSP *SymbolValsSP) procCollect(gameProp *GameProperty, _ sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {
	switch symbolValsSP.Config.CollectType {
	case SVSPCollectTypeNormal:
		return symbolValsSP.procCollectNormal(gameProp, nil, cd, gs, os)
	case SVSPCollectTypeSequence:
		return symbolValsSP.procCollectSequence(gameProp, nil, cd, gs, os)
	}

	goutils.Error("SymbolValsSP.procCollect:InvalidCollectType",
		slog.String("type", symbolValsSP.Config.StrCollectType),
		goutils.Err(ErrInvalidComponentConfig))

	return nil, nil, false, ErrInvalidComponentConfig
}

// procNormal runs the full processing flow for a play step: it applies
// multi processing then collect processing, adds any generated scenes
// or other-scenes into PlayResult via AddScene/AddOtherScene and
// returns flags indicating whether multi/collect were triggered.
func (symbolValsSP *SymbolValsSP) procNormal(gameProp *GameProperty, curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin, cd *SymbolValsSPData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (bool, bool, error) {

	isTriggerMulti := false
	isTriggerCollect := false

	multigs, multios, isTriggerMulti, err := symbolValsSP.procMulti(gameProp, plugin, cd, gs, os)
	if err != nil {
		goutils.Error("SymbolValsSP.procNormal:procMulti",
			goutils.Err(err))

		return false, false, err
	}

	if multigs != gs {
		symbolValsSP.AddScene(gameProp, curpr, multigs, &cd.BasicComponentData)
	}

	if multios != os {
		symbolValsSP.AddOtherScene(gameProp, curpr, multios, &cd.BasicComponentData)
	}

	colectgs, collectos, isTriggerCollect, err := symbolValsSP.procCollect(gameProp, plugin, cd, multigs, multios)
	if err != nil {
		goutils.Error("SymbolValsSP.procNormal:procCollect",
			goutils.Err(err))

		return false, false, err
	}

	if colectgs != multigs {
		symbolValsSP.AddScene(gameProp, curpr, colectgs, &cd.BasicComponentData)
	}

	if collectos != multios {
		symbolValsSP.AddOtherScene(gameProp, curpr, collectos, &cd.BasicComponentData)
	}

	return isTriggerMulti, isTriggerCollect, nil
}

// ProcControllers executes award controllers registered under the
// provided label (strVal). This allows configuration to trigger
// custom award logic on <trigger>/<multi>/<collect> events.
func (symbolValsSP *SymbolValsSP) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := symbolValsSP.Config.MapAwards[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

// playgame
// OnPlayGame is the main entry point invoked during a play. It expects
// icd to be a *SymbolValsSPData. It reads the current Scene and OtherScene
// from gameProp stacks, runs procNormal, triggers configured controllers
// and finalizes the step. If no multi or collect is triggered, it will
// return ErrComponentDoNothing.
func (symbolValsSP *SymbolValsSP) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd, ok := icd.(*SymbolValsSPData)
	if !ok || cd == nil {
		goutils.Error("SymbolValsSP.OnPlayGame:invalid icd",
			goutils.Err(ErrInvalidComponentConfig))
		return "", ErrInvalidComponentConfig
	}

	cd.onNewStep()

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	ogs := gameProp.OtherSceneStack.GetTopSceneEx(curpr, prs)

	switch symbolValsSP.Config.Type {
	case SVSPTypeNormal:
		isTriggerMulti, isTriggerCollect, err := symbolValsSP.procNormal(gameProp, curpr, plugin, cd, gs, ogs)
		if err != nil {
			goutils.Error("SymbolValsSP.OnPlayGame:procNormal",
				goutils.Err(err))

			return "", err
		}

		if !isTriggerMulti && !isTriggerCollect {
			nc := symbolValsSP.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		symbolValsSP.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")

		if isTriggerMulti {
			symbolValsSP.ProcControllers(gameProp, plugin, curpr, gp, 0, "<multi>")
		}

		if isTriggerCollect {
			symbolValsSP.ProcControllers(gameProp, plugin, curpr, gp, 0, "<collect>")
		}

		nc := symbolValsSP.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	goutils.Error("SymbolValsSP.OnPlayGame:InvalidType",
		slog.String("type", symbolValsSP.Config.StrType),
		goutils.Err(ErrInvalidComponentConfig))

	return "", ErrInvalidComponentConfig
}

// NewComponentData returns a fresh SymbolValsSPData instance for use
// in a new play. This satisfies the component factory contract.
func (symbolValsSP *SymbolValsSP) NewComponentData() IComponentData {
	return &SymbolValsSPData{}
}

// OnAsciiGame outputs the scene and other-scene modified by this
// component to the ASCII game output for debugging/visualization.
// It uses indices recorded in SymbolValsSPData. This method is best
// effort and returns nil when there is nothing to output.
func (symbolValsSP *SymbolValsSP) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd, ok := icd.(*SymbolValsSPData)
	if !ok || msd == nil {
		goutils.Error("SymbolValsSP.OnAsciiGame:invalid icd",
			goutils.Err(ErrInvalidComponentConfig))
		return ErrInvalidComponentConfig
	}

	if len(msd.UsedScenes) == 0 || len(msd.UsedOtherScenes) == 0 {
		// nothing to output
		return nil
	}

	// guard indexes
	si := msd.UsedScenes[0]
	oi := msd.UsedOtherScenes[0]
	if si < 0 || si >= len(pr.Scenes) || oi < 0 || oi >= len(pr.OtherScenes) {
		return nil
	}

	asciigame.OutputScene("after SymbolValsSP", pr.Scenes[si], mapSymbolColor)
	asciigame.OutputOtherScene("after SymbolValsSP", pr.OtherScenes[oi])

	return nil
}

// NewSymbolValsSP constructs a SymbolValsSP component with the provided name.
func NewSymbolValsSP(name string) IComponent {
	return &SymbolValsSP{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "normal",
// "coinSymbols": [
//     "CA"
// ],
// "mulSymbols": [
//     "MU"
// ],
// "mulType": "round",
// "collectSymbols": [
//     "CO"
// ],
// "collectTargetSymbol": "CA",
// "collectCoinSymbol": "BN",
// "collectMulSymbol": "BN"
// "collectType": "sequence"

// jsonSymbolValsSP is a lightweight structure matching the JSON
// representation used by the parsing helpers. The build() method
// converts it into the canonical SymbolValsSPConfig.
type jsonSymbolValsSP struct {
	StrType             string   `json:"type"`
	CoinSymbols         []string `json:"coinSymbols"`
	MulSymbols          []string `json:"mulSymbols"`
	StrMulType          string   `json:"mulType"`
	MulTargetSymbol     string   `json:"mulTargetSymbol"`
	CollectSymbols      []string `json:"collectSymbols"`
	CollectTargetSymbol string   `json:"collectTargetSymbol"`
	CollectCoinSymbol   string   `json:"collectCoinSymbol"`
	CollectMulSymbol    string   `json:"collectMulSymbol"`
	CollectType         string   `json:"collectType"`
}

// build converts the parsed JSON struct into a SymbolValsSPConfig with
// normalized string fields. Symbol names are not resolved here; that
// happens in InitEx using the paytable map.
func (jcfg *jsonSymbolValsSP) build() *SymbolValsSPConfig {
	cfg := &SymbolValsSPConfig{
		StrType:             strings.ToLower(jcfg.StrType),
		CoinSymbols:         slices.Clone(jcfg.CoinSymbols),
		MultiSymbols:        slices.Clone(jcfg.MulSymbols),
		StrMultiType:        strings.ToLower(jcfg.StrMulType),
		MultiTargetSymbol:   jcfg.MulTargetSymbol,
		CollectSymbols:      slices.Clone(jcfg.CollectSymbols),
		CollectTargetSymbol: jcfg.CollectTargetSymbol,
		CollectCoinSymbol:   jcfg.CollectCoinSymbol,
		CollectMultiSymbol:  jcfg.CollectMulSymbol,
		StrCollectType:      strings.ToLower(jcfg.CollectType),
	}

	return cfg
}

// parseSymbolValsSP parses a configuration cell (AST node) for the
// SymbolValsSP component and registers resulting configs into gamecfg.
// It returns the component label on success.
func parseSymbolValsSP(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolValsSP:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolValsSP:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolValsSP{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolValsSP:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseSymbolValsSP:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolValsSPTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
