package lowcode

import (
	"log/slog"
	"os"
	"slices"

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

const ReelsCollectorTypeName = "reelsCollector"

type ReelsCollectorTriggerType int

const (
	RCTTypeLeft  ReelsCollectorTriggerType = 0 // left
	RCTTypeRight ReelsCollectorTriggerType = 1 // right
)

func parseReelsCollectorTriggerType(str string) ReelsCollectorTriggerType {
	if str == "right" {
		return RCTTypeRight
	}

	return RCTTypeLeft
}

type ReelsCollectorPS struct {
	Collectors []int `json:"collectors"` // collectors
}

// SetPublicJson
func (ps *ReelsCollectorPS) SetPublicJson(str string) error {
	err := sonic.UnmarshalString(str, ps)
	if err != nil {
		goutils.Error("ReelsCollectorPS.SetPublicJson:UnmarshalString",
			goutils.Err(err))

		return err
	}

	return nil
}

// SetPrivateJson
func (ps *ReelsCollectorPS) SetPrivateJson(str string) error {
	return nil
}

// GetPublicJson
func (ps *ReelsCollectorPS) GetPublicJson() string {
	str, err := sonic.MarshalString(ps)
	if err != nil {
		goutils.Error("ReelsCollectorPS.GetPublicJson:MarshalString",
			goutils.Err(err))

		return ""
	}

	return str
}

// GetPrivateJson
func (ps *ReelsCollectorPS) GetPrivateJson() string {
	return ""
}

// Clone
func (ps *ReelsCollectorPS) Clone() IComponentPS {
	return &ReelsCollectorPS{
		Collectors: slices.Clone(ps.Collectors),
	}
}

type ReelsCollectorData struct {
	BasicComponentData
	Collectors []int
}

// OnNewGame -
func (reelsCollectorData *ReelsCollectorData) OnNewGame(gameProp *GameProperty, component IComponent) {
	reelsCollectorData.BasicComponentData.OnNewGame(gameProp, component)

	reelsCollectorData.Collectors = nil
}

// BuildPBComponentData
func (reelsCollectorData *ReelsCollectorData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.ReelsCollectorData{
		BasicComponentData: reelsCollectorData.BuildPBBasicComponentData(),
		Collectors:         make([]int32, len(reelsCollectorData.Collectors)),
	}

	for i, v := range reelsCollectorData.Collectors {
		pbcd.Collectors[i] = int32(v)
	}

	return pbcd
}

// Clone
func (reelsCollectorData *ReelsCollectorData) Clone() IComponentData {
	target := &ReelsCollectorData{
		BasicComponentData: reelsCollectorData.CloneBasicComponentData(),
		Collectors:         slices.Clone(reelsCollectorData.Collectors),
	}

	return target
}

// ReelsCollectorConfig - configuration for ReelsCollector
type ReelsCollectorConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrTriggerType       string                    `yaml:"triggerType" json:"triggerType"`     // triggerType
	TriggerType          ReelsCollectorTriggerType `yaml:"-" json:"-"`                         // triggerType
	MaxVal               int                       `yaml:"maxVal" json:"maxVal"`               // maxVal
	IsPlayerState        bool                      `yaml:"IsPlayerState" json:"IsPlayerState"` // IsPlayerState
	OutputMask           string                    `yaml:"outputMask" json:"outputMask"`       // outputMask
	MapAwards            map[string][]*Award       `yaml:"mapAwards" json:"mapAwards"`         // 新的奖励系统
}

// SetLinkComponent
func (cfg *ReelsCollectorConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type ReelsCollector struct {
	*BasicComponent `json:"-"`
	Config          *ReelsCollectorConfig `json:"config"`
}

// Init -
func (reelsCollector *ReelsCollector) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReelsCollector.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ReelsCollectorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReelsCollector.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return reelsCollector.InitEx(cfg, pool)
}

