package sgc7game

import "unicode"

func IsRngString(str string) bool {
	for _, r := range str {
		if !(unicode.IsDigit(r) || r == ',' || unicode.IsSpace(r)) {
			return false
		}
	}

	return true
}
