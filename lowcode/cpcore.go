package lowcode

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"

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

const CPCoreTypeName = "CPCore"

type mainSymbolInfo struct {
	symbolCode int
	level      int
	price      int
	moved      bool
	syms       []int
}

type CPCoreData struct {
	BasicComponentData
	lstMainSymbols       []*mainSymbolInfo
	isPopcornTriggered   bool
	gridSize             int
	isDontPressTriggered bool
	isEggTriggered       bool
	cfg                  *CPCoreConfig
	spSymVW              *sgc7game.ValWeights2
	isSpSymEgg           bool
	isSpSymPopcorn       bool
	isSpSymDontPress     bool
	spSymBonusNum        int
}

func (gcd *CPCoreData) clearMove() {
	for _, ms := range gcd.lstMainSymbols {
		ms.moved = false
	}
}

func (gcd *CPCoreData) onLevelUp(mainSymbol int, off int) {
	for _, ms := range gcd.lstMainSymbols {
		if ms.symbolCode == mainSymbol {
			ms.level += off

			if ms.level >= len(gcd.cfg.MapSymbolCode[ms.symbolCode]) {
				ms.level = len(gcd.cfg.MapSymbolCode[ms.symbolCode]) - 1
			}

			ms.price = gcd.cfg.MapSymbolCode[ms.symbolCode][ms.level]
		}
	}
}

func (gcd *CPCoreData) onAllLevelUp(off int) {
	for _, ms := range gcd.lstMainSymbols {
		ms.level += off

		if ms.level >= len(gcd.cfg.MapSymbolCode[ms.symbolCode]) {
			ms.level = len(gcd.cfg.MapSymbolCode[ms.symbolCode]) - 1
		}

		ms.price = gcd.cfg.MapSymbolCode[ms.symbolCode][ms.level]
	}
}

func (gcd *CPCoreData) getSymbolCode(ms int) int {

	for _, msi := range gcd.lstMainSymbols {
		if msi.symbolCode == ms {
			return gcd.cfg.MapSymbolCode[ms][msi.level]
		}
	}

	return -1
}

func (gcd *CPCoreData) getMainSymbolInfo(ms int) *mainSymbolInfo {

	for _, msi := range gcd.lstMainSymbols {
		if msi.symbolCode == ms {
			return msi
		}
	}

	return nil
}

func (gcd *CPCoreData) getNext() int {
	msc := -1
	msp := -1

	for _, ms := range gcd.lstMainSymbols {
		if ms.moved {
			continue
		}

		if ms.price > msp {
			msp = ms.price
			msc = ms.symbolCode
		}
	}

	return msc
}

func (gcd *CPCoreData) moveEnd(ms int) {
	for _, msi := range gcd.lstMainSymbols {
		if msi.symbolCode == ms {
			msi.moved = true

			return
		}
	}
}

func (gcd *CPCoreData) procSymbolsWithLevel(gs *sgc7game.GameScene) {
	for x, arr := range gs.Arr {
		for y, sc := range arr {
			for _, ms := range gcd.lstMainSymbols {
				if slices.Contains(ms.syms, sc) {
					gs.Arr[x][y] = ms.syms[ms.level]
				}
			}
		}
	}
}

func (gcd *CPCoreData) isCoin(sym int) bool {
	return slices.Contains(gcd.cfg.CoinSymbolCodes, sym)
}

func (gcd *CPCoreData) randCoin(plugin sgc7plugin.IPlugin) (int, error) {
	num, err := gcd.cfg.CoinWeightVW.RandVal(plugin)
	if err != nil {
		goutils.Error("CPCoreData.randCoin:RandVal",
			goutils.Err(err))

		return 0, err
	}

	return num.Int(), nil
}

func (gcd *CPCoreData) randSpNum(plugin sgc7plugin.IPlugin) (int, error) {
	num, err := gcd.cfg.SpSymbolNumWeightVW.RandVal(plugin)
	if err != nil {
		goutils.Error("CPCoreData.randSpNum:RandVal",
			goutils.Err(err))

		return 0, err
	}

	return num.Int(), nil
}

