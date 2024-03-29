package lowcode

import (
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// // findLastPBComponentData
// func findLastPBComponentData(lst []*sgc7game.PlayResult, componentName string) (proto.Message, *sgc7game.PlayResult) {
// 	for i := len(lst) - 1; i >= 0; i-- {
// 		pr := lst[i]

// 		gp := pr.CurGameModParams.(*GameParams)
// 		if gp != nil {
// 			pbcd := gp.MapComponentMsgs[componentName]
// 			if pbcd != nil {
// 				return pbcd, pr
// 			}
// 		}
// 	}

// 	return nil, nil
// }

// // findLastPBComponentDataEx
// func findLastPBComponentDataEx(lst []*sgc7game.PlayResult, respinComponentName string, componentName string) (proto.Message, *sgc7game.PlayResult) {
// 	for i := len(lst) - 1; i >= 0; i-- {
// 		pr := lst[i]

// 		gp := pr.CurGameModParams.(*GameParams)
// 		if gp != nil {
// 			pbRespin := gp.MapComponentMsgs[respinComponentName]
// 			pbcd := gp.MapComponentMsgs[componentName]
// 			if pbRespin != nil && pbcd != nil {
// 				return pbcd, pr
// 			}
// 		}
// 	}

// 	return nil, nil
// }

// // findFirstPBComponentData
// func findFirstPBComponentData(lst []*sgc7game.PlayResult, componentName string) (proto.Message, *sgc7game.PlayResult) {
// 	for _, pr := range lst {
// 		gp := pr.CurGameModParams.(*GameParams)
// 		if gp != nil {
// 			pbcd := gp.MapComponentMsgs[componentName]
// 			if pbcd != nil {
// 				return pbcd, pr
// 			}
// 		}
// 	}

// 	return nil, nil
// }

// // findFirstPBComponentDataEx
// func findFirstPBComponentDataEx(lst []*sgc7game.PlayResult, respinComponentName string, componentName string) (proto.Message, *sgc7game.PlayResult) {
// 	for _, pr := range lst {
// 		gp := pr.CurGameModParams.(*GameParams)
// 		if gp != nil {
// 			pbRespin := gp.MapComponentMsgs[respinComponentName]
// 			pbcd := gp.MapComponentMsgs[componentName]
// 			if pbRespin != nil && pbcd != nil {
// 				return pbcd, pr
// 			}
// 		}
// 	}

// 	return nil, nil
// }

// // findAllPBComponentDataEx
// func findAllPBComponentDataEx(lst []*sgc7game.PlayResult, respinComponentName string, componentName string) ([]proto.Message, []*sgc7game.PlayResult) {
// 	pbs := []proto.Message{}
// 	prs := []*sgc7game.PlayResult{}

// 	for _, pr := range lst {
// 		gp := pr.CurGameModParams.(*GameParams)
// 		if gp != nil {
// 			pbRespin := gp.MapComponentMsgs[respinComponentName]
// 			pbcd := gp.MapComponentMsgs[componentName]
// 			if pbRespin != nil && pbcd != nil {
// 				pbs = append(pbs, pbcd)
// 				prs = append(prs, pr)
// 			}
// 		}
// 	}

// 	if len(pbs) == 0 {
// 		return nil, nil
// 	}

// 	return pbs, prs
// }

// func calcTotalCashWins(lst []*sgc7game.PlayResult) int64 {
// 	wins := int64(0)

// 	for _, v := range lst {
// 		wins += v.CashWin
// 	}

// 	return wins
// }

func isSameBoolSlice(src []bool, dest []bool) bool {
	if len(src) == len(dest) {
		for i, v := range src {
			if v != dest[i] {
				return false
			}
		}

		return true
	}

	return false
}

func GetExcludeSymbols(pt *sgc7game.PayTables, symbols []int) []int {
	es := []int{}

	for s := range pt.MapPay {
		if goutils.IndexOfIntSlice(symbols, s, 0) < 0 {
			es = append(es, s)
		}
	}

	return es
}

func IsInPosArea(x, y int, posArea []int) bool {
	return x >= posArea[0] && x <= posArea[1] && y >= posArea[2] && y <= posArea[3]
}

// ProcCheat -
func ProcCheat(plugin sgc7plugin.IPlugin, cheat string) (*ForceOutcome, error) {
	if cheat != "" {
		if sgc7game.IsRngString(cheat) {
			str := goutils.AppendString("[", cheat, "]")

			rngs := []int{}
			err := sonic.Unmarshal([]byte(str), &rngs)
			if err != nil {
				return nil, err
			}

			plugin.SetCache(rngs)
		} else {
			if gAllowForceOutcome {
				return ParseForceOutcome(cheat), nil
			}
		}
	}

	return nil, nil
}

func procSpin(game *Game, ips sgc7game.IPlayerState, plugin sgc7plugin.IPlugin, stake *sgc7game.Stake, cmd string, params string) ([]*sgc7game.PlayResult, error) {
	results := []*sgc7game.PlayResult{}
	gameData := game.NewGameData(stake)
	defer game.DeleteGameData(gameData)

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := game.Play(plugin, cmd, params, ips, stake, results, gameData)
		if err != nil {
			goutils.Error("Spin:Play",
				slog.Int("results", len(results)),
				goutils.Err(err))

			return nil, err
		}

		if pr == nil {
			break
		}

		results = append(results, pr)
		if pr.IsFinish {
			break
		}

		if pr.IsWait {
			break
		}

		if len(pr.NextCmds) > 0 {
			cmd = pr.NextCmds[0]
		} else {
			cmd = ""
		}

		if len(results) >= MaxStepNum {
			goutils.Error("procSpin",
				slog.Int("steps", len(results)),
				goutils.Err(ErrTooManySteps))

			return nil, ErrTooManySteps
		}
	}

	return results, nil
}

func Spin(game *Game, ips sgc7game.IPlayerState, plugin sgc7plugin.IPlugin, stake *sgc7game.Stake, cmd string, params string, cheat string) ([]*sgc7game.PlayResult, error) {
	fo, err := ProcCheat(plugin, cheat)
	if err != nil {
		goutils.Error("Spin:ProcCheat",
			slog.String("cheat", cheat),
			goutils.Err(err))

		return nil, err
	}

	err = game.CheckStake(stake)
	if err != nil {
		goutils.Error("Spin:CheckStake",
			slog.Any("stake", stake),
			goutils.Err(err))

		return nil, err
	}

	if fo == nil {
		return procSpin(game, ips, plugin, stake, cmd, params)
	}

	for tryi := 0; tryi < gMaxForceOutcomeTimes; tryi++ {
		newips := ips.Clone()
		plugin.ClearCache()
		plugin.ClearUsedRngs()

		lst, err := procSpin(game, newips, plugin, stake, cmd, params)
		if err != nil {
			goutils.Error("Spin:procSpin",
				goutils.Err(err))

			return nil, err
		}

		if fo.IsValid(lst) {
			return lst, nil
		}
	}

	goutils.Error("Spin",
		goutils.Err(ErrCannotForceOutcome))

	return nil, ErrCannotForceOutcome
}

func InsStringSliceNonRep(arr []string, str string) []string {
	if goutils.IndexOfStringSlice(arr, str, 0) >= 0 {
		return arr
	}

	return append(arr, str)
}

func InsSliceNonRep(arr []string, src []string) []string {
	for _, v := range src {
		arr = InsStringSliceNonRep(arr, v)
	}

	return arr
}
