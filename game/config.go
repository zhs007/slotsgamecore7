package sgc7game

import jsoniter "github.com/json-iterator/go"

// Config - config
type Config struct {
	Line         *LineData             `json:"line"`
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

	cfg.Line = ld

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

	cfg.DefaultScene = NewGameScene(len(arr), len(arr[0]))

	return nil
}
