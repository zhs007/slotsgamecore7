package lowcode

import (
	"fmt"
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
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const ReelTriggerTypeName = "reelTrigger"

type ReelTriggerType int

const (
	RTTypeRow          ReelTriggerType = 0 // row
	RTTypeColumn       ReelTriggerType = 1 // column
	RTTypeRowNumber    ReelTriggerType = 2 // row number
	RTTypeColumnNumber ReelTriggerType = 3 // column number
	RTTypeFullScreen   ReelTriggerType = 4 // full screen
)

func parseReelTriggerType(str string) ReelTriggerType {
	switch str {
	case "row":
		return RTTypeRow
	case "column":
		return RTTypeColumn
	case "rownumber":
		return RTTypeRowNumber
	case "columnnumber":
		return RTTypeColumnNumber
	case "fullscreen":
		return RTTypeFullScreen
	}

	return RTTypeRow
}

type ReelTriggerData struct {
	BasicComponentData
	NextComponent string
	Masks         []bool
	Number        int
	SymbolCode    int
}

// OnNewGame -
func (reelTriggerData *ReelTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	reelTriggerData.BasicComponentData.OnNewGame(gameProp, component)

	reelTrigger := component.(*ReelTrigger)

	switch reelTrigger.Config.Type {
	case RTTypeRow, RTTypeRowNumber:
		reelTriggerData.Masks = make([]bool, gameProp.GetVal(GamePropHeight))
	case RTTypeColumn, RTTypeColumnNumber:
		reelTriggerData.Masks = make([]bool, gameProp.GetVal(GamePropWidth))
	}

	reelTriggerData.Number = 0
}

// onNewStep -
func (reelTriggerData *ReelTriggerData) onNewStep() {
	reelTriggerData.UsedResults = nil
	reelTriggerData.NextComponent = ""
	reelTriggerData.SymbolCode = -1
}

// Clone
func (reelTriggerData *ReelTriggerData) Clone() IComponentData {
	target := &ReelTriggerData{
		BasicComponentData: reelTriggerData.CloneBasicComponentData(),
		NextComponent:      reelTriggerData.NextComponent,
		Masks:              make([]bool, len(reelTriggerData.Masks)),
		Number:             reelTriggerData.Number,
	}

	copy(target.Masks, reelTriggerData.Masks)

	return target
}

// BuildPBComponentData
func (reelTriggerData *ReelTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.ReelTriggerData{
		BasicComponentData: reelTriggerData.BuildPBBasicComponentData(),
		NextComponent:      reelTriggerData.NextComponent,
		Masks:              make([]bool, len(reelTriggerData.Masks)),
		Number:             int32(reelTriggerData.Number),
	}

	copy(pbcd.Masks, reelTriggerData.Masks)

	return pbcd
}

// GetValEx -
func (reelTriggerData *ReelTriggerData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVNumber || key == CVOutputInt {
		return reelTriggerData.Number, true
	}

	return 0, false
}

// ReelTriggerConfig - configuration for ReelTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type ReelTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string            `yaml:"symbols" json:"symbols"`                       // symbols
	SymbolCodes          []int               `yaml:"-" json:"-"`                                   // symbol codes
	StrType              string              `yaml:"type" json:"type"`                             // ReelTriggerType
	Type                 ReelTriggerType     `yaml:"-" json:"-"`                                   // ReelTriggerType
	WildSymbols          []string            `yaml:"wildSymbols" json:"wildSymbols"`               // wild etc
	WildSymbolCodes      []int               `yaml:"-" json:"-"`                                   // wild symbolCode
	MinSymbolNum         int                 `yaml:"minSymbolNum" json:"minSymbolNum"`             // minSymbolNum
	TargetMask           string              `yaml:"targetMask" json:"targetMask"`                 // 可以把结果传递给一个mask
	MapBranchs           map[int]*BranchNode `yaml:"mapBranchs" json:"mapBranchs"`                 // mapBranchs
	IsCheckEmptySymbol   bool                `yaml:"isCheckEmptySymbol" json:"isCheckEmptySymbol"` //
}

