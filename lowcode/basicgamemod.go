package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

// BasicGameMod - basic gamemod
type BasicGameMod struct {
	*sgc7game.BasicGameMod
	Pool          *GamePropertyPool
	MapComponents map[int]*ComponentList
}

// OnPlay - on play
func (bgm *BasicGameMod) newPlayResult(prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, *GameParams) {
	gp := NewGameParam()
	gp.MapComponents = make(map[string]*anypb.Any)

	pr := &sgc7game.PlayResult{
		IsFinish:         true,
		CurGameMod:       "bg",
		NextGameMod:      "bg",
		CurGameModParams: gp,
	}

	if len(prs) > 0 {
		lastrs := prs[len(prs)-1]
		lastgp := lastrs.CurGameModParams.(*GameParams)

		gp.FirstComponent = lastgp.NextStepFirstComponent
	}

	return pr, gp
}

// OnPlay - on play
func (bgm *BasicGameMod) OnPlay(game sgc7game.IGame, plugin sgc7plugin.IPlugin, cmd string, param string,
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, gameData any) (*sgc7game.PlayResult, error) {

	gameProp, isok := gameData.(*GameProperty)
	if !isok {
		goutils.Error("BasicGameMod.OnPlay",
			zap.Error(ErrIvalidGameData))

		return nil, ErrIvalidGameData
	}

	components := bgm.MapComponents[int(stake.CashBet/stake.CoinBet)]
	gameProp.Components = components

	if len(prs) == 0 {
		bgm.OnNewGame(gameProp, stake)
	}

	bgm.OnNewStep(gameProp, stake)

	pr, gp := bgm.newPlayResult(prs)

	gameProp.SceneStack.onStepStart(pr)

	curComponent := components.Components[0]

	if gp.FirstComponent != "" {
		c, isok := components.MapComponents[gp.FirstComponent]
		if !isok {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
				zap.String("FirstComponent", gp.FirstComponent),
				zap.Error(ErrIvalidComponentName))

			return nil, ErrIvalidComponentName
		}

		curComponent = c
	} else {
		if cmd != DefaultCmd {
			cn, isok := bgm.Pool.Config.MapCmdComponent[cmd]
			if isok {
				c, isok1 := components.MapComponents[cn]
				if !isok1 {
					goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
						zap.String("cmd", cmd),
						zap.String("MapCmdComponent", cn),
						zap.Error(ErrIvalidComponentName))

					return nil, ErrIvalidComponentName
				}

				curComponent = c
			}
		} else {
			startComponent, isok := gameProp.Pool.Config.StartComponents[int(stake.CashBet/stake.CoinBet)]
			if isok {
				c, isok := components.MapComponents[startComponent]
				if !isok {
					goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
						zap.String("FirstComponent", startComponent),
						zap.Error(ErrIvalidComponentName))

					return nil, ErrIvalidComponentName
				}

				curComponent = c
			}
		}
	}

	for {
		isComponentDoNothing := false
		err := curComponent.OnPlayGame(gameProp, pr, gp, plugin, cmd, param, ps, stake, prs)
		if err != nil {
			if err != ErrComponentDoNothing {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					zap.Error(err))

				return nil, err
			}

			isComponentDoNothing = true
		}

		if !isComponentDoNothing {
			gameProp.HistoryComponents = append(gameProp.HistoryComponents, curComponent)
			gp.HistoryComponents = append(gp.HistoryComponents, curComponent.GetName())
		}

		respinComponent := gameProp.GetStrVal(GamePropRespinComponent)
		nextComponentName := gameProp.GetStrVal(GamePropNextComponent)
		if respinComponent != "" {
			// 一般来说，第一次触发respin才走这个分支
			pr.IsFinish = false

			if nextComponentName == "" {
				break
			}
		}

		if nextComponentName == "" {
			break
		}

		c, isok := components.MapComponents[nextComponentName]
		if !isok {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
				zap.String("nextComponentName", nextComponentName),
				zap.Error(ErrIvalidComponentName))

			return nil, ErrIvalidComponentName
		}

		curComponent = c
	}

	gameProp.BuildGameParam(gp)

	for _, v := range gameProp.HistoryComponents {
		err := v.OnPlayGameEnd(gameProp, pr, gp, plugin, cmd, param, ps, stake, prs)
		if err != nil {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGameEnd",
				zap.Error(err))

			return nil, err
		}

		cn := v.GetName()
		gp.AddComponentData(cn, gameProp.MapComponentData[cn])

		if gAllowStats2 {
			components.Stats2.onStepStats(v, gameProp.MapComponentData[cn])
			gameProp.stats2SpinData.onStepStats(v, gameProp.MapComponentData[cn])
		}
	}

	gameProp.ProcRespin(pr, gp)

	gameProp.onStepEnd(pr, prs)

	if gAllowStats2 && pr.IsFinish {
		components.Stats2.pushBetEnding()

		gameProp.stats2SpinData.onBetEnding(components.Stats2)

		// for _, curpr := range prs {
		// 	curpr.
		// }
	}
	// if pr.IsFinish {
	// 	for _, curpr := range prs {
	// 		for _, s := range curpr.Scenes {
	// 			gameProp.PoolScene.Put(s)
	// 		}

	// 		for _, s := range curpr.OtherScenes {
	// 			gameProp.PoolScene.Put(s)
	// 		}

	// 		for _, s := range curpr.PrizeScenes {
	// 			gameProp.PoolScene.Put(s)
	// 		}
	// 	}

	// 	for _, s := range pr.Scenes {
	// 		gameProp.PoolScene.Put(s)
	// 	}

	// 	for _, s := range pr.OtherScenes {
	// 		gameProp.PoolScene.Put(s)
	// 	}

	// 	for _, s := range pr.PrizeScenes {
	// 		gameProp.PoolScene.Put(s)
	// 	}
	// }

	return pr, nil
}

