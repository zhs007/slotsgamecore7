package lowcode

import (
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
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const SymbolValWinsTypeName = "symbolValWins"

type SymbolValWinsType int

const (
	svwTypeNormal        SymbolValWinsType = 0
	svwTypeCollector     SymbolValWinsType = 1
	svwTypeReelCollector SymbolValWinsType = 2
)

func parseSymbolValWinsType(strType string) SymbolValWinsType {
	switch strType {
	case "collector":
		return svwTypeCollector
	case "reelcollector":
		return svwTypeReelCollector
	}

	return svwTypeNormal
}

const (
	SVWDVWins      string = "wins"      // 中奖的数值，线注的倍数
	SVWDVSymbolNum string = "symbolNum" // 符号数量
)

type SymbolValWinsData struct {
	BasicComponentData
	SymbolNum int
	Wins      int
}

// OnNewGame -
func (symbolValWinsData *SymbolValWinsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	symbolValWinsData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (symbolValWinsData *SymbolValWinsData) onNewStep() {
	symbolValWinsData.UsedResults = nil
	symbolValWinsData.SymbolNum = 0
	symbolValWinsData.Wins = 0
}

// Clone
func (symbolValWinsData *SymbolValWinsData) Clone() IComponentData {
	target := &SymbolValWinsData{
		BasicComponentData: symbolValWinsData.CloneBasicComponentData(),
		SymbolNum:          symbolValWinsData.SymbolNum,
		Wins:               symbolValWinsData.Wins,
	}

	return target
}

// BuildPBComponentData
func (symbolValWinsData *SymbolValWinsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolValWinsData{
		BasicComponentData: symbolValWinsData.BuildPBBasicComponentData(),
	}

	if !gIsReleaseMode {
		pbcd.SymbolNum = int32(symbolValWinsData.SymbolNum)
		pbcd.Wins = int32(symbolValWinsData.Wins)
	}

	return pbcd
}

// GetValEx -
func (symbolValWinsData *SymbolValWinsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	switch key {
	case SVWDVSymbolNum:
		return symbolValWinsData.SymbolNum, true
	case SVWDVWins:
		return symbolValWinsData.Wins, true
	case CVResultNum, CVWinResultNum:
		return len(symbolValWinsData.UsedResults), true
	}

	return 0, false
}

// SymbolValWinsConfig - configuration for SymbolValWins
type SymbolValWinsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetTypeString        string            `yaml:"betType" json:"betType"`   // bet or totalBet or noPay
	BetType              BetType           `yaml:"-" json:"-"`               // bet or totalBet or noPay
	WinMulti             int               `yaml:"winMulti" json:"winMulti"` // bet or totalBet
	Symbols              []string          `yaml:"symbols" json:"symbols"`   // like collect
	SymbolCodes          []int             `yaml:"-" json:"-"`               //
	StrType              string            `yaml:"type" json:"type"`
	Type                 SymbolValWinsType `yaml:"-" json:"-"`
	CoinSymbols          []string          `yaml:"coinSymbols" json:"coinSymbols"` // coin symbols
	CoinSymbolCodes      []int             `yaml:"-" json:"-"`                     // coin symbols
}

// SetLinkComponent
func (cfg *SymbolValWinsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type SymbolValWins struct {
	*BasicComponent `json:"-"`
	Config          *SymbolValWinsConfig `json:"config"`
}

// Init -
func (svw *SymbolValWins) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolValWins.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolValWinsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolValWins.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return svw.InitEx(cfg, pool)
}

// InitEx -

func (svw *SymbolValWins) InitEx(cfg any, pool *GamePropertyPool) error {
	svw.Config = cfg.(*SymbolValWinsConfig)
	svw.Config.ComponentType = SymbolValWinsTypeName

	svw.Config.BetType = ParseBetType(svw.Config.BetTypeString)
	svw.Config.Type = parseSymbolValWinsType(svw.Config.StrType)

	for _, s := range svw.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolValWins.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		svw.Config.SymbolCodes = append(svw.Config.SymbolCodes, sc)
	}

	for _, s := range svw.Config.CoinSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolValWins.InitEx:CoinSymbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		svw.Config.CoinSymbolCodes = append(svw.Config.CoinSymbolCodes, sc)
	}

	svw.onInit(&svw.Config.BasicComponentConfig)

	return nil
}