// SetLinkComponent
func (cfg *ReelTriggerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else {
		if cfg.MapBranchs == nil {
			cfg.MapBranchs = make(map[int]*BranchNode)
		}

		if strings.ToLower(link) == "fullscreen" {
			cfg.MapBranchs[-1] = &BranchNode{
				JumpToComponent: componentName,
			}

			return
		}

		i64, err := goutils.String2Int64(link)
		if err != nil {
			goutils.Error("ReelTrigger.SetLinkComponent",
				slog.String("link", link),
				goutils.Err(err),
			)

			return
		}

		if cfg.MapBranchs[int(i64)] == nil {
			cfg.MapBranchs[int(i64)] = &BranchNode{
				JumpToComponent: componentName,
			}
		} else {
			cfg.MapBranchs[int(i64)].JumpToComponent = componentName
		}
	}
}

type ReelTrigger struct {
	*BasicComponent `json:"-"`
	Config          *ReelTriggerConfig `json:"config"`
}

// Init -
func (reelTrigger *ReelTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReelTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ReelTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReelTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return reelTrigger.InitEx(cfg, pool)
}

// InitEx -
func (reelTrigger *ReelTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	reelTrigger.Config = cfg.(*ReelTriggerConfig)
	reelTrigger.Config.ComponentType = ReelTriggerTypeName

	reelTrigger.Config.Type = parseReelTriggerType(reelTrigger.Config.StrType)

	for _, v := range reelTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("ReelTrigger.InitEx:Symbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		reelTrigger.Config.SymbolCodes = append(reelTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range reelTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ReelTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		reelTrigger.Config.WildSymbolCodes = append(reelTrigger.Config.WildSymbolCodes, sc)
	}

	for _, branch := range reelTrigger.Config.MapBranchs {
		for _, award := range branch.Awards {
			award.Init()
		}
	}

	switch reelTrigger.Config.Type {
	case RTTypeRow, RTTypeRowNumber:
		if reelTrigger.Config.MinSymbolNum <= 0 || reelTrigger.Config.MinSymbolNum > pool.Config.Width {
			reelTrigger.Config.MinSymbolNum = pool.Config.Width
		}
	case RTTypeColumn, RTTypeColumnNumber:
		if reelTrigger.Config.MinSymbolNum <= 0 || reelTrigger.Config.MinSymbolNum > pool.Config.Height {
			reelTrigger.Config.MinSymbolNum = pool.Config.Height
		}
	case RTTypeFullScreen:
		if reelTrigger.Config.MinSymbolNum <= 0 || reelTrigger.Config.MinSymbolNum > pool.Config.Height*pool.Config.Width {
			reelTrigger.Config.MinSymbolNum = pool.Config.Height * pool.Config.Width
		}
	}

	reelTrigger.onInit(&reelTrigger.Config.BasicComponentConfig)

	return nil
}

func (reelTrigger *ReelTrigger) isValidSymbolCode(rtdata *ReelTriggerData, sc int) bool {
	if rtdata.SymbolCode == -1 {
		if slices.Contains(reelTrigger.Config.WildSymbolCodes, sc) {
			return true
		}

		si := slices.Index(reelTrigger.Config.SymbolCodes, sc)
		if si < 0 {
			return false
		}

		rtdata.SymbolCode = reelTrigger.Config.SymbolCodes[si]

		return true
	} else {
		if sc == rtdata.SymbolCode {
			return true
		}

		if slices.Contains(reelTrigger.Config.WildSymbolCodes, sc) {
			return true
		}
	}

	return false
}

func (reelTrigger *ReelTrigger) calcRow(rtdata *ReelTriggerData, gs *sgc7game.GameScene) ([]bool, int) {
	triggerArr := make([]bool, len(rtdata.Masks))
	triggerNum := 0

	for y := 0; y < gs.Height; y++ {
		num := 0
		for x := 0; x < gs.Width; x++ {
			if reelTrigger.Config.IsCheckEmptySymbol {
				if gs.Arr[x][y] < 0 {
					num++
				}
			} else if reelTrigger.isValidSymbolCode(rtdata, gs.Arr[x][y]) {
				num++
			}
		}

		if num >= reelTrigger.Config.MinSymbolNum {
			triggerArr[y] = true
			triggerNum++
		}
	}

	return triggerArr, triggerNum
}

