package sgc7utils

const RuneStr = "0"

// Rune2Int
func Rune2Int(r rune) int {
	return int(byte(r) - RuneStr[0])
}
