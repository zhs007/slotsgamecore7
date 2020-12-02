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

// LoadLine - load linedata for reels
func (cfg *Config) LoadLine(fn string, reels int) error {
	if reels == 5 {
		return cfg.LoadLine5(fn)
	} else if reels == 3 {
		return cfg.LoadLine3(fn)
	}

	return ErrInvalidReels
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

// LoadPayTables3 - load paytables for reels 3
func (cfg *Config) LoadPayTables3(fn string) error {
	pt, err := LoadPayTables3JSON(fn)
	if err != nil {
		return err
	}

	cfg.PayTables = pt

	return nil
}

// LoadPayTables - load paytables for reels
func (cfg *Config) LoadPayTables(fn string, reels int) error {
	if reels == 5 {
		return cfg.LoadPayTables5(fn)
	} else if reels == 3 {
		return cfg.LoadPayTables3(fn)
	}

	return ErrInvalidReels
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

// LoadReels3 - load reels 3
func (cfg *Config) LoadReels3(name string, fn string) error {
	reels, err := LoadReels3JSON(fn)
	if err != nil {
		return err
	}

	cfg.Reels[name] = reels

	return nil
}

// LoadReels - load reels for reels
func (cfg *Config) LoadReels(name string, fn string, reels int) error {
	if reels == 5 {
		return cfg.LoadReels5(name, fn)
	} else if reels == 3 {
		return cfg.LoadReels3(name, fn)
	}

	return ErrInvalidReels
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
