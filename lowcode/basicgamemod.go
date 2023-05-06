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
	Pool       *GamePropertyPool
	Components *ComponentList
}

// OnPlay - on play
func (bgm *BasicGameMod) newPlayResult(prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, *GameParams) {
	gp := &GameParams{}
	gp.MapComponents = make(map[string]*anypb.Any)

	pr := &sgc7game.PlayResult{
		IsFinish:         true,
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

	if len(prs) == 0 {
		bgm.OnNewGame(gameProp)
	}

	bgm.OnNewStep(gameProp)

	pr, gp := bgm.newPlayResult(prs)

	curComponent := bgm.Components.Components[0]

	if gp.FirstComponent != "" {
		c, isok := bgm.Components.MapComponents[gp.FirstComponent]
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
				c, isok1 := bgm.Components.MapComponents[cn]
				if !isok1 {
					goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
						zap.String("cmd", cmd),
						zap.String("MapCmdComponent", cn),
						zap.Error(ErrIvalidComponentName))

					return nil, ErrIvalidComponentName
				}

				curComponent = c
			}
		}
	}

	for {
		err := curComponent.OnPlayGame(gameProp, pr, gp, plugin, cmd, param, ps, stake, prs)
		if err != nil {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
				zap.Error(err))

			return nil, err
		}

		gameProp.HistoryComponents = append(gameProp.HistoryComponents, curComponent)
		gp.HistoryComponents = append(gp.HistoryComponents, curComponent.GetName())

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

		c, isok := bgm.Components.MapComponents[nextComponentName]
		if !isok {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
				zap.String("nextComponentName", nextComponentName),
				zap.Error(ErrIvalidComponentName))

			return nil, ErrIvalidComponentName
		}

		curComponent = c
	}

	// gameProp.ProcRespin(pr, gp)

	// if pr.IsFinish && gameProp.GetVal(GamePropFGNum) > 0 {
	// 	pr.IsFinish = false
	// } else if gameProp.GetVal(GamePropTriggerFG) > 0 && gameProp.GetVal(GamePropFGNum) <= 0 {
	// 	gameProp.SetVal(GamePropTriggerFG, 0)
	// }

	gameProp.BuildGameParam(gp)

	for _, v := range gameProp.HistoryComponents {
		err := v.OnPlayGameEnd(gameProp, pr, gp, plugin, cmd, param, ps, stake, prs)
		if err != nil {
			goutils.Error("BasicGameMod.OnPlay:OnPlayGameEnd",
				zap.Error(err))

			return nil, err
		}
	}

	gameProp.ProcRespin(pr, gp)

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
func (bgm *BasicGameMod) OnNewGame(gameProp *GameProperty) error {
	gameProp.OnNewGame()

	for i, v := range bgm.Components.Components {
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
func (bgm *BasicGameMod) OnNewStep(gameProp *GameProperty) error {
	gameProp.OnNewStep()

	for i, v := range bgm.Components.Components {
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
		BasicGameMod: sgc7game.NewBasicGameMod(cfgGameMod.Type, pool.Config.Width, pool.Config.Height),
		Pool:         pool,
		Components:   NewComponentList(),
	}

	for _, v := range cfgGameMod.Components {
		c := mgrComponent.NewComponent(v)
		err := c.Init(pool.Config.GetPath(v.Config), pool)
		if err != nil {
			goutils.Error("NewBasicGameMod:Init",
				zap.Error(err))

			return nil
		}

		bgm.Components.AddComponent(v.Name, c)
		pool.onAddComponent(v.Name, c)
	}

	return bgm
}
