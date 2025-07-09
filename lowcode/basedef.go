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

	switch str {
	case "lines":
		return STTypeLines
	case "ways":
		return STTypeWays
	case "scatters":
		return STTypeScatters
	case "countscatter":
		return STTypeCountScatter
	case "countscatterinarea":
		return STTypeCountScatterInArea
	case "checklines":
		return STTypeCheckLines
	case "checkways":
		return STTypeCheckWays
	case "cluster":
		return STTypeCluster
	case "reelscatters":
		return STTypeReelScatters
	case "countscatterreels":
		return STTypeCountScatterReels
	case "adjacentpay":
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
	switch str {
	case "bet":
		return BTypeBet
	case "totalBet":
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
	switch str {
	case "add":
		return OSMTAdd
	case "mul":
		return OSMTMul
	case "powof2add":
		return OSMTPowOf2Add
	case "powof2mul":
		return OSMTPowOf2Mul
	}

	return OSMTNone
}

func GetSymbolValMultiFunc(t OtherSceneMultiType) sgc7game.FuncCalcMulti {
	switch t {
	case OSMTAdd:
		return func(src int, target int) int {
			if target > 1 {
				if src == 1 {
					return target
				}

				return src + target
			}

			return src
		}
	case OSMTMul:
		return func(src int, target int) int {
			if target > 1 {
				return src * target
			}

			return src
		}
	case OSMTPowOf2Add:
		return func(src int, target int) int {
			if target >= 1 {
				if src == 1 {
					return PowInt(2, target)
				}

				return src + PowInt(2, target)
			}

			return src
		}
	case OSMTPowOf2Mul:
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

// FuncOnSettedComponentIntVal - 处理完后才调用这个接口，val 是最终数值
type FuncOnSettedComponentIntVal func(componentName string, valName string, val int)

// FuncOnChgedComponentIntVal - 处理完后才调用这个接口，val 是最终数值
type FuncOnChgedComponentIntVal func(componentName string, valName string, val int, off int)

var gAllowFullComponentHistory bool

func SetAllowFullComponentHistory() {
	gAllowFullComponentHistory = true
}

var gRngLibConfig string

func SetRngLibConfig(fn string) {
	gRngLibConfig = fn
}

var gIsIgnoreGenDefaultScene bool

func SetIgnoreGenDefaultScene() {
	gIsIgnoreGenDefaultScene = true
}

func init() {
	initCheckWinType()
}
