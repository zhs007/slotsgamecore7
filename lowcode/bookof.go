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
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const BookOfTypeName = "bookOf"

type BookOfData struct {
	BasicComponentData
	Symbols []int
}

// OnNewGame -
func (bookOfData *BookOfData) OnNewGame() {
	bookOfData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (bookOfData *BookOfData) OnNewStep() {
	bookOfData.BasicComponentData.OnNewStep()

	bookOfData.Symbols = nil
}

// BuildPBComponentData
func (bookOfData *BookOfData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BookOfData{
		BasicComponentData: bookOfData.BuildPBBasicComponentData(),
	}

	for _, v := range bookOfData.Symbols {
		pbcd.Symbols = append(pbcd.Symbols, int32(v))
	}

	return pbcd
}

// BookOfConfig - configuration for BookOf feature
type BookOfConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetType              string   `yaml:"betType" json:"betType"`         // bet or totalBet
	WildSymbols          []string `yaml:"wildSymbols" json:"wildSymbols"` // 可以不要wild
	WildSymbolCodes      []int    `yaml:"-" json:"-"`
	ForceTrigger         bool     `yaml:"forceTrigger" json:"forceTrigger"`
	WeightTrigger        string   `yaml:"weightTrigger" json:"weightTrigger"`
	WeightSymbolNum      string   `yaml:"weightSymbolNum" json:"weightSymbolNum"`
	WeightSymbol         string   `yaml:"weightSymbol" json:"weightSymbol"`
	ForceSymbolNum       int      `yaml:"forceSymbolNum" json:"forceSymbolNum"`
	SymbolRNG            string   `yaml:"symbolRNG" json:"symbolRNG"`               // 只在ForceSymbolNum为1时有效
	SymbolCollection     string   `yaml:"symbolCollection" json:"symbolCollection"` // 图标从一个SymbolCollection里获取
}

type BookOf struct {
	*BasicComponent `json:"-"`
	Config          *BookOfConfig         `json:"config"`
	WeightTrigger   *sgc7game.ValWeights2 `json:"-"`
	WeightSymbolNum *sgc7game.ValWeights2 `json:"-"`
	WeightSymbol    *sgc7game.ValWeights2 `json:"-"`
}

// Init -
func (bookof *BookOf) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BookOf.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &BookOfConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BookOf.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return bookof.InitEx(cfg, pool)
}

