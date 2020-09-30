package sgc7rtp

import (
	"time"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
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

			ps := sgc7game.NewBasicPlayerState("bg", sgc7game.NewBPSNoBoostData)
			results := []*sgc7game.PlayResult{}
			cmd := "SPIN"
			off := 0

			for i := int64(0); i < spinnums/int64(worknums); i++ {
				plugin.ClearUsedRngs()

				currtp.Bet(stake.CashBet)

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
					currtp.OnResult(v)
				}

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

	return elapsed
}
