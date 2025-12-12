package lowcode

import (
	"context"
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
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const RefillSymbols2TypeName = "refillSymbols2"

type RefillSymbols2Data struct {
	BasicComponentData
	cfg *RefillSymbols2Config
}

// OnNewGame -
func (rs2d *RefillSymbols2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	rs2d.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (rs2d *RefillSymbols2Data) onNewStep() {
	rs2d.UsedScenes = nil
	rs2d.UsedOtherScenes = nil
	rs2d.Pos = nil
}

// Clone
func (rs2d *RefillSymbols2Data) Clone() IComponentData {
	target := &RefillSymbols2Data{
		BasicComponentData: rs2d.CloneBasicComponentData(),
		cfg:                rs2d.cfg,
	}

	return target
}

// BuildPBComponentData
func (rs2d *RefillSymbols2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BasicComponentData{
		BasicComponentData: rs2d.BuildPBBasicComponentData(),
	}

	return pbcd
}

// ChgConfigIntVal -
func (rs2d *RefillSymbols2Data) ChgConfigIntVal(key string, off int) int {
	if key == CCVHeight {
		_, isok := rs2d.MapConfigIntVals[key]
		if !isok {
			if rs2d.cfg.Height > 0 {
				rs2d.MapConfigIntVals[key] = rs2d.cfg.Height
			}
		}
	}

	return rs2d.BasicComponentData.ChgConfigIntVal(key, off)
}

// GetValEx -
func (rs2d *RefillSymbols2Data) GetValEx(key string, getType GetComponentValType) (int, bool) {
	switch key {
	case CVHeight:
		return rs2d.getHeight(), true
	}

	return 0, false
}

func (rs2d *RefillSymbols2Data) getHeight() int {
	height, isok := rs2d.GetConfigIntVal(CCVHeight)
	if isok {
		return height
	}

	return rs2d.cfg.Height
}

// RefillSymbols2Config - configuration for RefillSymbols2
type RefillSymbols2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	IsNeedProcSymbolVals bool   `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int    `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	DefaultSymbolVal     int    `yaml:"defaultSymbolVal" json:"defaultSymbolVal"`         // 重新填充的symbolVal是什么
	OutputToComponent    string `yaml:"outputToComponent" json:"outputToComponent"`       // 输出到哪个组件
	Height               int    `yaml:"height" json:"height"`                             // 重新填充的symbolVal是什么
	MaskX                string `yaml:"maskX" json:"maskX"`                               // maskX
	MaskY                string `yaml:"maskY" json:"maskY"`                               // maskY
	CollectorPay         string `yaml:"collectorPay" json:"collectorPay"`                 // collectorPay
	CPCore               string `yaml:"cpCore" json:"cpCore"`                             // cpCore
}

// SetLinkComponent
func (cfg *RefillSymbols2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RefillSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *RefillSymbols2Config `json:"config"`
}

// Init -
func (refillSymbols2 *RefillSymbols2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RefillSymbols2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RefillSymbols2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RefillSymbols2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return refillSymbols2.InitEx(cfg, pool)
}

// InitEx -
func (refillSymbols2 *RefillSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	refillSymbols2.Config = cfg.(*RefillSymbols2Config)
	refillSymbols2.Config.ComponentType = RefillSymbols2TypeName

	refillSymbols2.onInit(&refillSymbols2.Config.BasicComponentConfig)

	return nil
}

func (refillSymbols2 *RefillSymbols2) getSymbol(rd *sgc7game.ReelsData, x int, index int) int {
	index--

	for ; index < 0; index += len(rd.Reels[x]) {
	}

	return rd.Reels[x][index]
}

func (refillSymbols2 *RefillSymbols2) getHeight(basicCD *BasicComponentData) int {
	height, isok := basicCD.GetConfigIntVal(CCVHeight)
	if isok {
		return height
	}

	return refillSymbols2.Config.Height
}

func (refillSymbols2 *RefillSymbols2) getMaskX(_ *GameProperty, basicCD *BasicComponentData) string {
	str := basicCD.GetConfigVal(CCVMaskX)
	if str != "" {
		return str
	}

	return refillSymbols2.Config.MaskX
}

func (refillSymbols2 *RefillSymbols2) refillHeightAndMaskX(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, os *sgc7game.GameScene,
	height int, maskX string, outputCD IComponentData, rcd *RefillSymbols2Data) (*sgc7game.GameScene, *sgc7game.GameScene, error) {

	var maskArr []bool
	if maskX == "<empty>" {
		maskArr = make([]bool, 0, gs.Width)

		for i := 0; i < gs.Width; i++ {
			maskArr = append(maskArr, true)
		}
	} else {
		imaskd := gameProp.GetComponentDataWithName(maskX)
		if imaskd != nil {
			maskArr = imaskd.GetMask()
			if len(maskArr) != gs.Width {
				goutils.Error("RefillSymbols2.refillHeightAndMaskX:MaskX:len(arr)!=gs.Width",
					goutils.Err(ErrInvalidComponentConfig))

				return nil, nil, ErrInvalidComponentConfig
			}
		} else {
			goutils.Error("RefillSymbols2.refillHeightAndMaskX:MaskX",
				slog.String("maskX", maskX),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, nil, ErrInvalidComponentConfig
		}
	}

	ngs := gs
	cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

	if os != nil {
		nos := os

		for x := 0; x < gs.Width; x++ {
			if !maskArr[x] {
				continue
			} else {
				if ngs.Indexes[x] < 0 {
					ci, err := plugin.Random(context.Background(), len(cr.Reels[x]))
					if err != nil {
						goutils.Error("RefillSymbols2.refillHeightAndMaskX:Random",
							slog.Int("len", len(cr.Reels[x])),
							goutils.Err(err))

						return nil, nil, err
					}

					ngs.Indexes[x] = ci
				}
			}

			for y := gs.Height - 1; y >= 0; y-- {
				if y < gs.Height-height {
					continue
				}

				if ngs.Arr[x][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
						nos = os.CloneEx(gameProp.PoolScene)
					}

					ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
					ngs.Indexes[x]--

					nos.Arr[x][y] = refillSymbols2.Config.DefaultSymbolVal

					if outputCD != nil {
						outputCD.AddPos(x, y)
					}

					rcd.AddPos(x, y)
				}
			}
		}

		return ngs, nos, nil
	}

	for x := 0; x < gs.Width; x++ {
		if !maskArr[x] {
			continue
		} else {
			if ngs.Indexes[x] < 0 {
				ci, err := plugin.Random(context.Background(), len(cr.Reels[x]))
				if err != nil {
					goutils.Error("RefillSymbols2.refillHeightAndMaskX:Random",
						slog.Int("len", len(cr.Reels[x])),
						goutils.Err(err))

					return nil, nil, err
				}

				ngs.Indexes[x] = ci
			}
		}

		for y := gs.Height - 1; y >= 0; y-- {
			if y < gs.Height-height {
				continue
			}

			if ngs.Arr[x][y] == -1 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
				ngs.Indexes[x]--

				if outputCD != nil {
					outputCD.AddPos(x, y)
				}

				rcd.AddPos(x, y)
			}
		}
	}

	return ngs, nil, nil
}