func (gcd *CPCoreData) randSpSym(plugin sgc7plugin.IPlugin) (int, error) {
	num, err := gcd.spSymVW.RandVal(plugin)
	if err != nil {
		goutils.Error("CPCoreData.randSpSym:RandVal",
			goutils.Err(err))

		return 0, err
	}

	return num.Int(), nil
}

func (gcd *CPCoreData) onInitGenSpSym(sym int) {
	switch sym {
	case gcd.cfg.EggSymbolCode:
		gcd.isSpSymEgg = true

		gcd.spSymVW = gcd.spSymVW.CloneWithoutIntArray([]int{gcd.cfg.EggSymbolCode})
	case gcd.cfg.PopcornSymbolCode:
		gcd.isSpSymPopcorn = true

		gcd.spSymVW = gcd.spSymVW.CloneWithoutIntArray([]int{gcd.cfg.PopcornSymbolCode})
	case gcd.cfg.DontPressSymbolCode:
		gcd.isSpSymDontPress = true

		gcd.spSymVW = gcd.spSymVW.CloneWithoutIntArray([]int{gcd.cfg.DontPressSymbolCode})
	case gcd.cfg.BonusSymbolCode:
		gcd.spSymBonusNum++

		if gcd.spSymBonusNum >= 3 {
			gcd.spSymVW = gcd.spSymVW.CloneWithoutIntArray([]int{gcd.cfg.BonusSymbolCode})
		}
	}

}

// OnNewGame -
func (gcd *CPCoreData) OnNewGame(gameProp *GameProperty, component IComponent) {
	gcd.BasicComponentData.OnNewGame(gameProp, component)

	gcd.isPopcornTriggered = false
	gcd.isDontPressTriggered = false
	gcd.isEggTriggered = false
	gcd.gridSize = 6

	gcd.isSpSymDontPress = false
	gcd.isSpSymEgg = false
	gcd.isSpSymPopcorn = false
	gcd.spSymBonusNum = 0

	gcd.spSymVW = gcd.cfg.SpSymbolWeightVW

	for _, ms := range gcd.lstMainSymbols {
		ms.level = 0
		sc := ms.syms[0]
		ms.price = gameProp.Pool.DefaultPaytables.MapPay[sc][0]
	}
}

// OnNewStep -
func (gcd *CPCoreData) onNewStep() {
	gcd.UsedScenes = nil
}

// Clone
func (gcd *CPCoreData) Clone() IComponentData {
	target := &CPCoreData{
		BasicComponentData:   gcd.CloneBasicComponentData(),
		isPopcornTriggered:   gcd.isPopcornTriggered,
		gridSize:             gcd.gridSize,
		isEggTriggered:       gcd.isEggTriggered,
		isDontPressTriggered: gcd.isDontPressTriggered,
		cfg:                  gcd.cfg,
		spSymVW:              gcd.spSymVW,
		isSpSymEgg:           gcd.isSpSymEgg,
		isSpSymPopcorn:       gcd.isSpSymPopcorn,
		isSpSymDontPress:     gcd.isSpSymDontPress,
		spSymBonusNum:        gcd.spSymBonusNum,
		lstMainSymbols:       slices.Clone(gcd.lstMainSymbols),
	}

	return target
}

// BuildPBComponentData
func (gcd *CPCoreData) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: gcd.BuildPBBasicComponentData(),
	}
}

