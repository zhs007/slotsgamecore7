package sgc7game

import (
	"log/slog"
	"os"

	goutils "github.com/zhs007/goutils"
	"gopkg.in/yaml.v2"
)

// BasicGameConfig - configuration for basic game
type BasicGameConfig struct {
	LineData  string            `yaml:"linedata"`
	PayTables string            `yaml:"paytables"`
	Reels     map[string]string `yaml:"reels"`
}

// LoadGameConfig - load configuration
func LoadGameConfig(fn string, cfg any) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}

// Init5 - initial with 5 reels
func (bgc BasicGameConfig) Init5(ig IGame) error {
	cfg := ig.GetConfig()
	if cfg == nil {
		goutils.Error("BasicGameConfig.Init5:GetConfig",
			goutils.Err(ErrNullConfig))

		return ErrNullConfig
	}

	err := cfg.LoadLine5(bgc.LineData)
	if err != nil {
		goutils.Error("BasicGameConfig.Init5:LoadLine5",
			slog.String("LineData", bgc.LineData),
			goutils.Err(err))

		return err
	}

	err = cfg.LoadPayTables5(bgc.PayTables)
	if err != nil {
		goutils.Error("BasicGameConfig.Init5:LoadPayTables5",
			slog.String("PayTables", bgc.PayTables),
			goutils.Err(err))

		return err
	}

	for k, v := range bgc.Reels {
		err = cfg.LoadReels5(k, v)
		if err != nil {
			goutils.Error("BasicGameConfig.Init5:LoadReels5",
				slog.String("reels", v),
				goutils.Err(err))

			return err
		}
	}

	return nil
}
