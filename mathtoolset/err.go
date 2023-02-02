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
)
