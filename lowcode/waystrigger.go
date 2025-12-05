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
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WaysTriggerTypeName = "waysTrigger"

type WaysTriggerData struct {
	BasicComponentData
	NextComponent string
	SymbolNum     int
	WildNum       int
	RespinNum     int
	Wins          int
	WinMulti      int
	SymbolCodes   []int
}

// OnNewGame -
func (waysTriggerData *WaysTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	waysTriggerData.BasicComponentData.OnNewGame(gameProp, component)

	waysTriggerData.SymbolCodes = nil
}

// onNewStep -
func (waysTriggerData *WaysTriggerData) onNewStep() {
	waysTriggerData.UsedResults = nil

	waysTriggerData.NextComponent = ""
	waysTriggerData.SymbolNum = 0
	waysTriggerData.WildNum = 0
	waysTriggerData.RespinNum = 0
	waysTriggerData.Wins = 0
	waysTriggerData.WinMulti = 1
}

// Clone
func (waysTriggerData *WaysTriggerData) Clone() IComponentData {
	target := &WaysTriggerData{
		BasicComponentData: waysTriggerData.CloneBasicComponentData(),
		SymbolNum:          waysTriggerData.SymbolNum,
		WildNum:            waysTriggerData.WildNum,
		RespinNum:          waysTriggerData.RespinNum,
		Wins:               waysTriggerData.Wins,
		WinMulti:           waysTriggerData.WinMulti,
		SymbolCodes:        slices.Clone(waysTriggerData.SymbolCodes),
	}

	return target
}

// BuildPBComponentData
func (waysTriggerData *WaysTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WaysTriggerData{
		BasicComponentData: waysTriggerData.BuildPBBasicComponentData(),
		NextComponent:      waysTriggerData.NextComponent,
		SymbolNum:          int32(waysTriggerData.SymbolNum),
		WildNum:            int32(waysTriggerData.WildNum),
		RespinNum:          int32(waysTriggerData.RespinNum),
		Wins:               int32(waysTriggerData.Wins),
		WinMulti:           int32(waysTriggerData.WinMulti),
	}

	return pbcd
}

// GetValEx -
func (waysTriggerData *WaysTriggerData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	switch key {
	case CVSymbolNum:
		return waysTriggerData.SymbolNum, true
	case CVWildNum:
		return waysTriggerData.WildNum, true
	case CVRespinNum:
		return waysTriggerData.RespinNum, true
	case CVWins:
		return waysTriggerData.Wins, true
	case CVResultNum, CVWinResultNum:
		return len(waysTriggerData.UsedResults), true
	}

	return 0, false
}

func (waysTriggerData *WaysTriggerData) SetSymbolCodes(symbolCodes []int) {
	if len(symbolCodes) == 0 {
		waysTriggerData.SymbolCodes = nil

		return
	}

	waysTriggerData.SymbolCodes = slices.Clone(symbolCodes)
}

func (waysTriggerData *WaysTriggerData) GetSymbolCodes() []int {
	return waysTriggerData.SymbolCodes
}

