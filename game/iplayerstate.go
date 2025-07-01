package sgc7game

// IPlayerState - player state
type IPlayerState interface {
	// SetPublic - set player public state
	SetPublic(pub any) error
	// SetPrivate - set player private state
	SetPrivate(pri any) error

	// SetPublicJson - set player public state
	SetPublicJson(pub string) error
	// SetPrivateJson - set player private state
	SetPrivateJson(pri string) error

	// GetPublic - get player public state
	GetPublic() any
	// GetPrivate - get player private state
	GetPrivate() any

	// GetPublicJson - set player public state
	GetPublicJson() string
	// GetPrivateJson - set player private state
	GetPrivateJson() string

	// SetCurGameMod - set current game module
	SetCurGameMod(gamemod string)
	// GetCurGameMod - get current game module
	GetCurGameMod() string

	// OnOutput - on output
	OnOutput()

	Clone() IPlayerState
}
