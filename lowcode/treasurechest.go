package lowcode

import (
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

const TreasureChestTypeName = "treasureChest"

type TreasureChestType int

const (
	TreasureChestTypeFragmentCollection TreasureChestType = 0
)

func parseTreasureChestType(str string) TreasureChestType {
	if str == "fragmentcollection" {
		return TreasureChestTypeFragmentCollection
	}

	return TreasureChestTypeFragmentCollection
}

type TreasureChestData struct {
	BasicComponentData
	Selected []int
}

// OnNewGame -
func (treasureChestData *TreasureChestData) OnNewGame(gameProp *GameProperty, component IComponent) {
	treasureChestData.BasicComponentData.OnNewGame(gameProp, component)

	treasureChestData.Selected = nil
}

// onNewStep -
func (treasureChestData *TreasureChestData) onNewStep(gameProp *GameProperty, component IComponent) {
	treasureChestData.Selected = nil
}

// Clone
func (treasureChestData *TreasureChestData) Clone() IComponentData {
	target := &TreasureChestData{
		BasicComponentData: treasureChestData.CloneBasicComponentData(),
		Selected:           make([]int, len(treasureChestData.Selected)),
	}

	copy(target.Selected, treasureChestData.Selected)

	return target
}

// BuildPBComponentData
func (treasureChestData *TreasureChestData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.TreasureChestData{
		BasicComponentData: treasureChestData.BuildPBBasicComponentData(),
		Selected:           make([]int32, len(treasureChestData.Selected)),
	}

	for i, v := range treasureChestData.Selected {
		pbcd.Selected[i] = int32(v)
	}

	return pbcd
}

// TreasureChestConfig - configuration for TreasureChest
type TreasureChestConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`
	Type                 TreasureChestType     `yaml:"-" json:"-"`
	StrWeight            string                `yaml:"weight" json:"weight"` // weight
	Weight               *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	FragmentNum          int                   `yaml:"fragmentNum" json:"fragmentNum"` // fragmentNum
	MapBranchs           map[int]string        `yaml:"mapBranchs" json:"mapBranchs"`
	MapControllers       map[int][]*Award      `yaml:"mapControllers" json:"mapControllers"`
}

// SetLinkComponent
func (cfg *TreasureChestConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else {
		v, err := goutils.String2Int64(link)
		if err == nil {
			if cfg.MapBranchs == nil {
				cfg.MapBranchs = make(map[int]string)
			}

			cfg.MapBranchs[int(v)] = componentName
		}
	}
}

type TreasureChest struct {
	*BasicComponent `json:"-"`
	Config          *TreasureChestConfig `json:"config"`
}

// Init -
func (treasureChest *TreasureChest) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("TreasureChest.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &TreasureChestConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("TreasureChest.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return treasureChest.InitEx(cfg, pool)
}

// InitEx -
func (treasureChest *TreasureChest) InitEx(cfg any, pool *GamePropertyPool) error {
	treasureChest.Config = cfg.(*TreasureChestConfig)
	treasureChest.Config.ComponentType = TreasureChestTypeName

	treasureChest.Config.Type = parseTreasureChestType(treasureChest.Config.StrType)

	vw2, err := pool.LoadIntWeights(treasureChest.Config.StrWeight, true)
	if err != nil {
		goutils.Error("TreasureChest.InitEx:LoadIntWeights",
			slog.String("weight", treasureChest.Config.StrWeight),
			goutils.Err(err))

		return err
	}

	treasureChest.Config.Weight = vw2

	if treasureChest.Config.FragmentNum <= 0 {
		goutils.Error("TreasureChest.InitEx:FragmentNum",
			slog.Int("FragmentNum", treasureChest.Config.FragmentNum),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	for _, awards := range treasureChest.Config.MapControllers {
		for _, award := range awards {
			award.Init()
		}
	}

	treasureChest.onInit(&treasureChest.Config.BasicComponentConfig)

	return nil
}

// playgame
func (treasureChest *TreasureChest) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*TreasureChestData)
	cd.onNewStep(gameProp, treasureChest)

	vw2 := treasureChest.getWeight(gameProp, &cd.BasicComponentData)
	if vw2 == nil {
		goutils.Error("TreasureChest.OnPlayGame:getWeight",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	cd.Selected = make([]int, 0, len(vw2.Vals)*treasureChest.Config.FragmentNum)
	mapVals := make(map[int]int)
	for _, v := range vw2.Vals {
		mapVals[v.Int()] = treasureChest.Config.FragmentNum
	}

	for {
		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("TreasureChest.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		cv := cr.Int()
		cd.Selected = append(cd.Selected, cv)

		mapVals[cv]--
		if mapVals[cv] <= 0 {
			cd.Output = cv

			break
		}
	}

	// winMulti := jackpot.GetWinMulti(&cd.BasicComponentData)
	// wins := jackpot.GetWins(&cd.BasicComponentData)

	// cd.WinMulti = winMulti

	// cd.Wins = wins * winMulti

	// bet := gameProp.GetBet3(stake, jackpot.Config.BetType)

	// ret := &sgc7game.Result{
	// 	Symbol:    -1,
	// 	Type:      sgc7game.RTBonus,
	// 	LineIndex: -1,
	// 	CoinWin:   cd.Wins,
	// 	CashWin:   cd.Wins * bet,
	// }

	// jackpot.AddResult(curpr, ret, &cd.BasicComponentData)

	nc := treasureChest.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (treasureChest *TreasureChest) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData -
func (treasureChest *TreasureChest) NewComponentData() IComponentData {
	return &TreasureChestData{}
}

func (treasureChest *TreasureChest) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return treasureChest.Config.Weight
}

// OnStats2
func (treasureChest *TreasureChest) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	treasureChest.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*TreasureChestData)

	s2.ProcStatsIntVal(treasureChest.GetName(), cd.Output)
}

// NewStats2 -
func (treasureChest *TreasureChest) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptIntVal})
}

func NewTreasureChest(name string) IComponent {
	return &TreasureChest{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "fragmentCollection",
// "fragmentNum": 3,
// "weight": "weight_coin"
type jsonTreasureChest struct {
	Type        string `json:"type"`
	FragmentNum int    `json:"fragmentNum"`
	Weight      string `json:"weight"`
}

func (jcfg *jsonTreasureChest) build() *TreasureChestConfig {
	cfg := &TreasureChestConfig{
		StrType:     jcfg.Type,
		StrWeight:   jcfg.Weight,
		FragmentNum: jcfg.FragmentNum,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseTreasureChest(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseTreasureChest:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseTreasureChest:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonTreasureChest{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseTreasureChest:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseFeatureBarControllers(ctrls)
		if err != nil {
			goutils.Error("parseBasicReels:parseMapControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: TreasureChestTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