// CPCoreConfig - placeholder configuration for CPCore
type CPCoreConfig struct {
	BasicComponentConfig    `yaml:",inline" json:",inline"`
	CategoryCount           int                   `yaml:"categoryCount" json:"categoryCount"`
	MapSymbol               map[string][]string   `yaml:"mapSymbol" json:"mapSymbol"`
	MapSymbolCode           map[int][]int         `yaml:"-" json:"-"`
	BlankSymbol             string                `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode         int                   `yaml:"-" json:"-"`
	WildSymbol              string                `yaml:"wildSymbol" json:"wildSymbol"`
	WildSymbolCode          int                   `yaml:"-" json:"-"`
	WildUsedSymbolCode      int                   `yaml:"-" json:"-"`
	CoinSymbols             []string              `yaml:"coinSymbols" json:"coinSymbols"`
	CoinSymbolCodes         []int                 `yaml:"-" json:"-"`
	UpLevelSymbols          []string              `yaml:"upLevelSymbols" json:"upLevelSymbols"`
	UpLevelSymbolCodes      []int                 `yaml:"-" json:"-"`
	AllUpLevelSymbols       []string              `yaml:"allUpLevelSymbol" json:"allUpLevelSymbol"`
	AllUpLevelSymbolCodes   []int                 `yaml:"-" json:"-"`
	SwitcherSymbol          string                `yaml:"switcherSymbol" json:"switcherSymbol"`
	SwitcherSymbolCode      int                   `yaml:"-" json:"-"`
	PopcornSymbol           string                `yaml:"popcornSymbol" json:"popcornSymbol"`
	PopcornSymbolCode       int                   `yaml:"-" json:"-"`
	EggSymbol               string                `yaml:"eggSymbol" json:"eggSymbol"`
	EggSymbolCode           int                   `yaml:"-" json:"-"`
	EggUsedSymbolCode       int                   `yaml:"-" json:"-"`
	DontPressSymbol         string                `yaml:"dontpressSymbol" json:"dontpressSymbol"`
	DontPressSymbolCode     int                   `yaml:"-" json:"-"`
	DontPressUsedSymbolCode int                   `yaml:"-" json:"-"`
	mapSymbolValues         map[int]int           `yaml:"-" json:"-"`
	lstMainSymbols          []int                 `yaml:"-" json:"-"`
	CollectorMaxLevel       int                   `yaml:"collectorMaxLevel" json:"collectorMaxLevel"`
	CollectorNumber         int                   `yaml:"collectorNumber" json:"collectorNumber"`
	CoinWeight              string                `yaml:"coinWeight" json:"coinWeight"`
	CoinWeightVW            *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	CollectorReleasedNum    int                   `yaml:"collectorReleasedNum" json:"collectorReleasedNum"`
	SpSymbolWeight          string                `yaml:"spSymbolWeight" json:"spSymbolWeight"`
	SpSymbolWeightVW        *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	BonusSymbol             string                `yaml:"bonusSymbol" json:"bonusSymbol"`
	BonusSymbolCode         int                   `yaml:"-" json:"-"`
	SpSymbolNumWeight       string                `yaml:"spSymbolNumWeight" json:"spSymbolNumWeight"`
	SpSymbolNumWeightVW     *sgc7game.ValWeights2 `yaml:"-" json:"-"`

	MapControllers map[string][]*Award `yaml:"mapControllers" json:"mapControllers"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *CPCoreConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type CPCore struct {
	*BasicComponent `json:"-"`
	Config          *CPCoreConfig `json:"config"`
}

// Init - load from file (placeholder)
func (cpc *CPCore) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CPCore.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CPCoreConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CPCore.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return cpc.InitEx(cfg, pool)
}

