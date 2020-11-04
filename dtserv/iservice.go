package dtserv

import (
	any "github.com/golang/protobuf/ptypes/any"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// IService - service
type IService interface {
	// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
	BuildPlayerStateFromPB(ps *sgc7pb.PlayerState) (sgc7game.IPlayerState, error)
	// BuildPBPlayerState - sgc7game.IPlayerState -> *sgc7pb.PlayerState
	BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error)
	// BuildPBGameModParam - interface{} -> *any.Any
	BuildPBGameModParam(gp interface{}) (*any.Any, error)
}
