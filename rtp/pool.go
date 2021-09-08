package sgc7rtp

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// FuncOnPlayer - onPlayer(*PlayerPoolData, sgc7game.IPlayerState)
type FuncOnPlayer func(pd *PlayerPoolData, ps *sgc7game.IPlayerState) bool

// PlayerPoolData -
type PlayerPoolData struct {
	TagName    string
	Total      int64
	PlayerNums int64
	OnPlayer   FuncOnPlayer
}

// NewPlayerPoolData - new PlayerPoolData
func NewPlayerPoolData(tag string, onPlayer FuncOnPlayer) *PlayerPoolData {
	return &PlayerPoolData{
		TagName:  tag,
		OnPlayer: onPlayer,
	}
}

// Clone - clone
func (pd *PlayerPoolData) Clone() *PlayerPoolData {
	pd1 := &PlayerPoolData{
		TagName:    pd.TagName,
		Total:      pd.Total,
		PlayerNums: pd.PlayerNums,
		OnPlayer:   pd.OnPlayer,
	}

	return pd1
}

// Add - add
func (pd *PlayerPoolData) Add(pd1 *PlayerPoolData) {
	if pd.TagName == pd1.TagName {
		pd.Total += pd1.Total
		pd.PlayerNums += pd1.PlayerNums
	}
}
