package lowcode

import (
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7rtp "github.com/zhs007/slotsgamecore7/rtp"
	"github.com/zhs007/slotsgamecore7/stats2"
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

// func buildRTPSymbolsData(pool *GamePropertyPool) ([]int, []int) {
// 	symbols := []int{}
// 	nums := []int{}

// 	for _, v := range pool.Config.StatsSymbolCodes {
// 		symbols = append(symbols, int(v))
// 	}

// 	for i := range pool.DefaultPaytables.MapPay[0] {
// 		nums = append(nums, i+1)
// 	}

// 	return symbols, nums
// }

// func newFuncOnGameMod(cfgGameMod *RTPSymbolModule) sgc7rtp.FuncOnResult {
// 	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult, gameData any) bool {
// 		if len(cfgGameMod.Components) == 0 {
// 			return true
// 		}

// 		gp, isok := pr.CurGameModParams.(*GameParams)
// 		if isok {
// 			for _, v := range cfgGameMod.Components {
// 				_, hasComponent := gp.MapComponentMsgs[v]
// 				if hasComponent {
// 					return true
// 				}
// 			}
// 		}

// 		return false
// 	}
// }

// func newFuncOnResult(cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncOnResult {
// 	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult, gameData any) bool {
// 		if len(cfgSymbolFeature.Components) == 0 {
// 			return true
// 		}

// 		gp, isok := pr.CurGameModParams.(*GameParams)
// 		if isok {
// 			for _, v := range cfgSymbolFeature.Components {
// 				_, hasComponent := gp.MapComponentMsgs[v]
// 				if hasComponent {
// 					return true
// 				}
// 			}
// 		}

// 		return false
// 	}
// }

// func newFuncSymbolOnResult(pool *GamePropertyPool, cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncOnResult {
// 	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult, gameData any) bool {
// 		if len(cfgSymbolFeature.Components) == 0 {

// 			for _, v := range pr.Results {
// 				if v.Symbol == node.Symbol {
// 					node.TriggerNums++
// 					node.TotalWin += int64(v.CashWin)
// 				}
// 			}

// 			return true
// 		}

// 		ismine := false

// 		gp, isok := pr.CurGameModParams.(*GameParams)
// 		if isok {
// 			for _, componentName := range cfgSymbolFeature.Components {
// 				c, hasComponent := gp.MapComponents[componentName]
// 				if hasComponent {
// 					components := pool.mapComponents[pool.Stats.Bet]
// 					component := components.MapComponents[componentName]

// 					component.EachUsedResults(pr, c, func(ret *sgc7game.Result) {
// 						if ret.Symbol == node.Symbol {
// 							node.TriggerNums++
// 							node.TotalWin += int64(ret.CashWin)
// 						}
// 					})

// 					// for _, ri := range c.UsedResults {
// 					// 	ret := pr.Results[ri]

// 					// 	if ret.Symbol == node.Symbol {
// 					// 		node.TriggerNums++
// 					// 		node.TotalWin += int64(ret.CashWin)
// 					// 	}
// 					// }

// 					ismine = true
// 				}
// 			}
// 		}

// 		return ismine
// 	}
// }

// func newFuncSymbolNumOnResult(pool *GamePropertyPool, cfgSymbolFeature *RTPSymbolFeature) sgc7rtp.FuncOnResult {
// 	return func(node *sgc7rtp.RTPNode, pr *sgc7game.PlayResult, gameData any) bool {
// 		if len(cfgSymbolFeature.Components) == 0 {

// 			for _, v := range pr.Results {
// 				if v.Symbol == node.Symbol && v.SymbolNums == node.SymbolNums {
// 					node.TriggerNums++
// 					node.TotalWin += int64(v.CashWin)
// 				}
// 			}

// 			return true
// 		}

// 		ismine := false

// 		gp, isok := pr.CurGameModParams.(*GameParams)
// 		if isok {
// 			for _, componentName := range cfgSymbolFeature.Components {
// 				c, hasComponent := gp.MapComponents[componentName]
// 				if hasComponent {
// 					components := pool.mapComponents[pool.Stats.Bet]
// 					component := components.MapComponents[componentName]

// 					component.EachUsedResults(pr, c, func(ret *sgc7game.Result) {
// 						if ret.Symbol == node.Symbol && ret.SymbolNums == node.SymbolNums {
// 							node.TriggerNums++
// 							node.TotalWin += int64(ret.CashWin)
// 						}
// 					})

// 					// for _, ri := range c.UsedResults {
// 					// 	ret := pr.Results[ri]

// 					// 	if ret.Symbol == node.Symbol && ret.SymbolNums == node.SymbolNums {
// 					// 		node.TriggerNums++
// 					// 		node.TotalWin += int64(ret.CashWin)
// 					// 	}
// 					// }

// 					ismine = true
// 				}
// 			}
// 		}

// 		return ismine
// 	}
// }

// func newRTPGameModule(rtp *sgc7rtp.RTP, pool *GamePropertyPool, cfgGameModule *RTPSymbolModule) *sgc7rtp.RTPNode {
// 	gm := sgc7rtp.NewRTPGameModEx(cfgGameModule.Name, newFuncOnGameMod(cfgGameModule))

// 	symbols, nums := buildRTPSymbolsData(pool)
// 	names := []string{}
// 	funcOnResults := []sgc7rtp.FuncOnResult{}
// 	funcSymbolOnResults := []sgc7rtp.FuncOnResult{}
// 	funcSymbolNumOnResults := []sgc7rtp.FuncOnResult{}

// 	for _, v := range cfgGameModule.Features {
// 		feature := v

// 		names = append(names, v.Name)
// 		funcOnResults = append(funcOnResults, newFuncOnResult(feature))
// 		funcSymbolOnResults = append(funcSymbolOnResults, newFuncSymbolOnResult(pool, feature))
// 		funcSymbolNumOnResults = append(funcSymbolNumOnResults, newFuncSymbolNumOnResult(pool, feature))
// 	}

// 	sgc7rtp.InitGameMod3(gm, names, funcOnResults,
// 		symbols, nums,
// 		funcSymbolOnResults,
// 		funcSymbolNumOnResults)

// 	rtp.Root.AddChild(cfgGameModule.Name, gm)

// 	return gm
// }

// func hasComponent(i int, prs []*sgc7game.PlayResult, component string) bool {
// 	gp, isok := prs[i].CurGameModParams.(*GameParams)
// 	if isok {
// 		_, hasComponent := gp.MapComponentMsgs[component]
// 		if hasComponent {
// 			return true
// 		}
// 	}

// 	return false
// }

// func hasComponentEx(i int, prs []*sgc7game.PlayResult, components []string) bool {
// 	gp, isok := prs[i].CurGameModParams.(*GameParams)
// 	if isok {
// 		for _, v := range components {
// 			_, hasComponent := gp.MapComponentMsgs[v]
// 			if hasComponent {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }

// func newFuncHitRate(cfgHitRateFeature *RTPHitRateFeature) sgc7rtp.FuncHROnResult {
// 	return func(rtp *sgc7rtp.RTP, node *sgc7rtp.HitRateNode, i int, prs []*sgc7game.PlayResult) bool {
// 		if len(cfgHitRateFeature.Components) == 0 {
// 			return true
// 		}

// 		if hasComponentEx(i, prs, cfgHitRateFeature.Components) {
// 			node.TotalNums++
// 		} else {
// 			if i < len(prs)-1 && hasComponentEx(i+1, prs, cfgHitRateFeature.Components) {
// 				node.TriggerNums++

// 				return true
// 			}
// 		}

// 		return false
// 	}
// }

// func procHitRate(rtp *sgc7rtp.RTP, pool *GamePropertyPool, cfgHitRateFeature *RTPHitRateFeature) {
// 	rtp.AddHitRateNode(cfgHitRateFeature.Name, newFuncHitRate(cfgHitRateFeature))
// }

func StartRTP(gamecfg string, icore int, ispinnums int64, outputPath string, bet int64, funcNewRNG FuncNewRNG, wincap int64) error {
	// SetRTPMode()
	// IsStatsComponentMsg = true
	sgc7plugin.IsNoRNGCache = true

	game, err := NewGame2(gamecfg, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, funcNewRNG)
	if err != nil {
		goutils.Error("StartRTP:NewGame3",
			slog.String("gamecfg", gamecfg),
			goutils.Err(err))

		return err
	}

	rtp := sgc7rtp.NewRTP()

	if bet <= 0 {
		bet = int64(game.Pool.Config.Bets[0])
	}

	stake := &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  int64(bet),
		Currency: "EUR",
	}

	// if game.Pool.Stats != nil {
	// 	rtp.FuncRTPResults = func(lst []*sgc7game.PlayResult, gameData any) {
	// 		game.Pool.Stats.Push(stake, lst)
	// 	}
	// }

	// if game.Pool.Config.RTP != nil {
	// 	for _, m := range game.Pool.Config.RTP.Modules {
	// 		newRTPGameModule(rtp, game.Pool, m)
	// 	}

	// 	for _, hr := range game.Pool.Config.RTP.HitRateFeatures {
	// 		procHitRate(rtp, game.Pool, hr)
	// 	}
	// }

	d := sgc7rtp.StartRTP2(game, rtp, icore, ispinnums, stake, 100000, func(totalnums int64, curnums int64, curtime time.Duration) {
		goutils.Info("processing...",
			slog.Int64("total nums", totalnums),
			slog.Int64("current nums", curnums),
			slog.Duration("cost time", curtime))
	}, true, wincap)

	goutils.Info("finish.",
		slog.Int64("total nums", ispinnums),
		slog.Float64("rtp", float64(rtp.TotalWins)/float64(rtp.TotalBet)),
		slog.Duration("cost time", d))

	curtime := time.Now()

	rtp.Save2CSV(path.Join(outputPath, fmt.Sprintf("%v-%v.csv", game.Pool.Config.Name, curtime.Format("2006-01-02_15_04_05"))))

	// if game.Pool.Stats != nil {
	// 	game.Pool.Stats.Wait()

	// 	game.Pool.Stats.Root.SaveExcel(path.Join(outputPath, fmt.Sprintf("%v-stats-%v.xlsx", game.Pool.Config.Name, curtime.Format("2006-01-02_15_04_05"))))

	// 	goutils.Info("finish.",
	// 		slog.Int64("total nums", game.Pool.Stats.TotalNum))
	// }

	if gAllowStats2 {
		components := game.Pool.mapComponents[int(bet)]
		components.Stats2.WaitEnding()

		components.Stats2.SaveExcel(path.Join(outputPath, fmt.Sprintf("%v-%v-stats-%v.xlsx", game.Pool.Config.Name, bet, curtime.Format("2006-01-02_15_04_05"))))

		goutils.Info("finish.",
			slog.Int64("total nums", components.Stats2.BetTimes))
	}

	return nil
}

