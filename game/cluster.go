package sgc7game

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
			if isValidSymbol(scene0.Arr[x][y]) {
				cr := calcClusterResult(scene0, x, y, getSymbol(scene0.Arr[x][y]), pt, bet, isSameSymbol)
				// if err != nil {
				// 	sgc7utils.Error("sgc7game.CalcClusterResult:calcClusterResult",
				// 		zap.Int("x", x),
				// 		zap.Int("y", y),
				// 		sgc7utils.JSON("scene", scene),
				// 		sgc7utils.JSON("scene0", scene0),
				// 		zap.Error(err))

				// 	return nil, err
				// }

				if cr != nil {
					results = append(results, cr)
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

	if isSameSymbol(scene.Arr[x][y], symbol) {
		pos = append(pos, x, y)

		if x > 0 {
			if y > 0 {
				pos = calcClusterSymbol(scene, x-1, y-1, symbol, pos, isSameSymbol)
			}

			pos = calcClusterSymbol(scene, x-1, y, symbol, pos, isSameSymbol)

			if y < scene.Height-1 {
				pos = calcClusterSymbol(scene, x-1, y+1, symbol, pos, isSameSymbol)
			}
		}

		if y > 0 {
			pos = calcClusterSymbol(scene, x, y-1, symbol, pos, isSameSymbol)
		}

		if y < scene.Height-1 {
			pos = calcClusterSymbol(scene, x, y+1, symbol, pos, isSameSymbol)
		}

		if x < scene.Width-1 {
			if y > 0 {
				pos = calcClusterSymbol(scene, x+1, y-1, symbol, pos, isSameSymbol)
			}

			pos = calcClusterSymbol(scene, x+1, y, symbol, pos, isSameSymbol)

			if y < scene.Height-1 {
				pos = calcClusterSymbol(scene, x+1, y+1, symbol, pos, isSameSymbol)
			}
		}
	}

	return pos
}

// calcClusterResult - cluster
func calcClusterResult(scene *GameScene, x, y int, symbol int, pt *PayTables, bet int,
	isSameSymbol FuncIsSameSymbol) *Result {

	pos := calcClusterSymbol(scene, x, y, symbol, []int{}, isSameSymbol)

	nums := len(pos) / 2

	if pt.MapPay[symbol][nums] > 0 {
		r := &Result{
			Symbol:     symbol,
			Type:       RTCluster,
			LineIndex:  -1,
			Pos:        pos,
			SymbolNums: nums,
		}

		return r
	}

	return nil
}