package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WeightBranchTypeName = "weightBranch"

const (
	WBDVValue string = "value" // 权重表最终的value
)

type WeightBranchData struct {
	BasicComponentData
	Value string
}

// OnNewGame -
func (weightBranchData *WeightBranchData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightBranchData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (weightBranchData *WeightBranchData) OnNewStep(gameProp *GameProperty, component IComponent) {
	weightBranchData.BasicComponentData.OnNewStep(gameProp, component)
}

// BuildPBComponentData
func (weightBranchData *WeightBranchData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WeightBranchData{
		BasicComponentData: weightBranchData.BuildPBBasicComponentData(),
		Value:              weightBranchData.Value,
	}

	return pbcd
}

// GetVal -
func (weightBranchData *WeightBranchData) GetVal(key string) int {
	return 0
}

// SetVal -
func (weightBranchData *WeightBranchData) SetVal(key string, val int) {
}

// BranchNode -
type BranchNode struct {
	Awards          []*Award `yaml:"awards" json:"awards"` // 新的奖励系统
	JumpToComponent string   `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// WeightBranchConfig - configuration for WeightBranch
type WeightBranchConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Weight               string                 `yaml:"weight" json:"weight"`
	WeightVW             *sgc7game.ValWeights2  `json:"-"`
	MapBranchs           map[string]*BranchNode `yaml:"mapBranchs" json:"mapBranchs"` // 可以不用配置全，如果没有配置的，就跳转默认的next
}

type WeightBranch struct {
	*BasicComponent `json:"-"`
	Config          *WeightBranchConfig `json:"config"`
}

// Init -
func (weightBranch *WeightBranch) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightBranch.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WeightBranchConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightBranch.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return weightBranch.InitEx(cfg, pool)
}

// InitEx -
func (weightBranch *WeightBranch) InitEx(cfg any, pool *GamePropertyPool) error {
	weightBranch.Config = cfg.(*WeightBranchConfig)
	weightBranch.Config.ComponentType = WeightBranchTypeName

	if weightBranch.Config.Weight != "" {
		vw2, err := pool.LoadStrWeights(weightBranch.Config.Weight, weightBranch.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightBranch.Init:LoadStrWeights",
				zap.String("Weight", weightBranch.Config.Weight),
				zap.Error(err))

			return err
		}

		weightBranch.Config.WeightVW = vw2
	} else {
		goutils.Error("WeightBranch.InitEx:Weight",
			zap.Error(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	for _, node := range weightBranch.Config.MapBranchs {
		for _, award := range node.Awards {
			award.Init()
		}
	}

	weightBranch.onInit(&weightBranch.Config.BasicComponentConfig)

	return nil
}

// playgame
func (weightBranch *WeightBranch) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// weightBranch.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	wbd := icd.(*WeightBranchData)

	cr, err := weightBranch.Config.WeightVW.RandVal(plugin)
	if err != nil {
		goutils.Error("WeightBranch.OnPlayGame:RandVal",
			zap.Error(err))

		return "", err
	}

	wbd.Value = cr.String()

	nextComponent := ""

	branch, isok := weightBranch.Config.MapBranchs[wbd.Value]
	if isok {
		if len(branch.Awards) > 0 {
			gameProp.procAwards(plugin, branch.Awards, curpr, gp)
		}

		nextComponent = branch.JumpToComponent
	}

	nc := weightBranch.onStepEnd(gameProp, curpr, gp, nextComponent)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (weightBranch *WeightBranch) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	wbd := icd.(*WeightBranchData)

	fmt.Printf("weightBranch %v, got %v\n", weightBranch.GetName(), wbd.Value)

	return nil
}

// OnStats
func (weightBranch *WeightBranch) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewWeightBranch(name string) IComponent {
	return &WeightBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
