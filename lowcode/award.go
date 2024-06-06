package lowcode

type Award struct {
	AwardType       string   `yaml:"awardType" json:"awardType"`
	Type            int      `yaml:"-" json:"-"`
	Val             int      `yaml:"val" json:"-"`                           // 弃用，代码里已经不用了，初始化时会把数据转存到Vals里，为了兼容性保留配置
	StrParam        string   `yaml:"strParam" json:"-"`                      // 弃用，代码里已经不用了，初始化时会把数据转存到StrParams里，为了兼容性保留配置
	Vals            []int    `yaml:"vals" json:"vals"`                       // 数值参数
	StrParams       []string `yaml:"strParams" json:"strParams"`             // 字符串参数
	ComponentVals   []string `yaml:"componentVals" json:"componentVals"`     // 可以用component数值来替代常量，如果val长度为2，需要替换第二个参数，那么第一个参数应该给空字符串
	OnTriggerRespin string   `yaml:"onTriggerRespin" json:"onTriggerRespin"` // 在这个respin再次触发时才生效，这个时候会用当前respin的LastTriggerNum+CurTriggerNum作为TriggerIndex记下，当TriggerIndex==CurTriggerNum时才生效
	TriggerIndex    int      `yaml:"-" json:"-"`                             // 见上
}

func (cfg *Award) getType() int {
	if cfg.AwardType == "cash" {
		return AwardCash
	} else if cfg.AwardType == "collector" {
		return AwardCollector
	} else if cfg.AwardType == "respinTimes" {
		return AwardRespinTimes
	} else if cfg.AwardType == "gameMulti" {
		return AwardGameMulti
	} else if cfg.AwardType == "stepMulti" {
		return AwardStepMulti
	} else if cfg.AwardType == "initMask" {
		return AwardInitMask
	} else if cfg.AwardType == "triggerRespin" {
		return AwardTriggerRespin
	} else if cfg.AwardType == "noLevelUpCollector" {
		return AwardNoLevelUpCollector
	} else if cfg.AwardType == "weightGameRNG" {
		return AwardWeightGameRNG
	} else if cfg.AwardType == "pushSymbolCollection" {
		return AwardPushSymbolCollection
	} else if cfg.AwardType == "gameCoinMulti" {
		return AwardGameCoinMulti
	} else if cfg.AwardType == "stepCoinMulti" {
		return AwardStepCoinMulti
	} else if cfg.AwardType == "retriggerRespin" {
		return AwardRetriggerRespin
	} else if cfg.AwardType == "addRetriggerRespinNum" {
		return AwardAddRetriggerRespinNum
	} else if cfg.AwardType == "setMaskVal" {
		return AwardSetMaskVal
	} else if cfg.AwardType == "triggerRespin2" {
		return AwardTriggerRespin2
	} else if cfg.AwardType == "setComponentConfigVal" {
		return AwardSetComponentConfigVal
	} else if cfg.AwardType == "setComponentConfigIntVal" {
		return AwardSetComponentConfigIntVal
	} else if cfg.AwardType == "chgComponentConfigIntVal" {
		return AwardChgComponentConfigIntVal
	}

	return AwardUnknow
}

func (cfg *Award) Init() {
	if len(cfg.Vals) == 0 && cfg.Val > 0 {
		cfg.Vals = append(cfg.Vals, cfg.Val)
	}

	if len(cfg.StrParams) == 0 && cfg.StrParam != "" {
		cfg.StrParams = append(cfg.StrParams, cfg.StrParam)
	}

	cfg.Type = cfg.getType()
}

func (cfg *Award) GetVal(gameProp *GameProperty, i int) int {
	val := 0
	if i < len(cfg.Vals) {
		val = cfg.Vals[i]
	}

	if i < len(cfg.ComponentVals) {
		if cfg.ComponentVals[i] != "" {
			val, _ = gameProp.GetComponentVal(cfg.ComponentVals[i])
		}
	}

	return val
}

// func (cfg *Award) GetStringVal(gameProp *GameProperty, i int) string {
// 	val := ""
// 	if i < len(cfg.StrParams) {
// 		val = cfg.StrParams[i]
// 	}

// 	if i < len(cfg.ComponentVals) {
// 		if cfg.ComponentVals[i] != "" {
// 			val, _ = gameProp.GetComponentVal(cfg.ComponentVals[i])
// 		}
// 	}

// 	return val
// }

const (
	AwardUnknow                   int = 0  // 未知的奖励
	AwardCash                     int = 1  // 直接奖励cash
	AwardCollector                int = 2  // 奖励收集器
	AwardRespinTimes              int = 3  // 奖励respin次数
	AwardGameMulti                int = 4  // 奖励游戏整体倍数
	AwardStepMulti                int = 5  // 奖励这个step里的倍数
	AwardInitMask                 int = 6  // 初始化mask
	AwardTriggerRespin            int = 7  // 弃用，触发respin，理论上，在respin外面应该用AwardTriggerRespin，在respin里面应该用AwardRespinTimes，如果分不清楚，就统一用AwardTriggerRespin
	AwardNoLevelUpCollector       int = 8  // 奖励收集器，但不会触发升级奖励
	AwardWeightGameRNG            int = 9  // 权重产生一个rng，供后续逻辑用，全局用，不同step不会reset这个rng
	AwardPushSymbolCollection     int = 10 // 根据SymbolCollection自己的逻辑，产生一定数量的Symbol到SymbolCollection里
	AwardGameCoinMulti            int = 11 // 奖励游戏整体的coin倍数
	AwardStepCoinMulti            int = 12 // 奖励这个step里的coin倍数
	AwardRetriggerRespin          int = 13 // 奖励再次触发respin，这种只会用前面记录下的retrigger次数
	AwardAddRetriggerRespinNum    int = 14 // 奖励再次触发respin次数，这种会在前面的基础上增加
	AwardSetMaskVal               int = 15 // 设置mask的值
	AwardTriggerRespin2           int = 16 // 新的触发respin，不需要考虑trigger、retrigger、respinTimes，直接用这个就行，如果次数给-1，就会用当前的retriggerRespinNum
	AwardSetComponentConfigVal    int = 17 // 设置组件的configVal
	AwardSetComponentConfigIntVal int = 18 // 设置组件的configIntVal
	AwardChgComponentConfigIntVal int = 19 // 改变组件的configIntVal
)

// AwardCash
// Vals[0] - 增加的具体数值

// AwardCollector
// StrParams[0] - collector的组件名
// Vals[0] - 增加的具体数值

// AwardRespinTimes
// StrParams[0] - respin的组件名
// Vals[0] - 增加的具体数值

// AwardGameMulti
// Vals[0] - 绝对值

// AwardStepMulti
// Vals[0] - 绝对值

// AwardInitMask
// StrParams[0] - mask的组件名
// StrParams[1] - 用来初始化mask的scene

// AwardTriggerRespin
// StrParams[0] - respin的组件名
// Vals[0] - 增加的具体数值

// AwardNoLevelUpCollector
// StrParams[0] - collector的组件名
// Vals[0] - 增加的具体数值

// AwardWeightGameRNG
// StrParams[0] - 权重表
// StrParams[1] - rng名

// AwardPushSymbolCollection
// StrParams[0] - SymbolCollection的组件名
// Vals[0] - 增加的具体数值
