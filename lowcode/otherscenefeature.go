package lowcode

type OtherSceneFeatureConfig struct {
	Type string `yaml:"type"`
}

func (cfg *OtherSceneFeatureConfig) GetType() int {
	if cfg.Type == "gameMulti" {
		return OtherSceneFeatureGameMulti
	} else if cfg.Type == "gameMultiSum" {
		return OtherSceneFeatureGameMultiSum
	} else if cfg.Type == "stepMulti" {
		return OtherSceneFeatureStepMulti
	} else if cfg.Type == "stepMultiSum" {
		return OtherSceneFeatureStepMultiSum
	}

	return OtherSceneFeatureUnknow
}

const (
	OtherSceneFeatureUnknow       int = 0
	OtherSceneFeatureGameMulti    int = 1 // GameMulti，默认用乘法
	OtherSceneFeatureGameMultiSum int = 2 // GameMulti，默认用加法
	OtherSceneFeatureStepMulti    int = 3 // StepMulti，默认用乘法
	OtherSceneFeatureStepMultiSum int = 4 // StepMulti，默认用加法
)

type OtherSceneFeature struct {
	Type   int
	Config *OtherSceneFeatureConfig
}

func NewOtherSceneFeature(cfg *OtherSceneFeatureConfig) *OtherSceneFeature {
	return &OtherSceneFeature{
		Type:   cfg.GetType(),
		Config: cfg,
	}
}
