package mathtoolset

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrInvalidReel - invalid reel
	ErrInvalidReel = errors.New("invalid reel")

	// ErrInvalidInReelSymbolType - invalid InReelSymbolType
	ErrInvalidInReelSymbolType = errors.New("invalid InReelSymbolType")

	// ErrInvalidReelsStatsExcelFile - invalid ReelsStats excel file
	ErrInvalidReelsStatsExcelFile = errors.New("invalid ReelsStats excel file")

	// ErrInvalidReelsStatsExcelColname - invalid ReelsStats excel colname
	ErrInvalidReelsStatsExcelColname = errors.New("invalid ReelsStats excel colname")

	// ErrNoValidSymbols - no valid symbols
	ErrNoValidSymbols = errors.New("no valid symbols")

	// ErrValidParamInAutoChgWeights - invalid param in AutoChgWeights
	ErrValidParamInAutoChgWeights = errors.New("invalid param in AutoChgWeights")

	// ErrInvalidDataInAGRDataList - invalid data in agrDataList
	ErrInvalidDataInAGRDataList = errors.New("invalid data in agrDataList")

	// ErrNoResultInAutoChgWeights - no result in AutoChgWeights
	ErrNoResultInAutoChgWeights = errors.New("no result in AutoChgWeights")

	// ErrInvalidReelsWithMinOff - invalid reels with minoff
	ErrInvalidReelsWithMinOff = errors.New("invalid reels with minoff")

	// ErrInvalidCode - invalid code
	ErrInvalidCode = errors.New("invalid code")

	// ErrInvalidFunctionParams - invalid function params
	ErrInvalidFunctionParams = errors.New("invalid function params")
)
