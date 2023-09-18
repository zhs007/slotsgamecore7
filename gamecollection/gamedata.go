package gamecollection

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/zhs007/slotsgamecore7/lowcode"
)

type GameData struct {
	Name     string
	HashCode string
	Data     string
	Game     *lowcode.Game
}

func NewGameData(name string, data string) (*GameData, error) {

	gameD := &GameData{
		Name: name,
		Data: data,
	}

	hasher := sha1.New()
	hasher.Write([]byte(data))
	gameD.HashCode = hex.EncodeToString(hasher.Sum(nil))

	return gameD, nil
}
