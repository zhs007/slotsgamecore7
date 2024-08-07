package sgc7game

// CalcAdjacentPay - AdjacentPay
func CalcAdjacentPay(scene *GameScene, pt *PayTables, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol) ([]*Result, error) {
	results := []*Result{}

	scene0 := scene.Clone()

	for x, arr := range scene.Arr {
		for y := range arr {
			if scene0.Arr[x][y] >= 0 && isValidSymbol(scene0.Arr[x][y]) {
				crx := calcAdjacentPayWithX(scene0, x, y, getSymbol(scene0.Arr[x][y]), pt, bet, isSameSymbol)
				cry := calcAdjacentPayWithY(scene0, x, y, getSymbol(scene0.Arr[x][y]), pt, bet, isSameSymbol)

				if crx != nil {
					results = append(results, crx)

					for i := 0; i < len(crx.Pos)/2; i++ {
						if !isWild(scene0.Arr[crx.Pos[i*2]][crx.Pos[i*2+1]]) {
							scene0.Arr[crx.Pos[i*2]][crx.Pos[i*2+1]] = -1
						}
					}
				}

				if cry != nil {
					results = append(results, cry)

					for i := 0; i < len(cry.Pos)/2; i++ {
						if !isWild(scene0.Arr[cry.Pos[i*2]][cry.Pos[i*2+1]]) {
							scene0.Arr[cry.Pos[i*2]][cry.Pos[i*2+1]] = -1
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

// calcAdjacentPayWithX - AdjacentPay
func calcAdjacentPayWithX(scene *GameScene, x, y int, symbol int, pt *PayTables, bet int, isSameSymbol FuncIsSameSymbol) *Result {
	pos := []int{x, y}

	for tx := 1; x+tx < scene.Width; tx++ {
		if isSameSymbol(scene.Arr[x+tx][y], symbol) {
			pos = append(pos, x+tx, y)
		} else {
			break
		}
	}

	nums := len(pos) / 2

	if nums > len(pt.MapPay[symbol]) {
		nums = len(pt.MapPay[symbol])
	}

	if pt.MapPay[symbol][nums-1] > 0 {
		r := &Result{
			Symbol:     symbol,
			Type:       RTAdjacentPay,
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

// calcAdjacentPayWithY - AdjacentPay
func calcAdjacentPayWithY(scene *GameScene, x, y int, symbol int, pt *PayTables, bet int, isSameSymbol FuncIsSameSymbol) *Result {
	pos := []int{x, y}

	for ty := 1; y+ty < scene.Height; ty++ {
		if isSameSymbol(scene.Arr[x][y+ty], symbol) {
			pos = append(pos, x, y+ty)
		} else {
			break
		}
	}

	nums := len(pos) / 2

	if nums > len(pt.MapPay[symbol]) {
		nums = len(pt.MapPay[symbol])
	}

	if pt.MapPay[symbol][nums-1] > 0 {
		r := &Result{
			Symbol:     symbol,
			Type:       RTAdjacentPay,
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
