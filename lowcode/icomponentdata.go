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
	// GetVal -
	GetVal(key string) int
	// SetVal -
	SetVal(key string, val int)
	// GetConfigVal -
	GetConfigVal(key string) string
	// SetConfigVal -
	SetConfigVal(key string, val string)
	// GetConfigIntVal -
	GetConfigIntVal(key string) (int, bool)
	// SetConfigIntVal -
	SetConfigIntVal(key string, val int)
	// ChgConfigIntVal -
	ChgConfigIntVal(key string, off int)
	// GetResults -
	GetResults() []int
}
