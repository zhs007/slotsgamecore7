package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

// Game - game
type Game struct {
	*sgc7game.BasicGame
	Prop         *GameProperty
	MgrComponent *ComponentMgr
}

// Init - initial game
func (game *Game) Init(cfgfn string) error {
	prop, err := InitGameProperty(cfgfn)
	if err != nil {
		goutils.Error("Game.Init:InitGameProperty",
			zap.String("fn", cfgfn),
			zap.Error(err))

		return err
	}

	game.Prop = prop

	game.Cfg.PayTables = prop.CurPaytables
	game.SetVer(sgc7ver.Version)

	game.Cfg.SetDefaultSceneString(game.Prop.Config.DefaultScene)

	for _, v := range prop.Config.GameMods {
		game.AddGameMod(NewBasicGameMod(prop, v, game.MgrComponent))
	}

	return nil
}

// CheckStake - check stake
func (game *Game) CheckStake(stake *sgc7game.Stake) error {
	if goutils.IndexOfIntSlice(game.Prop.Config.Bets, int(stake.CashBet)/int(stake.CoinBet), 0) < 0 {
		return sgc7game.ErrInvalidStake
	}

	return nil
}

// NewPlayerState - new playerstate
func (game *Game) NewPlayerState() sgc7game.IPlayerState {
	bps := sgc7game.NewBasicPlayerState("bg")

	return bps
}

// ResetConfig
func (game *Game) ResetConfig(cfg interface{}) {
	ncfg := cfg.(*Config)

	for _, v := range game.Prop.Config.GameMods {
		gm := game.MapGameMods[v.Type].(*BasicGameMod)
		gm.ResetConfig(ncfg)
	}
}

// OnAsciiGame - outpur to asciigame
func (game *Game) OnAsciiGame(stake *sgc7game.Stake, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) error {
	for _, v := range game.Prop.Config.GameMods {
		gm := game.MapGameMods[v.Type].(*BasicGameMod)
		gm.OnAsciiGame(game.Prop, pr, lst)
	}

	if pr.IsFinish {
		if game.Prop.Stats != nil {
			game.Prop.Stats.OnResults(stake, lst)
		}
	}

	return nil
}

// NewGame - new a Game
func NewGame(cfgfn string) (*Game, error) {
	game := &Game{
		BasicGame: sgc7game.NewBasicGame(func() sgc7plugin.IPlugin {
			return sgc7plugin.NewBasicPlugin()
		}),
		MgrComponent: NewComponentMgr(),
	}

	err := game.Init(cfgfn)
	if err != nil {
		return nil, err
	}

	return game, nil
}

// NewGame - new a Game
func NewGameEx(cfgfn string, funcNewPlugin sgc7plugin.FuncNewPlugin) (*Game, error) {
	game := &Game{
		BasicGame:    sgc7game.NewBasicGame(funcNewPlugin),
		MgrComponent: NewComponentMgr(),
	}

	err := game.Init(cfgfn)
	if err != nil {
		return nil, err
	}

	return game, nil
}
