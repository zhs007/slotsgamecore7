package lowcode

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
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
				goutils.Err(err))

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
				goutils.Err(err))

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
				goutils.Err(err))

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
				goutils.Err(err))

			return 0, false
		}

		if val == "wins" {
			return int(msg.Wins), true
		} else if val == "symbolNum" {
			return int(msg.SymbolNum), true
		}
	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.RespinData" {
		var msg sgc7pb.RespinData

		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
		if err != nil {
			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:RespinData",
				goutils.Err(err))

			return 0, false
		}

		if val == "curRespinNum" {
			return int(msg.CurRespinNum), true
		}
	}

	return 0, false
}

// func GetComponentDataVal2(icd IComponentData, val string) (int, bool) {
// 	pbany, isok := pb.(*anypb.Any)
// 	if !isok {
// 		return 0, false
// 	}

// 	if pbany.TypeUrl == "type.googleapis.com/sgc7pb.LinesTriggerData" {
// 		var msg sgc7pb.LinesTriggerData

// 		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
// 		if err != nil {
// 			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:LinesTriggerData",
// 				goutils.Err(err))

// 			return 0, false
// 		}

// 		if val == "wins" {
// 			return int(msg.Wins), true
// 		} else if val == "symbolNum" {
// 			return int(msg.SymbolNum), true
// 		}
// 	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.ScatterTriggerData" {
// 		var msg sgc7pb.ScatterTriggerData

// 		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
// 		if err != nil {
// 			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:ScatterTriggerData",
// 				goutils.Err(err))

// 			return 0, false
// 		}

// 		if val == "wins" {
// 			return int(msg.Wins), true
// 		} else if val == "symbolNum" {
// 			return int(msg.SymbolNum), true
// 		}
// 	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.WaysTriggerData" {
// 		var msg sgc7pb.WaysTriggerData

// 		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
// 		if err != nil {
// 			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:WaysTriggerData",
// 				goutils.Err(err))

// 			return 0, false
// 		}

// 		if val == "wins" {
// 			return int(msg.Wins), true
// 		} else if val == "symbolNum" {
// 			return int(msg.SymbolNum), true
// 		}
// 	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.ClusterTriggerData" {
// 		var msg sgc7pb.ClusterTriggerData

// 		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
// 		if err != nil {
// 			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:ClusterTriggerData",
// 				goutils.Err(err))

// 			return 0, false
// 		}

// 		if val == "wins" {
// 			return int(msg.Wins), true
// 		} else if val == "symbolNum" {
// 			return int(msg.SymbolNum), true
// 		}
// 	} else if pbany.TypeUrl == "type.googleapis.com/sgc7pb.RespinData" {
// 		var msg sgc7pb.RespinData

// 		err := anypb.UnmarshalTo(pbany, &msg, proto.UnmarshalOptions{})
// 		if err != nil {
// 			goutils.Error("GetComponentDataVal:anypb.UnmarshalTo:RespinData",
// 				goutils.Err(err))

// 			return 0, false
// 		}

// 		if val == "curRespinNum" {
// 			return int(msg.CurRespinNum), true
// 		}
// 	}

// 	return 0, false
// }
