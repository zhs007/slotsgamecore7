package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type FuncNewComponent func(fn string) IComponent

type IComponent interface {
	// playgame
	OnPlayGame(curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin, cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error
	// pay
	OnPay(curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin, cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error
}
