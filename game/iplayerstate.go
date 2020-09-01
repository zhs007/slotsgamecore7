package sgc7game

// IPlayerState - player state
type IPlayerState interface {
	// SetPublic - set player public state
	SetPublic(pub interface{}) error
	// SetPrivate - set player private state
	SetPrivate(pri interface{}) error

	// SetPublicString - set player public state
	SetPublicString(pub string) error
	// SetPrivateString - set player private state
	SetPrivateString(pri string) error

	// GetPublic - get player public state
	GetPublic() interface{}
	// GetPrivate - get player private state
	GetPrivate() interface{}
}
