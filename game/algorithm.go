package sgc7game

// FuncIsSymbol - is symbol
type FuncIsSymbol func(cursymbol int) bool

// FuncIsScatter - cursymbol == scatter
type FuncIsScatter func(scatter int, cursymbol int) bool

// FuncIsWild - cursymbol == wild
type FuncIsWild func(cursymbol int) bool

// FuncIsSameSymbol - cursymbol == startsymbol
type FuncIsSameSymbol func(cursymbol int, startsymbol int) bool

// FuncIsSameSymbolEx - cursymbol == startsymbol
type FuncIsSameSymbolEx func(cursymbol int, startsymbol int, scene *GameScene, x, y int) bool

// FuncIsValidSymbol - is it a valid symbol?
type FuncIsValidSymbol func(cursymbol int) bool

// FuncIsValidSymbolEx - is it a valid symbol?
type FuncIsValidSymbolEx func(cursymbol int, scene *GameScene, x, y int) bool

// FuncCountSymbolInReel - count symbol nums in a reel
type FuncCountSymbolInReel func(cursymbol int, scene *GameScene, x int) int

// CalcScatter - calc scatter
func CalcScatter(scene *GameScene, pt *PayTables, scatter int, bet int, coins int,
	isScatter FuncIsScatter) *Result {

	nums := 0
	pos := make([]int, 0, len(scene.Arr)*len(scene.Arr[0])*2)
	for x := 0; x < len(scene.Arr); x++ {
		for y := 0; y < len(scene.Arr[x]); y++ {
			if isScatter(scatter, scene.Arr[x][y]) {
				nums++

				pos = append(pos, x, y)
			}
		}
	}

	if nums > len(scene.Arr) {
		nums = len(scene.Arr)
	}

	if nums > 0 && pt.MapPay[scatter][nums-1] > 0 {
		r := &Result{
			Symbol:    scatter,
			Type:      RTScatter,
			LineIndex: -1,
			Mul:       pt.MapPay[scatter][nums-1],
			CoinWin:   pt.MapPay[scatter][nums-1] * coins,
			CashWin:   pt.MapPay[scatter][nums-1] * coins * bet,
			Pos:       pos,
		}

		return r
	}

	return nil
}

// CalcLine - calc line
func CalcLine(scene *GameScene, pt *PayTables, ld []int, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol) *Result {

	s0 := scene.Arr[0][ld[0]]
	if !isValidSymbol(s0) {
		return nil
	}

	nums := 1
	pos := make([]int, 0, len(ld)*2)

	pos = append(pos, 0, ld[0])

	if isWild(s0) {
		wilds := 0
		ws := -1
		wnums := 1
		wpos := make([]int, 0, len(ld)*2)

		wpos = append(wpos, 0, ld[0])

		for x := 1; x < len(ld); x++ {
			cs := scene.Arr[x][ld[x]]

			if !isValidSymbol(cs) {
				break
			}

			if ws == -1 {
				if isWild(cs) {
					wilds++

					wnums++
					nums++

					pos = append(pos, x, ld[x])
					wpos = append(wpos, x, ld[x])
				} else {
					ws = cs

					nums++
					pos = append(pos, x, ld[x])
				}
			} else {
				if isWild(cs) {
					wilds++
				}

				if isSameSymbol(cs, ws) {
					nums++

					pos = append(pos, x, ld[x])
				} else {
					break
				}
			}
		}

		if ws == -1 {
			if wnums > 0 && pt.MapPay[s0][wnums-1] > 0 {
				r := &Result{
					Symbol:  s0,
					Type:    RTLine,
					Mul:     pt.MapPay[s0][wnums-1],
					CoinWin: pt.MapPay[s0][wnums-1],
					CashWin: pt.MapPay[s0][wnums-1] * bet,
					Pos:     wpos,
					Wilds:   wilds,
				}

				return r
			}

			return nil
		}

		wmul := 0
		mul := 0

		if wnums > 0 {
			wmul = pt.MapPay[s0][wnums-1]
		}

		if nums > 0 {
			mul = pt.MapPay[ws][nums-1]
		}

		if wmul == 0 && mul == 0 {
			return nil
		}

		if wmul >= mul {
			r := &Result{
				Symbol:  s0,
				Type:    RTLine,
				Mul:     pt.MapPay[s0][wnums-1],
				CoinWin: pt.MapPay[s0][wnums-1],
				CashWin: pt.MapPay[s0][wnums-1] * bet,
				Pos:     wpos,
				Wilds:   wilds,
			}

			return r
		}

		r := &Result{
			Symbol:  ws,
			Type:    RTLine,
			Mul:     pt.MapPay[ws][nums-1],
			CoinWin: pt.MapPay[ws][nums-1],
			CashWin: pt.MapPay[ws][nums-1] * bet,
			Pos:     pos,
			Wilds:   wilds,
		}

		return r
	}

	wilds := 0
	for x := 1; x < len(ld); x++ {
		cs := scene.Arr[x][ld[x]]

		if !isValidSymbol(cs) {
			break
		}

		if isSameSymbol(cs, s0) {
			if isWild(cs) {
				wilds++
			}

			nums++

			pos = append(pos, x, ld[x])
		} else {
			break
		}
	}

	if nums > 0 && pt.MapPay[s0][nums-1] > 0 {
		r := &Result{
			Symbol:  s0,
			Type:    RTLine,
			Mul:     pt.MapPay[s0][nums-1],
			CoinWin: pt.MapPay[s0][nums-1],
			CashWin: pt.MapPay[s0][nums-1] * bet,
			Pos:     pos,
			Wilds:   wilds,
		}

		return r
	}

	return nil
}

