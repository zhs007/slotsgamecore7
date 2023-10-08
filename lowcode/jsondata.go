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

type basicReelsData struct {
	ReelSet       string `json:"reelSet"`
	ReelSetWeight string `json:"reelSetWeight"`
	IsExpandReel  string `json:"isExpandReel"`
}

func (basicReels *basicReelsData) build() *BasicReelsConfig {
	return &BasicReelsConfig{
		ReelSet:        basicReels.ReelSet,
		ReelSetsWeight: basicReels.ReelSetWeight,
		IsExpandReel:   basicReels.IsExpandReel == "true",
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
	AfterMain      []string `json:"afterMain"`
	BeforMain      []string `json:"beforMain"`
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
		BeforMainTriggerName: basicWins.BeforMain,
		AfterMainTriggerName: basicWins.AfterMain,
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
	RespinNumWeight               string   `json:"respinNumWeight"`
	RespinNumWithScatterNum       int      `json:"respinNumWithScatterNum"`
	RespinNumWeightWithScatterNum int      `json:"respinNumWeightWithScatterNum"`
}

func (triggerFeature *triggerFeatureData) build() *TriggerFeatureConfig {
	return &TriggerFeatureConfig{
		Symbol:          triggerFeature.Symbol[0],
		Type:            triggerFeature.Type,
		MinNum:          triggerFeature.MinNum,
		WildSymbols:     triggerFeature.WildSymbols,
		SIWMSymbols:     triggerFeature.SIWMSymbols,
		SIWMMul:         triggerFeature.SIWMMul,
		RespinNum:       triggerFeature.RespinNum,
		BetType:         triggerFeature.BetType,
		RespinNumWeight: triggerFeature.RespinNumWeight,
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

type symbolMultiData struct {
	Symbols     []string `json:"symbols"`
	StaticMulti int      `json:"staticMulti"`
}

func (symbolMulti *symbolMultiData) build() *SymbolMultiConfig {
	return &SymbolMultiConfig{
		Symbols:     symbolMulti.Symbols,
		StaticMulti: symbolMulti.StaticMulti,
		// WeightMulti          string                   `yaml:"weightMulti"`    // 倍数权重
		// OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature"`
	}
}

type symbolValData struct {
	Symbols    []string `json:"symbols"`
	DefaultVal int      `json:"defaultVal"`
	WeightVal  string   `json:"weightVal"`
}

func (symbolVal *symbolValData) build() *SymbolValConfig {
	return &SymbolValConfig{
		Symbol:     symbolVal.Symbols[0],
		WeightVal:  symbolVal.WeightVal,
		DefaultVal: symbolVal.DefaultVal,
		// OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature"`
	}
}

type symbolVal2Data struct {
	Symbols    []string `json:"symbols"`
	DefaultVal int      `json:"defaultVal"`
	WeightVal  string   `json:"weightVal"`
	WeightSet  string   `json:"weightSet"`
	WeightsVal string   `json:"weightsVal"`
}

func (symbolVal2 *symbolVal2Data) build() *SymbolVal2Config {
	return &SymbolVal2Config{
		Symbol:     symbolVal2.Symbols[0],
		WeightSet:  symbolVal2.WeightSet,
		DefaultVal: symbolVal2.DefaultVal,
		// WeightsVal           []string                 `yaml:"weightsVal"`
		// RNGSet               string                   `yaml:"RNGSet"`
		// OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature"`
	}
}

type bookOfData struct {
	BetType        string   `json:"betType"`
	WildSymbols    []string `json:"wildSymbols"`
	WeightTrigger  string   `json:"weightTrigger"`
	ForceSymbolNum int      `json:"forceSymbolNum"`
	WeightSymbol   string   `json:"weightSymbol"`
}

func (bookOfData *bookOfData) build() *BookOfConfig {
	return &BookOfConfig{
		BetType:       bookOfData.BetType,
		WeightTrigger: bookOfData.WeightTrigger,
		WeightSymbol:  bookOfData.WeightSymbol,
		WildSymbols:   bookOfData.WildSymbols,
		// ForceTrigger     :bookOfData.ForceTrigger,
		// WeightSymbolNum      string   `yaml:"weightSymbolNum" json:"weightSymbolNum"`
		ForceSymbolNum: bookOfData.ForceSymbolNum,
		// SymbolRNG            string   `yaml:"symbolRNG" json:"symbolRNG"`               // 只在ForceSymbolNum为1时有效
		// SymbolCollection     string   `yaml:"symbolCollection" json:"symbolCollection"` // 图标从一个SymbolCollection里获取
	}
}
