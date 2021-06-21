package sgc7game

// IPlayerState - player state
type IPlayerState interface {
	// SetPublic - set player public state
	SetPublic(pub interface{}) error
	// SetPrivate - set player private state
	SetPrivate(pri interface{}) error

	// SetPublicJson - set player public state
	SetPublicJson(pub string) error
	// SetPrivateJson - set player private state
	SetPrivateJson(pri string) error

	// GetPublic - get player public state
	GetPublic() interface{}
	// GetPrivate - get player private state
	GetPrivate() interface{}

	// GetPublicJson - set player public state
	GetPublicJson() string
	// GetPrivateJson - set player private state
	GetPrivateJson() string
}
