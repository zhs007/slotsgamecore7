package lowcode

import (
	"fmt"
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

const MergePositionCollectionTypeName = "mergePositionCollection"

type MergePositionCollectionData struct {
	BasicComponentData
	Pos []int
}

// OnNewGame -
func (mergePositionCollectionData *MergePositionCollectionData) OnNewGame(gameProp *GameProperty, component IComponent) {
	mergePositionCollectionData.BasicComponentData.OnNewGame(gameProp, component)
	mergePositionCollectionData.Pos = nil
}

func (mergePositionCollectionData *MergePositionCollectionData) clear() {
	mergePositionCollectionData.Pos = nil
}

// Clone
func (mergePositionCollectionData *MergePositionCollectionData) Clone() IComponentData {
	target := &MergePositionCollectionData{
		BasicComponentData: mergePositionCollectionData.CloneBasicComponentData(),
		Pos:                slices.Clone(mergePositionCollectionData.Pos),
	}

	return target
}

// BuildPBComponentData
func (mergePositionCollectionData *MergePositionCollectionData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MergePositionCollectionData{
		BasicComponentData: mergePositionCollectionData.BuildPBBasicComponentData(),
		Pos:                make([]int32, len(mergePositionCollectionData.Pos)),
	}

	for i, s := range mergePositionCollectionData.Pos {
		pbcd.Pos[i] = int32(s)
	}

	return pbcd
}

// GetPos -
func (mergePositionCollectionData *MergePositionCollectionData) GetPos() []int {
	return mergePositionCollectionData.Pos
}

// HasPos -
func (mergePositionCollectionData *MergePositionCollectionData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(mergePositionCollectionData.Pos, x, y, 0) >= 0
}

// AddPos -
func (mergePositionCollectionData *MergePositionCollectionData) AddPos(x int, y int) {
	mergePositionCollectionData.Pos = append(mergePositionCollectionData.Pos, x, y)
}

// MergePositionCollectionConfig - configuration for MergePositionCollection feature
type MergePositionCollectionConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SrcSymbols           []string `yaml:"srcSymbols" json:"srcSymbols"`               // srcSymbols
	SrcSymbolCodes       []int    `yaml:"-" json:"-"`                                 // srcSymbols
	SrcComponents        []string `yaml:"srcComponents" json:"srcComponents"`         // srcComponents
	OutputToComponent    string   `yaml:"outputToComponent" json:"outputToComponent"` // outputToComponent
	IsClearOutput        bool     `yaml:"isClearOutput" json:"isClearOutput"`         // isClearOutput
	Awards               []*Award `yaml:"awards" json:"awards"`                       // 新的奖励系统
}

// SetLinkComponent
func (cfg *MergePositionCollectionConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

// MergePositionCollection - 也是一个非常特殊的组件，symbol集合
type MergePositionCollection struct {
	*BasicComponent `json:"-"`
	Config          *MergePositionCollectionConfig `json:"config"`
}

// Init -
func (mergePositionCollection *MergePositionCollection) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MergePositionCollection.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &MergePositionCollectionConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MergePositionCollection.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return mergePositionCollection.InitEx(cfg, pool)
}

// InitEx -
func (mergePositionCollection *MergePositionCollection) InitEx(cfg any, pool *GamePropertyPool) error {
	mergePositionCollection.Config = cfg.(*MergePositionCollectionConfig)
	mergePositionCollection.Config.ComponentType = MergePositionCollectionTypeName

	for _, s := range mergePositionCollection.Config.SrcSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("MergePositionCollection.InitEx:SrcSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		mergePositionCollection.Config.SrcSymbolCodes = append(mergePositionCollection.Config.SrcSymbolCodes, sc)
	}

	for _, award := range mergePositionCollection.Config.Awards {
		award.Init()
	}

	mergePositionCollection.onInit(&mergePositionCollection.Config.BasicComponentConfig)

	return nil
}

// procSymbols -
func (mergePositionCollection *MergePositionCollection) procSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult, pos []int, cd *MergePositionCollectionData) ([]int, error) {

	gs := mergePositionCollection.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("MergePositionCollection.procSymbols:GetTargetScene3",
			slog.String("componentName", mergePositionCollection.Name),
			goutils.Err(ErrInvalidScene))

		return nil, ErrInvalidScene
	}

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if slices.Contains(mergePositionCollection.Config.SrcSymbolCodes, s) {
				cd.Pos = append(cd.Pos, x, y)
				pos = append(pos, x, y)
			}
		}
	}

	return pos, nil
}