// WaysTriggerConfig - configuration for WaysTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WaysTriggerConfig struct {
	BasicComponentConfig            `yaml:",inline" json:",inline"`
	Symbols                         []string                      `yaml:"symbols" json:"symbols"`                                             // like scatter
	SymbolCodes                     []int                         `yaml:"-" json:"-"`                                                         // like scatter
	Type                            string                        `yaml:"type" json:"type"`                                                   // like scatters
	TriggerType                     SymbolTriggerType             `yaml:"-" json:"-"`                                                         // SymbolTriggerType
	BetTypeString                   string                        `yaml:"betType" json:"betType"`                                             // bet or totalBet or noPay
	BetType                         BetType                       `yaml:"-" json:"-"`                                                         // bet or totalBet or noPay
	OSMulTypeString                 string                        `yaml:"symbolValsMulti" json:"symbolValsMulti"`                             // OtherSceneMultiType
	OSMulType                       OtherSceneMultiType           `yaml:"-" json:"-"`                                                         // OtherSceneMultiType
	MinNum                          int                           `yaml:"minNum" json:"minNum"`                                               // like 3，countscatter 或 countscatterInArea 或 checkLines 或 checkWays 时生效
	WildSymbols                     []string                      `yaml:"wildSymbols" json:"wildSymbols"`                                     // wild etc
	WildSymbolCodes                 []int                         `yaml:"-" json:"-"`                                                         // wild symbolCode
	StrCheckWinType                 string                        `yaml:"checkWinType" json:"checkWinType"`                                   // left2right or right2left or all
	CheckWinType                    CheckWinType                  `yaml:"-" json:"-"`                                                         //
	WinMulti                        int                           `yaml:"winMulti" json:"winMulti"`                                           // winMulti，最后的中奖倍数，默认为1
	JumpToComponent                 string                        `yaml:"jumpToComponent" json:"jumpToComponent"`                             // jump to
	ForceToNext                     bool                          `yaml:"forceToNext" json:"forceToNext"`                                     // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards                          []*Award                      `yaml:"awards" json:"awards"`                                               // 新的奖励系统
	TargetMask                      string                        `yaml:"targetMask" json:"targetMask"`                                       // 如果是scatter这一组判断，可以把结果传递给一个mask
	IsReverse                       bool                          `yaml:"isReverse" json:"isReverse"`                                         // 如果isReverse，表示判定为否才触发
	PiggyBankComponent              string                        `yaml:"piggyBankComponent" json:"piggyBankComponent"`                       // piggyBank component
	OutputToComponent               string                        `yaml:"outputToComponent" json:"outputToComponent"`                         // 将结果给到一个 positionCollection
	IsAddRespinMode                 bool                          `yaml:"isAddRespinMode" json:"isAddRespinMode"`                             // 是否是增加respinNum模式，默认是增加triggerNum模式
	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWeightVW               *sgc7game.ValWeights2         `yaml:"-" json:"-"`                                                         // respin number weight
	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
	SetWinSymbols                   []string                      `yaml:"setWinSymbols" json:"setWinSymbols"`
	GenGigaSymbols2                 string                        `yaml:"genGigaSymbols2" json:"genGigaSymbols2"`
	RowMask                         string                        `yaml:"rowMask" json:"rowMask"`
}

// SetLinkComponent
func (cfg *WaysTriggerConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type WaysTrigger struct {
	*BasicComponent `json:"-"`
	Config          *WaysTriggerConfig `json:"config"`
}

// Init -
func (waysTrigger *WaysTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WaysTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WaysTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WaysTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return waysTrigger.InitEx(cfg, pool)
}

// InitEx -
func (waysTrigger *WaysTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	waysTrigger.Config = cfg.(*WaysTriggerConfig)
	waysTrigger.Config.ComponentType = WaysTriggerTypeName

	waysTrigger.Config.OSMulType = ParseOtherSceneMultiType(waysTrigger.Config.OSMulTypeString)

	for _, s := range waysTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WaysTrigger.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))
		}

		waysTrigger.Config.SymbolCodes = append(waysTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range waysTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WaysTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		waysTrigger.Config.WildSymbolCodes = append(waysTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(waysTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("WaysTrigger.InitEx:ParseSymbolTriggerType",
			slog.String("SymbolTriggerType", waysTrigger.Config.Type),
			goutils.Err(ErrInvalidSymbolTriggerType))

		return ErrInvalidSymbolTriggerType
	}

	waysTrigger.Config.TriggerType = stt

	waysTrigger.Config.BetType = ParseBetType(waysTrigger.Config.BetTypeString)

	for _, award := range waysTrigger.Config.Awards {
		award.Init()
	}

	waysTrigger.Config.CheckWinType = ParseCheckWinType(waysTrigger.Config.StrCheckWinType)

	if waysTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(waysTrigger.Config.RespinNumWeight, waysTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WaysTrigger.InitEx:LoadIntWeights",
				slog.String("Weight", waysTrigger.Config.RespinNumWeight),
				goutils.Err(err))

			return err
		}

		waysTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(waysTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range waysTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, waysTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("WaysTrigger.InitEx:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))

				return err
			}

			waysTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	if waysTrigger.Config.WinMulti <= 0 {
		waysTrigger.Config.WinMulti = 1
	}

	waysTrigger.onInit(&waysTrigger.Config.BasicComponentConfig)

	return nil
}

func (waysTrigger *WaysTrigger) getRowMask(basicCD *BasicComponentData) string {
	str := basicCD.GetConfigVal(CCVRowMask)
	if str != "" {
		return str
	}

	return waysTrigger.Config.RowMask
}

// playgame
func (waysTrigger *WaysTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if waysTrigger.Config.TargetMask != "" {
		gameProp.UseComponent(waysTrigger.Config.TargetMask)

		mask := make([]bool, gs.Width)

		if ret == nil {
			return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, waysTrigger.Config.TargetMask, mask, false)
		}

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, waysTrigger.Config.TargetMask, mask, false)
	}

	return nil
}

