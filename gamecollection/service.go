package gamecollection

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/zhs007/slotsgamecore7/grpcserv"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/types/known/anypb"
)

// Service - service
type Service struct {
	*grpcserv.BasicService2
}

// BuildPBGameModParam - interface{} -> *any.Any
func (sv *Service) BuildPBGameModParam(gp interface{}) (*any.Any, error) {
	mygp, isok := gp.(*lowcode.GameParams)
	if !isok {
		return nil, ErrInvalidGameParams
	}

	return anypb.New(&mygp.GameParam)
}

// BuildPBGameModParamFromAny - interface{} -> *any.Any
func (sv *Service) BuildPBGameModParamFromAny(msg *any.Any) (interface{}, error) {
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
