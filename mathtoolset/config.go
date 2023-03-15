package mathtoolset

import (
	"os"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Type      string  `yaml:"type"`
	Code      string  `yaml:"code"`
	TargetRTP float64 `yaml:"targetRTP"`
}

func LoadConfig(fn string) (*Config, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("LoadConfig:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("LoadConfig:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return cfg, nil
}
