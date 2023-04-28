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
	AwardUnknow             int = 0
	AwardCash               int = 1
	AwardCollector          int = 2
	AwardRespinTimes        int = 3
	AwardGameMulti          int = 4
	AwardStepMulti          int = 5
	AwardInitMask           int = 6
	AwardTriggerRespin      int = 7
	AwardNoLevelUpCollector int = 8
)

// type Award struct {
// 	AwardType int
// 	Config    *AwardConfig
// }

// func NewArard(cfg *AwardConfig) *Award {
// 	if len(cfg.Vals) == 0 && cfg.Val > 0 {
// 		cfg.Vals = append(cfg.Vals, cfg.Val)
// 	}

// 	if len(cfg.StrParams) == 0 && cfg.StrParam != "" {
// 		cfg.StrParams = append(cfg.StrParams, cfg.StrParam)
// 	}

// 	return &Award{
// 		AwardType: cfg.GetType(),
// 		Config:    cfg,
// 	}
// }
