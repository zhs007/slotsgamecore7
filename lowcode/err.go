package lowcode

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrMustHaveMainPaytables - must have main paytables
	ErrMustHaveMainPaytables = errors.New("must have main paytables")

	// ErrInvalidGameMod - invalid gamemod
	ErrInvalidGameMod = errors.New("invalid gamemod")

	// ErrInvalidComponent - invalid component
	ErrInvalidComponent = errors.New("invalid component")

	// ErrInvalidReels - invalid reels
	ErrInvalidReels = errors.New("invalid reels")

	// ErrInvalidSymbol - invalid symbol
	ErrInvalidSymbol = errors.New("invalid symbol")

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

	// ErrIvalidCurGameModParams - invalid CurGameModParams
	ErrIvalidCurGameModParams = errors.New("invalid CurGameModParams")

	// ErrIvalidPlayResultLength - invalid PlayResult Length
	ErrIvalidPlayResultLength = errors.New("invalid PlayResult Length")

	// ErrIvalidMultiLevelReelsConfig - invalid MultiLevelReels config
	ErrIvalidMultiLevelReelsConfig = errors.New("invalid MultiLevelReels config")

	// ErrIvalidStatsSymbolsInConfig - invalid StatsSymbols in config
	ErrIvalidStatsSymbolsInConfig = errors.New("invalid StatsSymbols in config")
	// ErrIvalidStatsComponentInConfig - invalid Stats's component in config
	ErrIvalidStatsComponentInConfig = errors.New("invalid Stats's component in config")

	// ErrIvalidComponentConfig - invalid component config
	ErrIvalidComponentConfig = errors.New("invalid component config")

	// ErrIvalidGameData - invalid gameData
	ErrIvalidGameData = errors.New("invalid gameData")

	// ErrIvalidSimpleRNG - invalid SimpleRNG
	ErrIvalidSimpleRNG = errors.New("invalid SimpleRNG")

	// ErrIvalidCmd - invalid cmd
	ErrIvalidCmd = errors.New("invalid cmd")
	// ErrIvalidCmdParam - invalid cmdparam
	ErrIvalidCmdParam = errors.New("invalid cmdparam")

	// ErrIvalidTagCurReels - invalid TagCurReels
	ErrIvalidTagCurReels = errors.New("invalid TagCurReels")

	// ErrIvalidSymbolCollection - invalid SymbolColletion
	ErrIvalidSymbolCollection = errors.New("invalid SymbolColletion")

	// ErrIvalidCustomNode - invalid custom-node
	ErrIvalidCustomNode = errors.New("invalid custom-node")
	// ErrIvalidTriggerLabel - invalid trigger label
	ErrIvalidTriggerLabel = errors.New("invalid trigger label")
	// ErrIvalidPayTables - invalid paytables
	ErrIvalidPayTables = errors.New("invalid paytables")
	// ErrIvalidSymbolInReels - invalid symbol in reels
	ErrIvalidSymbolInReels = errors.New("invalid symbol in reels")
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
	// ErrIvalidReels - invalid reels
	ErrIvalidReels = errors.New("invalid reels")
	// ErrUnsupportedOtherList - unsupported otherList
	ErrUnsupportedOtherList = errors.New("unsupported otherList")

	// ErrIvalidDefaultScene - invalid default scene
	ErrIvalidDefaultScene = errors.New("invalid default scene")
	// ErrIvalidWidth - invalid width
	ErrIvalidWidth = errors.New("invalid width")
	// ErrIvalidHeight - invalid height
	ErrIvalidHeight = errors.New("invalid height")

	// ErrIvalidProto - invalid proto
	ErrIvalidProto = errors.New("invalid proto")

	// ErrIvalidSymbol - invalid symbol
	ErrIvalidSymbol = errors.New("invalid symbol")

	// ErrIvalidSymbolTriggerType - invalid SymbolTriggerType
	ErrIvalidSymbolTriggerType = errors.New("invalid SymbolTriggerType")

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

	// ErrNoComponent - no component
	ErrNoComponent = errors.New("no component")
)
