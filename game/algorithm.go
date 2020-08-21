package sgc7game

import "encoding/json"

// FuncIsScatter - cursymbol == scatter
type FuncIsScatter func(scatter int, cursymbol int) bool

// CalcScatter - calc scatter
func CalcScatter(scene *GameScene, pt PayTables, scatter int, isScatter FuncIsScatter) (*Result, error) {
	nums := 0
	pos := []int{}
	for x := 0; x < len(scene.Arr); x++ {
		hass := false

		for y := 0; y < len(scene.Arr[x]); y++ {
			if isScatter(scatter, scene.Arr[x][y]) {
				hass = true

				pos = append(pos, x)
				pos = append(pos, y)
			}
		}

		if hass {
			nums++
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
		}

		str, err := json.Marshal(pos)
		if err != nil {
			return nil, err
		}

		r.Data = string(str)

		return r, nil
	}

	return nil, nil
}
