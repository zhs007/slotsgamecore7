package asciigame

import (
	"fmt"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type Color int

const (
	Red    Color = 1
	Green  Color = 2
	Yellow Color = 3
	Blue   Color = 4
	Purple Color = 5
	Cyan   Color = 6
	White  Color = 7
)

const ColorReset = "\033[0m"
const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorYellow = "\033[33m"
const ColorBlue = "\033[34m"
const ColorPurple = "\033[35m"
const ColorCyan = "\033[36m"
const ColorWhite = "\033[37m"

func OutputScene(str string, scene *sgc7game.GameScene, mapSymbolColor *SymbolColorMap) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			fmt.Printf("%v ", mapSymbolColor.GetSymbolString(scene.Arr[x][y]))
		}

		fmt.Print("\n")
	}
}

func OutputResults(result *sgc7game.PlayResult, mapSymbolColor *SymbolColorMap) {
	if result.CoinWin == 0 {
		fmt.Print("No Wins\n")
	} else {
		fmt.Printf("%v Wins\n", FormatColorString(fmt.Sprintf("%v", result.CoinWin), Yellow))
	}

	for _, v := range result.Results {
		fmt.Printf("%vx%v Wins %v\n", mapSymbolColor.GetSymbolString(v.Symbol),
			FormatColorString(fmt.Sprintf("%v", v.SymbolNums), Yellow),
			FormatColorString(fmt.Sprintf("%v", v.CoinWin), Yellow))
	}
}

func FormatColorString(str string, color Color) string {
	if color == Red {
		return fmt.Sprint(ColorRed, str, ColorReset)
	} else if color == Green {
		return fmt.Sprint(ColorGreen, str, ColorReset)
	} else if color == Yellow {
		return fmt.Sprint(ColorYellow, str, ColorReset)
	} else if color == Blue {
		return fmt.Sprint(ColorBlue, str, ColorReset)
	} else if color == Purple {
		return fmt.Sprint(ColorPurple, str, ColorReset)
	} else if color == Cyan {
		return fmt.Sprint(ColorCyan, str, ColorReset)
	} else if color == White {
		return fmt.Sprint(ColorWhite, str, ColorReset)
	}

	return str
}
