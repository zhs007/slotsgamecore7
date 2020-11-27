package sgc7game

// FuncNewBasicPlayerState - new BasicPlayerState and set PlayerBoostData
type FuncNewBasicPlayerState func() *BasicPlayerState

// NewBPSNoBoostData - new a BasicPlayerState without boostdata
func NewBPSNoBoostData() *BasicPlayerState {
	return &BasicPlayerState{}
}

// BasicPlayerPublicState - basic PlayerPublicState
type BasicPlayerPublicState struct {
	CurGameMod string `json:"curgamemod"`
	NextM      int    `json:"nextm"`
}

// BasicPlayerPrivateState - basic PlayerPrivateState
type BasicPlayerPrivateState struct {
}

// BasicPlayerState - basic PlayerState
type BasicPlayerState struct {
	Public  *BasicPlayerPublicState
	Private *BasicPlayerPrivateState
}

// NewBasicPlayerState - new BasicPlayerState
func NewBasicPlayerState(curgamemod string) *BasicPlayerState {
	bps := &BasicPlayerState{
		Public: &BasicPlayerPublicState{
			CurGameMod: curgamemod,
		},
		Private: &BasicPlayerPrivateState{},
	}

	return bps
}

// SetPublic - set player public state
func (ps *BasicPlayerState) SetPublic(pub interface{}) error {
	bpub, isok := pub.(*BasicPlayerPublicState)
	if isok {
		ps.Public = bpub

		return nil
	}

	return ErrInvalidPlayerPublicState
}

// SetPrivate - set player private state
func (ps *BasicPlayerState) SetPrivate(pri interface{}) error {
	bpri, isok := pri.(*BasicPlayerPrivateState)
	if isok {
		ps.Private = bpri

		return nil
	}

	return ErrInvalidPlayerPrivateState
}

// GetPublic - get player public state
func (ps *BasicPlayerState) GetPublic() interface{} {
	return ps.Public
}

// GetPrivate - get player private state
func (ps *BasicPlayerState) GetPrivate() interface{} {
	return ps.Private
}