func (refillSymbols2 *RefillSymbols2) refillMaskX(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, os *sgc7game.GameScene,
	maskX string, outputCD IComponentData, rcd *RefillSymbols2Data) (*sgc7game.GameScene, *sgc7game.GameScene, error) {

	var maskArr []bool
	if maskX == "<empty>" {
		maskArr = make([]bool, 0, gs.Width)

		for i := 0; i < gs.Width; i++ {
			maskArr = append(maskArr, true)
		}
	} else {
		imaskd := gameProp.GetComponentDataWithName(maskX)
		if imaskd != nil {
			maskArr = imaskd.GetMask()
			if len(maskArr) != gs.Width {
				goutils.Error("RefillSymbols2.refillMaskX:MaskX:len(arr)!=gs.Width",
					goutils.Err(ErrInvalidComponentConfig))

				return nil, nil, ErrInvalidComponentConfig
			}
		} else {
			goutils.Error("RefillSymbols2.refillMaskX:MaskX",
				slog.String("maskX", maskX),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, nil, ErrInvalidComponentConfig
		}
	}

	ngs := gs
	cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

	if os != nil {
		nos := os

		for x := 0; x < gs.Width; x++ {
			if !maskArr[x] {
				continue
			} else {
				if ngs.Indexes[x] < 0 {
					ci, err := plugin.Random(context.Background(), len(cr.Reels[x]))
					if err != nil {
						goutils.Error("RefillSymbols2.refillMaskX:Random",
							slog.Int("len", len(cr.Reels[x])),
							goutils.Err(err))

						return nil, nil, err
					}

					ngs.Indexes[x] = ci
				}
			}

			for y := gs.Height - 1; y >= 0; y-- {
				if ngs.Arr[x][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
						nos = os.CloneEx(gameProp.PoolScene)
					}

					ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
					ngs.Indexes[x]--

					nos.Arr[x][y] = refillSymbols2.Config.DefaultSymbolVal

					if outputCD != nil {
						outputCD.AddPos(x, y)
					}

					rcd.AddPos(x, y)
				}
			}
		}

		return ngs, nos, nil
	}

	for x := 0; x < gs.Width; x++ {
		if !maskArr[x] {
			continue
		} else {
			if ngs.Indexes[x] < 0 {
				ci, err := plugin.Random(context.Background(), len(cr.Reels[x]))
				if err != nil {
					goutils.Error("RefillSymbols2.refillMaskX:Random",
						slog.Int("len", len(cr.Reels[x])),
						goutils.Err(err))

					return nil, nil, err
				}

				ngs.Indexes[x] = ci
			}
		}

		for y := gs.Height - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
				ngs.Indexes[x]--

				if outputCD != nil {
					outputCD.AddPos(x, y)
				}

				rcd.AddPos(x, y)
			}
		}
	}

	return ngs, nil, nil
}

