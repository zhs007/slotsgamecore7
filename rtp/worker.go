package sgc7rtp

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// StartRTP - start RTP
func StartRTP(game sgc7game.IGame, rtp *RTP, worknums int, spinnums int64, stake *sgc7game.Stake) error {
	lastnums := worknums
	ch := make(chan *RTP)

	for i := 0; i < worknums; i++ {
		go func() {
			currtp := rtp.Clone()

			plugin := game.NewPlugin()
			defer game.FreePlugin(plugin)

			ps := sgc7game.NewBasicPlayerState("bg")
			results := []*sgc7game.PlayResult{}
			cmd := "SPIN"

			for i := int64(0); i < spinnums/int64(worknums); i++ {
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

	return nil
}
