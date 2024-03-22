package mathtoolset

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"gopkg.in/yaml.v2"
)

type GenReelsConfig struct {
	ReelsStatsFilename string   `yaml:"reelsStatsFilename"`
	ReelsFilename      string   `yaml:"reelsFilename"`
	MainSymbols        []string `yaml:"mainSymbols"`
	Offset             int      `yaml:"offset"`
}

type CodeConfig struct {
	Name           string `yaml:"name"`
	Code           string `yaml:"code"`
	DisableAutoRun bool   `yaml:"disableAutoRun"`
}

type Config struct {
	Type           string          `yaml:"type"`
	Code           string          `yaml:"code"`
	Codes          []*CodeConfig   `yaml:"codes"`
	TargetRTP      float64         `yaml:"targetRTP"`
	Paytables      string          `yaml:"paytables"`
	GenReelsConfig *GenReelsConfig `yaml:"genReelsConfig"`
}

func LoadConfig(fn string) (*Config, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("LoadConfig:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("LoadConfig:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	return cfg, nil
}
