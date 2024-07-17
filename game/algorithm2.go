package sgc7game

// CalcLine2 - calc line
func CalcLine2(scene *GameScene, pt *PayTables, ld []int, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol,
	getMulti FuncGetMulti) *Result {

	sx := 0

	s0 := getSymbol(scene.Arr[sx][ld[sx]])
	if !isValidSymbol(s0) {
		return nil
	}

	nums := 1
	pos := make([]int, 0, len(ld)*2)

	pos = append(pos, 0, ld[sx])
	otherMul := getMulti(sx, ld[sx])

	if isWild(s0) {
		wilds := 1
		ws := -1
		wnums := 1
		wpos := make([]int, 0, len(ld)*2)

		wpos = append(wpos, sx, ld[sx])
		wotherMul := getMulti(sx, ld[sx])

		for x := 1; x < len(ld); x++ {
			cs := scene.Arr[sx+x][ld[sx+x]]

			if !isValidSymbol(cs) && !isWild(cs) {
				break
			}

			if ws == -1 {
				if isWild(cs) {
					wilds++

					wnums++
					nums++

					pos = append(pos, sx+x, ld[sx+x])
					wpos = append(wpos, sx+x, ld[sx+x])
					wotherMul *= getMulti(sx+x, ld[sx+x])
					otherMul *= getMulti(sx+x, ld[sx+x])
				} else {
					ws = cs

					nums++
					pos = append(pos, sx+x, ld[sx+x])
					otherMul *= getMulti(sx+x, ld[sx+x])
				}
			} else {
				if isWild(cs) {
					wilds++
				}

				if isSameSymbol(cs, ws) {
					nums++

					pos = append(pos, sx+x, ld[sx+x])
					otherMul *= getMulti(sx+x, ld[sx+x])
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
					CoinWin:    pt.MapPay[s0][wnums-1] * wotherMul,
					CashWin:    pt.MapPay[s0][wnums-1] * bet * wotherMul,
					Pos:        wpos,
					Wilds:      wilds,
					SymbolNums: wnums,
					OtherMul:   wotherMul,
				}

				return r
			}

			return nil
		}

		wmul := 0
		mul := 0

		if wnums > 0 {
			wmul = pt.MapPay[s0][wnums-1] * wotherMul
		}

		if nums > 0 {
			mul = pt.MapPay[ws][nums-1] * otherMul
		}

		if wmul == 0 && mul == 0 {
			return nil
		}

		if wmul >= mul {
			r := &Result{
				Symbol:     s0,
				Type:       RTLine,
				Mul:        pt.MapPay[s0][wnums-1],
				CoinWin:    pt.MapPay[s0][wnums-1] * wotherMul,
				CashWin:    pt.MapPay[s0][wnums-1] * bet * wotherMul,
				Pos:        wpos,
				Wilds:      wilds,
				SymbolNums: wnums,
				OtherMul:   wotherMul,
			}

			return r
		}

		r := &Result{
			Symbol:     ws,
			Type:       RTLine,
			Mul:        pt.MapPay[ws][nums-1],
			CoinWin:    pt.MapPay[ws][nums-1] * otherMul,
			CashWin:    pt.MapPay[ws][nums-1] * bet * otherMul,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   otherMul,
		}

		return r
	}

	wilds := 0
	for x := 1; x < len(ld); x++ {
		cs := scene.Arr[sx+x][ld[sx+x]]

		if !isValidSymbol(cs) && !isWild(cs) {
			break
		}

		if isSameSymbol(cs, s0) {
			if isWild(cs) {
				wilds++
			}

			nums++

			pos = append(pos, sx+x, ld[sx+x])
			otherMul *= getMulti(sx+x, ld[sx+x])
		} else {
			break
		}
	}

	if nums > 0 && pt.MapPay[s0][nums-1] > 0 {
		r := &Result{
			Symbol:     s0,
			Type:       RTLine,
			Mul:        pt.MapPay[s0][nums-1],
			CoinWin:    pt.MapPay[s0][nums-1] * otherMul,
			CashWin:    pt.MapPay[s0][nums-1] * bet * otherMul,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   otherMul,
		}

		return r
	}

	return nil
}

