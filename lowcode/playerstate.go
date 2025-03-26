package lowcode

import (
	"github.com/bytedance/sonic"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type BetPS struct {
	MapComponentData map[string]IComponentPS `json:"mapComponentData"`
}

type BetMethodPS struct {
	MapBet map[int]*BetPS `json:"mapBet"`
}

func (bmps *BetMethodPS) GetBetPS(bet int) *BetPS {
	bps, isok := bmps.MapBet[bet]
	if !isok {
		bmps.MapBet[bet] = &BetPS{
			MapComponentData: make(map[string]IComponentPS),
		}

		return bmps.MapBet[bet]
	}

	return bps
}

func (bmps *BetMethodPS) GetBetCPS(bet int, componentName string) IComponentPS {
	bps, isok := bmps.MapBet[bet]
	if !isok {
		return nil
	}

	cps, isok := bps.MapComponentData[componentName]
	if !isok {
		return nil
	}

	return cps
}

func (bmps *BetMethodPS) HasBetPS(bet int) bool {
	_, isok := bmps.MapBet[bet]
	return isok
}

// PlayerState - player state
type PlayerState struct {
	MapBetMothodPub map[int]*BetMethodPS
	MapBetMothodPri map[int]*BetMethodPS
}

// SetPublic - set player public state
func (ps *PlayerState) SetPublic(pub any) error {
	goutils.Error("PlayerState.SetPublic",
		goutils.Err(ErrDeprecatedAPI))

	return ErrDeprecatedAPI
}

// SetPrivate - set player private state
func (ps *PlayerState) SetPrivate(pri any) error {
	goutils.Error("PlayerState.SetPrivate",
		goutils.Err(ErrDeprecatedAPI))

	return ErrDeprecatedAPI
}

// SetPublicJson - set player public state
func (ps *PlayerState) SetPublicJson(pub string) error {
	err := sonic.UnmarshalString(pub, &ps.MapBetMothodPub)
	if err != nil {
		goutils.Error("PlayerState.SetPublicJson",
			goutils.Err(err))

		return err
	}

	return nil
}

// SetPrivateJson - set player private state
func (ps *PlayerState) SetPrivateJson(pri string) error {
	err := sonic.UnmarshalString(pri, &ps.MapBetMothodPri)
	if err != nil {
		goutils.Error("PlayerState.SetPrivateJson",
			goutils.Err(err))

		return err
	}

	return nil
}

// GetPublic - get player public state
func (ps *PlayerState) GetPublic() any {
	goutils.Error("PlayerState.GetPublic",
		goutils.Err(ErrDeprecatedAPI))

	return nil
}

// GetPrivate - get player private state
func (ps *PlayerState) GetPrivate() any {
	goutils.Error("PlayerState.GetPrivate",
		goutils.Err(ErrDeprecatedAPI))

	return nil
}

// GetPublicJson - set player public state
func (ps *PlayerState) GetPublicJson() string {
	str, err := sonic.MarshalString(ps.MapBetMothodPub)
	if err != nil {
		goutils.Error("PlayerState.GetPublicJson",
			goutils.Err(err))
	}

	return str
}

// GetPrivateJson - set player private state
func (ps *PlayerState) GetPrivateJson() string {
	str, err := sonic.MarshalString(ps.MapBetMothodPri)
	if err != nil {
		goutils.Error("PlayerState.GetPrivateJson",
			goutils.Err(err))
	}

	return str
}

// SetCurGameMod - set current game module
func (ps *PlayerState) SetCurGameMod(gamemod string) {
}

// GetCurGameMod - get current game module
func (ps *PlayerState) GetCurGameMod() string {
	return BasicGameModName
}

func (ps *PlayerState) Clone() sgc7game.IPlayerState {
	dest := &PlayerState{}

	dest.SetPublicJson(ps.GetPublicJson())
	dest.SetPrivateJson(ps.GetPrivateJson())

	return dest
}

func (ps *PlayerState) GetBetMethodPub(betMethod int) *BetMethodPS {
	betmps, isok := ps.MapBetMothodPub[betMethod]
	if !isok {
		ps.MapBetMothodPub[betMethod] = &BetMethodPS{
			MapBet: make(map[int]*BetPS),
		}

		return ps.MapBetMothodPub[betMethod]
	}

	return betmps
}

func NewPlayerState() *PlayerState {
	return &PlayerState{
		MapBetMothodPub: make(map[int]*BetMethodPS),
		MapBetMothodPri: make(map[int]*BetMethodPS),
	}
}
