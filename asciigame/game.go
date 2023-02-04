package asciigame

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type FuncOnResult func(*sgc7game.PlayResult)

func StartGame(game sgc7game.IGame, stake *sgc7game.Stake, onResult FuncOnResult) error {
	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	cmd := "SPIN"
	ps := game.NewPlayerState()
	results := []*sgc7game.PlayResult{}

	curgamenum := 1
	balance := 10000
	totalmoney := 10000

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("please press S to start spin, or press Q to quit.")

		for {
			b, err := reader.ReadByte()
			if err != nil {
				goutils.Error("StartGame.ReadByte",
					zap.Error(err))

				return err
			}

			if b == 's' || b == 'S' {
				break
			}

			if b == 'q' || b == 'Q' {
				goto end
			}
		}

		step := 1
		fmt.Printf("#%v spin start -->\n", curgamenum)
		balance -= int(stake.CoinBet)
		fmt.Printf("bet %v, balance %v\n", stake.CoinBet, balance)

		for {
			pr, err := game.Play(plugin, cmd, "", ps, stake, results)
			if err != nil {
				goutils.Error("StartGame.Play",
					zap.Int("results", len(results)),
					zap.Error(err))

				break
			}

			if pr == nil {
				break
			}

			balance += pr.CoinWin
			results = append(results, pr)

			onResult(pr)

			if pr.IsFinish {
				break
			}

			fmt.Printf("step %v. please press SPACE to jump the next step ...", step)

			step++

			if pr.IsWait {
				break
			}

			if len(pr.NextCmds) > 0 {
				cmd = pr.NextCmds[0]
			} else {
				cmd = "SPIN"
			}
		}

		fmt.Printf("#%v spin end <--\n", curgamenum)

		curgamenum++
	}

end:

	fmt.Printf("you sipn %v, balance %v, win %v \n", curgamenum, balance, balance-totalmoney)

	return nil
}
