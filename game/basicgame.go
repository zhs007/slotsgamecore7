package sgc7game

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

// BasicGame - basic game
type BasicGame struct {
	Cfg         *Config
	MapGameMods map[string]IGameMod
	MgrPlugins  *sgc7plugin.PluginsMgr
}

// NewBasicGame - new a BasicGame
func NewBasicGame(funcNewPlugin sgc7plugin.FuncNewPlugin) *BasicGame {
	return &BasicGame{
		Cfg:         NewConfig(),
		MapGameMods: make(map[string]IGameMod),
		MgrPlugins:  sgc7plugin.NewPluginsMgr(funcNewPlugin),
	}
}

// GetConfig - get config
func (game *BasicGame) GetConfig() *Config {
	return game.Cfg
}

// NewPlugin - new a plugin
func (game *BasicGame) NewPlugin() sgc7plugin.IPlugin {
	return game.MgrPlugins.NewPlugin()
}

// FreePlugin - free a plugin
func (game *BasicGame) FreePlugin(plugin sgc7plugin.IPlugin) {
	game.MgrPlugins.FreePlugin(plugin)
}

// SetVer - set server version
func (game *BasicGame) SetVer(ver string) {
	game.Cfg.Ver = ver
	game.Cfg.CoreVer = sgc7ver.Version
}

// NewPlayerState - new playerstate
func (game *BasicGame) NewPlayerState() IPlayerState {
	return &BasicPlayerState{}
}

// Initialize - initialize PlayerState
func (game *BasicGame) Initialize() IPlayerState {
	return NewBasicPlayerState("bg")
}

// AddGameMod - add a gamemod
func (game *BasicGame) AddGameMod(gmod IGameMod) error {
	_, isok := game.MapGameMods[gmod.GetName()]
	if isok {
		return ErrDuplicateGameMod
	}

	game.MapGameMods[gmod.GetName()] = gmod

	return nil
}

// Play - play
func (game *BasicGame) Play(plugin sgc7plugin.IPlugin, cmd string, param string, ps IPlayerState, stake *Stake, prs []*PlayResult) (*PlayResult, error) {
	// bps, isok := ps.(*BasicPlayerState)
	// if !isok {
	// 	return nil, ErrInvalidBasicPlayerState
	// }

	curgamemod, isok := game.MapGameMods[ps.GetCurGameMod()]
	if !isok {
		sgc7utils.Error("sgc7game.BasicGame.Play:MapGameMods[CurGameMod]",
			zap.String("CurGameMod", ps.GetCurGameMod()),
			zap.Error(ErrInvalidGameMod))

		return nil, ErrInvalidGameMod
	}

	pr, err := curgamemod.OnPlay(game, plugin, cmd, param, ps, stake, prs)
	if err != nil {
		return nil, err
	}

	ps.SetCurGameMod(pr.NextGameMod)
	// bps.Public.CurGameMod = pr.NextGameMod

	return pr, nil
}

// CheckStake - check stake
func (game *BasicGame) CheckStake(stake *Stake) error {
	return ErrInvalidStake
}
