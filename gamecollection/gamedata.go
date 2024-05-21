package gamecollection

import (
	"log/slog"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/grpcserv"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

type GameData struct {
	GameCode string
	HashCode string
	Data     []byte
	Game     *lowcode.Game
	Service  grpcserv.IService
}

// Play - play game
func (gameD *GameData) Play(req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	ips := gameD.Game.NewPlayerState()
	if req.PlayerState != nil {
		err := gameD.Service.BuildPlayerStateFromPB(ips, req.PlayerState)
		if err != nil {
			goutils.Error("GameData.Play:BuildPlayerStateFromPB",
				goutils.Err(err))

			return nil, err
		}
	}

	plugin := gameD.Game.NewPlugin()
	defer gameD.Game.FreePlugin(plugin)

	stake := sgc7pbutils.BuildStake(req.Stake)

	results, err := lowcode.Spin(gameD.Game, ips, plugin, stake, req.Command, req.ClientParams, req.Cheat)
	if err != nil {
		goutils.Error("GameData.Play:Spin",
			goutils.Err(err))

		return nil, err
	}

	pr := &sgc7pb.ReplyPlay{
		RandomNumbers: sgc7pbutils.BuildPBRngs(plugin.GetUsedRngs()),
	}

	ps, err := gameD.Service.BuildPBPlayerState(ips)
	if err != nil {
		goutils.Error("GameData.Play:BuildPlayerState",
			goutils.Err(err))

		return nil, err
	}

	pr.PlayerState = ps

	if len(results) > 0 {
		AddPlayResult(gameD.Service, pr, results)

		lastr := results[len(results)-1]

		pr.Finished = lastr.IsFinish
		pr.NextCommands = lastr.NextCmds
		pr.NextCommandParams = lastr.NextCmdParams
	}

	return pr, nil
}

func NewGameData(gameCode string, data []byte, funcNewRNG lowcode.FuncNewRNG, funcNewFeatureLevel lowcode.FuncNewFeatureLevel) (*GameData, error) {
	game, err := lowcode.NewGame2WithData(data, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, funcNewRNG, funcNewFeatureLevel)
	if err != nil {
		goutils.Error("NewGameData:NewGame2WithData",
			slog.String("gameCode", gameCode),
			slog.String("data", string(data)),
			goutils.Err(err))

		return nil, err
	}

	gameD := &GameData{
		GameCode: gameCode,
		Data:     data,
		Game:     game,
		Service:  NewService(),
	}

	gameD.HashCode = Hash(data)

	return gameD, nil
}

func NewGameDataWithHash(gameCode string, data []byte, hash string, funcNewRNG lowcode.FuncNewRNG, funcNewFeatureLevel lowcode.FuncNewFeatureLevel) (*GameData, error) {
	game, err := lowcode.NewGame2WithData(data, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, funcNewRNG, funcNewFeatureLevel)
	if err != nil {
		goutils.Error("NewGameDataWithHash:NewGame2WithData",
			slog.String("gameCode", gameCode),
			slog.String("data", string(data)),
			goutils.Err(err))

		return nil, err
	}

	gameD := &GameData{
		GameCode: gameCode,
		Data:     data,
		Game:     game,
		HashCode: hash,
		Service:  NewService(),
	}

	return gameD, nil
}
