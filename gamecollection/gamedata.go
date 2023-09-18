package gamecollection

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

type GameData struct {
	GameCode string
	HashCode string
	Data     []byte
	Game     *lowcode.Game
}

func NewGameData(gameCode string, data []byte) (*GameData, error) {
	game, err := lowcode.NewGame2WithData(data, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	})
	if err != nil {
		goutils.Error("NewGameData:NewGame2WithData",
			zap.String("gameCode", gameCode),
			zap.String("data", string(data)),
			zap.Error(err))

		return nil, err
	}

	gameD := &GameData{
		GameCode: gameCode,
		Data:     data,
		Game:     game,
	}

	gameD.HashCode = Hash(data)

	return gameD, nil
}

func NewGameDataWithHash(gameCode string, data []byte, hash string) (*GameData, error) {
	game, err := lowcode.NewGame2WithData(data, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	})
	if err != nil {
		goutils.Error("NewGameDataWithHash:NewGame2WithData",
			zap.String("gameCode", gameCode),
			zap.String("data", string(data)),
			zap.Error(err))

		return nil, err
	}

	gameD := &GameData{
		GameCode: gameCode,
		Data:     data,
		Game:     game,
		HashCode: hash,
	}

	return gameD, nil
}
