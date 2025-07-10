package lowcode

import (
	"log/slog"
	"os"
	"slices"
	"strings"

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

const FlowDownSymbolsTypeName = "flowDownSymbols"

type FlowDownSymbolsType int

const (
	FDSTypeFlowToRight FlowDownSymbolsType = 0 // flow to right
)

func parseFlowDownSymbolsType(str string) FlowDownSymbolsType {
	if str == "flowtoright" {
		return FDSTypeFlowToRight
	}

	return FDSTypeFlowToRight
}

type FlowDownSymbolsData struct {
	BasicComponentData
	Pos []int
}

// OnNewGame -
func (flowDownSymbolsData *FlowDownSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	flowDownSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (flowDownSymbolsData *FlowDownSymbolsData) OnNewStep() {
	flowDownSymbolsData.UsedScenes = nil
	flowDownSymbolsData.Pos = nil
}

// Clone
func (flowDownSymbolsData *FlowDownSymbolsData) Clone() IComponentData {
	target := &FlowDownSymbolsData{
		BasicComponentData: flowDownSymbolsData.CloneBasicComponentData(),
		Pos:                slices.Clone(flowDownSymbolsData.Pos),
	}

	return target
}

// BuildPBComponentData
func (flowDownSymbolsData *FlowDownSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.FlowDownSymbolsData{
		BasicComponentData: flowDownSymbolsData.BuildPBBasicComponentData(),
		Pos:                make([]int32, len(flowDownSymbolsData.Pos)),
	}

	for i, v := range flowDownSymbolsData.Pos {
		pbcd.Pos[i] = int32(v)
	}

	return pbcd
}

// GetPos -
func (flowDownSymbolsData *FlowDownSymbolsData) GetPos() []int {
	return flowDownSymbolsData.Pos
}

// HasPos -
func (flowDownSymbolsData *FlowDownSymbolsData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(flowDownSymbolsData.Pos, x, y, 0) >= 0
}

// AddPos -
func (flowDownSymbolsData *FlowDownSymbolsData) AddPos(x int, y int) {
	flowDownSymbolsData.Pos = append(flowDownSymbolsData.Pos, x, y)
}

// AddPosEx -
func (flowDownSymbolsData *FlowDownSymbolsData) AddPosEx(x int, y int) {
	if !flowDownSymbolsData.HasPos(x, y) {
		flowDownSymbolsData.AddPos(x, y)
	}
}

// FlowDownSymbolsConfig - configuration for FlowDownSymbols
type FlowDownSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string              `yaml:"type" json:"type"`
	Type                 FlowDownSymbolsType `yaml:"-" json:"-"`
	SrcMask              string              `yaml:"srcMask" json:"srcMask"`
	FillSymbol           string              `yaml:"fillSymbol" json:"fillSymbol"`
	FillSymbolCode       int                 `yaml:"-" json:"-"`
	BlockSymbols         []string            `yaml:"blockSymbols" json:"blockSymbols"`
	BlockSymbolCodes     []int               `yaml:"-" json:"-"`
	Number               int                 `yaml:"number" json:"number"`
	JumpToComponent      string              `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
	Controllers          []*Award            `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *FlowDownSymbolsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type FlowDownSymbols struct {
	*BasicComponent `json:"-"`
	Config          *FlowDownSymbolsConfig `json:"config"`
}

// Init -
func (flowDownSymbols *FlowDownSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("FlowDownSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &FlowDownSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("FlowDownSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return flowDownSymbols.InitEx(cfg, pool)
}

// InitEx -
func (flowDownSymbols *FlowDownSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	flowDownSymbols.Config = cfg.(*FlowDownSymbolsConfig)
	flowDownSymbols.Config.ComponentType = FlowDownSymbolsTypeName

	flowDownSymbols.Config.Type = parseFlowDownSymbolsType(flowDownSymbols.Config.StrType)

	for _, v := range flowDownSymbols.Config.BlockSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("FlowDownSymbols.InitEx:BlockSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		flowDownSymbols.Config.BlockSymbolCodes = append(flowDownSymbols.Config.BlockSymbolCodes, sc)
	}

	if len(flowDownSymbols.Config.FillSymbol) > 0 {
		sc, isok := pool.DefaultPaytables.MapSymbols[flowDownSymbols.Config.FillSymbol]
		if !isok {
			goutils.Error("FlowDownSymbols.InitEx:FillSymbol",
				slog.String("symbol", flowDownSymbols.Config.FillSymbol),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		flowDownSymbols.Config.FillSymbolCode = sc
	}

	for _, award := range flowDownSymbols.Config.Controllers {
		award.Init()
	}

	flowDownSymbols.onInit(&flowDownSymbols.Config.BasicComponentConfig)

	return nil
}

// procFlowToRightOne -
func (flowDownSymbols *FlowDownSymbols) procFlowToRightOne(x int, gameProp *GameProperty, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult, cd *FlowDownSymbolsData, gs *sgc7game.GameScene) error {

	for y := 0; y < gs.Height; y++ {
		if goutils.IndexOfIntSlice(flowDownSymbols.Config.BlockSymbolCodes, gs.Arr[x][y], 0) >= 0 {
			if y == 0 {
				if x < gs.Width-1 {
					return flowDownSymbols.procFlowToRightOne(x+1, gameProp, curpr, prs, cd, gs)
				}
			} else {
				gs.Arr[x][y-1] = flowDownSymbols.Config.FillSymbolCode
			}

			return nil
		}

		if y == gs.Height-1 {
			gs.Arr[x][y] = flowDownSymbols.Config.FillSymbolCode

			return nil
		}
	}

	return nil
}

func (flowDownSymbols *FlowDownSymbols) getNumber(cd *FlowDownSymbolsData) int {
	number, isok := cd.GetConfigIntVal(CCVNumber)
	if isok {
		return number
	}

	return flowDownSymbols.Config.Number
}

// procFlowToRight -
func (flowDownSymbols *FlowDownSymbols) procFlowToRight(x int, gameProp *GameProperty, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult, cd *FlowDownSymbolsData, gs *sgc7game.GameScene) (*sgc7game.GameScene, error) {

	if flowDownSymbols.Config.Number <= 0 {
		return gs, nil
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	number := flowDownSymbols.getNumber(cd)

	for n := range number {
		err := flowDownSymbols.procFlowToRightOne(x, gameProp, curpr, prs, cd, ngs)
		if err != nil {
			goutils.Error("FlowDownSymbols.procFlowToRight:procFlowToRightOne",
				slog.Int("x", x),
				slog.Int("n", n+1),
				goutils.Err(err))

			return nil, err
		}
	}

	return ngs, nil
}

// OnProcControllers -
func (flowDownSymbols *FlowDownSymbols) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(flowDownSymbols.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, flowDownSymbols.Config.Controllers, curpr, gp)
	}
}

// playgame
func (flowDownSymbols *FlowDownSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*FlowDownSymbolsData)

	cd.OnNewStep()

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	sc2 := gs

	if flowDownSymbols.Config.Type == FDSTypeFlowToRight {
		if flowDownSymbols.Config.SrcMask == "" {
			for x := range gs.Width {
				newgs, err := flowDownSymbols.procFlowToRight(x, gameProp, curpr, prs, cd, sc2)
				if err != nil {
					goutils.Error("FlowDownSymbols.OnPlayGame:procFlowToRight",
						goutils.Err(err))

					return "", err
				}

				if newgs != gs {
					flowDownSymbols.AddScene(gameProp, curpr, newgs, &cd.BasicComponentData)
				}

				sc2 = newgs
			}
		} else {
			mask, err := gameProp.GetMask(flowDownSymbols.Config.SrcMask)
			if err != nil {
				goutils.Error("FlowDownSymbols.OnPlayGame:GetMask",
					slog.String("mask", flowDownSymbols.Config.SrcMask),
					goutils.Err(err))

				return "", err
			}

			for x, v := range mask {
				if v {
					newgs, err := flowDownSymbols.procFlowToRight(x, gameProp, curpr, prs, cd, sc2)
					if err != nil {
						goutils.Error("FlowDownSymbols.OnPlayGame:procFlowToRight",
							goutils.Err(err))

						return "", err
					}

					if newgs != gs {
						flowDownSymbols.AddScene(gameProp, curpr, newgs, &cd.BasicComponentData)
					}

					sc2 = newgs
				}
			}
		}

		if sc2 == gs {
			nc := flowDownSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		flowDownSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

		nc := flowDownSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	goutils.Error("FlowDownSymbols.OnPlayGame:InvalidType",
		slog.String("type", flowDownSymbols.Config.StrType),
		goutils.Err(ErrIvalidComponentConfig))

	return "", ErrIvalidComponentConfig
}

// NewComponentData -
func (flowDownSymbols *FlowDownSymbols) NewComponentData() IComponentData {
	return &FlowDownSymbolsData{}
}

// OnAsciiGame - outpur to asciigame
func (flowDownSymbols *FlowDownSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd := icd.(*FlowDownSymbolsData)

	asciigame.OutputScene("after flowDownSymbols", pr.Scenes[msd.UsedScenes[0]], mapSymbolColor)

	return nil
}

func NewFlowDownSymbols(name string) IComponent {
	return &FlowDownSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "flowToRight",
// "srcMask": "bg-mask",
// "fillSymbol": "WL",
// "blockSymbols": [
// 	"WL"
// ],
// "number": 1

type jsonFlowDownSymbols struct {
	StrType      string   `json:"type"`
	SrcMask      string   `json:"srcMask"`
	FillSymbol   string   `json:"fillSymbol"`
	BlockSymbols []string `json:"blockSymbols"`
	Number       int      `json:"number"`
}

func (jcfg *jsonFlowDownSymbols) build() *FlowDownSymbolsConfig {
	cfg := &FlowDownSymbolsConfig{
		StrType:      strings.ToLower(jcfg.StrType),
		SrcMask:      jcfg.SrcMask,
		BlockSymbols: slices.Clone(jcfg.BlockSymbols),
		FillSymbol:   jcfg.FillSymbol,
		Number:       jcfg.Number,
	}

	return cfg
}

func parseFlowDownSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseFlowDownSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseFlowDownSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonFlowDownSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseFlowDownSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseFlowDownSymbols:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: FlowDownSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