// ResetConfig
func (bgm *BasicGameMod) ResetConfig(cfg *Config) {
	bgm.Pool.Config = cfg
}

// OnAsciiGame - outpur to asciigame
func (bgm *BasicGameMod) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) error {
	for _, v := range gameProp.HistoryComponents {
		v.OnAsciiGame(gameProp, pr, lst, gameProp.Pool.MapSymbolColor)
	}

	return nil
}

// OnNewGame -
func (bgm *BasicGameMod) OnNewGame(gameProp *GameProperty, stake *sgc7game.Stake) error {
	if gAllowStats2 {
		gameProp.Components.Stats2.pushBet(stake.CashBet / stake.CoinBet)

		gameProp.stats2SpinData = newStats2SpinData()
	}

	gameProp.OnNewGame(stake)

	// components := bgm.MapComponents[int(stake.CashBet/stake.CoinBet)]

	for i, v := range gameProp.Components.Components {
		err := v.OnNewGame(gameProp)
		if err != nil {
			goutils.Error("BasicGameMod.OnNewGame:OnNewGame",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}
	}

	return nil
}

// OnNewStep -
func (bgm *BasicGameMod) OnNewStep(gameProp *GameProperty, stake *sgc7game.Stake) error {
	if gAllowStats2 {
		gameProp.Components.Stats2.pushStep()
	}

	gameProp.OnNewStep()

	// components := bgm.MapComponents[int(stake.CashBet/stake.CoinBet)]

	for i, v := range gameProp.Components.Components {
		err := v.OnNewStep(gameProp)
		if err != nil {
			goutils.Error("BasicGameMod.OnNewStep:OnNewStep",
				zap.Int("i", i),
				zap.Error(err))

			return err
		}
	}

	return nil
}

// NewBasicGameMod - new BaseGame
func NewBasicGameMod(pool *GamePropertyPool, cfgGameMod *GameModConfig, mgrComponent *ComponentMgr) *BasicGameMod {
	bgm := &BasicGameMod{
		BasicGameMod:  sgc7game.NewBasicGameMod(cfgGameMod.Type, pool.Config.Width, pool.Config.Height),
		Pool:          pool,
		MapComponents: make(map[int]*ComponentList),
	}

	for _, bet := range pool.Config.Bets {
		components := NewComponentList()
		mapComponentMapping := pool.Config.ComponentsMapping[bet]

		for _, v := range cfgGameMod.Components {
			c := mgrComponent.NewComponent(v)
			configfn := v.Config
			if mapComponentMapping != nil {
				mappingfn, isok := mapComponentMapping[v.Name]
				if isok {
					configfn = mappingfn
				}
			}

			err := c.Init(pool.Config.GetPath(configfn, false), pool)
			if err != nil {
				goutils.Error("NewBasicGameMod:Init",
					zap.Error(err))

				return nil
			}

			components.AddComponent(v.Name, c)
		}

		bgm.MapComponents[bet] = components
		pool.onAddComponentList(bet, components)
	}

	return bgm
}

// NewBasicGameMod2 - new BaseGame
func NewBasicGameMod2(pool *GamePropertyPool, cfgGameMod *GameModConfig, mgrComponent *ComponentMgr) (*BasicGameMod, error) {
	bgm := &BasicGameMod{
		BasicGameMod:  sgc7game.NewBasicGameMod(cfgGameMod.Type, pool.Config.Width, pool.Config.Height),
		Pool:          pool,
		MapComponents: make(map[int]*ComponentList),
	}

	for _, bet := range pool.Config.Bets {
		components := NewComponentList()

		for _, v := range cfgGameMod.Components {
			c := mgrComponent.NewComponent(v)

			err := c.InitEx(pool.Config.mapConfig[v.Name], pool)
			if err != nil {
				goutils.Error("NewBasicGameMod:Init",
					zap.Error(err))

				return nil, err
			}

			components.AddComponent(v.Name, c)
		}

		bgm.MapComponents[bet] = components
		pool.onAddComponentList(bet, components)
	}

	return bgm, nil
}
