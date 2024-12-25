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

const BombTypeName = "bomb"

type BombSourceType int

const (
	BSTypePositionCollection BombSourceType = 0
	BSTypeSymbols            BombSourceType = 1
)

func parseBombSourceType(str string) BombSourceType {
	if str == "positioncollection" {
		return BSTypePositionCollection
	}

	return BSTypeSymbols
}

type BombTargetType int

const (
	BTTypeSymbols BombTargetType = 0
	BTTypeRemove  BombTargetType = 1
)

func parseBombTargetType(str string) BombTargetType {
	if str == "remove" {
		return BTTypeRemove
	}

	return BTTypeSymbols
}

type BombData struct {
	BasicComponentData
	Pos [][]int
}

// OnNewGame -
func (bombData *BombData) OnNewGame(gameProp *GameProperty, component IComponent) {
	bombData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (bombData *BombData) onNewStep() {
	bombData.UsedScenes = nil
	bombData.Pos = nil
}

// Clone
func (bombData *BombData) Clone() IComponentData {
	target := &BombData{
		BasicComponentData: bombData.CloneBasicComponentData(),
	}

	target.Pos = make([][]int, len(bombData.Pos))
	for _, arr := range bombData.Pos {
		dstarr := make([]int, len(arr))
		copy(dstarr, arr)
		target.Pos = append(target.Pos, dstarr)
	}

	return target
}

// BuildPBComponentData
func (bombData *BombData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BombData{
		BasicComponentData: bombData.BuildPBBasicComponentData(),
	}

	num := 0
	for _, arr := range bombData.Pos {
		num += len(arr)
		num++
	}

	pbcd.Pos = make([]int32, 0, num)

	for _, arr := range bombData.Pos {
		for _, s := range arr {
			pbcd.Pos = append(pbcd.Pos, int32(s))
		}

		pbcd.Pos = append(pbcd.Pos, -1)
	}

	return pbcd
}

// GetPos -
func (bombData *BombData) GetPos() []int {
	num := 0
	for _, arr := range bombData.Pos {
		num += len(arr)
	}

	newpos := make([]int, 0, num)

	for _, arr := range bombData.Pos {
		newpos = append(newpos, arr...)
	}

	return newpos
}

// HasPos -
func (bombData *BombData) HasPos(x int, y int) bool {
	for _, arr := range bombData.Pos {
		if goutils.IndexOfInt2Slice(arr, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

// AddPos -
func (bombData *BombData) AddPos(x int, y int) {
	if len(bombData.Pos) == 0 {
		bombData.Pos = append(bombData.Pos, []int{})
	}

	bombData.Pos[len(bombData.Pos)-1] = append(bombData.Pos[len(bombData.Pos)-1], x, y)
}

// newData -
func (bombData *BombData) newData() {
	bombData.Pos = append(bombData.Pos, []int{})
}

// BombConfig - configuration for Bomb
type BombConfig struct {
	BasicComponentConfig     `yaml:",inline" json:",inline"`
	BombWidth                int            `json:"bombWidth"`
	BombHeight               int            `json:"bombHeight"`
	StrBombSourceType        string         `json:"bombSourceType"`
	BombSourceType           BombSourceType `json:"-"`
	SourceSymbols            []string       `json:"sourceSymbols"`
	SourceSymbolCodes        []int          `json:"-"`
	BombData                 [][]int        `json:"bombData"`
	SourcePositionCollection string         `json:"sourcePositionCollection"`
	SelectSourceNumber       int            `json:"selectSourceNumber"`
	IgnoreSymbols            []string       `json:"ignoreSymbols"`
	IgnoreSymbolCodes        []int          `json:"-"`
	StrBombTargetType        string         `json:"bombTargetType"`
	BombTargetType           BombTargetType `json:"-"`
	TargetSymbol             string         `json:"targetSymbol"`
	TargetSymbolCode         int            `json:"-"`
	OutputToComponent        string         `json:"outputToComponent"`
	IgnoreWinResults         []string       `json:"ignoreWinResults"`
	IgnorePostionCollections []string       `json:"ignorePostionCollections"`
	Controllers              []*Award       `yaml:"controllers" json:"controllers"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *BombConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type Bomb struct {
	*BasicComponent `json:"-"`
	Config          *BombConfig `json:"config"`
}

// Init -
func (bomb *Bomb) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Bomb.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &BombConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Bomb.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return bomb.InitEx(cfg, pool)
}

// InitEx -
func (bomb *Bomb) InitEx(cfg any, pool *GamePropertyPool) error {
	bomb.Config = cfg.(*BombConfig)
	bomb.Config.ComponentType = BombTypeName

	bomb.Config.BombSourceType = parseBombSourceType(bomb.Config.StrBombSourceType)
	bomb.Config.BombTargetType = parseBombTargetType(bomb.Config.StrBombTargetType)

	for _, s := range bomb.Config.SourceSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("Bomb.InitEx:SourceSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		bomb.Config.SourceSymbolCodes = append(bomb.Config.SourceSymbolCodes, sc)
	}

	for _, s := range bomb.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("Bomb.InitEx:IgnoreSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		bomb.Config.IgnoreSymbolCodes = append(bomb.Config.IgnoreSymbolCodes, sc)
	}

	if bomb.Config.TargetSymbol != "" {
		sc, isok := pool.DefaultPaytables.MapSymbols[bomb.Config.TargetSymbol]
		if !isok {
			goutils.Error("Bomb.InitEx:IgnoreSymbols.TargetSymbol",
				slog.String("symbol", bomb.Config.TargetSymbol),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		bomb.Config.TargetSymbolCode = sc
	}

	for _, ctrl := range bomb.Config.Controllers {
		ctrl.Init()
	}

	bomb.onInit(&bomb.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (bomb *Bomb) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(bomb.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, bomb.Config.Controllers, curpr, gp)
	}
}

func (bomb *Bomb) canBomb(gameProp *GameProperty, x int, y int, curpr *sgc7game.PlayResult, gp *GameParams) bool {
	for _, v := range bomb.Config.IgnorePostionCollections {
		pos := gameProp.GetComponentPos(v)
		if goutils.IndexOfInt2Slice(pos, x, y, 0) >= 0 {
			return false
		}
	}

	for _, v := range bomb.Config.IgnoreWinResults {
		// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
		if goutils.IndexOfStringSlice(gp.HistoryComponents, v, 0) < 0 {
			continue
		}

		ccd := gameProp.GetCurComponentDataWithName(v)
		if ccd != nil {
			lst := ccd.GetResults()
			for _, ri := range lst {

				if goutils.IndexOfInt2Slice(curpr.Results[ri].Pos, x, y, 0) >= 0 {
					return false
				}
			}
		}
	}

	return true
}

func (bomb *Bomb) bomb(gameProp *GameProperty, gs *sgc7game.GameScene, x int, y int, curpr *sgc7game.PlayResult, gp *GameParams, bsd *BombData) {
	cx := bomb.Config.BombWidth / 2
	cy := bomb.Config.BombHeight / 2

	bsd.newData()

	for sx, arr := range bomb.Config.BombData {
		for sy, v := range arr {
			tx := x + sx - cx
			ty := y + sy - cy
			if v != 0 && tx >= 0 && ty >= 0 && tx < gs.Width && ty < gs.Height {
				if bomb.canBomb(gameProp, tx, ty, curpr, gp) {
					bsd.AddPos(tx, ty)

					if bomb.Config.BombTargetType == BTTypeRemove {
						gs.Arr[tx][ty] = -1
					} else {
						gs.Arr[tx][ty] = bomb.Config.TargetSymbolCode
					}
				}
			}
		}
	}
}

func (bomb *Bomb) getSourcePos(ctx context.Context, gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, bsd *BombData) ([]int, error) {
	var npos []int

	if bomb.Config.BombSourceType == BSTypePositionCollection {
		npos = make([]int, 0, gs.Width*gs.Height*2)

		pos := gameProp.GetComponentPos(bomb.Config.SourcePositionCollection)
		if len(pos) >= 2 {
			for i := 0; i < len(pos)/2; i++ {
				x := pos[i*2]
				y := pos[i*2+1]
				if bomb.canBomb(gameProp, x, y, curpr, gp) {
					npos = append(npos, x, y)
				}
			}
		}
	} else if bomb.Config.BombSourceType == BSTypeSymbols {
		npos = make([]int, 0, gs.Width*gs.Height*2)

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(bomb.Config.SourceSymbolCodes, s, 0) >= 0 && bomb.canBomb(gameProp, x, y, curpr, gp) {
					npos = append(npos, x, y)
				}
			}
		}
	} else {
		return nil, nil
	}

	if bomb.Config.SelectSourceNumber <= 0 {
		return npos, nil
	}

	if bomb.Config.SelectSourceNumber < len(npos)/2 {
		retpos := make([]int, bomb.Config.SelectSourceNumber*2)

		for i := 0; i < bomb.Config.SelectSourceNumber; i++ {
			cr, err := plugin.Random(ctx, len(npos)/2)
			if err != nil {
				goutils.Error("Bomb.getSourcePos:Random",
					goutils.Err(err))

				return nil, err
			}

			retpos[i*2] = npos[cr*2]
			retpos[i*2+1] = npos[cr*2+1]

			npos = append(npos[:cr*2], npos[(cr+1)*2:]...)
		}

		return retpos, nil
	}

	return npos, nil
}

// playgame
func (bomb *Bomb) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	bsd := cd.(*BombData)

	bsd.onNewStep()

	gs := bomb.GetTargetScene3(gameProp, curpr, prs, 0)

	pos, err := bomb.getSourcePos(context.Background(), gameProp, gs, curpr, gp, plugin, bsd)
	if err != nil {
		goutils.Error("Bomb.OnPlayGame:getSourcePos",
			goutils.Err(err))

		return "", err
	}

	if len(pos) >= 2 {
		sc2 := gs.CloneEx(gameProp.PoolScene)

		for i := 0; i < len(pos)/2; i++ {
			x := pos[i*2]
			y := pos[i*2+1]

			bomb.bomb(gameProp, sc2, x, y, curpr, gp, bsd)
		}

		bomb.AddScene(gameProp, curpr, sc2, &bsd.BasicComponentData)

		bomb.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

		nc := bomb.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := bomb.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (burstSymbols *Bomb) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bsd := cd.(*BurstSymbolsData)

	asciigame.OutputScene("after bomb", pr.Scenes[bsd.UsedScenes[0]], mapSymbolColor)

	return nil
}

// NewComponentData -
func (burstSymbols *Bomb) NewComponentData() IComponentData {
	return &BombData{}
}

func NewBomb(name string) IComponent {
	return &Bomb{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "bombWidth": 3,
// "bombHeight": 3,
// "bombSourceType": "positionCollection",
// "bombData": [
//
//	[
//		0,
//		0,
//		0
//	],
//	[
//		0,
//		1,
//		0
//	],
//	[
//		0,
//		0,
//		0
//	]
//
// ],
// "sourcePositionCollection": "bg-pos-hotzone",
// "selectSourceNumber": 1,
// "ignoreSymbols": [
//
//	"WL"
//
// ],
// "bombTargetType": "remove",
// "outputToComponent": "bg-pos-rmoved",
// "ignoreWinResults": [
//
//	"bg-pay",
//	"bg-pay-upd",
//	"bg-pay-wl"
//
// ]
type jsonBomb struct {
	BombWidth                int      `json:"bombWidth"`
	BombHeight               int      `json:"bombHeight"`
	StrBombSourceType        string   `json:"bombSourceType"`
	SourceSymbols            []string `json:"sourceSymbols"`
	BombData                 [][]int  `json:"bombData"`
	SourcePositionCollection string   `json:"sourcePositionCollection"`
	SelectSourceNumber       int      `json:"selectSourceNumber"`
	IgnoreSymbols            []string `json:"ignoreSymbols"`
	StrBombTargetType        string   `json:"bombTargetType"`
	TargetSymbol             string   `json:"targetSymbol"`
	OutputToComponent        string   `json:"outputToComponent"`
	IgnoreWinResults         []string `json:"ignoreWinResults"`
	IgnorePostionCollections []string `json:"ignorePostionCollections"`
}

func (jcfg *jsonBomb) build() *BombConfig {
	cfg := &BombConfig{
		BombWidth:                jcfg.BombWidth,
		BombHeight:               jcfg.BombHeight,
		StrBombSourceType:        jcfg.StrBombSourceType,
		SourcePositionCollection: jcfg.SourcePositionCollection,
		SelectSourceNumber:       jcfg.SelectSourceNumber,
		StrBombTargetType:        jcfg.StrBombTargetType,
		OutputToComponent:        jcfg.OutputToComponent,
		TargetSymbol:             jcfg.TargetSymbol,
	}

	if len(jcfg.SourceSymbols) > 0 {
		cfg.SourceSymbols = make([]string, len(jcfg.SourceSymbols))
		copy(cfg.SourceSymbols, jcfg.SourceSymbols)
	}

	if len(jcfg.BombData) > 0 {
		cfg.BombData = make([][]int, len(jcfg.BombData))

		for _, arr := range jcfg.BombData {
			narr := make([]int, len(arr))
			copy(narr, arr)
			cfg.BombData = append(cfg.BombData, narr)
		}
	}

	if len(jcfg.IgnoreSymbols) > 0 {
		cfg.IgnoreSymbols = make([]string, len(jcfg.IgnoreSymbols))
		copy(cfg.IgnoreSymbols, jcfg.IgnoreSymbols)
	}

	if len(jcfg.IgnoreWinResults) > 0 {
		cfg.IgnoreWinResults = make([]string, len(jcfg.IgnoreWinResults))
		copy(cfg.IgnoreWinResults, jcfg.IgnoreWinResults)
	}

	if len(jcfg.IgnorePostionCollections) > 0 {
		cfg.IgnorePostionCollections = make([]string, len(jcfg.IgnorePostionCollections))
		copy(cfg.IgnorePostionCollections, jcfg.IgnorePostionCollections)
	}

	return cfg
}

func parseBomb(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseBurstSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseBurstSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonBomb{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseBurstSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseBurstSymbols:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: BombTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
