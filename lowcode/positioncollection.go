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

const PositionCollectionTypeName = "positionCollection"

type PositionCollectionType int

const (
	PCTypeNormal        PositionCollectionType = 0
	PCTypeNonRepeatable PositionCollectionType = 1
)

func parsePositionCollectionType(strType string) PositionCollectionType {
	if strType == "nonRepeatable" {
		return PCTypeNonRepeatable
	}

	return PCTypeNormal
}

type PositionCollectionData struct {
	BasicComponentData
	Pos []int
}

// OnNewGame -
func (positionCollectionData *PositionCollectionData) OnNewGame(gameProp *GameProperty, component IComponent) {
	positionCollectionData.BasicComponentData.OnNewGame(gameProp, component)

	positionCollection := component.(*PositionCollection)

	positionCollectionData.Pos = nil

	positionCollectionData.Pos = append(positionCollectionData.Pos, positionCollection.Config.InitPositions...)
}

func (positionCollectionData *PositionCollectionData) clear(pos []int) {
	positionCollectionData.Pos = nil
	positionCollectionData.Pos = append(positionCollectionData.Pos, pos...)
}

// // OnNewStep -
// func (positionCollectionData *PositionCollectionData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	positionCollectionData.BasicComponentData.OnNewStep(gameProp, component)
// }

// Clone
func (positionCollectionData *PositionCollectionData) Clone() IComponentData {
	target := &PositionCollectionData{
		BasicComponentData: positionCollectionData.CloneBasicComponentData(),
	}

	target.Pos = make([]int, len(positionCollectionData.Pos))
	copy(target.Pos, positionCollectionData.Pos)

	return target
}

// BuildPBComponentData
func (positionCollectionData *PositionCollectionData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.PositionCollectionData{
		BasicComponentData: positionCollectionData.BuildPBBasicComponentData(),
	}

	for _, s := range positionCollectionData.Pos {
		pbcd.Pos = append(pbcd.Pos, int32(s))
	}

	return pbcd
}

// GetPos -
func (positionCollectionData *PositionCollectionData) GetPos() []int {
	return positionCollectionData.Pos
}

// HasPos -
func (positionCollectionData *PositionCollectionData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(positionCollectionData.Pos, x, y, 0) >= 0
}

// AddPos -
func (positionCollectionData *PositionCollectionData) AddPos(x int, y int) {
	positionCollectionData.Pos = append(positionCollectionData.Pos, x, y)
}

// PositionCollectionConfig - configuration for PositionCollection feature
type PositionCollectionConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                 `yaml:"type" json:"type"`                         // type
	Type                 PositionCollectionType `yaml:"-" json:"-"`                               // type
	IsNeedClear          bool                   `yaml:"isNeedClear" json:"isNeedClear"`           // isNeedClear
	InitPositions        []int                  `yaml:"initPositions" json:"initPositions"`       // 初始化
	ForeachComponent     string                 `yaml:"foreachComponent" json:"foreachComponent"` // foreach
	Children             []string               `yaml:"-" json:"-"`                               //
}

// SetLinkComponent
func (cfg *PositionCollectionConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "foreach" {
		cfg.ForeachComponent = componentName
	}
}

// PositionCollection - 也是一个非常特殊的组件，symbol集合
type PositionCollection struct {
	*BasicComponent `json:"-"`
	Config          *PositionCollectionConfig `json:"config"`
}

// Init -
func (positionCollection *PositionCollection) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("PositionCollection.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &PositionCollectionConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("PositionCollection.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return positionCollection.InitEx(cfg, pool)
}

// InitEx -
func (positionCollection *PositionCollection) InitEx(cfg any, pool *GamePropertyPool) error {
	positionCollection.Config = cfg.(*PositionCollectionConfig)
	positionCollection.Config.ComponentType = PositionCollectionTypeName

	positionCollection.Config.Type = parsePositionCollectionType(positionCollection.Config.StrType)

	positionCollection.onInit(&positionCollection.Config.BasicComponentConfig)

	return nil
}

// // OnNewGame -
// func (symbolCollection2 *SymbolCollection2) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

// 	cd.OnNewGame()

// 	cd.SymbolCodes = append(cd.SymbolCodes, symbolCollection2.Config.InitSymbolCodes...)

// 	return nil
// }

func (positionCollection *PositionCollection) isClear(basicCD *BasicComponentData) bool {
	clear, isok := basicCD.GetConfigIntVal(CCVClear)
	if isok {
		return clear != 0
	}

	return false
}

// playgame
func (positionCollection *PositionCollection) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*PositionCollectionData)

	if positionCollection.Config.IsNeedClear || positionCollection.isClear(&cd.BasicComponentData) {
		cd.clear(positionCollection.Config.InitPositions)

		cd.SetConfigIntVal(CCVClear, 0)
	}

	nc := positionCollection.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// AddPos -
func (positionCollection *PositionCollection) AddPos(icd IComponentData, x int, y int) {
	if positionCollection.Config.Type == PCTypeNonRepeatable {
		if icd.HasPos(x, y) {
			return
		}
	}

	icd.AddPos(x, y)
}

// ClearData -
func (positionCollection *PositionCollection) ClearData(icd IComponentData, bForceNow bool) {
	if bForceNow {
		cd := icd.(*PositionCollectionData)
		cd.clear(positionCollection.Config.InitPositions)
	}
}

// OnAsciiGame - outpur to asciigame
func (positionCollection *PositionCollection) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult,
	mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*PositionCollectionData)

	if len(cd.Pos) > 0 {
		fmt.Printf("pos is %v\n", cd.Pos)
	}

	return nil
}

// // OnStats
// func (positionCollection *PositionCollection) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (positionCollection *PositionCollection) NewComponentData() IComponentData {
	return &PositionCollectionData{}
}

// OnGameInited - on game inited
func (positionCollection *PositionCollection) OnGameInited(components *ComponentList) error {
	if positionCollection.Config.ForeachComponent != "" {
		positionCollection.Config.Children = components.GetAllLinkComponents(positionCollection.Config.ForeachComponent)
	}

	return nil
}

// IsForeach -
func (positionCollection *PositionCollection) IsForeach() bool {
	return true
}

// GetAllLinkComponents - get all link components
func (positionCollection *PositionCollection) GetAllLinkComponents() []string {
	return []string{positionCollection.Config.DefaultNextComponent, positionCollection.Config.ForeachComponent}
}

// GetChildLinkComponents - get next link components
func (positionCollection *PositionCollection) GetChildLinkComponents() []string {
	return []string{positionCollection.Config.ForeachComponent}
}

func NewPositionCollection(name string) IComponent {
	return &PositionCollection{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "type": "nonRepeatable"
// "initPositions": [
//
//	[
//		1,
//		1
//	]
//
// ]
type jsonPositionCollection struct {
	Type          string  `json:"type"`          // type
	IsNeedClear   bool    `json:"isNeedClear"`   // isNeedClear
	InitPositions [][]int `json:"initPositions"` // initPositions
}

func (jcfg *jsonPositionCollection) build() *PositionCollectionConfig {
	cfg := &PositionCollectionConfig{
		StrType:     jcfg.Type,
		IsNeedClear: jcfg.IsNeedClear,
	}

	for _, arr := range jcfg.InitPositions {
		cfg.InitPositions = append(cfg.InitPositions, arr[0]-1, arr[1]-1)
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parsePositionCollection(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
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

	data := &jsonPositionCollection{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parsePositionCollection:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: PositionCollectionTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
