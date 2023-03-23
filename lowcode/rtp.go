package lowcode

import (
	"fmt"
	"path"
	"time"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7rtp "github.com/zhs007/slotsgamecore7/rtp"
	"go.uber.org/zap"
)

type RTPSymbolFeature struct {
	Name       string   `yaml:"name"`
	Components []string `yaml:"components"`
}

type RTPSymbolModule struct {
	Name       string              `yaml:"name"`
	Components []string            `yaml:"components"`
	Features   []*RTPSymbolFeature `yaml:"features"`
}

type RTPHitRateFeature struct {
	Name       string   `yaml:"name"`
	Components []string `yaml:"components"`
}

type RTPConfig struct {
	Modules         []*RTPSymbolModule   `yaml:"modules"`
	HitRateFeatures []*RTPHitRateFeature `yaml:"hitRateFeatures"`
}

func buildRTPSymbolsData(gameProp *GameProperty) ([]int, []int) {
	symbols := []int{}
	nums := []int{}

	for _, v := range gameProp.Config.StatsSymbolCodes {
		symbols = append(symbols, int(v))
	}

	for i := range gameProp.CurPaytables.MapPay[0] {
		nums = append(nums, i+1)
	}

	return symbols, nums
}

func newFuncOnGameMod(gameProp *GameProperty, cfgGameMod *RTPSymbolModule) sgc7rtp.FuncOnResult {
	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult) bool {
		if len(cfgGameMod.Components) == 0 {
			return true
		}

		gp, isok := pr.CurGameModParams.(*GameParams)
		if isok {
			for _, v := range cfgGameMod.Components {
				_, hasComponent := gp.MapComponents[v]
				if hasComponent {
					return true
				}
			}
		}

		return false
	}
}

func newFuncOnResult(gameProp *GameProperty, cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncOnResult {
	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult) bool {
		if len(cfgSymbolFeature.Components) == 0 {
			return true
		}

		gp, isok := pr.CurGameModParams.(*GameParams)
		if isok {
			for _, v := range cfgSymbolFeature.Components {
				_, hasComponent := gp.MapComponents[v]
				if hasComponent {
					return true
				}
			}
		}

		return false
	}
}

func newFuncSymbolOnResult(gameProp *GameProperty, cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncOnResult {
	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult) bool {
		if len(cfgSymbolFeature.Components) == 0 {

			for _, v := range pr.Results {
				if v.Symbol == node.Symbol {
					node.TriggerNums++
					node.TotalWin += int64(v.CashWin)
				}
			}

			return true
		}

		ismine := false

		gp, isok := pr.CurGameModParams.(*GameParams)
		if isok {
			for _, componentName := range cfgSymbolFeature.Components {
				c, hasComponent := gp.MapComponents[componentName]
				if hasComponent {
					for _, ri := range c.UsedResults {
						ret := pr.Results[ri]

						if ret.Symbol == node.Symbol {
							node.TriggerNums++
							node.TotalWin += int64(ret.CashWin)
						}
					}

					ismine = true
				}
			}
		}

		return ismine
	}
}

func newFuncSymbolNumOnResult(gameProp *GameProperty, cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncOnResult {
	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult) bool {
		if len(cfgSymbolFeature.Components) == 0 {

			for _, v := range pr.Results {
				if v.Symbol == node.Symbol && v.SymbolNums == node.SymbolNums {
					node.TriggerNums++
					node.TotalWin += int64(v.CashWin)
				}
			}

			return true
		}

		ismine := false

		gp, isok := pr.CurGameModParams.(*GameParams)
		if isok {
			for _, componentName := range cfgSymbolFeature.Components {
				c, hasComponent := gp.MapComponents[componentName]
				if hasComponent {
					for _, ri := range c.UsedResults {
						ret := pr.Results[ri]

						if ret.Symbol == node.Symbol && ret.SymbolNums == node.SymbolNums {
							node.TriggerNums++
							node.TotalWin += int64(ret.CashWin)
						}
					}

					ismine = true
				}
			}
		}

		return ismine
	}
}

// func newFuncHROnResult(gameProp *GameProperty, cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncHROnResult {
// 	return func(rtp *sgc7rtp.RTP, node *sgc7rtp.HitRateNode, pr *sgc7game.PlayResult) bool {
// 		gp, isok := pr.CurGameModParams.(*GameParams)
// 		if isok {
// 			// if gp.FGStartNum == 3 {
// 			if pr.CurGameMod == "fg" {
// 				node.TotalNums++
// 			}

