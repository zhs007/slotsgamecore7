package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"google.golang.org/protobuf/proto"
)

// findLastPBComponentData
func findLastPBComponentData(lst []*sgc7game.PlayResult, componentName string) (proto.Message, *sgc7game.PlayResult) {
	for i := len(lst) - 1; i >= 0; i-- {
		pr := lst[i]

		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbcd := gp.MapComponentMsgs[componentName]
			if pbcd != nil {
				return pbcd, pr
			}
		}
	}

	return nil, nil
}

// findLastPBComponentDataEx
func findLastPBComponentDataEx(lst []*sgc7game.PlayResult, respinComponentName string, componentName string) (proto.Message, *sgc7game.PlayResult) {
	for i := len(lst) - 1; i >= 0; i-- {
		pr := lst[i]

		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbRespin := gp.MapComponentMsgs[respinComponentName]
			pbcd := gp.MapComponentMsgs[componentName]
			if pbRespin != nil && pbcd != nil {
				return pbcd, pr
			}
		}
	}

	return nil, nil
}

// findFirstPBComponentData
func findFirstPBComponentData(lst []*sgc7game.PlayResult, componentName string) (proto.Message, *sgc7game.PlayResult) {
	for _, pr := range lst {
		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbcd := gp.MapComponentMsgs[componentName]
			if pbcd != nil {
				return pbcd, pr
			}
		}
	}

	return nil, nil
}

// findFirstPBComponentDataEx
func findFirstPBComponentDataEx(lst []*sgc7game.PlayResult, respinComponentName string, componentName string) (proto.Message, *sgc7game.PlayResult) {
	for _, pr := range lst {
		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbRespin := gp.MapComponentMsgs[respinComponentName]
			pbcd := gp.MapComponentMsgs[componentName]
			if pbRespin != nil && pbcd != nil {
				return pbcd, pr
			}
		}
	}

	return nil, nil
}

// findAllPBComponentDataEx
func findAllPBComponentDataEx(lst []*sgc7game.PlayResult, respinComponentName string, componentName string) ([]proto.Message, []*sgc7game.PlayResult) {
	pbs := []proto.Message{}
	prs := []*sgc7game.PlayResult{}

	for _, pr := range lst {
		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbRespin := gp.MapComponentMsgs[respinComponentName]
			pbcd := gp.MapComponentMsgs[componentName]
			if pbRespin != nil && pbcd != nil {
				pbs = append(pbs, pbcd)
				prs = append(prs, pr)
			}
		}
	}

	if len(pbs) == 0 {
		return nil, nil
	}

	return pbs, prs
}

func calcTotalCashWins(lst []*sgc7game.PlayResult) int64 {
	wins := int64(0)

	for _, v := range lst {
		wins += v.CashWin
	}

	return wins
}
