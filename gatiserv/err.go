package gatiserv

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrInvalidPlayerState - invalid PlayerState
	ErrInvalidPlayerState = errors.New("invalid PlayerState")

	// ErrInvalidCriticalComponentID - invalid CriticalComponentID
	ErrInvalidCriticalComponentID = errors.New("invalid CriticalComponentID")

	// ErrNonStatusOK - non statusOK
	ErrNonStatusOK = errors.New("non statusOK")
)
