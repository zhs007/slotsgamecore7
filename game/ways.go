package sgc7game

import (
	goutils "github.com/zhs007/goutils"
)

// CheckWays3
func CheckWays3(scene *GameScene, pt *PayTables, bet int,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getMulti FuncGetMulti) []*Result {

	results := []*Result{}

	arrSymbol := make([]int, 0, scene.Height)

	for y0 := 0; y0 < len(scene.Arr[0]); y0++ {
		cs := scene.Arr[0][y0]
		if !isValidSymbolEx(cs, scene, 0, y0) {
			continue
		}

		if goutils.IndexOfIntSlice(arrSymbol, cs, 0) >= 0 {
			continue
		}

		arrSymbol = append(arrSymbol, cs)

		arrpos := make([]int, 0, scene.Height*scene.Width*2)
		symbolnums := 0
		wildnums := 0
		mul := 1

		for x := 0; x < scene.Width; x++ {
			curnums := 0
			curmul := 0
			for y := 0; y < len(scene.Arr[x]); y++ {
				if !isValidSymbolEx(cs, scene, x, y) {
					continue
				}

				if isSameSymbol(scene.Arr[x][y], cs) {

					arrpos = append(arrpos, x, y)

					if isWild(scene.Arr[x][y]) {
						wildnums++
					}

					if curnums == 0 {
						symbolnums++
					}

					curnums++
					curmul += getMulti(x, y)
				}
			}

			if curnums == 0 {
				break
			}

			mul *= curmul
		}

		if symbolnums > 0 && pt.MapPay[cs][symbolnums-1] > 0 {
			r := &Result{
				Symbol:     cs,
				Type:       RTFullLineEx,
				Mul:        pt.MapPay[cs][symbolnums-1],
				CoinWin:    pt.MapPay[cs][symbolnums-1] * mul,
				CashWin:    pt.MapPay[cs][symbolnums-1] * bet * mul,
				Pos:        arrpos,
				Wilds:      wildnums,
				SymbolNums: symbolnums,
			}

			results = append(results, r)
		}
	}

	return results
}

// checkWays5Wild
func checkWays5Wild(scene *GameScene, pt *PayTables, bet int,
	cursc int,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	getMulti FuncGetMulti) *Result {

	arrpos := make([]int, 0, scene.Height*scene.Width*2)

	symbolnums := 0
	wildnums := 0
	mul := 1

	for x := 0; x < scene.Width; x++ {
		curnums := 0
		curmul := 0
		for y := 0; y < len(scene.Arr[x]); y++ {
			if !isValidSymbolEx(cursc, scene, x, y) {
				continue
			}

			if isWild(scene.Arr[x][y]) {
				arrpos = append(arrpos, x, y)

				if isWild(scene.Arr[x][y]) {
					wildnums++
				}

				if curnums == 0 {
					symbolnums++
				}

				curnums++
				curmul += getMulti(x, y)
			}
		}

		if curnums == 0 {
			break
		}

		mul *= curmul
	}

	if symbolnums > 0 && pt.MapPay[cursc][symbolnums-1] > 0 {
		r := &Result{
			Symbol:     cursc,
			Type:       RTFullLineEx,
			Mul:        pt.MapPay[cursc][symbolnums-1],
			CoinWin:    pt.MapPay[cursc][symbolnums-1] * mul,
			CashWin:    pt.MapPay[cursc][symbolnums-1] * bet * mul,
			Pos:        arrpos,
			Wilds:      wildnums,
			SymbolNums: symbolnums,
		}

		return r
	}

	return nil
}

// checkWays5
func checkWays5(scene *GameScene, pt *PayTables, bet int,
	cursc int,
	isValidSymbolEx FuncIsValidSymbolEx,
	getSymbolXY FuncGetSymbolXY,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getMulti FuncGetMulti) *Result {

	if isWild(cursc) {
		return checkWays5Wild(scene, pt, bet, cursc, isValidSymbolEx, isWild, getMulti)
	}

	arrpos := make([]int, 0, scene.Height*scene.Width*2)

	symbolnums := 0
	wildnums := 0
	mul := 1

	for x := 0; x < scene.Width; x++ {
		curnums := 0
		curmul := 0
		for y := 0; y < len(scene.Arr[x]); y++ {
			if !isValidSymbolEx(cursc, scene, x, y) {
				continue
			}

			if isSameSymbol(getSymbolXY(scene.Arr[x][y], x, y), cursc) {

				arrpos = append(arrpos, x, y)

				if isWild(scene.Arr[x][y]) {
					wildnums++
				}

				if curnums == 0 {
					symbolnums++
				}

				curnums++
				curmul += getMulti(x, y)
			}
		}

		if curnums == 0 {
			break
		}

		mul *= curmul
	}

	if symbolnums > 0 && pt.MapPay[cursc][symbolnums-1] > 0 {
		r := &Result{
			Symbol:     cursc,
			Type:       RTFullLineEx,
			Mul:        pt.MapPay[cursc][symbolnums-1],
			CoinWin:    pt.MapPay[cursc][symbolnums-1] * mul,
			CashWin:    pt.MapPay[cursc][symbolnums-1] * bet * mul,
			Pos:        arrpos,
			Wilds:      wildnums,
			SymbolNums: symbolnums,
		}

		return r
	}

	return nil
}

// CheckWays5
func CheckWays5(scene *GameScene, pt *PayTables, bet int,
	isValidSymbol FuncIsValidSymbol,
	isValidSymbolEx FuncIsValidSymbolEx,
	getSymbolXY FuncGetSymbolXY,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getMulti FuncGetMulti) []*Result {

	results := []*Result{}

	for sc := range pt.MapPay {
		if isValidSymbol(sc) {
			r := checkWays5(scene, pt, bet, sc, isValidSymbolEx, getSymbolXY, isWild, isSameSymbol, getMulti)
			if r != nil {
				results = append(results, r)
			}
		}
	}

	return results
}
