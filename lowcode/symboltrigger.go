package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const SymbolTriggerTypeName = "symbolTrigger"

type SymbolTriggerType int

const (
	STTypeUnknow             SymbolTriggerType = 0
	STTypeLines              SymbolTriggerType = 1
	STTypeWays               SymbolTriggerType = 2
	STTypeScatters           SymbolTriggerType = 3
	STTypeCountScatter       SymbolTriggerType = 4
	STTypeCountScatterInArea SymbolTriggerType = 5
)

func ParseSymbolTriggerType(str string) SymbolTriggerType {
	if str == "lines" {
		return STTypeLines
	} else if str == "ways" {
		return STTypeWays
	} else if str == "scatters" {
		return STTypeScatters
	} else if str == "countscatter" {
		return STTypeCountScatter
	} else if str == "countscatterInArea" {
		return STTypeCountScatterInArea
	}

	return STTypeUnknow
}

type BetType int

const (
	BTypeNoPay    BetType = 0
	BTypeBet      BetType = 1
	BTypeTotalBet BetType = 2
)

func ParseBetType(str string) BetType {
	if str == "bet" {
		return BTypeBet
	} else if str == "totalBet" {
		return BTypeTotalBet
	}

	return BTypeNoPay
}

type SymbolTriggerData struct {
	BasicComponentData
	NextComponent string
}