func (refillSymbols2 *RefillSymbols2) refillOnlyHeight(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene,
	height int, outputCD IComponentData, rcd *RefillSymbols2Data) (*sgc7game.GameScene, *sgc7game.GameScene) {

	ngs := gs

	if os != nil {
		nos := os

		for x := 0; x < gs.Width; x++ {
			for y := gs.Height - 1; y >= 0; y-- {
				if y < gs.Height-height {
					continue
				}

				if ngs.Arr[x][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
						nos = os.CloneEx(gameProp.PoolScene)
					}

					cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

					ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
					ngs.Indexes[x]--

					nos.Arr[x][y] = refillSymbols2.Config.DefaultSymbolVal

					if outputCD != nil {
						outputCD.AddPos(x, y)
					}

					rcd.AddPos(x, y)
				}
			}
		}

		return ngs, nos
	}

	for x := 0; x < gs.Width; x++ {
		for y := gs.Height - 1; y >= 0; y-- {
			if y < gs.Height-height {
				continue
			}

			if ngs.Arr[x][y] == -1 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

				ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
				ngs.Indexes[x]--

				if outputCD != nil {
					outputCD.AddPos(x, y)
				}

				rcd.AddPos(x, y)
			}
		}
	}

	return ngs, nil
}

func (refillSymbols2 *RefillSymbols2) refill(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene,
	outputCD IComponentData, rcd *RefillSymbols2Data) (*sgc7game.GameScene, *sgc7game.GameScene) {

	ngs := gs

	if os != nil {
		nos := os

		for x := 0; x < gs.Width; x++ {
			for y := gs.Height - 1; y >= 0; y-- {
				if ngs.Arr[x][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
						nos = os.CloneEx(gameProp.PoolScene)
					}

					cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

					ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
					ngs.Indexes[x]--

					nos.Arr[x][y] = refillSymbols2.Config.DefaultSymbolVal

					if outputCD != nil {
						outputCD.AddPos(x, y)
					}

					rcd.AddPos(x, y)
				}
			}
		}

		return ngs, nos
	}

	for x := 0; x < gs.Width; x++ {
		for y := gs.Height - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

				ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
				ngs.Indexes[x]--

				if outputCD != nil {
					outputCD.AddPos(x, y)
				}
			}
		}
	}

	return ngs, nil
}

