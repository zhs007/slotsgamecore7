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

const RollNumberTypeName = "rollNumber"

type RollNumberData struct {
	BasicComponentData
	Number int
}

// OnNewGame -
func (rollNumberData *RollNumberData) OnNewGame(gameProp *GameProperty, component IComponent) {
	rollNumberData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (rollNumberData *RollNumberData) Clone() IComponentData {
	target := &RollNumberData{
		BasicComponentData: rollNumberData.CloneBasicComponentData(),
		Number:             rollNumberData.Number,
	}

	return target
}

// BuildPBComponentData
func (rollNumberData *RollNumberData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RollNumberData{
		BasicComponentData: rollNumberData.BuildPBBasicComponentData(),
		Number:             int32(rollNumberData.Number),
	}

	return pbcd
}

// GetValEx -
func (rollNumberData *RollNumberData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVNumber || key == CVOutputInt {
		return rollNumberData.Number, true
	}

	return 0, false
}

// RollNumberConfig - configuration for RollNumber
type RollNumberConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Weight               string                `yaml:"weight" json:"weight"`
	WeightVW             *sgc7game.ValWeights2 `json:"-"`
	Awards               []*Award              `yaml:"awards" json:"awards"`             // 新的奖励系统
	MapValAwards         map[int][]*Award      `yaml:"mapValAwards" json:"mapValAwards"` // 新的奖励系统
	ForceVal             int                   `yaml:"forceVal" json:"forceVal"`
}

// SetLinkComponent
func (cfg *RollNumberConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RollNumber struct {
	*BasicComponent `json:"-"`
	Config          *RollNumberConfig `json:"config"`
}

// Init -
func (rollNumber *RollNumber) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RollNumber.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RollNumberConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RollNumber.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return rollNumber.InitEx(cfg, pool)
}

// InitEx -
func (rollNumber *RollNumber) InitEx(cfg any, pool *GamePropertyPool) error {
	rollNumber.Config = cfg.(*RollNumberConfig)
	rollNumber.Config.ComponentType = RollNumberTypeName

	if rollNumber.Config.Weight != "" {
		vw2, err := pool.LoadSymbolWeights(rollNumber.Config.Weight, "val", "weight", pool.DefaultPaytables, rollNumber.Config.UseFileMapping)
		if err != nil {
			goutils.Error("RollNumber.Init:LoadStrWeights",
				slog.String("Weight", rollNumber.Config.Weight),
				goutils.Err(err))

			return err
		}

		rollNumber.Config.WeightVW = vw2
	} else {
		goutils.Error("RollNumber.InitEx:Weight",
			goutils.Err(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	for _, award := range rollNumber.Config.Awards {
		award.Init()
	}

	rollNumber.onInit(&rollNumber.Config.BasicComponentConfig)

	return nil
}

func (rollNumber *RollNumber) getForceVal(basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVForceVal)
	if isok && v != -1 {
		return v
	}

	v, isok = basicCD.GetConfigIntVal(CCVForceValNow)
	if isok && v != -1 {
		return v
	}

	if rollNumber.Config.ForceVal != -1 {
		return rollNumber.Config.ForceVal
	}

	return -1
}

func (rollNumber *RollNumber) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	val := basicCD.GetConfigVal(CCVWeight)
	if val != "" {
		vw2, err := gameProp.Pool.LoadStrWeights(val, rollNumber.Config.UseFileMapping)
		if err != nil {
			goutils.Error("RollNumber.getWeight:LoadStrWeights",
				slog.String("Weight", val),
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return rollNumber.Config.WeightVW
}

// playgame
func (rollNumber *RollNumber) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	rnd := icd.(*RollNumberData)

	rnd.Number = 0

	forceVal := rollNumber.getForceVal(&rnd.BasicComponentData)
	if forceVal == -1 {
		vw := rollNumber.getWeight(gameProp, &rnd.BasicComponentData)
		if vw.MaxWeight == 0 {
			nc := rollNumber.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		cr, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("RollNumber.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		rnd.Number = cr.Int()
	} else {
		rnd.Number = forceVal
	}

	if len(rollNumber.Config.Awards) > 0 {
		gameProp.procAwards(plugin, rollNumber.Config.Awards, curpr, gp)
	}

	awards, isok := rollNumber.Config.MapValAwards[rnd.Number]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}

	nc := rollNumber.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (rollNumber *RollNumber) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	rsd := icd.(*RollNumberData)

	fmt.Printf("rollNumber %v, got %v\n", rollNumber.GetName(), rsd.Number)

	return nil
}

// NewComponentData -
func (rollNumber *RollNumber) NewComponentData() IComponentData {
	return &RollNumberData{}
}

func NewRollNumber(name string) IComponent {
	return &RollNumber{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"forceVal": -1,
//		"weight": "fgbookofsymbol",
//	},
type jsonRollNumber struct {
	Weight   string `json:"weight"`
	ForceVal int    `json:"forceVal"`
}

func (jcfg *jsonRollNumber) build() *RollNumberConfig {
	cfg := &RollNumberConfig{
		Weight:   jcfg.Weight,
		ForceVal: jcfg.ForceVal,
	}

	return cfg
}

func parseRollNumber(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRollNumber:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRollNumber:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRollNumber{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRollNumber:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, mapAwards, err := parseIntValAndAllControllers(ctrls)
		if err != nil {
			goutils.Error("parseRollNumber:parseIntValAndAllControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
		cfgd.MapValAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RollNumberTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
