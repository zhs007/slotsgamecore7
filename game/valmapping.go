package sgc7game

import (
	"strings"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// ValMapping
type ValMapping[T int, V int] struct {
	MapVals map[T]V
}

func (vm *ValMapping[T, V]) Clone() *ValMapping[T, V] {
	nvm := &ValMapping[T, V]{
		MapVals: make(map[T]V),
	}

	for k, v := range vm.MapVals {
		nvm.MapVals[k] = v
	}

	return nvm
}

func NewValMapping[T int, V int](typevals []T, vals []V) (*ValMapping[T, V], error) {
	if len(typevals) != len(vals) {
		goutils.Error("NewValMapping",
			zap.Int("typevals", len(typevals)),
			zap.Int("vals", len(vals)),
			zap.Error(ErrInvalidValMapping))

		return nil, ErrInvalidValMapping
	}

	vm := &ValMapping[T, V]{
		MapVals: make(map[T]V),
	}

	for i, k := range typevals {
		vm.MapVals[k] = vals[i]
	}

	return vm, nil
}

// LoadValMappingFromExcel - load xlsx file
func LoadValMappingFromExcel[T int, V int](fn string, headerType string, headerVal string) (*ValMapping[T, V], error) {
	typevals := []T{}
	vals := []V{}

	err := LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		if header == headerType {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadValMappingFromExcel:String2Int64",
					zap.String("header", header),
					zap.String("data", data),
					zap.Error(err))

				return err
			}

			typevals = append(typevals, T(v))
		} else if header == headerVal {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadValMappingFromExcel:String2Int64",
					zap.String("header", header),
					zap.String("data", data),
					zap.Error(err))

				return err
			}

			vals = append(vals, V(v))
		}
		return nil
	})
	if err != nil {
		goutils.Error("LoadValMappingFromExcel:LoadExcel",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return NewValMapping(typevals, vals)
}
