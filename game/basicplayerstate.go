package sgc7game

// BasicPlayerPublicState - basic PlayerPublicState
type BasicPlayerPublicState struct {
	CurGameMod string
}

// BasicPlayerPrivateState - basic PlayerPrivateState
type BasicPlayerPrivateState struct {
}

// BasicPlayerState - basic PlayerState
type BasicPlayerState struct {
	Public  BasicPlayerPublicState
	Private BasicPlayerPrivateState
}

// NewBasicPlayerState - new BasicPlayerState
func NewBasicPlayerState(curgamemod string) *BasicPlayerState {
	return &BasicPlayerState{
		Public: BasicPlayerPublicState{
			CurGameMod: curgamemod,
		},
		Private: BasicPlayerPrivateState{},
	}
}

// SetPublic - set player public state
func (ps *BasicPlayerState) SetPublic(pub interface{}) error {
	bpub, isok := pub.(*BasicPlayerPublicState)
	if isok {
		ps.Public = *bpub

		return nil
	}

	return ErrInvalidPlayerPublicState
}

// SetPrivate - set player private state
func (ps *BasicPlayerState) SetPrivate(pri interface{}) error {
	bpri, isok := pri.(*BasicPlayerPrivateState)
	if isok {
		ps.Private = *bpri

		return nil
	}

	return ErrInvalidPlayerPrivateState
}
