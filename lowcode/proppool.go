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
		return curComponent.OnStats(f, s, lst)
	}, pool.Config.Width, pool.Config.StatsSymbolCodes)

	for _, v := range cfg.Children {
		_, err := pool.NewStatsWithConfig(feature, v)
		if err != nil {
			goutils.Error("GameProperty.NewStatsWithConfig:NewStatsWithConfig",
				goutils.JSON("v", v),
				zap.Error(err))

			return nil, err
		}
	}

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

		pool.Stats = NewStats(statsTotal)

		go pool.Stats.StartWorker()
	}

	return nil
}

func NewGamePropertyPool(cfgfn string) (*GamePropertyPool, error) {
	cfg, err := LoadConfig(cfgfn)
	if err != nil {
		goutils.Error("NewGamePropertyPool:LoadConfig",
			zap.String("cfgfn", cfgfn),
			zap.Error(err))

		return nil, err
	}

	pool := &GamePropertyPool{
		Config:           cfg,
		DefaultPaytables: cfg.GetDefaultPaytables(),
		DefaultLineData:  cfg.GetDefaultLineData(),
		MapComponents:    make(map[string]IComponent),
	}

	sv, err := LoadSymbolsViewer(cfg.SymbolsViewer)
	if err != nil {
		goutils.Error("NewGamePropertyPool:LoadSymbolsViewer",
			zap.String("fn", cfg.SymbolsViewer),
			zap.Error(err))

		return nil, err
	}

	pool.SymbolsViewer = sv
	pool.MapSymbolColor = asciigame.NewSymbolColorMap(pool.DefaultPaytables)
	wColor := color.New(color.BgRed, color.FgHiWhite)
	hColor := color.New(color.BgBlue, color.FgHiWhite)
	mColor := color.New(color.BgGreen, color.FgHiWhite)
	sColor := color.New(color.BgMagenta, color.FgHiWhite)
	for k, v := range sv.MapSymbols {
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