func (refillSymbols2 *RefillSymbols2) procCollectorPay(gameProp *GameProperty, gs *sgc7game.GameScene) {
	if refillSymbols2.Config.CollectorPay == "" {
		return
	}

	icd := gameProp.GetComponentDataWithName(refillSymbols2.Config.CollectorPay)
	if icd == nil {
		goutils.Error("RefillSymbols2.procCollectorPay:CollectorPay",
			slog.String("collectorPay", refillSymbols2.Config.CollectorPay),
			goutils.Err(ErrInvalidComponent))

		return
	}

	cptd, isok := icd.(*CollectorPayTriggerData)
	if !isok {
		goutils.Error("RefillSymbols2.procCollectorPay:CollectorPay",
			slog.String("collectorPay", refillSymbols2.Config.CollectorPay),
			goutils.Err(ErrInvalidComponentData))

		return
	}

	cptd.procSymbolsWithLevel(gs)
}

func (refillSymbols2 *RefillSymbols2) getCPCoreData(gameProp *GameProperty) (*CPCoreData, error) {
	icd := gameProp.GetComponentDataWithName(refillSymbols2.Config.CPCore)
	if icd == nil {
		goutils.Error("RefillSymbols2.getCPCoreData:GetComponentDataWithName",
			slog.String("CPCore", refillSymbols2.Config.CPCore),
			goutils.Err(ErrInvalidComponentData))

		return nil, ErrInvalidComponentData
	}

	cd, isok := icd.(*CPCoreData)
	if !isok {
		goutils.Error("RefillSymbols2.getCPCoreData:ComponentDataType",
			slog.String("CPCore", refillSymbols2.Config.CPCore),
			goutils.Err(ErrInvalidComponentData))

		return nil, ErrInvalidComponentData
	}

	return cd, nil
}

func (refillSymbols2 *RefillSymbols2) onPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, _ *GameParams, plugin sgc7plugin.IPlugin,
	_ string, _ string, _ sgc7game.IPlayerState, _ *sgc7game.Stake, prs []*sgc7game.PlayResult, rcd *RefillSymbols2Data) (
	*sgc7game.GameScene, *sgc7game.GameScene, bool, error) {

	gs := refillSymbols2.GetTargetScene3(gameProp, curpr, prs, 0)

	var os *sgc7game.GameScene
	if refillSymbols2.Config.IsNeedProcSymbolVals {
		os = refillSymbols2.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	var outputCD IComponentData
	if refillSymbols2.Config.OutputToComponent != "" {
		outputCD = gameProp.GetComponentDataWithName(refillSymbols2.Config.OutputToComponent)
		if outputCD == nil {
			goutils.Error("RefillSymbols2.OnPlayGame:OutputToComponent",
				slog.String("outputToComponent", refillSymbols2.Config.OutputToComponent),
				goutils.Err(ErrInvalidComponent))

			return nil, nil, false, ErrInvalidComponent
		}

		outputCD.ClearPos()
	}

	height := refillSymbols2.getHeight(&rcd.BasicComponentData)
	maskX := refillSymbols2.getMaskX(gameProp, &rcd.BasicComponentData)

	if maskX != "" {
		if height > 0 && height < gs.Height {
			ngs, nos, err := refillSymbols2.refillHeightAndMaskX(gameProp, plugin, gs, os, height, maskX, outputCD, rcd)
			if err != nil {
				goutils.Error("RefillSymbols2.OnPlayGame:refillHeightAndMaskX",
					goutils.Err(err))

				return nil, nil, false, err
			}

			if ngs == gs {
				return nil, nil, false, nil
			}

			refillSymbols2.procCollectorPay(gameProp, ngs)

			refillSymbols2.AddScene(gameProp, curpr, ngs, &rcd.BasicComponentData)
			if nos != nil {
				refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rcd.BasicComponentData)
			}

			return ngs, nos, true, nil
		}

		ngs, nos, err := refillSymbols2.refillMaskX(gameProp, plugin, gs, os, maskX, outputCD, rcd)
		if err != nil {
			goutils.Error("RefillSymbols2.OnPlayGame:refillMaskX",
				goutils.Err(err))

			return nil, nil, false, err
		}

		if ngs == gs {
			return nil, nil, false, nil
		}

		if nos != nil {
			refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rcd.BasicComponentData)
		}

		refillSymbols2.procCollectorPay(gameProp, ngs)

		refillSymbols2.AddScene(gameProp, curpr, ngs, &rcd.BasicComponentData)

		return ngs, nos, true, nil
	}

	if height > 0 && height < gs.Height {
		ngs, nos := refillSymbols2.refillOnlyHeight(gameProp, gs, os, height, outputCD, rcd)
		if ngs == gs {
			return nil, nil, false, nil
		}

		refillSymbols2.procCollectorPay(gameProp, ngs)

		refillSymbols2.AddScene(gameProp, curpr, ngs, &rcd.BasicComponentData)
		if nos != nil {
			refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rcd.BasicComponentData)
		}

		return ngs, nos, true, nil
	}

	ngs, nos := refillSymbols2.refill(gameProp, gs, os, outputCD, rcd)
	if ngs == gs {
		return nil, nil, false, nil
	}

	if nos != nil {
		refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rcd.BasicComponentData)
	}

	refillSymbols2.procCollectorPay(gameProp, ngs)

	refillSymbols2.AddScene(gameProp, curpr, ngs, &rcd.BasicComponentData)

	return ngs, nos, true, nil
}

