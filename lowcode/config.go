package lowcode

import (
	"log/slog"
	"path"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
)

type ComponentConfig struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Config string `yaml:"config"`
}

type BetConfig struct {
	Bet            int                              `yaml:"bet"`
	TotalBetInWins int                              `yaml:"totalBetInWins"`
	Start          string                           `yaml:"start"`
	Components     []*ComponentConfig               `yaml:"components"`
	mapConfig      map[string]IComponentConfig      `yaml:"-"`
	mapBasicConfig map[string]*BasicComponentConfig `yaml:"-"`
	ForceEndings   []string                         `yaml:"-"`
}

func (betCfg *BetConfig) Reset(start string, endings []string) {
	betCfg.Start = start
	betCfg.ForceEndings = endings
}

type Config struct {
	Name              string                           `yaml:"name"`
	Width             int                              `yaml:"width"`
	Height            int                              `yaml:"height"`
	Linedata          map[string]string                `yaml:"linedata"`
	MapLinedate       map[string]*sgc7game.LineData    `yaml:"-"`
	Paytables         map[string]string                `yaml:"paytables"`
	MapPaytables      map[string]*sgc7game.PayTables   `yaml:"-"`
	Reels             map[string]string                `yaml:"reels"`
	MapReels          map[string]*sgc7game.ReelsData   `yaml:"-"`
	FileMapping       map[string]string                `yaml:"fileMapping"`
	SymbolsViewer     string                           `yaml:"symbolsViewer"`
	DefaultScene      string                           `yaml:"defaultScene"`
	DefaultPaytables  string                           `yaml:"defaultPaytables"`
	DefaultLinedata   string                           `yaml:"defaultLinedata"`
	Bets              []int                            `yaml:"bets"`
	TotalBetInWins    []int                            `yaml:"totalBetInWins"`
	StatsSymbols      []string                         `yaml:"statsSymbols"`
	StatsSymbolCodes  []mathtoolset.SymbolType         `yaml:"-"`
	MainPath          string                           `yaml:"mainPath"`
	MapCmdComponent   map[string]string                `yaml:"mapCmdComponent"`
	ComponentsMapping map[int]map[string]string        `yaml:"componentsMapping"`
	MapBetConfigs     map[int]*BetConfig               `yaml:"mapBetConfigs"`
	mapValWeights     map[string]*sgc7game.ValWeights2 `yaml:"-"`
	mapReelSetWeights map[string]*sgc7game.ValWeights2 `yaml:"-"`
	mapStrWeights     map[string]*sgc7game.ValWeights2 `yaml:"-"`
	mapIntMapping     map[string]*sgc7game.ValMapping2 `yaml:"-"`
}

func (cfg *Config) Reset(bet int, start string, endings []string) {
	betCfg, isok := cfg.MapBetConfigs[bet]
	if isok {
		betCfg.Reset(start, endings)
	}
}

func (cfg *Config) GetPath(fn string, useFileMapping bool) string {
	if useFileMapping {
		curfn, isok := cfg.FileMapping[fn]
		if isok {
			fn = curfn
		}
	}

	if cfg.MainPath != "" {
		return path.Join(cfg.MainPath, fn)
	}

	return fn
}

func (cfg *Config) BuildStatsSymbolCodes(paytables *sgc7game.PayTables) error {
	cfg.StatsSymbolCodes = nil
	for _, v := range cfg.StatsSymbols {
		symbolCode, isok := paytables.MapSymbols[v]
		if !isok {
			goutils.Error("Config.BuildStatsSymbolCodes",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidStatsSymbolsInConfig))

			return ErrInvalidStatsSymbolsInConfig
		}

		cfg.StatsSymbolCodes = append(cfg.StatsSymbolCodes, mathtoolset.SymbolType(symbolCode))
	}

	return nil
}

func (cfg *Config) GetDefaultPaytables() *sgc7game.PayTables {
	name := cfg.DefaultPaytables
	if name == "" {
		name = "main"
	}

	pt, isok := cfg.MapPaytables[name]
	if isok {
		return pt
	}

	return nil
}

func (cfg *Config) GetDefaultLineData() *sgc7game.LineData {
	name := cfg.DefaultLinedata
	if name == "" {
		name = "main"
	}

	ld, isok := cfg.MapLinedate[name]
	if isok {
		return ld
	}

	return nil
}