// InitEx -
func (reelsCollector *ReelsCollector) InitEx(cfg any, pool *GamePropertyPool) error {
	reelsCollector.Config = cfg.(*ReelsCollectorConfig)
	reelsCollector.Config.ComponentType = ReelsCollectorTypeName

	reelsCollector.Config.TriggerType = parseReelsCollectorTriggerType(reelsCollector.Config.StrTriggerType)

	for _, awards := range reelsCollector.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	reelsCollector.onInit(&reelsCollector.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (reelsCollector *ReelsCollector) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(reelsCollector.Config.MapAwards) > 0 {
		awards, isok := reelsCollector.Config.MapAwards[strVal]
		if isok {
			gameProp.procAwards(plugin, awards, curpr, gp)
		}
	}
}

// procMask
func (reelsCollector *ReelsCollector) procMask(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, triggerReelIndex int) error {

	if reelsCollector.Config.OutputMask != "" {
		gameProp.UseComponent(reelsCollector.Config.OutputMask)

		mask := make([]bool, gameProp.GetVal(GamePropWidth))
		mask[triggerReelIndex] = true

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, reelsCollector.Config.OutputMask, mask, true)
	}

	return nil
}

// func (reelsCollector *ReelsCollector) isClear(basicCD *BasicComponentData) bool {
// 	clear, isok := basicCD.GetConfigIntVal(CCVClear)
// 	if isok {
// 		return clear != 0
// 	}

// 	return false
// }

