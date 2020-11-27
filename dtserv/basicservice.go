package dtserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// BasicService - BasicService
type BasicService struct {
}

// NewBasicService - new BasicService
func NewBasicService() *BasicService {
	return &BasicService{}
}

// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
func (bs *BasicService) BuildPlayerStateFromPB(ps *sgc7pb.PlayerState) (sgc7game.IPlayerState, error) {
	ips := &sgc7game.BasicPlayerState{}

	pub := &sgc7pb.BasicPlayerPublicState{}
	pri := &sgc7pb.BasicPlayerPrivateState{}

	if ps.Public != nil {
		err := ps.Public.UnmarshalTo(pub)
		if err != nil {
			sgc7utils.Error("BasicService.BuildPlayerStateFromPB:Public.UnmarshalTo",
				zap.Error(err))

			return nil, err
		}

		ips.SetPublic(&sgc7game.BasicPlayerPublicState{
			CurGameMod: pub.CurGameMod,
			NextM:      int(pub.NextM),
		})
	}

	if ps.Private != nil {
		err := ps.Private.UnmarshalTo(pri)
		if err != nil {
			sgc7utils.Error("BasicService.BuildPlayerStateFromPB:Private.UnmarshalTo",
				zap.Error(err))

			return nil, err
		}

		ips.SetPrivate(&sgc7game.BasicPlayerPrivateState{})
	}

	return ips, nil
}

// BuildPBPlayerState - sgc7game.IPlayerState -> *sgc7pb.PlayerState
func (bs *BasicService) BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error) {
	curps, isok := ps.(*sgc7game.BasicPlayerState)
	if !isok {
		return nil, ErrInvalidBasicPlayerState
	}

	pub := &sgc7pb.BasicPlayerPublicState{
		CurGameMod: curps.Public.CurGameMod,
		NextM:      int32(curps.Public.NextM),
	}
	pri := &sgc7pb.BasicPlayerPrivateState{}

	pbpub, err := anypb.New(pub)
	if err != nil {
		sgc7utils.Error("BasicService.BuildPBPlayerState:New(pub)",
			zap.Error(err))

		return nil, err
	}

	pbpri, err := anypb.New(pri)
	if err != nil {
		sgc7utils.Error("BasicService.BuildPBPlayerState:New(pri)",
			zap.Error(err))

		return nil, err
	}

	return &sgc7pb.PlayerState{
		Public:  pbpub,
		Private: pbpri,
	}, nil
}
