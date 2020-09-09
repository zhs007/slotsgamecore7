package sgc7game

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrInvalidReels - invalid reels name
	ErrInvalidReels = errors.New("invalid reels name")

	// ErrDuplicateGameMod - duplicate gamemod
	ErrDuplicateGameMod = errors.New("duplicate gamemod")
	// ErrInvalidGameMod - invalid GameMod
	ErrInvalidGameMod = errors.New("invalid GameMod")
	// ErrInvalidWHGameMod - invalid Width or Height in GameMod
	ErrInvalidWHGameMod = errors.New("invalid Width or Height in GameMod")

	// ErrInvalidCommand - invalid command
	ErrInvalidCommand = errors.New("invalid command")

	// ErrInvalidBasicPlayerState - invalid BasicPlayerState
	ErrInvalidBasicPlayerState = errors.New("invalid BasicPlayerState")
	// ErrInvalidPlayerPublicState - invalid PlayerPublicState
	ErrInvalidPlayerPublicState = errors.New("invalid PlayerPublicState")
	// ErrInvalidPlayerPrivateState - invalid PlayerPrivateState
	ErrInvalidPlayerPrivateState = errors.New("invalid PlayerPrivateState")

	// ErrNonGameModCalcScene - non CalcScene in GameMod
	ErrNonGameModCalcScene = errors.New("non CalcScene in GameMod")
	// ErrNonGameModPayout - non Payout in GameMod
	ErrNonGameModPayout = errors.New("non Payout in GameMod")

	// ErrInvalidWeights - invalid weights
	ErrInvalidWeights = errors.New("invalid weights")

	// ErrNullConfig - null config
	ErrNullConfig = errors.New("null config")

	// ErrInvalidArray - invalid array
	ErrInvalidArray = errors.New("invalid array")
)
