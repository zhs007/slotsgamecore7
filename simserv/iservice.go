package simserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// IService - service
type IService interface {
	// Config - get configuration
	Config() *sgc7game.Config
	// Initialize - initialize a player
	Initialize() sgc7game.IPlayerState
	// Play - play game
	Play(params *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error)
}
