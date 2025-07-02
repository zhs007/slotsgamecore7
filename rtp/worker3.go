package sgc7rtp

import (
	"context"
	"log/slog"
	"time"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gonum.org/v1/gonum/stat"
)

// FuncOnRTPTimer3 - on timer for rtp
type FuncOnRTPTimer3 func(totalnums int64, curnums int64, curtime time.Duration, win int64, bet int64)

// startWorker3 -
func startWorker3(game sgc7game.IGame, rtp *RTP, spinnums int64, stake *sgc7game.Stake,
	needVariance bool, limitPayout int64, ch chan *RTP) {

	go func() {
		currtp := rtp.Clone()

		plugin := game.NewPlugin()
		defer game.FreePlugin(plugin)

		ps := game.Initialize()
		results := []*sgc7game.PlayResult{}
		gameData := game.NewGameData(stake)
		if gameData == nil {
			goutils.Error("startWorker3:NewGameData",
				goutils.Err(sgc7game.ErrInvalidStake))

			ch <- currtp

			return
		}

		defer game.DeleteGameData(gameData)
		cmd := "SPIN"
		cmdparam := ""

		for i := int64(0); i < spinnums; i++ {
			game.OnBet(plugin, cmd, cmdparam, ps, stake, results, gameData)

			ps.OnOutput()

			pbsjson := ps.GetPublicJson()
			ppsjson := ps.GetPrivateJson()
			iserrturn := false

			plugin.ClearUsedRngs()

			totalReturn := int64(0)
			for {
				pr, err := game.Play(plugin, cmd, cmdparam, ps, stake, results, gameData)
				if err != nil {
					iserrturn = true

					goutils.Error("startWorker3.Play",
						slog.Int("results", len(results)),
						goutils.Err(err))

					break
				}

				if pr == nil {
					break
				}

				results = append(results, pr)
				if pr.IsFinish {

					if currtp.Stats2 != nil {
						currtp.Stats2.OnResults(stake, results)
					}

					break
				}

				if len(pr.NextCmds) > 0 {
					if len(pr.NextCmds) > 1 {
						cr, err := plugin.Random(context.Background(), len(pr.NextCmds))
						if err != nil {
							goutils.Error("startWorker3.Random",
								goutils.Err(err))

							break
						}

						cmd = pr.NextCmds[cr]

						if len(pr.NextCmdParams) > cr {
							cmdparam = pr.NextCmdParams[cr]
						} else {
							cmdparam = ""
						}
					} else {
						cmd = pr.NextCmds[0]

						if len(pr.NextCmdParams) > 0 {
							cmdparam = pr.NextCmdParams[0]
						} else {
							cmdparam = ""
						}
					}
				} else {
					cmd = "SPIN"
					cmdparam = ""
				}
			}

			if iserrturn {
				ps.SetPublicJson(pbsjson)
				ps.SetPrivateJson(ppsjson)

				results = nil

				i--
			} else {
				if limitPayout > 0 {
					cp := int64(0)
					for _, v := range results {
						if cp+v.CashWin < limitPayout {
							cp += v.CashWin
						} else {
							v.CashWin = limitPayout - cp
							cp = limitPayout
						}
					}
				}

				currtp.Bet(stake.CashBet)

				for i, v := range results {
					currtp.TotalWins += v.CashWin
					totalReturn += v.CashWin

					currtp.OnResult(stake, i, results, gameData)
				}

				currtp.OnResults(results, gameData)

				if needVariance {
					rngs := sgc7plugin.GetRngs(plugin)
					currtp.AddReturns(float64(totalReturn/stake.CashBet), rngs)
				}

				results = nil
			}
		}

		currtp.OnPlayerPoolData(ps)

		ch <- currtp
	}()
}

// StartRTP3 - start RTP
func StartRTP3(game sgc7game.IGame, rtp *RTP, worknums int, spinnums int64, stake *sgc7game.Stake, numsTimer int,
	ontimer FuncOnRTPTimer3, needVariance bool, limitPayout int64) time.Duration {

	if spinnums < 10000 {
		return StartRTP(game, rtp, 1, spinnums, stake, numsTimer, nil, needVariance, limitPayout)
	}

	zerortp := rtp.Clone()

	tasknum := 100

	t1 := time.Now()

	lastnums := tasknum
	ch := make(chan *RTP)
	perspinnum := spinnums / int64(tasknum)

	for range worknums {
		startWorker3(game, zerortp, perspinnum, stake, needVariance, limitPayout, ch)
		lastnums--
	}

	donenum := tasknum
	curspinnum := int64(0)
	for {
		currtp := <-ch

		curspinnum += perspinnum

		rtp.Add(currtp)

		donenum--
		if donenum <= 0 {
			break
		}

		ontimer(spinnums, curspinnum, time.Since(t1), rtp.TotalWins, rtp.TotalBet)

		if lastnums > 0 {
			lastnums--
			startWorker3(game, zerortp, perspinnum, stake, needVariance, limitPayout, ch)
		}
	}

	elapsed := time.Since(t1)

	if needVariance {
		rtp.Variance = stat.Variance(rtp.Returns, rtp.ReturnWeights)
		rtp.StdDev = stat.StdDev(rtp.Returns, rtp.ReturnWeights)
	}

	return elapsed
}
