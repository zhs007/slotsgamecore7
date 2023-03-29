package simserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/types/known/anypb"
)

// IService - service
type IService interface {
	// GetGame - get game
	GetGame() sgc7game.IGame
	// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
	BuildPlayerStateFromPB(ps sgc7game.IPlayerState, pspb *sgc7pb.PlayerState) error
	// BuildPBPlayerState - sgc7game.IPlayerState -> *sgc7pb.PlayerState
	BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error)
	// BuildPBGameModParam - any -> *any.Any
	BuildPBGameModParam(gp any) (*anypb.Any, error)
	// BuildPBGameModParamFromAny - any -> *any.Any
	BuildPBGameModParamFromAny(msg *anypb.Any) (any, error)
}