// InitEx -
func (bookof *BookOf) InitEx(cfg any, pool *GamePropertyPool) error {
	bookof.Config = cfg.(*BookOfConfig)
	bookof.Config.ComponentType = BookOfTypeName

	if bookof.Config.WeightTrigger != "" {
		vw2, err := pool.LoadStrWeights(bookof.Config.WeightTrigger, bookof.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BookOf.Init:LoadValWeights",
				zap.String("Weight", bookof.Config.WeightTrigger),
				zap.Error(err))

			return err
		}

		bookof.WeightTrigger = vw2
	}

	if bookof.Config.WeightSymbolNum != "" {
		vw2, err := pool.LoadStrWeights(bookof.Config.WeightSymbolNum, bookof.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BookOf.Init:LoadValWeights",
				zap.String("Weight", bookof.Config.WeightSymbolNum),
				zap.Error(err))

			return err
		}

		bookof.WeightSymbolNum = vw2
	}

	if bookof.Config.WeightSymbol != "" {
		vw2, err := pool.LoadSymbolWeights(bookof.Config.WeightSymbol, "val", "weight", pool.DefaultPaytables, bookof.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BookOf.Init:LoadSymbolWeights",
				zap.String("Weight", bookof.Config.WeightSymbol),
				zap.Error(err))

			return err
		}

		bookof.WeightSymbol = vw2
	}

	for _, v := range bookof.Config.WildSymbols {
		bookof.Config.WildSymbolCodes = append(bookof.Config.WildSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	bookof.onInit(&bookof.Config.BasicComponentConfig)

	return nil
}

// playgame
func (bookof *BookOf) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	bookof.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[bookof.Name].(*BookOfData)

	isTrigger := bookof.Config.ForceTrigger

	if !isTrigger && bookof.WeightTrigger != nil {
		iv, err := bookof.WeightTrigger.RandVal(plugin)
		if err != nil {
			goutils.Error("bookof.OnPlayGame:WeightTrigger.RandVal",
				zap.Error(err))

			return err
		}

		if iv.Int() != 0 {
			isTrigger = true
		}
	}

	if isTrigger {
		if bookof.Config.SymbolCollection != "" {
			sccd := gameProp.MapComponentData[bookof.Config.SymbolCollection].(*SymbolCollectionData)
			if sccd == nil {
				goutils.Error("BookOf.OnPlayGame",
					zap.Error(ErrIvalidSymbolCollection))

				return ErrIvalidSymbolCollection
			}

			if len(sccd.SymbolCodes) > 0 {
				cd.Symbols = make([]int, len(sccd.SymbolCodes))
				copy(cd.Symbols, sccd.SymbolCodes)
			} else {
				cd.Symbols = nil
			}

		} else {
			symbolNum := bookof.Config.ForceSymbolNum

			if symbolNum <= 0 && bookof.WeightSymbolNum != nil {
				iv, err := bookof.WeightSymbolNum.RandVal(plugin)
				if err != nil {
					goutils.Error("bookof.OnPlayGame:WeightSymbolNum.RandVal",
						zap.Error(err))

					return err
				}

				symbolNum = iv.Int()
			}

			if bookof.Config.ForceSymbolNum == 1 && bookof.Config.SymbolRNG != "" {
				rng := gameProp.GetTagInt(bookof.Config.SymbolRNG)
				cs := bookof.WeightSymbol.Vals[rng]

				cd.Symbols = append(cd.Symbols, cs.Int())
			} else {
				curWeight := bookof.WeightSymbol.Clone()

				for i := 0; i < symbolNum; i++ {
					cs, err := curWeight.RandVal(plugin)
					if err != nil {
						goutils.Error("bookof.OnPlayGame:curWeight.RandVal",
							zap.Error(err))

						return err
					}

					cd.Symbols = append(cd.Symbols, cs.Int())

					err = curWeight.RemoveVal(cs)
					if err != nil {
						goutils.Error("bookof.OnPlayGame:curWeight.RemoveVal",
							zap.Error(err))

						return err
					}
				}
			}
		}

		gs := bookof.GetTargetScene3(gameProp, curpr, &cd.BasicComponentData, bookof.Name, "", 0)

		for _, s := range cd.Symbols {
			ngs, err := bookof.procBookOfScene(gameProp, gs, s)
			if err != nil {
				goutils.Error("bookof.OnPlayGame:procBookOfScene",
					zap.Error(err))

				return err
			}

			bookof.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

			scr := sgc7game.CalcScatter3(gs, gameProp.CurPaytables, s, gameProp.GetBet(stake, bookof.Config.BetType)*gameProp.GetVal(GamePropCurLineNum), 1,
				func(scatter int, cursymbol int) bool {
					return cursymbol == s || goutils.IndexOfIntSlice(bookof.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, true)
			if scr != nil {
				bookof.AddResult(curpr, scr, &cd.BasicComponentData)
			}
		}

		bookof.onStepEnd(gameProp, curpr, gp, "")
	} else {
		bookof.onStepEnd(gameProp, curpr, gp, "")
	}

	return nil
}

// procBookOfScene - outpur to asciigame
func (bookof *BookOf) procBookOfScene(gameProp *GameProperty, gs *sgc7game.GameScene, symbol int) (*sgc7game.GameScene, error) {
	ngs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range gs.Arr {
		hass := false
		for _, s := range arr {
			if s == symbol || goutils.IndexOfIntSlice(bookof.Config.WildSymbolCodes, s, 0) >= 0 {
				hass = true

				break
			}
		}

		if hass {
			for y := range arr {
				ngs.Arr[x][y] = symbol
			}
		}
	}

	return ngs, nil
}

// OnAsciiGame - outpur to asciigame
func (bookof *BookOf) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[bookof.Name].(*BookOfData)

	if len(cd.Symbols) > 0 {
		strsymbols := ""
		for i, v := range cd.Symbols {
			if i > 0 {
				strsymbols += ", "
			}

			strsymbols += gameProp.CurPaytables.GetStringFromInt(v)
		}

		fmt.Printf("The BookOf Symbols is %v\n", strsymbols)

		for i, si := range cd.UsedScenes {
			asciigame.OutputScene(fmt.Sprintf("The symbols for BookOf - %v ", i+1), pr.Scenes[si], mapSymbolColor)
		}

		asciigame.OutputResults(fmt.Sprintf("%v wins", bookof.Name), pr, func(i int, ret *sgc7game.Result) bool {
			return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
		}, mapSymbolColor)
	}

	return nil
}

// OnStats
func (bookof *BookOf) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	wins := int64(0)
	isTrigger := false

	for _, v := range lst {
		gp, isok := v.CurGameModParams.(*GameParams)
		if isok {
			curComponent, isok := gp.MapComponentMsgs[bookof.Name]
			if isok {
				curwins, err := bookof.OnStatsWithPB(feature, curComponent, v)
				if err != nil {
					goutils.Error("BookOf.OnStats",
						zap.Error(err))

					continue
				}

				isTrigger = true
				wins += curwins
			}
		}
	}

	feature.CurWins.AddWin(int(wins) * 100 / int(stake.CashBet))

	if feature.Parent != nil {
		totalwins := int64(0)

		for _, v := range lst {
			totalwins += v.CashWin
		}

		feature.AllWins.AddWin(int(totalwins) * 100 / int(stake.CashBet))
	}

	return isTrigger, stake.CashBet, wins
}

// OnStatsWithPB -
func (bookof *BookOf) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.BookOfData)
	if !isok {
		goutils.Error("BookOf.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return bookof.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (bookof *BookOf) NewComponentData() IComponentData {
	return &BookOfData{}
}

// EachUsedResults -
func (bookof *BookOf) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.BookOfData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BookOf.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewBookOf(name string) IComponent {
	return &BookOf{
		BasicComponent: NewBasicComponent(name),
	}
}
