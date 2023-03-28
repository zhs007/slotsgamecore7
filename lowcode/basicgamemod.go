package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
)

// BasicGameMod - basic gamemod
type BasicGameMod struct {
	*sgc7game.BasicGameMod
	Pool              *GamePropertyPool
	Components        *ComponentList
	HistoryComponents []IComponent
}

// OnPlay - on play
func (bgm *BasicGameMod) newPlayResult(prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, *GameParams) {
	gp := &GameParams{}
	gp.MapComponents = make(map[string]*sgc7pb.ComponentData)

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
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, gameData interface{}) (*sgc7game.PlayResult, error) {

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

	if cmd == "SPIN" {
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
		}

		for {
			err := curComponent.OnPlayGame(gameProp, pr, gp, plugin, cmd, param, ps, stake, prs)
			if err != nil {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					zap.Error(err))

				return nil, err
			}

			bgm.HistoryComponents = append(bgm.HistoryComponents, curComponent)

			respinComponent := gameProp.GetStrVal(GamePropRespinComponent)
			if respinComponent != "" {
				pr.IsFinish = false

				break
			}

			nextComponentName := gameProp.GetStrVal(GamePropNextComponent)
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

		if pr.IsFinish && gameProp.GetVal(GamePropFGNum) > 0 {
			pr.IsFinish = false
		} else if gameProp.GetVal(GamePropTriggerFG) > 0 && gameProp.GetVal(GamePropFGNum) <= 0 {
			gameProp.SetVal(GamePropTriggerFG, 0)
		}

		return pr, nil
	}

	return nil, sgc7game.ErrInvalidCommand
}

// ResetConfig
func (bgm *BasicGameMod) ResetConfig(cfg *Config) {
	bgm.Pool.Config = cfg
}

// OnAsciiGame - outpur to asciigame
func (bgm *BasicGameMod) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) error {
	for _, v := range bgm.HistoryComponents {
		v.OnAsciiGame(gameProp, pr, lst, gameProp.Pool.MapSymbolColor)
	}

	return nil
}

// OnNewGame -
func (bgm *BasicGameMod) OnNewGame(gameProp *GameProperty) error {
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
	bgm.HistoryComponents = nil
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
		err := c.Init(v.Config, pool)
		if err != nil {
			goutils.Error("NewBasicGameMod:Init",
				zap.Error(err))

			return nil
		}

		bgm.Components.AddComponent(v.Name, c)
		pool.onAddComponent(v.Name, c)
	}

	// err := pool.InitStats()
	// if err != nil {
	// 	goutils.Error("NewBasicGameMod:InitStats",
	// 		zap.Error(err))

	// 	return nil
	// }

	return bgm
}