func StartRTPWithData(gamecfg []byte, icore int, ispinnums int64, bet int64, ontimer sgc7rtp.FuncOnRTPTimer, funcNewRNG FuncNewRNG) (*stats2.Stats, error) {
	// SetRTPMode()
	// IsStatsComponentMsg = true
	sgc7plugin.IsNoRNGCache = true

	game, err := NewGame2WithData(gamecfg, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, funcNewRNG)
	if err != nil {
		goutils.Error("StartRTPWithData:NewGame3",
			goutils.Err(err))

		return nil, err
	}

	rtp := sgc7rtp.NewRTP()

	if bet <= 0 {
		bet = int64(game.Pool.Config.Bets[0])
	}

	stake := &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  int64(bet),
		Currency: "EUR",
	}

	// if game.Pool.Stats != nil {
	// 	rtp.FuncRTPResults = func(lst []*sgc7game.PlayResult, gameData any) {
	// 		game.Pool.Stats.Push(stake, lst)
	// 	}
	// }

	// if game.Pool.Config.RTP != nil {
	// 	for _, m := range game.Pool.Config.RTP.Modules {
	// 		newRTPGameModule(rtp, game.Pool, m)
	// 	}

	// 	for _, hr := range game.Pool.Config.RTP.HitRateFeatures {
	// 		procHitRate(rtp, game.Pool, hr)
	// 	}
	// }

	d := sgc7rtp.StartRTP2(game, rtp, icore, ispinnums, stake, int(ispinnums/100), ontimer, true, 0)

	goutils.Info("finish.",
		slog.Int64("total nums", ispinnums),
		slog.Float64("rtp", float64(rtp.TotalWins)/float64(rtp.TotalBet)),
		slog.Duration("cost time", d))

	// curtime := time.Now()

	// rtp.Save2CSV(path.Join(outputPath, fmt.Sprintf("%v-%v.csv", game.Pool.Config.Name, curtime.Format("2006-01-02_15_04_05"))))

	// if game.Pool.Stats != nil {
	// 	game.Pool.Stats.Wait()

	// 	game.Pool.Stats.Root.SaveExcel(path.Join(outputPath, fmt.Sprintf("%v-stats-%v.xlsx", game.Pool.Config.Name, curtime.Format("2006-01-02_15_04_05"))))

	// 	goutils.Info("finish.",
	// 		slog.Int64("total nums", game.Pool.Stats.TotalNum))
	// }

	// if gAllowStats2 {
	components := game.Pool.mapComponents[int(bet)]
	components.Stats2.WaitEnding()

	// components.Stats2.SaveExcel(path.Join(outputPath, fmt.Sprintf("%v-%v-stats-%v.xlsx", game.Pool.Config.Name, bet, curtime.Format("2006-01-02_15_04_05"))))

	// goutils.Info("finish.",
	// 	slog.Int64("total nums", components.Stats2.BetTimes))
	// }

	return components.Stats2, nil
}
