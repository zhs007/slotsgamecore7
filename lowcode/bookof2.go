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

const BookOf2TypeName = "bookOf2"

type BookOf2Data struct {
	BasicComponentData
	Symbols []int
}

// OnNewGame -
func (bookOf2Data *BookOf2Data) OnNewGame() {
	bookOf2Data.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (bookOf2Data *BookOf2Data) OnNewStep() {
	bookOf2Data.BasicComponentData.OnNewStep()

	bookOf2Data.Symbols = nil
}

// BuildPBComponentData
func (bookOf2Data *BookOf2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BookOf2Data{
		BasicComponentData: bookOf2Data.BuildPBBasicComponentData(),
	}

	for _, v := range bookOf2Data.Symbols {
		pbcd.Symbols = append(pbcd.Symbols, int32(v))
	}

	return pbcd
}

// BookOf2Config - configuration for BookOf feature
type BookOf2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetType              string                `yaml:"betType" json:"betType"`         // bet or totalBet
	WildSymbols          []string              `yaml:"wildSymbols" json:"wildSymbols"` // 可以不要wild
	WildSymbolCodes      []int                 `yaml:"-" json:"-"`
	ForceTrigger         bool                  `yaml:"forceTrigger" json:"forceTrigger"`
	WeightTrigger        string                `yaml:"weightTrigger" json:"weightTrigger"`
	WeightTriggerVW      *sgc7game.ValWeights2 `json:"-"`
	WeightSymbolNum      string                `yaml:"weightSymbolNum" json:"weightSymbolNum"`
	WeightSymbolNumVW    *sgc7game.ValWeights2 `json:"-"`
	WeightSymbol         string                `yaml:"weightSymbol" json:"weightSymbol"`
	WeightSymbolVW       *sgc7game.ValWeights2 `json:"-"`
	ForceSymbolNum       int                   `yaml:"forceSymbolNum" json:"forceSymbolNum"`
	SymbolRNG            string                `yaml:"symbolRNG" json:"symbolRNG"`               // 只在ForceSymbolNum为1时有效
	SymbolCollection     string                `yaml:"symbolCollection" json:"symbolCollection"` // 图标从一个SymbolCollection里获取
	MapSymbols           map[string][]string   `yaml:"mapSymbols" json:"mapSymbols"`             // 一个符号对应1+wild数量个变化，第1个是被选中，第2个是wild[0]变化效果，第3个是wild[1]变化效果
	MapSymbolCodes       map[int][]int         `yaml:"-" json:"-"`
}

type BookOf2 struct {
	*BasicComponent `json:"-"`
	Config          *BookOf2Config `json:"config"`
}

// Init -
func (bookof2 *BookOf2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BookOf2.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &BookOf2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BookOf2.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return bookof2.InitEx(cfg, pool)
}