// procPositionCollection
func (waysTrigger *WaysTrigger) procPositionCollection(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	cd *WaysTriggerData) error {

	if waysTrigger.Config.OutputToComponent != "" {
		pcd := gameProp.GetComponentDataWithName(waysTrigger.Config.OutputToComponent)
		if pcd != nil {
			gameProp.UseComponent(waysTrigger.Config.OutputToComponent)
			pc := gameProp.Components.MapComponents[waysTrigger.Config.OutputToComponent]

			for _, ri := range cd.UsedResults {
				ret := curpr.Results[ri]

				for i := 0; i < len(ret.Pos)/2; i++ {
					pc.AddPos(pcd, ret.Pos[i*2], ret.Pos[i*2+1])
				}
			}
		}
	}

	return nil
}

func (waysTrigger *WaysTrigger) getSymbols(gameProp *GameProperty, cd *WaysTriggerData) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	if len(cd.SymbolCodes) == 0 {
		return waysTrigger.Config.SymbolCodes
	}

	return cd.SymbolCodes
}

// CanTriggerWithScene -
func (waysTrigger *WaysTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) {
	return waysTrigger.canTrigger(gameProp, gs, nil, curpr, stake, icd.(*WaysTriggerData))
}

func (waysTrigger *WaysTrigger) getGigaData(gameProp *GameProperty) (*GenGigaSymbols2Data, error) {
	icd := gameProp.GetComponentDataWithName(waysTrigger.Config.GenGigaSymbols2)
	if icd == nil {
		return nil, nil
	}

	gigacd, isok := icd.(*GenGigaSymbols2Data)
	if !isok {
		goutils.Error("WaysTrigger.getGigaData:TypeAssert",
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	return gigacd, nil
}

// CanTrigger -
func (waysTrigger *WaysTrigger) canTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene, _ *sgc7game.PlayResult, stake *sgc7game.Stake, cd *WaysTriggerData) (bool, []*sgc7game.Result) {
	isTrigger := false
	lst := []*sgc7game.Result{}
	symbols := waysTrigger.getSymbols(gameProp, cd)
	if len(symbols) == 0 {
		if waysTrigger.Config.IsReverse {
			isTrigger = !isTrigger
		}

		return isTrigger, lst
	}

	rowMask := waysTrigger.getRowMask(&cd.BasicComponentData)
	gigacd, err := waysTrigger.getGigaData(gameProp)
	if err != nil {
		goutils.Error("WaysTrigger.canTrigger:getGigaData",
			goutils.Err(err))

		return false, nil
	}

	if rowMask != "" {
		maskCompData := gameProp.GetComponentDataWithName(rowMask)
		if maskCompData == nil {
			goutils.Error("WaysTrigger.canTrigger:GetComponentDataWithName",
				goutils.Err(ErrInvalidComponentConfig))

			return false, nil
		}

		validMask := maskCompData.GetMask()

		switch waysTrigger.Config.TriggerType {
		case STTypeWays:
			if os != nil {
				if gigacd != nil {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						},
						func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							if !validMask[y] {
								return false
							}

							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							gigad := gigacd.getGigaData(x, y)
							if gigad != nil {
								return gigad.SymbolCode
							}

							return cursymbol
						},
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return os.Arr[x][y]
						})

					lst = append(lst, currets...)
				} else {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						},
						func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							if !validMask[y] {
								return false
							}

							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							return cursymbol
						},
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return os.Arr[x][y]
						})

					lst = append(lst, currets...)
				}
			} else {
				if gigacd != nil {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							if !validMask[y] {
								return false
							}

							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							gigad := gigacd.getGigaData(x, y)
							if gigad != nil {
								return gigad.SymbolCode
							}

							return cursymbol
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return 1
						})

					lst = append(lst, currets...)
				} else {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						},
						func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							if !validMask[y] {
								return false
							}

							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							return cursymbol
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return 1
						})

					lst = append(lst, currets...)
				}
			}

			if len(lst) > 0 {
				isTrigger = true
			}
		case STTypeCheckWays:
			currets := sgc7game.CheckWays(gs, waysTrigger.Config.MinNum,
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					if !validMask[y] {
						return false
					}

					return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				})

			lst = append(lst, currets...)

			if len(lst) > 0 {
				isTrigger = true
			}
		}
	} else {
		switch waysTrigger.Config.TriggerType {
		case STTypeWays:
			if os != nil {
				if gigacd != nil {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							gigad := gigacd.getGigaData(x, y)
							if gigad != nil {
								return gigad.SymbolCode
							}

							return cursymbol
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return os.Arr[x][y]
						})

					lst = append(lst, currets...)
				} else {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							gigad := gigacd.getGigaData(x, y)
							if gigad != nil {
								return gigad.SymbolCode
							}

							return cursymbol
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return os.Arr[x][y]
						})

					lst = append(lst, currets...)
				}
			} else {
				if gigacd != nil {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							gigad := gigacd.getGigaData(x, y)
							if gigad != nil {
								return gigad.SymbolCode
							}

							return cursymbol
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return 1
						})

					lst = append(lst, currets...)
				} else {
					currets := sgc7game.CheckWays5(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, waysTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
							return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
						}, func(cursymbol int, x, y int) int {
							return cursymbol
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(x, y int) int {
							return 1
						})

					lst = append(lst, currets...)
				}
			}

			if len(lst) > 0 {
				isTrigger = true
			}
		case STTypeCheckWays:
			currets := sgc7game.CheckWays(gs, waysTrigger.Config.MinNum,
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				})

			lst = append(lst, currets...)

			if len(lst) > 0 {
				isTrigger = true
			}
		}
	}

	if waysTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

