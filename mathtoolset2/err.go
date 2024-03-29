package mathtoolset2

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

	// ErrInvalidTargetRTP - invalid targetRTP
	ErrInvalidTargetRTP = errors.New("invalid targetRTP")

	// ErrUnimplementedCode - unimplemented code
	ErrUnimplementedCode = errors.New("unimplemented code")

	// ErrInvalidScatterNumber - invalid scatter number
	ErrInvalidScatterNumber = errors.New("invalid scatter number")

	// ErrCannotBeConverged - cannot be converged
	ErrCannotBeConverged = errors.New("cannot be converged")

	// ErrWinWeightMerge - WinWeight.Merge error
	ErrWinWeightMerge = errors.New("WinWeight.Merge error")
	// ErrWinWeightScale - WinWeight.scale error
	ErrWinWeightScale = errors.New("WinWeight.scale error")
	// ErrDuplicateAvgWin - duplicate avgwin
	ErrDuplicateAvgWin = errors.New("duplicate avgwin")

	// ErrInvalidReelsStats2File - invalid reelsstats2 file
	ErrInvalidReelsStats2File = errors.New("invalid reelsstats2 file")
	// ErrGenStackReel - genStackReel error
	ErrGenStackReel = errors.New("genStackReel error")

	// ErrReturnNotOK -
	ErrReturnNotOK = errors.New("return not ok")

	// ErrRunError -
	ErrRunError = errors.New("run error")

	// ErrInvalidFileData - invalid filedata
	ErrInvalidFileData = errors.New("invalid filedata")
)
