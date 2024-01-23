package lowcode

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func GetComponentDataVal(pb proto.Message, val string) (int, bool) {
	pbany, isok := pb.(*anypb.Any)
	if !isok {
		return 0, false
	}

	if pbany.TypeUrl == "type.googleapis.com/sgc7pb.LinesTriggerData" {
		var msg sgc7pb.LinesTriggerData

		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
		if err != nil {
			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo",
				zap.Error(err))

			return 0, false
		}

		if val == "wins" {
			return int(msg.Wins), true
		}
	}

	return 0, false
}