func (waysTrigger *WaysTrigger) procWinSymbols(gameProp *GameProperty, lst []*sgc7game.Result) {
	if len(waysTrigger.Config.SetWinSymbols) > 0 {
		if len(lst) == 0 {
			for _, v := range waysTrigger.Config.SetWinSymbols {
				curicd := gameProp.GetComponentDataWithName(v)
				if curicd != nil {
					curicd.SetSymbolCodes(nil)
				}
			}

			return
		}

		symbolCodes := make([]int, 0, len(lst))

		for _, v := range lst {
			if !slices.Contains(symbolCodes, v.Symbol) {
				symbolCodes = append(symbolCodes, v.Symbol)
			}
		}

		for _, v := range waysTrigger.Config.SetWinSymbols {
			curicd := gameProp.GetComponentDataWithName(v)
			if curicd != nil {
				curicd.SetSymbolCodes(symbolCodes)
			}
		}
	}
}

// procWins
func (waysTrigger *WaysTrigger) procWins(gameProp *GameProperty, curpr *sgc7game.PlayResult, std *WaysTriggerData, lst []*sgc7game.Result) (int, error) {
	if waysTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0

			waysTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds
		}

		return 0, nil
	}

	std.WinMulti = waysTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti
		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin

		waysTrigger.AddResult(curpr, v, &std.BasicComponentData)

		std.SymbolNum += v.SymbolNums
		std.WildNum += v.Wilds
	}

	if std.Wins > 0 {
		if waysTrigger.Config.PiggyBankComponent != "" {
			cd := gameProp.GetCurComponentDataWithName(waysTrigger.Config.PiggyBankComponent)
			if cd == nil {
				goutils.Error("ScatterTrigger.procWins:GetCurComponentDataWithName",
					slog.String("PiggyBankComponent", waysTrigger.Config.PiggyBankComponent),
					goutils.Err(ErrInvalidComponent))

				return 0, ErrInvalidComponent
			}

			cd.ChgConfigIntVal(CCVSavedMoney, std.Wins)

			for _, v := range lst {
				v.IsNoPayNow = true
			}

			gameProp.UseComponent(waysTrigger.Config.PiggyBankComponent)
		}
	}

	return std.Wins, nil
}

