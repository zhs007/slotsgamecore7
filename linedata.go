package slotsgamecore7

type lineInfo struct {
	R1   int `json:"R1"`
	R2   int `json:"R2"`
	R3   int `json:"R3"`
	R4   int `json:"R4"`
	R5   int `json:"R5"`
	Line int `json:"line"`
}

// LineData - line data
type LineData struct {
	Lines [][]int
}

// LoadLineJSON - load json file
func LoadLineJSON(fn string) (*LineData, error) {
	return nil, nil
}