func (reelTrigger *ReelTrigger) procFullScreen(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	rtdata *ReelTriggerData, gs *sgc7game.GameScene) bool {

	triggerNum := 0

	for y := 0; y < gs.Height; y++ {
		for x := 0; x < gs.Width; x++ {
			if reelTrigger.Config.IsCheckEmptySymbol {
				if gs.Arr[x][y] < 0 {
					triggerNum++
				}
			} else if reelTrigger.isValidSymbolCode(rtdata, gs.Arr[x][y]) {
				triggerNum++
			}
		}
	}

	if triggerNum >= reelTrigger.Config.MinSymbolNum {
		n, isok := reelTrigger.Config.MapBranchs[-1]
		if isok {
			gameProp.procAwards(plugin, n.Awards, curpr, gp)

			rtdata.NextComponent = n.JumpToComponent
		}

		return true
	}

	return false
}

func (reelTrigger *ReelTrigger) procRow(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	rtdata *ReelTriggerData, gs *sgc7game.GameScene) bool {

	triggerArr, _ := reelTrigger.calcRow(rtdata, gs)

	rtdata.NextComponent = ""
	isTrigger := false

	for i, v := range triggerArr {
		if !rtdata.Masks[i] && v {
			rtdata.Masks[i] = true

			if reelTrigger.Config.MapBranchs[i+1] != nil {
				gameProp.procAwards(plugin, reelTrigger.Config.MapBranchs[i+1].Awards, curpr, gp)

				rtdata.NextComponent = reelTrigger.Config.MapBranchs[i+1].JumpToComponent
			}

			rtdata.Number++

			isTrigger = true
		}
	}

	return isTrigger
}

func (reelTrigger *ReelTrigger) procRowNumber(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	rtdata *ReelTriggerData, gs *sgc7game.GameScene) bool {

	triggerArr, triggerNum := reelTrigger.calcRow(rtdata, gs)

	rtdata.NextComponent = ""

	rtdata.Masks = triggerArr

	if triggerNum > rtdata.Number {
		for i := rtdata.Number; i < triggerNum; i++ {
			if reelTrigger.Config.MapBranchs[i+1] != nil {
				gameProp.procAwards(plugin, reelTrigger.Config.MapBranchs[i+1].Awards, curpr, gp)

				rtdata.NextComponent = reelTrigger.Config.MapBranchs[i+1].JumpToComponent
			}
		}

		rtdata.Number = triggerNum

		return true
	}

	return false
}

func (reelTrigger *ReelTrigger) calcColumn(rtdata *ReelTriggerData, gs *sgc7game.GameScene) ([]bool, int) {
	triggerArr := make([]bool, len(rtdata.Masks))
	triggerNum := 0

	for x := 0; x < gs.Width; x++ {
		num := 0
		for y := 0; y < gs.Height; y++ {
			if reelTrigger.Config.IsCheckEmptySymbol {
				if gs.Arr[x][y] < 0 {
					num++
				}
			} else if reelTrigger.isValidSymbolCode(rtdata, gs.Arr[x][y]) {
				num++
			}
		}

		if num >= reelTrigger.Config.MinSymbolNum {
			triggerArr[x] = true
			triggerNum++
		}
	}

	return triggerArr, triggerNum
}

func (reelTrigger *ReelTrigger) procColumn(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	rtdata *ReelTriggerData, gs *sgc7game.GameScene) bool {

	triggerArr, _ := reelTrigger.calcColumn(rtdata, gs)

	rtdata.NextComponent = ""
	isTrigger := false

	for i, v := range triggerArr {
		if !rtdata.Masks[i] && v {
			rtdata.Masks[i] = true

			if reelTrigger.Config.MapBranchs[i+1] != nil {
				gameProp.procAwards(plugin, reelTrigger.Config.MapBranchs[i+1].Awards, curpr, gp)

				rtdata.NextComponent = reelTrigger.Config.MapBranchs[i+1].JumpToComponent
			}

			rtdata.Number++

			isTrigger = true
		}
	}

	return isTrigger
}

func (reelTrigger *ReelTrigger) procColumnNumber(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	rtdata *ReelTriggerData, gs *sgc7game.GameScene) bool {

	triggerArr, triggerNum := reelTrigger.calcColumn(rtdata, gs)

	rtdata.NextComponent = ""

	rtdata.Masks = triggerArr

	if triggerNum > rtdata.Number {
		for i := rtdata.Number; i < triggerNum; i++ {
			if reelTrigger.Config.MapBranchs[i+1] != nil {
				gameProp.procAwards(plugin, reelTrigger.Config.MapBranchs[i+1].Awards, curpr, gp)

				rtdata.NextComponent = reelTrigger.Config.MapBranchs[i+1].JumpToComponent
			}
		}

		rtdata.Number = triggerNum

		return true
	}

	return false
}