// calcRespinNum
func (waysTrigger *WaysTrigger) calcRespinNum(plugin sgc7plugin.IPlugin, ret *sgc7game.Result) (int, error) {

	if len(waysTrigger.Config.RespinNumWeightWithScatterNumVW) > 0 {
		vw2, isok := waysTrigger.Config.RespinNumWeightWithScatterNumVW[ret.SymbolNums]
		if isok {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
					slog.Int("SymbolNum", ret.SymbolNums),
					goutils.Err(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(waysTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := waysTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWithScatterNum",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if waysTrigger.Config.RespinNumWeightVW != nil {
		cr, err := waysTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightVW",
				goutils.Err(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if waysTrigger.Config.RespinNum > 0 {
		return waysTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// OnProcControllers -
func (waysTrigger *WaysTrigger) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(waysTrigger.Config.Awards) > 0 {
		gameProp.procAwards(plugin, waysTrigger.Config.Awards, curpr, gp)
	}
}

// playgame
func (waysTrigger *WaysTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	std := icd.(*WaysTriggerData)
	std.onNewStep()

	gs := waysTrigger.GetTargetScene3(gameProp, curpr, prs, 0)

	var os *sgc7game.GameScene
	if waysTrigger.Config.OSMulType != OSMTNone {
		os = waysTrigger.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	isTrigger, lst := waysTrigger.canTrigger(gameProp, gs, os, curpr, stake, std)

	if isTrigger {
		waysTrigger.procWins(gameProp, curpr, std, lst)
		waysTrigger.procWinSymbols(gameProp, lst)

		respinNum, err := waysTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:calcRespinNum",
				goutils.Err(err))

			return "", nil
		}

		std.RespinNum = respinNum

		err = waysTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}

		err = waysTrigger.procPositionCollection(gameProp, curpr, std)
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:procPositionCollection",
				goutils.Err(err))

			return "", err
		}

		waysTrigger.ProcControllers(gameProp, plugin, curpr, gp, 0, "")

		if waysTrigger.Config.JumpToComponent != "" {
			if gameProp.IsRespin(waysTrigger.Config.JumpToComponent) {
				// 如果jumpto是一个respin，那么就需要trigger respin
				if std.RespinNum == 0 {
					if waysTrigger.Config.ForceToNext {
						std.NextComponent = waysTrigger.Config.DefaultNextComponent
					} else {
						rn := gameProp.GetLastRespinNum(waysTrigger.Config.JumpToComponent)
						if rn > 0 {
							gameProp.TriggerRespin(plugin, curpr, gp, 0, waysTrigger.Config.JumpToComponent, !waysTrigger.Config.IsAddRespinMode)

							lst[0].Type = sgc7game.RTFreeGame
							lst[0].Value = rn
						}
					}
				} else {
					// 如果jumpto是respin，需要treigger这个respin
					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, waysTrigger.Config.JumpToComponent, !waysTrigger.Config.IsAddRespinMode)

					lst[0].Type = sgc7game.RTFreeGame
					lst[0].Value = std.RespinNum
				}
			}

			std.NextComponent = waysTrigger.Config.JumpToComponent

			nc := waysTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nc, nil
		}

		nc := waysTrigger.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	} else {
		err := waysTrigger.procMask(gs, gameProp, curpr, gp, plugin, nil)
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}
	}

	nc := waysTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (waysTrigger *WaysTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WaysTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", waysTrigger.Name, std.NextComponent)
	}

	return nil
}

