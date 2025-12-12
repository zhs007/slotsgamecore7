package sgc7game

import "github.com/zhs007/goutils"

// isValidAdjacentPayResult - 这个接口只用于内部调用，为了避免重复判断中奖的，所以传入的 gamescene 是一个特殊用途的 scene，不能乱用
// gs 只有 2 种值，0 和 -1， 0 表示已经处理过中奖了，-1 表示没有处理过
// 考虑到 A W W W B，这样的局面，是 2 个中奖（分别是 A W W W 和 W W W B），所以只要有一个位置是新的就算有效中奖
func isValidAdjacentPayResult(gs *GameScene, ret *Result) bool {
	for i := 0; i < len(ret.Pos)/2; i++ {
		if gs.Arr[ret.Pos[i*2]][ret.Pos[i*2+1]] != 0 {
			return true
		}
	}

	return false
}

// CalcAdjacentPay - AdjacentPay
func CalcAdjacentPay(scene *GameScene, pt *PayTables, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol) ([]*Result, error) {

	results := []*Result{}

	scene0 := scene.Clone()
	// gsx 是一个默认全 -1 的 scene
	gsx, err := NewGameScene(scene.Width, scene.Height)
	if err != nil {
		goutils.Error("CalcAdjacentPay:NewGameScene",
			goutils.Err(err))

		return nil, err
	}

	isprocx := false

	// 先判断 x 方向，注意，scene 是 arr[x][y]，所以 x 其实是竖向的
	for x, arr := range scene.Arr {
		for y := range arr {
			if scene0.Arr[x][y] >= 0 && isValidSymbol(scene0.Arr[x][y]) {
				crx := calcAdjacentPayWithX(scene0, x, y, getSymbol(scene0.Arr[x][y]), pt, bet, isSameSymbol, isWild)

				// 这里的 crx 是一个特殊的结果，不能直接用来计算结果
				// 其实这里只是需要保证不重复判断即可
				// 考虑到 wild，所以只要有一个位置是新的就算有效中奖
				if crx != nil && isValidAdjacentPayResult(gsx, crx) {
					results = append(results, crx)

					for i := 0; i < len(crx.Pos)/2; i++ {
						// 这里也是为了不重复判断，把除了 wild 以外的设置为空
						if !isWild(scene0.Arr[crx.Pos[i*2]][crx.Pos[i*2+1]]) {
							scene0.Arr[crx.Pos[i*2]][crx.Pos[i*2+1]] = -1

							isprocx = true
						}

						// 这里将 gsx 设置为 0，就是不能重复判断
						gsx.Arr[crx.Pos[i*2]][crx.Pos[i*2+1]] = 0
					}
				}
			}
		}
	}

	if isprocx {
		// 当清空 scene0 的一部分以后才需要重新 clone，否则可以节省一个 clone
		scene0 = scene.Clone()
	}

	gsy, err := NewGameScene(scene.Width, scene.Height)
	if err != nil {
		goutils.Error("CalcAdjacentPay:NewGameScene",
			goutils.Err(err))

		return nil, err
	}

	for x, arr := range scene.Arr {
		for y := range arr {
			if scene0.Arr[x][y] >= 0 && isValidSymbol(scene0.Arr[x][y]) {
				cry := calcAdjacentPayWithY(scene0, x, y, getSymbol(scene0.Arr[x][y]), pt, bet, isSameSymbol, isWild)

				if cry != nil && isValidAdjacentPayResult(gsy, cry) {
					results = append(results, cry)

					for i := 0; i < len(cry.Pos)/2; i++ {
						if !isWild(scene0.Arr[cry.Pos[i*2]][cry.Pos[i*2+1]]) {
							scene0.Arr[cry.Pos[i*2]][cry.Pos[i*2+1]] = -1
						}

						gsy.Arr[cry.Pos[i*2]][cry.Pos[i*2+1]] = 0
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
func calcAdjacentPayWithX(scene *GameScene, x, y int, symbol int, pt *PayTables, bet int, isSameSymbol FuncIsSameSymbol, isWild FuncIsWild) *Result {
	pos := []int{x, y}

	if isWild(symbol) {
		wpos := []int{x, y}
		curs := -1
		for tx := 1; x+tx < scene.Width; tx++ {
			if scene.Arr[x+tx][y] < 0 {
				break
			}

			if curs == -1 {
				if isWild(scene.Arr[x+tx][y]) {
					wpos = append(wpos, x+tx, y)
					pos = append(pos, x+tx, y)
				} else {
					pos = append(pos, x+tx, y)

					curs = scene.Arr[x+tx][y]
				}
			} else if isSameSymbol(scene.Arr[x+tx][y], curs) {
				pos = append(pos, x+tx, y)
			} else {
				break
			}
		}

		wnums := len(wpos) / 2
		if wnums > len(pt.MapPay[symbol]) {
			wnums = len(pt.MapPay[symbol])
		}

		if curs == -1 {
			if pt.MapPay[symbol][wnums-1] > 0 {
				r := &Result{
					Symbol:     symbol,
					Type:       RTAdjacentPay,
					LineIndex:  -1,
					Pos:        wpos,
					SymbolNums: wnums,
					Mul:        pt.MapPay[symbol][wnums-1],
					CoinWin:    pt.MapPay[symbol][wnums-1],
					CashWin:    pt.MapPay[symbol][wnums-1] * bet,
				}

				return r
			}

			return nil
		}

		nums := len(pos) / 2

		if nums > len(pt.MapPay[curs]) {
			nums = len(pt.MapPay[curs])
		}

		if pt.MapPay[curs][nums-1] > 0 && pt.MapPay[symbol][wnums-1] > 0 {
			// if pt.MapPay[curs][nums-1] > pt.MapPay[symbol][wnums-1] {
			r := &Result{
				Symbol:     curs,
				Type:       RTAdjacentPay,
				LineIndex:  -1,
				Pos:        pos,
				SymbolNums: nums,
				Mul:        pt.MapPay[curs][nums-1],
				CoinWin:    pt.MapPay[curs][nums-1],
				CashWin:    pt.MapPay[curs][nums-1] * bet,
			}

			return r
			// }

			// r := &Result{
			// 	Symbol:     symbol,
			// 	Type:       RTAdjacentPay,
			// 	LineIndex:  -1,
			// 	Pos:        wpos,
			// 	SymbolNums: wnums,
			// 	Mul:        pt.MapPay[symbol][wnums-1],
			// 	CoinWin:    pt.MapPay[symbol][wnums-1],
			// 	CashWin:    pt.MapPay[symbol][wnums-1] * bet,
			// }

			// return r
		}

		if pt.MapPay[symbol][wnums-1] > 0 {
			r := &Result{
				Symbol:     symbol,
				Type:       RTAdjacentPay,
				LineIndex:  -1,
				Pos:        wpos,
				SymbolNums: wnums,
				Mul:        pt.MapPay[symbol][wnums-1],
				CoinWin:    pt.MapPay[symbol][wnums-1],
				CashWin:    pt.MapPay[symbol][wnums-1] * bet,
			}

			return r
		}

		if pt.MapPay[curs][nums-1] > 0 {
			r := &Result{
				Symbol:     curs,
				Type:       RTAdjacentPay,
				LineIndex:  -1,
				Pos:        pos,
				SymbolNums: nums,
				Mul:        pt.MapPay[curs][nums-1],
				CoinWin:    pt.MapPay[curs][nums-1],
				CashWin:    pt.MapPay[curs][nums-1] * bet,
			}

			return r
		}

		return nil
	}

	for tx := 1; x+tx < scene.Width; tx++ {
		if scene.Arr[x+tx][y] < 0 {
			break
		}

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
func calcAdjacentPayWithY(scene *GameScene, x, y int, symbol int, pt *PayTables, bet int, isSameSymbol FuncIsSameSymbol, isWild FuncIsWild) *Result {
	pos := []int{x, y}

	if isWild(symbol) {
		wpos := []int{x, y}
		curs := -1

		for ty := 1; y+ty < scene.Height; ty++ {
			if scene.Arr[x][y+ty] < 0 {
				break
			}

			if curs == -1 {
				if isWild(scene.Arr[x][y+ty]) {
					wpos = append(wpos, x, y+ty)
					pos = append(pos, x, y+ty)
				} else {
					pos = append(pos, x, y+ty)

					curs = scene.Arr[x][y+ty]
				}
			} else if isSameSymbol(scene.Arr[x][y+ty], curs) {
				pos = append(pos, x, y+ty)
			} else {
				break
			}
		}

		wnums := len(wpos) / 2
		if wnums > len(pt.MapPay[symbol]) {
			wnums = len(pt.MapPay[symbol])
		}

		if curs == -1 {
			if pt.MapPay[symbol][wnums-1] > 0 {
				r := &Result{
					Symbol:     symbol,
					Type:       RTAdjacentPay,
					LineIndex:  -1,
					Pos:        wpos,
					SymbolNums: wnums,
					Mul:        pt.MapPay[symbol][wnums-1],
					CoinWin:    pt.MapPay[symbol][wnums-1],
					CashWin:    pt.MapPay[symbol][wnums-1] * bet,
				}

				return r
			}

			return nil
		}

		nums := len(pos) / 2

		if nums > len(pt.MapPay[curs]) {
			nums = len(pt.MapPay[curs])
		}

		if pt.MapPay[curs][nums-1] > 0 && pt.MapPay[symbol][wnums-1] > 0 {
			// if pt.MapPay[curs][nums-1] > pt.MapPay[symbol][wnums-1] {
			r := &Result{
				Symbol:     curs,
				Type:       RTAdjacentPay,
				LineIndex:  -1,
				Pos:        pos,
				SymbolNums: nums,
				Mul:        pt.MapPay[curs][nums-1],
				CoinWin:    pt.MapPay[curs][nums-1],
				CashWin:    pt.MapPay[curs][nums-1] * bet,
			}

			return r
			// }

			// r := &Result{
			// 	Symbol:     symbol,
			// 	Type:       RTAdjacentPay,
			// 	LineIndex:  -1,
			// 	Pos:        wpos,
			// 	SymbolNums: wnums,
			// 	Mul:        pt.MapPay[symbol][wnums-1],
			// 	CoinWin:    pt.MapPay[symbol][wnums-1],
			// 	CashWin:    pt.MapPay[symbol][wnums-1] * bet,
			// }

			// return r
		}

		if pt.MapPay[symbol][wnums-1] > 0 {
			r := &Result{
				Symbol:     symbol,
				Type:       RTAdjacentPay,
				LineIndex:  -1,
				Pos:        wpos,
				SymbolNums: wnums,
				Mul:        pt.MapPay[symbol][wnums-1],
				CoinWin:    pt.MapPay[symbol][wnums-1],
				CashWin:    pt.MapPay[symbol][wnums-1] * bet,
			}

			return r
		}

		if pt.MapPay[curs][nums-1] > 0 {
			r := &Result{
				Symbol:     curs,
				Type:       RTAdjacentPay,
				LineIndex:  -1,
				Pos:        pos,
				SymbolNums: nums,
				Mul:        pt.MapPay[curs][nums-1],
				CoinWin:    pt.MapPay[curs][nums-1],
				CashWin:    pt.MapPay[curs][nums-1] * bet,
			}

			return r
		}

		return nil
	}

	for ty := 1; y+ty < scene.Height; ty++ {
		if scene.Arr[x][y+ty] < 0 {
			break
		}

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
