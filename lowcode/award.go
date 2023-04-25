package lowcode

type AwardConfig struct {
	AwardType string `yaml:"awardType"`
	Val       int    `yaml:"val"`
	StrParam  string `yaml:"strParam"`
}

func (cfg *AwardConfig) GetType() int {
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
	}

	return AwardUnknow
}

const (
	AwardUnknow      int = 0
	AwardCash        int = 1
	AwardCollector   int = 2
	AwardRespinTimes int = 3
	AwardGameMulti   int = 4
	AwardStepMulti   int = 5
	AwardInitMask    int = 6
)

type Award struct {
	AwardType int
	Config    *AwardConfig
}

func NewArard(cfg *AwardConfig) *Award {
	return &Award{
		AwardType: cfg.GetType(),
		Config:    cfg,
	}
}
