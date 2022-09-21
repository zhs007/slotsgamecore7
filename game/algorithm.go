package sgc7game

import (
	goutils "github.com/zhs007/goutils"
)

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

// FuncCalcOtherMul - calc other multi
type FuncCalcOtherMul func(scene *GameScene, result *Result) int

// FuncCalcOtherMulEx - calc other multi
type FuncCalcOtherMulEx func(scene *GameScene, symbol int, pos []int) int

// FuncGetSymbol - get symbol
type FuncGetSymbol func(cursymbol int) int

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
			Symbol:     scatter,
			Type:       RTScatter,
			LineIndex:  -1,
			Mul:        pt.MapPay[scatter][nums-1],
			CoinWin:    pt.MapPay[scatter][nums-1] * coins,
			CashWin:    pt.MapPay[scatter][nums-1] * coins * bet,
			Pos:        pos,
			SymbolNums: nums,
		}

		return r
	}

	return nil
}

// CalcScatter2 - calc scatter
func CalcScatter2(scene *GameScene, pt *PayTables, scatter int, bet int, coins int,
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

	if nums > len(pt.MapPay[scatter]) {
		nums = len(pt.MapPay[scatter])
	}

	if nums > 0 && pt.MapPay[scatter][nums-1] > 0 {
		r := &Result{
			Symbol:     scatter,
			Type:       RTScatter,
			LineIndex:  -1,
			Mul:        pt.MapPay[scatter][nums-1],
			CoinWin:    pt.MapPay[scatter][nums-1] * coins,
			CashWin:    pt.MapPay[scatter][nums-1] * coins * bet,
			Pos:        pos,
			SymbolNums: nums,
		}

		return r
	}

	return nil
}

// CalcScatterEx - calc scatter
func CalcScatterEx(scene *GameScene, scatter int, nums int, isScatter FuncIsScatter) *Result {
	curnums := 0
	pos := make([]int, 0, len(scene.Arr)*len(scene.Arr[0])*2)
	for x := 0; x < len(scene.Arr); x++ {
		for y := 0; y < len(scene.Arr[x]); y++ {
			if isScatter(scatter, scene.Arr[x][y]) {
				curnums++

				pos = append(pos, x, y)
			}
		}
	}

	if curnums >= nums {
		r := &Result{
			Symbol:     scatter,
			Type:       RTScatterEx,
			LineIndex:  -1,
			Pos:        pos,
			SymbolNums: curnums,
		}

		return r
	}

	return nil
}

// CalcScatterOnReels - calc scatter
func CalcScatterOnReels(scene *GameScene, scatter int, nums int, isScatter FuncIsScatter) *Result {
	curnums := 0
	reelnums := 0
	pos := make([]int, 0, len(scene.Arr)*len(scene.Arr[0])*2)
	for x := 0; x < len(scene.Arr); x++ {
		hasscatter := false
		for y := 0; y < len(scene.Arr[x]); y++ {
			if isScatter(scatter, scene.Arr[x][y]) {
				curnums++

				if !hasscatter {
					reelnums++

					hasscatter = true
				}

				pos = append(pos, x, y)
			}
		}
	}

	if reelnums >= nums {
		r := &Result{
			Symbol:     scatter,
			Type:       RTScatterOnReels,
			LineIndex:  -1,
			Pos:        pos,
			SymbolNums: reelnums,
		}

		return r
	}

	return nil
}

// CalcLine - calc line
func CalcLine(scene *GameScene, pt *PayTables, ld []int, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol) *Result {

	s0 := getSymbol(scene.Arr[0][ld[0]])
	if !isValidSymbol(s0) {
		return nil
	}

	nums := 1
	pos := make([]int, 0, len(ld)*2)

	pos = append(pos, 0, ld[0])

	if isWild(s0) {
		wilds := 1
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
					Symbol:     s0,
					Type:       RTLine,
					Mul:        pt.MapPay[s0][wnums-1],
					CoinWin:    pt.MapPay[s0][wnums-1],
					CashWin:    pt.MapPay[s0][wnums-1] * bet,
					Pos:        wpos,
					Wilds:      wilds,
					SymbolNums: wnums,
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
				Symbol:     s0,
				Type:       RTLine,
				Mul:        pt.MapPay[s0][wnums-1],
				CoinWin:    pt.MapPay[s0][wnums-1],
				CashWin:    pt.MapPay[s0][wnums-1] * bet,
				Pos:        wpos,
				Wilds:      wilds,
				SymbolNums: wnums,
			}

			return r
		}

		r := &Result{
			Symbol:     ws,
			Type:       RTLine,
			Mul:        pt.MapPay[ws][nums-1],
			CoinWin:    pt.MapPay[ws][nums-1],
			CashWin:    pt.MapPay[ws][nums-1] * bet,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
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
			Symbol:     s0,
			Type:       RTLine,
			Mul:        pt.MapPay[s0][nums-1],
			CoinWin:    pt.MapPay[s0][nums-1],
			CashWin:    pt.MapPay[s0][nums-1] * bet,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
		}

		return r
	}

	return nil
}

