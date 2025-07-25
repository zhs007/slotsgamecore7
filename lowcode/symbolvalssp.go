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

const SymbolValsSPTypeName = "symbolValsSP"

type SymbolValsSPType int

const (
	SVSPTypeNormal SymbolValsSPType = 0 // normal
)

func parseSymbolValsSPType(str string) SymbolValsSPType {
	return SVSPTypeNormal
}

type SymbolValsSPMultiType int

const (
	SVSPMultiTypeNormal SymbolValsSPMultiType = 0 // normal
	SVSPMultiTypeRound  SymbolValsSPMultiType = 1 // round
)

func parseSymbolValsSPMultiType(str string) SymbolValsSPMultiType {
	if str == "round" {
		return SVSPMultiTypeRound
	}

	return SVSPMultiTypeNormal
}

type SymbolValsSPCollectType int

const (
	SVSPCollectTypeNormal   SymbolValsSPCollectType = 0 // normal
	SVSPCollectTypeSequence SymbolValsSPCollectType = 1 // sequence
)

func parseSymbolValsSPCollectType(str string) SymbolValsSPCollectType {
	if str == "sequence" {
		return SVSPCollectTypeSequence
	}

	return SVSPCollectTypeNormal
}

type SymbolValsSPData struct {
	BasicComponentData
	Pos                  []int // 位置
	UsedScenes           []int // 使用的场景, -1分隔，分隔后，第一个位置是触发图标，后面是影响到的图标
	MultiSymbolNum       int   // 倍数图标个数
	MultiCoinSymbolNum   int   // 倍数影响的图标个数
	Multi                int   // 倍数
	CollectSymbolNum     int   // 收集图标个数
	CollectCoinSymbolNum int   // 收集到 coin 图标个数
	CollectCoin          int   // 收集到的 coin 金额
}

// OnNewGame -
func (symbolValsSPData *SymbolValsSPData) OnNewGame(gameProp *GameProperty, component IComponent) {
	symbolValsSPData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (symbolValsSPData *SymbolValsSPData) onNewStep() {
	symbolValsSPData.UsedScenes = nil
	symbolValsSPData.UsedOtherScenes = nil

	symbolValsSPData.Pos = nil
	symbolValsSPData.MultiSymbolNum = 0
	symbolValsSPData.MultiCoinSymbolNum = 0
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
type SymbolValsSPConfig struct {
	BasicComponentConfig    `yaml:",inline" json:",inline"`
	StrType                 string                  `yaml:"type" json:"type"`
	Type                    SymbolValsSPType        `yaml:"-" json:"-"`
	CoinSymbols             []string                `yaml:"coinSymbols" json:"coinSymbols"`
	CoinSymbolCodes         []int                   `yaml:"-" json:"-"`
	MultiSymbols            []string                `yaml:"multiSymbols" json:"multiSymbols"`
	MultiSymbolCodes        []int                   `yaml:"-" json:"-"`
	StrMultiType            string                  `yaml:"multiType" json:"multiType"`
	MultiType               SymbolValsSPMultiType   `yaml:"-" json:"-"`
	MultiTargetSymbol       string                  `yaml:"multiTargetSymbol" json:"multiTargetSymbol"`
	MultiTargetSymbolCode   int                     `yaml:"-" json:"-"`
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
	MapAwards               map[string][]*Award     `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *SymbolValsSPConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	}
}

type SymbolValsSP struct {
	*BasicComponent `json:"-"`
	Config          *SymbolValsSPConfig `json:"config"`
}

// Init -
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
func (symbolValsSP *SymbolValsSP) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolValsSP.Config = cfg.(*SymbolValsSPConfig)
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

// procMultiNormal -
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

// procMultiRoundXY -
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

// procMultiRound -
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

// procMulti -
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

// procCollectNormal -
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

// procCollectSequence -
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

// procCollect -
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

// procNormal -
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

// OnProcControllers -
func (symbolValsSP *SymbolValsSP) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := symbolValsSP.Config.MapAwards[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

// playgame
func (symbolValsSP *SymbolValsSP) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*SymbolValsSPData)

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

// NewComponentData -
func (symbolValsSP *SymbolValsSP) NewComponentData() IComponentData {
	return &SymbolValsSPData{}
}

// OnAsciiGame - outpur to asciigame
func (symbolValsSP *SymbolValsSP) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd := icd.(*SymbolValsSPData)

	asciigame.OutputScene("after SymbolValsSP", pr.Scenes[msd.UsedScenes[0]], mapSymbolColor)
	asciigame.OutputOtherScene("after SymbolValsSP", pr.OtherScenes[msd.UsedOtherScenes[0]])

	return nil
}

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
