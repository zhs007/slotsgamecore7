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
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WeightBranchTypeName = "weightBranch"

type WeightBranchData struct {
	BasicComponentData
	Value         string
	WeightVW      *sgc7game.ValWeights2
	IgnoreBranchs []string
}

// OnNewGame -
func (weightBranchData *WeightBranchData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightBranchData.BasicComponentData.OnNewGame(gameProp, component)

	weightBranchData.WeightVW = nil
	weightBranchData.IgnoreBranchs = nil
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

// GetValEx -
func (weightBranchData *WeightBranchData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	return 0, false
}

// GetStrVal -
func (weightBranchData *WeightBranchData) GetStrVal(key string) (string, bool) {
	if key == CSVValue {
		return weightBranchData.Value, true
	}

	return "", false
}

// SetConfigVal -
func (weightBranchData *WeightBranchData) SetConfigVal(key string, val string) {
	if key == CCVWeight {
		weightBranchData.WeightVW = nil
	}

	weightBranchData.BasicComponentData.SetConfigVal(key, val)
}

// SetConfigIntVal - CCVValueNum的set和chg逻辑不太一样，等于的时候不会触发任何的 controllers
func (weightBranchData *WeightBranchData) SetConfigIntVal(key string, val int) {
	if key == CCVClearForceTriggerOnceCache {
		weightBranchData.WeightVW = nil
		weightBranchData.IgnoreBranchs = nil
	} else {
		weightBranchData.BasicComponentData.SetConfigIntVal(key, val)
	}
}

// ChgConfigIntVal -
func (weightBranchData *WeightBranchData) ChgConfigIntVal(key string, off int) int {
	if key == CCVClearForceTriggerOnceCache {
		weightBranchData.IgnoreBranchs = nil

		return 0
	}

	return weightBranchData.BasicComponentData.ChgConfigIntVal(key, off)
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
	ForceTriggerOnce     []string               `yaml:"forceTriggerOnce" json:"forceTriggerOnce"`
	IsNeedPlayerSelect   bool                   `yaml:"isNeedPlayerSelect" json:"isNeedPlayerSelect"`
}

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
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
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
	if wbd.WeightVW != nil {
		return wbd.WeightVW
	}

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

func (weightBranch *WeightBranch) onBranch(branch string, wbd *WeightBranchData, vw2 *sgc7game.ValWeights2) error {
	if len(weightBranch.Config.ForceTriggerOnce) > 0 {
		if goutils.IndexOfStringSlice(weightBranch.Config.ForceTriggerOnce, branch, 0) >= 0 {
			if goutils.IndexOfStringSlice(wbd.IgnoreBranchs, branch, 0) >= 0 {
				goutils.Error("WeightBranch.onBranch",
					slog.String("branch", branch),
					goutils.Err(ErrInvalidBranch))

				return ErrInvalidBranch
			}

			wbd.IgnoreBranchs = append(wbd.IgnoreBranchs, branch)

			if wbd.WeightVW != nil {
				nvw2, err := wbd.WeightVW.CloneExcludeVal(sgc7game.NewStrValEx(branch))
				if err != nil {
					goutils.Error("WeightBranch.CloneExcludeVal",
						slog.String("branch", branch),
						goutils.Err(err))

					return err
				}

				wbd.WeightVW = nvw2
			} else {
				nvw2 := vw2.Clone()

				for _, v := range wbd.IgnoreBranchs {
					err := nvw2.RemoveVal(sgc7game.NewStrValEx(v))
					if err != nil {
						goutils.Error("WeightBranch.RemoveVal",
							slog.String("branch", branch),
							goutils.Err(err))

						return err
					}
				}

				wbd.WeightVW = nvw2
			}
		}
	}

	return nil
}

// playgame
func (weightBranch *WeightBranch) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	wbd := icd.(*WeightBranchData)

	curBetMode := int(stake.CashBet / stake.CoinBet)

	forceBranch := weightBranch.getForceBrach(wbd)
	if forceBranch == "" {
		vw2 := weightBranch.getWeight(gameProp, wbd)

		if weightBranch.Config.IsNeedPlayerSelect {
			if cmd == DefaultCmd {
				lstcmd := []string{}
				lstparam := []string{}

				for i, v := range vw2.Vals {
					if vw2.Weights[i] > 0 {
						lstcmd = append(lstcmd, weightBranch.Name)
						lstparam = append(lstparam, v.String())
					}
				}

				curpr.NextCmds = lstcmd
				curpr.NextCmdParams = lstparam
				curpr.IsFinish = false
				curpr.IsWait = true

				nc := weightBranch.onStepEnd(gameProp, curpr, gp, "")

				return nc, nil
			} else if cmd == weightBranch.Name {
				isSelectOK := false
				for i, v := range vw2.Vals {
					w := vw2.Weights[i]

					if w > 0 && param == v.String() {
						wbd.Value = param

						isSelectOK = true

						break
					}
				}

				if !isSelectOK {
					goutils.Error("WeightBranch.OnPlayGame:IsNeedPlayerSelect",
						slog.String("branch", param),
						goutils.Err(ErrInvalidBranch))

					return "", ErrInvalidBranch
				}
			} else {
				goutils.Error("WeightBranch.OnPlayGame:IsNeedPlayerSelect",
					slog.String("cmd", cmd),
					goutils.Err(ErrInvalidCommand))

				return "", ErrInvalidCommand
			}
		} else {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("WeightBranch.OnPlayGame:RandVal",
					goutils.Err(err))

				return "", err
			}

			wbd.Value = cr.String()

			weightBranch.onBranch(wbd.Value, wbd, vw2)
		}
	} else {
		wbd.Value = forceBranch
	}

	if gameProp.rng != nil {
		gameProp.rng.OnChoiceBranch(curBetMode, weightBranch, wbd.Value)
	}

	nextComponent := ""

	weightBranch.ProcControllers(gameProp, plugin, curpr, gp, -1, "<any>")
	weightBranch.ProcControllers(gameProp, plugin, curpr, gp, -1, wbd.Value)

	branch, isok := weightBranch.Config.MapBranchs[wbd.Value]
	if isok {
		// if len(branch.Awards) > 0 {
		// 	gameProp.procAwards(plugin, branch.Awards, curpr, gp)
		// }

		nextComponent = branch.JumpToComponent
	}

	nc := weightBranch.onStepEnd(gameProp, curpr, gp, nextComponent)

	return nc, nil
}

