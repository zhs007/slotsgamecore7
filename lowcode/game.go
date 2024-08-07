package lowcode

import (
	"github.com/bytedance/sonic"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

// Game - game
type Game struct {
	*sgc7game.BasicGame
	Pool         *GamePropertyPool
	MgrComponent *ComponentMgr
}

// Init - initial game
func (game *Game) Init2(cfg *Config, funcNewRNG FuncNewRNG, funcNewFeatureLevel FuncNewFeatureLevel) error {
	pool, err := newGamePropertyPool2(cfg, funcNewRNG, funcNewFeatureLevel)
	if err != nil {
		goutils.Error("Game.Init2:NewGamePropertyPool2",
			goutils.Err(err))

		return err
	}

	game.Pool = pool

	game.Cfg.PayTables = pool.DefaultPaytables
	game.SetVer(sgc7ver.Version)

	game.Cfg.SetDefaultSceneString(cfg.DefaultScene)

	pool.loadAllWeights()

	gamemod, err := NewBasicGameMod2(pool, game.MgrComponent)
	if err != nil {
		goutils.Error("Game.Init2:NewBasicGameMod2",
			goutils.Err(err))

		return err
	}

	game.AddGameMod(gamemod)

	err = pool.InitStats(pool.Config.Bets[0])
	if err != nil {
		goutils.Error("Game.Init2:InitStats",
			goutils.Err(err))

		return nil
	}

	err = game.BuildGameConfigData()
	if err != nil {
		goutils.Error("Game.Init2:BuildGameConfigData",
			goutils.Err(err))

		return nil
	}

	pool.onInit()

	if cfg.DefaultScene == "" {
		gs, err := GenDefaultScene(game, cfg.Bets[0])
		if err != nil {
			goutils.Error("Game.Init2:GenDefaultScene",
				goutils.Err(err))

			return nil
		}

		cfg.DefaultScene = gs.ToString()

		game.Cfg.SetDefaultSceneString(cfg.DefaultScene)
	}

	return nil
}

// CheckStake - check stake
func (game *Game) CheckStake(stake *sgc7game.Stake) error {
	if goutils.IndexOfIntSlice(game.Pool.Config.Bets, int(stake.CashBet)/int(stake.CoinBet), 0) < 0 {
		return sgc7game.ErrInvalidStake
	}

	return nil
}

// NewPlayerState - new playerstate
func (game *Game) NewPlayerState() sgc7game.IPlayerState {
	return &sgc7game.BasicPlayerState{}
}

// Initialize - initialize PlayerState
func (game *Game) Initialize() sgc7game.IPlayerState {
	bps := sgc7game.NewBasicPlayerState(BasicGameModName)

	return bps
}

// ResetConfig
func (game *Game) ResetConfig(cfg any) {
	ncfg := cfg.(*Config)

	gm := game.MapGameMods[BasicGameModName].(*BasicGameMod)
	gm.ResetConfig(ncfg)
}

// OnAsciiGame - outpur to asciigame
func (game *Game) OnAsciiGame(gameProp *GameProperty, stake *sgc7game.Stake, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) error {
	gm := game.MapGameMods[BasicGameModName].(*BasicGameMod)
	gm.OnAsciiGame(gameProp, pr, lst)

	return nil
}

// NewGameData - new GameData
func (game *Game) NewGameData(stake *sgc7game.Stake) sgc7game.IGameData {
	pool, isok := game.Pool.MapGamePropPool[int(stake.CashBet)/int(stake.CoinBet)]
	if !isok {
		return nil
	}

	gameProp := pool.Get().(*GameProperty)

	return gameProp
}

// DeleteGameData - delete GameData
func (game *Game) DeleteGameData(gamed sgc7game.IGameData) {
	game.Pool.MapGamePropPool[gamed.GetBetMul()].Put(gamed)
}

// BuildGameConfigData - build game configration data
func (game *Game) BuildGameConfigData() error {
	buf, err := sonic.Marshal(game.Pool.mapComponents)
	if err != nil {
		goutils.Error("Game.BuildGameConfigData:Marshal",
			goutils.Err(err))

		return err
	}

	game.Cfg.Data = string(buf)

	return nil
}
