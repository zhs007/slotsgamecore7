package lowcode

import (
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/types/known/anypb"
)

// BasicGameMod - basic gamemod
type BasicGameMod struct {
	*sgc7game.BasicGameMod
	Pool          *GamePropertyPool
	MapComponents map[int]*ComponentList
}

// OnPlay - on play
func (bgm *BasicGameMod) newPlayResult(prs []*sgc7game.PlayResult, stake *sgc7game.Stake, ps *PlayerState) (*sgc7game.PlayResult, *GameParams) {
	gp := NewGameParam(stake, ps)
	gp.MapComponents = make(map[string]*anypb.Any)

	pr := &sgc7game.PlayResult{
		IsFinish:         true,
		CurGameMod:       BasicGameModName,
		NextGameMod:      BasicGameModName,
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
			goutils.Err(ErrInvalidGameData))

		return nil, ErrInvalidGameData
	}

	curBetMode := int(stake.CashBet / stake.CoinBet)

	components := bgm.MapComponents[curBetMode]
	gameProp.Components = components

	if len(prs) == 0 {
		bgm.OnNewGame(gameProp, stake, plugin)

		gameProp.featureLevel.Init()
	}

	bgm.OnNewStep(gameProp, stake)

	pr, gp := bgm.newPlayResult(prs, stake, ps.(*PlayerState))

	gameProp.SceneStack.onStepStart(pr)
	gameProp.OtherSceneStack.onStepStart(pr)

	var curComponent IComponent

	if gp.FirstComponent != "" {
		c, isok := components.MapComponents[gp.FirstComponent]
		if !isok {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
				slog.String("FirstComponent", gp.FirstComponent),
				goutils.Err(ErrInvalidComponentName))

			return nil, ErrInvalidComponentName
		}

		curComponent = c
	} else {
		if cmd != DefaultCmd {
			cn, isok := bgm.Pool.Config.MapCmdComponent[cmd]
			if isok {
				c, isok1 := components.MapComponents[cn]
				if !isok1 {
					goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
						slog.String("cmd", cmd),
						slog.String("MapCmdComponent", cn),
						goutils.Err(ErrInvalidComponentName))

					return nil, ErrInvalidComponentName
				}

				curComponent = c
			} else {
				c, isok1 := components.MapComponents[cmd]
				if !isok1 {
					goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
						slog.String("cmd", cmd),
						goutils.Err(ErrInvalidComponentName))

					return nil, ErrInvalidComponentName
				}

				curComponent = c
			}
		} else {
			startComponent := gameProp.Pool.Config.MapBetConfigs[int(stake.CashBet/stake.CoinBet)].Start
			// if isok {
			c, isok := components.MapComponents[startComponent]
			if !isok {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					slog.String("FirstComponent", startComponent),
					goutils.Err(ErrInvalidComponentName))

				return nil, ErrInvalidComponentName
			}

			curComponent = c
			// }
		}
	}

	for {
		isComponentDoNothing := false
		isSetMode, set, currng, newComponent := gameProp.rng.GetCurRNG(curBetMode, gameProp, curComponent, gameProp.callStack.GetCurComponentData(gameProp, curComponent), gameProp.featureLevel)
		if newComponent != "" {
			c, isok := components.MapComponents[newComponent]
			if !isok {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					slog.String("newComponent", newComponent),
					goutils.Err(ErrInvalidComponentName))

				return nil, ErrInvalidComponentName
			}

			curComponent = c
		}

		cd := gameProp.callStack.GetCurComponentData(gameProp, curComponent)

		nextComponentName := ""
		var err error

		if isSetMode {
			nextComponentName, err = curComponent.OnPlayGameWithSet(gameProp, pr, gp, currng, cmd, param, ps, stake, prs, cd, set)
		} else {
			nextComponentName, err = curComponent.OnPlayGame(gameProp, pr, gp, currng, cmd, param, ps, stake, prs, cd)
		}

		if err != nil {
			if err != ErrComponentDoNothing {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					goutils.Err(err))

				return nil, err
			}

			isComponentDoNothing = true
		}

		if !isComponentDoNothing {
			gameProp.OnCallEnd(curComponent, cd, gp, pr)

			err := curComponent.EachSymbols(gameProp, pr, gp, currng, ps, stake, prs, cd)
			if err != nil {
				goutils.Error("BasicGameMod.OnPlay:EachSymbols",
					goutils.Err(err))

				return nil, err
			}
		} else if gAllowFullComponentHistory {
			gp.HistoryComponentsEx = append(gp.HistoryComponentsEx, curComponent.GetName())
		}

		if !gIsReleaseMode {
			if goutils.IndexOfStringSlice(gameProp.Pool.Config.MapBetConfigs[int(stake.CashBet/stake.CoinBet)].ForceEndings, nextComponentName, 0) >= 0 {
				// nextComponentName = ""
				break
			}
		}

		if nextComponentName == "" {
			nc, err := gameProp.procRespinBeforeStepEnding(pr, gp)
			if err != nil {
				goutils.Error("BasicGameMod.OnPlay:procRespinBeforeStepEnding",
					goutils.Err(err))

				return nil, err
			}

			if nc == "" {
				break
			}

			// 这里可能会由于respin循环回滚导致异常,循环回滚是为了处理前端渲染的问题,待查
			// 后续可以考虑去掉循环回滚
			if !gameProp.IsRespin(nc) && gameProp.callStack.IsInCurCallStack(nc) {
				goutils.Error("BasicGameMod.OnPlay:procRespinBeforeStepEnding:IsInCurCallStack",
					slog.String("component", nc),
					goutils.Err(ErrInvalidComponentConfig))

				break
			}

			nextComponentName = nc
		}

		if gameProp.IsRespin(nextComponentName) {
			if !gameProp.IsEndingRespin(nextComponentName) {
				gameProp.onTriggerRespin(nextComponentName)
				gp.NextStepFirstComponent = nextComponentName

				pr.IsFinish = false

				break
			} else {
				cr, isok := gameProp.Components.MapComponents[nextComponentName]
				if isok {
					cd := gameProp.GetGlobalComponentData(cr)
					nc, err := cr.ProcRespinOnStepEnd(gameProp, pr, gp, cd, true)
					if err != nil {
						goutils.Error("BasicGameMod.OnPlay:IsRespin:ProcRespinOnStepEnd",
							slog.String("respin", nextComponentName),
							goutils.Err(err))

						return nil, err
					}

					if nc == "" {
						break
					}

					nextComponentName = nc
				}
			}
		}

		c, isok := components.MapComponents[nextComponentName]
		if !isok {
			goutils.Error("BasicGameMod.OnPlay:MapComponents",
				slog.String("nextComponentName", nextComponentName),
				goutils.Err(ErrInvalidComponentName))

			return nil, ErrInvalidComponentName
		}

		curComponent = c

		curComponentNum := gameProp.callStack.GetComponentNum()
		if curComponentNum >= MaxComponentNumInStep {
			goutils.Error("BasicGameMod.OnPlay",
				slog.Int("components", curComponentNum),
				goutils.Err(ErrTooManyComponentsInStep))

			return nil, ErrTooManyComponentsInStep
		}
	}

	gameProp.BuildGameParam(gp)

	err := gameProp.callStack.Each(gameProp, func(tag string, gameProp *GameProperty, ic IComponent, cd IComponentData) error {
		gp.AddComponentData(tag, cd)

		return nil
	})
	if err != nil {
		goutils.Error("BasicGameMod.OnPlay:gameProp.callStack.Each",
			goutils.Err(err))

		return nil, err
	}

	gameProp.ProcRespin(pr, gp)

	gameProp.onStepEnd(curBetMode, gp, pr, prs)

	if gAllowStats2 && pr.IsFinish {
		totalwins := int64(pr.CoinWin)

		for _, cpr := range prs {
			totalwins += int64(cpr.CoinWin)
		}

		rngs := sgc7plugin.GetRngs(plugin)
		gameProp.stats2Cache.ProcStatsOnEnding(totalwins, rngs)

		components.Stats2.PushCache(gameProp.stats2Cache)

		if components.RngLib != nil {
			rets := make([]*sgc7game.PlayResult, len(prs)+1)
			copy(rets, prs)
			rets[len(prs)] = pr

			rngname := components.RngLib.onResults(rets)
			if rngname != "" {
				components.Stats2.PushRNGs(rngname, rngs)
			}
		}
	}

	return pr, nil
}

