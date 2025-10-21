package lowcode

import (
	"log/slog"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const DropDownSymbols2TypeName = "dropDownSymbols2"

type DropDownSymbols2Type int

const (
	DDS2TypeNormal           DropDownSymbols2Type = 0
	DDS2TypeHexGridStaggered DropDownSymbols2Type = 1
)

func parseDropDownSymbols2Type(strType string) DropDownSymbols2Type {
	if strType == "hexgridstaggered" {
		return DDS2TypeHexGridStaggered
	}

	return DDS2TypeNormal
}

// DropDownSymbols2Config - configuration for DropDownSymbols2
type DropDownSymbols2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	HoldSymbols          []string             `yaml:"holdSymbols" json:"holdSymbols"`                   // 不需要下落的symbol
	HoldSymbolCodes      []int                `yaml:"-" json:"-"`                                       // 不需要下落的symbol
	IsNeedProcSymbolVals bool                 `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int                  `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	StrType              string               `yaml:"type" json:"type"`                                 // 类型
	Type                 DropDownSymbols2Type `yaml:"-" json:"-"`                                       // 类型
	RowMask              string               `yaml:"rowMask" json:"rowMask"`                           // rowMask
	MapAwards            map[string][]*Award  `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *DropDownSymbols2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type DropDownSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *DropDownSymbols2Config `json:"config"`
}

// Init -
func (dropDownSymbols *DropDownSymbols2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("DropDownSymbols2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &DropDownSymbols2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("DropDownSymbols2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return dropDownSymbols.InitEx(cfg, pool)
}

// InitEx -
func (dropDownSymbols *DropDownSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	dropDownSymbols.Config = cfg.(*DropDownSymbols2Config)
	dropDownSymbols.Config.ComponentType = DropDownSymbols2TypeName

	dropDownSymbols.Config.Type = parseDropDownSymbols2Type(dropDownSymbols.Config.StrType)

	for _, v := range dropDownSymbols.Config.HoldSymbols {
		dropDownSymbols.Config.HoldSymbolCodes = append(dropDownSymbols.Config.HoldSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, awards := range dropDownSymbols.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	dropDownSymbols.onInit(&dropDownSymbols.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (dropDownSymbols *DropDownSymbols2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := dropDownSymbols.Config.MapAwards[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

func (dropDownSymbols *DropDownSymbols2) procNormalWithOS(ngs *sgc7game.GameScene, nos *sgc7game.GameScene) error {

	for x, arr := range ngs.Arr {
		for y := len(arr) - 1; y >= 0; {
			if arr[y] == -1 {
				hass := false
				for y1 := y - 1; y1 >= 0; y1-- {
					if arr[y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
						arr[y] = arr[y1]
						arr[y1] = -1

						nos.Arr[x][y] = nos.Arr[x][y1]
						nos.Arr[x][y1] = dropDownSymbols.Config.EmptySymbolVal

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

	return nil
}

func (dropDownSymbols *DropDownSymbols2) procNormal(gameProp *GameProperty, gs *sgc7game.GameScene) (*sgc7game.GameScene, error) {

	ngs := gs

	for x := range ngs.Arr {
		for y := len(ngs.Arr[x]) - 1; y >= 0; {
			if ngs.Arr[x][y] == -1 {
				hass := false
				for y1 := y - 1; y1 >= 0; y1-- {
					if ngs.Arr[x][y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

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

	return ngs, nil
}

func (dropDownSymbols *DropDownSymbols2) procHexGridStaggeredWithOS(ngs *sgc7game.GameScene, nos *sgc7game.GameScene) error {

	for x, arr := range ngs.Arr {
		for y := len(arr) - 1; y >= 0; {
			if arr[y] == -1 {
				hass := false
				for y1 := y - 1; y1 >= 0; y1-- {
					if arr[y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
						arr[y] = arr[y1]
						arr[y1] = -1

						nos.Arr[x][y] = nos.Arr[x][y1]
						nos.Arr[x][y1] = dropDownSymbols.Config.EmptySymbolVal

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

	return nil
}

func (dropDownSymbols *DropDownSymbols2) procHexGridStaggered(gameProp *GameProperty, gs *sgc7game.GameScene) (bool, *sgc7game.GameScene, error) {

	ngs := gs

	// 有 rowMask 时很复杂,最后的下落不能算trigger,应该算refill
	if dropDownSymbols.Config.RowMask != "" {
		imaskd := gameProp.GetComponentDataWithName(dropDownSymbols.Config.RowMask)
		if imaskd == nil {
			goutils.Error("DropDownSymbols2.getSrcPos:RowMask:imaskd==nil",
				goutils.Err(ErrInvalidComponentConfig))

			return false, nil, ErrInvalidComponentConfig
		}

		maskarr := imaskd.GetMask()

		// 先正常下落，再处理滚动
		for x := range ngs.Arr {
			for y := len(ngs.Arr[x]) - 1; y >= 0; {
				if !maskarr[y] {
					y--

					continue
				}

				if ngs.Arr[x][y] == -1 {
					hass := false
					for y1 := y - 1; y1 >= 0; y1-- {
						if ngs.Arr[x][y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
							if ngs == gs {
								ngs = gs.CloneEx(gameProp.PoolScene)
							}

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

		isRoll := false

		for x := 1; x < ngs.Width; x++ {
			for y := len(ngs.Arr[x]) - 2; y >= 0; y-- {
				if ngs.Arr[x][y] == -1 {
					break
				}

				if !maskarr[y] {
					continue
				}

				if x%2 == 1 {
					if ngs.Arr[x-1][y] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						isRoll = true

						ngs.Arr[x-1][y] = ngs.Arr[x][y]

						ngs.Arr[x][y] = -1

						for ty := y - 1; ty >= 0; ty-- {
							if ngs.Arr[x][ty] == -1 {
								break
							}

							ngs.Arr[x-1][ty] = ngs.Arr[x][ty]

							ngs.Arr[x][ty] = -1
						}
					}
				} else {
					if ngs.Arr[x-1][y+1] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						isRoll = true

						ngs.Arr[x-1][y+1] = ngs.Arr[x][y]

						ngs.Arr[x][y] = -1

						for ty := y - 1; ty >= 0; ty-- {
							if ngs.Arr[x][ty] == -1 {
								break
							}

							ngs.Arr[x-1][ty+1] = ngs.Arr[x][ty]

							ngs.Arr[x][ty] = -1
						}
					}
				}
			}
		}

		for x := ngs.Width - 2; x >= 0; x-- {
			for y := len(ngs.Arr[x]) - 2; y >= 0; y-- {
				if ngs.Arr[x][y] == -1 {
					break
				}

				if !maskarr[y] {
					continue
				}

				if x%2 == 1 {
					if ngs.Arr[x+1][y] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						isRoll = true

						ngs.Arr[x+1][y] = ngs.Arr[x][y]

						ngs.Arr[x][y] = -1

						for ty := y - 1; ty >= 0; ty-- {
							if ngs.Arr[x][ty] == -1 {
								break
							}

							ngs.Arr[x+1][ty] = ngs.Arr[x][ty]

							ngs.Arr[x][ty] = -1
						}
					}
				} else {
					if ngs.Arr[x+1][y+1] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						isRoll = true

						ngs.Arr[x+1][y+1] = ngs.Arr[x][y]

						ngs.Arr[x][y] = -1

						for ty := y - 1; ty >= 0; ty-- {
							if ngs.Arr[x][ty] == -1 {
								break
							}

							ngs.Arr[x+1][ty+1] = ngs.Arr[x][ty]

							ngs.Arr[x][ty] = -1
						}
					}
				}
			}
		}

		if isRoll {
			return true, ngs, nil
		}

		for x := range ngs.Arr {
			for y := len(ngs.Arr[x]) - 1; y >= 0; {
				if ngs.Arr[x][y] == -1 {
					hass := false
					for y1 := y - 1; y1 >= 0; y1-- {
						if ngs.Arr[x][y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
							if ngs == gs {
								ngs = gs.CloneEx(gameProp.PoolScene)
							}

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

		return false, ngs, nil
	}

	// 先正常下落，再处理滚动
	for x := range ngs.Arr {
		for y := len(ngs.Arr[x]) - 1; y >= 0; {
			if ngs.Arr[x][y] == -1 {
				hass := false
				for y1 := y - 1; y1 >= 0; y1-- {
					if ngs.Arr[x][y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

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

	// 滚动时先 x 从 1 开始扫(从下往上)，看能不能向左滚，如果能滚就直接处理，空的位置可以留下来，后面的就可以一个方向滚动; 如果一个图标滚动了,上面的图标也都应该一起动;滚动马上执行,这样后面的图标才有位置动
	// 再 x 从 0 开始扫，前面已经滚动过的轴跳过，看能不能向右滚，如果能滚就直接处理，空的位置可以留下来，后面的就可以一个方向滚动
	// 就这个顺序迭代
	isRoll := false

	for x := 1; x < ngs.Width; x++ {
		for y := len(ngs.Arr[x]) - 2; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				break
			}

			if x%2 == 1 {
				if ngs.Arr[x-1][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					isRoll = true

					ngs.Arr[x-1][y] = ngs.Arr[x][y]

					ngs.Arr[x][y] = -1

					for ty := y - 1; ty >= 0; ty-- {
						if ngs.Arr[x][ty] == -1 {
							break
						}

						ngs.Arr[x-1][ty] = ngs.Arr[x][ty]

						ngs.Arr[x][ty] = -1
					}
				}
			} else {
				if ngs.Arr[x-1][y+1] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					isRoll = true

					ngs.Arr[x-1][y+1] = ngs.Arr[x][y]

					ngs.Arr[x][y] = -1

					for ty := y - 1; ty >= 0; ty-- {
						if ngs.Arr[x][ty] == -1 {
							break
						}

						ngs.Arr[x-1][ty+1] = ngs.Arr[x][ty]

						ngs.Arr[x][ty] = -1
					}
				}
			}
		}
	}

	for x := ngs.Width - 2; x >= 0; x-- {
		for y := len(ngs.Arr[x]) - 2; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				break
			}

			if x%2 == 1 {
				if ngs.Arr[x+1][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					isRoll = true

					ngs.Arr[x+1][y] = ngs.Arr[x][y]

					ngs.Arr[x][y] = -1

					for ty := y - 1; ty >= 0; ty-- {
						if ngs.Arr[x][ty] == -1 {
							break
						}

						ngs.Arr[x+1][ty] = ngs.Arr[x][ty]

						ngs.Arr[x][ty] = -1
					}
				}
			} else {
				if ngs.Arr[x+1][y+1] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					isRoll = true

					ngs.Arr[x+1][y+1] = ngs.Arr[x][y]

					ngs.Arr[x][y] = -1

					for ty := y - 1; ty >= 0; ty-- {
						if ngs.Arr[x][ty] == -1 {
							break
						}

						ngs.Arr[x+1][ty+1] = ngs.Arr[x][ty]

						ngs.Arr[x][ty] = -1
					}
				}
			}
		}
	}

	return isRoll, ngs, nil
}

// playgame
func (dropDownSymbols *DropDownSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil
	bcd.UsedOtherScenes = nil

	gs := dropDownSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("DropDownSymbols2.OnPlayGame",
			goutils.Err(ErrInvalidScene))

		return "", ErrInvalidScene
	}

	if !gs.HasSymbol(-1) {
		nc := dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	var os *sgc7game.GameScene
	if dropDownSymbols.Config.IsNeedProcSymbolVals {
		os = dropDownSymbols.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	if os != nil {
		ngs := gs.CloneEx(gameProp.PoolScene)
		nos := os.CloneEx(gameProp.PoolScene)

		switch dropDownSymbols.Config.Type {
		case DDS2TypeNormal:
			err := dropDownSymbols.procNormalWithOS(ngs, nos)
			if err != nil {
				goutils.Error("DropDownSymbols2.OnPlayGame:procNormalWithOS",
					goutils.Err(err))

				return "", err
			}
		case DDS2TypeHexGridStaggered:
			err := dropDownSymbols.procHexGridStaggeredWithOS(ngs, nos)
			if err != nil {
				goutils.Error("DropDownSymbols2.OnPlayGame:procHexGridStaggeredWithOS",
					goutils.Err(err))

				return "", err
			}
		}

		dropDownSymbols.AddOtherScene(gameProp, curpr, nos, bcd)

		dropDownSymbols.AddScene(gameProp, curpr, ngs, bcd)

		dropDownSymbols.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")
	} else {
		switch dropDownSymbols.Config.Type {
		case DDS2TypeNormal:
			ngs, err := dropDownSymbols.procNormal(gameProp, gs)
			if err != nil {
				goutils.Error("DropDownSymbols2.OnPlayGame:procNormal",
					goutils.Err(err))

				return "", err
			}

			if ngs != gs {
				dropDownSymbols.AddScene(gameProp, curpr, ngs, bcd)

				dropDownSymbols.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")

				return dropDownSymbols.onStepEnd(gameProp, curpr, gp, ""), nil
			}
		case DDS2TypeHexGridStaggered:
			istrigger, ngs, err := dropDownSymbols.procHexGridStaggered(gameProp, gs)
			if err != nil {
				goutils.Error("DropDownSymbols2.OnPlayGame:procHexGridStaggered",
					goutils.Err(err))

				return "", err
			}

			if ngs != gs {
				dropDownSymbols.AddScene(gameProp, curpr, ngs, bcd)

				if istrigger {
					dropDownSymbols.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")
				}

				return dropDownSymbols.onStepEnd(gameProp, curpr, gp, ""), nil
			}
		}

		nc := dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (dropDownSymbols *DropDownSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after dropDownSymbols2", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// EachUsedResults -
func (dropDownSymbols *DropDownSymbols2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewDropDownSymbols2(name string) IComponent {
	return &DropDownSymbols2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "isNeedProcSymbolVals": false,
// "type": "hexGridStaggered"
// "rowMask": "mask-height4"
type jsonDropDownSymbols2 struct {
	HoldSymbols          []string `json:"holdSymbols"`          // 不需要下落的symbol
	IsNeedProcSymbolVals bool     `json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int      `json:"emptySymbolVal"`       // 空的symbolVal是什么
	Type                 string   `json:"type"`                 // 类型
	RowMask              string   `json:"rowMask"`              // rowMask
}

func (jcfg *jsonDropDownSymbols2) build() *DropDownSymbols2Config {
	cfg := &DropDownSymbols2Config{
		HoldSymbols:          jcfg.HoldSymbols,
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
		StrType:              strings.ToLower(jcfg.Type),
		RowMask:              jcfg.RowMask,
	}

	return cfg
}

func parseDropDownSymbols2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseDropDownSymbols2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseDropDownSymbols2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonDropDownSymbols2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseDropDownSymbols2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseDropDownSymbols2:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: DropDownSymbols2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
