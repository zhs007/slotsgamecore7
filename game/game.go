package sgc7game

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// IGame - game
type IGame interface {
	// GetPlugin - get plugin
	GetPlugin() sgc7plugin.IPlugin

	// GetConfig - get config
	GetConfig() *Config
	// Initialize - initialize PlayerState
	Initialize() IPlayerState

	// AddGameMod - add a gamemod
	AddGameMod(name string, gmod IGameMod) error
}