// NewComponentData -
func (waysTrigger *WaysTrigger) NewComponentData() IComponentData {
	return &WaysTriggerData{}
}

func (waysTrigger *WaysTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		if winMulti <= 0 {
			return 1
		}

		return winMulti
	}

	if waysTrigger.Config.WinMulti <= 0 {
		return 1
	}

	return waysTrigger.Config.WinMulti
}

// NewStats2 -
func (waysTrigger *WaysTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (waysTrigger *WaysTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	waysTrigger.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*WaysTriggerData)

	s2.ProcStatsWins(waysTrigger.Name, int64(cd.Wins))
}

// GetAllLinkComponents - get all link components
func (waysTrigger *WaysTrigger) GetAllLinkComponents() []string {
	return []string{waysTrigger.Config.DefaultNextComponent, waysTrigger.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (waysTrigger *WaysTrigger) GetNextLinkComponents() []string {
	return []string{waysTrigger.Config.DefaultNextComponent, waysTrigger.Config.JumpToComponent}
}

func NewWaysTrigger(name string) IComponent {
	return &WaysTrigger{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"triggerType": "lines",
//	"betType": "bet",
//	"checkWinType": "left2right",
//	"symbols": [
//		"WL",
//		"A",
//		"B",
//		"C",
//		"D",
//		"E",
//		"F",
//		"G",
//		"H",
//		"J",
//		"K",
//		"L"
//	],
//	"wildSymbols": [
//		"WL"
//	]
//
// "rowMask": "mask-height4"
// "genGigaSymbols2": "bg-gengiga"
type jsonWaysTrigger struct {
	Symbols             []string `json:"symbols"`
	TriggerType         string   `json:"triggerType"`
	BetType             string   `json:"betType"`
	SymbolValsMulti     string   `json:"symbolValsMulti"`
	MinNum              int      `json:"minNum"`
	WildSymbols         []string `json:"wildSymbols"`
	WinMulti            int      `json:"winMulti"`
	PutMoneyInPiggyBank string   `json:"putMoneyInPiggyBank"`
	OutputToComponent   string   `json:"outputToComponent"`
	RowMask             string   `json:"rowMask"`
	GenGigaSymbols2     string   `json:"genGigaSymbols2"`
}

func (jcfg *jsonWaysTrigger) build() *WaysTriggerConfig {
	cfg := &WaysTriggerConfig{
		Symbols:            jcfg.Symbols,
		Type:               jcfg.TriggerType,
		BetTypeString:      jcfg.BetType,
		MinNum:             jcfg.MinNum,
		WildSymbols:        jcfg.WildSymbols,
		WinMulti:           jcfg.WinMulti,
		PiggyBankComponent: jcfg.PutMoneyInPiggyBank,
		OSMulTypeString:    jcfg.SymbolValsMulti,
		OutputToComponent:  jcfg.OutputToComponent,
		RowMask:            jcfg.RowMask,
		GenGigaSymbols2:    jcfg.GenGigaSymbols2,
	}

	return cfg
}

func parseWaysTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWaysTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWaysTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWaysTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWaysTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseWaysTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WaysTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
