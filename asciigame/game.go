package asciigame

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

func readStdin(out chan byte, in chan bool) {
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
			out <- b[0]
		}
	}
}

// if return true, then break
type FuncOnGetChar func(c byte) bool

func getchar(onchar FuncOnGetChar) {
	chanOutput := make(chan byte)
	chanEnd := make(chan bool)

	go readStdin(chanOutput, chanEnd)
	for {
		c := <-chanOutput

		if onchar(c) {
			chanEnd <- true

			break
		}
	}
}

type FuncOnResult func(*sgc7game.PlayResult)

func StartGame(game sgc7game.IGame, stake *sgc7game.Stake, onResult FuncOnResult) error {
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
		isend := false
		getchar(func(c byte) bool {
			if c == 's' || c == 'S' {
				return true
			}

			if c == 'q' || c == 'Q' {
				isend = true

				return true
			}

			return false
		})
		if isend {
			goto end
		}

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
			getchar(func(c byte) bool {
				if c == 'n' || c == 'N' {
					return true
				}

				return false
			})

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