// CalcFullLineEx - calc fullline & no wild in reel0
func CalcFullLineEx(scene *GameScene, pt *PayTables, bet int, cursymbol int,
	countSymbolInReel FuncCountSymbolInReel) *Result {
	return nil
}

func buildFullLineResult(scene *GameScene, pt *PayTables, bet int, s0 int, arry []int) *Result {
	nums := len(arry)

	if nums > 0 && pt.MapPay[s0][nums-1] > 0 {
		r := &Result{
			Symbol:  s0,
			Type:    RTFullLine,
			Mul:     pt.MapPay[s0][nums-1],
			CoinWin: pt.MapPay[s0][nums-1],
			CashWin: pt.MapPay[s0][nums-1] * bet,
		}

		for x, y := range arry {
			r.Pos = append(r.Pos, x, y)
		}

		return r
	}

	return nil
}

// calcDeepFullLine - calc deep fullline
func calcDeepFullLine(scene *GameScene, pt *PayTables, bet int, s0 int, arry []int, results []*Result,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol) ([]*Result, bool) {

	iswin := false
	cx := len(arry)
	for y, cs := range scene.Arr[cx] {
		if isValidSymbolEx(cs, scene, cx, y) && isSameSymbol(cs, s0) {
			arry = append(arry, y)

			if cx < scene.Width-1 {
				curiswin := false
				results, iswin = calcDeepFullLine(scene, pt, bet, s0, arry, results, isValidSymbolEx, isWild, isSameSymbol)
				if curiswin {
					iswin = true
				}
			} else {
				r := buildFullLineResult(scene, pt, bet, s0, arry)

				if r != nil {
					results = append(results, r)
				}

				iswin = true
			}

			arry = arry[0:cx]
		}
	}

	if !iswin {
		r := buildFullLineResult(scene, pt, bet, s0, arry)

		if r != nil {
			results = append(results, r)
		}

		iswin = true
	}

	return results, iswin
}

// CalcFullLine - calc fullline
func CalcFullLine(scene *GameScene, pt *PayTables, bet int,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol) []*Result {

	results := []*Result{}

	yarr := make([]int, 0, scene.Width)

	for y := 0; y < scene.Height; y++ {
		s0 := scene.Arr[0][y]

		if isValidSymbolEx(s0, scene, 0, y) {
			yarr = append(yarr, y)

			results, _ = calcDeepFullLine(scene, pt, bet, s0, yarr, results, isValidSymbolEx, isWild, isSameSymbol)

			yarr = yarr[0:0]
		}
	}

	return results
}

// // CountSymbols - count symbol number
// func CountSymbols(scene *GameScene, isSymbol FuncIsSymbol) int {
// 	nums := 0

// 	for x := 0; x < len(scene.Arr); x++ {
// 		for y := 0; y < len(scene.Arr[x]); y++ {
// 			if isSymbol(scene.Arr[x][y]) {
// 				nums++
// 			}
// 		}
// 	}

// 	return nums
// }
