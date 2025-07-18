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
	"gopkg.in/yaml.v2"
)

const CatchSymbolsTypeName = "catchSymbols"

type CatchSymbolsType int

const (
	CSTypeRule1 CatchSymbolsType = 0
)

func parseCatchSymbolsType(str string) CatchSymbolsType {
	if str == "rule1" {
		return CSTypeRule1
	}

	return CSTypeRule1
}

func mergePosWithoutSelected(srcpos []int, targetpos []int, selectindex int) []int {
	newpos := append(srcpos[0:selectindex*2], srcpos[(selectindex+1)*2:]...)
	return append(targetpos, newpos...)
}

func findNearest(sx int, sy int, targetpos []int) int {
	curlen := 9999999
	ci := -1
	var cx int
	var cy int

	for i := 0; i < len(targetpos)/2; i++ {
		tx := targetpos[i*2]
		ty := targetpos[i*2+1]

		cl := AbsInt(sx-tx) + AbsInt(sy-ty)
		if ci == -1 {
			curlen = cl
			ci = i
			cx = tx
			cy = ty
		} else if cl < curlen {
			curlen = cl
			ci = i
			cx = tx
			cy = ty
		} else if cl == curlen {
			// 如果距离一样，优先选x小的
			if tx < cx {
				curlen = cl
				ci = i
				cx = tx
				cy = ty
			} else if cx == tx {
				// 如果x也一样，优先选y小的，y不可能一样了
				if ty < cy {
					curlen = cl
					ci = i
					cx = tx
					cy = ty
				}
			}
		}
	}

	return ci
}

func moveToX(gs *sgc7game.GameScene, sx int, sy int, tx int, overrideSym int, endingSym int, ignoreSyms []int, csd *CatchSymbolsData, ignoreStart bool) {
	if sx > tx {
		for x := sx; x >= tx; x-- {
			if ignoreStart {
				ignoreStart = false

				continue
			}

			if goutils.IndexOfIntSlice(ignoreSyms, gs.Arr[x][sy], 0) < 0 {
				if x == tx {
					gs.Arr[x][sy] = endingSym
				} else {
					gs.Arr[x][sy] = overrideSym
				}

				csd.AddPos(x, sy)
			}
		}
	} else if sx < tx {
		for x := sx; x <= tx; x++ {
			if ignoreStart {
				ignoreStart = false

				continue
			}

			if goutils.IndexOfIntSlice(ignoreSyms, gs.Arr[x][sy], 0) < 0 {
				if x == tx {
					gs.Arr[x][sy] = endingSym
				} else {
					gs.Arr[x][sy] = overrideSym
				}

				csd.AddPos(x, sy)
			}
		}
	} else {
		gs.Arr[sx][sy] = endingSym

		csd.AddPos(sx, sy)
	}
}

func moveToY(gs *sgc7game.GameScene, sx int, sy int, ty int, overrideSym int, endingSym int, ignoreSyms []int, csd *CatchSymbolsData, ignoreStart bool) {
	if sy > ty {
		for y := sy; y >= ty; y-- {
			if ignoreStart {
				ignoreStart = false

				continue
			}

			if goutils.IndexOfIntSlice(ignoreSyms, gs.Arr[sx][y], 0) < 0 {
				if y == ty {
					gs.Arr[sx][y] = endingSym
				} else {
					gs.Arr[sx][y] = overrideSym
				}

				csd.AddPos(sx, y)
			}
		}
	} else if sy < ty {
		for y := sy; y <= ty; y++ {
			if ignoreStart {
				ignoreStart = false

				continue
			}

			if goutils.IndexOfIntSlice(ignoreSyms, gs.Arr[sx][y], 0) < 0 {
				if y == ty {
					gs.Arr[sx][y] = endingSym
				} else {
					gs.Arr[sx][y] = overrideSym
				}

				csd.AddPos(sx, y)
			}
		}
	} else {
		gs.Arr[sx][sy] = endingSym

		csd.AddPos(sx, sy)
	}
}

func moveTo(gs *sgc7game.GameScene, sx int, sy int, tx int, ty int, overrideSym int, endingSym int, ignoreSyms []int, csd *CatchSymbolsData, ignoreStart bool) {
	if sx != tx && sy != ty {
		moveToX(gs, sx, sy, tx, overrideSym, overrideSym, ignoreSyms, csd, ignoreStart)
		sx = tx

		moveToY(gs, sx, sy, ty, overrideSym, endingSym, ignoreSyms, csd, true)
	} else if sx != tx {
		moveToX(gs, sx, sy, tx, overrideSym, endingSym, ignoreSyms, csd, ignoreStart)
	} else if sy != ty {
		moveToY(gs, sx, sy, ty, overrideSym, endingSym, ignoreSyms, csd, ignoreStart)
	}
}