// CalcLineEx - calc line
func CalcLineEx(scene *GameScene, pt *PayTables, ld []int, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	calcOtherMul FuncCalcOtherMul,
	getSymbol FuncGetSymbol) *Result {
	r := CalcLine(scene, pt, ld, bet, isValidSymbol, isWild, isSameSymbol, getSymbol)
	if r != nil {
		r.OtherMul = calcOtherMul(scene, r)

		if r.OtherMul > 1 {
			r.CoinWin = r.CoinWin * r.OtherMul
			r.CashWin = r.CashWin * r.OtherMul
		}
	}

	return r
}

// CalcLineOtherMul - calc line with otherMul
func CalcLineOtherMul(scene *GameScene, pt *PayTables, ld []int, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	calcOtherMul FuncCalcOtherMulEx,
	getSymbol FuncGetSymbol) *Result {
	s0 := getSymbol(scene.Arr[0][ld[0]])
	if !isValidSymbol(s0) {
		return nil
	}

	nums := 1
	pos := make([]int, 0, len(ld)*2)

	pos = append(pos, 0, ld[0])

	if isWild(s0) {
		wilds := 1
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
				wothermul := calcOtherMul(scene, s0, wpos)

				r := &Result{
					Symbol:     s0,
					Type:       RTLine,
					Mul:        pt.MapPay[s0][wnums-1],
					CoinWin:    pt.MapPay[s0][wnums-1] * wothermul,
					CashWin:    pt.MapPay[s0][wnums-1] * wothermul * bet,
					Pos:        wpos,
					Wilds:      wilds,
					SymbolNums: wnums,
					OtherMul:   wothermul,
				}

				return r
			}

			return nil
		}

		wmul := 0
		mul := 0
		wothermul := 1
		othermul := 1

		if wnums > 0 {
			wothermul = calcOtherMul(scene, s0, wpos)
			wmul = pt.MapPay[s0][wnums-1]
		}

		if nums > 0 {
			othermul = calcOtherMul(scene, ws, pos)
			mul = pt.MapPay[ws][nums-1]
		}

		if wmul == 0 && mul == 0 {
			return nil
		}

		if wmul*wothermul >= mul*othermul {
			r := &Result{
				Symbol:     s0,
				Type:       RTLine,
				Mul:        pt.MapPay[s0][wnums-1],
				CoinWin:    pt.MapPay[s0][wnums-1] * wothermul,
				CashWin:    pt.MapPay[s0][wnums-1] * wothermul * bet,
				Pos:        wpos,
				Wilds:      wilds,
				SymbolNums: wnums,
				OtherMul:   wothermul,
			}

			return r
		}

		r := &Result{
			Symbol:     ws,
			Type:       RTLine,
			Mul:        pt.MapPay[ws][nums-1],
			CoinWin:    pt.MapPay[ws][nums-1] * othermul,
			CashWin:    pt.MapPay[ws][nums-1] * othermul * bet,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   othermul,
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
		othermul := calcOtherMul(scene, s0, pos)

		r := &Result{
			Symbol:     s0,
			Type:       RTLine,
			Mul:        pt.MapPay[s0][nums-1],
			CoinWin:    pt.MapPay[s0][nums-1] * othermul,
			CashWin:    pt.MapPay[s0][nums-1] * othermul * bet,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   othermul,
		}

		return r
	}

	return nil
}

// CalcFullLineEx - calc fullline & no wild in reel0
//		用数个数的方式来计算全线游戏，第一轴不能有wild
func CalcFullLineEx(scene *GameScene, pt *PayTables, bet int,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol) []*Result {

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
				}
			}

			if curnums == 0 {
				break
			}

			mul *= curnums
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

// calcSymbolFullLineEx2 - calc fullline
//		用数个数的方式来计算全线游戏
func calcSymbolFullLineEx2(scene *GameScene, pt *PayTables, symbol int, bet int,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol) *Result {

	arrpos := make([]int, 0, scene.Height*scene.Width*2)
	symbolnums := 0
	wildnums := 0
	mul := 1

	for x := 0; x < scene.Width; x++ {
		curnums := 0
		for y := 0; y < len(scene.Arr[x]); y++ {
			if !isValidSymbolEx(symbol, scene, x, y) {
				continue
			}

			if isSameSymbol(scene.Arr[x][y], symbol) {

				arrpos = append(arrpos, x, y)

				if isWild(scene.Arr[x][y]) {
					wildnums++
				}

				if curnums == 0 {
					symbolnums++
				}

				curnums++
			}
		}

		if curnums == 0 {
			break
		}

		mul *= curnums
	}

	if symbolnums > 0 && pt.MapPay[symbol][symbolnums-1] > 0 {
		r := &Result{
			Symbol:     symbol,
			Type:       RTFullLineEx,
			Mul:        pt.MapPay[symbol][symbolnums-1],
			CoinWin:    pt.MapPay[symbol][symbolnums-1] * mul,
			CashWin:    pt.MapPay[symbol][symbolnums-1] * bet * mul,
			Pos:        arrpos,
			Wilds:      wildnums,
			SymbolNums: symbolnums,
		}

		return r
	}

	return nil
}

