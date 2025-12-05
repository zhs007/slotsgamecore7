package lowcode

import (
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
	"gopkg.in/yaml.v2"
)

// CollectorPayTriggerTypeName is the component type name for the collector pay trigger.
const CollectorPayTriggerTypeName = "collectorPayTrigger"

type mainSymbolInfo struct {
	symbolCode int
	level      int
	price      int
	moved      bool
	syms       []int
}

type CollectorPayTriggerData struct {
	BasicComponentData
	lstMainSymbols []*mainSymbolInfo
	cfg            *CollectorPayTriggerConfig
}

func (cpt *CollectorPayTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	cpt.BasicComponentData.OnNewGame(gameProp, component)

	for _, ms := range cpt.lstMainSymbols {
		ms.level = 0
		sc := ms.syms[0]
		ms.price = gameProp.Pool.DefaultPaytables.MapPay[sc][0]
	}
}

func (cpt *CollectorPayTriggerData) onLevelUp(mainSymbol int, off int) {
	for _, ms := range cpt.lstMainSymbols {
		if ms.symbolCode == mainSymbol {
			ms.level += off

			if ms.level >= len(cpt.cfg.MapSymbolCode[ms.symbolCode]) {
				ms.level = len(cpt.cfg.MapSymbolCode[ms.symbolCode]) - 1
			}

			ms.price = cpt.cfg.MapSymbolCode[ms.symbolCode][ms.level]
		}
	}
}

func (cpt *CollectorPayTriggerData) onAllLevelUp(off int) {
	for _, ms := range cpt.lstMainSymbols {
		ms.level += off

		if ms.level >= len(cpt.cfg.MapSymbolCode[ms.symbolCode]) {
			ms.level = len(cpt.cfg.MapSymbolCode[ms.symbolCode]) - 1
		}

		ms.price = cpt.cfg.MapSymbolCode[ms.symbolCode][ms.level]
	}
}

func (cpt *CollectorPayTriggerData) getSymbolCode(ms int) int {

	for _, msi := range cpt.lstMainSymbols {
		if msi.symbolCode == ms {
			return cpt.cfg.MapSymbolCode[ms][msi.level]
		}
	}

	return -1
}

func (cpt *CollectorPayTriggerData) getMainSymbolInfo(ms int) *mainSymbolInfo {

	for _, msi := range cpt.lstMainSymbols {
		if msi.symbolCode == ms {
			return msi
		}
	}

	return nil
}

