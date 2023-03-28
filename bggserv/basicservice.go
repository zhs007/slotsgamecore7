package bggserv

import (
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
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
func (bs *BasicService) BuildPlayerStateFromPB(ps sgc7game.IPlayerState, pspb *sgc7pb.PlayerState) error {
	// ips := &sgc7game.BasicPlayerState{}

	pub := &sgc7pb.BasicPlayerPublicState{}
	pri := &sgc7pb.BasicPlayerPrivateState{}

	if pspb.Public != nil {
		err := pspb.Public.UnmarshalTo(pub)
		if err != nil {
			goutils.Error("BasicService.BuildPlayerStateFromPB:Public.UnmarshalTo",
				zap.Error(err))

			return err
		}

		ps.SetPublic(&sgc7game.BasicPlayerPublicState{
			CurGameMod: pub.CurGameMod,
			NextM:      int(pub.NextM),
		})
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
