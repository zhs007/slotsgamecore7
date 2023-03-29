package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

type GameParams struct {
	sgc7pb.GameParam `json:",inline"`
	LastScene        *sgc7game.GameScene `json:"-"`
	LastOtherScene   *sgc7game.GameScene `json:"-"`
}

func (gp *GameParams) AddComponentData(name string, cd IComponentData) error {
	pbmsg := cd.BuildPBComponentData()

	pbany, err := anypb.New(pbmsg)
	if err != nil {
		goutils.Error("GameParams.AddComponentData:New",
			zap.Error(err))

		return err
	}

	gp.MapComponents[name] = pbany

	return nil
}

// gIsForceDisableStats - disable stats
var gIsForceDisableStats bool

// SetForceDisableStats - disable stats
func SetForceDisableStats() {
	gIsForceDisableStats = true
}
