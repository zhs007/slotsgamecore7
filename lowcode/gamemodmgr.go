package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type FuncNewGameMod func(gameProp *GameProperty, cfgGameMod *GameModConfig, mgrComponent *ComponentMgr) sgc7game.IGameMod

type GameModMgr struct {
	MapGameMod map[string]FuncNewGameMod
}

func (mgr *GameModMgr) Reg(gamemod string, funcNew FuncNewGameMod) {
	mgr.MapGameMod[gamemod] = funcNew
}

func (mgr *GameModMgr) NewGameMod(gameProp *GameProperty, cfgGameMod *GameModConfig, mgrComponent *ComponentMgr) sgc7game.IGameMod {
	funcNew, isok := mgr.MapGameMod[cfgGameMod.Type]
	if isok {
		return funcNew(gameProp, cfgGameMod, mgrComponent)
	}

	goutils.Error("GameModMgr.NewGameMod",
		zap.String("gamemod", cfgGameMod.Type),
		zap.Error(ErrInvalidGameMod))

	return nil
}

func NewGameModMgr() *GameModMgr {
	mgr := &GameModMgr{
		MapGameMod: make(map[string]FuncNewGameMod),
	}

	mgr.Reg("bg", NewBaseGame)

	return mgr
}
