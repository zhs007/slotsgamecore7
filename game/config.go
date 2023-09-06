package sgc7game

import (
	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// Config - config
type Config struct {
	Lines         *LineData                     `json:"lines"`
	Reels         map[string]*ReelsData         `json:"reels"`
	PayTables     *PayTables                    `json:"paytables"`
	Width         int                           `json:"width"`
	Height        int                           `json:"height"`
	DefaultScene  *GameScene                    `json:"defaultscene"`
	Ver           string                        `json:"ver"`
	CoreVer       string                        `json:"corever"`
	SWReels       map[string]*SymbolWeightReels `json:"-"`
	DefaultScene2 []*GameScene                  `json:"defaultscene2"`
	BetMuls       []int32                       `json:"betMuls"`
}

// NewConfig - new a Config
func NewConfig() *Config {
	return &Config{
		Reels:   make(map[string]*ReelsData),
		SWReels: make(map[string]*SymbolWeightReels),
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

// LoadLine6 - load linedata for reels 6
func (cfg *Config) LoadLine6(fn string) error {
	ld, err := LoadLine6JSON(fn)
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
	} else if reels == 6 {
		return cfg.LoadLine6(fn)
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

// LoadPayTables6 - load paytables for reels 6
func (cfg *Config) LoadPayTables6(fn string) error {
	pt, err := LoadPayTables6JSON(fn)
	if err != nil {
		return err
	}

	cfg.PayTables = pt

	return nil
}

// LoadPayTables15 - load paytables for reels 15
func (cfg *Config) LoadPayTables15(fn string) error {
	pt, err := LoadPayTables15JSON(fn)
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
	} else if reels == 6 {
		return cfg.LoadPayTables6(fn)
	} else if reels == 15 {
		return cfg.LoadPayTables15(fn)
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

// LoadSymboloWeightReels - load reels for SymbolWeightReels
func (cfg *Config) LoadSymboloWeightReels(name string, fn string, reels int) error {
	swreels, err := LoadSymbolWeightReels5JSON(fn)
	if err != nil {
		return err
	}

	cfg.SWReels[name] = swreels

	return nil
}

// SetDefaultSceneString - [][]int in json
func (cfg *Config) SetDefaultSceneString(str string) error {
	var arr [][]int
	err := sonic.Unmarshal([]byte(str), &arr)
	if err != nil {
		return err
	}

	ds, err := NewGameSceneWithArr2Ex(arr)
	if err != nil {
		goutils.Error("sgc7game.Config.SetDefaultSceneString:NewGameSceneWithArr2Ex",
			zap.String("str", str),
			zap.Error(err))

		return err
	}

	cfg.DefaultScene = ds

	return nil
}

// AddDefaultSceneString2 - [][]int in json
func (cfg *Config) AddDefaultSceneString2(str string) error {
	var arr [][]int
	err := sonic.Unmarshal([]byte(str), &arr)
	if err != nil {
		return err
	}

	ds, err := NewGameSceneWithArr2Ex(arr)
	if err != nil {
		goutils.Error("sgc7game.Config.AddDefaultSceneString2:NewGameSceneWithArr2Ex",
			zap.String("str", str),
			zap.Error(err))

		return err
	}

	cfg.DefaultScene2 = append(cfg.DefaultScene2, ds)

	return nil
}

// LoadLineDataFromExcel - load linedata for reels
func (cfg *Config) LoadLineDataFromExcel(fn string) error {
	ld, err := LoadLineDataFromExcel(fn)
	if err != nil {
		goutils.Error("Config.LoadLineDataFromExcel",
			zap.Error(err))

		return err
	}

	cfg.Lines = ld

	return nil
}

// LoadReelsFromExcel - load linedata for reels
func (cfg *Config) LoadReelsFromExcel(tag string, fn string) error {
	rd, err := LoadReelsFromExcel(fn)
	if err != nil {
		goutils.Error("Config.LoadReelsFromExcel",
			zap.Error(err))

		return err
	}

	cfg.Reels[tag] = rd

	return nil
}

// LoadPaytablesFromExcel - load linedata for reels
func (cfg *Config) LoadPaytablesFromExcel(fn string) error {
	pt, err := LoadPaytablesFromExcel(fn)
	if err != nil {
		goutils.Error("Config.LoadPaytablesFromExcel",
			zap.Error(err))

		return err
	}

	cfg.PayTables = pt

	return nil
}
