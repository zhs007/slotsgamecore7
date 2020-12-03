package sgc7game

// RemoveSymbolWithResult - remove symbol with win result
func RemoveSymbolWithResult(scene *GameScene, result *PlayResult) error {
	for _, v := range result.Results {
		for i := 0; i < len(v.Pos)/2; i++ {
			scene.Arr[v.Pos[i*2]][v.Pos[i*2+1]] = -1
		}
	}

	return nil
}

// DropDownSymbols - drop down symbols
func DropDownSymbols(scene *GameScene) error {
	for _, arr := range scene.Arr {
		for y := len(arr) - 1; y >= 0; {
			if arr[y] == -1 {
				hass := false
				for y1 := y - 1; y1 >= 0; y1-- {
					if arr[y1] != -1 {
						arr[y] = arr[y1]
						arr[y1] = -1

						hass = true
						y = y1
						break
					}
				}

				if !hass {
					break
				}
			} else {
				y--
			}
		}
	}

	return nil
}
