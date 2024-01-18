package stats2

type Option int

const (
	OptWins       Option = 1
	OptSymbolWins Option = 2
)

type Options []Option

func (opts Options) Has(opt Option) bool {
	for _, v := range opts {
		if v == opt {
			return true
		}
	}

	return false
}