// CalcLineRL2 - calc line with right->left
func CalcLineRL2(scene *GameScene, pt *PayTables, ld []int, bet int,
	isValidSymbol FuncIsValidSymbol,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol,
	getMulti FuncGetMulti) *Result {

	sx := len(scene.Arr) - 1

	s0 := getSymbol(scene.Arr[sx][ld[sx]])
	if !isValidSymbol(s0) {
		return nil
	}

	nums := 1
	pos := make([]int, 0, len(ld)*2)

	pos = append(pos, sx, ld[sx])
	otherMul := getMulti(sx, ld[sx])

	if isWild(s0) {
		wilds := 1
		ws := -1
		wnums := 1
		wpos := make([]int, 0, len(ld)*2)

		wpos = append(wpos, sx, ld[sx])
		wotherMul := getMulti(sx, ld[sx])

		for x := 1; x < len(ld); x++ {
			cs := scene.Arr[sx-x][ld[sx-x]]

			if !isValidSymbol(cs) && !isWild(cs) {
				break
			}

			if ws == -1 {
				if isWild(cs) {
					wilds++

					wnums++
					nums++

					pos = append(pos, sx-x, ld[sx-x])
					wpos = append(wpos, sx-x, ld[sx-x])
					otherMul *= getMulti(sx-x, ld[sx-x])
					wotherMul *= getMulti(sx-x, ld[sx-x])
				} else {
					ws = cs

					nums++
					pos = append(pos, sx-x, ld[sx-x])
					otherMul *= getMulti(sx-x, ld[sx-x])
				}
			} else {
				if isWild(cs) {
					wilds++
				}

				if isSameSymbol(cs, ws) {
					nums++

					pos = append(pos, sx-x, ld[sx-x])
					otherMul *= getMulti(sx-x, ld[sx-x])
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
					CoinWin:    pt.MapPay[s0][wnums-1] * wotherMul,
					CashWin:    pt.MapPay[s0][wnums-1] * bet * wotherMul,
					Pos:        wpos,
					Wilds:      wilds,
					SymbolNums: wnums,
					OtherMul:   wotherMul,
				}

				return r
			}

			return nil
		}

		wmul := 0
		mul := 0

		if wnums > 0 {
			wmul = pt.MapPay[s0][wnums-1] * wotherMul
		}

		if nums > 0 {
			mul = pt.MapPay[ws][nums-1] * otherMul
		}

		if wmul == 0 && mul == 0 {
			return nil
		}

		if wmul >= mul {
			r := &Result{
				Symbol:     s0,
				Type:       RTLine,
				Mul:        pt.MapPay[s0][wnums-1],
				CoinWin:    pt.MapPay[s0][wnums-1] * wotherMul,
				CashWin:    pt.MapPay[s0][wnums-1] * bet * wotherMul,
				Pos:        wpos,
				Wilds:      wilds,
				SymbolNums: wnums,
				OtherMul:   wotherMul,
			}

			return r
		}

		r := &Result{
			Symbol:     ws,
			Type:       RTLine,
			Mul:        pt.MapPay[ws][nums-1],
			CoinWin:    pt.MapPay[ws][nums-1] * otherMul,
			CashWin:    pt.MapPay[ws][nums-1] * bet * otherMul,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   otherMul,
		}

		return r
	}

	wilds := 0
	for x := 1; x < len(ld); x++ {
		cs := scene.Arr[sx-x][ld[sx-x]]

		if !isValidSymbol(cs) && !isWild(cs) {
			break
		}

		if isSameSymbol(cs, s0) {
			if isWild(cs) {
				wilds++
			}

			nums++

			pos = append(pos, sx-x, ld[sx-x])
			otherMul *= getMulti(sx-x, ld[sx-x])
		} else {
			break
		}
	}

	if nums > 0 && pt.MapPay[s0][nums-1] > 0 {
		r := &Result{
			Symbol:     s0,
			Type:       RTLine,
			Mul:        pt.MapPay[s0][nums-1],
			CoinWin:    pt.MapPay[s0][nums-1] * otherMul,
			CashWin:    pt.MapPay[s0][nums-1] * bet * otherMul,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   otherMul,
		}

		return r
	}

	return nil
}

