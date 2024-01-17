package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var IsStatsComponentMsg bool

const (
	TagCurReels = "reels"
)

const DefaultCmd = "SPIN"

type GameParams struct {
	sgc7pb.GameParam `json:",inline"`
	LastScene        *sgc7game.GameScene      `json:"-"`
	LastOtherScene   *sgc7game.GameScene      `json:"-"`
	MapComponentMsgs map[string]proto.Message `json:"-"`
}

func (gp *GameParams) AddComponentData(name string, cd IComponentData) error {
	if IsStatsComponentMsg {
		pbmsg := cd.BuildPBComponentData()

		gp.MapComponentMsgs[name] = pbmsg

		return nil
	}

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

func (gp *GameParams) SetGameProp(gameProp *GameProperty) error {
	if len(gameProp.MapVals) > 0 {
		gp.MapVals = make(map[int32]int32)

		for k, v := range gameProp.MapVals {
			gp.MapVals[int32(k)] = int32(v)
		}
	}

	return nil
}

func NewGameParam() *GameParams {
	return &GameParams{
		MapComponentMsgs: make(map[string]proto.Message),
	}
}

// gIsForceDisableStats - disable stats
var gIsForceDisableStats bool

// SetForceDisableStats - disable stats
func SetForceDisableStats() {
	gIsForceDisableStats = true
}

// gIsReleaseMode - release mode
var gIsReleaseMode bool

// SetReleaseMode - release mode
func SetReleaseMode() {
	gIsReleaseMode = true
}

type CheckWinType int

const (
	// CheckWinTypeLeftRight - left -> right
	CheckWinTypeLeftRight CheckWinType = 0
	// CheckWinTypeRightLeft - right -> left
	CheckWinTypeRightLeft CheckWinType = 1
	// CheckWinTypeAll - left -> right & right -> left
	CheckWinTypeAll CheckWinType = 2
)

var strCheckWinType map[string]CheckWinType

func ParseCheckWinType(str string) CheckWinType {
	v, isok := strCheckWinType[str]
	if isok {
		return CheckWinType(v)
	}

	return CheckWinTypeLeftRight
}

func initCheckWinType() {
	strCheckWinType = make(map[string]CheckWinType)

	strCheckWinType["left2right"] = CheckWinTypeLeftRight
	strCheckWinType["right2left"] = CheckWinTypeRightLeft
	strCheckWinType["all"] = CheckWinTypeAll
}

var gJsonMode bool

func SetJsonMode() {
	gJsonMode = true
}

func init() {
	initCheckWinType()

	gJsonMode = true
}
