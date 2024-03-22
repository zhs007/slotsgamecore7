package grpcserv

import (
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// BasicService2 - BasicService2
type BasicService2 struct {
}

// NewBasicService - new BasicService2
func NewBasicService2() *BasicService2 {
	return &BasicService2{}
}

// BuildPlayerStateFromPB - *sgc7pb.PlayerState -> sgc7game.IPlayerState
func (bs *BasicService2) BuildPlayerStateFromPB(ps sgc7game.IPlayerState, pspb *sgc7pb.PlayerState) error {
	pub := &sgc7pb.BasicPlayerPublicState2{}
	pri := &sgc7pb.BasicPlayerPrivateState2{}

	if pspb.Public != nil {
		err := pspb.Public.UnmarshalTo(pub)
		if err != nil {
			goutils.Error("BasicService2.BuildPlayerStateFromPB:Public.UnmarshalTo",
				goutils.Err(err))

			return err
		}

		ps.SetPublicJson(pub.Json)
	}

	if pspb.Private != nil {
		err := pspb.Private.UnmarshalTo(pri)
		if err != nil {
			goutils.Error("BasicService.BuildPlayerStateFromPB:Private.UnmarshalTo",
				goutils.Err(err))

			return err
		}

		ps.SetPrivate(&sgc7game.BasicPlayerPrivateState{})
	}

	return nil
}

// BuildPBPlayerState - sgc7game.IPlayerState -> *sgc7pb.PlayerState
func (bs *BasicService2) BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error) {
	pub := &sgc7pb.BasicPlayerPublicState2{
		Json: ps.GetPublicJson(),
	}
	pri := &sgc7pb.BasicPlayerPrivateState2{
		Json: ps.GetPrivateJson(),
	}

	pbpub, err := anypb.New(pub)
	if err != nil {
		goutils.Error("BasicService2.BuildPBPlayerState:New(pub)",
			goutils.Err(err))

		return nil, err
	}

	pbpri, err := anypb.New(pri)
	if err != nil {
		goutils.Error("BasicService.BuildPBPlayerState:New(pri)",
			goutils.Err(err))

		return nil, err
	}

	return &sgc7pb.PlayerState{
		Public:  pbpub,
		Private: pbpri,
	}, nil
}