// playgame
func (svw *SymbolValWins) OnPlayGame(gameProp *GameProperty, pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	svwd := icd.(*SymbolValWinsData)
	svwd.onNewStep()

	gs := svw.GetTargetScene3(gameProp, pr, prs, 0)
	if gs == nil {
		goutils.Error("SymbolValWins.OnPlayGame:GetTargetScene3",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	os := svw.GetTargetOtherScene3(gameProp, pr, prs, 0)

	if os != nil {
		collectorpos := []int{}
		mul := 0
		switch svw.Config.Type {
		case svwTypeCollector:
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(svw.Config.SymbolCodes, s, 0) >= 0 {
						mul++

						collectorpos = append(collectorpos, x, y)
					}
				}
			}
		case svwTypeReelCollector:
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(svw.Config.SymbolCodes, s, 0) >= 0 {
						mul++

						collectorpos = append(collectorpos, x, y)

						break
					}
				}
			}
		default:
			mul = 1
		}

		totalvals := 0
		pos := make([]int, 0, len(os.Arr)*len(os.Arr[0])*2)

		if len(svw.Config.CoinSymbolCodes) > 0 {
			for x := 0; x < len(os.Arr); x++ {
				for y := 0; y < len(os.Arr[x]); y++ {
					if slices.Contains(svw.Config.CoinSymbolCodes, gs.Arr[x][y]) && os.Arr[x][y] > 0 {
						totalvals += os.Arr[x][y]
						pos = append(pos, x, y)

						svwd.SymbolNum++
					}
				}
			}
		} else {
			for x := 0; x < len(os.Arr); x++ {
				for y := 0; y < len(os.Arr[x]); y++ {
					if os.Arr[x][y] > 0 {
						totalvals += os.Arr[x][y]
						pos = append(pos, x, y)

						svwd.SymbolNum++
					}
				}
			}
		}

		if totalvals > 0 && mul > 0 {
			bet := gameProp.GetBet3(stake, svw.Config.BetType)
			othermul := svw.GetWinMulti(&svwd.BasicComponentData)

			for i := 0; i < mul; i++ {
				newpos := make([]int, 0, len(pos)+2)

				if svw.isCollectorType() {

					newpos = append(newpos, collectorpos[i*2], collectorpos[i*2+1])
				}

				newpos = append(newpos, pos...)

				ret := &sgc7game.Result{
					Type:       sgc7game.RTCoins,
					LineIndex:  -1,
					Pos:        newpos,
					SymbolNums: len(pos) / 2,
					Mul:        1,
				}

				if svw.isCollectorType() {

					ret.Symbol = gs.Arr[newpos[0]][newpos[1]]
				}

				ret.CoinWin = totalvals * othermul
				ret.CashWin = ret.CoinWin * bet
				ret.OtherMul = othermul

				svwd.Wins += ret.CoinWin

				svw.AddResult(pr, ret, &svwd.BasicComponentData)
			}

			nc := svw.onStepEnd(gameProp, pr, gp, "")

			return nc, nil
		}
	}

	nc := svw.onStepEnd(gameProp, pr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (svw *SymbolValWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*SymbolValWinsData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

// NewComponentData -
func (svw *SymbolValWins) NewComponentData() IComponentData {
	return &SymbolValWinsData{}
}

// NewStats2 -
func (svw *SymbolValWins) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins, stats2.OptIntVal})
}

// OnStats2
func (svw *SymbolValWins) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	svw.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	svwd := icd.(*SymbolValWinsData)

	s2.ProcStatsWins(svw.Name, int64(svwd.Wins))

	multi := svw.GetWinMulti(&svwd.BasicComponentData)

	s2.ProcStatsIntVal(svw.GetName(), multi)
}

func (svw *SymbolValWins) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		return winMulti
	}

	return svw.Config.WinMulti
}

// isCollectorType reports whether the config type is a collector variant
func (svw *SymbolValWins) isCollectorType() bool {
	return svw.Config.Type == svwTypeCollector || svw.Config.Type == svwTypeReelCollector
}

func NewSymbolValWins(name string) IComponent {
	return &SymbolValWins{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "betType": "bet",
// "winMulti": 1,
// "type": "normal",
// "coinSymbols": [
//
//	"CA"
//
// ]
type jsonSymbolValWins struct {
	BetType     string   `json:"betType"`  // bet or totalBet or noPay
	WinMulti    int      `json:"winMulti"` // bet or totalBet
	Symbols     []string `json:"symbols"`  // like collect
	Type        string   `yaml:"type" json:"type"`
	CoinSymbols []string `json:"coinSymbols"` // coin symbols
}

func (jcfg *jsonSymbolValWins) build() *SymbolValWinsConfig {
	cfg := &SymbolValWinsConfig{
		BetTypeString: jcfg.BetType,
		WinMulti:      jcfg.WinMulti,
		Symbols:       jcfg.Symbols,
		StrType:       strings.ToLower(jcfg.Type),
		CoinSymbols:   jcfg.CoinSymbols,
	}

	return cfg
}

func parseSymbolValWins(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolValWins:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolValWins:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolValWins{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolValWins:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolValWinsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