// CalcFullLineEx2 - calc fullline
//		用数个数的方式来计算全线游戏
func CalcFullLineEx2(scene *GameScene, pt *PayTables, bet int,
	isValidSymbolEx FuncIsValidSymbolEx,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol) []*Result {

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

		// is wild
		if isWild(cs) {
			arrpos := make([]int, 0, scene.Height*scene.Width*2)
			symbolnums := 0
			wildnums := 0
			mul := 1

			for x := 0; x < scene.Width; x++ {
				curnums := 0

				for y := 0; y < len(scene.Arr[x]); y++ {
					if isWild(scene.Arr[x][y]) {
						arrpos = append(arrpos, x, y)

						if isWild(scene.Arr[x][y]) {
							wildnums++
						}

						if curnums == 0 {
							symbolnums++
						}

						curnums++
					}
				}

				if curnums == 0 {
					break
				}

				mul *= curnums
			}

			wx := symbolnums

			if wx == scene.Width {
				if pt.MapPay[cs][wx-1] > 0 {
					r := &Result{
						Symbol:     cs,
						Type:       RTFullLineEx,
						Mul:        pt.MapPay[cs][wx-1],
						CoinWin:    pt.MapPay[cs][wx-1] * mul,
						CashWin:    pt.MapPay[cs][wx-1] * bet * mul,
						Pos:        arrpos,
						Wilds:      wildnums,
						SymbolNums: wx,
					}

					results = append(results, r)
				}

				continue
			}

			for ty := 0; ty < len(scene.Arr[wx]); ty++ {
				ws := scene.Arr[wx][ty]
				if !isValidSymbolEx(ws, scene, wx, ty) {
					continue
				}

				if goutils.IndexOfIntSlice(arrSymbol, ws, 0) >= 0 {
					continue
				}

				arrSymbol = append(arrSymbol, ws)

				r := calcSymbolFullLineEx2(scene, pt, ws, bet, isValidSymbolEx, isWild, isSameSymbol)
				if r != nil {
					results = append(results, r)
				}
			}

			continue
		}

		r := calcSymbolFullLineEx2(scene, pt, cs, bet, isValidSymbolEx, isWild, isSameSymbol)
		if r != nil {
			results = append(results, r)
		}
	}

	return results
}

func buildFullLineResult(scene *GameScene, pt *PayTables, bet int, s0 int, arry []int, wildNums int) *Result {
	nums := len(arry)

	if nums > 0 && pt.MapPay[s0][nums-1] > 0 {
		r := &Result{
			Symbol:     s0,
			Type:       RTFullLine,
			Mul:        pt.MapPay[s0][nums-1],
			CoinWin:    pt.MapPay[s0][nums-1],
			CashWin:    pt.MapPay[s0][nums-1] * bet,
			SymbolNums: nums,
			Wilds:      wildNums,
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
	isSameSymbol FuncIsSameSymbol, wildNums int) ([]*Result, bool) {

	iswin := false
	cx := len(arry)
	// wildnums := 0

	for y, cs := range scene.Arr[cx] {
		if isValidSymbolEx(cs, scene, cx, y) && isSameSymbol(cs, s0) {
			if isWild(cs) {
				wildNums++
			}

			arry = append(arry, y)

			if cx < scene.Width-1 {
				curiswin := false
				results, iswin = calcDeepFullLine(scene, pt, bet, s0, arry, results, isValidSymbolEx, isWild, isSameSymbol, wildNums)
				if curiswin {
					iswin = true
				}
			} else {
				r := buildFullLineResult(scene, pt, bet, s0, arry, wildNums)

				if r != nil {
					results = append(results, r)
				}

				iswin = true
			}

			arry = arry[0:cx]
		}
	}

	if !iswin {
		r := buildFullLineResult(scene, pt, bet, s0, arry, wildNums)

		if r != nil {
			results = append(results, r)
		}

		iswin = true
	}

	return results, iswin
}

// CalcFullLine - calc fullline & no wild in reel0
// 		还是用算线的方式来计算全线游戏，效率会低一些，数据量也会大一些，但某些特殊类型的游戏只能这样计算
//		也没有考虑第一轴有wild的情况（后续可调整算法）
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

			results, _ = calcDeepFullLine(scene, pt, bet, s0, yarr, results, isValidSymbolEx, isWild, isSameSymbol, 0)

			yarr = yarr[0:0]
		}
	}

	return results
}
