package sgc7game

import sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"

// IGameMod - game
type IGameMod interface {
	// GetName - get mode name
	GetName() string

	// OnPlay - on play
	OnPlay(game IGame, plugin sgc7plugin.IPlugin, cmd string, param string, stake *Stake, prs []*PlayResult) (*PlayResult, error)
}
