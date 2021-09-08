package sgc7game

import (
	"io/ioutil"

	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// BasicGameConfig - configuration for basic game
type BasicGameConfig struct {
	LineData  string            `yaml:"linedata"`
	PayTables string            `yaml:"paytables"`
	Reels     map[string]string `yaml:"reels"`
}

// LoadGameConfig - load configuration
func LoadGameConfig(fn string, cfg interface{}) error {
	data, err := ioutil.ReadFile(fn)
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
			zap.Error(ErrNullConfig))

		return ErrNullConfig
	}

	err := cfg.LoadLine5(bgc.LineData)
	if err != nil {
		goutils.Error("BasicGameConfig.Init5:LoadLine5",
			zap.String("LineData", bgc.LineData),
			zap.Error(err))

		return err
	}

	err = cfg.LoadPayTables5(bgc.PayTables)
	if err != nil {
		goutils.Error("BasicGameConfig.Init5:LoadPayTables5",
			zap.String("PayTables", bgc.PayTables),
			zap.Error(err))

		return err
	}

	for k, v := range bgc.Reels {
		err = cfg.LoadReels5(k, v)
		if err != nil {
			goutils.Error("BasicGameConfig.Init5:LoadReels5",
				zap.String("reels", v),
				zap.Error(err))

			return err
		}
	}

	return nil
}
