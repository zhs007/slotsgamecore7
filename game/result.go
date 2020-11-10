package sgc7game

// ResultType - result type
type ResultType int

const (
	// RTScatter - scatter
	RTScatter = 1
	// RTLine - line
	RTLine = 2
	// RTFullLine - full line
	RTFullLine = 3
	// RTFullLineEx - full line ex
	RTFullLineEx = 4
)

// Result - result for slots game
type Result struct {
	Type      ResultType `json:"type"`
	LineIndex int        `json:"lineindex"`
	Symbol    int        `json:"symbol"`
	Mul       int        `json:"mul"`
	CoinWin   int        `json:"coinwin"`
	CashWin   int        `json:"cashwin"`
	Pos       []int      `json:"pos"`
	OtherMul  int        `json:"othermul"`
	Wilds     int        `json:"wilds"`
}
