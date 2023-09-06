package sgc7game

import (
	"fmt"
	"os"
	"strings"

	"github.com/bytedance/sonic"
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

func (rd *ReelsData) SetReel(ri int, reel []int) {
	rd.Reels[ri] = reel
}

// DropDownIntoGameScene - 用轮子当前位置处理下落
//
//	注意：
//		1. 这个接口需要特别注意，传入indexes是上一次用过的，所以实际用应该-1
//		2. 这个接口按道理只会对index做减法操作，所以不会考虑向下越界问题，只处理向上的越界
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

// DropDownIntoGameScene2 - 用轮子当前位置处理下落
//
//	注意：
//		1. 这个接口需要特别注意，传入indexes是上一次用过的，所以实际用应该-1
//		2. 这个接口按道理只会对index做减法操作，所以不会考虑向下越界问题，只处理向上的越界
func (rd *ReelsData) DropDownIntoGameScene2(scene *GameScene, indexes []int) ([]int, error) {
	narr := []int{}
	for x, arr := range scene.Arr {

		ci := indexes[x]

		for y := len(arr) - 1; y >= 0; y-- {
			v := arr[y]
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

func (rd *ReelsData) SaveExcel(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "line")
	for i := range rd.Reels {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("R%v", i+1))
	}

	maxj := 0

	for i, reel := range rd.Reels {
		if maxj < len(reel) {
			maxj = len(reel)
		}

		for j, v := range reel {
			f.SetCellInt(sheet, goutils.Pos2Cell(i+1, j+1), v)
		}
	}

	for i := 0; i < maxj; i++ {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, i+1), i)
	}

	return f.SaveAs(fn)
}

func (rd *ReelsData) SaveExcelEx(fn string, paytables *PayTables) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "line")
	for i := range rd.Reels {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("R%v", i+1))
	}

	maxj := 0

	for i, reel := range rd.Reels {
		if maxj < len(reel) {
			maxj = len(reel)
		}

		for j, v := range reel {
			f.SetCellStr(sheet, goutils.Pos2Cell(i+1, j+1), paytables.GetStringFromInt(v))
		}
	}

	for i := 0; i < maxj; i++ {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, i+1), i)
	}

	return f.SaveAs(fn)
}

// LoadReels5JSON - load json file
func LoadReels5JSON(fn string) (*ReelsData, error) {
	w := 5

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var ri []reelsInfo5
	err = sonic.Unmarshal(data, &ri)
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
	w := 3

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var ri []reelsInfo5
	err = sonic.Unmarshal(data, &ri)
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
	defer f.Close()

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
					colCell = strings.TrimSpace(colCell)
					if len(colCell) > 0 {
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
					} else {
						isend[ri] = true
					}
				}
			}
		}
	}

	return p, nil
}

// LoadReelsFromExcel2 - load xlsx file
func LoadReelsFromExcel2(fn string, paytables *PayTables) (*ReelsData, error) {
	p := &ReelsData{
		Reels: [][]int{},
	}

	// x -> ri
	mapri := make(map[int]int)
	maxri := 0
	isend := []bool{}
	isfirst := true

	err := LoadExcel(fn, "", func(x int, str string) string {
		header := strings.ToLower(strings.TrimSpace(str))
		if header[0] == 'r' {
			iv, err := goutils.String2Int64(header[1:])
			if err != nil {
				goutils.Error("LoadReelsFromExcel2:LoadExcel:String2Int64",
					zap.String("fn", fn),
					zap.String("header", header),
					zap.Error(err))

				return ""
			}

			if iv <= 0 {
				goutils.Error("LoadReelsFromExcel2:LoadExcel",
					zap.String("info", "check iv"),
					zap.String("fn", fn),
					zap.String("header", header),
					zap.Error(ErrInvalidReelsExcelFile))

				return ""
			}

			mapri[x] = int(iv) - 1
			if int(iv) > maxri {
				maxri = int(iv)
			}
		}

		return header
	}, func(x int, y int, header string, data string) error {
		if isfirst {
			isfirst = false

			if maxri != len(mapri) {
				goutils.Error("LoadReelsFromExcel2",
					zap.String("info", "check len"),
					zap.String("fn", fn),
					zap.Int("maxri", maxri),
					goutils.JSON("mapri", mapri),
					zap.Error(ErrInvalidReelsExcelFile))

				return ErrInvalidReelsExcelFile
			}

			if maxri <= 0 {
				goutils.Error("LoadReelsFromExcel2",
					zap.String("info", "check empty"),
					zap.String("fn", fn),
					zap.Int("maxri", maxri),
					goutils.JSON("mapri", mapri),
					zap.Error(ErrInvalidReelsExcelFile))

				return ErrInvalidReelsExcelFile
			}

			for i := 0; i < maxri; i++ {
				p.Reels = append(p.Reels, []int{})
				isend = append(isend, false)
			}
		}

		ri, isok := mapri[x]
		if isok {
			data = strings.TrimSpace(data)
			if len(data) > 0 {
				s, isok := paytables.MapSymbols[data]
				if isok {
					if isend[ri] {
						goutils.Error("LoadReelsFromExcel2",
							zap.String("info", "check already finished."),
							zap.String("val", data),
							zap.Int("y", y),
							zap.Int("x", x),
							zap.Error(ErrInvalidReelsExcelFile))

						return ErrInvalidReelsExcelFile
					}

					p.Reels[ri] = append(p.Reels[ri], s)
				} else {
					isend[ri] = true
				}
			} else {
				isend[ri] = true
			}
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadReelsFromExcel2:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return p, nil
}

// NewReelsData - new ReelsData
func NewReelsData(num int) *ReelsData {
	rd := &ReelsData{
		Reels: make([][]int, num),
	}

	return rd
}
