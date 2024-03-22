package sgc7game

import (
	"fmt"
	"log/slog"

	goutils "github.com/zhs007/goutils"
)

const FloatValType string = "floatval"

func NewFloatVal[T float32 | float64]() IVal {
	return &FloatVal[T]{}
}

func NewFloatValEx[T float32 | float64](v T) IVal {
	return &FloatVal[T]{
		Val: v,
	}
}

// FloatVal
type FloatVal[T float32 | float64] struct {
	Val T `json:"val"`
}

func (val *FloatVal[T]) Type() string {
	return FloatValType
}

func (val *FloatVal[T]) IsSame(right IVal) bool {
	if right.Type() == FloatValType {
		return val.Float64() == right.Float64()
	}

	return false
}

func (val *FloatVal[T]) ParseString(str string) error {
	v, err := goutils.String2Float64(str)
	if err != nil {
		goutils.Error("FloatVal[T].ParseString:String2Float64",
			slog.String("str", str),
			goutils.Err(err))

		return err
	}

	val.Val = T(v)

	return nil
}

func (val *FloatVal[T]) Int32() int32 {
	return int32(val.Val)
}

func (val *FloatVal[T]) Int64() int64 {
	return int64(val.Val)
}

func (val *FloatVal[T]) Int() int {
	return int(val.Val)
}

func (val *FloatVal[T]) Float32() float32 {
	return float32(val.Val)
}

func (val *FloatVal[T]) Float64() float64 {
	return float64(val.Val)
}

func (val *FloatVal[T]) String() string {
	return fmt.Sprintf("%v", val.Val)
}

// Int32Arr - return a []int32
func (val *FloatVal[T]) Int32Arr() []int32 {
	return []int32{int32(val.Val)}
}

// Int64Arr - return a []int64
func (val *FloatVal[T]) Int64Arr() []int64 {
	return []int64{int64(val.Val)}
}

// IntArr - return a []int
func (val *FloatVal[T]) IntArr() []int {
	return []int{int(val.Val)}
}

// Float32Arr - return a []float32
func (val *FloatVal[T]) Float32Arr() []float32 {
	return []float32{float32(val.Val)}
}

// Float64Arr - return a []float64
func (val *FloatVal[T]) Float64Arr() []float64 {
	return []float64{float64(val.Val)}
}

// StringArr - return a []string
func (val *FloatVal[T]) StringArr() []string {
	return []string{val.String()}
}

// GetInt - return val[index]
func (val *FloatVal[T]) GetInt(index int) int {
	if index == 0 {
		return int(val.Val)
	}

	return 0
}
