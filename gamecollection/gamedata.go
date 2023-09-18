package gamecollection

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/grpcserv"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
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
				zap.Error(err))

			return nil, err
		}
	}

	plugin := gameD.Game.NewPlugin()
	defer gameD.Game.FreePlugin(plugin)

	ProcCheat(plugin, req.Cheat)

	stake := sgc7pbutils.BuildStake(req.Stake)
	err := gameD.Game.CheckStake(stake)
	if err != nil {
		goutils.Error("GameData.Play:CheckStake",
			goutils.JSON("stake", stake),
			zap.Error(err))

		return nil, err
	}

	results := []*sgc7game.PlayResult{}
	gameData := gameD.Game.NewGameData()

	cmd := req.Command

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := gameD.Game.Play(plugin, cmd, req.ClientParams, ips, stake, results, gameData)
		if err != nil {
			goutils.Error("GameData.Play:Play",
				zap.Int("results", len(results)),
				zap.Error(err))

			return nil, err
		}

		if pr == nil {
			break
		}

		results = append(results, pr)
		if pr.IsFinish {
			break
		}

		if pr.IsWait {
			break
		}

		if len(pr.NextCmds) > 0 {
			cmd = pr.NextCmds[0]
		} else {
			cmd = ""
		}
	}

	pr := &sgc7pb.ReplyPlay{
		RandomNumbers: sgc7pbutils.BuildPBRngs(plugin.GetUsedRngs()),
	}

	ps, err := gameD.Service.BuildPBPlayerState(ips)
	if err != nil {
		goutils.Error("GameData.Play:BuildPlayerState",
			zap.Error(err))

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
		Service:  NewService(),
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
