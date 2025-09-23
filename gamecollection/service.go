package gamecollection

import (
	"github.com/zhs007/slotsgamecore7/grpcserv"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/types/known/anypb"
)

// Service - service
type Service struct {
	*grpcserv.BasicService2
}

// BuildPBGameModParam - interface{} -> *anypb.Any
func (sv *Service) BuildPBGameModParam(gp interface{}) (*anypb.Any, error) {
	mygp, isok := gp.(*lowcode.GameParams)
	if !isok {
		return nil, ErrInvalidGameParams
	}

	return anypb.New(&mygp.GameParam)
}

// BuildPBGameModParamFromAny - interface{} -> *anypb.Any
func (sv *Service) BuildPBGameModParamFromAny(msg *anypb.Any) (interface{}, error) {
	mygp := &sgc7pb.GameParam{}

	err := msg.UnmarshalTo(mygp)
	if err != nil {
		return nil, err
	}

	return mygp, nil
}

func NewService() *Service {
	return &Service{}
}
