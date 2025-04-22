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

// startWorker -
func startWorker(game sgc7game.IGame, rtp *RTP, spinnums int64, stake *sgc7game.Stake,
	needVariance bool, limitPayout int64, ch chan *RTP) {

	go func() {
		currtp := rtp.Clone()

		plugin := game.NewPlugin()
		defer game.FreePlugin(plugin)

		ps := game.Initialize()
		// ps := sgc7game.NewBasicPlayerState("bg")
		results := []*sgc7game.PlayResult{}
		gameData := game.NewGameData(stake)
		if gameData == nil {
			goutils.Error("startWorker:NewGameData",
				goutils.Err(sgc7game.ErrInvalidStake))

			ch <- currtp

			return
		}

		defer game.DeleteGameData(gameData)
		cmd := "SPIN"
		cmdparam := ""
		// off := 0

		for i := int64(0); i < spinnums; i++ {
			game.OnBet(plugin, cmd, cmdparam, ps, stake, results, gameData)

			pbsjson := ps.GetPublicJson()
			ppsjson := ps.GetPrivateJson()
			iserrturn := false

			plugin.ClearUsedRngs()

			// currtp.Bet(stake.CashBet)

			totalReturn := int64(0)
			for {
				pr, err := game.Play(plugin, cmd, cmdparam, ps, stake, results, gameData)
				if err != nil {
					iserrturn = true

					goutils.Error("StartRTP.Play",
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

				// if pr.IsWait {
				// 	break
				// }

				if len(pr.NextCmds) > 0 {
					if len(pr.NextCmds) > 1 {
						cr, err := plugin.Random(context.Background(), len(pr.NextCmds))
						if err != nil {
							goutils.Error("StartRTP.Random",
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

				results = nil //results[:0]

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

				results = nil //results[:0]

				// off++
				// if off >= numsTimer {
				// 	chTimer <- off

				// 	off = 0
				// }
			}
		}

		currtp.OnPlayerPoolData(ps)

		ch <- currtp
	}()
}

// StartRTP2 - start RTP
func StartRTP2(game sgc7game.IGame, rtp *RTP, worknums int, spinnums int64, stake *sgc7game.Stake, numsTimer int,
	ontimer FuncOnRTPTimer, needVariance bool, limitPayout int64) time.Duration {

	if spinnums < 100 {
		return StartRTP(game, rtp, 1, spinnums, stake, numsTimer, ontimer, needVariance, limitPayout)
	}

	zerortp := rtp.Clone()

	tasknum := 100

	t1 := time.Now()

	lastnums := tasknum
	ch := make(chan *RTP)
	// chTimer := make(chan int)

	// go func() {
	// 	lastspinnums := spinnums

	// 	for {
	// 		curnums := <-chTimer

	// 		lastspinnums -= int64(curnums)

	// 		ontimer(spinnums, spinnums-lastspinnums, time.Since(t1))

	// 		if lastspinnums <= 0 {
	// 			break
	// 		}
	// 	}
	// }()

	for i := 0; i < worknums; i++ {
		startWorker(game, zerortp, spinnums/int64(tasknum), stake, needVariance, limitPayout, ch)
	}

	// lastspinnums := 0
	curspinnum := int64(0)
	for {
		currtp := <-ch

		curspinnum += spinnums / int64(tasknum)

		rtp.Add(currtp)

		lastnums--

		if lastnums <= 0 {
			break
		}

		ontimer(spinnums, curspinnum, time.Since(t1))

		if lastnums >= worknums {
			startWorker(game, zerortp, spinnums/int64(tasknum), stake, needVariance, limitPayout, ch)
		}
	}

	elapsed := time.Since(t1)

	if needVariance {
		rtp.Variance = stat.Variance(rtp.Returns, rtp.ReturnWeights)
		rtp.StdDev = stat.StdDev(rtp.Returns, rtp.ReturnWeights)
	}

	return elapsed
}