// OnProcControllers -
func (weightBranch *WeightBranch) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	branch, isok := weightBranch.Config.MapBranchs[strVal]
	if isok {
		if len(branch.Awards) > 0 {
			gameProp.procAwards(plugin, branch.Awards, curpr, gp)
		}
	}
}

// OnAsciiGame - outpur to asciigame
func (weightBranch *WeightBranch) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	wbd := icd.(*WeightBranchData)

	fmt.Printf("weightBranch %v, got %v\n", weightBranch.GetName(), wbd.Value)

	return nil
}

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

// OnStats2
func (weightBranch *WeightBranch) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	weightBranch.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*WeightBranchData)

	s2.ProcStatsStrVal(weightBranch.GetName(), cd.Value)
}

// NewStats2 -
func (weightBranch *WeightBranch) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptStrVal})
}

func NewWeightBranch(name string) IComponent {
	return &WeightBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "configuration": {
// "weight": "greenweight"
// "forceBranch": "continue"
// isNeedPlayerSelect
// }
type jsonWeightBranch struct {
	Weight             string   `json:"weight"`
	ForceBranch        string   `json:"forceBranch"`
	ForceTriggerOnce   []string `json:"forceTriggerOnce"`
	IsNeedPlayerSelect bool     `json:"isNeedPlayerSelect"`
}

func (jwr *jsonWeightBranch) build() *WeightBranchConfig {
	cfg := &WeightBranchConfig{
		Weight:             jwr.Weight,
		ForceBranch:        jwr.ForceBranch,
		ForceTriggerOnce:   jwr.ForceTriggerOnce,
		MapBranchs:         make(map[string]*BranchNode),
		IsNeedPlayerSelect: jwr.IsNeedPlayerSelect,
	}

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
