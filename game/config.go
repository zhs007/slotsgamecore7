package sgc7game

import (
	jsoniter "github.com/json-iterator/go"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// Config - config
type Config struct {
	Lines        *LineData             `json:"lines"`
	Reels        map[string]*ReelsData `json:"reels"`
	PayTables    *PayTables            `json:"paytables"`
	Width        int                   `json:"width"`
	Height       int                   `json:"height"`
	DefaultScene *GameScene            `json:"defaultscene"`
	Ver          string                `json:"ver"`
	CoreVer      string                `json:"corever"`
}

// NewConfig - new a Config
func NewConfig() *Config {
	return &Config{
		Reels: make(map[string]*ReelsData),
	}
}

// LoadLine5 - load linedata for reels 5
func (cfg *Config) LoadLine5(fn string) error {
	ld, err := LoadLine5JSON(fn)
	if err != nil {
		return err
	}

	cfg.Lines = ld

	return nil
}

// LoadLine3 - load linedata for reels 3
func (cfg *Config) LoadLine3(fn string) error {
	ld, err := LoadLine3JSON(fn)
	if err != nil {
		return err
	}

	cfg.Lines = ld

	return nil
}

// LoadPayTables5 - load paytables for reels 5
func (cfg *Config) LoadPayTables5(fn string) error {
	pt, err := LoadPayTables5JSON(fn)
	if err != nil {
		return err
	}

	cfg.PayTables = pt

	return nil
}

// LoadReels5 - load reels 5
func (cfg *Config) LoadReels5(name string, fn string) error {
	reels, err := LoadReels5JSON(fn)
	if err != nil {
		return err
	}

	cfg.Reels[name] = reels

	return nil
}

// SetDefaultSceneString - [][]int in json
func (cfg *Config) SetDefaultSceneString(str string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	var arr [][]int
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		return err
	}

	ds, err := NewGameSceneWithArr2(arr)
	if err != nil {
		sgc7utils.Error("sgc7game.Config.SetDefaultSceneString:NewGameSceneWithArr2",
			zap.String("str", str),
			zap.Error(err))

		return err
	}

	cfg.DefaultScene = ds

	return nil
}
