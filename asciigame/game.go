package asciigame

import (
	"fmt"

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

type FuncOnResult func(*sgc7game.PlayResult, []*sgc7game.PlayResult)

func StartGame(game sgc7game.IGame, stake *sgc7game.Stake, onResult FuncOnResult) error {
	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	cmd := "SPIN"
	ps := game.NewPlayerState()
	results := []*sgc7game.PlayResult{}

	curgamenum := 1
	balance := 10000
	totalmoney := 10000

	autotimes := 0

	for {
		if autotimes <= 0 {
			fmt.Printf("please press %v to start spin, or press %v to spin 100 times, or press %v to spin 1000 times, or press %v to quit.\n",
				FormatColorString("S", ColorKey), FormatColorString("H", ColorExitKey), FormatColorString("K", ColorExitKey), FormatColorString("Q", ColorExitKey))

			isend := false

			getchar(func(c getch.KeyCode) bool {
				if c == getch.KeyS {
					return true
				}

				if c == getch.KeyH {
					autotimes = 100

					return true
				}

				if c == getch.KeyK {
					autotimes = 1000

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
		}

		step := 1

		fmt.Printf("%v spin start -->\n",
			FormatColorString(fmt.Sprintf("#%v", curgamenum), ColorNumber))

		balance -= int(stake.CashBet)
		spinwins := 0

		fmt.Printf("bet %v, balance %v\n",
			FormatColorString(fmt.Sprintf("%v", stake.CashBet), ColorNumber),
			FormatColorString(fmt.Sprintf("%v", balance), ColorNumber))

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

			balance += int(pr.CashWin)
			spinwins += int(pr.CashWin)
			results = append(results, pr)

			onResult(pr, results)

			if pr.IsFinish {
				break
			}

			fmt.Printf("balance %v , win %v \n",
				FormatColorString(fmt.Sprintf("%v", balance), ColorNumber),
				FormatColorString(fmt.Sprintf("%v", pr.CashWin), ColorNumber))

			if autotimes <= 0 {
				fmt.Printf("step %v. please press %v to jump to the next step.\n",
					FormatColorString(fmt.Sprintf("%v", step), ColorNumber),
					FormatColorString("N", ColorKey))

				getchar(func(c getch.KeyCode) bool {
					return c == getch.KeyN
				})
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

		fmt.Printf("balance %v , win %v \n",
			FormatColorString(fmt.Sprintf("%v", balance), ColorNumber),
			FormatColorString(fmt.Sprintf("%v", spinwins), ColorNumber))

		fmt.Printf("%v spin end <--\n",
			FormatColorString(fmt.Sprintf("#%v", curgamenum), ColorNumber))

		curgamenum++
		autotimes--

		results = nil
	}

end:

	fmt.Printf("you sipn %v, balance %v, win %v \n",
		FormatColorString(fmt.Sprintf("%v", curgamenum), ColorNumber),
		FormatColorString(fmt.Sprintf("%v", balance), ColorNumber),
		FormatColorString(fmt.Sprintf("%v", balance-totalmoney), SelectColor(func() bool {
			return balance > totalmoney
		}, ColorWin, ColorLose)))

	return nil
}