// playgame
func (refillSymbols2 *RefillSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	rs2d := cd.(*RefillSymbols2Data)

	rs2d.onNewStep()

	if refillSymbols2.Config.CPCore != "" {
		ngs, nos, isTriggered, err := refillSymbols2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, rs2d)
		if err != nil {
			goutils.Error("RefillSymbols2.OnPlayGame:onPlayGame",
				goutils.Err(err))

			return "", err
		}

		if isTriggered {
			cpcorecd, err := refillSymbols2.getCPCoreData(gameProp)
			if err != nil {
				goutils.Error("RefillSymbols2.OnPlayGame:getCPCoreData",
					goutils.Err(err))

				return "", err
			}

			err = cpcorecd.onRefillSymbols2(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, rs2d, refillSymbols2, ngs, nos)
			if err != nil {
				goutils.Error("RefillSymbols2.OnPlayGame:onRefillSymbols2",
					goutils.Err(err))

				return "", err
			}

			nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

			return nc, nil
		}

		nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	_, _, isTriggered, err := refillSymbols2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, rs2d)
	if err != nil {
		goutils.Error("RefillSymbols2.OnPlayGame:onPlayGame",
			goutils.Err(err))

		return "", err
	}

	if !isTriggered {
		nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil

	// gs := refillSymbols2.GetTargetScene3(gameProp, curpr, prs, 0)

	// var os *sgc7game.GameScene
	// if refillSymbols2.Config.IsNeedProcSymbolVals {
	// 	os = refillSymbols2.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	// }

	// var outputCD IComponentData
	// if refillSymbols2.Config.OutputToComponent != "" {
	// 	outputCD = gameProp.GetComponentDataWithName(refillSymbols2.Config.OutputToComponent)
	// 	if outputCD == nil {
	// 		goutils.Error("RefillSymbols2.OnPlayGame:OutputToComponent",
	// 			slog.String("outputToComponent", refillSymbols2.Config.OutputToComponent),
	// 			goutils.Err(ErrInvalidComponent))

	// 		return "", ErrInvalidComponent
	// 	}

	// 	outputCD.ClearPos()
	// }

	// height := refillSymbols2.getHeight(&rs2d.BasicComponentData)
	// maskX := refillSymbols2.getMaskX(gameProp, &rs2d.BasicComponentData)

	// if maskX != "" {
	// 	if height > 0 && height < gs.Height {
	// 		ngs, nos, err := refillSymbols2.refillHeightAndMaskX(gameProp, plugin, gs, os, height, maskX, outputCD)
	// 		if err != nil {
	// 			goutils.Error("RefillSymbols2.OnPlayGame:refillHeightAndMaskX",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		if ngs == gs {
	// 			nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 			return nc, ErrComponentDoNothing
	// 		}

	// 		refillSymbols2.procCollectorPay(gameProp, ngs)

	// 		refillSymbols2.AddScene(gameProp, curpr, ngs, &rs2d.BasicComponentData)
	// 		if nos != nil {
	// 			refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rs2d.BasicComponentData)
	// 		}

	// 		nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 		return nc, nil
	// 	}

	// 	ngs, nos, err := refillSymbols2.refillMaskX(gameProp, plugin, gs, os, maskX, outputCD)
	// 	if err != nil {
	// 		goutils.Error("RefillSymbols2.OnPlayGame:refillMaskX",
	// 			goutils.Err(err))

	// 		return "", err
	// 	}

	// 	if ngs == gs {
	// 		nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 		return nc, ErrComponentDoNothing
	// 	}

	// 	if nos != nil {
	// 		refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rs2d.BasicComponentData)
	// 	}

	// 	refillSymbols2.procCollectorPay(gameProp, ngs)

	// 	refillSymbols2.AddScene(gameProp, curpr, ngs, &rs2d.BasicComponentData)

	// 	nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, nil
	// }

	// if height > 0 && height < gs.Height {
	// 	ngs, nos := refillSymbols2.refillOnlyHeight(gameProp, gs, os, height, outputCD)
	// 	if ngs == gs {
	// 		nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 		return nc, ErrComponentDoNothing
	// 	}

	// 	refillSymbols2.procCollectorPay(gameProp, ngs)

	// 	refillSymbols2.AddScene(gameProp, curpr, ngs, &rs2d.BasicComponentData)
	// 	if nos != nil {
	// 		refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rs2d.BasicComponentData)
	// 	}

	// 	nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, nil
	// }

	// ngs, nos := refillSymbols2.refill(gameProp, gs, os, outputCD)
	// if ngs == gs {
	// 	nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, ErrComponentDoNothing
	// }

	// if nos != nil {
	// 	refillSymbols2.AddOtherScene(gameProp, curpr, nos, &rs2d.BasicComponentData)
	// }

	// refillSymbols2.procCollectorPay(gameProp, ngs)

	// refillSymbols2.AddScene(gameProp, curpr, ngs, &rs2d.BasicComponentData)

	// nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	// return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (refillSymbols2 *RefillSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after refillSymbols2", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// EachUsedResults -
func (refillSymbols2 *RefillSymbols2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

// NewComponentData -
func (refillSymbols2 *RefillSymbols2) NewComponentData() IComponentData {
	return &RefillSymbols2Data{
		cfg: refillSymbols2.Config,
	}
}

func NewRefillSymbols2(name string) IComponent {
	return &RefillSymbols2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "isNeedProcSymbolVals": false,
// "emptySymbolVal": -1,
// "defaultSymbolVal": 0,
// "Height": 6,
// "maskX": "mask-6"
// "collectorPay": "bg-pay"
// "cpCore": "bg-core"
type jsonRefillSymbols2 struct {
	IsNeedProcSymbolVals bool   `json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int    `json:"emptySymbolVal"`       // 空的symbolVal是什么
	DefaultSymbolVal     int    `json:"defaultSymbolVal"`     // 重新填充的symbolVal是什么
	OutputToComponent    string `json:"outputToComponent"`    // 输出到哪个组件
	Height               int    `json:"Height"`               // height, <=0 is ignore
	MaskX                string `json:"maskX"`                // maskX
	MaskY                string `json:"maskY"`                // maskY
	CollectorPay         string `json:"collectorPay"`         // collectorPay
	CPCore               string `json:"cpCore"`               // cpCore
}

func (jcfg *jsonRefillSymbols2) build() *RefillSymbols2Config {
	cfg := &RefillSymbols2Config{
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
		DefaultSymbolVal:     jcfg.DefaultSymbolVal,
		OutputToComponent:    jcfg.OutputToComponent,
		Height:               jcfg.Height,
		MaskX:                jcfg.MaskX,
		MaskY:                jcfg.MaskY,
		CollectorPay:         jcfg.CollectorPay,
	}

	return cfg
}

func parseRefillSymbols2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRefillSymbols2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRefillSymbols2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRefillSymbols2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRefillSymbols2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RefillSymbols2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
