package sgc7game

// ResultType - result type
type ResultType int

const (
	// RTScatter - scatter
	RTScatter = 1
	// RTLine - line
	RTLine = 2
)

// Result - result for slots game
type Result struct {
	Type      ResultType
	LineIndex int
	Symbol    int
	Mul       int
	Win       int
	RealWin   int
	Pos       []int
}
