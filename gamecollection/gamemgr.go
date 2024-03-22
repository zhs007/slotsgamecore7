package gamecollection

import (
	"log/slog"
	"sync"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
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
				slog.String("gameCode", gameCode),
				slog.String("hash", hash))

			return nil
		}

		// goutils.Info("GameMgr.InitGame",
		// 	slog.String("data", string(data)))

		gameD1, err := NewGameDataWithHash(gameCode, data, hash)
		if err != nil {
			goutils.Error("GameMgr.InitGame:NewGameDataWithHash",
				goutils.Err(err))

			return err
		}

		mgr.MapGames[gameCode] = gameD1

		goutils.Info("GameMgr.InitGame:OK!",
			slog.String("gameCode", gameCode))

		return nil
	}

	// goutils.Info("GameMgr.InitGame",
	// 	slog.String("data", string(data)))

	gameD1, err := NewGameData(gameCode, data)
	if err != nil {
		goutils.Error("GameMgr.InitGame:NewGameData",
			goutils.Err(err))

		return err
	}

	mgr.MapGames[gameCode] = gameD1

	goutils.Info("GameMgr.InitGame:OK!",
		slog.String("gameCode", gameCode))

	return nil
}

func (mgr *GameMgr) GetGameConfig(gameCode string) (*sgc7game.Config, error) {
	mgr.Lock()
	defer mgr.Unlock()

	gameD, isok := mgr.MapGames[gameCode]
	if !isok {
		goutils.Error("GameMgr.GetGameConfig",
			slog.String("gameCode", gameCode),
			slog.Int("game number", len(mgr.MapGames)),
			goutils.Err(ErrInvalidGameCode))

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
			slog.String("gameCode", gameCode),
			slog.Bool("gameD", gameD != nil),
			slog.Bool("gameD.Game", gameD.Game != nil),
			slog.Bool("gameD.Service", gameD.Service != nil),
			slog.Int("game number", len(mgr.MapGames)),
			goutils.Err(ErrInvalidGameCode))

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
			slog.String("gameCode", gameCode),
			goutils.Err(ErrInvalidGameCode))

		return nil, ErrInvalidGameCode
	}

	reply, err := gameD.Play(req)
	if err != nil {
		goutils.Error("GameMgr.PlayGame",
			slog.String("gameCode", gameCode),
			goutils.Err(err))

		return nil, err
	}

	return reply, nil
}

func NewGameMgr() *GameMgr {
	return &GameMgr{
		MapGames: make(map[string]*GameData),
	}
}
