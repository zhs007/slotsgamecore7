package sgc7game

// IPlayerState - player state
type IPlayerState interface {
	// SetPublic - set player public state
	SetPublic(pub interface{})
	// SetPrivate - set player private state
	SetPrivate(pri interface{})
}