func (cpt *CollectorPayTriggerData) getNext() int {
	msc := -1
	msp := -1

	for _, ms := range cpt.lstMainSymbols {
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

func (cpt *CollectorPayTriggerData) moveEnd(ms int) {
	for _, msi := range cpt.lstMainSymbols {
		if msi.symbolCode == ms {
			msi.moved = true

			return
		}
	}
}

func (cpt *CollectorPayTriggerData) procSymbolsWithLevel(gs *sgc7game.GameScene) {
	for x, arr := range gs.Arr {
		for y, sc := range arr {
			for _, ms := range cpt.lstMainSymbols {
				if slices.Contains(ms.syms, sc) {
					gs.Arr[x][y] = ms.syms[ms.level]
				}
			}
		}
	}
}

func (cpt *CollectorPayTriggerData) onNewStep() {
	cpt.UsedScenes = nil
	cpt.UsedResults = nil
	cpt.Pos = nil

	for _, ms := range cpt.lstMainSymbols {
		ms.moved = false
	}
}

func (cpt *CollectorPayTriggerData) Clone() IComponentData {
	target := &CollectorPayTriggerData{
		BasicComponentData: cpt.CloneBasicComponentData(),
		lstMainSymbols:     slices.Clone(cpt.lstMainSymbols),
	}

	return target
}

// CollectorPayTriggerConfig - configuration for CollectorPayTrigger
// CollectorPayTriggerConfig is the configuration for the CollectorPayTrigger component.
type CollectorPayTriggerConfig struct {
	BasicComponentConfig   `yaml:",inline" json:",inline"`
	CategoryCount          int                 `yaml:"categoryCount" json:"categoryCount"`
	MapSymbol              map[string][]string `yaml:"mapSymbol" json:"mapSymbol"`
	MapSymbolCode          map[int][]int       `yaml:"-" json:"-"`
	BlankSymbol            string              `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode        int                 `yaml:"-" json:"-"`
	WildSymbol             string              `yaml:"wildSymbol" json:"wildSymbol"`
	WildSymbolCode         int                 `yaml:"-" json:"-"`
	CoinSymbols            []string            `yaml:"coinSymbols" json:"coinSymbols"`
	CoinSymbolCodes        []int               `yaml:"-" json:"-"`
	UpLevelSymbols         []string            `yaml:"upLevelSymbols" json:"upLevelSymbols"`
	UpLevelSymbolCodes     []int               `yaml:"-" json:"-"`
	AllUpLevelSymbols      []string            `yaml:"allUpLevelSymbol" json:"allUpLevelSymbol"`
	AllUpLevelSymbolCodes  []int               `yaml:"-" json:"-"`
	SwitcherSymbol         string              `yaml:"switcherSymbol" json:"switcherSymbol"`
	SwitcherSymbolCode     int                 `yaml:"-" json:"-"`
	PopcornSymbol          string              `yaml:"popcornSymbol" json:"popcornSymbol"`
	PopcornSymbolCode      int                 `yaml:"-" json:"-"`
	EggSymbol              string              `yaml:"eggSymbol" json:"eggSymbol"`
	EggSymbolCode          int                 `yaml:"-" json:"-"`
	DontPressSymbol        string              `yaml:"dontpressSymbol" json:"dontpressSymbol"`
	DontPressSymbolCode    int                 `yaml:"-" json:"-"`
	TriggerOnlySymbols     []string            `yaml:"triggerOnlySymbols" json:"triggerOnlySymbols"`
	TriggerOnlySymbolCodes []int               `yaml:"-" json:"-"`
	HighLevelSPSymbolCount int                 `yaml:"highLevelSPSymbolCount" json:"highLevelSPSymbolCount"`
	HighLevelSPSymbols     []string            `yaml:"highLevelSPSymbols" json:"highLevelSPSymbols"`
	HighLevelSPSymbolCodes []int               `yaml:"-" json:"-"`
	LowLevelSPSymbolCount  int                 `yaml:"lowLevelSPSymbolCount" json:"lowLevelSPSymbolCount"`
	LowLevelSPSymbols      []string            `yaml:"lowLevelSPSymbols" json:"lowLevelSPSymbols"`
	LowLevelSPSymbolCodes  []int               `yaml:"-" json:"-"`
	JumpToComponent        string              `yaml:"jump" json:"jump"`
	mapSymbolValues        map[int]int         `yaml:"-" json:"-"`
	lstMainSymbols         []int               `yaml:"-" json:"-"`
	MapControllers         map[string][]*Award `yaml:"mapControllers" json:"mapControllers"` // 新的奖励系统
}

// SetLinkComponent sets a link ("next" or "jump") to another component by name.
func (cfg *CollectorPayTriggerConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

// CollectorPayTrigger is the runtime component that implements collector pay trigger logic.
type CollectorPayTrigger struct {
	*BasicComponent `json:"-"`
	Config          *CollectorPayTriggerConfig `json:"config"`
}

// Init - load from file
// Init loads the component configuration from a YAML file and initializes the component.
func (cpt *CollectorPayTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CollectorPayTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CollectorPayTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CollectorPayTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return cpt.InitEx(cfg, pool)
}

// InitEx - initialize from config object
// InitEx initializes the component from an in-memory config object (usually unmarshaled from YAML).
func (cpt *CollectorPayTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	cpt.Config = cfg.(*CollectorPayTriggerConfig)
	cpt.Config.ComponentType = CollectorPayTriggerTypeName

	if cpt.Config.BlankSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpt.Config.BlankSymbol]
		if !isok {
			goutils.Error("CollectorPayTrigger.InitEx:BlankSymbol",
				slog.String("BlankSymbol", cpt.Config.BlankSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpt.Config.BlankSymbolCode = sc
	} else {
		cpt.Config.BlankSymbolCode = -1
	}

	if cpt.Config.WildSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpt.Config.WildSymbol]
		if !isok {
			goutils.Error("CollectorPayTrigger.InitEx:WildSymbol",
				slog.String("WildSymbol", cpt.Config.WildSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpt.Config.WildSymbolCode = sc
	} else {
		cpt.Config.WildSymbolCode = -1
	}

	if cpt.Config.SwitcherSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpt.Config.SwitcherSymbol]
		if !isok {
			goutils.Error("CollectorPayTrigger.InitEx:SwitcherSymbol",
				slog.String("SwitcherSymbol", cpt.Config.SwitcherSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpt.Config.SwitcherSymbolCode = sc
	} else {
		cpt.Config.SwitcherSymbolCode = -1
	}

	if cpt.Config.PopcornSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpt.Config.PopcornSymbol]
		if !isok {
			goutils.Error("CollectorPayTrigger.InitEx:PopcornSymbol",
				slog.String("PopcornSymbol", cpt.Config.PopcornSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpt.Config.PopcornSymbolCode = sc
	} else {
		cpt.Config.PopcornSymbolCode = -1
	}

	if cpt.Config.EggSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpt.Config.EggSymbol]
		if !isok {
			goutils.Error("CollectorPayTrigger.InitEx:EggSymbol",
				slog.String("EggSymbol", cpt.Config.EggSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpt.Config.EggSymbolCode = sc
	} else {
		cpt.Config.EggSymbolCode = -1
	}

	if cpt.Config.DontPressSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cpt.Config.DontPressSymbol]
		if !isok {
			goutils.Error("CollectorPayTrigger.InitEx:DontPressSymbol",
				slog.String("DontPressSymbol", cpt.Config.DontPressSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		cpt.Config.DontPressSymbolCode = sc
	} else {
		cpt.Config.DontPressSymbolCode = -1
	}

	if len(cpt.Config.CoinSymbols) > 0 {
		cpt.Config.CoinSymbolCodes = make([]int, len(cpt.Config.CoinSymbols))
		for i, cs := range cpt.Config.CoinSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:CoinSymbols",
					slog.String("CoinSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpt.Config.CoinSymbolCodes[i] = sc
		}
	}

	if len(cpt.Config.UpLevelSymbols) > 0 {
		cpt.Config.UpLevelSymbolCodes = make([]int, len(cpt.Config.UpLevelSymbols))
		for i, cs := range cpt.Config.UpLevelSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:UpLevelSymbols",
					slog.String("UpLevelSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpt.Config.UpLevelSymbolCodes[i] = sc
		}
	}

	if len(cpt.Config.AllUpLevelSymbols) > 0 {
		cpt.Config.AllUpLevelSymbolCodes = make([]int, len(cpt.Config.AllUpLevelSymbols))
		for i, cs := range cpt.Config.AllUpLevelSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:AllUpLevelSymbols",
					slog.String("AllUpLevelSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpt.Config.AllUpLevelSymbolCodes[i] = sc
		}
	}

	if len(cpt.Config.TriggerOnlySymbols) > 0 {
		cpt.Config.TriggerOnlySymbolCodes = make([]int, len(cpt.Config.TriggerOnlySymbols))
		for i, cs := range cpt.Config.TriggerOnlySymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:TriggerOnlySymbols",
					slog.String("TriggerOnlySymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpt.Config.TriggerOnlySymbolCodes[i] = sc
		}
	}

	if len(cpt.Config.HighLevelSPSymbols) > 0 {
		cpt.Config.HighLevelSPSymbolCodes = make([]int, len(cpt.Config.HighLevelSPSymbols))
		for i, cs := range cpt.Config.HighLevelSPSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:HighLevelSPSymbols",
					slog.String("HighLevelSPSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpt.Config.HighLevelSPSymbolCodes[i] = sc
		}
	}

	if len(cpt.Config.LowLevelSPSymbols) > 0 {
		cpt.Config.LowLevelSPSymbolCodes = make([]int, len(cpt.Config.LowLevelSPSymbols))
		for i, cs := range cpt.Config.LowLevelSPSymbols {
			sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:LowLevelSPSymbols",
					slog.String("LowLevelSPSymbol", cs),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cpt.Config.LowLevelSPSymbolCodes[i] = sc
		}
	}

	if len(cpt.Config.MapSymbol) > 0 {
		cpt.Config.MapSymbolCode = make(map[int][]int)
		for ms, css := range cpt.Config.MapSymbol {
			mssc, isok := pool.Config.GetDefaultPaytables().MapSymbols[ms]
			if !isok {
				goutils.Error("CollectorPayTrigger.InitEx:MapSymbol",
					slog.String("MainSymbol", ms),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}

			cssc := make([]int, len(css))

			for i, cs := range css {
				sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
				if !isok {
					goutils.Error("CollectorPayTrigger.InitEx:MapSymbol:CollectedSymbols",
						slog.String("CollectedSymbol", cs),
						goutils.Err(ErrInvalidComponentConfig))

					return ErrInvalidComponentConfig
				}

				cssc[i] = sc
			}

			cpt.Config.MapSymbolCode[mssc] = cssc

			cpt.Config.lstMainSymbols = append(cpt.Config.lstMainSymbols, mssc)
		}

		sort.Slice(cpt.Config.lstMainSymbols, func(i, j int) bool { return cpt.Config.lstMainSymbols[i] < cpt.Config.lstMainSymbols[j] })
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

	cpt.Config.mapSymbolValues = make(map[int]int)

	cpt.Config.mapSymbolValues[cpt.Config.DontPressSymbolCode] = 2
	cpt.Config.mapSymbolValues[cpt.Config.EggSymbolCode] = 3

	for k, arr := range cpt.Config.MapSymbolCode {
		cpt.Config.mapSymbolValues[k] = 1

		for _, a := range arr {
			cpt.Config.mapSymbolValues[a] = 4
		}
	}

	cpt.Config.mapSymbolValues[cpt.Config.WildSymbolCode] = 5
	cpt.Config.mapSymbolValues[cpt.Config.PopcornSymbolCode] = 6
	cpt.Config.mapSymbolValues[cpt.Config.SwitcherSymbolCode] = 7
	for _, v := range cpt.Config.UpLevelSymbolCodes {
		cpt.Config.mapSymbolValues[v] = 8
	}
	for _, v := range cpt.Config.AllUpLevelSymbolCodes {
		cpt.Config.mapSymbolValues[v] = 9
	}
	for _, v := range cpt.Config.CoinSymbolCodes {
		cpt.Config.mapSymbolValues[v] = 10
	}

	for _, arr := range cpt.Config.MapControllers {
		for _, a := range arr {
			if a != nil {
				a.Init()
			}
		}
	}

	cpt.onInit(&cpt.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (cpt *CollectorPayTrigger) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := cpt.Config.MapControllers[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

func (cpt *CollectorPayTrigger) calcSymbolCodeValue(symbol int, mainSymbol int) int {
	v, isok := cpt.Config.mapSymbolValues[symbol]
	if isok {
		if v == 4 && !slices.Contains(cpt.Config.MapSymbolCode[mainSymbol], symbol) {
			return -1
		}

		if v == 1 && symbol != mainSymbol {
			return -1
		}

		return v
	}

	return -1
}

// 计算这个位置的价值,如果已经走过了,返回-1,如果是空位置,且isIgnoreEmpty为true,也返回-1,如果isIgnoreEmpty为false,返回0
func (cpt *CollectorPayTrigger) calcTileValue(gs *sgc7game.GameScene, mainSymbol int, x, y int, pd *PosData, isIgnoreEmpty bool) int {
	if pd.Has(x, y) {
		return -1
	}

	if gs.Arr[x][y] == -3 {
		return -1
	}

	if gs.Arr[x][y] == -2 {
		return 0
	}

	if gs.Arr[x][y] == -1 {
		if isIgnoreEmpty {
			return 0
		}

		return -1
	}

	return cpt.calcSymbolCodeValue(gs.Arr[x][y], mainSymbol)
}

func (cpt *CollectorPayTrigger) isCanEnding(gs *sgc7game.GameScene, mainSymbol int, x, y int, pd *PosData, isIgnoreEmpty bool) bool {
	if y < gs.Height-1 {
		cv := cpt.calcTileValue(gs, mainSymbol, x, y+1, pd, isIgnoreEmpty)
		if cv >= 0 {
			return false
		}
	}

	if y > 0 {
		cv := cpt.calcTileValue(gs, mainSymbol, x, y-1, pd, isIgnoreEmpty)
		if cv >= 0 {
			return false
		}
	}

	if x > 0 {
		cv := cpt.calcTileValue(gs, mainSymbol, x-1, y, pd, isIgnoreEmpty)
		if cv >= 0 {
			return false
		}
	}

	if x < gs.Width-1 {
		cv := cpt.calcTileValue(gs, mainSymbol, x+1, y, pd, isIgnoreEmpty)
		if cv >= 0 {
			return false
		}
	}

	return true
}

// 找到一条路径,返回该路径上的最大价值以及路径数组,找到第一个大于普通符号的位置就返回
func (cpt *CollectorPayTrigger) findPathDeep(gameProp *GameProperty, gs *sgc7game.GameScene, mainSymbol int, x, y int, pd *PosData, isIgnoreEmpty bool) (int, *PosData) {
	csv := cpt.calcTileValue(gs, mainSymbol, x, y, pd, isIgnoreEmpty)
	if csv < 0 {
		return -1, nil
	}

	if csv > 4 {
		npd := gameProp.posPool.Clone(pd)
		npd.Add(x, y)

		return csv, npd
	}

	if csv != 4 && csv != 0 {
		npd := gameProp.posPool.Clone(pd)
		npd.Add(x, y)

		return csv, npd
	}

	tmppd := gameProp.posPool.Clone(pd)
	tmppd.Add(x, y)

	if cpt.isCanEnding(gs, mainSymbol, x, y, tmppd, isIgnoreEmpty) {
		return csv, tmppd
	}

	maxv := -1
	var maxvpd *PosData

	if y < gs.Height-1 {
		cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x, y+1, tmppd, isIgnoreEmpty)
		if cdv >= 0 {
			if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
				maxv = cdv

				maxvpd = cdpd
			}
		}
	}

	if y > 0 {
		if maxvpd == nil || !maxvpd.Has(x, y-1) {
			cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x, y-1, tmppd, isIgnoreEmpty)
			if cdv >= 0 {
				if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
					maxv = cdv

					maxvpd = cdpd
				}
			}
		}
	}

	if x > 0 {
		if maxvpd == nil || !maxvpd.Has(x-1, y) {
			cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x-1, y, tmppd, isIgnoreEmpty)
			if cdv >= 0 {
				if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
					maxv = cdv

					maxvpd = cdpd
				}
			}
		}
	}

	if x < gs.Width-1 {
		if maxvpd == nil || !maxvpd.Has(x+1, y) {
			cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x+1, y, tmppd, isIgnoreEmpty)
			if cdv >= 0 {
				if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
					maxv = cdv

					maxvpd = cdpd
				}
			}
		}
	}

	if maxv >= 0 {
		return maxv, maxvpd
	}

	return -1, nil
}

// 找到一条路径,返回该路径上的最大价值以及路径数组,找到第一个大于普通符号的位置就返回
func (cpt *CollectorPayTrigger) findPath(gameProp *GameProperty, gs *sgc7game.GameScene, mainSymbol int, x, y int, isIgnoreEmpty bool) (int, *PosData) {
	pd := gameProp.posPool.Get()

	pd.Add(x, y)

	maxv := -1
	var maxvpd *PosData

	if y < gs.Height-1 {
		cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x, y+1, pd, isIgnoreEmpty)
		if cdv > 0 {
			if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
				maxv = cdv

				maxvpd = cdpd
			}
		}
	}

	if y > 0 {
		if maxvpd == nil || !maxvpd.Has(x, y-1) {
			cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x, y-1, pd, isIgnoreEmpty)
			if cdv > 0 {
				if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
					maxv = cdv

					maxvpd = cdpd
				}
			}
		}
	}

	if x < gs.Width-1 {
		if maxvpd == nil || !maxvpd.Has(x+1, y) {
			cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x+1, y, pd, isIgnoreEmpty)
			if cdv > 0 {
				if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
					maxv = cdv

					maxvpd = cdpd
				}
			}
		}
	}

	if x > 0 {
		if maxvpd == nil || !maxvpd.Has(x-1, y) {
			cdv, cdpd := cpt.findPathDeep(gameProp, gs, mainSymbol, x-1, y, pd, isIgnoreEmpty)
			if cdv > 0 {
				if cdv > maxv || (maxv == cdv && len(maxvpd.pos) > len(cdpd.pos)) {
					maxv = cdv

					maxvpd = cdpd
				}
			}
		}
	}

	return maxv, maxvpd
}

func (cpt *CollectorPayTrigger) rechgScene(gs *sgc7game.GameScene) {
	for x := 0; x < gs.Width; x++ {
		for y := 0; y < gs.Height; y++ {
			if gs.Arr[x][y] < 0 {
				gs.Arr[x][y] = -1
			}
		}
	}
}

func (cpt *CollectorPayTrigger) reinitScene(gs *sgc7game.GameScene, cd *CollectorPayTriggerData) {
	for x := 0; x < gs.Width; x++ {
		for y := 0; y < gs.Height; y++ {
			if gs.Arr[x][y] < 0 {
				gs.Arr[x][y] = -3

				continue
			}

			for _, msi := range cd.lstMainSymbols {
				if slices.Contains(msi.syms, gs.Arr[x][y]) {
					gs.Arr[x][y] = msi.syms[msi.level]

					break
				}
			}
		}
	}
}

func (cpt *CollectorPayTrigger) calcVal(gs *sgc7game.GameScene, x, y int, ms int, syms []int) int {
	if gs.Arr[x][y] == -2 {
		return 0
	}

	if gs.Arr[x][y] == cpt.Config.SwitcherSymbolCode {
		return 9
	}

	if slices.Contains(cpt.Config.AllUpLevelSymbolCodes, gs.Arr[x][y]) {
		return 8
	}

	if slices.Contains(cpt.Config.UpLevelSymbolCodes, gs.Arr[x][y]) {
		return 7
	}

	if slices.Contains(cpt.Config.CoinSymbolCodes, gs.Arr[x][y]) {
		return 6
	}

	if gs.Arr[x][y] == cpt.Config.PopcornSymbolCode {
		return 5
	}

	if slices.Contains(syms, gs.Arr[x][y]) {
		return 2
	}

	if gs.Arr[x][y] == ms {
		return 1
	}

	if gs.Arr[x][y] == cpt.Config.EggSymbolCode {
		return 3
	}

	if gs.Arr[x][y] == cpt.Config.DontPressSymbolCode {
		return 4
	}

	return 0
}

func (cpt *CollectorPayTrigger) isNearVal(cv int) bool {
	if cv >= 4 {
		return true
	}

	return false
}

func (cpt *CollectorPayTrigger) checkDistente(vs *sgc7game.GameScene, ds *sgc7game.GameScene, x, y int, d int) {
	if vs.Arr[x][y] != -1 {
		if ds.Arr[x][y] > d {
			ds.Arr[x][y] = d
			// up
			if y < vs.Height-1 {
				cpt.checkDistente(vs, ds, x, y+1, d+1)
			}

			// down
			if y > 0 {
				cpt.checkDistente(vs, ds, x, y-1, d+1)
			}

			// left
			if x > 0 {
				cpt.checkDistente(vs, ds, x-1, y, d+1)
			}

			// right
			if x < vs.Width-1 {
				cpt.checkDistente(vs, ds, x+1, y, d+1)
			}
		}
	}
}

func (cpt *CollectorPayTrigger) checkVal(gs *sgc7game.GameScene, vs *sgc7game.GameScene, ds *sgc7game.GameScene, ms int, x, y int, syms []int, d int) bool {
	if vs.Arr[x][y] != -1 {
		cpt.checkDistente(vs, ds, x, y, d)

		return false
	}

	val := cpt.calcVal(gs, x, y, ms, syms)
	if val == 0 {
		return false
	}

	vs.Arr[x][y] = val
	ds.Arr[x][y] = d

	// up
	if y < gs.Height-1 {
		cpt.checkVal(gs, vs, ds, ms, x, y+1, syms, d+1)
	}

	// down
	if y > 0 {
		cpt.checkVal(gs, vs, ds, ms, x, y-1, syms, d+1)
	}

	// left
	if x > 0 {
		cpt.checkVal(gs, vs, ds, ms, x-1, y, syms, d+1)
	}

	// right
	if x < gs.Width-1 {
		cpt.checkVal(gs, vs, ds, ms, x+1, y, syms, d+1)
	}

	return true
}

func (cpt *CollectorPayTrigger) genValData(gs *sgc7game.GameScene, vs *sgc7game.GameScene, ds *sgc7game.GameScene, mx, my int, ms int, syms []int) (*sgc7game.GameScene, *sgc7game.GameScene) {
	vs.Clear(-1)
	ds.Clear(-1)

	cpt.checkVal(gs, vs, ds, ms, mx, my, syms, 0)

	return vs, ds
}

func (cpt *CollectorPayTrigger) findTarget(vs *sgc7game.GameScene, ds *sgc7game.GameScene) (int, int) {
	cv := -1
	cd := -1
	cx := -1
	cy := -1

	for x := 0; x < vs.Width; x++ {
		for y := 0; y < vs.Height; y++ {
			if vs.Arr[x][y] > cv {
				cv = vs.Arr[x][y]
				cd = ds.Arr[x][y]
				cx = x
				cy = y
			} else if cv == 2 && vs.Arr[x][y] == cv {
				if ds.Arr[x][y] > cd {
					cd = ds.Arr[x][y]
					cx = x
					cy = y
				}
			} else if vs.Arr[x][y] == cv && cv > 0 {
				if ds.Arr[x][y] < cd {
					cd = ds.Arr[x][y]
					cx = x
					cy = y
				}
			}
		}
	}

	return cx, cy
}

func (cpt *CollectorPayTrigger) findNearPath(ds *sgc7game.GameScene, mainSymbol int, x, y int, pos *PosData) bool {
	pos.Add(x, y)

	d := ds.Arr[x][y]

	if d == 0 {
		return true
	}

	// up
	if y < ds.Height-1 {
		if ds.Arr[x][y+1] == d-1 {
			if cpt.findNearPath(ds, mainSymbol, x, y+1, pos) {
				return true
			}
		}
	}

	// down
	if y > 0 {
		if ds.Arr[x][y-1] == d-1 {
			if cpt.findNearPath(ds, mainSymbol, x, y-1, pos) {
				return true
			}
		}
	}

	// left
	if x > 0 {
		if ds.Arr[x-1][y] == d-1 {
			if cpt.findNearPath(ds, mainSymbol, x-1, y, pos) {
				return true
			}
		}
	}

	// right
	if x < ds.Width-1 {
		if ds.Arr[x+1][y] == d-1 {
			if cpt.findNearPath(ds, mainSymbol, x+1, y, pos) {
				return true
			}
		}
	}

	return false
}

func (cpt *CollectorPayTrigger) procAllUpLevel(gs *sgc7game.GameScene, off int, cd *CollectorPayTriggerData) {

	cd.onAllLevelUp(off)

	for x := 0; x < gs.Width; x++ {
		for y := 0; y < gs.Height; y++ {
			for _, msi := range cd.lstMainSymbols {
				cursym := msi.syms[msi.level]

				if slices.Contains(msi.syms, gs.Arr[x][y]) {
					gs.Arr[x][y] = cursym
				}
			}
		}
	}
}

func (cpt *CollectorPayTrigger) procUpLevel(gs *sgc7game.GameScene, ms int, off int, cd *CollectorPayTriggerData) {

	cd.onLevelUp(ms, off)

	msi := cd.getMainSymbolInfo(ms)

	for x := 0; x < gs.Width; x++ {
		for y := 0; y < gs.Height; y++ {

			if slices.Contains(msi.syms, gs.Arr[x][y]) {
				gs.Arr[x][y] = msi.syms[msi.level]
			}
		}
	}
}

func (cpt *CollectorPayTrigger) procNear(gameProp *GameProperty, bet int, ms int, sx, sy int, dx, dy int, vs *sgc7game.GameScene,
	ds *sgc7game.GameScene, gs *sgc7game.GameScene, os *sgc7game.GameScene, syms []int, curpr *sgc7game.PlayResult,
	cd *CollectorPayTriggerData) (*sgc7game.GameScene, int, int, error) {

	pos := gameProp.posPool.Get()

	isok := cpt.findNearPath(ds, vs.Arr[sx][sy], dx, dy, pos)
	if !isok {
		goutils.Error("CollectorPayTrigger.procNear:findNearPath",
			slog.Int("sx", sx),
			slog.Int("sy", sy),
			slog.Int("dx", dx),
			slog.Int("dy", dy))

		return nil, -1, -1, ErrInvalidComponentData
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	if len(cd.Pos) != 0 {
		cd.AddPos(-1, -1)
	}

	ex := sx
	ey := sy

	curRet := &sgc7game.Result{
		Symbol:    ms,
		Type:      sgc7game.RTCollectorPay,
		LineIndex: -1,
	}

	for i := pos.Len() - 1; i >= 0; i-- {
		x := pos.pos[i*2]
		y := pos.pos[i*2+1]

		cd.AddPos(x, y)

		if slices.Contains(syms, ngs.Arr[x][y]) || ngs.Arr[x][y] == ms {
			ngs.Arr[x][y] = -2
			os.Arr[x][y] = -1

			ex = x
			ey = y

			curRet.Pos = append(curRet.Pos, x, y)
		} else if ngs.Arr[x][y] == cpt.Config.WildSymbolCode {
			curRet.Pos = append(curRet.Pos, x, y)
		}

		alluplevelIndex := slices.Index(cpt.Config.AllUpLevelSymbolCodes, ngs.Arr[x][y])
		if alluplevelIndex >= 0 {
			cpt.procAllUpLevel(ngs, alluplevelIndex+1, cd)

			ngs.Arr[x][y] = -2
			os.Arr[x][y] = -1

			ex = x
			ey = y

			curRet.Pos = append(curRet.Pos, x, y)
		}

		uplevelIndex := slices.Index(cpt.Config.UpLevelSymbolCodes, ngs.Arr[x][y])
		if uplevelIndex >= 0 {
			cpt.procUpLevel(ngs, ms, uplevelIndex+1, cd)

			ngs.Arr[x][y] = -2
			os.Arr[x][y] = -1

			ex = x
			ey = y

			curRet.Pos = append(curRet.Pos, x, y)
		}

		if slices.Contains(cpt.Config.CoinSymbolCodes, ngs.Arr[x][y]) {
			ngs.Arr[x][y] = -2

			if os.Arr[x][y] <= 0 {
				goutils.Error("CollectorPayTrigger.procNear:InvalidCoinValue",
					slog.Int("x", x),
					slog.Int("y", y),
					slog.Int("CoinSymbol", ngs.Arr[x][y]),
					slog.Int("CoinValue", os.Arr[x][y]))

				return nil, -1, -1, ErrInvalidComponentData
			}

			ret := &sgc7game.Result{
				Symbol:    -1,
				Type:      sgc7game.RTBonus,
				LineIndex: -1,
				CoinWin:   os.Arr[x][y],
				CashWin:   os.Arr[x][y] * bet,
			}

			os.Arr[x][y] = -1

			ex = x
			ey = y

			curRet.Pos = append(curRet.Pos, x, y)

			cpt.AddResult(curpr, ret, &cd.BasicComponentData)
		}

		if ngs.Arr[x][y] == cpt.Config.SwitcherSymbolCode || ngs.Arr[x][y] == cpt.Config.PopcornSymbolCode || ngs.Arr[x][y] == cpt.Config.EggSymbolCode || ngs.Arr[x][y] == cpt.Config.DontPressSymbolCode {
			ngs.Arr[x][y] = -2
			os.Arr[x][y] = -1

			ex = x
			ey = y

			curRet.Pos = append(curRet.Pos, x, y)
		}
	}

	if len(curRet.Pos) > 0 {
		curRet.CoinWin = len(curRet.Pos) * bet
		cpt.AddResult(curpr, curRet, &cd.BasicComponentData)
	}

	ngs.Arr[ex][ey] = ms

	cpt.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

	return ngs, ex, ey, nil
}

// procSymbolsWithPos
func (cpt *CollectorPayTrigger) procCollect(gameProp *GameProperty, bet int, curpr *sgc7game.PlayResult, gs *sgc7game.GameScene, os *sgc7game.GameScene, cd *CollectorPayTriggerData) error {
	ngs := gs

	vs := gameProp.PoolScene.New2(gs.Width, gs.Height, -1)
	ds := gameProp.PoolScene.New2(gs.Width, gs.Height, -1)

	for {
		mainSymbol := cd.getNext()
		if mainSymbol == -1 {
			break
		}

		arr := cpt.Config.MapSymbolCode[mainSymbol]

		for x := 0; x < ngs.Width; x++ {
			for y := 0; y < ngs.Height; y++ {
				if ngs.Arr[x][y] == mainSymbol {
					sx := x
					sy := y

					for {
						// 先计算价值和距离,然后按价值高低,一个个的获取,有走最短路径和最长路径2个方案,可能对局面造成影响的几个符号走最短路径,剩下的走最长路径
						vs, ds := cpt.genValData(ngs, vs, ds, sx, sy, mainSymbol, arr)

						tx, ty := cpt.findTarget(vs, ds)
						if tx == -1 && ty == -1 {
							break
						}

						if sx == tx && sy == ty {
							break
						}

						// cv := vs.Arr[tx][ty]

						// if cpt.isNearVal(cv) {
						cngs, ex, ey, err := cpt.procNear(gameProp, bet, mainSymbol, sx, sy, tx, ty, vs, ds, ngs, os, arr, curpr, cd)
						if err != nil {
							goutils.Error("CollectorPayTrigger.procCollect:procNear",
								slog.Int("sx", sx),
								slog.Int("sy", sy),
								slog.Int("tx", tx),
								slog.Int("ty", ty),
								goutils.Err(err))

							return err
						}

						ngs = cngs

						sx = ex
						sy = ey
						// } else {
						// 	break
						// }

						// sx = tx
						// sy = ty
					}
				}
			}
		}

		cd.moveEnd(mainSymbol)
	}

	cpt.rechgScene(ngs)

	return nil
}

// OnPlayGame - check collector value and proc awards when reach threshold
// OnPlayGame processes a play event for the component and returns the next component name (if any).
func (cpt *CollectorPayTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*CollectorPayTriggerData)
	cd.onNewStep()

	gs := cpt.GetTargetScene3(gameProp, curpr, prs, 0)
	os := cpt.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	nos := os.CloneEx(gameProp.PoolScene)

	ngs := gs.CloneEx(gameProp.PoolScene)
	cpt.reinitScene(ngs, cd)

	bet := gameProp.GetBet3(stake, BTypeBet)

	cpt.procCollect(gameProp, bet, curpr, ngs, nos, cd)

	if len(cd.UsedResults) > 0 {
		cpt.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")
	}

	cpt.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)

	nc := cpt.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
// OnAsciiGame outputs the component state to an asciigame representation (no-op for this component).
func (cpt *CollectorPayTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func (cpt *CollectorPayTrigger) NewComponentData() IComponentData {
	cd := &CollectorPayTriggerData{
		cfg: cpt.Config,
	}

	for _, ms := range cpt.Config.lstMainSymbols {
		cd.lstMainSymbols = append(cd.lstMainSymbols, &mainSymbolInfo{
			symbolCode: ms,
			level:      0,
			price:      0,
			syms:       cpt.Config.MapSymbolCode[ms],
		})
	}

	return cd
}

// NewCollectorPayTrigger creates a new CollectorPayTrigger component instance.
func NewCollectorPayTrigger(name string) IComponent {
	return &CollectorPayTrigger{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "categoryCount": 4,
// "highLevelSPSymbolCount": 1,
// "lowLevelSPSymbolCount": 1,
// "mapSymbol": [
//
//	{
//	    "mainSymbol": "RP",
//	    "collectedSymbols": [
//	        "R1",
//	        "R2",
//	        "R3",
//	        "R4",
//	        "R5",
//	        "R6",
//	        "R7"
//	    ]
//	},
//	{
//	    "mainSymbol": "PP",
//	    "collectedSymbols": [
//	        "P1",
//	        "P2",
//	        "P3",
//	        "P4",
//	        "P5",
//	        "P6",
//	        "P7"
//	    ]
//	},
//	{
//	    "mainSymbol": "GP",
//	    "collectedSymbols": [
//	        "G1",
//	        "G2",
//	        "G3",
//	        "G4",
//	        "G6",
//	        "G5",
//	        "G7"
//	    ]
//	},
//	{
//	    "mainSymbol": "BP",
//	    "collectedSymbols": [
//	        "B1",
//	        "B2",
//	        "B3",
//	        "B4",
//	        "B5",
//	        "B6",
//	        "B7"
//	    ]
//	}
//
// ],
// "wildSymbol": "WL",
// "coinSymbols": [
//
//	"CN"
//
// ],
// "upLevelSymbol": [
//
//	"L1",
//	"L2",
//	"L3"
//
// ],
// "allUpLevelSymbol": [
//
//	"AL1",
//	"AL2",
//	"AL3"
//
// ],
// "switcherSymbol": "MR",
// "popcornSymbol": "PC",
// "eggSymbol": "EG",
// "dontpressSymbol": "DP",
// "highLevelSPSymbol": [
//
//	"MR"
//
// ],
// "lowLevelSPSymbol": [
//
//	"DP"
//
// ],
// "triggerOnlySymbols": [
//
//	"WL",
//	"DP",
//	"EG"
//
// ]

type jsonCPTSymbolData struct {
	MainSymbol       string   `json:"mainSymbol"`
	CollectedSymbols []string `json:"collectedSymbols"`
}

type jsonCollectorPayTrigger struct {
	CategoryCount          int                 `json:"categoryCount"`
	MapSymbol              []jsonCPTSymbolData `json:"mapSymbol"`
	BlankSymbol            string              `json:"blankSymbol"`
	WildSymbol             string              `json:"wildSymbol"`
	CoinSymbols            []string            `json:"coinSymbols"`
	UpLevelSymbol          []string            `json:"upLevelSymbol"`
	AllUpLevelSymbol       []string            `json:"allUpLevelSymbol"`
	SwitcherSymbol         string              `json:"switcherSymbol"`
	PopcornSymbol          string              `json:"popcornSymbol"`
	EggSymbol              string              `json:"eggSymbol"`
	DontPressSymbol        string              `json:"dontpressSymbol"`
	TriggerOnlySymbols     []string            `json:"triggerOnlySymbols"`
	HighLevelSPSymbolCount int                 `json:"highLevelSPSymbolCount"`
	HighLevelSPSymbol      []string            `json:"highLevelSPSymbol"`
	LowLevelSPSymbolCount  int                 `json:"lowLevelSPSymbolCount"`
	LowLevelSPSymbol       []string            `json:"lowLevelSPSymbol"`
}

func (j *jsonCollectorPayTrigger) build() *CollectorPayTriggerConfig {
	cfg := &CollectorPayTriggerConfig{
		CategoryCount:          j.CategoryCount,
		MapSymbol:              make(map[string][]string),
		WildSymbol:             j.WildSymbol,
		CoinSymbols:            j.CoinSymbols,
		UpLevelSymbols:         j.UpLevelSymbol,
		AllUpLevelSymbols:      j.AllUpLevelSymbol,
		SwitcherSymbol:         j.SwitcherSymbol,
		PopcornSymbol:          j.PopcornSymbol,
		EggSymbol:              j.EggSymbol,
		DontPressSymbol:        j.DontPressSymbol,
		TriggerOnlySymbols:     j.TriggerOnlySymbols,
		HighLevelSPSymbolCount: j.HighLevelSPSymbolCount,
		HighLevelSPSymbols:     j.HighLevelSPSymbol,
		LowLevelSPSymbolCount:  j.LowLevelSPSymbolCount,
		LowLevelSPSymbols:      j.LowLevelSPSymbol,
	}

	for _, ms := range j.MapSymbol {
		cfg.MapSymbol[ms.MainSymbol] = ms.CollectedSymbols
	}

	return cfg
}

func parseCollectorPayTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCollectorPayTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCollectorPayTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCollectorPayTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCollectorPayTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseDropDownSymbols2:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CollectorPayTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
