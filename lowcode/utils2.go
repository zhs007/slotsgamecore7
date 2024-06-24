package lowcode

import (
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type WeightResults struct {
	Results []*sgc7game.PlayResult
	Weight  int
}

func SpinWithSeed(game *Game, ips sgc7game.IPlayerState, seed int, stake *sgc7game.Stake) ([]*WeightResults, error) {
	lst := []*WeightResults{}

	gameData := game.NewGameData(stake)
	defer game.DeleteGameData(gameData)

	gameProp, isok := gameData.(*GameProperty)
	if !isok {
		goutils.Error("SpinWithSeed",
			goutils.Err(ErrIvalidGameData))

		return nil, ErrIvalidGameData
	}

	cmd := "SPIN"

	for {
		results := []*sgc7game.PlayResult{}
		plugin := sgc7plugin.NewPRNGPlugin()
		plugin.SetSeed(seed)

		for {
			if cmd == "" {
				cmd = "SPIN"
			}

			pr, err := game.Play(plugin, cmd, "", ips, stake, results, gameData)
			if err != nil {
				goutils.Error("SpinWithSeed:Play",
					slog.Int("results", len(results)),
					goutils.Err(err))

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

			if len(results) >= MaxStepNum {
				goutils.Error("SpinWithSeed",
					slog.Int("steps", len(results)),
					goutils.Err(ErrTooManySteps))

				return nil, ErrTooManySteps
			}
		}

		currng, isok := gameProp.rng.(*SimpleRNG)
		if !isok {
			goutils.Error("SpinWithSeed",
				goutils.Err(ErrIvalidSimpleRNG))

			return nil, ErrIvalidSimpleRNG
		}

		if !currng.IsNeedIterate() {
			lst = append(lst, &WeightResults{
				Weight:  0,
				Results: results,
			})

			break
		} else {
			if currng.curIndex > 0 {
				lst = append(lst, &WeightResults{
					Weight:  currng.weights[currng.curIndex-1],
					Results: results,
				})

				if currng.IsIterateEnding() {
					break
				}
			}
		}
	}

	return lst, nil
}
