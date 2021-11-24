package simserv

import (
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
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

// // GetConfig - get configuration
// func (serv *BasicService) GetConfig() *sgc7game.Config {
// 	return serv.game.GetConfig()
// }

// // Initialize - initialize a player
// func (serv *BasicService) Initialize() sgc7game.IPlayerState {
// 	return serv.game.Initialize()
// }

// // Play - play game
// func (serv *BasicService) Play(params *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {

// }

// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
func (bs *BasicService) BuildPlayerStateFromPB(ps sgc7game.IPlayerState, pspb *sgc7pb.PlayerState) error {
	pub := &sgc7pb.BasicPlayerPublicState2{}
	pri := &sgc7pb.BasicPlayerPrivateState2{}

	if pspb.Public != nil {
		err := pspb.Public.UnmarshalTo(pub)
		if err != nil {
			goutils.Error("BasicService.BuildPlayerStateFromPB:Public.UnmarshalTo",
				zap.Error(err))

			return err
		}

		ps.SetPublicJson(pub.Json)
	}

	if pspb.Private != nil {
		err := pspb.Private.UnmarshalTo(pri)
		if err != nil {
			goutils.Error("BasicService.BuildPlayerStateFromPB:Private.UnmarshalTo",
				zap.Error(err))

			return err
		}

		ps.SetPrivate(&sgc7game.BasicPlayerPrivateState{})
	}

	return nil
}

// BuildPBPlayerState - sgc7game.IPlayerState -> *sgc7pb.PlayerState
func (bs *BasicService) BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error) {
	pub := &sgc7pb.BasicPlayerPublicState2{
		Json: ps.GetPublicJson(),
	}
	pri := &sgc7pb.BasicPlayerPrivateState2{
		Json: ps.GetPrivateJson(),
	}

	pbpub, err := anypb.New(pub)
	if err != nil {
		goutils.Error("BasicService.BuildPBPlayerState:New(pub)",
			zap.Error(err))

		return nil, err
	}

	pbpri, err := anypb.New(pri)
	if err != nil {
		goutils.Error("BasicService.BuildPBPlayerState:New(pri)",
			zap.Error(err))

		return nil, err
	}

	return &sgc7pb.PlayerState{
		Public:  pbpub,
		Private: pbpri,
	}, nil
}
