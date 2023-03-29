package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"google.golang.org/protobuf/proto"
)

type FuncOnEachUsedResult func(*sgc7game.Result)

type IComponentData interface {
	// OnNewGame -
	OnNewGame()
	// OnNewStep -
	OnNewStep()
	// BuildPBComponentData
	BuildPBComponentData() proto.Message
}
