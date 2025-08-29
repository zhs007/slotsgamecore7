package lowcode

import (
	"context"
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
	"gopkg.in/yaml.v2"
)

const GenMaskTypeName = "genMask"

type GenMaskType int

const (
	GMTypeSet    GenMaskType = 0 // set
	GMTypeNot    GenMaskType = 1 // not
	GMTypeAnd    GenMaskType = 2 // and
	GMTypeOr     GenMaskType = 3 // or
	GMTypeXor    GenMaskType = 4 // xor
	GMTypeRandom GenMaskType = 5 // random
)

func (gmt GenMaskType) isSingleMask() bool {
	return gmt == GMTypeSet || gmt == GMTypeRandom || gmt == GMTypeNot
}

func parseGenMaskType(str string) GenMaskType {
	str = strings.ToLower(str)

	switch str {
	case "not":
		return GMTypeNot
	case "and":
		return GMTypeAnd
	case "or":
		return GMTypeOr
	case "xor":
		return GMTypeXor
	case "random":
		return GMTypeRandom
	}

	return GMTypeSet
}

// GenMaskConfig - configuration for GenMask
type GenMaskConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string      `yaml:"type" json:"type"`
	Type                 GenMaskType `yaml:"-" json:"-"`
	MaskLen              int         `yaml:"maskLen" json:"maskLen"`
	OutputMask           string      `yaml:"outputMask" json:"outputMask"`
	SrcMask              []string    `yaml:"srcMask" json:"srcMask"`
	WeightValue          int         `yaml:"weightValue" json:"weightValue"`
	InitMask             []bool      `yaml:"initMask" json:"initMask"`
	Controllers          []*Award    `yaml:"controllers" json:"controllers"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *GenMaskConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenMask struct {
	*BasicComponent `json:"-"`
	Config          *GenMaskConfig `json:"config"`
}

// Init -
func (gm *GenMask) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenMask.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenMaskConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenMask.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gm.InitEx(cfg, pool)
}

// InitEx -
func (gm *GenMask) InitEx(cfg any, pool *GamePropertyPool) error {
	cfgd, ok := cfg.(*GenMaskConfig)
	if !ok {
		goutils.Error("GenMask.InitEx:InvalidConfigType",
			slog.String("got", fmt.Sprintf("%T", cfg)))

		return ErrInvalidComponent
	}

	gm.Config = cfgd
	gm.Config.ComponentType = GenMaskTypeName

	gm.Config.Type = parseGenMaskType(gm.Config.StrType)

	if gm.Config.OutputMask == "" {
		goutils.Error("GenMask.InitEx:OutputMask",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if gm.Config.Type.isSingleMask() {
		if gm.Config.Type == GMTypeRandom {
			if gm.Config.WeightValue < 0 || gm.Config.WeightValue > 10000 {
				goutils.Error("GenMask.InitEx:WeightValue",
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}
		} else {
			if len(gm.Config.SrcMask) > 1 {
				goutils.Error("GenMask.InitEx:SrcMask",
					slog.Int("len", len(gm.Config.SrcMask)),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			} else if len(gm.Config.SrcMask) == 1 {
				if len(gm.Config.InitMask) > 0 {
					goutils.Error("GenMask.InitEx:InitMask",
						slog.Int("len", len(gm.Config.InitMask)),
						goutils.Err(ErrInvalidComponentConfig))

					return ErrInvalidComponentConfig
				}
			} else {
				if len(gm.Config.InitMask) == 0 {
					goutils.Error("GenMask.InitEx:InitMask",
						slog.Int("len", len(gm.Config.InitMask)),
						goutils.Err(ErrInvalidComponentConfig))

					return ErrInvalidComponentConfig
				}
			}
		}
	} else {
		if len(gm.Config.InitMask) > 0 {
			if len(gm.Config.SrcMask) == 0 {
				goutils.Error("GenMask.InitEx:SrcMask",
					slog.Int("len", len(gm.Config.SrcMask)),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}
		} else {
			if len(gm.Config.SrcMask) < 2 {
				goutils.Error("GenMask.InitEx:SrcMask",
					slog.Int("len", len(gm.Config.SrcMask)),
					goutils.Err(ErrInvalidComponentConfig))

				return ErrInvalidComponentConfig
			}
		}
	}

	if len(gm.Config.InitMask) > 0 && len(gm.Config.InitMask) != gm.Config.MaskLen {
		goutils.Error("GenMask.InitEx:InitMask",
			slog.Int("len", len(gm.Config.InitMask)),
			slog.Int("masklen", gm.Config.MaskLen),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	for _, ctrl := range gm.Config.Controllers {
		ctrl.Init()
	}

	gm.onInit(&gm.Config.BasicComponentConfig)

	return nil
}

func (gm *GenMask) getFirstMask(gameProp *GameProperty, basicCD *BasicComponentData) []bool {
	if len(gm.Config.InitMask) > 0 {
		return gm.Config.InitMask
	}

	if len(gm.Config.SrcMask) > 0 {
		maskVal, err := gameProp.GetMask(gm.Config.SrcMask[0])
		if err != nil {
			goutils.Error("GenMask.getFirstMask:GetMask",
				goutils.Err(err))

			return nil
		}
		if len(maskVal) != gm.Config.MaskLen {
			goutils.Error("GenMask.getFirstMask:MaskLen",
				slog.Int("got", len(maskVal)),
				slog.Int("want", gm.Config.MaskLen),
				goutils.Err(ErrInvalidComponentConfig))

			return nil
		}

		return maskVal
	}

	return nil
}

func (gm *GenMask) getAllMask(gameProp *GameProperty, _ *BasicComponentData) [][]bool {
	size := len(gm.Config.SrcMask)

	if len(gm.Config.InitMask) > 0 {
		size++
	}

	masks := make([][]bool, 0, size)

	if len(gm.Config.InitMask) > 0 {
		if len(gm.Config.InitMask) != gm.Config.MaskLen {
			goutils.Error("GenMask.getAllMask:InitMaskLen",
				slog.Int("got", len(gm.Config.InitMask)),
				slog.Int("want", gm.Config.MaskLen),
				goutils.Err(ErrInvalidComponentConfig))

			return nil
		}
		masks = append(masks, gm.Config.InitMask)
	}

	for _, name := range gm.Config.SrcMask {
		maskVal, err := gameProp.GetMask(name)
		if err != nil {
			goutils.Error("GenMask.getAllMask:GetMask",
				slog.String("name", name),
				goutils.Err(err))

			return nil
		}
		if len(maskVal) != gm.Config.MaskLen {
			goutils.Error("GenMask.getAllMask:MaskLen",
				slog.String("name", name),
				slog.Int("got", len(maskVal)),
				slog.Int("want", gm.Config.MaskLen),
				goutils.Err(ErrInvalidComponentConfig))

			return nil
		}

		masks = append(masks, maskVal)
	}

	return masks
}

// OnProcControllers -
func (gm *GenMask) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {

	gameProp.procAwards(plugin, gm.Config.Controllers, curpr, gp)
}

// playgame
func (gm *GenMask) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*BasicComponentData)

	switch gm.Config.Type {
	case GMTypeSet:
		mask := gm.getFirstMask(gameProp, cd)

		if mask == nil {
			goutils.Error("GenMask.OnPlayGame:getFirstMask",
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

		gameProp.UseComponent(gm.Config.OutputMask)

		gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, mask, false)
	case GMTypeNot:
		mask := gm.getFirstMask(gameProp, cd)

		if mask == nil {
			goutils.Error("GenMask.OnPlayGame:getFirstMask",
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

		nmask := make([]bool, len(mask))
		for i := 0; i < len(mask); i++ {
			nmask[i] = !mask[i]
		}

		gameProp.UseComponent(gm.Config.OutputMask)

		gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, nmask, false)
	case GMTypeRandom:
		mask := gm.getFirstMask(gameProp, cd)
		if mask == nil {
			nmask := make([]bool, gm.Config.MaskLen)
			for i := 0; i < gm.Config.MaskLen; i++ {
				if gm.Config.WeightValue >= 10000 {
					nmask[i] = true
				} else if gm.Config.WeightValue <= 0 {
					nmask[i] = false
				} else {
					cr, err := plugin.Random(context.Background(), 10000)
					if err != nil {
						goutils.Error("GenMask.OnPlayGame:Random",
							goutils.Err(err))
						return "", err
					}

					nmask[i] = cr < gm.Config.WeightValue
				}
			}

			gameProp.UseComponent(gm.Config.OutputMask)

			gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, nmask, false)
		} else {
			if len(mask) != gm.Config.MaskLen {
				goutils.Error("GenMask.OnPlayGame:MaskLen",
					slog.Int("got", len(mask)),
					slog.Int("want", gm.Config.MaskLen),
					goutils.Err(ErrInvalidComponentConfig))
				return "", ErrInvalidComponentConfig
			}
			nmask := make([]bool, gm.Config.MaskLen)
			for i := 0; i < gm.Config.MaskLen; i++ {
				if mask[i] {
					if gm.Config.WeightValue >= 10000 {
						nmask[i] = true
					} else if gm.Config.WeightValue <= 0 {
						nmask[i] = false
					} else {
						cr, err := plugin.Random(context.Background(), 10000)
						if err != nil {
							goutils.Error("GenMask.OnPlayGame:Random",
								goutils.Err(err))
							return "", err
						}

						nmask[i] = cr < gm.Config.WeightValue
					}
				} else {
					nmask[i] = false
				}
			}

			gameProp.UseComponent(gm.Config.OutputMask)

			gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, nmask, false)
		}

	case GMTypeAnd:
		nmask := make([]bool, gm.Config.MaskLen)
		masks := gm.getAllMask(gameProp, cd)
		if masks == nil {
			return "", ErrInvalidComponentConfig
		}

		for i, curmask := range masks {
			if i == 0 {
				copy(nmask, curmask)
			} else {
				for j := 0; j < gm.Config.MaskLen; j++ {
					nmask[j] = nmask[j] && curmask[j]
				}
			}
		}

		gameProp.UseComponent(gm.Config.OutputMask)

		gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, nmask, false)
	case GMTypeOr:
		nmask := make([]bool, gm.Config.MaskLen)
		masks := gm.getAllMask(gameProp, cd)
		if masks == nil {
			return "", ErrInvalidComponentConfig
		}

		for i, curmask := range masks {
			if i == 0 {
				copy(nmask, curmask)
			} else {
				for j := 0; j < gm.Config.MaskLen; j++ {
					nmask[j] = nmask[j] || curmask[j]
				}
			}
		}

		gameProp.UseComponent(gm.Config.OutputMask)

		gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, nmask, false)
	case GMTypeXor:
		nmask := make([]bool, gm.Config.MaskLen)
		masks := gm.getAllMask(gameProp, cd)
		if masks == nil {
			return "", ErrInvalidComponentConfig
		}

		for i, curmask := range masks {
			if i == 0 {
				copy(nmask, curmask)
			} else {
				for j := 0; j < gm.Config.MaskLen; j++ {
					nmask[j] = nmask[j] != curmask[j]
				}
			}
		}

		gameProp.UseComponent(gm.Config.OutputMask)

		gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, gm.Config.OutputMask, nmask, false)
	}

	gm.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := gm.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame
func (gm *GenMask) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewGenMask(name string) IComponent {
	gm := &GenMask{
		BasicComponent: NewBasicComponent(name, 0),
	}

	return gm
}

// "type": "random",
// "maskLen": 6,
// "outputMask": "bg-mask-vs",
// "srcMask": [],
// "weightValue": 5000,
// "initMask": [
//
//	0,
//	1,
//	1,
//	1,
//	1,
//	0
//
// ]
type jsonGenMask struct {
	Type        string   `json:"type"`
	MaskLen     int      `json:"maskLen"`
	OutputMask  string   `json:"outputMask"`
	SrcMask     []string `json:"srcMask"`
	WeightValue int      `json:"weightValue"`
	InitMask    []int    `json:"initMask"`
}

func (jcfg *jsonGenMask) build() *GenMaskConfig {
	cfg := &GenMaskConfig{
		StrType:     strings.ToLower(jcfg.Type),
		MaskLen:     jcfg.MaskLen,
		OutputMask:  jcfg.OutputMask,
		SrcMask:     slices.Clone(jcfg.SrcMask),
		WeightValue: jcfg.WeightValue,
		InitMask:    make([]bool, len(jcfg.InitMask)),
	}

	for i, v := range jcfg.InitMask {
		if v != 0 {
			cfg.InitMask[i] = true
		} else {
			cfg.InitMask[i] = false
		}
	}

	return cfg
}

func parseGenMask(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenMask:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenMask:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenMask{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenMask:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		controllers, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseBasicReels:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = controllers
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenMaskTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
