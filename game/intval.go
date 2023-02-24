package sgc7game

import (
	"fmt"

	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

const IntValType string = "intval"

func NewIntVal[T int | int32 | int64]() IVal {
	return &IntVal[T]{}
}

func NewIntValEx[T int | int32 | int64](v T) IVal {
	return &IntVal[T]{
		Val: v,
	}
}

// IntVal
type IntVal[T int | int32 | int64] struct {
	Val T `json:"val"`
}

func (val *IntVal[T]) Type() string {
	return IntValType
}

func (val *IntVal[T]) IsSame(right IVal) bool {
	if right.Type() == IntValType {
		return val.Int64() == right.Int64()
	}

	return false
}

func (val *IntVal[T]) ParseString(str string) error {
	v, err := goutils.String2Int64(str)
	if err != nil {
		goutils.Error("IntVal[T].ParseString:String2Int64",
			zap.String("str", str),
			zap.Error(err))

		return err
	}

	val.Val = T(v)

	return nil
}

func (val *IntVal[T]) Int32() int32 {
	return int32(val.Val)
}

func (val *IntVal[T]) Int64() int64 {
	return int64(val.Val)
}

func (val *IntVal[T]) Int() int {
	return int(val.Val)
}

func (val *IntVal[T]) Float32() float32 {
	return float32(val.Val)
}

func (val *IntVal[T]) Float64() float64 {
	return float64(val.Val)
}

func (val *IntVal[T]) String() string {
	return fmt.Sprintf("%v", val.Val)
}

// Int32Arr - return a []int32
func (val *IntVal[T]) Int32Arr() []int32 {
	return []int32{int32(val.Val)}
}

// Int64Arr - return a []int64
func (val *IntVal[T]) Int64Arr() []int64 {
	return []int64{int64(val.Val)}
}

// IntArr - return a []int
func (val *IntVal[T]) IntArr() []int {
	return []int{int(val.Val)}
}

// Float32Arr - return a []float32
func (val *IntVal[T]) Float32Arr() []float32 {
	return []float32{float32(val.Val)}
}

// Float64Arr - return a []float64
func (val *IntVal[T]) Float64Arr() []float64 {
	return []float64{float64(val.Val)}
}

// StringArr - return a []string
func (val *IntVal[T]) StringArr() []string {
	return []string{val.String()}
}
