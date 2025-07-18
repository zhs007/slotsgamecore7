package lowcode

import (
	"log/slog"
	"sync"

	"github.com/fatih/color"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
)

type GamePropertyPool struct {
	MapGamePropPool     map[int]*sync.Pool
	Config              *Config
	DefaultPaytables    *sgc7game.PayTables
	DefaultLineData     *sgc7game.LineData
	SymbolsViewer       *SymbolsViewer
	MapSymbolColor      *asciigame.SymbolColorMap
	mapComponents       map[int]*ComponentList
	mapStrValWeights    map[string]*sgc7game.ValWeights2
	mapIntValWeights    map[string]*sgc7game.ValWeights2
	mapSymbolValWeights map[string]*sgc7game.ValWeights2
	newRNG              FuncNewRNG
	newFeatureLevel     FuncNewFeatureLevel
}

func (pool *GamePropertyPool) newGameProp(betMul int) *GameProperty {
	gameProp := &GameProperty{
		CurBetMul:        betMul,
		Pool:             pool,
		MapVals:          make(map[int]int),
		MapStrVals:       make(map[int]string),
		MapIntValWeights: make(map[string]*sgc7game.ValWeights2),
		MapStats:         make(map[string]*sgc7stats.Feature),
		mapInt:           make(map[string]int),
		CurPaytables:     pool.DefaultPaytables,
		CurLineData:      pool.DefaultLineData,
		PoolScene:        sgc7game.NewGameScenePoolEx(),
		SceneStack:       NewSceneStack(false),
		OtherSceneStack:  NewSceneStack(true),
		callStack:        NewCallStack(),
		rng:              pool.newRNG(),
		featureLevel:     pool.newFeatureLevel(betMul),
	}

	if gameProp.CurLineData != nil {
		gameProp.SetVal(GamePropCurLineNum, len(gameProp.CurLineData.Lines))
	}

	// mapc, isok := pool.mapComponents[betMul]
	// if !isok {
	// 	goutils.Error("GamePropertyPool.newGameProp:mapComponents",
	// 		slog.Int("betMul", betMul),
	// 		goutils.Err(ErrInvalidBet))

	// 	return nil
	// }

	// for k, v := range mapc.MapComponents {
	// 	gameProp.MapComponentData[k] = v.NewComponentData()
	// }

	gameProp.SetVal(GamePropWidth, pool.Config.Width)
	gameProp.SetVal(GamePropHeight, pool.Config.Height)

	return gameProp
}

// func (pool *GamePropertyPool) NewGameProp(betMul int) (*GameProperty, error) {
// 	gameProp := pool.newGameProp()

// 	gameProp.SetVal(GamePropWidth, pool.Config.Width)
// 	gameProp.SetVal(GamePropHeight, pool.Config.Height)

// 	return gameProp, nil
// }

func (pool *GamePropertyPool) onAddComponentList(betMul int, components *ComponentList) {
	pool.mapComponents[betMul] = components
}

// func (pool *GamePropertyPool) NewStatsWithConfig(betMul int, parent *sgc7stats.Feature, cfg *StatsConfig) (*sgc7stats.Feature, error) {
// 	components, isok := pool.mapComponents[betMul]
// 	if !isok {
// 		goutils.Error("GameProperty.NewStatsWithConfig",
// 			slog.Int("bet", betMul),
// 			goutils.Err(ErrInvalidBet))

// 		return nil, ErrInvalidBet
// 	}

// 	curComponent, isok := components.MapComponents[cfg.Component]
// 	if !isok {
// 		goutils.Error("GameProperty.NewStatsWithConfig",
// 			goutils.Err(ErrInvalidStatsComponentInConfig))

// 		return nil, ErrInvalidStatsComponentInConfig
// 	}

// 	feature := NewStatsFeature(parent, cfg.Name, func(f *sgc7stats.Feature, s *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 		if cfg.IsNeedForceStats {
// 			return true, s.CashBet, calcTotalCashWins(lst)
// 		}

// 		return curComponent.OnStats(f, s, lst)
// 	}, pool.Config.Width, pool.Config.StatsSymbolCodes, StatusTypeUnknow, "")

