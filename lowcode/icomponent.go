package lowcode

import (
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"google.golang.org/protobuf/types/known/anypb"
)

type FuncNewComponent func(name string) IComponent

type IComponent interface {
	// Init -
	Init(fn string, pool *GamePropertyPool) error
	// InitEx -
	InitEx(cfg any, pool *GamePropertyPool) error
	// OnNewGame - 这个一定要注意处理正确，为了节省cpu，没有主动处理componentData的该接口，如果确定需要，要自己调用
	OnNewGame(gameProp *GameProperty) error
	// OnNewStep -
	OnNewStep(gameProp *GameProperty) error
	// OnPlayGame - on playgame
	OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
		cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error
	// OnAsciiGame - outpur to asciigame
	OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error
	// OnStats -
	OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64)
	// NewComponentData -
	NewComponentData() IComponentData
	// EachUsedResults -
	EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult)
	// OnPlayGame - on playgame
	OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
		cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error
	// GetName -
	GetName() string

	// IsMask -
	IsMask() bool

	// IsRespin -
	IsRespin() bool
}
