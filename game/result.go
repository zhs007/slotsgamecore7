package sgc7game

// ResultType - result type
type ResultType int

const (
	// RTScatter - scatter
	RTScatter = 1
	// RTLine - line
	RTLine = 2
	// RTFullLine - full line
	// 243线游戏，最多会出现243条记录，某些游戏必须要这种方式
	RTFullLine = 3
	// RTFullLineEx - full line ex
	// 全线游戏的汇总记录，243线最多3条记录
	RTFullLineEx = 4
	// RTScatterEx - scatter只计数量
	RTScatterEx = 5
)

// Result - result for slots game
type Result struct {
	Type       ResultType `json:"type"`
	LineIndex  int        `json:"lineindex"`
	Symbol     int        `json:"symbol"`
	Mul        int        `json:"mul"`
	CoinWin    int        `json:"coinwin"`
	CashWin    int        `json:"cashwin"`
	Pos        []int      `json:"pos"`
	OtherMul   int        `json:"othermul"`
	Wilds      int        `json:"wilds"`
	SymbolNums int        `json:"symbolnums"`
}
