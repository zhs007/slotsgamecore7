package lowcode

import (
	"github.com/bytedance/sonic"
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
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

		return err
	}

	err = game.BuildGameConfigData()
	if err != nil {
		goutils.Error("Game.Init2:BuildGameConfigData",
			goutils.Err(err))

		return err
	}

	pool.onInit()

	if cfg.DefaultScene == "" {
		gs, err := GenDefaultScene(game, cfg.Bets[0])
		if err != nil {
			goutils.Error("Game.Init2:GenDefaultScene",
				goutils.Err(err))

			return err
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
// NewPlayerState 用于 new 一个空的 playerstate，不需要 initial ，后面会 reset 数据
// Initialize 用于直接 生成一个 playerstate，并初始化它
func (game *Game) NewPlayerState() sgc7game.IPlayerState {
	ps, err := game.Pool.NewPlayerState()
	if err != nil {
		goutils.Error("Game.NewPlayerState:NewPlayerState",
			goutils.Err(err))

		return nil
	}

	return ps
}

// Initialize - initialize PlayerState
func (game *Game) Initialize() sgc7game.IPlayerState {
	ps, err := game.Pool.InitPlayerState()
	if err != nil {
		goutils.Error("Game.Initialize:InitPlayerState",
			goutils.Err(err))

		return nil
	}

	return ps
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

func (game *Game) SaveParSheet(f *excelize.File) error {
	err := SavePaytable(f, "paytable", game.Cfg.PayTables)
	if err != nil {
		goutils.Error("Game.SaveParSheet:SavePaytable",
			goutils.Err(err))

		return err
	}

	for rn, r := range game.Pool.Config.MapReels {
		err = SaveReels(f, rn, game.Cfg.PayTables, r)
		if err != nil {
			goutils.Error("Game.SaveParSheet:SaveReels",
				goutils.Err(err))

			return err
		}
	}

	for ln, l := range game.Pool.Config.MapLinedate {
		err = SaveLineData(f, ln, l)
		if err != nil {
			goutils.Error("Game.SaveParSheet:SaveLineData",
				goutils.Err(err))

			return err
		}
	}

	for vn, vw := range game.Pool.mapSymbolValWeights {
		err = SaveSymbolWeights(f, vn, game.Cfg.PayTables, vw)
		if err != nil {
			goutils.Error("Game.SaveParSheet:SaveSymbolWeights",
				goutils.Err(err))

			return err
		}
	}

	for vn, vw := range game.Pool.mapIntValWeights {
		err = SaveIntWeights(f, vn, vw)
		if err != nil {
			goutils.Error("Game.SaveParSheet:SaveIntWeights",
				goutils.Err(err))

			return err
		}
	}

	for vn, vw := range game.Pool.mapStrValWeights {
		err = SaveStrWeights(f, vn, vw)
		if err != nil {
			goutils.Error("Game.SaveParSheet:SaveStrWeights",
				goutils.Err(err))

			return err
		}
	}

	return nil
}

// OnBet
func (game *Game) OnBet(plugin sgc7plugin.IPlugin, cmd string, param string, ips sgc7game.IPlayerState,
	stake *sgc7game.Stake, prs []*sgc7game.PlayResult, gameData any) error {
	gameProp, isok := gameData.(*GameProperty)
	if !isok {
		goutils.Error("Game.OnBet:GameProperty",
			goutils.Err(ErrIvalidGameData))

		return ErrIvalidGameData
	}

	ps, isok := ips.(*PlayerState)
	if !isok {
		goutils.Error("Game.OnBet:PlayerState",
			goutils.Err(ErrIvalidPlayerState))

		return ErrIvalidPlayerState
	}

	game.Pool.InitPlayerStateOnBet(gameProp, plugin, ps, stake)

	return nil
}
