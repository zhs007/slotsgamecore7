package sgc7game

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

type reelsInfo5 struct {
	R1   int `json:"R1"`
	R2   int `json:"R2"`
	R3   int `json:"R3"`
	R4   int `json:"R4"`
	R5   int `json:"R5"`
	Line int `json:"line"`
}

// ReelsData - reels data
type ReelsData struct {
	Reels [][]int `json:"reels"`
}

// isValidRI5 - is it valid reelsInfo5
func isValidRI5(ri5s []reelsInfo5) bool {
	if len(ri5s) <= 0 {
		return false
	}

	// alllinezero := true
	for _, v := range ri5s {
		if v.Line > 0 {
			// alllinezero = false

			return true
		}
	}

	return false
}

// LoadReels5JSON - load json file
func LoadReels5JSON(fn string) (*ReelsData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	w := 5

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var ri []reelsInfo5
	err = json.Unmarshal(data, &ri)
	if err != nil {
		return nil, err
	}

	if !isValidRI5(ri) {
		return nil, nil
	}

	p := &ReelsData{
		Reels: [][]int{},
	}

	for i := 0; i < w; i++ {
		p.Reels = append(p.Reels, []int{})
	}

	for _, v := range ri {
		if v.R1 >= 0 {
			p.Reels[0] = append(p.Reels[0], v.R1)
		}

		if v.R2 >= 0 {
			p.Reels[1] = append(p.Reels[1], v.R2)
		}

		if v.R3 >= 0 {
			p.Reels[2] = append(p.Reels[2], v.R3)
		}

		if v.R4 >= 0 {
			p.Reels[3] = append(p.Reels[3], v.R4)
		}

		if v.R5 >= 0 {
			p.Reels[4] = append(p.Reels[4], v.R5)
		}
	}

	return p, nil
}

// LoadReels3JSON - load json file
func LoadReels3JSON(fn string) (*ReelsData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	w := 3

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var ri []reelsInfo5
	err = json.Unmarshal(data, &ri)
	if err != nil {
		return nil, err
	}

	if !isValidRI5(ri) {
		return nil, nil
	}

	p := &ReelsData{
		Reels: [][]int{},
	}

	for i := 0; i < w; i++ {
		p.Reels = append(p.Reels, []int{})
	}

	for _, v := range ri {
		if v.R1 >= 0 {
			p.Reels[0] = append(p.Reels[0], v.R1)
		}

		if v.R2 >= 0 {
			p.Reels[1] = append(p.Reels[1], v.R2)
		}

		if v.R3 >= 0 {
			p.Reels[2] = append(p.Reels[2], v.R3)
		}
	}

	return p, nil
}

// DropDownIntoGameScene - 用轮子当前位置处理下落
//		注意：
//			1. 这个接口需要特别注意，传入indexes是上一次用过的，所以实际用应该-1
//			2. 这个接口按道理只会对index做减法操作，所以不会考虑向下越界问题，只处理向上的越界
func (rd *ReelsData) DropDownIntoGameScene(scene *GameScene, indexes []int) ([]int, error) {
	narr := []int{}
	for x, arr := range scene.Arr {

		ci := indexes[x]

		for y, v := range arr {
			if v == -1 {
				ci--
				if ci < 0 {
					ci += len(rd.Reels[x])
				}

				scene.Arr[x][y] = rd.Reels[x][ci]
			}
		}

		narr = append(narr, ci)
	}

	return narr, nil
}
