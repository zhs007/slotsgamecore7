package lowcode

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrMustHaveMainPaytables - must have main paytables
	ErrMustHaveMainPaytables = errors.New("must have main paytables")

	// ErrInvalidGameMod - invalid gamemod
	ErrInvalidGameMod = errors.New("invalid gamemod")

	// ErrInvalidGameConfig - invalid game config
	ErrInvalidGameConfig = errors.New("invalid game config")

	// ErrInvalidComponent - invalid component
	ErrInvalidComponent = errors.New("invalid component")

	// ErrInvalidPaytables - invalid paytables
	ErrInvalidPaytables = errors.New("invalid paytables")

	// ErrInvalidLineData - invalid line data
	ErrInvalidLineData = errors.New("invalid line data")

	// ErrInvalidGamePropertyString - invalid gameProperty string
	ErrInvalidGamePropertyString = errors.New("invalid gameProperty string")

	// ErrParseScript - parse script error
	ErrParseScript = errors.New("parse script error")
	// ErrNoFunctionInScript - no function in script
	ErrNoFunctionInScript = errors.New("no function in script")
	// ErrWrongFunctionInScript - wrong function in script
	ErrWrongFunctionInScript = errors.New("wrong function in script")

	// ErrInvalidComponentName - invalid component name
	ErrInvalidComponentName = errors.New("invalid component name")

	// ErrInvalidCurGameModParams - invalid CurGameModParams
	ErrInvalidCurGameModParams = errors.New("invalid CurGameModParams")

	// ErrInvalidPlayResultLength - invalid PlayResult Length
	ErrInvalidPlayResultLength = errors.New("invalid PlayResult Length")

	// ErrInvalidMultiLevelReelsConfig - invalid MultiLevelReels config
	ErrInvalidMultiLevelReelsConfig = errors.New("invalid MultiLevelReels config")

	// ErrInvalidStatsSymbolsInConfig - invalid StatsSymbols in config
	ErrInvalidStatsSymbolsInConfig = errors.New("invalid StatsSymbols in config")
	// ErrInvalidStatsComponentInConfig - invalid Stats's component in config
	ErrInvalidStatsComponentInConfig = errors.New("invalid Stats's component in config")

	// ErrInvalidComponentConfig - invalid component config
	ErrInvalidComponentConfig = errors.New("invalid component config")

	// ErrInvalidGameData - invalid gameData
	ErrInvalidGameData = errors.New("invalid gameData")

	// ErrInvalidPlayerState - invalid playerState
	ErrInvalidPlayerState = errors.New("invalid playerState")

	// ErrInvalidSimpleRNG - invalid SimpleRNG
	ErrInvalidSimpleRNG = errors.New("invalid SimpleRNG")

	// ErrInvalidCmd - invalid cmd
	ErrInvalidCmd = errors.New("invalid cmd")
	// ErrInvalidCmdParam - invalid cmdparam
	ErrInvalidCmdParam = errors.New("invalid cmdparam")

	// ErrInvalidTagCurReels - invalid TagCurReels
	ErrInvalidTagCurReels = errors.New("invalid TagCurReels")

	// ErrInvalidSymbolCollection - invalid SymbolColletion
	ErrInvalidSymbolCollection = errors.New("invalid SymbolColletion")

	// ErrInvalidCustomNode - invalid custom-node
	ErrInvalidCustomNode = errors.New("invalid custom-node")
	// ErrInvalidTriggerLabel - invalid trigger label
	ErrInvalidTriggerLabel = errors.New("invalid trigger label")
	// ErrInvalidPayTables - invalid paytables
	ErrInvalidPayTables = errors.New("invalid paytables")
	// ErrInvalidSymbolInReels - invalid symbol in reels
	ErrInvalidSymbolInReels = errors.New("invalid symbol in reels")
	// ErrNoComponentValues - no componentValues
	ErrNoComponentValues = errors.New("no componentValues")
	// ErrUnsupportedComponentType - unsupported componentType
	ErrUnsupportedComponentType = errors.New("unsupported componentType")
	// ErrUnsupportedLinkType - unsupported link type
	ErrUnsupportedLinkType = errors.New("unsupported link type")
	// ErrUnsupportedControllerType - unsupported ControllerType
	ErrUnsupportedControllerType = errors.New("unsupported ControllerType")
	// ErrInvalidJsonNode - invalid json node
	ErrInvalidJsonNode = errors.New("invalid json node")
	// ErrInvalidReels - invalid reels
	ErrInvalidReels = errors.New("invalid reels")
	// ErrUnsupportedOtherList - unsupported otherList
	ErrUnsupportedOtherList = errors.New("unsupported otherList")

	// ErrInvalidDefaultScene - invalid default scene
	ErrInvalidDefaultScene = errors.New("invalid default scene")
	// ErrInvalidWidth - invalid width
	ErrInvalidWidth = errors.New("invalid width")
	// ErrInvalidHeight - invalid height
	ErrInvalidHeight = errors.New("invalid height")

	// ErrInvalidProto - invalid proto
	ErrInvalidProto = errors.New("invalid proto")

	// ErrInvalidSymbol - invalid symbol
	ErrInvalidSymbol = errors.New("invalid symbol")

	// ErrInvalidSymbolTriggerType - invalid SymbolTriggerType
	ErrInvalidSymbolTriggerType = errors.New("invalid SymbolTriggerType")

	// ErrNotMask - not mask
	ErrNotMask = errors.New("not mask")
	// ErrNotRespin - not respin
	ErrNotRespin = errors.New("not respin")

	// ErrInvalidSymbolNum - invalid SymbolNum
	ErrInvalidSymbolNum = errors.New("invalid SymbolNum")
	// ErrInvalidComponentVal - invalid ComponentVal
	ErrInvalidComponentVal = errors.New("invalid ComponentVal")
	// ErrInvalidBet - invalid Bet
	ErrInvalidBet = errors.New("invalid Bet")

	// ErrInvalidIntValMappingFile - invalid IntValMappingFile
	ErrInvalidIntValMappingFile = errors.New("invalid IntValMappingFile")
	// ErrInvalidIntValMappingValue - invalid IntValMapping value
	ErrInvalidIntValMappingValue = errors.New("invalid IntValMapping value")

	// ErrInvalidWeightVal - invalid weight value
	ErrInvalidWeightVal = errors.New("invalid weight value")

	// ErrComponentDoNothing - component do nothing
	ErrComponentDoNothing = errors.New("component do nothing")

	// ErrTooManySteps - too many steps
	ErrTooManySteps = errors.New("too many steps")
	// ErrTooManyComponentsInStep - too many components in step
	ErrTooManyComponentsInStep = errors.New("too many components in step")

	// ErrCannotForceOutcome - cannot force outcome
	ErrCannotForceOutcome = errors.New("cannot force outcome")

	// ErrInvalidCallStackNode - invalid callstack node
	ErrInvalidCallStackNode = errors.New("invalid callstack node")

	// ErrInvalidComponentChildren - invalid component children
	ErrInvalidComponentChildren = errors.New("invalid component children")

	// ErrInvalidForceOutcome2Code - invalid ForceOutcome2 code
	ErrInvalidForceOutcome2Code = errors.New("invalid ForceOutcome2 code")
	// ErrInvalidForceOutcome2ReturnVal - invalid ForceOutcome2 return value
	ErrInvalidForceOutcome2ReturnVal = errors.New("invalid ForceOutcome2 return value")

	// ErrInvalidOtherScene - invalid OtherScene
	ErrInvalidOtherScene = errors.New("invalid OtherScene")

	// ErrInvalidScene - invalid Scene
	ErrInvalidScene = errors.New("invalid Scene")

	// ErrInvalidSetComponent - invalid a set component
	ErrInvalidSetComponent = errors.New("invalid a set component")

	// ErrInvalidScriptParamsNumber - invalid script params number
	ErrInvalidScriptParamsNumber = errors.New("invalid script params number")
	// ErrInvalidScriptParamType - invalid script param type
	ErrInvalidScriptParamType = errors.New("invalid param type")

	// ErrInvalidPosition - invalid position
	ErrInvalidPosition = errors.New("invalid position")

	// ErrInvalidPositionCollection - invalid positionCollection
	ErrInvalidPositionCollection = errors.New("invalid positionCollection")

	// ErrInvalidCollectorVal - invalid Collector.Val
	ErrInvalidCollectorVal = errors.New("invalid Collector.Val")
	// ErrInvalidCollectorLogic - invalid Collector logic
	ErrInvalidCollectorLogic = errors.New("invalid Collector logic")

	// ErrInvalidAnyProtoBuf - invalid AnyProtoBuf
	ErrInvalidAnyProtoBuf = errors.New("invalid AnyProtoBuf")
	// ErrInvalidPBComponentData - invalid invalid PB ComponentData
	ErrInvalidPBComponentData = errors.New("invalid PB ComponentData")
	// ErrInvalidFuncNewComponentData - invalid FuncNewComponentData
	ErrInvalidFuncNewComponentData = errors.New("invalid FuncNewComponentData")
	// ErrInvalidAnypbTypeURL - invalid anypb TypeURL
	ErrInvalidAnypbTypeURL = errors.New("invalid anypb TypeURL")
	// ErrNoWeight - no weight
	ErrNoWeight = errors.New("no weight")

	// ErrInvalidComponentData - invalid invalid ComponentData
	ErrInvalidComponentData = errors.New("invalid ComponentData")

	// ErrInvalidBranch - invalid branch
	ErrInvalidBranch = errors.New("invalid branch")

	// ErrInvalidCommand - invalid command
	ErrInvalidCommand = errors.New("invalid command")

	// ErrNoComponent - no component
	ErrNoComponent = errors.New("no component")

	// ErrCanNotGenDefaultScene - can not gen default scene
	ErrCanNotGenDefaultScene = errors.New("can not gen default scene")

	// ErrDeprecatedAPI - deprecated API
	ErrDeprecatedAPI = errors.New("deprecated API")
)
