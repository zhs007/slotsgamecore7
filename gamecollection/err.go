package gamecollection

import "errors"

var (
	// ErrInvalidGameCode - invalid gameCode
	ErrInvalidGameCode = errors.New("invalid gameCode")
	// ErrInvalidGameParams - invalid GameParams
	ErrInvalidGameParams = errors.New("invalid GameParams")
)
