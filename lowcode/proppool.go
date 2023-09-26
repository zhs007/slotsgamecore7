package lowcode

import (
	"sync"

	"github.com/fatih/color"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
)

type GamePropertyPool struct {
	Pool             sync.Pool
	Config           *Config
	DefaultPaytables *sgc7game.PayTables
	DefaultLineData  *sgc7game.LineData
	SymbolsViewer    *SymbolsViewer
	MapSymbolColor   *asciigame.SymbolColorMap
	MapComponents    map[string]IComponent
	Stats            *Stats
}

func (pool *GamePropertyPool) newGameProp() *GameProperty {
	gameProp := &GameProperty{
		Pool:             pool,
		MapVals:          make(map[int]int),
		MapStrVals:       make(map[int]string),
		MapIntValWeights: make(map[string]*sgc7game.ValWeights2),
		MapStats:         make(map[string]*sgc7stats.Feature),
		mapInt:           make(map[string]int),
		CurPaytables:     pool.DefaultPaytables,
		CurLineData:      pool.DefaultLineData,
		MapComponentData: make(map[string]IComponentData),
		PoolScene:        sgc7game.NewGameScenePoolEx(),
	}

	if gameProp.CurLineData != nil {
		gameProp.SetVal(GamePropCurLineNum, len(gameProp.CurLineData.Lines))
	}

	for k, v := range pool.MapComponents {
		gameProp.MapComponentData[k] = v.NewComponentData()
	}

	return gameProp
}

func (pool *GamePropertyPool) NewGameProp() (*GameProperty, error) {
	gameProp := pool.newGameProp()

	gameProp.SetVal(GamePropWidth, pool.Config.Width)
	gameProp.SetVal(GamePropHeight, pool.Config.Height)

	return gameProp, nil
}

func (pool *GamePropertyPool) onAddComponent(name string, component IComponent) {
	pool.MapComponents[name] = component
}

