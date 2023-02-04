package asciigame

import (
	"fmt"
	"os"

	"devt.de/krotik/common/termutil/getch"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

// if return true, then break
type FuncOnGetChar func(k getch.KeyCode) bool

func getchar(onchar FuncOnGetChar) {
	err := getch.Start()
	if err != nil {
		goutils.Error("getchar:Start",
			zap.Error(err))

		return
	}
	defer getch.Stop()

	for {
		e, err := getch.Getch()
		if err != nil {
			goutils.Error("getchar:Start",
				zap.Error(err))

			return
		}

		if onchar(e.Code) {
			return
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
		getchar(func(c getch.KeyCode) bool {
			if c == getch.KeyS {
				return true
			}

			if c == getch.KeyQ {
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
			getchar(func(c getch.KeyCode) bool {
				if c == getch.KeyN {
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