// 			if pr.CurGameMod == "bg" && pr.NextGameMod == "fg" {
// 				node.TriggerNums++

// 				return true
// 			}
// 			// }
// 		}

// 		return false
// 	}
// }

func newRTPGameModule(rtp *sgc7rtp.RTP, gameProp *GameProperty, cfgGameModule *RTPSymbolModule) *sgc7rtp.RTPNode {
	gm := sgc7rtp.NewRTPGameModEx(cfgGameModule.Name, newFuncOnGameMod(gameProp, cfgGameModule))

	symbols, nums := buildRTPSymbolsData(gameProp)
	names := []string{}
	funcOnResults := []sgc7rtp.FuncOnResult{}
	funcSymbolOnResults := []sgc7rtp.FuncOnResult{}
	funcSymbolNumOnResults := []sgc7rtp.FuncOnResult{}

	for _, v := range cfgGameModule.Features {
		feature := v

		names = append(names, v.Name)
		funcOnResults = append(funcOnResults, newFuncOnResult(gameProp, feature))
		funcSymbolOnResults = append(funcSymbolOnResults, newFuncSymbolOnResult(gameProp, feature))
		funcSymbolNumOnResults = append(funcSymbolNumOnResults, newFuncSymbolNumOnResult(gameProp, feature))
	}

	sgc7rtp.InitGameMod3(gm, names, funcOnResults,
		symbols, nums,
		funcSymbolOnResults,
		funcSymbolNumOnResults)

	rtp.Root.AddChild(cfgGameModule.Name, gm)

	return gm
}

func StartRTP(gamecfg string, icore int, ispinnums int64, outputPath string) error {
	game, err := NewGame(gamecfg)
	if err != nil {
		goutils.Error("StartRTP:NewGame",
			zap.String("gamecfg", gamecfg),
			zap.Error(err))

		return err
	}

	rtp := sgc7rtp.NewRTP()

	rtp.Stats2 = game.Prop.Stats

	for _, m := range game.Prop.Config.RTP.Modules {
		newRTPGameModule(rtp, game.Prop, m)
	}

	// symbols, nums := buildRTPSymbolsData(game.Prop)

	// bg := sgc7rtp.NewRTPGameMod("bg")
	// sgc7rtp.InitGameMod3(bg, []string{"normal", "expsyms"}, []sgc7rtp.FuncOnResult{onBG, onBGExpSyms},
	// 	symbols, nums,
	// 	[]sgc7rtp.FuncOnResult{OnBGSymbolResult, OnBGSymbolResultExpSyms},
	// 	[]sgc7rtp.FuncOnResult{OnBGSymbolNumsResult, OnBGSymbolNumsResultExpSyms})
	// rtp.Root.AddChild("bg", bg)

	// fg := sgc7rtp.NewRTPGameMod("fg")
	// sgc7rtp.InitGameMod3(fg, []string{"normal", "expsyms"}, []sgc7rtp.FuncOnResult{onFG, onFGExpSyms},
	// 	[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, []int{1, 2, 3, 4, 5},
	// 	[]sgc7rtp.FuncOnResult{OnFGSymbolResult, OnFGSymbolResultExpSyms},
	// 	[]sgc7rtp.FuncOnResult{OnFGSymbolNumsResult, OnFGSymbolNumsResultExpSyms})
	// rtp.Root.AddChild("fg", fg)

	// rtp.AddHitRateNode("fg", OnFGHitRate)
	// rtp.AddHitRateNode("jackpot", OnJackpotHitRate)

	bet := game.Prop.Config.Bets[0]

	d := sgc7rtp.StartRTP(game, rtp, icore, ispinnums, &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  int64(bet),
		Currency: "EUR",
	}, 100000, func(totalnums int64, curnums int64, curtime time.Duration) {
		goutils.Info("processing...",
			zap.Int64("total nums", totalnums),
			zap.Int64("current nums", curnums),
			zap.Duration("cost time", curtime))
	}, true, 0)

	goutils.Info("finish.",
		zap.Int64("total nums", ispinnums),
		zap.Float64("rtp", float64(rtp.TotalWins)/float64(rtp.TotalBet)),
		zap.Duration("cost time", d))

	curtime := time.Now()

	rtp.Save2CSV(path.Join(outputPath, fmt.Sprintf("%v-%v.csv", game.Prop.Config.Name, curtime.Format("2006-01-02 15:04:05"))))
	rtp.Stats2.SaveExcel(path.Join(outputPath, fmt.Sprintf("%v-stats-%v.xlsx", game.Prop.Config.Name, curtime.Format("2006-01-02 15:04:05"))))

	return nil
}
