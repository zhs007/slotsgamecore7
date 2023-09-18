package gamecollection

import (
	"sync"

	"github.com/zhs007/slotsgamecore7/lowcode"
)

type GameMgr struct {
	sync.Mutex
	MapGames map[string]*lowcode.Game
}

func NewGameMgr() *GameMgr {
	return &GameMgr{
		MapGames: make(map[string]*lowcode.Game),
	}
}
