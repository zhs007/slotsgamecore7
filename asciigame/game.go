package asciigame

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

func readStdin(out chan string, in chan bool) {
	//no buffering
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	//no visible output
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 1)
	for {
		select {
		case <-in:
			return
		default:
			os.Stdin.Read(b)
			fmt.Printf(">>> %v: ", b)
			out <- string(b)
		}
	}
}

type FuncOnResult func(*sgc7game.PlayResult)

func StartGame(game sgc7game.IGame, stake *sgc7game.Stake, onResult FuncOnResult) error {
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	b := make([]byte, 1)

	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	cmd := "SPIN"
	ps := game.NewPlayerState()
	results := []*sgc7game.PlayResult{}

	curgamenum := 1
	balance := 10000
	totalmoney := 10000

	for {
		fmt.Print("please press S to start spin, or press Q to quit.")

		for {
			os.Stdin.Read(b)

			if b[0] == 's' || b[0] == 'S' {
				break
			}

			if b[0] == 'q' || b[0] == 'Q' {
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

			fmt.Printf("step %v. please press N to jump the next step ...", step)
			for {
				os.Stdin.Read(b)

				if b[0] == 'n' || b[0] == 'N' {
					break
				}
			}

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
