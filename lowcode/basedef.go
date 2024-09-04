package lowcode

import (
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/types/known/anypb"
)

// var IsStatsComponentMsg bool

const (
	TagCurReels = "reels"
)

const DefaultCmd = "SPIN"

type SymbolTriggerType int

const (
	STTypeUnknow             SymbolTriggerType = 0  // 非法
	STTypeLines              SymbolTriggerType = 1  // 线中奖判断，一定是判断全部线，且读paytable来判断是否可以中奖
	STTypeWays               SymbolTriggerType = 2  // ways中奖判断，且读paytable来判断是否可以中奖
	STTypeScatters           SymbolTriggerType = 3  // scatter中奖判断，且读paytable来判断是否可以中奖
	STTypeCountScatter       SymbolTriggerType = 4  // scatter判断，需要传入minnum，不读paytable
	STTypeCountScatterInArea SymbolTriggerType = 5  // 区域内的scatter判断，需要传入minnum，不读paytable
	STTypeCheckLines         SymbolTriggerType = 6  // 线判断，一定是判断全部线，需要传入minnum，不读paytable
	STTypeCheckWays          SymbolTriggerType = 7  // ways判断，需要传入minnum，不读paytable
	STTypeCluster            SymbolTriggerType = 8  // cluster，且读paytable来判断是否可以中奖
	STTypeReelScatters       SymbolTriggerType = 9  // scatter中奖判断，且一轴上只算1个scatter，且读paytable来判断是否可以中奖
	STTypeCountScatterReels  SymbolTriggerType = 10 // scatter中奖判断，且一轴上只算1个scatter，不读paytable
	STTypeAdjacentPay        SymbolTriggerType = 11 // adjacentPay，且读paytable来判断是否可以中奖
)

func ParseSymbolTriggerType(str string) SymbolTriggerType {
	str = strings.ToLower(str)

	if str == "lines" {
		return STTypeLines
	} else if str == "ways" {
		return STTypeWays
	} else if str == "scatters" {
		return STTypeScatters
	} else if str == "countscatter" {
		return STTypeCountScatter
	} else if str == "countscatterinarea" {
		return STTypeCountScatterInArea
	} else if str == "checklines" {
		return STTypeCheckLines
	} else if str == "checkways" {
		return STTypeCheckWays
	} else if str == "cluster" {
		return STTypeCluster
	} else if str == "reelscatters" {
		return STTypeReelScatters
	} else if str == "countscatterreels" {
		return STTypeCountScatterReels
	} else if str == "adjacentpay" {
		return STTypeAdjacentPay
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

type OtherSceneMultiType int

const (
	OSMTNone      OtherSceneMultiType = 0
	OSMTAdd       OtherSceneMultiType = 1 // 每个位置用加来计算总倍数
	OSMTMul       OtherSceneMultiType = 2 // 每个位置用乘来计算总倍数
	OSMTPowOf2Add OtherSceneMultiType = 3 // 每个位置用2的次方之和来计算总倍数
	OSMTPowOf2Mul OtherSceneMultiType = 4 // 每个位置用2的次方之积来计算总倍数
)

func ParseOtherSceneMultiType(str string) OtherSceneMultiType {
	if str == "add" {
		return OSMTAdd
	} else if str == "mul" {
		return OSMTMul
	} else if str == "powof2add" {
		return OSMTPowOf2Add
	} else if str == "powof2mul" {
		return OSMTPowOf2Mul
	}

	return OSMTNone
}

func GetSymbolValMultiFunc(t OtherSceneMultiType) sgc7game.FuncCalcMulti {
	if t == OSMTAdd {
		return func(src int, target int) int {
			if target > 1 {
				if src == 1 {
					return target
				}

				return src + target
			}

			return src
		}
	} else if t == OSMTMul {
		return func(src int, target int) int {
			if target > 1 {
				return src * target
			}

			return src
		}
	} else if t == OSMTPowOf2Add {
		return func(src int, target int) int {
			if target >= 1 {
				if src == 1 {
					return PowInt(2, target)
				}

				return src + PowInt(2, target)
			}

			return src
		}
	} else if t == OSMTPowOf2Mul {
		return func(src int, target int) int {
			if target >= 1 {
				return src * PowInt(2, target)
			}

			return src
		}
	}

	return func(src int, target int) int {
		return 1
	}
}

type GameParams struct {
	sgc7pb.GameParam `json:",inline"`
	LastScene        *sgc7game.GameScene       `json:"-"`
	LastOtherScene   *sgc7game.GameScene       `json:"-"`
	MapComponentData map[string]IComponentData `json:"-"`
}

func (gp *GameParams) AddComponentData(name string, cd IComponentData) error {
	if !gIsReleaseMode {
		gp.MapComponentData[name] = cd.Clone()
	}

	if gIsRTPMode {
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
		MapComponentData: make(map[string]IComponentData),
	}
}

// gIsReleaseMode - release mode
// release模式下，效率会高一些，正式服务器、校验rtp默认都是release模式
var gIsReleaseMode bool

// SetReleaseMode - release mode
func SetReleaseMode() {
	gIsReleaseMode = true
}

// gIsRTPMode - RTP mode
// rtp模式下，不考虑前端数据，所以会更快一些
var gIsRTPMode bool

// SetRTPMode - rtp mode
func SetRTPMode() {
	gIsRTPMode = true
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

// FuncOnChgComponentIntVal - 当这个接口处理完数据，需要返回true，这时底层就不会再处理了
type FuncOnChgComponentIntVal func(componentName string, valName string, off int) bool

var gAllowFullComponentHistory bool


func init() {
	initCheckWinType()
}
