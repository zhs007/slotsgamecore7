package lowcode

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type BasicComponent struct {
	UsedScenes  []int
	UsedResults []int
}

// AddScene -
func (basicComponent *BasicComponent) AddScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, sc *sgc7game.GameScene, tag string) {
	si := len(curpr.Scenes)
	basicComponent.UsedScenes = append(basicComponent.UsedScenes, si)

	curpr.Scenes = append(curpr.Scenes, sc)

	if tag != "" {
		gameProp.TagScene(curpr, tag, si)
	}
}

// AddResult -
func (basicComponent *BasicComponent) AddResult(curpr *sgc7game.PlayResult, ret *sgc7game.Result) {
	curpr.CashWin += int64(ret.CashWin)
	curpr.CoinWin += ret.CoinWin

	basicComponent.UsedResults = append(basicComponent.UsedResults, len(curpr.Results))

	curpr.Results = append(curpr.Results, ret)
}

func NewBasicComponent() *BasicComponent {
	return &BasicComponent{}
}
