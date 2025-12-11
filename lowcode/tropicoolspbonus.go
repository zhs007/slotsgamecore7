package lowcode

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"log/slog"

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

const TropiCoolSPBonusTypeName = "TropiCoolSPBonus"

type TropiCoolSPBonusData struct {
	BasicComponentData
	spBonusX int
	cfg      *TropiCoolSPBonusConfig
}

// OnNewGame - reset state
func (cd *TropiCoolSPBonusData) OnNewGame(gameProp *GameProperty, component IComponent) {
	cd.BasicComponentData.OnNewGame(gameProp, component)

	cd.spBonusX = -1
}

// Clone - shallow clone
func (cd *TropiCoolSPBonusData) Clone() IComponentData {
	return &TropiCoolSPBonusData{
		BasicComponentData: cd.CloneBasicComponentData(),
		cfg:                cd.cfg,
	}
}

// BuildPBComponentData - build protobuf data
func (cd *TropiCoolSPBonusData) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: cd.BuildPBBasicComponentData(),
	}
}

// TropiCoolSPBonusConfig - placeholder configuration
type TropiCoolSPBonusConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SpBonusSymbol        string   `yaml:"spBonusSymbol" json:"spBonusSymbol"`
	SpBonusSymbolCode    int      `yaml:"-" json:"-"`
	SpBonusSymbolCode2   int      `yaml:"-" json:"-"`
	SpSymbols            []string `yaml:"spSymbols" json:"spSymbols"`
	SpSymbolCodes        []int    `yaml:"-" json:"-"`
	GenGigaSymbols2      string   `yaml:"genGigaSymbols2" json:"genGigaSymbols2"`
	HoldSymbols          []string `yaml:"holdSymbols" json:"holdSymbols"`
	HoldSymbolCodes      []int    `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *TropiCoolSPBonusConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type TropiCoolSPBonus struct {
	*BasicComponent `json:"-"`
	Config          *TropiCoolSPBonusConfig `json:"config"`
}

// Init - load from file
func (gen *TropiCoolSPBonus) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("TropiCoolSPBonus.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &TropiCoolSPBonusConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("TropiCoolSPBonus.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *TropiCoolSPBonus) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*TropiCoolSPBonusConfig)
	gen.Config.ComponentType = TropiCoolSPBonusTypeName

	sc, isok := pool.DefaultPaytables.MapSymbols[gen.Config.SpBonusSymbol]
	if !isok {
		goutils.Error("GenTropiCoolSPSymbols.InitEx:SpBonusSymbol",
			slog.String("SpBonusSymbol", gen.Config.SpBonusSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.SpBonusSymbolCode = sc

	sc, isok = pool.DefaultPaytables.MapSymbols[fmt.Sprintf("%v_%v", gen.Config.SpBonusSymbol, 2)]
	if !isok {
		goutils.Error("GenTropiCoolSPSymbols.InitEx:SpBonusSymbol",
			slog.String("SpBonusSymbol", gen.Config.SpBonusSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.SpBonusSymbolCode2 = sc

	for _, s := range gen.Config.SpSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenTropiCoolSPSymbols.InitEx:Symbol",
				slog.String("SpSymbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		gen.Config.SpSymbolCodes = append(gen.Config.SpSymbolCodes, sc)
	}

	for _, s := range gen.Config.HoldSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenTropiCoolSPSymbols.InitEx:Symbol",
				slog.String("HoldSymbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		gen.Config.HoldSymbolCodes = append(gen.Config.HoldSymbolCodes, sc)
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *TropiCoolSPBonus) canDownSymbol(tgs *sgc7game.GameScene, x, y int, ggcd *GenGigaSymbols2Data) bool {
	if y+1 >= tgs.Height {
		return false
	}

	ny := y + 1
	if tgs.Arr[x][ny] != -1 && !ggcd.cfg.isBroken(tgs.Arr[x][ny]) {
		newgigadata := ggcd.getGigaData(x, ny)
		if newgigadata != nil {
			return gen.canDownGiga(tgs, newgigadata, ggcd)
		}

		return gen.canDownSymbol(tgs, x, ny, ggcd)
	}

	return true
}

func (gen *TropiCoolSPBonus) canDownGiga(tgs *sgc7game.GameScene, gigadata *gigaData, ggcd *GenGigaSymbols2Data) bool {
	if gigadata.Y+gigadata.Height-1 >= tgs.Height-1 {
		return false
	}

	y := gigadata.Y + gigadata.Height
	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		if tgs.Arr[x][y] != -1 && !ggcd.cfg.isBroken(tgs.Arr[x][y]) {

			newgigadata := ggcd.getGigaData(x, y)
			if newgigadata != nil {
				candrop := gen.canDownGiga(tgs, newgigadata, ggcd)
				if !candrop {
					return false
				}
			} else {
				candrop := gen.canDownSymbol(tgs, x, y, ggcd)
				if !candrop {
					return false
				}
			}
		}
	}

	return true
}

func (gen *TropiCoolSPBonus) downSymbol(tgs *sgc7game.GameScene, x, y int, ggcd *GenGigaSymbols2Data) {
	ny := y + 1
	if tgs.Arr[x][ny] != -1 && !ggcd.cfg.isBroken(tgs.Arr[x][ny]) {
		newgigadata := ggcd.getGigaData(x, ny)
		if newgigadata != nil {
			gen.downGiga(tgs, newgigadata, ggcd)
		} else {
			gen.downSymbol(tgs, x, ny, ggcd)
		}
	} else {
		tgs.Arr[x][ny] = tgs.Arr[x][y]
		tgs.Arr[x][y] = -1
	}
}

func (gen *TropiCoolSPBonus) downGiga(tgs *sgc7game.GameScene, gigadata *gigaData, ggcd *GenGigaSymbols2Data) {
	y := gigadata.Y + gigadata.Height
	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		if tgs.Arr[x][y] != -1 && !ggcd.cfg.isBroken(tgs.Arr[x][y]) {
			newgigadata := ggcd.getGigaData(x, y)
			if newgigadata != nil {
				gen.downGiga(tgs, newgigadata, ggcd)
			} else {
				gen.downSymbol(tgs, x, y, ggcd)
			}
		}
	}

	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		for y := gigadata.Y; y < gigadata.Y+gigadata.Height; y++ {
			tgs.Arr[x][y] = -1
		}
	}

	gigadata.Y++

	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		for y := gigadata.Y; y < gigadata.Y+gigadata.Height; y++ {
			tgs.Arr[x][y] = gigadata.CurSymbolCode
		}
	}
}

func (gen *TropiCoolSPBonus) dropdownGiga(_ *GameProperty, ngs *sgc7game.GameScene, gigacd *GenGigaSymbols2Data, lst []*gigaData) bool {
	isdown := false

	sort.Slice(lst, func(i, j int) bool {
		return lst[i].getBottom() > lst[j].getBottom()
	})

	for _, v := range lst {
		for {
			candrop := gen.canDownGiga(ngs, v, gigacd)
			if !candrop {
				break
			}

			gen.downGiga(ngs, v, gigacd)

			isdown = true
		}
	}

	// for _, v := range gigacd.gigaData {
	// 	cy := gigacd.calcDropdown(ngs, v)
	// 	if cy != v.Y {
	// 		isdown = true

	// 		for tx := v.X; tx <= v.X+v.Width-1; tx++ {
	// 			for ty := v.Y; ty <= v.Y+v.Height-1; ty++ {
	// 				ngs.Arr[tx][ty] = -1
	// 			}
	// 		}

	// 		v.Y = cy

	// 		for tx := v.X; tx <= v.X+v.Width-1; tx++ {
	// 			for ty := v.Y; ty <= v.Y+v.Height-1; ty++ {
	// 				ngs.Arr[tx][ty] = v.CurSymbolCode
	// 			}
	// 		}
	// 	}
	// }

	return isdown
}

func (gen *TropiCoolSPBonus) refillGiga(gameProp *GameProperty, ngs *sgc7game.GameScene, gigacd *GenGigaSymbols2Data) {
	gigacd.sortGigaData()

	for _, v := range gigacd.gigaData {
		for ty := v.Y + v.Height; ty < ngs.Height; ty++ {
			for tx := v.X; tx < v.X+v.Width-1; tx++ {
				if ngs.Arr[tx][ty] == -1 {

					ngs.Arr[tx][ty] = v.CurSymbolCode
				}
			}
		}
	}
}

func (gen *TropiCoolSPBonus) dropdown(gameProp *GameProperty, ngs *sgc7game.GameScene) error {
	gigaicd := gameProp.GetComponentDataWithName(gen.Config.GenGigaSymbols2)
	if gigaicd == nil {
		goutils.Error("TropiCoolSPBonus.dropdown:GetComponentDataWithName",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	gigacd, isok := gigaicd.(*GenGigaSymbols2Data)
	if !isok {
		goutils.Error("TropiCoolSPBonus.dropdown:GenGigaSymbols2Data",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	for {
		lst := []*gigaData{}

		for _, v := range gigacd.gigaData {
			isget := false

			for tx := v.X; tx < v.X+v.Width; tx++ {
				for ty := ngs.Height - 1; ty >= v.Y+v.Height; ty-- {
					if ngs.Arr[tx][ty] == -1 {
						lst = append(lst, v)

						isget = true

						break
					}
				}

				if isget {
					break
				}
			}
		}

		for x := range ngs.Arr {
			for y := len(ngs.Arr[x]) - 1; y >= 0; {
				if ngs.Arr[x][y] == -1 {
					hass := false
					for y1 := y - 1; y1 >= 0; y1-- {
						// if giga then break
						if gigacd.getGigaData(x, y1) != nil {
							break
						}

						if ngs.Arr[x][y1] != -1 {
							ngs.Arr[x][y] = ngs.Arr[x][y1]
							ngs.Arr[x][y1] = -1

							hass = true
							y--
							break
						}
					}

					if !hass {
						break
					}
				} else {
					y--
				}
			}
		}

		isdrop := gen.dropdownGiga(gameProp, ngs, gigacd, lst)
		if !isdrop {
			break
		}
	}

	gen.refillGiga(gameProp, ngs, gigacd)

	return nil
}

func (gen *TropiCoolSPBonus) procBonus(gameProp *GameProperty, ngs *sgc7game.GameScene, gp *GameParams, cd *TropiCoolSPBonusData) error {

	if cd.spBonusX < 0 {
		spx := -1
		for x, arr := range ngs.Arr {
			for y, symbolCode := range arr {
				if symbolCode == gen.Config.SpBonusSymbolCode || symbolCode == gen.Config.SpBonusSymbolCode2 {
					spx = x
				} else if !slices.Contains(gen.Config.HoldSymbolCodes, symbolCode) {
					ngs.Arr[x][y] = -1
				}
			}
		}

		if spx < 0 {
			goutils.Error("TropiCoolSPBonus.procBonus:NoSPBonusSymbol",
				slog.String("component", gen.GetName()),
				goutils.Err(ErrInvalidComponentData))

			return ErrInvalidComponentData
		}

		cd.spBonusX = spx

		err := gen.dropdown(gameProp, ngs)
		if err != nil {
			goutils.Error("TropiCoolSPBonus.procBonus:dropdown",
				slog.String("component", gen.GetName()),
				goutils.Err(err))

			return err
		}
	} else {
		for x, arr := range ngs.Arr {
			for y, symbolCode := range arr {
				if !slices.Contains(gen.Config.HoldSymbolCodes, symbolCode) {
					ngs.Arr[x][y] = -1
				}
			}
		}

		cd.spBonusX++

		err := gen.dropdown(gameProp, ngs)
		if err != nil {
			goutils.Error("TropiCoolSPBonus.procBonus:dropdown",
				slog.String("component", gen.GetName()),
				goutils.Err(err))

			return err
		}

		for y := ngs.Height - 1; y >= 0; y-- {
			if ngs.Arr[cd.spBonusX][y] == -1 {
				if y-1 >= 0 {
					ngs.Arr[cd.spBonusX][y] = gen.Config.SpBonusSymbolCode2
					ngs.Arr[cd.spBonusX][y-1] = gen.Config.SpBonusSymbolCode2
				} else {
					ngs.Arr[cd.spBonusX][y] = gen.Config.SpBonusSymbolCode
				}

				break
			}
		}
	}

	return nil
}

// OnPlayGame - placeholder: do nothing
func (gen *TropiCoolSPBonus) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd, isok := icd.(*TropiCoolSPBonusData)
	if !isok {
		goutils.Error("TropiCoolSPBonus.OnPlayGame:TropiCoolSPBonusData",
			slog.String("component", gen.GetName()),
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	if gs == nil {
		goutils.Error("TropiCoolSPBonus.OnPlayGame:GetTargetScene3",
			slog.String("component", gen.GetName()),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	err := gen.procBonus(gameProp, ngs, gp, cd)
	if err != nil {
		goutils.Error("TropiCoolSPBonus.OnPlayGame:procBonus",
			slog.String("component", gen.GetName()),
			goutils.Err(err))

		return "", err
	}

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - no-op
func (gen *TropiCoolSPBonus) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData - return base component data
func (gen *TropiCoolSPBonus) NewComponentData() IComponentData {
	return &TropiCoolSPBonusData{
		cfg: gen.Config,
	}
}

func NewTropiCoolSPBonus(name string) IComponent {
	return &TropiCoolSPBonus{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "spBonusSymbol": "EL",
// "spSymbols": [
//
//	"B1",
//	"B2"
//
// ],
// "srcGenGigaSymbols2": "bgsp-gengiga"
// "holdSymbols": [
//
//	"WL",
//	"LW2",
//	"LW3",
//	"SC"
//
// ]
type jsonTropiCoolSPBonus struct {
	SpBonusSymbol   string   `json:"spBonusSymbol"`
	SpSymbols       []string `json:"spSymbols"`
	GenGigaSymbols2 string   `json:"genGigaSymbols2"`
	HoldSymbols     []string `json:"holdSymbols"`
}

func (j *jsonTropiCoolSPBonus) build() *TropiCoolSPBonusConfig {
	return &TropiCoolSPBonusConfig{
		SpBonusSymbol:   j.SpBonusSymbol,
		SpSymbols:       slices.Clone(j.SpSymbols),
		GenGigaSymbols2: j.GenGigaSymbols2,
	}
}

func parseTropiCoolSPBonus(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseTropiCoolSPBonus:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseTropiCoolSPBonus:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonTropiCoolSPBonus{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseTropiCoolSPBonus:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: TropiCoolSPBonusTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