func procOneCatchAll(gs *sgc7game.GameScene, sx int, sy int, targetpos []int, overrideSym int, endingSym int, ignoreSyms []int, csd *CatchSymbolsData) {
	csd.newData()
	ignoreStart := false

	for {
		if len(targetpos) == 0 {
			break
		}

		ti := findNearest(sx, sy, targetpos)
		tx := targetpos[ti*2]
		ty := targetpos[ti*2+1]

		if len(targetpos) == 2 {
			moveTo(gs, sx, sy, tx, ty, overrideSym, endingSym, ignoreSyms, csd, ignoreStart)

			break
		}

		moveTo(gs, sx, sy, tx, ty, overrideSym, overrideSym, ignoreSyms, csd, ignoreStart)

		sx = tx
		sy = ty

		targetpos = append(targetpos[0:ti*2], targetpos[(ti+1)*2:]...)

		ignoreStart = true
	}
}

func procAllCatchOne(gs *sgc7game.GameScene, srcpos []int, tx int, ty int, overrideSym int, endingSym int, ignoreSyms []int, csd *CatchSymbolsData) {
	for i := 0; i < len(srcpos)/2; i++ {
		csd.newData()

		sx := srcpos[i*2]
		sy := srcpos[i*2+1]

		moveTo(gs, sx, sy, tx, ty, overrideSym, endingSym, ignoreSyms, csd, false)
	}
}

type CatchSymbolsData struct {
	BasicComponentData
	Pos       [][]int
	SymbolNum int
}

