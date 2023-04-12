package sgc7game

import (
	"fmt"
	"strings"

	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

const IntArrValType string = "intarrval"

func NewIntArrVal[T int | int32 | int64]() IVal {
	return &IntArrVal[T]{}
}

// IntArrVal
type IntArrVal[T int | int32 | int64] struct {
	Vals []T `json:"vals"`
}

func (val *IntArrVal[T]) Type() string {
	return IntArrValType
}

func (val *IntArrVal[T]) IsSame(right IVal) bool {
	if right.Type() == IntArrValType {
		i64arr := right.Int64Arr()
		if len(val.Vals) == len(i64arr) {
			for i, v := range i64arr {
				if val.Vals[i] != T(v) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func (val *IntArrVal[T]) ParseString(str string) error {
	varr := []T{}

	arr := strings.Split(str, ",")
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v != "" {
			iv, err := goutils.String2Int64(v)
			if err != nil {
				goutils.Error("IntArrVal[T].ParseString:String2Int64",
					zap.String("str", str),
					zap.Error(err))

				return err
			}

			varr = append(varr, T(iv))
		}
	}

	val.Vals = varr

	return nil
}

func (val *IntArrVal[T]) Int32() int32 {
	if len(val.Vals) > 0 {
		return int32(val.Vals[0])
	}

	return 0
}

func (val *IntArrVal[T]) Int64() int64 {
	if len(val.Vals) > 0 {
		return int64(val.Vals[0])
	}

	return 0
}

func (val *IntArrVal[T]) Int() int {
	if len(val.Vals) > 0 {
		return int(val.Vals[0])
	}

	return 0
}

func (val *IntArrVal[T]) Float32() float32 {
	if len(val.Vals) > 0 {
		return float32(val.Vals[0])
	}

	return 0
}

func (val *IntArrVal[T]) Float64() float64 {
	if len(val.Vals) > 0 {
		return float64(val.Vals[0])
	}

	return 0
}

func (val *IntArrVal[T]) String() string {
	str := ""
	for i, v := range val.Vals {
		if i == 0 {
			str = fmt.Sprintf("%v", v)
		} else {
			str += fmt.Sprintf(",%v", v)
		}
	}

	return str
}

// Int32Arr - return a []int32
func (val *IntArrVal[T]) Int32Arr() []int32 {
	arr := make([]int32, len(val.Vals))

	for i, v := range val.Vals {
		arr[i] = int32(v)
	}

	return arr
}

// Int64Arr - return a []int64
func (val *IntArrVal[T]) Int64Arr() []int64 {
	arr := make([]int64, len(val.Vals))

	for i, v := range val.Vals {
		arr[i] = int64(v)
	}

	return arr
}

// IntArr - return a []int
func (val *IntArrVal[T]) IntArr() []int {
	arr := make([]int, len(val.Vals))

	for i, v := range val.Vals {
		arr[i] = int(v)
	}

	return arr
}

// Float32Arr - return a []float32
func (val *IntArrVal[T]) Float32Arr() []float32 {
	arr := make([]float32, len(val.Vals))

	for i, v := range val.Vals {
		arr[i] = float32(v)
	}

	return arr
}

// Float64Arr - return a []float64
func (val *IntArrVal[T]) Float64Arr() []float64 {
	arr := make([]float64, len(val.Vals))

	for i, v := range val.Vals {
		arr[i] = float64(v)
	}

	return arr
}

// StringArr - return a []string
func (val *IntArrVal[T]) StringArr() []string {
	arr := make([]string, len(val.Vals))

	for i, v := range val.Vals {
		arr[i] = fmt.Sprintf("%v", v)
	}

	return arr
}

// GetInt - return val[index]
func (val *IntArrVal[T]) GetInt(index int) int {
	return int(val.Vals[index])
}