// InitEx -
func (bookof2 *BookOf2) InitEx(cfg any, pool *GamePropertyPool) error {
	bookof2.Config = cfg.(*BookOf2Config)
	bookof2.Config.ComponentType = BookOf2TypeName

	if bookof2.Config.WeightTrigger != "" {
		vw2, err := pool.LoadStrWeights(bookof2.Config.WeightTrigger, bookof2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BookOf2.InitEx:LoadValWeights",
				zap.String("Weight", bookof2.Config.WeightTrigger),
				zap.Error(err))

			return err
		}

		bookof2.Config.WeightTriggerVW = vw2
	}

	if bookof2.Config.WeightSymbolNum != "" {
		vw2, err := pool.LoadStrWeights(bookof2.Config.WeightSymbolNum, bookof2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BookOf2.InitEx:LoadValWeights",
				zap.String("Weight", bookof2.Config.WeightSymbolNum),
				zap.Error(err))

			return err
		}

		bookof2.Config.WeightSymbolNumVW = vw2
	}

	if bookof2.Config.WeightSymbol != "" {
		vw2, err := pool.LoadSymbolWeights(bookof2.Config.WeightSymbol, "val", "weight", pool.DefaultPaytables, bookof2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BookOf2.InitEx:LoadSymbolWeights",
				zap.String("Weight", bookof2.Config.WeightSymbol),
				zap.Error(err))

			return err
		}

		bookof2.Config.WeightSymbolVW = vw2
	}

	for _, v := range bookof2.Config.WildSymbols {
		bookof2.Config.WildSymbolCodes = append(bookof2.Config.WildSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	if len(bookof2.Config.MapSymbols) > 0 {
		bookof2.Config.MapSymbolCodes = make(map[int][]int)

		for strs, arrs := range bookof2.Config.MapSymbols {
			arrsc := []int{}
			k, isok := pool.DefaultPaytables.MapSymbols[strs]
			if !isok {
				goutils.Error("BookOf2.InitEx:mapSymbols",
					zap.String("symbol", strs),
					zap.Error(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			if len(arrs) != len(bookof2.Config.WildSymbols)+1 {
				goutils.Error("BookOf2.InitEx:mapSymbols",
					zap.String("symbol", strs),
					zap.Error(ErrIvalidComponentConfig))

				return ErrIvalidComponentConfig
			}

			for _, cs := range arrs {
				ck, isok := pool.DefaultPaytables.MapSymbols[cs]
				if !isok {
					goutils.Error("BookOf2.InitEx:mapSymbols",
						zap.String("symbol", strs),
						zap.Error(ErrInvalidSymbol))

					return ErrInvalidSymbol
				}

				arrsc = append(arrsc, ck)
			}

			bookof2.Config.MapSymbolCodes[k] = arrsc
		}
	}

	bookof2.onInit(&bookof2.Config.BasicComponentConfig)

	return nil
}

// playgame
func (bookof2 *BookOf2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	bookof2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[bookof2.Name].(*BookOf2Data)

	isTrigger := bookof2.Config.ForceTrigger

	if !isTrigger && bookof2.Config.WeightTriggerVW != nil {
		iv, err := bookof2.Config.WeightTriggerVW.RandVal(plugin)
		if err != nil {
			goutils.Error("bookof2.OnPlayGame:WeightTrigger.RandVal",
				zap.Error(err))

			return err
		}

		if iv.Int() != 0 {
			isTrigger = true
		}
	}

	if isTrigger {
		if bookof2.Config.SymbolCollection != "" {
			sccd := gameProp.MapComponentData[bookof2.Config.SymbolCollection].(*SymbolCollectionData)
			if sccd == nil {
				goutils.Error("BookOf2.OnPlayGame",
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
			symbolNum := bookof2.Config.ForceSymbolNum

			if symbolNum <= 0 && bookof2.Config.WeightSymbolNumVW != nil {
				iv, err := bookof2.Config.WeightSymbolNumVW.RandVal(plugin)
				if err != nil {
					goutils.Error("bookof2.OnPlayGame:WeightSymbolNum.RandVal",
						zap.Error(err))

					return err
				}

				symbolNum = iv.Int()
			}

			if bookof2.Config.ForceSymbolNum == 1 && bookof2.Config.SymbolRNG != "" {
				rng := gameProp.GetTagInt(bookof2.Config.SymbolRNG)
				cs := bookof2.Config.WeightSymbolVW.Vals[rng]

				cd.Symbols = append(cd.Symbols, cs.Int())
			} else {
				curWeight := bookof2.Config.WeightSymbolVW.Clone()

				for i := 0; i < symbolNum; i++ {
					cs, err := curWeight.RandVal(plugin)
					if err != nil {
						goutils.Error("bookof2.OnPlayGame:curWeight.RandVal",
							zap.Error(err))

						return err
					}

					cd.Symbols = append(cd.Symbols, cs.Int())

					err = curWeight.RemoveVal(cs)
					if err != nil {
						goutils.Error("bookof2.OnPlayGame:curWeight.RemoveVal",
							zap.Error(err))

						return err
					}
				}
			}
		}

		gs := bookof2.GetTargetScene3(gameProp, curpr, prs, &cd.BasicComponentData, bookof2.Name, "", 0)

		bookof2.procInitScene(gameProp, gs, curpr, cd)

		for _, s := range cd.Symbols {
			ngs, err := bookof2.procBookOfScene(gameProp, gs, s)
			if err != nil {
				goutils.Error("bookof2.OnPlayGame:procBookOfScene",
					zap.Error(err))

				return err
			}

			bookof2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

			scr := sgc7game.CalcScatter3(gs, gameProp.CurPaytables, s, gameProp.GetBet(stake, bookof2.Config.BetType)*gameProp.GetVal(GamePropCurLineNum), 1,
				func(scatter int, cursymbol int) bool {
					return cursymbol == s || goutils.IndexOfIntSlice(bookof2.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, true)
			if scr != nil {
				bookof2.AddResult(curpr, scr, &cd.BasicComponentData)
			}
		}

		bookof2.onStepEnd(gameProp, curpr, gp, "")
	} else {
		bookof2.onStepEnd(gameProp, curpr, gp, "")
	}

	return nil
}

// procBookOfScene - outpur to asciigame
func (bookof2 *BookOf2) procInitScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, cd *BookOf2Data) {
	if len(bookof2.Config.MapSymbolCodes) > 0 {
		ngs := gs.CloneEx(gameProp.PoolScene)

		for x, arr := range gs.Arr {
			for y, s := range arr {
				ci := goutils.IndexOfIntSlice(cd.Symbols, s, 0)
				if ci >= 0 {
					ngs.Arr[x][y] = bookof2.Config.MapSymbolCodes[cd.Symbols[ci]][0]
				}
			}
		}

		bookof2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	}
}

// procBookOfScene - outpur to asciigame
func (bookof2 *BookOf2) procBookOfScene(gameProp *GameProperty, gs *sgc7game.GameScene, symbol int) (*sgc7game.GameScene, error) {
	ngs := gs.CloneEx(gameProp.PoolScene)

	if len(bookof2.Config.MapSymbolCodes) > 0 {
		for x, arr := range gs.Arr {
			hass := false
			wi := -1
			for _, s := range arr {
				ci := goutils.IndexOfIntSlice(bookof2.Config.WildSymbolCodes, s, 0)
				if ci >= 0 {
					if wi > ci {
						wi = ci
					}
				} else if s == symbol {
					hass = true
				}
			}

			if wi >= 0 {
				for y := range arr {
					ngs.Arr[x][y] = bookof2.Config.MapSymbolCodes[symbol][wi+1]
				}
			}
			if hass {
				for y := range arr {
					ngs.Arr[x][y] = bookof2.Config.MapSymbolCodes[symbol][0]
				}
			}
		}
	} else {
		for x, arr := range gs.Arr {
			hass := false
			for _, s := range arr {
				if s == symbol || goutils.IndexOfIntSlice(bookof2.Config.WildSymbolCodes, s, 0) >= 0 {
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
	}

	return ngs, nil
}

// OnAsciiGame - outpur to asciigame
func (bookof2 *BookOf2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[bookof2.Name].(*BookOfData)

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

		asciigame.OutputResults(fmt.Sprintf("%v wins", bookof2.Name), pr, func(i int, ret *sgc7game.Result) bool {
			return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
		}, mapSymbolColor)
	}

	return nil
}

// OnStats
func (bookof2 *BookOf2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	wins := int64(0)
	isTrigger := false

	for _, v := range lst {
		gp, isok := v.CurGameModParams.(*GameParams)
		if isok {
			curComponent, isok := gp.MapComponentMsgs[bookof2.Name]
			if isok {
				curwins, err := bookof2.OnStatsWithPB(feature, curComponent, v)
				if err != nil {
					goutils.Error("BookOf2.OnStats",
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
func (bookof2 *BookOf2) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.BookOfData)
	if !isok {
		goutils.Error("BookOf2.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return bookof2.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (bookof2 *BookOf2) NewComponentData() IComponentData {
	return &BookOf2Data{}
}

// EachUsedResults -
func (bookof2 *BookOf2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.BookOfData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BookOf2.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewBookOf2(name string) IComponent {
	return &BookOf2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
