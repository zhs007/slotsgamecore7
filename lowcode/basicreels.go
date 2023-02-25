package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// TriggerFeatureConfig - configuration for trigger feature
type TriggerFeatureConfig struct {
	Symbol  string   `yaml:"symbol"`  // like scatter
	Type    string   `yaml:"type"`    // like scatters
	Scripts []string `yaml:"scripts"` // scripts
}

// BasicReelsConfig - configuration for BasicReels
type BasicReelsConfig struct {
	MainType       string                  `yaml:"mainType"`       // lines or ways
	ExcludeSymbols []string                `yaml:"excludeSymbols"` // w/s etc
	ReelSetsWeight string                  `yaml:"reelSetWeight"`
	MysteryWeight  string                  `yaml:"mysteryWeight"`
	BeforMain      []*TriggerFeatureConfig `yaml:"beforMain"` // befor the maintype
	AfterMain      []*TriggerFeatureConfig `yaml:"afterMain"` // after the maintype
}

type BasicReels struct {
}

// playgame
func (basicReels *BasicReels) OnPlayGame(curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	return nil
}

// pay
func (basicReels *BasicReels) OnPay(curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	return nil
}

func NewBasicReels(fn string) IComponent {
	return &BasicReels{}
}
