package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WeightBranchTypeName = "weightBranch"

type WeightBranchData struct {
	BasicComponentData
	Value string
}

// OnNewGame -
func (weightBranchData *WeightBranchData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightBranchData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (weightBranchData *WeightBranchData) Clone() IComponentData {
	target := &WeightBranchData{
		BasicComponentData: weightBranchData.CloneBasicComponentData(),
		Value:              weightBranchData.Value,
	}

	return target
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
func (weightBranchData *WeightBranchData) GetVal(key string) (int, bool) {
	return 0, false
}

// // SetVal -
// func (weightBranchData *WeightBranchData) SetVal(key string, val int) {
// }

// GetStrVal -
func (weightBranchData *WeightBranchData) GetStrVal(key string) (string, bool) {
	if key == CSVValue {
		return weightBranchData.Value, true
	}

	return "", false
}

// BranchNode -
type BranchNode struct {
	Awards          []*Award `yaml:"awards" json:"awards"` // 新的奖励系统
	JumpToComponent string   `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// WeightBranchConfig - configuration for WeightBranch
type WeightBranchConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ForceBranch          string                 `yaml:"forceBranch" json:"forceBranch"`
	Weight               string                 `yaml:"weight" json:"weight"`
	WeightVW             *sgc7game.ValWeights2  `json:"-"`
	MapBranchs           map[string]*BranchNode `yaml:"mapBranchs" json:"mapBranchs"` // 可以不用配置全，如果没有配置的，就跳转默认的next
}

// // SetLinkComponent
// func (cfg *WeightBranchConfig) hasValWeight(val string) bool {
// 	for _, v := range cfg.WeightVW.Vals {
// 		if v.String() == val {
// 			return true
// 		}
// 	}

// 	return false
// }

// SetLinkComponent
func (cfg *WeightBranchConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else {
		if cfg.MapBranchs == nil {
			cfg.MapBranchs = make(map[string]*BranchNode)
		}

		if cfg.MapBranchs[link] == nil {
			cfg.MapBranchs[link] = &BranchNode{
				JumpToComponent: componentName,
			}
		} else {
			cfg.MapBranchs[link].JumpToComponent = componentName
		}
	}

	// if cfg.hasValWeight(link) {
	// }
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WeightBranchConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightBranch.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
				slog.String("Weight", weightBranch.Config.Weight),
				goutils.Err(err))

			return err
		}

		weightBranch.Config.WeightVW = vw2
	} else {
		goutils.Error("WeightBranch.InitEx:Weight",
			goutils.Err(ErrIvalidComponentConfig))

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

func (weightBranch *WeightBranch) getForceBrach(wbd *WeightBranchData) string {
	val := wbd.BasicComponentData.GetConfigVal(CCVForceBranch)
	if val != "" {
		return val
	}

	return weightBranch.Config.ForceBranch
}

func (weightBranch *WeightBranch) getWeight(gameProp *GameProperty, wbd *WeightBranchData) *sgc7game.ValWeights2 {
	val := wbd.BasicComponentData.GetConfigVal(CCVWeight)
	if val != "" {
		vw2, err := gameProp.Pool.LoadStrWeights(val, weightBranch.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightBranch.getWeight:LoadStrWeights",
				slog.String("Weight", val),
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return weightBranch.Config.WeightVW
}

// playgame
func (weightBranch *WeightBranch) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// weightBranch.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	wbd := icd.(*WeightBranchData)

	curBetMode := int(stake.CashBet / stake.CoinBet)

	forceBranch := weightBranch.getForceBrach(wbd)
	if forceBranch == "" {
		vw2 := weightBranch.getWeight(gameProp, wbd)
		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("WeightBranch.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		wbd.Value = cr.String()
	} else {
		wbd.Value = forceBranch
	}

	if gameProp.rng != nil {
		gameProp.rng.OnChoiceBranch(curBetMode, weightBranch, wbd.Value)
	}

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

// // OnStats
// func (weightBranch *WeightBranch) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (weightBranch *WeightBranch) NewComponentData() IComponentData {
	return &WeightBranchData{}
}

// GetAllLinkComponents - get all link components
func (weightBranch *WeightBranch) GetAllLinkComponents() []string {
	lst := []string{}

	if weightBranch.Config.MapBranchs != nil {
		for k := range weightBranch.Config.MapBranchs {
			lst = append(lst, k)
		}
	}

	return lst
}

// GetNextLinkComponents - get next link components
func (weightBranch *WeightBranch) GetNextLinkComponents() []string {
	lst := []string{}

	if weightBranch.Config.MapBranchs != nil {
		for _, v := range weightBranch.Config.MapBranchs {
			lst = append(lst, v.JumpToComponent)
		}
	}

	return lst
}

func NewWeightBranch(name string) IComponent {
	return &WeightBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "configuration": {
// "weight": "greenweight"
// "forceBranch": "continue"
// }
type jsonWeightBranch struct {
	Weight      string `json:"weight"`
	ForceBranch string `json:"forceBranch"`
}

func (jwr *jsonWeightBranch) build() *WeightBranchConfig {
	cfg := &WeightBranchConfig{
		Weight:      jwr.Weight,
		ForceBranch: jwr.ForceBranch,
		MapBranchs:  make(map[string]*BranchNode),
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWeightBranch(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("WeightBranch:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("WeightBranch:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWeightBranch{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("WeightBranch:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseMapControllers(ctrls)
		if err != nil {
			goutils.Error("parseBasicReels:parseMapControllers",
				goutils.Err(err))

			return "", err
		}

		for k, arr := range mapAwards {
			if cfgd.MapBranchs[k] == nil {
				cfgd.MapBranchs[k] = &BranchNode{
					Awards: arr,
				}
			} else {
				cfgd.MapBranchs[k].Awards = arr
			}
		}

		// cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WeightBranchTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
