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
			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:LinesTriggerData",
				zap.Error(err))

			return 0, false
		}

		if val == "wins" {
			return int(msg.Wins), true
		} else if val == "symbolNum" {
			return int(msg.SymbolNum), true
		}
	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.ScatterTriggerData" {
		var msg sgc7pb.ScatterTriggerData

		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
		if err != nil {
			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:ScatterTriggerData",
				zap.Error(err))

			return 0, false
		}

		if val == "wins" {
			return int(msg.Wins), true
		} else if val == "symbolNum" {
			return int(msg.SymbolNum), true
		}
	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.WaysTriggerData" {
		var msg sgc7pb.WaysTriggerData

		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
		if err != nil {
			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:WaysTriggerData",
				zap.Error(err))

			return 0, false
		}

		if val == "wins" {
			return int(msg.Wins), true
		} else if val == "symbolNum" {
			return int(msg.SymbolNum), true
		}
	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.ClusterTriggerData" {
		var msg sgc7pb.ClusterTriggerData

		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
		if err != nil {
			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:ClusterTriggerData",
				zap.Error(err))

			return 0, false
		}

		if val == "wins" {
			return int(msg.Wins), true
		} else if val == "symbolNum" {
			return int(msg.SymbolNum), true
		}
	}

	return 0, false
}
