package asciigame

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

var ColorNumber *color.Color
var ColorKey *color.Color
var ColorExitKey *color.Color
var ColorWin *color.Color
var ColorLose *color.Color

func OutputReverseScene(str string, scene *sgc7game.GameScene, mapSymbolColor *SymbolColorMap) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			fmt.Printf("%v ", mapSymbolColor.GetSymbolString(scene.Arr[x][scene.Height-1-y]))
		}

		fmt.Print("\n")
	}
}

func IsPosInResult(x, y int, result *sgc7game.PlayResult) bool {
	for _, v := range result.Results {
		if goutils.IndexOfInt2Slice(v.Pos, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

func OutputReverseSceneEx(str string, scene *sgc7game.GameScene, result *sgc7game.PlayResult, c *color.Color, mapSymbolColor *SymbolColorMap) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			if result != nil && IsPosInResult(x, scene.Height-1-y, result) {
				fmt.Printf("%v ", FormatColorString(mapSymbolColor.PayTables.GetStringFromInt(scene.Arr[x][scene.Height-1-y]), c))
			} else {
				fmt.Printf("%v ", mapSymbolColor.PayTables.GetStringFromInt(scene.Arr[x][scene.Height-1-y]))
			}
		}

		fmt.Print("\n")
	}
}

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

func OutputOtherScene(str string, scene *sgc7game.GameScene) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			fmt.Printf("%v ", scene.Arr[x][y])
		}

		fmt.Print("\n")
	}
}

func OutputReverseOtherScene(str string, scene *sgc7game.GameScene) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			fmt.Printf("%v ", scene.Arr[x][scene.Height-1-y])
		}

		fmt.Print("\n")
	}
}

type FuncIsResult func(int, *sgc7game.Result) bool

func OutputResults(str string, result *sgc7game.PlayResult, isResult FuncIsResult, mapSymbolColor *SymbolColorMap) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	wins := 0

	for i, v := range result.Results {
		if isResult(i, v) {
			wins += v.CashWin
		}
	}

	if wins == 0 {
		fmt.Print("No Wins\n")
	} else {
		fmt.Printf("%v Wins\n", FormatColorString(fmt.Sprintf("%v", wins), ColorNumber))
	}

	for i, v := range result.Results {
		if isResult(i, v) {
			fmt.Printf("%vx%v Wins %v\n", mapSymbolColor.GetSymbolString(v.Symbol),
				FormatColorString(fmt.Sprintf("%v", v.SymbolNums), ColorNumber),
				FormatColorString(fmt.Sprintf("%v", v.CashWin), ColorNumber))
		}
	}
}

func FormatColorString(str string, c *color.Color) string {
	if c == nil {
		return str
	}

	return c.SprintFunc()(str)
}

func init() {
	ColorNumber = color.New(color.FgHiYellow)
	ColorKey = color.New(color.FgHiGreen)
	ColorExitKey = color.New(color.FgHiRed)
	ColorWin = color.New(color.FgHiGreen)
	ColorLose = color.New(color.FgHiRed)
}