// playgame
func (reelsCollector *ReelsCollector) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ips sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ReelsCollectorData)

	if reelsCollector.Config.IsPlayerState {
		ps, isok := ips.(*PlayerState)
		if !isok {
			goutils.Error("ReelsCollector.OnPlayGame:PlayerState",
				goutils.Err(ErrIvalidPlayerState))

			return "", ErrIvalidPlayerState
		}

		betMethod := stake.CashBet / stake.CoinBet
		bmd := ps.GetBetMethodPub(int(betMethod))
		if bmd == nil {
			goutils.Error("ReelsCollector.OnPlayGame:GetBetMethodPub",
				goutils.Err(ErrIvalidPlayerState))

			return "", ErrIvalidPlayerState
		}

		cps := bmd.GetBetCPS(int(stake.CoinBet), reelsCollector.GetName())
		if cps == nil {
			goutils.Error("ReelsCollector.OnPlayGame:GetBetCPS",
				goutils.Err(ErrIvalidPlayerState))

			return "", ErrIvalidPlayerState
		}

		cbps, isok := cps.(*ReelsCollectorPS)
		if !isok {
			goutils.Error("ReelsCollector.OnPlayGame:ReelsCollectorPS",
				goutils.Err(ErrIvalidPlayerState))

			return "", ErrIvalidPlayerState
		}

		if len(cbps.Collectors) == 0 {
			w := gameProp.GetVal(GamePropWidth)
			cbps.Collectors = make([]int, w)
		}

		if reelsCollector.Config.TriggerType == RCTTypeLeft {
			for i, v := range cbps.Collectors {
				if v == reelsCollector.Config.MaxVal {
					cbps.Collectors[i] = 0
					cd.Collectors = slices.Clone(cbps.Collectors)

					reelsCollector.procMask(gameProp, curpr, gp, plugin, i)
					reelsCollector.ProcControllers(gameProp, plugin, curpr, gp, i, "")

					nc := reelsCollector.onStepEnd(gameProp, curpr, gp, "")

					return nc, nil
				}
			}
		} else if reelsCollector.Config.TriggerType == RCTTypeRight {
			for i := len(cbps.Collectors) - 1; i >= 0; i-- {
				if cbps.Collectors[i] == reelsCollector.Config.MaxVal {
					cbps.Collectors[i] = 0
					cd.Collectors = slices.Clone(cbps.Collectors)

					reelsCollector.procMask(gameProp, curpr, gp, plugin, i)
					reelsCollector.ProcControllers(gameProp, plugin, curpr, gp, i, "")

					nc := reelsCollector.onStepEnd(gameProp, curpr, gp, "")

					return nc, nil
				}
			}
		} else {
			goutils.Error("ReelsCollector.OnPlayGame:InvalidTriggerType",
				slog.String("triggerType", reelsCollector.Config.StrTriggerType),
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

		cd.Collectors = slices.Clone(cbps.Collectors)
		nc := reelsCollector.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	if len(cd.Collectors) == 0 {
		w := gameProp.GetVal(GamePropWidth)
		cd.Collectors = make([]int, w)
	}

	if reelsCollector.Config.TriggerType == RCTTypeLeft {
		for i, v := range cd.Collectors {
			if v == reelsCollector.Config.MaxVal {
				cd.Collectors[i] = 0
				cd.Collectors = slices.Clone(cd.Collectors)

				reelsCollector.procMask(gameProp, curpr, gp, plugin, i)
				reelsCollector.ProcControllers(gameProp, plugin, curpr, gp, i, "")

				nc := reelsCollector.onStepEnd(gameProp, curpr, gp, "")

				return nc, nil
			}
		}
	} else if reelsCollector.Config.TriggerType == RCTTypeRight {
		for i := len(cd.Collectors) - 1; i >= 0; i-- {
			if cd.Collectors[i] == reelsCollector.Config.MaxVal {
				cd.Collectors[i] = 0
				cd.Collectors = slices.Clone(cd.Collectors)

				reelsCollector.procMask(gameProp, curpr, gp, plugin, i)
				reelsCollector.ProcControllers(gameProp, plugin, curpr, gp, i, "")

				nc := reelsCollector.onStepEnd(gameProp, curpr, gp, "")

				return nc, nil
			}
		}
	} else {
		goutils.Error("ReelsCollector.OnPlayGame:InvalidTriggerType",
			slog.String("triggerType", reelsCollector.Config.StrTriggerType),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	cd.Collectors = slices.Clone(cd.Collectors)
	nc := reelsCollector.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (featureBar2 *ReelsCollector) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult,
	mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	return nil
}

// NewComponentData -
func (featureBar2 *ReelsCollector) NewComponentData() IComponentData {
	return &ReelsCollectorData{}
}

// InitPlayerState -
func (featureBar2 *ReelsCollector) InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	ps *PlayerState, betMethod int, bet int) error {

	if featureBar2.Config.IsPlayerState {
		bmd := ps.GetBetMethodPub(betMethod)
		if bet <= 0 {
			return nil
		}

		bps := bmd.GetBetPS(bet)

		cname := featureBar2.GetName()

		_, isok := bps.MapComponentData[cname]
		if !isok {
			str, isok := bps.MapString[cname]
			if isok {
				cps := &ReelsCollectorPS{}
				cps.SetPublicJson(str)

				bps.MapComponentData[cname] = cps
			} else {
				w := gameProp.GetVal(GamePropWidth)
				cps := &ReelsCollectorPS{
					Collectors: make([]int, w),
				}

				bps.MapComponentData[featureBar2.GetName()] = cps
			}
		}
	}

	return nil
}

func NewReelsCollector(name string) IComponent {
	return &ReelsCollector{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "triggerType": "left",
// "IsPlayerState": false,
// "maxVal": 4
// "outputMask": "fg-mask-h3"
type jsonReelsCollector struct {
	StrTriggerType string `json:"triggerType"`   // triggerType
	MaxVal         int    `json:"maxVal"`        // maxVal
	IsPlayerState  bool   `json:"IsPlayerState"` // IsPlayerState
	OutputMask     string `json:"outputMask"`    // outputMask
}

func (jcfg *jsonReelsCollector) build() *ReelsCollectorConfig {
	cfg := &ReelsCollectorConfig{
		StrTriggerType: jcfg.StrTriggerType,
		MaxVal:         jcfg.MaxVal,
		IsPlayerState:  jcfg.IsPlayerState,
		OutputMask:     jcfg.OutputMask,
	}

	return cfg
}

func parseReelsCollector(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseReelsCollector:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseReelsCollector:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonReelsCollector{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseReelsCollector:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseMapStringAndAllControllers(ctrls)
		if err != nil {
			goutils.Error("parseReelsCollector:parseMapStringAndAllControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ReelsCollectorTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
