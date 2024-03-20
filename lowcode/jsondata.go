package lowcode

type paytableData struct {
	Code   int    `json:"Code"`
	Symbol string `json:"Symbol"`
	Data   []int  `json:"data"`
}

type weightData struct {
	Val    string `json:"val"`
	Weight int    `json:"weight"`
}

type mappingData struct {
	In  string `json:"in"`
	Out string `json:"out"`
}
