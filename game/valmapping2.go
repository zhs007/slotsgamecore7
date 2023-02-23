package sgc7game

import (
	"strings"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// ValMapping2
type ValMapping2 struct {
	MapVals map[int]IVal `json:"mapVals"`
}

func (vm *ValMapping2) Keys() []int {
	lst := []int{}

	for k := range vm.MapVals {
		lst = append(lst, k)
	}

	return lst
}

func (vm *ValMapping2) Clone() *ValMapping2 {
	nvm := &ValMapping2{
		MapVals: make(map[int]IVal),
	}

	for k, v := range vm.MapVals {
		nvm.MapVals[k] = v
	}

	return nvm
}

func NewValMapping2(typevals []int, vals []IVal) (*ValMapping2, error) {
	if len(typevals) != len(vals) {
		goutils.Error("NewValMapping",
			zap.Int("typevals", len(typevals)),
			zap.Int("vals", len(vals)),
			zap.Error(ErrInvalidValMapping))

		return nil, ErrInvalidValMapping
	}

	vm := &ValMapping2{
		MapVals: make(map[int]IVal),
	}

	for i, k := range typevals {
		vm.MapVals[k] = vals[i]
	}

	return vm, nil
}

func NewValMappingEx2() *ValMapping2 {
	return &ValMapping2{
		MapVals: make(map[int]IVal),
	}
}

// LoadValMapping2FromExcel - load xlsx file
func LoadValMapping2FromExcel(fn string, headerType string, headerVal string, funcNew FuncNewIVal) (*ValMapping2, error) {
	typevals := []int{}
	vals := []IVal{}

	err := LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		if header == headerType {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadValMapping2FromExcel:String2Int64",
					zap.String("header", header),
					zap.String("data", data),
					zap.Error(err))

				return err
			}

			typevals = append(typevals, int(v))
		} else if header == headerVal {
			cv := funcNew()
			err := cv.ParseString(data)
			if err != nil {
				goutils.Error("LoadValMapping2FromExcel:ParseString",
					zap.String("header", header),
					zap.String("data", data),
					zap.Error(err))

				return err
			}

			vals = append(vals, cv)
		}
		return nil
	})
	if err != nil {
		goutils.Error("LoadValMapping2FromExcel:LoadExcel",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return NewValMapping2(typevals, vals)
}
