package dtserv

import (
	"context"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// IService - service
type IService interface {
	// GetConfig - get configuration
	GetConfig() *sgc7game.Config
	// Initialize - initialize a player
	Initialize() sgc7game.IPlayerState
	// Play - play game
	Play(ctx context.Context, req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error)

	// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
	BuildPlayerStateFromPB(ps *sgc7pb.PlayerState) sgc7game.IPlayerState
	// BuildPlayerStatePB - sgc7game.IPlayerState -> *sgc7pb.PlayerState
	BuildPlayerStatePB(ps sgc7game.IPlayerState) *sgc7pb.PlayerState
}
