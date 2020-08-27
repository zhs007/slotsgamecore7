package sgc7game

// FuncIsScatter - cursymbol == scatter
type FuncIsScatter func(scatter int, cursymbol int) bool

// FuncIsWild - cursymbol == wild
type FuncIsWild func(cursymbol int) bool

// FuncIsSameSymbol - cursymbol == symbol
type FuncIsSameSymbol func(symbol int, cursymbol int) bool

// FuncIsValidSymbol - is it a valid symbol?
type FuncIsValidSymbol func(cursymbol int) bool

// CalcScatter - calc scatter
func CalcScatter(scene *GameScene, pt *PayTables, scatter int, bet int,
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

	if pt.MapPay[scatter][nums-1] > 0 {
		r := &Result{
			Symbol:    scatter,
			Type:      RTScatter,
			LineIndex: -1,
			Mul:       pt.MapPay[scatter][nums-1],
			Win:       pt.MapPay[scatter][nums-1] * bet,
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
			if pt.MapPay[s0][wnums-1] > 0 {
				r := &Result{
					Symbol: s0,
					Type:   RTLine,
					Mul:    pt.MapPay[s0][wnums-1],
					Win:    pt.MapPay[s0][wnums-1] * bet,
					Pos:    wpos,
				}

				return r
			}

			return nil
		}

		wmul := pt.MapPay[s0][wnums-1]
		mul := pt.MapPay[ws][nums-1]

		if wmul == 0 && mul == 0 {
			return nil
		}

		if wmul >= mul {
			r := &Result{
				Symbol: s0,
				Type:   RTLine,
				Mul:    pt.MapPay[s0][wnums-1],
				Win:    pt.MapPay[s0][wnums-1] * bet,
				Pos:    wpos,
			}

			return r
		}

		r := &Result{
			Symbol: ws,
			Type:   RTLine,
			Mul:    pt.MapPay[ws][nums-1],
			Win:    pt.MapPay[ws][nums-1] * bet,
			Pos:    pos,
		}

		return r
	}

	for x := 1; x < len(ld); x++ {
		cs := scene.Arr[x][ld[x]]

		if !isValidSymbol(cs) {
			break
		}

		if isSameSymbol(cs, s0) {
			nums++

			pos = append(pos, x, ld[x])
		} else {
			break
		}

	}

	if pt.MapPay[s0][nums-1] > 0 {
		r := &Result{
			Symbol: s0,
			Type:   RTLine,
			Mul:    pt.MapPay[s0][nums-1],
			Win:    pt.MapPay[s0][nums-1] * bet,
			Pos:    pos,
		}

		return r
	}

	return nil
}