// OnNewGame -
func (symbolTriggerData *SymbolTriggerData) OnNewGame() {
	symbolTriggerData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (symbolTriggerData *SymbolTriggerData) OnNewStep() {
	symbolTriggerData.BasicComponentData.OnNewStep()

	symbolTriggerData.NextComponent = ""
}

// BuildPBComponentData
func (symbolTriggerData *SymbolTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolTriggerData{
		BasicComponentData: symbolTriggerData.BuildPBBasicComponentData(),
		NextComponent:      symbolTriggerData.NextComponent,
	}

	return pbcd
}

// SymbolTriggerConfig - configuration for SymbolTrigger
type SymbolTriggerConfig struct {
	BasicComponentConfig        `yaml:",inline" json:",inline"`
	Symbol                      string            `yaml:"symbol" json:"symbol"`                           // like scatter
	SymbolCode                  int               `yaml:"-" json:"-"`                                     // like scatter
	Type                        string            `yaml:"type" json:"type"`                               // like scatters
	TriggerType                 SymbolTriggerType `yaml:"-" json:"-"`                                     // SymbolTriggerType
	BetTypeString               string            `yaml:"betType" json:"betType"`                         // bet or totalBet or noPay
	BetType                     BetType           `yaml:"-" json:"-"`                                     // bet or totalBet or noPay
	MinNum                      int               `yaml:"minNum" json:"minNum"`                           // like 3，STTypeCountScatter 或 STTypeCountScatterInArea 时生效
	WildSymbols                 []string          `yaml:"wildSymbols" json:"wildSymbols"`                 // wild etc
	WildSymbolCodes             []int             `yaml:"-" json:"-"`                                     // wild symbolCode
	PosArea                     []int             `yaml:"posArea" json:"posArea"`                         // 只在countscatterInArea时生效，[minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
	CountScatterPayAs           string            `yaml:"countScatterPayAs" json:"countScatterPayAs"`     // countscatter时，按什么符号赔付
	SymbolCodeCountScatterPayAs int               `yaml:"-" json:"-"`                                     // countscatter时，按什么符号赔付
	JumpToComponent             string            `yaml:"jumpToComponent" json:"jumpToComponent"`         // jump to
	ForcedJump                  bool              `yaml:"forcedJump" json:"forcedJump"`                   // 强制跳转，中断当前流程
	TagSymbolNum                string            `yaml:"tagSymbolNum" json:"tagSymbolNum"`               // 这里可以将symbol数量记下来，别的地方能获取到
	Awards                      []*Award          `yaml:"awards" json:"awards"`                           // 新的奖励系统
	SymbolAwardsWeights         *AwardsWeights    `yaml:"symbolAwardsWeights" json:"symbolAwardsWeights"` // 每个中奖符号随机一组奖励
	TargetMask                  string            `yaml:"targetMask" json:"targetMask"`                   // 如果是scatter这一组判断，可以把结果传递给一个mask
}

type SymbolTrigger struct {
	*BasicComponent `json:"-"`
	Config          *SymbolTriggerConfig `json:"config"`
}

// Init -
func (symbolTrigger *SymbolTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolTrigger.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolTrigger.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return symbolTrigger.InitEx(cfg, pool)
}

// InitEx -
func (symbolTrigger *SymbolTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolTrigger.Config = cfg.(*SymbolTriggerConfig)
	symbolTrigger.Config.ComponentType = SymbolTriggerTypeName

	sc, isok := pool.DefaultPaytables.MapSymbols[symbolTrigger.Config.Symbol]
	if !isok {
		goutils.Error("SymbolTrigger.InitEx:Symbol",
			zap.String("symbol", symbolTrigger.Config.Symbol),
			zap.Error(ErrIvalidSymbol))
	}

	symbolTrigger.Config.SymbolCode = sc

	sc, isok = pool.DefaultPaytables.MapSymbols[symbolTrigger.Config.CountScatterPayAs]
	if isok {
		symbolTrigger.Config.SymbolCodeCountScatterPayAs = sc
	}

	for _, s := range symbolTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolTrigger.InitEx:WildSymbols",
				zap.String("symbol", s),
				zap.Error(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		symbolTrigger.Config.WildSymbolCodes = append(symbolTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(symbolTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("SymbolTrigger.InitEx:WildSymbols",
			zap.String("SymbolTriggerType", symbolTrigger.Config.Type),
			zap.Error(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	symbolTrigger.Config.TriggerType = stt

	symbolTrigger.Config.BetType = ParseBetType(symbolTrigger.Config.BetTypeString)

	for _, award := range symbolTrigger.Config.Awards {
		award.Init()
	}

	if symbolTrigger.Config.SymbolAwardsWeights != nil {
		symbolTrigger.Config.SymbolAwardsWeights.Init()
	}

	symbolTrigger.onInit(&symbolTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (symbolTrigger *SymbolTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if symbolTrigger.Config.TargetMask != "" {
		mask := make([]bool, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, symbolTrigger.Config.TargetMask, mask)
	}

	return nil
}

// playgame
func (symbolTrigger *SymbolTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := gameProp.MapComponentData[symbolTrigger.Name].(*SymbolTriggerData)

	gs := symbolTrigger.GetTargetScene(gameProp, curpr, &std.BasicComponentData, "")

	isTrigger := false
	var ret *sgc7game.Result

	if symbolTrigger.Config.TriggerType == STTypeScatters {
		ret = sgc7game.CalcScatter4(gs, gameProp.CurPaytables, symbolTrigger.Config.SymbolCode, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, true)

		if ret != nil {
			if symbolTrigger.Config.BetType == BTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				gameProp.ProcMulti(ret)
			}

			symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
			isTrigger = true
		}
	} else if symbolTrigger.Config.TriggerType == STTypeCountScatter {
		ret = sgc7game.CalcScatterEx(gs, symbolTrigger.Config.SymbolCode, symbolTrigger.Config.MinNum, func(scatter int, cursymbol int) bool {
			return cursymbol == scatter || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		})

		if ret != nil {
			if symbolTrigger.Config.BetType == BTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				if symbolTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet2(stake, symbolTrigger.Config.BetType)
				}

				gameProp.ProcMulti(ret)
			}

			symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
			isTrigger = true
		}
	} else if symbolTrigger.Config.TriggerType == STTypeCountScatterInArea {
		ret = sgc7game.CountScatterInArea(gs, symbolTrigger.Config.SymbolCode, symbolTrigger.Config.MinNum,
			func(x, y int) bool {
				return x >= symbolTrigger.Config.PosArea[0] && x <= symbolTrigger.Config.PosArea[1] && y >= symbolTrigger.Config.PosArea[2] && y <= symbolTrigger.Config.PosArea[3]
			},
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			})

		if ret != nil {
			if symbolTrigger.Config.BetType == BTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				if symbolTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet2(stake, symbolTrigger.Config.BetType)
				}

				gameProp.ProcMulti(ret)
			}

			symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
			isTrigger = true
		}
	}

	if isTrigger {
		err := symbolTrigger.procMask(gs, gameProp, curpr, gp, plugin, ret)
		if err != nil {
			goutils.Error("SymbolTrigger.OnPlayGame:procMask",
				zap.Error(err))

			return err
		}

		if symbolTrigger.Config.TagSymbolNum != "" {
			gameProp.TagInt(symbolTrigger.Config.TagSymbolNum, ret.SymbolNums)
		}

		if len(symbolTrigger.Config.Awards) > 0 {
			gameProp.procAwards(plugin, symbolTrigger.Config.Awards, curpr, gp)
		}

		if symbolTrigger.Config.SymbolAwardsWeights != nil {
			for i := 0; i < ret.SymbolNums; i++ {
				node, err := symbolTrigger.Config.SymbolAwardsWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("SymbolTrigger.OnPlayGame:SymbolAwardsWeights.RandVal",
						zap.Error(err))

					return nil
				}

				gameProp.procAwards(plugin, node.Awards, curpr, gp)
			}
		}

		if symbolTrigger.Config.JumpToComponent != "" {
			std.NextComponent = symbolTrigger.Config.JumpToComponent

			symbolTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nil
		}
	}

	symbolTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolTrigger *SymbolTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	std := gameProp.MapComponentData[symbolTrigger.Name].(*SymbolTriggerData)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v", symbolTrigger.Name, std.NextComponent)
	}

	return nil
}

// OnStats
func (symbolTrigger *SymbolTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewComponentData -
func (symbolTrigger *SymbolTrigger) NewComponentData() IComponentData {
	return &SymbolTriggerData{}
}

func NewSymbolTrigger(name string) IComponent {
	return &SymbolTrigger{
		BasicComponent: NewBasicComponent(name),
	}
}