func (pool *GamePropertyPool) NewStatsWithConfig(parent *sgc7stats.Feature, cfg *StatsConfig) (*sgc7stats.Feature, error) {
	curComponent, isok := pool.MapComponents[cfg.Component]
	if !isok {
		goutils.Error("GameProperty.NewStatsWithConfig",
			zap.Error(ErrIvalidStatsComponentInConfig))

		return nil, ErrIvalidStatsComponentInConfig
	}

	feature := NewStatsFeature(parent, cfg.Name, func(f *sgc7stats.Feature, s *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
		if cfg.IsNeedForceStats {
			return true, s.CashBet, calcTotalCashWins(lst)
		}

		return curComponent.OnStats(f, s, lst)
	}, pool.Config.Width, pool.Config.StatsSymbolCodes, StatusTypeUnknow, "")

	for _, v := range cfg.Children {
		_, err := pool.NewStatsWithConfig(feature, v)
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:NewStatsWithConfig",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	for k, v := range cfg.RespinEndingStatus {
		_, err := pool.newStatusStats(feature, v, StatusTypeRespinEnding, k)
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	for k, v := range cfg.RespinStartStatus {
		_, err := pool.newStatusStats(feature, v, StatusTypeRespinStart, k)
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	for k, v := range cfg.RespinStartStatusEx {
		_, err := pool.newStatusStats(feature, v, StatusTypeRespinStartEx, k)
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	for _, v := range cfg.RespinNumStatus {
		_, err := pool.newStatusStats(feature, v, StatusTypeRespinNum, "")
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	for _, v := range cfg.RespinWinStatus {
		_, err := pool.newStatusStats(feature, v, StatusTypeRespinWin, "")
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	for _, v := range cfg.RespinStartNumStatus {
		_, err := pool.newStatusStats(feature, v, StatusTypeRespinStartNum, "")
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:newStatusStats",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

	return feature, nil
}

func (pool *GamePropertyPool) newStatusStats(parent *sgc7stats.Feature, componentName string, statusType int, respinName string) (*sgc7stats.Feature, error) {
	curComponent, isok := pool.MapComponents[componentName]
	if !isok {
		goutils.Error("GameProperty.NewStatsWithConfig",
			zap.Error(ErrIvalidStatsComponentInConfig))

		return nil, ErrIvalidStatsComponentInConfig
	}

	feature := NewStatsFeature(parent, componentName, func(f *sgc7stats.Feature, s *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
		return curComponent.OnStats(f, s, lst)
	}, pool.Config.Width, pool.Config.StatsSymbolCodes, statusType, respinName)

	return feature, nil
}

func (pool *GamePropertyPool) InitStats() error {
	err := pool.Config.BuildStatsSymbolCodes(pool.DefaultPaytables)
	if err != nil {
		goutils.Error("GamePropertyPool.InitStats:BuildStatsSymbolCodes",
			zap.Error(err))

		return err
	}

	if !gIsForceDisableStats && pool.Config.Stats != nil {
		statsTotal := sgc7stats.NewFeature("total", sgc7stats.FeatureBasic, func(f *sgc7stats.Feature, s *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
			totalWin := int64(0)

			for _, v := range lst {
				totalWin += v.CashWin
			}

			return true, s.CashBet, totalWin
		}, nil)

		_, err := pool.NewStatsWithConfig(statsTotal, pool.Config.Stats)
		if err != nil {
			goutils.Error("GameProperty.InitStats:BuildStatsSymbolCodes",
				zap.Error(err))

			return err
		}

		pool.Stats = NewStats(statsTotal, pool)

		go pool.Stats.StartWorker()
	}

	return nil
}

// LoadValWeights - load xlsx file
func (pool *GamePropertyPool) LoadValWeights(fn string, headerVal string, headerWeight string, funcNew sgc7game.FuncNewIVal, useFileMapping bool) (*sgc7game.ValWeights2, error) {
	if pool.Config.mapValWeights != nil {
		return pool.Config.mapValWeights[fn], nil
	}

	vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(fn, useFileMapping), headerVal, headerWeight, funcNew)
	if err != nil {
		goutils.Error("GamePropertyPool.LoadValWeights:LoadValWeights2FromExcel",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return vw2, nil
}

// LoadSymbolWeights - load xlsx file
func (pool *GamePropertyPool) LoadSymbolWeights(fn string, headerVal string, headerWeight string, paytables *sgc7game.PayTables, useFileMapping bool) (*sgc7game.ValWeights2, error) {
	if pool.Config.mapValWeights != nil {
		vw := pool.Config.mapValWeights[fn]

		vals := make([]sgc7game.IVal, len(vw.Vals))

		for i, v := range vw.Vals {
			vals[i] = sgc7game.NewIntValEx(paytables.MapSymbols[v.String()])
		}

		nvw, err := sgc7game.NewValWeights2(vals, vw.Weights)
		if err != nil {
			goutils.Error("GamePropertyPool.LoadValWeights:NewValWeights2",
				zap.Error(err))

			return nil, err
		}

		return nvw, nil
	}

	vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(pool.Config.GetPath(fn, useFileMapping), headerVal, headerWeight, paytables)
	if err != nil {
		goutils.Error("GamePropertyPool.LoadValWeights:LoadValWeights2FromExcel",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return vw2, nil
}

func NewGamePropertyPool(cfgfn string) (*GamePropertyPool, error) {
	cfg, err := LoadConfig(cfgfn)
	if err != nil {
		goutils.Error("NewGamePropertyPool:LoadConfig",
			zap.String("cfgfn", cfgfn),
			zap.Error(err))

		return nil, err
	}

	return NewGamePropertyPool2(cfg)
}

func NewGamePropertyPool2(cfg *Config) (*GamePropertyPool, error) {
	pool := &GamePropertyPool{
		Config:           cfg,
		DefaultPaytables: cfg.GetDefaultPaytables(),
		DefaultLineData:  cfg.GetDefaultLineData(),
		MapComponents:    make(map[string]IComponent),
	}

	if cfg.SymbolsViewer == "" {
		sv := NewSymbolViewerFromPaytables(pool.DefaultPaytables)

		pool.SymbolsViewer = sv
	} else {
		sv, err := LoadSymbolsViewer(cfg.GetPath(cfg.SymbolsViewer, false))
		if err != nil {
			goutils.Error("NewGamePropertyPool2:LoadSymbolsViewer",
				zap.String("fn", cfg.SymbolsViewer),
				zap.Error(err))

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
		if v.Color == "wild" {
			pool.MapSymbolColor.AddSymbolColor(k, wColor)
		} else if v.Color == "high" {
			pool.MapSymbolColor.AddSymbolColor(k, hColor)
		} else if v.Color == "medium" {
			pool.MapSymbolColor.AddSymbolColor(k, mColor)
		} else if v.Color == "scatter" {
			pool.MapSymbolColor.AddSymbolColor(k, sColor)
		}
	}

	pool.MapSymbolColor.OnGetSymbolString = func(s int) string {
		return pool.SymbolsViewer.MapSymbols[s].Output
	}

	pool.Pool = sync.Pool{
		New: func() any {
			return pool.newGameProp()
		},
	}

	return pool, nil
}