// ResetConfig
func (bgm *BasicGameMod) ResetConfig(cfg *Config) {
	bgm.Pool.Config = cfg
}

// OnAsciiGame - outpur to asciigame
func (bgm *BasicGameMod) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) error {
	return gameProp.callStack.Each(gameProp, func(tag string, gameProp *GameProperty, ic IComponent, cd IComponentData) error {
		ic.OnAsciiGame(gameProp, pr, lst, gameProp.Pool.MapSymbolColor, cd)

		return nil
	})
}

// OnNewGame -
func (bgm *BasicGameMod) OnNewGame(gameProp *GameProperty, stake *sgc7game.Stake, curPlugin sgc7plugin.IPlugin) error {
	if gAllowStats2 {
		gameProp.Components.Stats2.PushBet(int(stake.CashBet / stake.CoinBet))

		gameProp.stats2Cache = stats2.NewCache(int(stake.CashBet / stake.CoinBet))
	}

	gameProp.OnNewGame(stake, curPlugin)

	return nil
}

// OnNewStep -
func (bgm *BasicGameMod) OnNewStep(gameProp *GameProperty, stake *sgc7game.Stake) error {
	gameProp.OnNewStep()

	return nil
}

// NewBasicGameMod2 - new BaseGame
func NewBasicGameMod2(pool *GamePropertyPool, mgrComponent *ComponentMgr) (*BasicGameMod, error) {
	bgm := &BasicGameMod{
		BasicGameMod:  sgc7game.NewBasicGameMod(BasicGameModName, pool.Config.Width, pool.Config.Height),
		Pool:          pool,
		MapComponents: make(map[int]*ComponentList),
	}

	for bet, betcfg := range pool.Config.MapBetConfigs {
		components := NewComponentList()

		for _, v := range betcfg.Components {
			c := mgrComponent.NewComponent(v)

			err := c.InitEx(betcfg.mapConfig[v.Name], pool)
			if err != nil {
				goutils.Error("NewBasicGameMod2:Init",
					goutils.Err(err))

				return nil, err
			}

			components.AddComponent(v.Name, c)
		}

		bgm.MapComponents[bet] = components
		pool.onAddComponentList(bet, components)
	}

	return bgm, nil
}
