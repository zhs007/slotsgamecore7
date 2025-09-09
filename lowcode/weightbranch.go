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

// WeightBranchData represents runtime state for a WeightBranch component.
// It embeds BasicComponentData and stores the currently selected branch value,
// an optional per-instance ValWeights2 (WeightVW) and a list of branches
// to ignore when ForceTriggerOnce is configured.
type WeightBranchData struct {
	BasicComponentData
	Value          string
	WeightVW       *sgc7game.ValWeights2
	IgnoreBranches []string
}

// OnNewGame resets per-game state for the component data.
// It is called when a new game/play begins so cached per-data state is cleared.
func (weightBranchData *WeightBranchData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightBranchData.BasicComponentData.OnNewGame(gameProp, component)

	weightBranchData.WeightVW = nil
	weightBranchData.IgnoreBranches = nil
}

// Clone returns a copy of the WeightBranchData suitable for per-play storage.
func (weightBranchData *WeightBranchData) Clone() IComponentData {
	target := &WeightBranchData{
		BasicComponentData: weightBranchData.CloneBasicComponentData(),
		Value:              weightBranchData.Value,
	}

	return target
}

// BuildPBComponentData builds the protobuf representation of this component data.
func (weightBranchData *WeightBranchData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WeightBranchData{
		BasicComponentData: weightBranchData.BuildPBBasicComponentData(),
		Value:              weightBranchData.Value,
	}

	return pbcd
}

// GetValEx returns an integer configuration or runtime value by key.
// Currently not used by WeightBranch and returns (0,false).
func (weightBranchData *WeightBranchData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	return 0, false
}

// GetStrVal returns a string configuration or runtime value by key.
// Supports CSVValue which maps to the currently selected branch Value.
func (weightBranchData *WeightBranchData) GetStrVal(key string) (string, bool) {
	if key == CSVValue {
		return weightBranchData.Value, true
	}

	return "", false
}

// SetConfigVal applies a string configuration change to the component data.
// When the weight configuration is changed it clears any cached per-data WeightVW.
func (weightBranchData *WeightBranchData) SetConfigVal(key string, val string) {
	if key == CCVWeight {
		weightBranchData.WeightVW = nil
	}

	weightBranchData.BasicComponentData.SetConfigVal(key, val)
}

// SetConfigIntVal applies an integer configuration change to the component data.
// Special handling: when CCVClearForceTriggerOnceCache is set it clears per-data caches.
func (weightBranchData *WeightBranchData) SetConfigIntVal(key string, val int) {
	if key == CCVClearForceTriggerOnceCache {
		weightBranchData.WeightVW = nil
		weightBranchData.IgnoreBranches = nil
	} else {
		weightBranchData.BasicComponentData.SetConfigIntVal(key, val)
	}
}

// ChgConfigIntVal changes an integer config by offset and returns the new value.
// When CCVClearForceTriggerOnceCache is changed it clears per-data ignore list.
func (weightBranchData *WeightBranchData) ChgConfigIntVal(key string, off int) int {
	if key == CCVClearForceTriggerOnceCache {
		weightBranchData.IgnoreBranches = nil

		return 0
	}

	return weightBranchData.BasicComponentData.ChgConfigIntVal(key, off)
}

