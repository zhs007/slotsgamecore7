package gamecollection

import (
	"sync"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
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
			goutils.Info("GameMgr.InitGame:same hash",
				zap.String("gameCode", gameCode),
				zap.String("hash", hash))

			return nil
		}

		// goutils.Info("GameMgr.InitGame",
		// 	zap.String("data", string(data)))

		gameD1, err := NewGameDataWithHash(gameCode, data, hash)
		if err != nil {
			goutils.Error("GameMgr.InitGame:NewGameDataWithHash",
				zap.Error(err))

			return err
		}

		mgr.MapGames[gameCode] = gameD1

		goutils.Info("GameMgr.InitGame:OK!",
			zap.String("gameCode", gameCode))

		return nil
	}

	// goutils.Info("GameMgr.InitGame",
	// 	zap.String("data", string(data)))

	gameD1, err := NewGameData(gameCode, data)
	if err != nil {
		goutils.Error("GameMgr.InitGame:NewGameData",
			zap.Error(err))

		return err
	}

	mgr.MapGames[gameCode] = gameD1

	goutils.Info("GameMgr.InitGame:OK!",
		zap.String("gameCode", gameCode))

	return nil
}

func (mgr *GameMgr) GetGameConfig(gameCode string) (*sgc7game.Config, error) {
	mgr.Lock()
	defer mgr.Unlock()

	gameD, isok := mgr.MapGames[gameCode]
	if !isok {
		goutils.Error("GameMgr.GetGameConfig",
			zap.String("gameCode", gameCode),
			zap.Int("game number", len(mgr.MapGames)),
			zap.Error(ErrInvalidGameCode))

		return nil, ErrInvalidGameCode
	}

	return gameD.Game.GetConfig(), nil
}

func (mgr *GameMgr) InitializeGamePlayer(gameCode string) (*sgc7pb.PlayerState, error) {
	mgr.Lock()
	defer mgr.Unlock()

	gameD, isok := mgr.MapGames[gameCode]
	if !isok || gameD == nil || gameD.Game == nil || gameD.Service == nil {
		goutils.Error("GameMgr.InitializeGamePlayer",
			zap.String("gameCode", gameCode),
			zap.Bool("gameD", gameD != nil),
			zap.Bool("gameD.Game", gameD.Game != nil),
			zap.Bool("gameD.Service", gameD.Service != nil),
			zap.Int("game number", len(mgr.MapGames)),
			zap.Error(ErrInvalidGameCode))

		return nil, ErrInvalidGameCode
	}

	ps := gameD.Game.Initialize()

	return gameD.Service.BuildPBPlayerState(ps)
}

// PlayGame - play game
func (mgr *GameMgr) PlayGame(gameCode string, req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	mgr.Lock()
	defer mgr.Unlock()

	gameD, isok := mgr.MapGames[gameCode]
	if !isok {
		goutils.Error("GameMgr.PlayGame",
			zap.String("gameCode", gameCode),
			zap.Error(ErrInvalidGameCode))

		return nil, ErrInvalidGameCode
	}

	reply, err := gameD.Play(req)
	if err != nil {
		goutils.Error("GameMgr.PlayGame",
			zap.String("gameCode", gameCode),
			zap.Error(err))

		return nil, err
	}

	return reply, nil
}

func NewGameMgr() *GameMgr {
	return &GameMgr{
		MapGames: make(map[string]*GameData),
	}
}
