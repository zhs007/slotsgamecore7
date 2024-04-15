package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type SceneStackData struct {
	Component  string
	SceneIndex int
	Scene      *sgc7game.GameScene
}

type SceneStack struct {
	Scenes       []*SceneStackData
	IsOtherScene bool
	// CacheScenes []*SceneStackData
}

// func (stack *SceneStack) Push(scene string, index int) {
// 	stack.Scenes = append(stack.Scenes, &SceneStackData{
// 		Component:  scene,
// 		SceneIndex: index,
// 	})
// }

func (stack *SceneStack) Push(scene string, index int, gs *sgc7game.GameScene) {
	ssd := &SceneStackData{
		Component:  scene,
		SceneIndex: index,
		Scene:      gs,
	}

	stack.Scenes = append(stack.Scenes, ssd)

	// if isNeedCache {
	// 	stack.CacheScenes = append(stack.CacheScenes, ssd)
	// }
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

			stack.Push("", 0, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			curpr.OtherScenes = append(curpr.OtherScenes, prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1])
			// prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]

			return stack.GetTopScene(curpr, prs)
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		stack.Push("", 0, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		curpr.Scenes = append(curpr.Scenes, prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1])
		// prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]

		return stack.GetTopScene(curpr, prs)
	}

	return stack.Scenes[len(stack.Scenes)-1]
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

	// for _, v := range stack.CacheScenes {
	// 	v.SceneIndex = len(pr.Scenes)

	// 	pr.Scenes = append(pr.Scenes, v.Scene)

	// 	stack.Scenes = append(stack.Scenes, v)
	// }

	// stack.CacheScenes = nil
}

func NewSceneStack(isOtherScene bool) *SceneStack {
	return &SceneStack{
		IsOtherScene: isOtherScene,
	}
}
