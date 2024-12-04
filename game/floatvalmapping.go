package sgc7game

import (
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
)

// FloatValMapping - a mapping for float values
type FloatValMapping[T int, V float32 | float64] struct {
	MapVals map[T]V `json:"mapVals"`
}

func (vm *FloatValMapping[T, V]) Keys() []T {
	lst := []T{}

	for k := range vm.MapVals {
		lst = append(lst, k)
	}

	return lst
}

func (vm *FloatValMapping[T, V]) Clone() *FloatValMapping[T, V] {
	nvm := &FloatValMapping[T, V]{
		MapVals: make(map[T]V),
	}

	for k, v := range vm.MapVals {
		nvm.MapVals[k] = v
	}

	return nvm
}

func NewFloatValMapping[T int, V float32 | float64](typevals []T, vals []V) (*FloatValMapping[T, V], error) {
	if len(typevals) != len(vals) {
		goutils.Error("NewFloatValMapping",
			slog.Int("typevals", len(typevals)),
			slog.Int("vals", len(vals)),
			goutils.Err(ErrInvalidValMapping))

		return nil, ErrInvalidValMapping
	}

	vm := &FloatValMapping[T, V]{
		MapVals: make(map[T]V),
	}

	for i, k := range typevals {
		vm.MapVals[k] = vals[i]
	}

	return vm, nil
}

func NewFloatValMappingEx[T int, V float32 | float64]() *FloatValMapping[T, V] {
	return &FloatValMapping[T, V]{
		MapVals: make(map[T]V),
	}
}

// LoadFloatValMappingFromExcel - load xlsx file
func LoadFloatValMappingFromExcel[T int, V float32 | float64](fn string, headerType string, headerVal string) (*FloatValMapping[T, V], error) {
	typevals := []T{}
	vals := []V{}

	err := LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		if header == headerType {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadFloatValMappingFromExcel:String2Int64",
					slog.String("header", header),
					slog.String("data", data),
					goutils.Err(err))

				return err
			}

			typevals = append(typevals, T(v))
		} else if header == headerVal {
			v, err := goutils.String2Float64(data)
			if err != nil {
				goutils.Error("LoadFloatValMappingFromExcel:String2Float64",
					slog.String("header", header),
					slog.String("data", data),
					goutils.Err(err))

				return err
			}

			vals = append(vals, V(v))
		}
		return nil
	})
	if err != nil {
		goutils.Error("LoadFloatValMappingFromExcel:LoadExcel",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	return NewFloatValMapping(typevals, vals)
}
