package sgc7game

import (
	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
)

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
func (ps *BasicPlayerState) SetPublic(pub any) error {
	bpub, isok := pub.(*BasicPlayerPublicState)
	if isok {
		ps.Public = bpub

		return nil
	}

	return ErrInvalidPlayerPublicState
}

// SetPrivate - set player private state
func (ps *BasicPlayerState) SetPrivate(pri any) error {
	bpri, isok := pri.(*BasicPlayerPrivateState)
	if isok {
		ps.Private = bpri

		return nil
	}

	return ErrInvalidPlayerPrivateState
}

// GetPublic - get player public state
func (ps *BasicPlayerState) GetPublic() any {
	return ps.Public
}

// GetPrivate - get player private state
func (ps *BasicPlayerState) GetPrivate() any {
	return ps.Private
}

// SetPublicJson - set player public state
func (ps *BasicPlayerState) SetPublicJson(pubjson string) error {
	pub := &BasicPlayerPublicState{}
	err := sonic.Unmarshal([]byte(pubjson), pub)
	if err != nil {
		goutils.Warn("BasicPlayerState.SetPublicJson",
			goutils.Err(err))

		return err
	}

	ps.SetPublic(pub)

	return nil
}

// SetPrivateJson - set player private state
func (ps *BasicPlayerState) SetPrivateJson(prijson string) error {
	pri := &BasicPlayerPrivateState{}
	err := sonic.Unmarshal([]byte(prijson), pri)
	if err != nil {
		goutils.Warn("BasicPlayerState.SetPrivateJson",
			goutils.Err(err))

		return err
	}

	ps.SetPrivate(pri)

	return nil
}

// GetPublicJson - set player public state
func (ps *BasicPlayerState) GetPublicJson() string {
	bpub, err := sonic.Marshal(ps.GetPublic())
	if err != nil {
		goutils.Warn("BasicPlayerState.GetPublicJson",
			goutils.Err(err))

		return ""
	}

	return string(bpub)
}

// GetPrivateJson - set player private state
func (ps *BasicPlayerState) GetPrivateJson() string {
	bpri, err := sonic.Marshal(ps.GetPrivate())
	if err != nil {
		goutils.Warn("BasicPlayerState.GetPrivateJson",
			goutils.Err(err))

		return ""
	}

	return string(bpri)
}

// SetCurGameMod - set current game module
func (ps *BasicPlayerState) SetCurGameMod(gamemod string) {
	ps.Public.CurGameMod = gamemod
}

// GetCurGameMod - get current game module
func (ps *BasicPlayerState) GetCurGameMod() string {
	return ps.Public.CurGameMod
}

func (ps *BasicPlayerState) Clone() IPlayerState {
	nps := NewBasicPlayerState(ps.GetCurGameMod())

	nps.SetPublicJson(ps.GetPublicJson())
	nps.SetPrivateJson(ps.GetPrivateJson())

	return nps
}

// OnOutput - on output
func (ps *BasicPlayerState) OnOutput() {

}