// playgame
func (reelTrigger *ReelTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	rtd := cd.(*ReelTriggerData)
	rtd.onNewStep()

	gs := reelTrigger.GetTargetScene3(gameProp, curpr, prs, 0)
	isTrigger := false

	switch reelTrigger.Config.Type {
	case RTTypeRow:
		isTrigger = reelTrigger.procRow(gameProp, curpr, gp, plugin, rtd, gs)
	case RTTypeRowNumber:
		isTrigger = reelTrigger.procRowNumber(gameProp, curpr, gp, plugin, rtd, gs)
	case RTTypeColumn:
		isTrigger = reelTrigger.procColumn(gameProp, curpr, gp, plugin, rtd, gs)
	case RTTypeColumnNumber:
		isTrigger = reelTrigger.procColumnNumber(gameProp, curpr, gp, plugin, rtd, gs)
	case RTTypeFullScreen:
		isTrigger = reelTrigger.procFullScreen(gameProp, curpr, gp, plugin, rtd, gs)
	}

	if isTrigger {
		nc := reelTrigger.onStepEnd(gameProp, curpr, gp, rtd.NextComponent)

		return nc, nil
	}

	nc := reelTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (reelTrigger *ReelTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	rtd := cd.(*ReelTriggerData)

	if rtd.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", reelTrigger.Name, rtd.NextComponent)
	}

	return nil
}

// NewComponentData -
func (reelTrigger *ReelTrigger) NewComponentData() IComponentData {
	return &ReelTriggerData{}
}

// NewStats2 -
func (reelTrigger *ReelTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (reelTrigger *ReelTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	reelTrigger.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)
}

// GetAllLinkComponents - get all link components
func (reelTrigger *ReelTrigger) GetAllLinkComponents() []string {
	lst := []string{reelTrigger.Config.DefaultNextComponent}

	for _, v := range reelTrigger.Config.MapBranchs {
		lst = append(lst, v.JumpToComponent)
	}

	return lst
}

// GetNextLinkComponents - get next link components
func (reelTrigger *ReelTrigger) GetNextLinkComponents() []string {
	lst := []string{reelTrigger.Config.DefaultNextComponent}

	for _, v := range reelTrigger.Config.MapBranchs {
		lst = append(lst, v.JumpToComponent)
	}

	return lst
}

// func (reelTrigger *ReelTrigger) getSymbols(gameProp *GameProperty) []int {
// 	s := gameProp.GetCurCallStackSymbol()
// 	if s >= 0 {
// 		return []int{s}
// 	}

// 	return []int{reelTrigger.Config.SymbolCode}
// }

func NewReelTrigger(name string) IComponent {
	return &ReelTrigger{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "triggerType": "rowNumber",
// "minSymbolNum": 4,
// "symbols": "CASH"
// "IsCheckEmptySymbol": true
type jsonReelTrigger struct {
	Symbols            []string `json:"symbols"`
	TriggerType        string   `json:"triggerType"`
	MinSymbolNum       int      `json:"minSymbolNum"`
	WildSymbols        []string `json:"wildSymbols"`
	TargetMask         string   `json:"targetMask"` // 可以把结果传递给一个mask
	IsCheckEmptySymbol bool     `json:"IsCheckEmptySymbol"`
	TriggerSymbols     []string `json:"triggerSymbols"`
}

func (jcfg *jsonReelTrigger) build() *ReelTriggerConfig {
	cfg := &ReelTriggerConfig{
		StrType:            strings.ToLower(jcfg.TriggerType),
		WildSymbols:        jcfg.WildSymbols,
		MinSymbolNum:       jcfg.MinSymbolNum,
		TargetMask:         jcfg.TargetMask,
		MapBranchs:         make(map[int]*BranchNode),
		IsCheckEmptySymbol: jcfg.IsCheckEmptySymbol,
	}

	if len(jcfg.Symbols) > 0 {
		cfg.Symbols = jcfg.Symbols
	} else {
		cfg.Symbols = jcfg.TriggerSymbols
	}

	return cfg
}

func parseReelTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseReelTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseReelTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonReelTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseReelTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseReelTriggerControllers(ctrls)
		if err != nil {
			goutils.Error("parseReelTrigger:parseReelTriggerControllers",
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
		Type: ReelTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
