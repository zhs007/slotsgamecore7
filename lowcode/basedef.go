package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var IsStatsComponentMsg bool

const (
	TagCurReels = "reels"
)

const DefaultCmd = "SPIN"

type SymbolTriggerType int

const (
	STTypeUnknow             SymbolTriggerType = 0 // 非法
	STTypeLines              SymbolTriggerType = 1 // 线中奖判断，一定是判断全部线，且读paytable来判断是否可以中奖
	STTypeWays               SymbolTriggerType = 2 // ways中奖判断，且读paytable来判断是否可以中奖
	STTypeScatters           SymbolTriggerType = 3 // scatter中奖判断，且读paytable来判断是否可以中奖
	STTypeCountScatter       SymbolTriggerType = 4 // scatter判断，需要传入minnum，不读paytable
	STTypeCountScatterInArea SymbolTriggerType = 5 // 区域内的scatter判断，需要传入minnum，不读paytable
	STTypeCheckLines         SymbolTriggerType = 6 // 线判断，一定是判断全部线，需要传入minnum，不读paytable
	STTypeCheckWays          SymbolTriggerType = 7 // ways判断，需要传入minnum，不读paytable
	STTypeCluster            SymbolTriggerType = 8 // cluster，且读paytable来判断是否可以中奖
	STTypeReelScatters       SymbolTriggerType = 9 // scatter中奖判断，且一轴上只算1个scatter，且读paytable来判断是否可以中奖
)

func ParseSymbolTriggerType(str string) SymbolTriggerType {
	if str == "lines" {
		return STTypeLines
	} else if str == "ways" {
		return STTypeWays
	} else if str == "scatters" {
		return STTypeScatters
	} else if str == "countscatter" {
		return STTypeCountScatter
	} else if str == "countscatterInArea" {
		return STTypeCountScatterInArea
	} else if str == "checkLines" {
		return STTypeCheckLines
	} else if str == "checkWays" {
		return STTypeCheckWays
	} else if str == "cluster" {
		return STTypeCluster
	} else if str == "reelscatters" {
		return STTypeReelScatters
	}

	return STTypeUnknow
}

type BetType int

const (
	BTypeNoPay    BetType = 0
	BTypeBet      BetType = 1
	BTypeTotalBet BetType = 2
)

func ParseBetType(str string) BetType {
	if str == "bet" {
		return BTypeBet
	} else if str == "totalBet" {
		return BTypeTotalBet
	}

	return BTypeNoPay
}

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
			goutils.Err(err))

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

// // gIsForceDisableStats - disable stats
// var gIsForceDisableStats bool

// // SetForceDisableStats - disable stats
// func SetForceDisableStats() {
// 	gIsForceDisableStats = true
// }

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
	// CheckWinTypeCount - count
	CheckWinTypeCount CheckWinType = 3
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
	strCheckWinType["count"] = CheckWinTypeCount
}

// json mode - 是否json配置模式，json模式下会有些特殊处理，少一些互斥锁
var gJsonMode bool

func SetJsonMode() {
	gJsonMode = true
}

// statsv2 - 是否开启 stats ，默认不开启，有cpu消耗
var gAllowStats2 bool

func SetAllowStatsV2() {
	gAllowStats2 = true
}

var gAllowForceOutcome bool
var gMaxForceOutcomeTimes int

func SetAllowForceOutcome(maxTry int) {
	gAllowForceOutcome = true
	gMaxForceOutcomeTimes = maxTry
}

const MaxStepNum = 1000
const MaxComponentNumInStep = 100

const BasicGameModName = "basic"

func init() {
	initCheckWinType()

	// gIsForceDisableStats = true
}
