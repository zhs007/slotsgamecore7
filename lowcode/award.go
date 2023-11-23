package lowcode

type Award struct {
	AwardType string   `yaml:"awardType" json:"awardType"`
	Type      int      `yaml:"-" json:"-"`
	Val       int      `yaml:"val" json:"-"`      // 弃用，代码里已经不用了，初始化时会把数据转存到Vals里，为了兼容性保留配置
	StrParam  string   `yaml:"strParam" json:"-"` // 弃用，代码里已经不用了，初始化时会把数据转存到StrParams里，为了兼容性保留配置
	Vals      []int    `yaml:"vals" json:"vals"`
	StrParams []string `yaml:"strParams" json:"strParams"`
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

const (
	AwardUnknow               int = 0  // 未知的奖励
	AwardCash                 int = 1  // 直接奖励cash
	AwardCollector            int = 2  // 奖励收集器
	AwardRespinTimes          int = 3  // 奖励respin次数
	AwardGameMulti            int = 4  // 奖励游戏整体倍数
	AwardStepMulti            int = 5  // 奖励这个step里的倍数
	AwardInitMask             int = 6  // 初始化mask
	AwardTriggerRespin        int = 7  // 触发respin，理论上，在respin外面应该用AwardTriggerRespin，在respin里面应该用AwardRespinTimes，如果分不清楚，就统一用AwardTriggerRespin
	AwardNoLevelUpCollector   int = 8  // 奖励收集器，但不会触发升级奖励
	AwardWeightGameRNG        int = 9  // 权重产生一个rng，供后续逻辑用，全局用，不同step不会reset这个rng
	AwardPushSymbolCollection int = 10 // 根据SymbolCollection自己的逻辑，产生一定数量的Symbol到SymbolCollection里
	AwardGameCoinMulti        int = 11 // 奖励游戏整体的coin倍数
	AwardStepCoinMulti        int = 12 // 奖励这个step里的coin倍数
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
