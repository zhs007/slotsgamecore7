package sgc7game

// Config - config
type Config struct {
	Line         *LineData             `json:"line"`
	Reels        map[string]*ReelsData `json:"reels"`
	PayTables    *PayTables            `json:"paytables"`
	Width        int                   `json:"width"`
	Height       int                   `json:"height"`
	DefaultScene *GameScene            `json:"defaultscene"`
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
