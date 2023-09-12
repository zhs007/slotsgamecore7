package lowcode

type paytableData struct {
	Code   int    `json:"Code"`
	Symbol string `json:"Symbol"`
	Data   []int  `json:"data"`
}

type basicReelsData struct {
	ReelSet      string `json:"reelSet"`
	IsExpandReel string `json:"isExpandReel"`
}

func (basicReels *basicReelsData) build() *BasicReelsConfig {
	return &BasicReelsConfig{
		ReelSet:      basicReels.ReelSet,
		IsExpandReel: basicReels.IsExpandReel == "true",
	}
}

type basicWinsData struct {
	MainType       string   `json:"mainType"`
	BetType        string   `json:"betType"`
	ExcludeSymbols []string `json:"excludeSymbols"`
	WildSymbols    []string `json:"wildSymbols"`
	CheckWinType   string   `json:"checkWinType"`
	SIWMSymbols    []string `json:"SIWMSymbols"`
	SIWMMul        int      `json:"SIWMMul"`
	AfterMain      string   `json:"afterMain"`
	BeforMain      string   `json:"beforMain"`
}

func (basicWins *basicWinsData) build() *BasicWinsConfig {
	return &BasicWinsConfig{
		MainType:             basicWins.MainType,
		BetType:              basicWins.BetType,
		StrCheckWinType:      basicWins.CheckWinType,
		SIWMSymbols:          basicWins.SIWMSymbols,
		SIWMMul:              basicWins.SIWMMul,
		ExcludeSymbols:       basicWins.ExcludeSymbols,
		WildSymbols:          basicWins.WildSymbols,
		BeforMainTriggerName: []string{basicWins.BeforMain},
		AfterMainTriggerName: []string{basicWins.AfterMain},
	}
}

type triggerFeatureData struct {
	Label                         string   `json:"label"`
	Symbol                        []string `json:"symbol"`
	Type                          string   `json:"type"`
	BetType                       string   `json:"betType"`
	WildSymbols                   []string `json:"wildSymbols"`
	MinNum                        int      `json:"minNum"`
	CheckWinType                  string   `json:"checkWinType"`
	SIWMSymbols                   []string `json:"SIWMSymbols"`
	SIWMMul                       int      `json:"SIWMMul"`
	RespinNum                     int      `json:"respinNum"`
	RespinNumWithScatterNum       int      `json:"respinNumWithScatterNum"`
	RespinNumWeightWithScatterNum int      `json:"respinNumWeightWithScatterNum"`
}

func (triggerFeature *triggerFeatureData) build() *TriggerFeatureConfig {
	return &TriggerFeatureConfig{
		Symbol:      triggerFeature.Symbol[0],
		Type:        triggerFeature.Type,
		MinNum:      triggerFeature.MinNum,
		WildSymbols: triggerFeature.WildSymbols,
		SIWMSymbols: triggerFeature.SIWMSymbols,
		SIWMMul:     triggerFeature.SIWMMul,
		RespinNum:   triggerFeature.RespinNum,
		BetType:     triggerFeature.BetType,
		// RespinNumWeight               string         `yaml:"respinNumWeight"`               // respin number weight
		// RespinNumWithScatterNum       map[int]int    `yaml:"respinNumWithScatterNum"`       // respin number with scatter number
		// RespinNumWeightWithScatterNum map[int]string `yaml:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
		// CountScatterPayAs             string         `yaml:"countScatterPayAs"`             // countscatter时，按什么符号赔付
		// RespinComponent               string         `yaml:"respinComponent"`               // like fg-spin
		// NextComponent                 string         `yaml:"nextComponent"`                 // next component
		// TagSymbolNum                  string         `yaml:"tagSymbolNum"`                  // 这里可以将symbol数量记下来，别的地方能获取到
		// Awards                        []*Award       `yaml:"awards"`                        // 新的奖励系统
		// SymbolAwardsWeights           *AwardsWeights `yaml:"symbolAwardsWeights"`           // 每个中奖符号随机一组奖励
	}
}