// BranchNode defines configuration for a single branch value.
// It contains optional Awards to trigger and the component name to jump to.
type BranchNode struct {
	Awards          []*Award `yaml:"awards" json:"awards"` // 新的奖励系统
	JumpToComponent string   `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// WeightBranchConfig is the configuration for a WeightBranch component.
// It specifies the weight table name, optional forced branch, mapping of branches
// to target components and other behavior such as ForceTriggerOnce.
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
// SetLinkComponent links a branch key (or "next") to a component name.
// Use "next" to set the default next component; otherwise the link is treated
// as a branch value and stored in MapBranchs.
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

// WeightBranch is a component that chooses a branch based on configured weights.
// It supports forced branches, player selection, and one-time force triggers.

// Init initializes the WeightBranch from a YAML file specified by fn.
// It loads the configuration and delegates to InitEx.
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

// InitEx initializes the WeightBranch from an already-parsed configuration object.
// It validates the configuration and loads the referenced weight table into Config.WeightVW.
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

func (weightBranch *WeightBranch) getForceBranch(wbd *WeightBranchData) string {
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
			if goutils.IndexOfStringSlice(wbd.IgnoreBranches, branch, 0) >= 0 {
				goutils.Error("WeightBranch.onBranch",
					slog.String("branch", branch),
					goutils.Err(ErrInvalidBranch))

				return ErrInvalidBranch
			}

			wbd.IgnoreBranches = append(wbd.IgnoreBranches, branch)

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

				for _, v := range wbd.IgnoreBranches {
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

// OnPlayGame processes a play through this component.
// It selects a branch based on configuration, player selection, or forced branch,
// executes controllers and determines the next component to jump to.
func (weightBranch *WeightBranch) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	wbd, ok := icd.(*WeightBranchData)
	if !ok {
		goutils.Error("WeightBranch.OnPlayGame:invalid icd",
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	curBetMode := int(stake.CashBet / stake.CoinBet)

	forceBranch := weightBranch.getForceBranch(wbd)
	if forceBranch == "" {
		vw2 := weightBranch.getWeight(gameProp, wbd)
		if vw2 == nil {
			goutils.Error("WeightBranch.OnPlayGame:missing weight",
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

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
		nextComponent = branch.JumpToComponent
	}

	nc := weightBranch.onStepEnd(gameProp, curpr, gp, nextComponent)

	return nc, nil
}

// ProcControllers executes branch-linked controller awards for the given string value.
func (weightBranch *WeightBranch) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	branch, isok := weightBranch.Config.MapBranchs[strVal]
	if isok {
		if len(branch.Awards) > 0 {
			gameProp.procAwards(plugin, branch.Awards, curpr, gp)
		}
	}
}

// OnAsciiGame renders debugging/ascii output for this component.
// It returns an error when icd is not the expected type.
func (weightBranch *WeightBranch) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	wbd, ok := icd.(*WeightBranchData)
	if !ok {
		goutils.Error("WeightBranch.OnAsciiGame:invalid icd",
			goutils.Err(ErrInvalidComponentData))

		return ErrInvalidComponentData
	}

	fmt.Printf("weightBranch: name=%s, value=%s\n",
		weightBranch.GetName(),
		wbd.Value)

	return nil
}

// NewComponentData creates a new empty WeightBranchData for runtime use.
func (weightBranch *WeightBranch) NewComponentData() IComponentData {
	return &WeightBranchData{}
}

// GetAllLinkComponents returns all branch keys configured for this component.
// The returned slice is not guaranteed to be in any particular order.
func (weightBranch *WeightBranch) GetAllLinkComponents() []string {
	lst := []string{}

	if weightBranch.Config.MapBranchs != nil {
		for k := range weightBranch.Config.MapBranchs {
			lst = append(lst, k)
		}
	}

	return lst
}

// GetNextLinkComponents returns all next component names referenced by branches.
// Results may include empty strings for branches without a configured jump.
func (weightBranch *WeightBranch) GetNextLinkComponents() []string {
	lst := []string{}

	if weightBranch.Config.MapBranchs != nil {
		for _, v := range weightBranch.Config.MapBranchs {
			lst = append(lst, v.JumpToComponent)
		}
	}

	return lst
}

// OnStats2 collects runtime stats for this component into the provided cache.
func (weightBranch *WeightBranch) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	weightBranch.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd, ok := icd.(*WeightBranchData)
	if !ok {
		goutils.Error("WeightBranch.OnStats2:invalid icd",
			goutils.Err(ErrInvalidComponentData))

		return
	}

	s2.ProcStatsStrVal(weightBranch.GetName(), cd.Value)
}

// NewStats2 returns a new stats feature describing this component's statistics.
func (weightBranch *WeightBranch) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptStrVal})
}

func NewWeightBranch(name string) IComponent {
	return &WeightBranch{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// NewWeightBranch constructs a new WeightBranch component with the given name.

// "configuration": {
// "weight": "greenweight"
// "forceBranch": "continue"
// isNeedPlayerSelect
// }
// jsonWeightBranch is a lightweight struct used when parsing compact JSON
// configuration for a WeightBranch from the lowcode format.
type jsonWeightBranch struct {
	Weight             string   `json:"weight"`
	ForceBranch        string   `json:"forceBranch"`
	ForceTriggerOnce   []string `json:"forceTriggerOnce"`
	IsNeedPlayerSelect bool     `json:"isNeedPlayerSelect"`
}

// build converts the parsed jsonWeightBranch into a full WeightBranchConfig.
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

// parseWeightBranch parses a weightBranch configuration from the AST node used
// by the lowcode parser and registers the component configuration into gamecfg.
// It returns the generated component label on success.
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
