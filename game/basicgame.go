package sgc7game

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// BasicGame - basic game
type BasicGame struct {
	Cfg         *Config
	MapGameMods map[string]IGameMod
	Plugin      sgc7plugin.IPlugin
}

// NewBasicGame - new a BasicGame
func NewBasicGame() BasicGame {
	return BasicGame{
		Cfg:         NewConfig(),
		MapGameMods: make(map[string]IGameMod),
	}
}

// GetConfig - get config
func (game *BasicGame) GetConfig() *Config {
	return game.Cfg
}

// GetPlugin - get plugin
func (game *BasicGame) GetPlugin() sgc7plugin.IPlugin {
	return game.Plugin
}

// NewPlayerState - new playerstate
func (game *BasicGame) NewPlayerState() IPlayerState {
	return &BasicPlayerState{}
}

// Initialize - initialize PlayerState
func (game *BasicGame) Initialize() IPlayerState {
	return NewBasicPlayerState("BG")
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
func (game *BasicGame) Play(cmd string, param string, ps IPlayerState, stake *Stake, prs []*PlayResult) (*PlayResult, error) {
	bps, isok := ps.(*BasicPlayerState)
	if !isok {
		return nil, ErrInvalidBasicPlayerState
	}

	curgamemod, isok := game.MapGameMods[bps.Public.CurGameMod]
	if !isok {
		sgc7utils.Error("sgc7game.BasicGame.Play:MapGameMods[CurGameMod]", 
			zap.String("CurGameMod", bps.Public.CurGameMod),
			zap.Error(ErrInvalidGameMod))

		return nil, ErrInvalidGameMod
	}

	pr, err := curgamemod.OnPlay(game, cmd, param, stake, prs)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