// procCompnents -
func (mergePositionCollection *MergePositionCollection) procCompnents(gameProp *GameProperty, pos []int,
	cd *MergePositionCollectionData) ([]int, error) {

	for _, v := range mergePositionCollection.Config.SrcComponents {
		scd := gameProp.GetComponentDataWithName(v)
		if scd == nil {
			goutils.Error("MergePositionCollection.procCompnents:GetComponentDataWithName",
				slog.String("componentName", mergePositionCollection.Name),
				slog.String("positionCollection", v),
				goutils.Err(ErrInvalidPositionCollection))

			return nil, ErrInvalidPositionCollection
		}

		pos := scd.GetPos()

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			cd.Pos = append(cd.Pos, x, y)
			pos = append(pos, x, y)
		}
	}

	return pos, nil
}

// playgame
func (mergePositionCollection *MergePositionCollection) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*MergePositionCollectionData)
	cd.clear()

	pos := []int{}

	if len(mergePositionCollection.Config.SrcSymbolCodes) > 0 {
		npos, err := mergePositionCollection.procSymbols(gameProp, curpr, prs, pos, cd)
		if err != nil {
			goutils.Error("MergePositionCollection.OnPlayGame:procSymbols",
				slog.String("componentName", mergePositionCollection.Name),
				goutils.Err(err))

			return "", err
		}

		pos = npos
	}

	if len(mergePositionCollection.Config.SrcComponents) > 0 {
		npos, err := mergePositionCollection.procCompnents(gameProp, pos, cd)
		if err != nil {
			goutils.Error("MergePositionCollection.OnPlayGame:procCompnents",
				slog.String("componentName", mergePositionCollection.Name),
				goutils.Err(err))

			return "", err
		}

		pos = npos
	}

	if len(cd.Pos) == 0 {
		nc := mergePositionCollection.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	pcd := gameProp.GetComponentDataWithName(mergePositionCollection.Config.OutputToComponent)
	pc, isok := gameProp.Components.MapComponents[mergePositionCollection.Config.OutputToComponent]
	if !isok || pcd == nil {
		goutils.Error("MergePositionCollection.OnPlayGame:GetComponentDataWithName",
			slog.String("componentName", mergePositionCollection.Name),
			slog.String("positionCollection", mergePositionCollection.Config.OutputToComponent),
			goutils.Err(ErrInvalidPositionCollection))

		return "", ErrInvalidPositionCollection
	}

	if mergePositionCollection.Config.IsClearOutput {
		pc.ClearData(pcd, true)
	}

	for i := range len(pos) / 2 {
		x := pos[i*2]
		y := pos[i*2+1]

		pc.AddPos(pcd, x, y)
	}

	nc := mergePositionCollection.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (mergePositionCollection *MergePositionCollection) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult,
	mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*MergePositionCollectionData)

	if len(cd.Pos) > 0 {
		fmt.Printf("pos is %v\n", cd.Pos)
	}

	return nil
}

// NewComponentData -
func (mergePositionCollection *MergePositionCollection) NewComponentData() IComponentData {
	return &MergePositionCollectionData{}
}

// GetAllLinkComponents - get all link components
func (mergePositionCollection *MergePositionCollection) GetAllLinkComponents() []string {
	return []string{mergePositionCollection.Config.DefaultNextComponent}
}

func NewMergePositionCollection(name string) IComponent {
	return &MergePositionCollection{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "positionType": "positioncollection",
// "isClearOutput": true,
// "srcComponents": [
//
//	"rs-pos-wl",
//	"rs-pos-wm"
//
// ],
// "outputToComponent": "rs-pos-wm"
type jsonMergePositionCollection struct {
	SrcSymbols        []string `json:"srcSymbols"`        // srcSymbols
	SrcComponents     []string `json:"srcComponents"`     // srcComponents
	OutputToComponent string   `json:"outputToComponent"` // outputToComponent
	IsClearOutput     bool     `json:"isClearOutput"`     // isClearOutput
}

func (jcfg *jsonMergePositionCollection) build() *MergePositionCollectionConfig {
	cfg := &MergePositionCollectionConfig{
		SrcSymbols:        jcfg.SrcSymbols,
		SrcComponents:     jcfg.SrcComponents,
		OutputToComponent: jcfg.OutputToComponent,
		IsClearOutput:     jcfg.IsClearOutput,
	}

	return cfg
}

func parseMergePositionCollection(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parsePositionCollection:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parsePositionCollection:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonMergePositionCollection{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parsePositionCollection:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseLinesTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: MergePositionCollectionTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
