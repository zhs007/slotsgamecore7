package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"google.golang.org/protobuf/types/known/anypb"
)

// findLastPBComponentData
func findLastPBComponentData(lst []*sgc7game.PlayResult, componentName string) (*anypb.Any, *sgc7game.PlayResult) {
	for i := len(lst) - 1; i >= 0; i-- {
		pr := lst[i]

		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbcd := gp.MapComponents[componentName]
			if pbcd != nil {
				return pbcd, pr
			}
		}
	}

	return nil, nil
}

// findFirstPBComponentData
func findFirstPBComponentData(lst []*sgc7game.PlayResult, componentName string) (*anypb.Any, *sgc7game.PlayResult) {
	for _, pr := range lst {
		gp := pr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbcd := gp.MapComponents[componentName]
			if pbcd != nil {
				return pbcd, pr
			}
		}
	}

	return nil, nil
}
