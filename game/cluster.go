package sgc7game

import goutils "github.com/zhs007/goutils"

// CalcClusterResult - cluster
func CalcClusterResult(scene *GameScene, pt *PayTables, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol) ([]*Result, error) {
	results := []*Result{}

	scene0 := scene.Clone()

	for x, arr := range scene.Arr {
		for y := range arr {
			if scene0.Arr[x][y] >= 0 && isValidSymbol(scene0.Arr[x][y]) {
				cr := calcClusterResult(scene0, x, y, getSymbol(scene0.Arr[x][y]), pt, bet, isSameSymbol)
				// if err != nil {
				// 	goutils.Error("sgc7game.CalcClusterResult:calcClusterResult",
				// 		slog.Int("x", x),
				// 		slog.Int("y", y),
				// 		slog.Any("scene", scene),
				// 		slog.Any("scene0", scene0),
				// 		goutils.Err(err))

				// 	return nil, err
				// }

				if cr != nil {
					results = append(results, cr)

					for i := 0; i < len(cr.Pos)/2; i++ {
						if !isWild(scene0.Arr[cr.Pos[i*2]][cr.Pos[i*2+1]]) {
							scene0.Arr[cr.Pos[i*2]][cr.Pos[i*2+1]] = -1
						}
					}
				}
			}
		}
	}

	if len(results) > 0 {
		return results, nil
	}

	return nil, nil
}

// calcClusterSymbol - cluster
func calcClusterSymbol(scene *GameScene, x, y int, symbol int, pos []int,
	isSameSymbol FuncIsSameSymbol) []int {

	if goutils.IndexOfInt2Slice(pos, x, y, 0) >= 0 {
		return pos
	}

	if isSameSymbol(scene.Arr[x][y], symbol) {
		pos = append(pos, x, y)

		if x > 0 {
			// if y > 0 {
			// 	pos = calcClusterSymbol(scene, x-1, y-1, symbol, pos, isSameSymbol)
			// }

			pos = calcClusterSymbol(scene, x-1, y, symbol, pos, isSameSymbol)

			// if y < scene.Height-1 {
			// 	pos = calcClusterSymbol(scene, x-1, y+1, symbol, pos, isSameSymbol)
			// }
		}

		if y > 0 {
			pos = calcClusterSymbol(scene, x, y-1, symbol, pos, isSameSymbol)
		}

		if y < scene.Height-1 {
			pos = calcClusterSymbol(scene, x, y+1, symbol, pos, isSameSymbol)
		}

		if x < scene.Width-1 {
			// if y > 0 {
			// 	pos = calcClusterSymbol(scene, x+1, y-1, symbol, pos, isSameSymbol)
			// }

			pos = calcClusterSymbol(scene, x+1, y, symbol, pos, isSameSymbol)

			// if y < scene.Height-1 {
			// 	pos = calcClusterSymbol(scene, x+1, y+1, symbol, pos, isSameSymbol)
			// }
		}
	}

	return pos
}

// calcClusterResult - cluster
func calcClusterResult(scene *GameScene, x, y int, symbol int, pt *PayTables, bet int,
	isSameSymbol FuncIsSameSymbol) *Result {

	pos := calcClusterSymbol(scene, x, y, symbol, []int{}, isSameSymbol)

	nums := len(pos) / 2

	if nums > len(pt.MapPay[symbol]) {
		nums = len(pt.MapPay[symbol])
	}

	if pt.MapPay[symbol][nums-1] > 0 {
		r := &Result{
			Symbol:     symbol,
			Type:       RTCluster,
			LineIndex:  -1,
			Pos:        pos,
			SymbolNums: nums,
			Mul:        pt.MapPay[symbol][nums-1],
			CoinWin:    pt.MapPay[symbol][nums-1],
			CashWin:    pt.MapPay[symbol][nums-1] * bet,
		}

		return r
	}

	return nil
}