// 	for _, v := range cfg.Children {
// 		_, err := pool.NewStatsWithConfig(betMul, feature, v)
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:NewStatsWithConfig",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	for k, v := range cfg.RespinEndingStatus {
// 		_, err := pool.newStatusStats(betMul, feature, v, StatusTypeRespinEnding, k)
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	for k, v := range cfg.RespinStartStatus {
// 		_, err := pool.newStatusStats(betMul, feature, v, StatusTypeRespinStart, k)
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	for k, v := range cfg.RespinStartStatusEx {
// 		_, err := pool.newStatusStats(betMul, feature, v, StatusTypeRespinStartEx, k)
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	for _, v := range cfg.RespinNumStatus {
// 		_, err := pool.newStatusStats(betMul, feature, v, StatusTypeRespinNum, "")
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	for _, v := range cfg.RespinWinStatus {
// 		_, err := pool.newStatusStats(betMul, feature, v, StatusTypeRespinWin, "")
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	for _, v := range cfg.RespinStartNumStatus {
// 		_, err := pool.newStatusStats(betMul, feature, v, StatusTypeRespinStartNum, "")
// 		if err != nil {
// 			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
// 				slog.Any("v", v),
// 				goutils.Err(err))

// 			return nil, err
// 		}
// 	}

// 	return feature, nil
// }

// func (pool *GamePropertyPool) newStatusStats(betMul int, parent *sgc7stats.Feature, componentName string, statusType int, respinName string) (*sgc7stats.Feature, error) {
// 	components, isok := pool.mapComponents[betMul]
// 	if !isok {
// 		goutils.Error("GameProperty.newStatusStats",
// 			slog.Int("bet", betMul),
// 			goutils.Err(ErrInvalidBet))

// 		return nil, ErrInvalidBet
// 	}

// 	curComponent, isok := components.MapComponents[componentName]
// 	if !isok {
// 		goutils.Error("GameProperty.NewStatsWithConfig",
// 			goutils.Err(ErrInvalidStatsComponentInConfig))

// 		return nil, ErrInvalidStatsComponentInConfig
// 	}

// 	feature := NewStatsFeature(parent, componentName, func(f *sgc7stats.Feature, s *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 		return curComponent.OnStats(f, s, lst)
// 	}, pool.Config.Width, pool.Config.StatsSymbolCodes, statusType, respinName)

// 	return feature, nil
// }

func (pool *GamePropertyPool) loadAllWeights() {
	for v, vw2 := range pool.Config.mapValWeights {
		pool.mapIntValWeights[v] = vw2
	}

	for v, vw2 := range pool.Config.mapStrWeights {
		pool.mapStrValWeights[v] = vw2
	}

	for v, vw2 := range pool.Config.mapReelSetWeights {
		pool.mapStrValWeights[v] = vw2
	}
}

func (pool *GamePropertyPool) onInit() {
	for bet, v := range pool.mapComponents {
		v.onInit(pool.Config.MapBetConfigs[bet].Start)
	}
}

func (pool *GamePropertyPool) InitStats(betMul int) error {
	err := pool.Config.BuildStatsSymbolCodes(pool.DefaultPaytables)
	if err != nil {
		goutils.Error("GamePropertyPool.InitStats:BuildStatsSymbolCodes",
			goutils.Err(err))

		return err
	}

	// if !gIsForceDisableStats && pool.Config.Stats != nil {
	// 	statsTotal := sgc7stats.NewFeature("total", sgc7stats.FeatureBasic, func(f *sgc7stats.Feature, s *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	// 		totalWin := int64(0)

	// 		for _, v := range lst {
	// 			totalWin += v.CashWin
	// 		}

	// 		return true, s.CashBet, totalWin
	// 	}, nil)

	// 	_, err := pool.NewStatsWithConfig(betMul, statsTotal, pool.Config.Stats)
	// 	if err != nil {
	// 		goutils.Error("GameProperty.InitStats:BuildStatsSymbolCodes",
	// 			goutils.Err(err))

	// 		return err
	// 	}

	// 	pool.Stats = NewStats(statsTotal, pool)
	// 	pool.Stats.Bet = betMul

	// 	go pool.Stats.StartWorker()
	// }

	return nil
}

// LoadStrWeights - load xlsx file
func (pool *GamePropertyPool) LoadStrWeights(fn string, useFileMapping bool) (*sgc7game.ValWeights2, error) {
	// if gJsonMode {
	return pool.mapStrValWeights[fn], nil
	// }

	// pool.lock.RLock()
	// vw, isok := pool.mapStrValWeights[fn]
	// if isok {
	// 	pool.lock.RUnlock()

	// 	return vw, nil
	// }
	// pool.lock.RUnlock()

	// if pool.Config.mapValWeights != nil {
	// 	nvw := pool.Config.mapValWeights[fn]

	// 	pool.lock.Lock()
	// 	pool.mapStrValWeights[fn] = nvw
	// 	pool.lock.Unlock()

	// 	return nvw, nil
	// }

	// vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(fn, useFileMapping), "val", "weight", sgc7game.NewStrVal)
	// if err != nil {
	// 	goutils.Error("GamePropertyPool.LoadStrWeights:LoadValWeights2FromExcel",
	// 		slog.String("fn", fn),
	// 		goutils.Err(err))

	// 	return nil, err
	// }

	// pool.lock.Lock()
	// pool.mapStrValWeights[fn] = vw2
	// pool.lock.Unlock()

	// return vw2, nil
}

// LoadIntWeights - load xlsx file
func (pool *GamePropertyPool) LoadIntWeights(fn string, useFileMapping bool) (*sgc7game.ValWeights2, error) {
	// if gJsonMode {
	return pool.mapIntValWeights[fn], nil
	// }

	// pool.lock.RLock()
	// vw, isok := pool.mapIntValWeights[fn]
	// if isok {
	// 	pool.lock.RUnlock()

	// 	return vw, nil
	// }
	// pool.lock.RUnlock()

	// if pool.Config.mapValWeights != nil {
	// 	vw := pool.Config.mapValWeights[fn]

	// 	vals := make([]sgc7game.IVal, len(vw.Vals))

	// 	for _, v := range vw.Vals {
	// 		i64, err := goutils.String2Int64(v.String())
	// 		if err != nil {
	// 			goutils.Error("GamePropertyPool.LoadIntWeights:String2Int64",
	// 				goutils.Err(err))

	// 			return nil, err
	// 		}

	// 		vals = append(vals, sgc7game.NewIntValEx[int](int(i64)))
	// 	}

	// 	nvw, err := sgc7game.NewValWeights2(vals, vw.Weights)
	// 	if err != nil {
	// 		goutils.Error("GamePropertyPool.LoadIntWeights:NewValWeights2",
	// 			goutils.Err(err))

	// 		return nil, err
	// 	}

	// 	pool.lock.Lock()
	// 	pool.mapIntValWeights[fn] = nvw
	// 	pool.lock.Unlock()

	// 	return nvw, nil
	// }

	// vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(fn, useFileMapping), "val", "weight", sgc7game.NewIntVal[int])
	// if err != nil {
	// 	goutils.Error("GamePropertyPool.LoadIntWeights:LoadValWeights2FromExcel",
	// 		slog.String("fn", fn),
	// 		goutils.Err(err))

	// 	return nil, err
	// }

	// pool.lock.Lock()
	// pool.mapIntValWeights[fn] = vw2
	// pool.lock.Unlock()

	// return vw2, nil
}

// LoadSymbolWeights - load xlsx file
func (pool *GamePropertyPool) LoadIntMapping(fn string) *sgc7game.ValMapping2 {
	return pool.Config.mapIntMapping[fn]
}

// LoadSymbolWeights - load xlsx file
func (pool *GamePropertyPool) LoadSymbolWeights(fn string, headerVal string, headerWeight string, paytables *sgc7game.PayTables, useFileMapping bool) (*sgc7game.ValWeights2, error) {
	return pool.mapIntValWeights[fn], nil
}

func (pool *GamePropertyPool) SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, name string, index int, mask bool) error {
	ic, isok := gameProp.Components.MapComponents[name]
	if !isok || !ic.IsMask() {
		goutils.Error("GamePropertyPool.SetMaskVal",
			slog.String("name", name),
			goutils.Err(ErrInvalidComponentName))

		return ErrInvalidComponentName
	}

	cd := gameProp.GetComponentData(ic)
	// im, isok := ic.(IMask)
	// if !isok {
	// 	goutils.Error("GamePropertyPool.SetMaskVal",
	// 		slog.String("name", name),
	// 		goutils.Err(ErrNotMask))

	// 	return ErrNotMask
	// }

	if cd != nil {
		ic.SetMaskVal(plugin, gameProp, curpr, gp, cd, index, mask)

		return nil
	}

	goutils.Error("GamePropertyPool.GetMask",
		slog.String("name", name),
		goutils.Err(ErrInvalidComponent))

	return ErrInvalidComponent
}

func (pool *GamePropertyPool) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, name string, mask []bool, isOnlyTrue bool) error {
	ic, isok := gameProp.Components.MapComponents[name]
	if !isok || !ic.IsMask() {
		goutils.Error("GamePropertyPool.SetMask",
			slog.String("name", name),
			goutils.Err(ErrInvalidComponentName))

		return ErrInvalidComponentName
	}

	cd := gameProp.GetComponentData(ic)
	// im, isok := ic.(IMask)
	// if !isok {
	// 	goutils.Error("GamePropertyPool.SetMaskVal",
	// 		slog.String("name", name),
	// 		goutils.Err(ErrNotMask))

	// 	return ErrNotMask
	// }

	if cd != nil {
		if isOnlyTrue {
			ic.SetMaskOnlyTrue(plugin, gameProp, curpr, gp, cd, mask)
		} else {
			ic.SetMask(plugin, gameProp, curpr, gp, cd, mask)
		}

		return nil
	}

	goutils.Error("GamePropertyPool.GetMask",
		slog.String("name", name),
		goutils.Err(ErrInvalidComponent))

	return ErrInvalidComponent

	// im, isok := ic.(IMask)
	// if !isok {
	// 	goutils.Error("GamePropertyPool.SetMask",
	// 		slog.String("name", name),
	// 		goutils.Err(ErrNotMask))

	// 	return ErrNotMask
	// }

	// if isOnlyTrue {
	// 	im.SetMaskOnlyTrue(plugin, gameProp, curpr, gp, mask)
	// } else {
	// 	im.SetMask(plugin, gameProp, curpr, gp, mask)
	// }

	// return nil
}

func (pool *GamePropertyPool) GetMask(name string, gameProp *GameProperty) ([]bool, error) {
	ic, isok := gameProp.Components.MapComponents[name]
	if !isok || !ic.IsMask() {
		goutils.Error("GamePropertyPool.GetMask",
			slog.String("name", name),
			goutils.Err(ErrInvalidComponentName))

		return nil, ErrInvalidComponentName
	}

	cd := gameProp.GetComponentData(ic)
	if cd != nil {
		return cd.GetMask(), nil
	}

	goutils.Error("GamePropertyPool.GetMask",
		slog.String("name", name),
		goutils.Err(ErrInvalidComponent))

	// im, isok := ic.(IMask)
	// if !isok {
	// 	goutils.Error("GamePropertyPool.GetMask",
	// 		slog.String("name", name),
	// 		goutils.Err(ErrNotMask))

	// 	return nil, ErrNotMask
	// }

	// mask := im.GetMask(gameProp)

	return nil, ErrInvalidComponent
}

func (pool *GamePropertyPool) PushTrigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, name string, num int) error {
	ic, isok := gameProp.Components.MapComponents[name]
	if !isok || !ic.IsRespin() {
		goutils.Error("GamePropertyPool.PushTrigger",
			slog.String("name", name),
			goutils.Err(ErrInvalidComponentName))

		return ErrInvalidComponentName
	}

	pool.pushTrigger(gameProp, plugin, curpr, gp, ic, num)

	return nil
}

func (pool *GamePropertyPool) pushTrigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, ic IComponent, num int) error {
	cd := gameProp.GetGlobalComponentData(ic)
	cd.PushTriggerRespin(gameProp, plugin, curpr, gp, num)

	return nil
}

func (pool *GamePropertyPool) GetComponentList(bet int) *ComponentList {
	return pool.mapComponents[bet]
}

func (pool *GamePropertyPool) NewPlayerState() (*PlayerState, error) {
	ps := NewPlayerState()

	return ps, nil
}

func (pool *GamePropertyPool) InitPlayerState() (*PlayerState, error) {
	ps := NewPlayerState()

	for betMethod, components := range pool.mapComponents {
		for _, c := range components.Components {
			err := c.InitPlayerState(pool, nil, nil, ps, betMethod, 0)
			if err != nil {
				goutils.Error("GamePropertyPool.InitPlayerState:InitPlayerState",
					goutils.Err(err))

				return nil, err
			}
		}
	}

	return ps, nil
}

func (pool *GamePropertyPool) InitPlayerStateOnBet(gameProp *GameProperty, plugin sgc7plugin.IPlugin, ps *PlayerState,
	stake *sgc7game.Stake) error {

	betMethod := stake.CashBet / stake.CoinBet
	components, isok := pool.mapComponents[int(betMethod)]
	if isok {
		for _, c := range components.Components {
			err := c.InitPlayerState(pool, gameProp, plugin, ps, int(betMethod), int(stake.CoinBet))
			if err != nil {
				goutils.Error("GamePropertyPool.InitPlayerStateOnBet:InitPlayerState",
					goutils.Err(err))

				return err
			}
		}
	}

	return nil
}

func (pool *GamePropertyPool) ChgReelsCollector(gameProp *GameProperty, name string, ps *PlayerState, betMethod int, bet int, reelsData []int) error {
	ic, isok := gameProp.Components.MapComponents[name]
	if !isok {
		goutils.Error("GamePropertyPool.ChgReelsCollector",
			slog.String("name", name),
			goutils.Err(ErrInvalidComponentName))

		return ErrInvalidComponentName
	}

	icd := gameProp.GetComponentData(ic)

	ic.ChgReelsCollector(icd, ps, betMethod, bet, reelsData)

	return nil
}

func newGamePropertyPool2(cfg *Config, funcNewRNG FuncNewRNG, funcNewFeatureLevel FuncNewFeatureLevel) (*GamePropertyPool, error) {
	pool := &GamePropertyPool{
		MapGamePropPool:     make(map[int]*sync.Pool),
		Config:              cfg,
		DefaultPaytables:    cfg.GetDefaultPaytables(),
		DefaultLineData:     cfg.GetDefaultLineData(),
		mapComponents:       make(map[int]*ComponentList),
		mapStrValWeights:    make(map[string]*sgc7game.ValWeights2),
		mapIntValWeights:    make(map[string]*sgc7game.ValWeights2),
		mapSymbolValWeights: make(map[string]*sgc7game.ValWeights2),
		newRNG:              funcNewRNG,
		newFeatureLevel:     funcNewFeatureLevel,
	}

	if cfg.SymbolsViewer == "" {
		sv := NewSymbolViewerFromPaytables(pool.DefaultPaytables)

		pool.SymbolsViewer = sv
	} else {
		sv, err := LoadSymbolsViewer(cfg.GetPath(cfg.SymbolsViewer, false))
		if err != nil {
			goutils.Error("NewGamePropertyPool2:LoadSymbolsViewer",
				slog.String("fn", cfg.SymbolsViewer),
				goutils.Err(err))

			return nil, err
		}

		pool.SymbolsViewer = sv
	}

	pool.MapSymbolColor = asciigame.NewSymbolColorMap(pool.DefaultPaytables)
	wColor := color.New(color.BgRed, color.FgHiWhite)
	hColor := color.New(color.BgBlue, color.FgHiWhite)
	mColor := color.New(color.BgGreen, color.FgHiWhite)
	sColor := color.New(color.BgMagenta, color.FgHiWhite)
	for k, v := range pool.SymbolsViewer.MapSymbols {
		switch v.Color {
		case "wild":
			pool.MapSymbolColor.AddSymbolColor(k, wColor)
		case "high":
			pool.MapSymbolColor.AddSymbolColor(k, hColor)
		case "medium":
			pool.MapSymbolColor.AddSymbolColor(k, mColor)
		case "scatter":
			pool.MapSymbolColor.AddSymbolColor(k, sColor)
		}
	}

	pool.MapSymbolColor.OnGetSymbolString = func(s int) string {
		obj, isok := pool.SymbolsViewer.MapSymbols[s]
		if isok {
			return obj.Output
		}

		return " "
	}

	for _, bet := range cfg.Bets {
		pool.MapGamePropPool[bet] = &sync.Pool{
			New: func() any {
				return pool.newGameProp(bet)
			},
		}
	}

	return pool, nil
}
