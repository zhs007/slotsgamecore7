package asciigame

import (
	"fmt"

	"devt.de/krotik/common/termutil/getch"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

var gKeys []getch.KeyCode

// if return true, then break
type FuncOnGetChar func(k getch.KeyCode) bool

func getchar(onchar FuncOnGetChar) error {
	err := getch.Start()
	if err != nil {
		goutils.Error("getchar:Start",
			zap.Error(err))

		return err
	}
	defer getch.Stop()

	for {
		e, err := getch.Getch()
		if err != nil {
			goutils.Error("getchar:Start",
				zap.Error(err))

			return err
		}

		if onchar(e.Code) {
			return nil
		}
	}
}

type FuncOnResult func(*sgc7game.PlayResult, []*sgc7game.PlayResult, any)

func StartGame(game sgc7game.IGame, stake *sgc7game.Stake, onResult FuncOnResult, autogametimes int, isSkipGetChar bool, isBreakAtFeature bool) error {
	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	cmd := "SPIN"
	cmdparam := ""
	ps := game.NewPlayerState()
	results := []*sgc7game.PlayResult{}

	curgamenum := 1
	balance := 10000
	totalmoney := 10000

	autotimes := autogametimes

	for {
		if autotimes <= 0 && !isSkipGetChar {
			fmt.Printf("please press %v to start spin, or press %v to spin 100 times, or press %v to spin 1000 times, or press %v to quit.\n",
				FormatColorString("S", ColorKey), FormatColorString("H", ColorKey), FormatColorString("K", ColorKey), FormatColorString("Q", ColorExitKey))

			isend := false

			err := getchar(func(c getch.KeyCode) bool {
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
			if err != nil {
				goto end
			}

			if isend {
				goto end
			}
		}

		step := 1

		fmt.Printf("%v spin start -->\n",
			FormatColorString(fmt.Sprintf("#%v", curgamenum), ColorNumber))

		balance -= int(stake.CashBet)
		spinwins := 0
		gameData := game.NewGameData()

		fmt.Printf("bet %v, balance %v\n",
			FormatColorString(fmt.Sprintf("%v", stake.CashBet), ColorNumber),
			FormatColorString(fmt.Sprintf("%v", balance), ColorNumber))

		for {
			pr, err := game.Play(plugin, cmd, cmdparam, ps, stake, results, gameData)
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

			onResult(pr, results, gameData)

			if pr.IsFinish {
				break
			}

			if isBreakAtFeature {
				autotimes = 0
			}

			fmt.Printf("balance %v , win %v \n",
				FormatColorString(fmt.Sprintf("%v", balance), ColorNumber),
				FormatColorString(fmt.Sprintf("%v", pr.CashWin), ColorNumber))

			curSelected := 0

			if autotimes <= 0 && !isSkipGetChar {
				if len(pr.NextCmds) > 1 {
					keys := []getch.KeyCode{}

					fmt.Printf("step %v. please press the Number to select the feature or press %v to quit.\n",
						FormatColorString(fmt.Sprintf("%v", step), ColorNumber),
						FormatColorString("Q", ColorExitKey))

					for i, cmd := range pr.NextCmds {
						fmt.Printf("%v. %v\n",
							FormatColorString(fmt.Sprintf("%v", i+1), ColorKey),
							FormatColorString(cmd, ColorKey))

						keys = append(keys, gKeys[i])
					}

					isend := false
					getchar(func(c getch.KeyCode) bool {
						if c == getch.KeyQ {
							isend = true

							return true
						}

						for i, ck := range keys {
							if c == ck {
								curSelected = i

								return true
							}
						}

						return false
					})

					if isend {
						goto end
					}
				} else {
					fmt.Printf("step %v. please press %v to jump to the next step or press %v to quit.\n",
						FormatColorString(fmt.Sprintf("%v", step), ColorNumber),
						FormatColorString("N", ColorKey),
						FormatColorString("Q", ColorExitKey))

					isend := false
					getchar(func(c getch.KeyCode) bool {
						if c == getch.KeyN {
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
			}

			step++

			// if pr.IsWait {
			// 	break
			// }

			if len(pr.NextCmds) > 0 {
				cmd = pr.NextCmds[curSelected]
				
				if len(pr.NextCmdParams) > 0 {
					cmdparam = pr.NextCmdParams[curSelected]
				} else{
					cmdparam = ""	
				}
			} else {
				cmd = "SPIN"
				cmdparam = ""
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

	fmt.Printf("you spin %v, balance %v, win %v \n",
		FormatColorString(fmt.Sprintf("%v", curgamenum), ColorNumber),
		FormatColorString(fmt.Sprintf("%v", balance), ColorNumber),
		FormatColorString(fmt.Sprintf("%v", balance-totalmoney), SelectColor(func() bool {
			return balance > totalmoney
		}, ColorWin, ColorLose)))

	return nil
}

func init() {
	gKeys = []getch.KeyCode{getch.Key1, getch.Key2, getch.Key3, getch.Key4, getch.Key5, getch.Key6, getch.Key7, getch.Key8, getch.Key9}
}