// CountSymbolOnLine - count on line
func CountSymbolOnLine(scene *GameScene, pt *PayTables, ld []int, bet int, symbol int,
	isWild FuncIsWild,
	isSameSymbol FuncIsSameSymbol,
	getSymbol FuncGetSymbol,
	getMulti FuncGetMulti, calcMulti FuncCalcMulti) *Result {

	sx := 0

	s0 := getSymbol(scene.Arr[sx][ld[sx]])

	nums := 0
	pos := make([]int, 0, len(ld)*2)

	otherMul := 1

	if isSameSymbol(s0, symbol) {
		nums++

		pos = append(pos, 0, ld[sx])
		otherMul = getMulti(sx, ld[sx])

		if isWild(s0) {
			wilds := 1
			ws := -1
			wnums := 1
			wpos := make([]int, 0, len(ld)*2)

			wpos = append(wpos, sx, ld[sx])
			wotherMul := getMulti(sx, ld[sx])

			for x := 1; x < len(ld); x++ {
				cs := scene.Arr[sx+x][ld[sx+x]]

				if isSameSymbol(cs, symbol) {
					if ws == -1 {
						if isWild(cs) {
							wilds++

							wnums++
							nums++

							pos = append(pos, sx+x, ld[sx+x])
							wpos = append(wpos, sx+x, ld[sx+x])
							wotherMul = calcMulti(wotherMul, getMulti(sx+x, ld[sx+x]))
							otherMul = calcMulti(otherMul, getMulti(sx+x, ld[sx+x]))
						} else {
							ws = symbol

							nums++
							pos = append(pos, sx+x, ld[sx+x])
							otherMul = calcMulti(otherMul, getMulti(sx+x, ld[sx+x]))
						}
					} else {
						if isWild(cs) {
							wilds++
						}

						nums++

						pos = append(pos, sx+x, ld[sx+x])
						otherMul = calcMulti(otherMul, getMulti(sx+x, ld[sx+x]))
					}
				}
			}

			if ws == -1 {
				if wnums > 0 && pt.MapPay[s0][wnums-1] > 0 {
					r := &Result{
						Symbol:     s0,
						Type:       RTLine,
						Mul:        pt.MapPay[s0][wnums-1],
						CoinWin:    pt.MapPay[s0][wnums-1] * wotherMul,
						CashWin:    pt.MapPay[s0][wnums-1] * bet * wotherMul,
						Pos:        wpos,
						Wilds:      wilds,
						SymbolNums: wnums,
						OtherMul:   wotherMul,
					}

					return r
				}

				return nil
			}

			wmul := 0
			mul := 0

			if wnums > 0 {
				wmul = pt.MapPay[s0][wnums-1] * wotherMul
			}

			if nums > 0 {
				mul = pt.MapPay[ws][nums-1] * otherMul
			}

			if wmul == 0 && mul == 0 {
				return nil
			}

			if wmul >= mul {
				r := &Result{
					Symbol:     s0,
					Type:       RTLine,
					Mul:        pt.MapPay[s0][wnums-1],
					CoinWin:    pt.MapPay[s0][wnums-1] * wotherMul,
					CashWin:    pt.MapPay[s0][wnums-1] * bet * wotherMul,
					Pos:        wpos,
					Wilds:      wilds,
					SymbolNums: wnums,
					OtherMul:   wotherMul,
				}

				return r
			}

			r := &Result{
				Symbol:     ws,
				Type:       RTLine,
				Mul:        pt.MapPay[ws][nums-1],
				CoinWin:    pt.MapPay[ws][nums-1] * otherMul,
				CashWin:    pt.MapPay[ws][nums-1] * bet * otherMul,
				Pos:        pos,
				Wilds:      wilds,
				SymbolNums: nums,
				OtherMul:   otherMul,
			}

			return r
		}
	}

	wilds := 0
	for x := 1; x < len(ld); x++ {
		cs := scene.Arr[sx+x][ld[sx+x]]

		if isSameSymbol(cs, symbol) {
			if isWild(cs) {
				wilds++
			}

			nums++

			pos = append(pos, sx+x, ld[sx+x])
			otherMul = calcMulti(otherMul, getMulti(sx+x, ld[sx+x]))
		}
	}

	if nums > 0 && pt.MapPay[symbol][nums-1] > 0 {
		r := &Result{
			Symbol:     symbol,
			Type:       RTLine,
			Mul:        pt.MapPay[symbol][nums-1],
			CoinWin:    pt.MapPay[symbol][nums-1] * otherMul,
			CashWin:    pt.MapPay[symbol][nums-1] * bet * otherMul,
			Pos:        pos,
			Wilds:      wilds,
			SymbolNums: nums,
			OtherMul:   otherMul,
		}

		return r
	}

	return nil
}
