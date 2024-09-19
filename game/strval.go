package sgc7game

import (
	"log/slog"

	goutils "github.com/zhs007/goutils"
)

const StrValType string = "strval"

func NewStrVal() IVal {
	return &StrVal{}
}

func NewStrValEx(val string) IVal {
	return &StrVal{
		Val: val,
	}
}

// StrVal
type StrVal struct {
	Val string `json:"val"`
}

func (val *StrVal) Type() string {
	return StrValType
}

func (val *StrVal) IsSame(right IVal) bool {
	if right.Type() == StrValType {
		return val.Val == right.String()
	}

	return false
}

func (val *StrVal) ParseString(str string) error {
	val.Val = str

	return nil
}

func (val *StrVal) Int32() int32 {
	v, err := goutils.String2Int64(val.Val)
	if err != nil {
		goutils.Error("StrVal.Int32:String2Int64",
			slog.String("val", val.Val),
			goutils.Err(err))

		return 0
	}

	return int32(v)
}

func (val *StrVal) Int64() int64 {
	v, err := goutils.String2Int64(val.Val)
	if err != nil {
		goutils.Error("StrVal.Int64:String2Int64",
			slog.String("val", val.Val),
			goutils.Err(err))

		return 0
	}

	return v
}

func (val *StrVal) Int() int {
	v, err := goutils.String2Int64(val.Val)
	if err != nil {
		goutils.Error("StrVal.Int:String2Int64",
			slog.String("val", val.Val),
			goutils.Err(err))

		return 0
	}

	return int(v)
}

func (val *StrVal) Float32() float32 {
	v, err := goutils.String2Float64(val.Val)
	if err != nil {
		goutils.Error("StrVal.Float32:String2Float64",
			slog.String("val", val.Val),
			goutils.Err(err))

		return 0
	}

	return float32(v)
}

func (val *StrVal) Float64() float64 {
	v, err := goutils.String2Float64(val.Val)
	if err != nil {
		goutils.Error("StrVal.Float64:String2Float64",
			slog.String("val", val.Val),
			goutils.Err(err))

		return 0
	}

	return v
}

func (val *StrVal) String() string {
	return val.Val
}

// Int32Arr - return a []int32
func (val *StrVal) Int32Arr() []int32 {
	return []int32{val.Int32()}
}

// Int64Arr - return a []int64
func (val *StrVal) Int64Arr() []int64 {
	return []int64{val.Int64()}
}

// IntArr - return a []int
func (val *StrVal) IntArr() []int {
	return []int{val.Int()}
}

// Float32Arr - return a []float32
func (val *StrVal) Float32Arr() []float32 {
	return []float32{val.Float32()}
}

// Float64Arr - return a []float64
func (val *StrVal) Float64Arr() []float64 {
	return []float64{val.Float64()}
}

// StringArr - return a []string
func (val *StrVal) StringArr() []string {
	return []string{val.Val}
}

// GetInt - return val[index]
func (val *StrVal) GetInt(index int) int {
	if index == 0 {
		return val.Int()
	}

	return 0
}
