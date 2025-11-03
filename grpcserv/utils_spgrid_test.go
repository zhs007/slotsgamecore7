package grpcserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/types/known/anypb"
)

type fakeSvc struct{}

func (s *fakeSvc) BuildPlayerStateFromPB(ps sgc7game.IPlayerState, pspb *sgc7pb.PlayerState) error {
	return nil
}

func (s *fakeSvc) BuildPBPlayerState(ps sgc7game.IPlayerState) (*sgc7pb.PlayerState, error) {
	return nil, nil
}

func (s *fakeSvc) BuildPBGameModParam(gp any) (*anypb.Any, error) {
	return nil, nil
}

func (s *fakeSvc) BuildPBGameModParamFromAny(msg *anypb.Any) (any, error) {
	return nil, nil
}

func Test_addWinResult_SPGrid(t *testing.T) {
	// build scenes
	gs1 := &sgc7game.GameScene{}
	err := gs1.InitWithArr2([][]int{{1, 2, 3}, {4, 5, 6}})
	assert.NoError(t, err)

	gs2 := &sgc7game.GameScene{}
	err = gs2.InitWithArr2([][]int{{7, 8}, {9, 10}})
	assert.NoError(t, err)

	prInternal := sgc7game.NewPlayResult("bg", 0, 0, "bg")
	prInternal.SPGrid = make(map[string][]*sgc7game.GameScene)
	prInternal.SPGrid["key1"] = []*sgc7game.GameScene{gs1, gs2}

	reply := &sgc7pb.ReplyPlay{}

	svc := &fakeSvc{}

	err = addWinResult(svc, reply, prInternal)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(reply.Results))
	gr := reply.Results[0]
	assert.NotNil(t, gr.ClientData)

	sp := gr.ClientData.SpGrid
	assert.NotNil(t, sp)
	lst, ok := sp["key1"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(lst.Scenes))
}
