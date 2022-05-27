package sgc7game

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type reelsInfo5 struct {
	R1   int `json:"R1"`
	R2   int `json:"R2"`
	R3   int `json:"R3"`
	R4   int `json:"R4"`
	R5   int `json:"R5"`
	Line int `json:"line"`
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

// ReelsData - reels data
type ReelsData struct {
	Reels [][]int `json:"reels"`
}

// 主要用于BuildReelsPosData接口
type FuncReelsDataPos func(rd *ReelsData, x, y int) bool

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

// BuildReelsPosData - 构建轮子坐标数据，一般用于后续的转轮算法，主要起到随机优化效率用
func (rd *ReelsData) BuildReelsPosData(onpos FuncReelsDataPos) (*ReelsPosData, error) {
	rpd := NewReelsPosData(rd)

	for x, arr := range rd.Reels {
		for y := range arr {
			if onpos(rd, x, y) {
				rpd.AddPos(x, y)
			}
		}
	}

	return rpd, nil
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

// LoadReelsFromExcel - load xlsx file
func LoadReelsFromExcel(fn string) (*ReelsData, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadReelsFromExcel:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadReelsFromExcel:GetSheetList",
			goutils.JSON("SheetList", lstname),
			zap.String("fn", fn),
			zap.Error(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadReelsFromExcel:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	p := &ReelsData{
		Reels: [][]int{},
	}

	// x -> ri
	mapri := make(map[int]int)
	maxri := 0
	isend := []bool{}

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				if colCell[0] == 'r' || colCell[0] == 'R' {
					iv, err := goutils.String2Int64(colCell[1:])
					if err != nil {
						goutils.Error("LoadReelsFromExcel:String2Int64",
							zap.String("fn", fn),
							zap.String("header", colCell),
							zap.Error(err))

						return nil, err
					}

					if iv <= 0 {
						goutils.Error("LoadReelsFromExcel",
							zap.String("info", "check iv"),
							zap.String("fn", fn),
							zap.String("header", colCell),
							zap.Error(ErrInvalidReelsExcelFile))

						return nil, ErrInvalidReelsExcelFile
					}

					mapri[x] = int(iv) - 1
					if int(iv) > maxri {
						maxri = int(iv)
					}
				}
			}

			if maxri != len(mapri) {
				goutils.Error("LoadReelsFromExcel",
					zap.String("info", "check len"),
					zap.String("fn", fn),
					zap.Int("maxri", maxri),
					goutils.JSON("mapri", mapri),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			if maxri <= 0 {
				goutils.Error("LoadReelsFromExcel",
					zap.String("info", "check empty"),
					zap.String("fn", fn),
					zap.Int("maxri", maxri),
					goutils.JSON("mapri", mapri),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			for i := 0; i < maxri; i++ {
				p.Reels = append(p.Reels, []int{})
				isend = append(isend, false)
			}
		} else {
			for x, colCell := range row {
				ri, isok := mapri[x]
				if isok {
					v, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadReelsFromExcel:String2Int64",
							zap.String("val", colCell),
							zap.Error(err))

						return nil, err
					}

					if v < 0 {
						isend[ri] = true
					} else if isend[ri] {
						goutils.Error("LoadReelsFromExcel",
							zap.String("info", "check already finished."),
							zap.String("val", colCell),
							zap.Int("y", y),
							zap.Int("x", x),
							zap.Error(err))

						return nil, err
					} else {
						p.Reels[ri] = append(p.Reels[ri], int(v))
					}
				}
			}
		}
	}

	return p, nil
}
