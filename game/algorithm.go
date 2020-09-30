package sgc7game

// FuncIsSymbol - is symbol
type FuncIsSymbol func(cursymbol int) bool

// FuncIsScatter - cursymbol == scatter
type FuncIsScatter func(scatter int, cursymbol int) bool

// FuncIsWild - cursymbol == wild
type FuncIsWild func(cursymbol int) bool

// FuncIsSameSymbol - cursymbol == startsymbol
type FuncIsSameSymbol func(cursymbol int, startsymbol int) bool

// FuncIsValidSymbol - is it a valid symbol?
type FuncIsValidSymbol func(cursymbol int) bool

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
					HasW:    true,
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
				HasW:    true,
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
			HasW:    true,
		}

		return r
	}

	hasw := false
	for x := 1; x < len(ld); x++ {
		cs := scene.Arr[x][ld[x]]

		if !isValidSymbol(cs) {
			break
		}

		if isSameSymbol(cs, s0) {
			if isWild(cs) {
				hasw = true
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
			HasW:    hasw,
		}

		return r
	}

	return nil
}

// CountSymbols - count symbol number
func CountSymbols(scene *GameScene, isSymbol FuncIsSymbol) int {
	nums := 0

	for x := 0; x < len(scene.Arr); x++ {
		for y := 0; y < len(scene.Arr[x]); y++ {
			if isSymbol(scene.Arr[x][y]) {
				nums++
			}
		}
	}

	return nums
}