// OnNewGame -
func (catchSymbolsData *CatchSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	catchSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (catchSymbolsData *CatchSymbolsData) OnNewStep() {
	catchSymbolsData.UsedScenes = nil
	catchSymbolsData.Pos = nil
	catchSymbolsData.SymbolNum = 0
}

// Clone
func (catchSymbolsData *CatchSymbolsData) Clone() IComponentData {
	target := &CatchSymbolsData{
		BasicComponentData: catchSymbolsData.CloneBasicComponentData(),
		SymbolNum:          catchSymbolsData.SymbolNum,
	}

	target.Pos = make([][]int, len(catchSymbolsData.Pos))
	for i, arr := range catchSymbolsData.Pos {
		dstarr := make([]int, len(arr))
		copy(dstarr, arr)
		target.Pos[i] = dstarr
	}

	return target
}

// BuildPBComponentData
func (catchSymbolsData *CatchSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.CatchSymbolsData{
		BasicComponentData: catchSymbolsData.BuildPBBasicComponentData(),
	}

	num := 0
	for _, arr := range catchSymbolsData.Pos {
		num += len(arr)
		num++
	}

	pbcd.Pos = make([]int32, 0, num)

	for _, arr := range catchSymbolsData.Pos {
		for _, s := range arr {
			pbcd.Pos = append(pbcd.Pos, int32(s))
		}

		pbcd.Pos = append(pbcd.Pos, -1)
	}

	return pbcd
}

// GetPos -
func (catchSymbolsData *CatchSymbolsData) GetPos() []int {
	num := 0
	for _, arr := range catchSymbolsData.Pos {
		num += len(arr)
	}

	newpos := make([]int, 0, num)

	for _, arr := range catchSymbolsData.Pos {
		newpos = append(newpos, arr...)
	}

	return newpos
}

// HasPos -
func (catchSymbolsData *CatchSymbolsData) HasPos(x int, y int) bool {
	for _, arr := range catchSymbolsData.Pos {
		if goutils.IndexOfInt2Slice(arr, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

// AddPos -
func (catchSymbolsData *CatchSymbolsData) AddPos(x int, y int) {
	if len(catchSymbolsData.Pos) == 0 {
		catchSymbolsData.Pos = append(catchSymbolsData.Pos, []int{})
	}

	catchSymbolsData.Pos[len(catchSymbolsData.Pos)-1] = append(catchSymbolsData.Pos[len(catchSymbolsData.Pos)-1], x, y)
	catchSymbolsData.SymbolNum++
}

// ClearPos -
func (catchSymbolsData *CatchSymbolsData) ClearPos() {
	catchSymbolsData.Pos = nil
}

// GetValEx -
func (catchSymbolsData *CatchSymbolsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSymbolNum {
		return catchSymbolsData.SymbolNum, true
	}

	return 0, false
}

// newData -
func (catchSymbolsData *CatchSymbolsData) newData() {
	catchSymbolsData.Pos = append(catchSymbolsData.Pos, []int{})
}

// CatchSymbolsConfig - configuration for CatchSymbols
type CatchSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrCatchType         string           `yaml:"catchType" json:"catchType"`
	CatchType            CatchSymbolsType `yaml:"-" json:"-"`
	SourceSymbols        []string         `yaml:"sourceSymbols" json:"sourceSymbols"`
	SourceSymbolCodes    []int            `yaml:"-" json:"-"`
	TargetSymbols        []string         `yaml:"targetSymbols" json:"targetSymbols"`
	TargetSymbolCodes    []int            `yaml:"-" json:"-"`
	IgnoreSymbols        []string         `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes    []int            `yaml:"-" json:"-"`
	OverrideSymbol       string           `yaml:"overrideSymbol" json:"overrideSymbol"`
	OverrideSymbolCode   int              `yaml:"-" json:"-"`
	UpgradeSymbol        string           `yaml:"upgradeSymbol" json:"upgradeSymbol"`
	UpgradeSymbolCode    int              `yaml:"-" json:"-"`
	PositionCollection   string           `yaml:"positionCollection" json:"positionCollection"`
	Controllers          []*Award         `yaml:"controllers" json:"controllers"`         // 新的奖励系统
	JumpToComponent      string           `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
}

// SetLinkComponent
func (cfg *CatchSymbolsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type CatchSymbols struct {
	*BasicComponent `json:"-"`
	Config          *CatchSymbolsConfig `json:"config"`
}

// Init -
func (catchSymbols *CatchSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CatchSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CatchSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CatchSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return catchSymbols.InitEx(cfg, pool)
}

// InitEx -
func (catchSymbols *CatchSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	catchSymbols.Config = cfg.(*CatchSymbolsConfig)
	catchSymbols.Config.ComponentType = CatchSymbolsTypeName

	catchSymbols.Config.CatchType = parseCatchSymbolsType(catchSymbols.Config.StrCatchType)

	for _, s := range catchSymbols.Config.SourceSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("CatchSymbols.InitEx:SourceSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		catchSymbols.Config.SourceSymbolCodes = append(catchSymbols.Config.SourceSymbolCodes, sc)
	}

	for _, s := range catchSymbols.Config.TargetSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("CatchSymbols.InitEx:TargetSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		catchSymbols.Config.TargetSymbolCodes = append(catchSymbols.Config.TargetSymbolCodes, sc)
	}

	for _, s := range catchSymbols.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("CatchSymbols.InitEx:IgnoreSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		catchSymbols.Config.IgnoreSymbolCodes = append(catchSymbols.Config.IgnoreSymbolCodes, sc)
	}

	sc0, isok := pool.DefaultPaytables.MapSymbols[catchSymbols.Config.OverrideSymbol]
	if !isok {
		goutils.Error("CatchSymbols.InitEx:OverrideSymbol",
			slog.String("symbol", catchSymbols.Config.OverrideSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	catchSymbols.Config.OverrideSymbolCode = sc0

	sc1, isok := pool.DefaultPaytables.MapSymbols[catchSymbols.Config.UpgradeSymbol]
	if !isok {
		goutils.Error("CatchSymbols.InitEx:UpgradeSymbol",
			slog.String("symbol", catchSymbols.Config.UpgradeSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	catchSymbols.Config.UpgradeSymbolCode = sc1

	for _, ctrl := range catchSymbols.Config.Controllers {
		ctrl.Init()
	}

	catchSymbols.onInit(&catchSymbols.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (catchSymbols *CatchSymbols) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(catchSymbols.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, catchSymbols.Config.Controllers, curpr, gp)
	}
}

// playgame
func (catchSymbols *CatchSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// moveSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	csd := cd.(*CatchSymbolsData)

	csd.OnNewStep()

	gs := catchSymbols.GetTargetScene3(gameProp, curpr, prs, 0)

	srcpos := []int{}
	targetpos := []int{}

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if goutils.IndexOfIntSlice(catchSymbols.Config.SourceSymbolCodes, s, 0) >= 0 {
				srcpos = append(srcpos, x, y)
			} else if goutils.IndexOfIntSlice(catchSymbols.Config.TargetSymbolCodes, s, 0) >= 0 {
				targetpos = append(targetpos, x, y)
			}
		}
	}

	if len(srcpos) == 0 || len(targetpos) == 0 {
		nc := catchSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	sc2 := gs.CloneEx(gameProp.PoolScene)

	if catchSymbols.Config.CatchType == CSTypeRule1 {
		si := -1

		if len(srcpos) == 2 {
			si = 0
		} else {
			cr, err := plugin.Random(context.TODO(), len(srcpos)/2)
			if err != nil {
				goutils.Error("CatchSymbols.OnPlayGame:CSTypeRule1:Random",
					goutils.Err(ErrInvalidSymbol))

				return "", ErrInvalidSymbol
			}

			si = cr
		}

		if catchSymbols.Config.PositionCollection != "" {
			gameProp.AddComponentPos(catchSymbols.Config.PositionCollection, targetpos)
		}

		if len(targetpos) > 2 {
			sx := srcpos[si*2]
			sy := srcpos[si*2+1]

			if len(srcpos) > 2 {
				procOneCatchAll(sc2, sx, sy, mergePosWithoutSelected(srcpos, targetpos, si), catchSymbols.Config.OverrideSymbolCode, catchSymbols.Config.UpgradeSymbolCode, catchSymbols.Config.IgnoreSymbolCodes, csd)
			} else {
				procOneCatchAll(sc2, sx, sy, targetpos, catchSymbols.Config.OverrideSymbolCode, gs.Arr[sx][sy], catchSymbols.Config.IgnoreSymbolCodes, csd)
			}
		} else if len(targetpos) == 2 {
			if len(srcpos) > 2 {
				procAllCatchOne(sc2, srcpos, targetpos[0], targetpos[1], catchSymbols.Config.OverrideSymbolCode, catchSymbols.Config.UpgradeSymbolCode, catchSymbols.Config.IgnoreSymbolCodes, csd)
			} else {
				sx := srcpos[si*2]
				sy := srcpos[si*2+1]

				procOneCatchAll(sc2, sx, sy, targetpos, catchSymbols.Config.OverrideSymbolCode, gs.Arr[sx][sy], catchSymbols.Config.IgnoreSymbolCodes, csd)
			}
		}
	}

	catchSymbols.AddScene(gameProp, curpr, sc2, &csd.BasicComponentData)

	catchSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")
	// if len(catchSymbols.Config.Controllers) > 0 {
	// 	gameProp.procAwards(plugin, catchSymbols.Config.Controllers, curpr, gp)
	// }

	nc := catchSymbols.onStepEnd(gameProp, curpr, gp, catchSymbols.Config.JumpToComponent)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (catchSymbols *CatchSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	csd := cd.(*CatchSymbolsData)

	asciigame.OutputScene("after catchSymbols", pr.Scenes[csd.UsedScenes[0]], mapSymbolColor)

	return nil
}

// NewComponentData -
func (catchSymbols *CatchSymbols) NewComponentData() IComponentData {
	return &CatchSymbolsData{}
}

func NewCatchSymbols(name string) IComponent {
	return &CatchSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "catchType": "normal",
// "sourceSymbols": [
//
//	"RW",
//	"RW2"
//
// ],
// "targetSymbols": [
//
//	"MM"
//
// ],
// "ignoreSymbols": [
//
//	"SC"
//
// ],
// "overrideSymbol": "MY",
// "upgradeSymbol": "RW2",
// "positionCollection": "bg-burstpos"
type jsonCatchSymbols struct {
	CatchType          string   `json:"catchType"`
	SourceSymbols      []string `json:"sourceSymbols"`
	TargetSymbols      []string `json:"targetSymbols"`
	IgnoreSymbols      []string `json:"ignoreSymbols"`
	OverrideSymbol     string   `json:"overrideSymbol"`
	UpgradeSymbol      string   `json:"upgradeSymbol"`
	PositionCollection string   `json:"positionCollection"`
}

func (jcfg *jsonCatchSymbols) build() *CatchSymbolsConfig {
	cfg := &CatchSymbolsConfig{
		StrCatchType:       jcfg.CatchType,
		SourceSymbols:      jcfg.SourceSymbols,
		TargetSymbols:      jcfg.TargetSymbols,
		IgnoreSymbols:      jcfg.IgnoreSymbols,
		OverrideSymbol:     jcfg.OverrideSymbol,
		UpgradeSymbol:      jcfg.UpgradeSymbol,
		PositionCollection: jcfg.PositionCollection,
	}

	// for _, v := range jms.MoveData {
	// 	cmd := &MoveData{
	// 		Src:            v.Src,
	// 		Target:         v.Target,
	// 		MoveType:       v.MoveType,
	// 		TargetSymbol:   v.TargetSymbol,
	// 		OverrideSrc:    v.OverrideSrc == "true",
	// 		OverrideTarget: v.OverrideTarget == "true",
	// 		OverridePath:   v.OverridePath == "true",
	// 	}

	// 	if cmd.Src.X > 0 {
	// 		cmd.Src.X--
	// 	}

	// 	if cmd.Src.Y > 0 {
	// 		cmd.Src.Y--
	// 	}

	// 	if cmd.Target.X > 0 {
	// 		cmd.Target.X--
	// 	}

	// 	if cmd.Target.Y > 0 {
	// 		cmd.Target.Y--
	// 	}

	// 	cfg.MoveData = append(cfg.MoveData, cmd)
	// }

	// cfg.UseSceneV3 = true

	return cfg
}

func parseCatchSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCatchSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCatchSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCatchSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCatchSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseCatchSymbols:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CatchSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
