package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const MaskBranchTypeName = "maskBranch"

// MaskBranchNode -
type MaskBranchNode struct {
	MaskVal         []bool   `yaml:"mask" json:"mask"`
	Awards          []*Award `yaml:"awards" json:"awards"` // 新的奖励系统
	JumpToComponent string   `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// MaskBranchConfig - configuration for MaskBranch
type MaskBranchConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Mask                 string            `yaml:"mask" json:"mask"`   // mask
	Nodes                []*MaskBranchNode `yaml:"nodes" json:"nodes"` // 可以不用配置全，如果没有配置的，就跳转默认的next
}

type MaskBranch struct {
	*BasicComponent `json:"-"`
	Config          *MaskBranchConfig `json:"config"`
}

// Init -
func (maskBranch *MaskBranch) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MaskBranch.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MaskBranchConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MaskBranch.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return maskBranch.InitEx(cfg, pool)
}

// InitEx -
func (maskBranch *MaskBranch) InitEx(cfg any, pool *GamePropertyPool) error {
	maskBranch.Config = cfg.(*MaskBranchConfig)
	maskBranch.Config.ComponentType = MaskBranchTypeName

	for _, node := range maskBranch.Config.Nodes {
		for _, award := range node.Awards {
			award.Init()
		}
	}

	maskBranch.onInit(&maskBranch.Config.BasicComponentConfig)

	return nil
}

// playgame
func (maskBranch *MaskBranch) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	maskBranch.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	// cd := gameProp.MapComponentData[maskBranch.Name].(*BasicComponentData)

	maskdata, err := gameProp.Pool.GetMask(maskBranch.Config.Mask, gameProp)
	if err != nil {
		goutils.Error("MaskBranch.OnPlayGame:GetMask",
			zap.Error(err))

		return "", err
	}

	nextComponent := ""

	for _, n := range maskBranch.Config.Nodes {
		if isSameBoolSlice(n.MaskVal, maskdata) {
			if len(n.Awards) > 0 {
				gameProp.procAwards(plugin, n.Awards, curpr, gp)
			}

			nextComponent = n.JumpToComponent

			break
		}
	}

	nc := maskBranch.onStepEnd(gameProp, curpr, gp, nextComponent)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (maskBranch *MaskBranch) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	maskdata, err := gameProp.Pool.GetMask(maskBranch.Config.Mask, gameProp)
	if err != nil {
		goutils.Error("MaskBranch.OnPlayGame:GetMask",
			zap.Error(err))

		return err
	}

	if maskdata != nil {
		fmt.Printf("MaskBranch %v, got %v is %#v", maskBranch.GetName(), maskBranch.Config.Mask, maskdata)
	}

	return nil
}

// OnStats
func (maskBranch *MaskBranch) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMaskBranch(name string) IComponent {
	return &MaskBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
