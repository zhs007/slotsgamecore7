package simserv

import (
	"github.com/golang/protobuf/ptypes/any"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// IService - service
type IService interface {
	// GetGame - get game
	GetGame() sgc7game.IGame
	// GetConfig - get configuration
	GetConfig() *sgc7game.Config
	// Initialize - initialize a player
	Initialize() sgc7game.IPlayerState
	// Play - play game
	Play(params *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error)
	// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
	BuildPlayerStateFromPB(ps sgc7game.IPlayerState, pspb *sgc7pb.PlayerState) error
	// BuildPBPlayerState - sgc7game.IPlayerState -> *sgc7pb.PlayerState
	BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error)
	// BuildPBGameModParam - interface{} -> *any.Any
	BuildPBGameModParam(gp interface{}) (*any.Any, error)
	// BuildPBGameModParamFromAny - interface{} -> *any.Any
	BuildPBGameModParamFromAny(msg *any.Any) (interface{}, error)
}
