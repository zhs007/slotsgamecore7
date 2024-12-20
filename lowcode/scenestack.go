package lowcode

import (
	"slices"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type SceneStackData struct {
	Component string
	Scene     *sgc7game.GameScene
}

type SceneStack struct {
	Scenes       []*SceneStackData
	IsOtherScene bool
}

func (stack *SceneStack) Push(scene string, gs *sgc7game.GameScene) {
	ssd := &SceneStackData{
		Component: scene,
		Scene:     gs,
	}

	stack.Scenes = append(stack.Scenes, ssd)
}

func (stack *SceneStack) InsertPreScene(scene string, gs *sgc7game.GameScene) {
	ssd := &SceneStackData{
		Component: scene,
		Scene:     gs,
	}

	stack.Scenes = slices.Insert(stack.Scenes, len(stack.Scenes)-1, ssd)
}

func (stack *SceneStack) Pop() *SceneStackData {
	if len(stack.Scenes) == 0 {
		return nil
	}

	ssd := stack.Scenes[len(stack.Scenes)-1]

	stack.Scenes = stack.Scenes[:len(stack.Scenes)-1]

	return ssd
}

func (stack *SceneStack) GetTopScene(curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *SceneStackData {
	if len(stack.Scenes) == 0 {
		if len(prs) == 0 {
			return nil
		}

		if stack.IsOtherScene {
			if len(prs[len(prs)-1].OtherScenes) == 0 {
				return nil
			}

			stack.Push("", prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			curpr.OtherScenes = append(curpr.OtherScenes, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			// prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]

			return stack.Scenes[len(stack.Scenes)-1]
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		stack.Push("", prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		curpr.Scenes = append(curpr.Scenes, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		// prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]

		return stack.Scenes[len(stack.Scenes)-1]
	}

	return stack.Scenes[len(stack.Scenes)-1]
}

func (stack *SceneStack) GetTopSceneEx(curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *sgc7game.GameScene {
	if len(stack.Scenes) == 0 {
		if len(prs) == 0 {
			return nil
		}

		if stack.IsOtherScene {
			if len(prs[len(prs)-1].OtherScenes) == 0 {
				return nil
			}

			stack.Push("", prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			curpr.OtherScenes = append(curpr.OtherScenes, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			// prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]

			return stack.Scenes[len(stack.Scenes)-1].Scene
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		stack.Push("", prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		curpr.Scenes = append(curpr.Scenes, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		// prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]

		return stack.Scenes[len(stack.Scenes)-1].Scene
	}

	return stack.Scenes[len(stack.Scenes)-1].Scene
}

func (stack *SceneStack) GetPreTopScene(curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *SceneStackData {
	if len(stack.Scenes) >= 2 {
		return stack.Scenes[len(stack.Scenes)-2]
	} else if len(stack.Scenes) == 1 {
		if len(prs) == 0 {
			return nil
		}

		if stack.IsOtherScene {
			if len(prs[len(prs)-1].OtherScenes) == 0 {
				return nil
			}

			stack.InsertPreScene("", prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			// stack.Push("", 0, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			curpr.OtherScenes = append(curpr.OtherScenes, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])

			return stack.Scenes[len(stack.Scenes)-2]
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		stack.InsertPreScene("", prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		// stack.Push("", 0, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		curpr.Scenes = append(curpr.Scenes, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])

		return stack.Scenes[len(stack.Scenes)-2]
	}

	return nil
}

func (stack *SceneStack) GetPreTopSceneEx(curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *sgc7game.GameScene {
	if len(stack.Scenes) >= 2 {
		return stack.Scenes[len(stack.Scenes)-2].Scene
	} else if len(stack.Scenes) == 1 {
		if len(prs) == 0 {
			return nil
		}

		if stack.IsOtherScene {
			if len(prs[len(prs)-1].OtherScenes) == 0 {
				return nil
			}

			stack.InsertPreScene("", prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			// stack.Push("", 0, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			curpr.OtherScenes = append(curpr.OtherScenes, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])

			return stack.Scenes[len(stack.Scenes)-2].Scene
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		stack.InsertPreScene("", prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		// stack.Push("", 0, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		curpr.Scenes = append(curpr.Scenes, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])

		return stack.Scenes[len(stack.Scenes)-2].Scene
	}

	return nil
}

func (stack *SceneStack) Has(scene string) bool {
	for _, v := range stack.Scenes {
		if v.Component == scene {
			return true
		}
	}

	return false
}

func (stack *SceneStack) PopTo(scene string) {
	maxi := -1
	for i := len(stack.Scenes) - 1; i >= 0; i-- {
		if scene == stack.Scenes[i].Component {
			maxi = i

			break
		}
	}

	if maxi >= 0 {
		stack.Scenes = stack.Scenes[:maxi+1]
	}
}

func (stack *SceneStack) PopEx(num int) {
	stack.Scenes = stack.Scenes[:num]
}

func (stack *SceneStack) GetTargetScene3(gameProp *GameProperty, basicCfg *BasicComponentConfig, si int, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *sgc7game.GameScene {
	if stack.IsOtherScene {
		if len(basicCfg.TargetOtherScenes3) > si {
			for i := len(stack.Scenes) - 1; i >= 0; i-- {
				ci := goutils.IndexOfStringSlice(basicCfg.TargetOtherScenes3[si], stack.Scenes[i].Component, 0)
				if ci >= 0 {
					return stack.Scenes[i].Scene
				}
			}
		}
	} else {
		if len(basicCfg.TargetScenes3) > si {
			for i := len(stack.Scenes) - 1; i >= 0; i-- {
				ci := goutils.IndexOfStringSlice(basicCfg.TargetScenes3[si], stack.Scenes[i].Component, 0)
				if ci >= 0 {
					return stack.Scenes[i].Scene
				}
			}
		}
	}

	ssd := stack.GetTopScene(curpr, prs)
	if ssd == nil {
		return nil
	}

	return ssd.Scene
}

func (stack *SceneStack) onStepStart(_ *sgc7game.PlayResult) {
	stack.Scenes = nil
}

func NewSceneStack(isOtherScene bool) *SceneStack {
	return &SceneStack{
		IsOtherScene: isOtherScene,
	}
}
