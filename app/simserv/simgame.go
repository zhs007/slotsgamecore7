package main

import (
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/lowcode"
	"github.com/zhs007/slotsgamecore7/simserv"
	"go.uber.org/zap"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// SimService - simservice
type SimService struct {
	*simserv.BasicService
}

// NewSimService - new a SimService
func NewSimService(game sgc7game.IGame) (simserv.IService, error) {
	bs, err := simserv.NewBasicService(game)
	if err != nil {
		goutils.Error("NewSimService:NewBasicService",
			zap.Error(err))
	}

	return &SimService{bs}, nil
}

// BuildPBGameModParam - any -> *any.Any
func (sv *SimService) BuildPBGameModParam(gp any) (*anypb.Any, error) {
	mygp, isok := gp.(*lowcode.GameParams)
	if !isok {
		return nil, ErrInvalidGameParams
	}

	return anypb.New(&mygp.GameParam)
}

// BuildPBGameModParamFromAny - any -> *any.Any
func (sv *SimService) BuildPBGameModParamFromAny(msg *anypb.Any) (any, error) {
	mygp := &lowcode.GameParams{}

	err := msg.UnmarshalTo(&mygp.GameParam)
	if err != nil {
		return nil, err
	}

	return mygp, nil
}
