package sgc7game

import sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"

// BasicGame - basic game
type BasicGame struct {
	Cfg         Config
	MapGameMods map[string]IGameMod
	Plugin      sgc7plugin.IPlugin
}

// NewBasicGame - new a BasicGame
func NewBasicGame() BasicGame {
	return BasicGame{
		MapGameMods: make(map[string]IGameMod),
	}
}

// GetConfig - get config
func (game *BasicGame) GetConfig() *Config {
	return &game.Cfg
}

// GetPlugin - get plugin
func (game *BasicGame) GetPlugin() sgc7plugin.IPlugin {
	return game.Plugin
}

// Initialize - initialize PlayerState
func (game *BasicGame) Initialize() IPlayerState {
	return nil
}

// AddGameMod - add a gamemod
func (game *BasicGame) AddGameMod(name string, gmod IGameMod) error {
	_, isok := game.MapGameMods[name]
	if isok {
		return ErrDuplicateGameMod
	}

	game.MapGameMods[name] = gmod

	return nil
}
