package sgc7game

import sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"

// IGame - game
type IGame interface {
	// NewPlugin - new a plugin
	NewPlugin() sgc7plugin.IPlugin
	// FreePlugin - free a plugin
	FreePlugin(plugin sgc7plugin.IPlugin)

	// NewPlayerState - new playerstate
	NewPlayerState() IPlayerState
	// SetVer - set server version
	SetVer(ver string)

	// GetConfig - get config
	GetConfig() *Config
	// Initialize - initialize PlayerState
	Initialize() IPlayerState

	// CheckStake - check stake
	CheckStake(stake *Stake) error
	// Play - play
	Play(plugin sgc7plugin.IPlugin, cmd string, param string, ps IPlayerState, stake *Stake, prs []*PlayResult, gameData any) (*PlayResult, error)
	// NewGameData - new GameData
	NewGameData(stake *Stake) IGameData
	// DeleteGameData - delete GameData
	DeleteGameData(gamed IGameData)

	// AddGameMod - add a gamemod
	AddGameMod(gmod IGameMod) error

	// ResetConfig
	ResetConfig(cfg any)

	// OnBet
	OnBet(plugin sgc7plugin.IPlugin, cmd string, param string, ps IPlayerState, stake *Stake, prs []*PlayResult, gameData any) error
}
