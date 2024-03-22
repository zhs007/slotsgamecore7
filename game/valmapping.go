package sgc7game

import (
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
)

// ValMapping
type ValMapping[T int, V int | int32 | int64] struct {
	MapVals map[T]V `json:"mapVals"`
}

func (vm *ValMapping[T, V]) Keys() []T {
	lst := []T{}

	for k := range vm.MapVals {
		lst = append(lst, k)
	}

	return lst
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

func NewValMapping[T int, V int | int32 | int64](typevals []T, vals []V) (*ValMapping[T, V], error) {
	if len(typevals) != len(vals) {
		goutils.Error("NewValMapping",
			slog.Int("typevals", len(typevals)),
			slog.Int("vals", len(vals)),
			goutils.Err(ErrInvalidValMapping))

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

func NewValMappingEx[T int, V int | int32 | int64]() *ValMapping[T, V] {
	return &ValMapping[T, V]{
		MapVals: make(map[T]V),
	}
}

// LoadValMappingFromExcel - load xlsx file
func LoadValMappingFromExcel[T int, V int | int32 | int64](fn string, headerType string, headerVal string) (*ValMapping[T, V], error) {
	typevals := []T{}
	vals := []V{}

	err := LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		if header == headerType {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadValMappingFromExcel:String2Int64",
					slog.String("header", header),
					slog.String("data", data),
					goutils.Err(err))

				return err
			}

			typevals = append(typevals, T(v))
		} else if header == headerVal {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadValMappingFromExcel:String2Int64",
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
		goutils.Error("LoadValMappingFromExcel:LoadExcel",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	return NewValMapping(typevals, vals)
}
