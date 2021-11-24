package simserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// BasicService - basic service
type BasicService struct {
	game sgc7game.IGame
}

// NewBasicService - new a BasicService
func NewBasicService(game sgc7game.IGame) (*BasicService, error) {
	return &BasicService{
		game: game,
	}, nil
}

// GetGame - get game
func (serv *BasicService) GetGame() sgc7game.IGame {
	return serv.game
}

// GetConfig - get configuration
func (serv *BasicService) GetConfig() *sgc7game.Config {
	return serv.game.GetConfig()
}

// Initialize - initialize a player
func (serv *BasicService) Initialize() sgc7game.IPlayerState {
	return serv.game.Initialize()
}

// // Play - play game
// func (serv *BasicService) Play(params *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {

// }
