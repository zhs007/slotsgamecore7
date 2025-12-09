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
	"gopkg.in/yaml.v2"
)

const TropiCoolExchangeTypeName = "tropiCoolExchange"

// TropiCoolExchangeConfig - placeholder configuration for TropiCoolExchange
type TropiCoolExchangeConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	RowSymbol            string                `yaml:"rowSymbol" json:"rowSymbol"`
	RowSymbolCode        int                   `yaml:"-" json:"-"`
	ColSymbol            string                `yaml:"colSymbol" json:"colSymbol"`
	ColSymbolCode        int                   `yaml:"-" json:"-"`
	GenGigaSymbols2      string                `yaml:"genGigaSymbols2" json:"genGigaSymbols2"`
	Weight               string                `yaml:"weight" json:"weight"`
	WeightVW             *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	SymbolCodes          []int                 `yaml:"-" json:"-"`
	WildSymbol           string                `yaml:"wildSymbol" json:"wildSymbol"`
	WildSymbolCode       int                   `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *TropiCoolExchangeConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type TropiCoolExchange struct {
	*BasicComponent `json:"-"`
	Config          *TropiCoolExchangeConfig `json:"config"`
}

// Init - load from file
func (gen *TropiCoolExchange) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("TropiCoolExchange.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &TropiCoolExchangeConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("TropiCoolExchange.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *TropiCoolExchange) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*TropiCoolExchangeConfig)
	gen.Config.ComponentType = TropiCoolExchangeTypeName

	sc, isok := pool.DefaultPaytables.MapSymbols[gen.Config.RowSymbol]
	if !isok {
		goutils.Error("TropiCoolExchange.InitEx:RowSymbol",
			slog.String("symbol", gen.Config.RowSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.RowSymbolCode = sc

	sc, isok = pool.DefaultPaytables.MapSymbols[gen.Config.ColSymbol]
	if !isok {
		goutils.Error("TropiCoolExchange.InitEx:ColSymbol",
			slog.String("symbol", gen.Config.ColSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.ColSymbolCode = sc

	vw, err := pool.LoadIntWeights(gen.Config.Weight, true)
	if err != nil {
		goutils.Error("TropiCoolExchange.InitEx:LoadIntWeights",
			slog.String("weight", gen.Config.Weight),
			goutils.Err(err))

		return err
	}

	gen.Config.WeightVW = vw

	gen.Config.SymbolCodes = vw.GetIntVals()

	sc, isok = pool.DefaultPaytables.MapSymbols[gen.Config.WildSymbol]
	if !isok {
		goutils.Error("TropiCoolExchange.InitEx:WildSymbol",
			slog.String("symbol", gen.Config.WildSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.WildSymbolCode = sc

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *TropiCoolExchange) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, err := gameProp.Pool.LoadIntWeights(str, true)
		if err != nil {
			goutils.Error("TropiCoolExchange.getWeight:LoadIntWeights",
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return gen.Config.WeightVW
}

func (gen *TropiCoolExchange) getGenGigaSymbols2Data(gameProp *GameProperty) (*GenGigaSymbols2Data, error) {
	gigaicd := gameProp.GetComponentDataWithName(gen.Config.GenGigaSymbols2)
	if gigaicd == nil {
		goutils.Error("TropiCoolExchange.getGenGigaSymbols2Data:GetComponentDataWithName",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	ggcd, isok := gigaicd.(*GenGigaSymbols2Data)
	if !isok {
		goutils.Error("TropiCoolExchange.getGenGigaSymbols2Data:GenGigaSymbols2Data",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	return ggcd, nil
}

func (gen *TropiCoolExchange) procRowGiga(plugin sgc7plugin.IPlugin, ngs *sgc7game.GameScene,
	ts *sgc7game.GameScene, ggd *gigaData, ggcd *GenGigaSymbols2Data, vw *sgc7game.ValWeights2) error {
	syms := make([]int, 0, ngs.Width*ggd.Width)

	for cx := 0; cx < ngs.Width; cx++ {
		if cx >= ggd.X && cx < ggd.X+ggd.Width {
			continue
		}

		for cy := ggd.Y; cy < ggd.Y+ggd.Height; cy++ {
			cggd := ggcd.getGigaData(cx, cy)
			if cggd != nil {
				if slices.Contains(gen.Config.SymbolCodes, cggd.SymbolCode) {
					syms = append(syms, cggd.SymbolCode)
				}
			} else {
				if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[cx][cy]) {
					syms = append(syms, ngs.Arr[cx][cy])
				}
			}
		}
	}

	nvw := vw.CloneWithIntArray(syms)
	cr, err := nvw.RandVal(plugin)
	if err != nil {
		goutils.Error("TropiCoolExchange.procRowGiga:RandVal",
			goutils.Err(err))

		return err
	}

	sc := cr.Int()

	ggd.chgSymbol(ngs, sc, ggcd.cfg)
	for tx := ggd.X; tx < ggd.X+ggd.Width; tx++ {
		for ty := ggd.Y; ty < ggd.Y+ggd.Height; ty++ {
			ts.Arr[tx][ty]++
		}
	}

	for cx := 0; cx < ngs.Width; cx++ {
		if cx >= ggd.X && cx < ggd.X+ggd.Width {
			continue
		}

		for cy := ggd.Y; cy < ggd.Y+ggd.Height; cy++ {
			cggd := ggcd.getGigaData(cx, cy)
			if cggd != nil {
				if slices.Contains(gen.Config.SymbolCodes, cggd.SymbolCode) {
					cggd.chgSymbol(ngs, sc, ggcd.cfg)

					for tx := cggd.X; tx < cggd.X+cggd.Width; tx++ {
						for ty := cggd.Y; ty < cggd.Y+cggd.Height; ty++ {
							ts.Arr[tx][ty]++
						}
					}
				} else if cggd.SymbolCode == gen.Config.RowSymbolCode {
					cggd.chgSymbol(ngs, gen.Config.ColSymbolCode, ggcd.cfg)

					for tx := cggd.X; tx < cggd.X+cggd.Width; tx++ {
						for ty := cggd.Y; ty < cggd.Y+cggd.Height; ty++ {
							ts.Arr[tx][ty]++
						}
					}
				}
			} else {
				if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[cx][cy]) {
					ngs.Arr[cx][cy] = sc

					ts.Arr[cx][cy]++
				} else if ngs.Arr[cx][cy] == gen.Config.RowSymbolCode {
					ngs.Arr[cx][cy] = gen.Config.ColSymbolCode

					ts.Arr[cx][cy]++
				}
			}
		}
	}

	return nil
}

func (gen *TropiCoolExchange) procColGiga(plugin sgc7plugin.IPlugin, ngs *sgc7game.GameScene,
	ts *sgc7game.GameScene, ggd *gigaData, ggcd *GenGigaSymbols2Data, vw *sgc7game.ValWeights2) error {
	syms := make([]int, 0, ngs.Height*ggd.Width)

	for cy := 0; cy < ngs.Height; cy++ {
		if cy >= ggd.Y && cy < ggd.Y+ggd.Height {
			continue
		}

		for cx := ggd.X; cx < ggd.X+ggd.Width; cx++ {
			cggd := ggcd.getGigaData(cx, cy)
			if cggd != nil {
				if slices.Contains(gen.Config.SymbolCodes, cggd.SymbolCode) {
					syms = append(syms, cggd.SymbolCode)
				}
			} else {
				if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[cx][cy]) {
					syms = append(syms, ngs.Arr[cx][cy])
				}
			}
		}
	}

	nvw := vw.CloneWithIntArray(syms)
	cr, err := nvw.RandVal(plugin)
	if err != nil {
		goutils.Error("TropiCoolExchange.procColGiga:RandVal",
			goutils.Err(err))

		return err
	}

	sc := cr.Int()

	ggd.chgSymbol(ngs, sc, ggcd.cfg)
	for tx := ggd.X; tx < ggd.X+ggd.Width; tx++ {
		for ty := ggd.Y; ty < ggd.Y+ggd.Height; ty++ {
			ts.Arr[tx][ty]++
		}
	}

	for cy := 0; cy < ngs.Height; cy++ {
		if cy >= ggd.Y && cy < ggd.Y+ggd.Height {
			continue
		}

		for cx := ggd.X; cx < ggd.X+ggd.Width; cx++ {
			cggd := ggcd.getGigaData(cx, cy)
			if cggd != nil {
				if slices.Contains(gen.Config.SymbolCodes, cggd.SymbolCode) {
					cggd.chgSymbol(ngs, sc, ggcd.cfg)

					for tx := cggd.X; tx < cggd.X+cggd.Width; tx++ {
						for ty := cggd.Y; ty < cggd.Y+cggd.Height; ty++ {
							ts.Arr[tx][ty]++
						}
					}
				}
			} else {
				if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[cx][cy]) {
					ngs.Arr[cx][cy] = sc

					ts.Arr[cx][cy]++
				}
			}
		}
	}

	return nil
}

func (gen *TropiCoolExchange) procRowSymbol(plugin sgc7plugin.IPlugin, ngs *sgc7game.GameScene, ts *sgc7game.GameScene, x, y int,
	ggcd *GenGigaSymbols2Data, vw *sgc7game.ValWeights2) error {
	syms := make([]int, 0, ngs.Width)

	for cx := 0; cx < ngs.Width; cx++ {
		ggd := ggcd.getGigaData(cx, y)
		if ggd != nil {
			if slices.Contains(gen.Config.SymbolCodes, ggd.SymbolCode) {
				syms = append(syms, ggd.SymbolCode)
			}
		} else {
			if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[cx][y]) {
				syms = append(syms, ngs.Arr[cx][y])
			}
		}
	}

	nvw := vw.CloneWithIntArray(syms)
	cr, err := nvw.RandVal(plugin)
	if err != nil {
		goutils.Error("TropiCoolExchange.procRowSymbol:RandVal",
			goutils.Err(err))

		return err
	}

	sc := cr.Int()

	for cx := 0; cx < ngs.Width; cx++ {
		ggd := ggcd.getGigaData(cx, y)
		if ggd != nil {
			if slices.Contains(gen.Config.SymbolCodes, ggd.SymbolCode) {
				ggd.chgSymbol(ngs, sc, ggcd.cfg)

				for tx := ggd.X; tx < ggd.X+ggd.Width; tx++ {
					for ty := ggd.Y; ty < ggd.Y+ggd.Height; ty++ {
						ts.Arr[tx][ty]++
					}
				}
			} else if ggd.SymbolCode == gen.Config.RowSymbolCode {
				ggd.chgSymbol(ngs, gen.Config.ColSymbolCode, ggcd.cfg)

				for tx := ggd.X; tx < ggd.X+ggd.Width; tx++ {
					for ty := ggd.Y; ty < ggd.Y+ggd.Height; ty++ {
						ts.Arr[tx][ty]++
					}
				}
			}
		} else {
			if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[cx][y]) || x == cx {
				ngs.Arr[cx][y] = sc

				ts.Arr[cx][y]++
			} else if ngs.Arr[cx][y] == gen.Config.RowSymbolCode {
				ngs.Arr[cx][y] = gen.Config.ColSymbolCode

				ts.Arr[cx][y]++
			}
		}
	}

	return nil
}

func (gen *TropiCoolExchange) procColSymbol(plugin sgc7plugin.IPlugin, ngs *sgc7game.GameScene, ts *sgc7game.GameScene, x, y int,
	ggcd *GenGigaSymbols2Data, vw *sgc7game.ValWeights2) error {
	syms := make([]int, 0, ngs.Height)

	for cy := 0; cy < ngs.Height; cy++ {
		ggd := ggcd.getGigaData(x, cy)
		if ggd != nil {
			if slices.Contains(gen.Config.SymbolCodes, ggd.SymbolCode) {
				syms = append(syms, ggd.SymbolCode)
			}
		} else {
			if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[x][cy]) {
				syms = append(syms, ngs.Arr[x][cy])
			}
		}
	}

	nvw := vw.CloneWithIntArray(syms)
	cr, err := nvw.RandVal(plugin)
	if err != nil {
		goutils.Error("TropiCoolExchange.procRowSymbol:RandVal",
			goutils.Err(err))

		return err
	}

	sc := cr.Int()

	for cy := 0; cy < ngs.Height; cy++ {
		ggd := ggcd.getGigaData(x, cy)
		if ggd != nil {
			if slices.Contains(gen.Config.SymbolCodes, ggd.SymbolCode) {
				ggd.chgSymbol(ngs, sc, ggcd.cfg)

				for tx := ggd.X; tx < ggd.X+ggd.Width; tx++ {
					for ty := ggd.Y; ty < ggd.Y+ggd.Height; ty++ {
						ts.Arr[tx][ty]++
					}
				}
			}
		} else {
			if slices.Contains(gen.Config.SymbolCodes, ngs.Arr[x][cy]) || y == cy {
				ngs.Arr[x][cy] = sc

				ts.Arr[x][cy]++
			}
		}
	}

	return nil
}

// OnPlayGame - placeholder: no-op component
func (gen *TropiCoolExchange) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// do minimal advance/step-end processing and return do-nothing
	bcd, isok := icd.(*BasicComponentData)
	if !isok {
		goutils.Error("TropiCoolExchange.OnPlayGame:BasicComponentData",
			slog.String("component", gen.GetName()),
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	bcd.UsedScenes = nil

	ggcd, err := gen.getGenGigaSymbols2Data(gameProp)
	if err != nil {
		goutils.Error("TropiCoolExchange.OnPlayGame:getGenGigaSymbols2Data",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(err))

		return "", err
	}

	gs := gen.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("TropiCoolExchange.OnPlayGame:GetTargetScene3",
			slog.String("component", gen.GetName()),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	vw := gen.getWeight(gameProp, bcd)

	ts := gameProp.PoolScene.New(gs.Width, gs.Height)

	ngs := gs.CloneEx(gameProp.PoolScene)

	isTrigger := false

	for x := 0; x < gs.Width; x++ {
		for y := 0; y < gs.Height; y++ {
			sym := gs.Arr[x][y]

			ggd := ggcd.getGigaData(x, y)
			if ggd != nil {
				switch ggd.SymbolCode {
				case gen.Config.RowSymbolCode:
					err = gen.procRowGiga(plugin, ngs, ts, ggd, ggcd, vw)
					if err != nil {
						goutils.Error("TropiCoolExchange.OnPlayGame:procRowGiga",
							goutils.Err(err))

						return "", err
					}

					isTrigger = true
				case gen.Config.ColSymbolCode:
					err = gen.procColGiga(plugin, ngs, ts, ggd, ggcd, vw)
					if err != nil {
						goutils.Error("TropiCoolExchange.OnPlayGame:procColGiga",
							goutils.Err(err))

						return "", err
					}

					isTrigger = true
				}

				continue
			}

			switch sym {
			case gen.Config.RowSymbolCode:
				err = gen.procRowSymbol(plugin, ts, ngs, x, y, ggcd, vw)
				if err != nil {
					goutils.Error("TropiCoolExchange.OnPlayGame:procRowSymbol",
						goutils.Err(err))

					return "", err
				}

				isTrigger = true
			case gen.Config.ColSymbolCode:
				err = gen.procColSymbol(plugin, ts, ngs, x, y, ggcd, vw)
				if err != nil {
					goutils.Error("TropiCoolExchange.OnPlayGame:procColSymbol",
						goutils.Err(err))

					return "", err
				}

				isTrigger = true
			}
		}
	}

	if !isTrigger {
		nc := gen.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	for x := 0; x < ts.Width; x++ {
		for y := 0; y < ts.Height; y++ {
			if ts.Arr[x][y] >= 2 {
				ngs.Arr[x][y] = gen.Config.WildSymbolCode
			}
		}
	}

	gen.AddScene(gameProp, curpr, ngs, bcd)

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *TropiCoolExchange) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData - return base component data
func (gen *TropiCoolExchange) NewComponentData() IComponentData {
	return &BasicComponentData{}
}

func NewTropiCoolExchange(name string) IComponent {
	return &TropiCoolExchange{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "rowSymbol": "RS",
// "colSymbol": "CS",
// "genGigaSymbols2": "bg-gengiga"
// "weight": "mysteryweight"
// "wildSymbol": "WL"
type jsonTropiCoolExchange struct {
	RowSymbol       string `json:"rowSymbol"`
	ColSymbol       string `json:"colSymbol"`
	GenGigaSymbols2 string `json:"genGigaSymbols2"`
	Weight          string `json:"weight"`
	WildSymbol      string `json:"wildSymbol"`
}

func (j *jsonTropiCoolExchange) build() *TropiCoolExchangeConfig {
	return &TropiCoolExchangeConfig{
		RowSymbol:       j.RowSymbol,
		ColSymbol:       j.ColSymbol,
		GenGigaSymbols2: j.GenGigaSymbols2,
		Weight:          j.Weight,
		WildSymbol:      j.WildSymbol,
	}
}

func parseTropiCoolExchange(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseTropiCoolExchange:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseTropiCoolExchange:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonTropiCoolExchange{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseTropiCoolExchange:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: TropiCoolExchangeTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
