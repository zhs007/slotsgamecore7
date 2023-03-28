package lowcode

type CollectorData struct {
	Val          int // 当前总值, Current total value
	NewCollector int // 这一个step收集到的, The values collected in this step
}

// OnNewGame -
func (cd *CollectorData) onNewGame() {
	cd.Val = 0
}

// OnNewStep -
func (cd *CollectorData) onNewStep() {
	cd.NewCollector = 0
}

func NewCollectorData() *CollectorData {
	return &CollectorData{}
}
