package sgc7game

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrInvalidReels - invalid reels name
	ErrInvalidReels = errors.New("invalid reels name")

	// ErrDuplicateGameMod - duplicate gamemod
	ErrDuplicateGameMod = errors.New("duplicate gamemod")
)
