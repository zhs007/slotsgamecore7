package sgc7game

// BasicGame - basic game
type BasicGame struct {
	Cfg Config
}

// GetConfig - get config
func (game *BasicGame) GetConfig() *Config {
	return &game.Cfg
}
