package sgc7rtp

import (
	"context"
	"math"
	"time"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/stat"
)

// FuncOnRTPTimer - on timer for rtp
type FuncOnRTPTimer func(totalnums int64, curnums int64, curtime time.Duration)

// StartRTP - start RTP
func StartRTP(game sgc7game.IGame, rtp *RTP, worknums int, spinnums int64, stake *sgc7game.Stake, numsTimer int, ontimer FuncOnRTPTimer) time.Duration {
	t1 := time.Now()

	lastnums := worknums
	ch := make(chan *RTP)
	chTimer := make(chan int)

	go func() {
		lastspinnums := spinnums

		for {
			curnums := <-chTimer

			lastspinnums -= int64(curnums)

			ontimer(spinnums, spinnums-lastspinnums, time.Since(t1))

			if lastspinnums <= 0 {
				break
			}
		}
	}()

	for i := 0; i < worknums; i++ {
		go func() {
			currtp := rtp.Clone()

			plugin := game.NewPlugin()
			defer game.FreePlugin(plugin)

			ps := sgc7game.NewBasicPlayerState("bg")
			results := []*sgc7game.PlayResult{}
			cmd := "SPIN"
			off := 0

			for i := int64(0); i < spinnums/int64(worknums); i++ {
				plugin.ClearUsedRngs()

				currtp.Bet(stake.CashBet)

				totalReturn := int64(0)
				for {
					pr, err := game.Play(plugin, cmd, "", ps, stake, results)
					if err != nil {
						sgc7utils.Error("StartRTP.Play",
							zap.Int("results", len(results)),
							zap.Error(err))

						break
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
						cmd = "SPIN"
					}
				}

				for _, v := range results {
					totalReturn += v.CashWin

					currtp.OnResult(v)
				}

				currtp.OnResults(results)

				rtp.Returns = append(rtp.Returns, float64(totalReturn/stake.CashBet))

				results = results[:0]

				off++
				if off >= numsTimer {
					chTimer <- off

					off = 0
				}
			}

			ch <- currtp
		}()
	}

	for {
		currtp := <-ch

		rtp.Add(currtp)

		lastnums--

		if lastnums <= 0 {
			break
		}
	}

	elapsed := time.Since(t1)

	rtp.Variance = stat.Variance(rtp.Returns, nil)

	return elapsed
}

// StartRTP - start RTP
func StartScaleRTPDown(game sgc7game.IGame, rtp *RTP, worknums int, spinnums int64, stake *sgc7game.Stake, numsTimer int,
	ontimer FuncOnRTPTimer, hitFrequency float64, originalRTP float64, targetRTP float64) time.Duration {

	val := int((targetRTP/originalRTP - hitFrequency) / (1 - hitFrequency) * math.MaxInt32)

	t1 := time.Now()

	lastnums := worknums
	ch := make(chan *RTP)
	chTimer := make(chan int)

	go func() {
		lastspinnums := spinnums

		for {
			curnums := <-chTimer

			lastspinnums -= int64(curnums)

			ontimer(spinnums, spinnums-lastspinnums, time.Since(t1))

			if lastspinnums <= 0 {
				break
			}
		}
	}()

	for i := 0; i < worknums; i++ {
		go func() {
			currtp := rtp.Clone()

			plugin := game.NewPlugin()
			defer game.FreePlugin(plugin)

			ps := sgc7game.NewBasicPlayerState("bg")
			results := []*sgc7game.PlayResult{}
			cmd := "SPIN"
			off := 0

			for i := int64(0); i < spinnums/int64(worknums); {
				plugin.ClearUsedRngs()

				totalReturn := int64(0)
				for {
					pr, err := game.Play(plugin, cmd, "", ps, stake, results)
					if err != nil {
						sgc7utils.Error("StartScaleRTPDown.Play",
							zap.Int("results", len(results)),
							zap.Error(err))

						break
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
						cmd = "SPIN"
					}
				}

				iswin := false
				for _, v := range results {
					if v.CoinWin > 0 {
						iswin = true

						break
					}
				}

				if iswin {
					cr, err := plugin.Random(context.Background(), math.MaxInt32)
					if err != nil {
						sgc7utils.Error("StartScaleRTPDown.Random",
							zap.Int("results", len(results)),
							zap.Error(err))
					}

					if cr > val {
						results = results[:0]

						continue
					}
				}

				currtp.Bet(stake.CashBet)

				for _, v := range results {
					totalReturn += v.CashWin

					currtp.OnResult(v)
				}

				currtp.OnResults(results)

				rtp.Returns = append(rtp.Returns, float64(totalReturn/stake.CashBet))

				results = results[:0]

				off++
				if off >= numsTimer {
					chTimer <- off

					off = 0
				}

				i++
			}

			ch <- currtp
		}()
	}

	for {
		currtp := <-ch

		rtp.Add(currtp)

		lastnums--

		if lastnums <= 0 {
			break
		}
	}

	elapsed := time.Since(t1)

	rtp.Variance = stat.Variance(rtp.Returns, nil)

	return elapsed
}
