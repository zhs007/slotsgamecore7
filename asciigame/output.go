package asciigame

import (
	"fmt"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func OutputScene(str string, scene *sgc7game.GameScene, paytables *sgc7game.PayTables) {
	if len(str) > 0 {
		fmt.Printf("%v:\n", str)
	}

	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			fmt.Printf("%v ", paytables.GetStringFromInt(scene.Arr[x][y]))
		}

		fmt.Print("\n")
	}
}

func OutputResults(result *sgc7game.PlayResult, paytables *sgc7game.PayTables) {
	if result.CoinWin == 0 {
		fmt.Print("No Wins\n")
	} else {
		fmt.Printf("%v Wins\n", result.CoinWin)
	}

	for _, v := range result.Results {
		fmt.Printf("%vx%v Wins %v\n", paytables.GetStringFromInt(v.Symbol), v.SymbolNums, v.CoinWin)
	}
}