// InitEx - initialize from config object (placeholder)
func (cpc *CPCore) InitEx(cfg any, pool *GamePropertyPool) error {
	cpc.Config = cfg.(*CPCoreConfig)
	cpc.Config.ComponentType = CPCoreTypeName

	if cpc.Config.BlankSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpc.Config.BlankSymbol]
		if !isok {
			goutils.Error("CPCore.InitEx:BlankSymbol",
				slog.String("BlankSymbol", cpc.Config.BlankSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.BlankSymbolCode = sc
	} else {
		goutils.Error("CPCore.InitEx:BlankSymbol",
			slog.String("BlankSymbol", cpc.Config.BlankSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.WildSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpc.Config.WildSymbol]
		if !isok {
			goutils.Error("CPCore.InitEx:WildSymbol",
				slog.String("WildSymbol", cpc.Config.WildSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.WildSymbolCode = sc

		sc, isok = pool.Config.GetDefaultPaytables().MapSymbols[fmt.Sprintf("%vD", cpc.Config.WildSymbol)]
		if !isok {
			goutils.Error("CPCore.InitEx:WildSymbol",
				slog.String("WildSymbol", cpc.Config.WildSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.WildUsedSymbolCode = sc
	} else {
		goutils.Error("CPCore.InitEx:WildSymbol",
			slog.String("WildSymbol", cpc.Config.WildSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.SwitcherSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpc.Config.SwitcherSymbol]
		if !isok {
			goutils.Error("CPCore.InitEx:SwitcherSymbol",
				slog.String("SwitcherSymbol", cpc.Config.SwitcherSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.SwitcherSymbolCode = sc
	} else {
		goutils.Error("CPCore.InitEx:SwitcherSymbol",
			slog.String("SwitcherSymbol", cpc.Config.SwitcherSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.PopcornSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpc.Config.PopcornSymbol]
		if !isok {
			goutils.Error("CPCore.InitEx:PopcornSymbol",
				slog.String("PopcornSymbol", cpc.Config.PopcornSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.PopcornSymbolCode = sc
	} else {
		goutils.Error("CPCore.InitEx:PopcornSymbol",
			slog.String("PopcornSymbol", cpc.Config.PopcornSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.EggSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpc.Config.EggSymbol]
		if !isok {
			goutils.Error("CPCore.InitEx:EggSymbol",
				slog.String("EggSymbol", cpc.Config.EggSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.EggSymbolCode = sc

		sc, isok = pool.Config.GetDefaultPaytables().MapSymbols[fmt.Sprintf("%vD", cpc.Config.EggSymbol)]
		if !isok {
			goutils.Error("CPCore.InitEx:EggSymbol",
				slog.String("EggSymbol", cpc.Config.EggSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.EggUsedSymbolCode = sc
	} else {
		goutils.Error("CPCore.InitEx:EggSymbol",
			slog.String("EggSymbol", cpc.Config.EggSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.DontPressSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpc.Config.DontPressSymbol]
		if !isok {
			goutils.Error("CPCore.InitEx:DontPressSymbol",
				slog.String("DontPressSymbol", cpc.Config.DontPressSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.DontPressSymbolCode = sc

		sc, isok = pool.Config.GetDefaultPaytables().MapSymbols[fmt.Sprintf("%vD", cpc.Config.DontPressSymbol)]
		if !isok {
			goutils.Error("CPCore.InitEx:DontPressSymbol",
				slog.String("DontPressSymbol", cpc.Config.DontPressSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpc.Config.DontPressUsedSymbolCode = sc
	} else {
		goutils.Error("CPCore.InitEx:DontPressSymbol",
			slog.String("DontPressSymbol", cpc.Config.DontPressSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if len(cpc.Config.CoinSymbols) > 0 {
		cpc.Config.CoinSymbolCodes = make([]int, len(cpc.Config.CoinSymbols))
		for i, cs := range cpc.Config.CoinSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CPCore.InitEx:CoinSymbols",
					slog.String("CoinSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpc.Config.CoinSymbolCodes[i] = sc
		}
	} else {
		goutils.Error("CPCore.InitEx:CoinSymbols",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if len(cpc.Config.UpLevelSymbols) > 0 {
		cpc.Config.UpLevelSymbolCodes = make([]int, len(cpc.Config.UpLevelSymbols))
		for i, cs := range cpc.Config.UpLevelSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CPCore.InitEx:UpLevelSymbols",
					slog.String("UpLevelSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpc.Config.UpLevelSymbolCodes[i] = sc
		}
	} else {
		goutils.Error("CPCore.InitEx:UpLevelSymbols",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if len(cpc.Config.AllUpLevelSymbols) > 0 {
		cpc.Config.AllUpLevelSymbolCodes = make([]int, len(cpc.Config.AllUpLevelSymbols))
		for i, cs := range cpc.Config.AllUpLevelSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CPCore.InitEx:AllUpLevelSymbols",
					slog.String("AllUpLevelSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpc.Config.AllUpLevelSymbolCodes[i] = sc
		}
	} else {
		goutils.Error("CPCore.InitEx:AllUpLevelSymbols",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if len(cpc.Config.MapSymbol) > 0 {
		cpc.Config.MapSymbolCode = make(map[int][]int)
		for ms, css := range cpc.Config.MapSymbol {
			mssc, isok := pool.Config.GetDefaultPaytables().MapSymbols[ms]
			if !isok {
				goutils.Error("CPCore.InitEx:MapSymbol",
					slog.String("MainSymbol", ms),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cssc := make([]int, len(css))

			for i, cs := range css {
				sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
				if !isok {
					goutils.Error("CPCore.InitEx:MapSymbol:CollectedSymbols",
						slog.String("CollectedSymbol", cs),
						goutils.Err(ErrInvalidComponentConfig))

					return ErrInvalidComponentConfig
				}

				cssc[i] = sc
			}

			cpc.Config.MapSymbolCode[mssc] = cssc

			cpc.Config.lstMainSymbols = append(cpc.Config.lstMainSymbols, mssc)
		}

		sort.Slice(cpc.Config.lstMainSymbols, func(i, j int) bool { return cpc.Config.lstMainSymbols[i] < cpc.Config.lstMainSymbols[j] })
	}

	if cpc.Config.CoinWeight != "" {
		vw, err := pool.LoadIntWeights(cpc.Config.CoinWeight, true)
		if err != nil {
			goutils.Error("CPCore.InitEx:CoinWeight:GetValWeights2",
				slog.String("CoinWeight", cpc.Config.CoinWeight),
				goutils.Err(err))

			return err
		}

		cpc.Config.CoinWeightVW = vw
	} else {
		goutils.Error("CPCore.InitEx:CoinWeight",
			slog.String("CoinWeight", cpc.Config.CoinWeight),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.SpSymbolWeight != "" {
		vw, err := pool.LoadIntWeights(cpc.Config.SpSymbolWeight, true)
		if err != nil {
			goutils.Error("CPCore.InitEx:SpSymbolWeight:GetValWeights2",
				slog.String("SpSymbolWeight", cpc.Config.SpSymbolWeight),
				goutils.Err(err))

			return err
		}

		cpc.Config.SpSymbolWeightVW = vw
	} else {
		goutils.Error("CPCore.InitEx:SpSymbolWeight",
			slog.String("SpSymbolWeight", cpc.Config.SpSymbolWeight),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if cpc.Config.SpSymbolNumWeight != "" {
		vw, err := pool.LoadIntWeights(cpc.Config.SpSymbolNumWeight, true)
		if err != nil {
			goutils.Error("CPCore.InitEx:SpSymbolNumWeight:GetValWeights2",
				slog.String("SpSymbolNumWeight", cpc.Config.SpSymbolNumWeight),
				goutils.Err(err))

			return err
		}

		cpc.Config.SpSymbolNumWeightVW = vw
	} else {
		goutils.Error("CPCore.InitEx:SpSymbolNumWeight",
			slog.String("SpSymbolNumWeight", cpc.Config.SpSymbolNumWeight),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	// coin
	// alllevelup
	// levelup
	// switcher
	// popcorn
	// wild
	// normal symbols
	// egg
	// dontpress
	// mainSymbol

	cpc.Config.mapSymbolValues = make(map[int]int)

	cpc.Config.mapSymbolValues[cpc.Config.DontPressSymbolCode] = 2
	cpc.Config.mapSymbolValues[cpc.Config.EggSymbolCode] = 3

	for k, arr := range cpc.Config.MapSymbolCode {
		cpc.Config.mapSymbolValues[k] = 1

		for _, a := range arr {
			cpc.Config.mapSymbolValues[a] = 4
		}
	}

	cpc.Config.mapSymbolValues[cpc.Config.WildSymbolCode] = 5
	cpc.Config.mapSymbolValues[cpc.Config.PopcornSymbolCode] = 6
	cpc.Config.mapSymbolValues[cpc.Config.SwitcherSymbolCode] = 7
	for _, v := range cpc.Config.UpLevelSymbolCodes {
		cpc.Config.mapSymbolValues[v] = 8
	}
	for _, v := range cpc.Config.AllUpLevelSymbolCodes {
		cpc.Config.mapSymbolValues[v] = 9
	}
	for _, v := range cpc.Config.CoinSymbolCodes {
		cpc.Config.mapSymbolValues[v] = 10
	}

	for _, arr := range cpc.Config.MapControllers {
		for _, a := range arr {
			if a != nil {
				a.Init()
			}
		}
	}

	cpc.onInit(&cpc.Config.BasicComponentConfig)

	return nil
}

// OnPlayGame - placeholder
func (cpc *CPCore) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd, isok := icd.(*CPCoreData)
	if !isok {
		goutils.Error("CPCore.OnPlayGame:invalid component data",
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	cd.onNewStep()

	posd := gameProp.posPool.Get()

	gs := cpc.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("CPCore.OnPlayGame:GetTargetScene3",
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	os := cpc.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	if os == nil {
		os = gameProp.PoolScene.New(gs.Width, gs.Height)
	} else {
		os = os.CloneEx(gameProp.PoolScene)
	}

	for x, arr := range gs.Arr {
		for y, v := range arr {
			if v >= 0 {
				posd.Add(x, y)
			}
		}
	}

	for k := range cpc.Config.MapSymbolCode {
		if posd.IsEmpty() {
			goutils.Error("CPCore.OnPlayGame:Random",
				goutils.Err(ErrInvalidComponentData))

			return "", ErrInvalidComponentData
		}

		ci, err := plugin.Random(context.Background(), posd.Len())
		if err != nil {
			goutils.Error("CPCore.OnPlayGame:Random",
				goutils.Err(err))

			return "", err
		}

		x, y := posd.Get(ci)
		ngs.Arr[x][y] = k
		posd.Del(ci)
	}

	spnum, err := cd.randSpNum(plugin)
	if err != nil {
		goutils.Error("CPCore.OnPlayGame:randSpNum",
			goutils.Err(err))

		return "", err
	}

	for i := 0; i < spnum; i++ {
		if posd.IsEmpty() {
			break
		}

		ci, err := plugin.Random(context.Background(), posd.Len())
		if err != nil {
			goutils.Error("CPCore.OnPlayGame:Random2",
				goutils.Err(err))

			return "", err
		}

		x, y := posd.Get(ci)

		cs, err := cd.randSpSym(plugin)
		if err != nil {
			goutils.Error("CPCore.OnPlayGame:randSpSym",
				goutils.Err(err))

			return "", err
		}

		ngs.Arr[x][y] = cs

		posd.Del(ci)

		cd.onInitGenSpSym(cs)

		if cd.isCoin(cs) {
			coin, err := cd.randCoin(plugin)
			if err != nil {
				goutils.Error("CPCore.OnPlayGame:randCoin",
					goutils.Err(err))

				return "", err
			}

			os.Arr[x][y] = coin
		}
	}

	cpc.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	cpc.AddOtherScene(gameProp, curpr, os, &cd.BasicComponentData)

	nc := cpc.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - placeholder
func (cpc *CPCore) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData -
func (cpc *CPCore) NewComponentData() IComponentData {
	cd := &CPCoreData{
		cfg: cpc.Config,
	}

	for _, ms := range cpc.Config.lstMainSymbols {
		cd.lstMainSymbols = append(cd.lstMainSymbols, &mainSymbolInfo{
			symbolCode: ms,
			level:      0,
			price:      0,
			syms:       cpc.Config.MapSymbolCode[ms],
		})
	}

	return cd
}

func NewCPCore(name string) IComponent {
	return &CPCore{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "categoryCount": 4,
// "type": "normal",
// "collectorMaxLevel": 3,
// "collectorNumber": 20,
// "coinWeight": "coinweight",
// "mapSymbol": [
//     {
//         "mainSymbol": "RP",
//         "collectedSymbols": [
//             "R1",
//             "R2",
//             "R3",
//             "R4",
//             "R5",
//             "R6",
//             "R7"
//         ]
//     },
//     {
//         "mainSymbol": "PP",
//         "collectedSymbols": [
//             "P1",
//             "P2",
//             "P3",
//             "P4",
//             "P5",
//             "P6",
//             "P7"
//         ]
//     },
//     {
//         "mainSymbol": "GP",
//         "collectedSymbols": [
//             "G1",
//             "G2",
//             "G3",
//             "G4",
//             "G5",
//             "G6",
//             "G7"
//         ]
//     },
//     {
//         "mainSymbol": "BP",
//         "collectedSymbols": [
//             "B1",
//             "B2",
//             "B3",
//             "B4",
//             "B5",
//             "B6",
//             "B7"
//         ]
//     }
// ],
// "wildSymbol": "WL",
// "coinSymbols": [
//     "CN"
// ],
// "upLevelSymbols": [
//     "L1",
//     "L2",
//     "L3"
// ],
// "allUpLevelSymbols": [
//     "AL1",
//     "AL2",
//     "AL3"
// ],
// "switcherSymbol": "MR",
// "popcornSymbol": "PC",
// "eggSymbol": "EG",
// "dontpressSymbol": "DP",
// "collectorReleasedNum": 3,
// "spSymbolWeight": "bgspsymweight",
// "blankSymbol": "BN",
// "bonusSymbol": "FG",
// "spSymbolNumWeight": "bgspnumweight"

type jsonCPCoreSymbolData struct {
	MainSymbol       string   `json:"mainSymbol"`
	CollectedSymbols []string `json:"collectedSymbols"`
}

type jsonCPCore struct {
	CategoryCount        int                    `json:"categoryCount"`
	CollectorMaxLevel    int                    `json:"collectorMaxLevel"`
	CollectorNumber      int                    `json:"collectorNumber"`
	CoinWeight           string                 `json:"coinWeight"`
	MapSymbol            []jsonCPCoreSymbolData `json:"mapSymbol"`
	BlankSymbol          string                 `json:"blankSymbol"`
	WildSymbol           string                 `json:"wildSymbol"`
	CoinSymbols          []string               `json:"coinSymbols"`
	UpLevelSymbols       []string               `json:"upLevelSymbols"`
	AllUpLevelSymbols    []string               `json:"allUpLevelSymbols"`
	SwitcherSymbol       string                 `json:"switcherSymbol"`
	PopcornSymbol        string                 `json:"popcornSymbol"`
	EggSymbol            string                 `json:"eggSymbol"`
	DontPressSymbol      string                 `json:"dontpressSymbol"`
	CollectorReleasedNum int                    `json:"collectorReleasedNum"`
	SpSymbolWeight       string                 `json:"spSymbolWeight"`
	BonusSymbol          string                 `json:"bonusSymbol"`
	SpSymbolNumWeight    string                 `json:"spSymbolNumWeight"`
}

func (j *jsonCPCore) build() *CPCoreConfig {
	cfg := &CPCoreConfig{
		CategoryCount:        j.CategoryCount,
		CollectorMaxLevel:    j.CollectorMaxLevel,
		CollectorNumber:      j.CollectorNumber,
		CoinWeight:           j.CoinWeight,
		BlankSymbol:          j.BlankSymbol,
		WildSymbol:           j.WildSymbol,
		MapSymbol:            make(map[string][]string),
		CoinSymbols:          slices.Clone(j.CoinSymbols),
		UpLevelSymbols:       slices.Clone(j.UpLevelSymbols),
		AllUpLevelSymbols:    slices.Clone(j.AllUpLevelSymbols),
		SwitcherSymbol:       j.SwitcherSymbol,
		PopcornSymbol:        j.PopcornSymbol,
		EggSymbol:            j.EggSymbol,
		DontPressSymbol:      j.DontPressSymbol,
		CollectorReleasedNum: j.CollectorReleasedNum,
		SpSymbolWeight:       j.SpSymbolWeight,
		BonusSymbol:          j.BonusSymbol,
		SpSymbolNumWeight:    j.SpSymbolNumWeight,
	}

	for _, ms := range j.MapSymbol {
		cfg.MapSymbol[ms.MainSymbol] = ms.CollectedSymbols
	}

	return cfg
}

// parseCPCore - minimal JSON cell parser for CPCore (placeholder)
func parseCPCore(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCPCore:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCPCore:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCPCore{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCPCore:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseCPCore:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CPCoreTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
