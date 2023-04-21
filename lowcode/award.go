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
	}

	return AwardUnknow
}

const (
	AwardUnknow      int = 0
	AwardCash        int = 1
	AwardCollector   int = 2
	AwardRespinTimes int = 3
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
