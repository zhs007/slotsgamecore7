package sgc7game

type FuncNewIVal func() IVal

// IVal
type IVal interface {
	// Type - get type of IVal
	Type() string

	// ParseString - str -> IVal
	ParseString(str string) error
	// IsSame - return this == right
	IsSame(right IVal) bool

	// Int32 - return a int32
	Int32() int32
	// Int64 - return a int64
	Int64() int64
	// Int - return a int
	Int() int
	// Float32 - return a float32
	Float32() float32
	// Float64 - return a float64
	Float64() float64
	// String - return a string
	String() string

	// Int32Arr - return a []int32
	Int32Arr() []int32
	// Int64Arr - return a []int64
	Int64Arr() []int64
	// IntArr - return a []int
	IntArr() []int
	// Float32Arr - return a []float32
	Float32Arr() []float32
	// Float64Arr - return a []float64
	Float64Arr() []float64
	// StringArr - return a []string
	StringArr() []string
}
