package gatiserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// IService - service
type IService interface {
	// Config - get configuration
	Config() *sgc7game.Config
}
