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
}

// BasicPlayerPrivateState - basic PlayerPrivateState
type BasicPlayerPrivateState struct {
}

// BasicPlayerState - basic PlayerState
type BasicPlayerState struct {
	Public  *BasicPlayerPublicState
	Private *BasicPlayerPrivateState
}

// // NewBasicPlayerStateEx - new BasicPlayerState
// func NewBasicPlayerStateEx(pub string, pri string) *BasicPlayerState {
// 	ps := &BasicPlayerState{}

// 	err := ps.SetPublicString(pub)
// 	if err != nil {
// 		sgc7utils.Error("NewBasicPlayerStateEx:SetPublicString",
// 			zap.Error(err),
// 			zap.String("pub", pub),
// 			zap.String("pri", pri))

// 		return nil
// 	}

// 	err = ps.SetPrivateString(pri)
// 	if err != nil {
// 		sgc7utils.Error("NewBasicPlayerStateEx:SetPrivateString",
// 			zap.Error(err),
// 			zap.String("pub", pub),
// 			zap.String("pri", pri))

// 		return nil
// 	}

// 	return ps
// }

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

// // SetPublicString - set player public state
// func (ps *BasicPlayerState) SetPublicString(pub string) error {
// 	json := jsoniter.ConfigCompatibleWithStandardLibrary

// 	err := json.Unmarshal([]byte(pub), &ps.Public)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // SetPrivateString - set player private state
// func (ps *BasicPlayerState) SetPrivateString(pri string) error {
// 	json := jsoniter.ConfigCompatibleWithStandardLibrary

// 	err := json.Unmarshal([]byte(pri), &ps.Private)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
