package sgc7game

import (
	"strings"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// ValArrMapping
type ValArrMapping[T int, V int | int32 | int64] struct {
	MapVals map[T][]V `json:"mapVals"`
}

func (vm *ValArrMapping[T, V]) Keys() []T {
	lst := []T{}

	for k := range vm.MapVals {
		lst = append(lst, k)
	}

	return lst
}

func (vm *ValArrMapping[T, V]) Clone() *ValArrMapping[T, V] {
	nvm := &ValArrMapping[T, V]{
		MapVals: make(map[T][]V),
	}

	for k, v := range vm.MapVals {
		nvm.MapVals[k] = v
	}

	return nvm
}

func NewValArrMapping[T int, V int | int32 | int64](typevals []T, vals [][]V) (*ValArrMapping[T, V], error) {
	if len(typevals) != len(vals) {
		goutils.Error("NewValArrMapping",
			zap.Int("typevals", len(typevals)),
			zap.Int("vals", len(vals)),
			zap.Error(ErrInvalidValMapping))

		return nil, ErrInvalidValMapping
	}

	vm := &ValArrMapping[T, V]{
		MapVals: make(map[T][]V),
	}

	for i, k := range typevals {
		vm.MapVals[k] = make([]V, len(vals[i]))
		copy(vm.MapVals[k], vals[i])
	}

	return vm, nil
}

func NewValArrMappingEx[T int, V int | int32 | int64]() *ValArrMapping[T, V] {
	return &ValArrMapping[T, V]{
		MapVals: make(map[T][]V),
	}
}

// LoadValArrMappingFromExcel - load xlsx file
func LoadValArrMappingFromExcel[T int, V int | int32 | int64](fn string, headerType string, headerVal string) (*ValArrMapping[T, V], error) {
	typevals := []T{}
	vals := [][]V{}

	err := LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		if header == headerType {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadValArrMappingFromExcel:String2Int64",
					zap.String("header", header),
					zap.String("data", data),
					zap.Error(err))

				return err
			}

			typevals = append(typevals, T(v))
		} else if header == headerVal {
			varr := []V{}

			arr := strings.Split(data, ",")
			for _, v := range arr {
				v = strings.TrimSpace(v)
				if v != "" {
					iv, err := goutils.String2Int64(v)
					if err != nil {
						goutils.Error("LoadValArrMappingFromExcel:String2Int64",
							zap.String("val", v),
							zap.Error(err))

						return err
					}

					varr = append(varr, V(iv))
				}
			}

			vals = append(vals, varr)
		}
		return nil
	})
	if err != nil {
		goutils.Error("LoadValArrMappingFromExcel:LoadExcel",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return NewValArrMapping(typevals, vals)
}
