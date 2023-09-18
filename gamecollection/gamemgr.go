package gamecollection

import (
	"sync"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type GameMgr struct {
	sync.Mutex
	MapGames map[string]*GameData
}

func (mgr *GameMgr) InitGame(gameCode string, data []byte) error {
	mgr.Lock()
	defer mgr.Unlock()

	gameD, isok := mgr.MapGames[gameCode]
	if isok {
		hash := Hash(data)

		if hash == gameD.HashCode {
			return nil
		}

		gameD1, err := NewGameDataWithHash(gameCode, data, hash)
		if err != nil {
			goutils.Error("GameMgr.InitGame:NewGameDataWithHash",
				zap.Error(err))

			return err
		}

		mgr.MapGames[gameCode] = gameD1

		return nil
	}

	gameD1, err := NewGameData(gameCode, data)
	if err != nil {
		goutils.Error("GameMgr.InitGame:NewGameData",
			zap.Error(err))

		return err
	}

	mgr.MapGames[gameCode] = gameD1

	return nil
}

func NewGameMgr() *GameMgr {
	return &GameMgr{
		MapGames: make(map[string]*GameData),
	}
}
