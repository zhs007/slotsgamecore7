// Package lowcode provides runtime helpers for low-code game components.
//
// SceneStack manages a stack of scenes that components can push/pop during
// game execution. It also supports restoring scenes from previous
// PlayResult entries when the current stack is empty.
package lowcode

import (
	"slices"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// SceneStackData pairs a component name (or placeholder) with a GameScene
// instance. Component is the logical owner/identifier of the scene; Scene
// holds the concrete scene data.
type SceneStackData struct {
	Component string
	Scene     *sgc7game.GameScene
}

// SceneStack holds a slice of SceneStackData and a flag indicating whether
// this stack is used for "other scenes" (IsOtherScene=true) or the main
// scenes. Many helper methods may push/pop or restore scenes from cached
// PlayResult values when the stack is temporarily empty.
type SceneStack struct {
	Scenes       []*SceneStackData
	IsOtherScene bool
}

// Push pushes a new SceneStackData onto the stack.
// 'scene' is the component name associated with the GameScene 'gs'.
func (stack *SceneStack) Push(scene string, gs *sgc7game.GameScene) {
	ssd := &SceneStackData{
		Component: scene,
		Scene:     gs,
	}

	stack.Scenes = append(stack.Scenes, ssd)
}

// InsertPreScene inserts a SceneStackData just before the current top of
// the stack. If the stack has N elements, the new element will be placed
// at index N-1 (i.e., immediately below the top). Note: calling this when
// the stack is empty will cause an index error; callers should ensure the
// stack length if necessary.
func (stack *SceneStack) InsertPreScene(scene string, gs *sgc7game.GameScene) {
	ssd := &SceneStackData{
		Component: scene,
		Scene:     gs,
	}

	// Insert before the last element (top). Keep behavior consistent with
	// the original implementation.
	stack.Scenes = slices.Insert(stack.Scenes, len(stack.Scenes)-1, ssd)
}

// Pop removes and returns the top SceneStackData. Returns nil if the
// stack is empty.
func (stack *SceneStack) Pop() *SceneStackData {
	if len(stack.Scenes) == 0 {
		return nil
	}

	ssd := stack.Scenes[len(stack.Scenes)-1]

	stack.Scenes = stack.Scenes[:len(stack.Scenes)-1]

	return ssd
}

// GetTopScene returns the top SceneStackData. If the stack is empty, it
// attempts to restore a scene from the last entry in prs and pushes that
// scene onto the stack. In that restoration path, it also appends the
// restored scene into curpr.Scenes or curpr.OtherScenes (side-effect).
// Returns nil if neither the stack nor prs provides a scene.
func (stack *SceneStack) GetTopScene(curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *SceneStackData {
	if len(stack.Scenes) == 0 {
		if len(prs) == 0 {
			return nil
		}

		// When the stack is empty, restore the most recent scene from the
		// provided play results. Note this has side-effects: it pushes the
		// restored scene into the stack and also appends it to curpr.
		if stack.IsOtherScene {
			if len(prs[len(prs)-1].OtherScenes) == 0 {
				return nil
			}

			restored := prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1]
			stack.Push("", restored)
			curpr.OtherScenes = append(curpr.OtherScenes, restored)

			return stack.Scenes[len(stack.Scenes)-1]
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		restored := prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]
		stack.Push("", restored)
		curpr.Scenes = append(curpr.Scenes, restored)

		return stack.Scenes[len(stack.Scenes)-1]
	}

	return stack.Scenes[len(stack.Scenes)-1]
}

// GetTopSceneEx is like GetTopScene but returns the GameScene instance
// directly. It shares the same restoration side-effects when the stack is
// empty.
func (stack *SceneStack) GetTopSceneEx(curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) *sgc7game.GameScene {
	if len(stack.Scenes) == 0 {
		if len(prs) == 0 {
			return nil
		}

		if stack.IsOtherScene {
			if len(prs[len(prs)-1].OtherScenes) == 0 {
				return nil
			}

			restored := prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1]
			stack.Push("", restored)
			curpr.OtherScenes = append(curpr.OtherScenes, restored)

			return stack.Scenes[len(stack.Scenes)-1].Scene
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		restored := prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]
		stack.Push("", restored)
		curpr.Scenes = append(curpr.Scenes, restored)

		return stack.Scenes[len(stack.Scenes)-1].Scene
	}

	return stack.Scenes[len(stack.Scenes)-1].Scene
}

// GetPreTopScene returns the scene just below the top. If the stack has at
// least two entries, it returns the second-from-top element. If the stack
// has exactly one entry and prs is provided, it inserts a pre-scene from
// the last PlayResult (with the same side-effects as GetTopScene) so the
// pre-top becomes available. Returns nil if no pre-top scene can be
// determined.
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

			restored := prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1]
			// Insert restored scene just below the top.
			stack.InsertPreScene("", restored)
			curpr.OtherScenes = append(curpr.OtherScenes, restored)

			return stack.Scenes[len(stack.Scenes)-2]
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		restored := prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]
		stack.InsertPreScene("", restored)
		curpr.Scenes = append(curpr.Scenes, restored)

		return stack.Scenes[len(stack.Scenes)-2]
	}

	return nil
}

// GetPreTopSceneEx is like GetPreTopScene but returns the GameScene
// instance directly.
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

			restored := prs[len(prs)-1].OtherScenes[len(prs[len(prs)-1].OtherScenes)-1]
			stack.InsertPreScene("", restored)
			curpr.OtherScenes = append(curpr.OtherScenes, restored)

			return stack.Scenes[len(stack.Scenes)-2].Scene
		}

		if len(prs[len(prs)-1].Scenes) == 0 {
			return nil
		}

		restored := prs[len(prs)-1].Scenes[len(prs[len(prs)-1].Scenes)-1]
		stack.InsertPreScene("", restored)
		curpr.Scenes = append(curpr.Scenes, restored)

		return stack.Scenes[len(stack.Scenes)-2].Scene
	}

	return nil
}

// Has reports whether a scene with the given component name exists in the
// stack. Comparison is by exact string equality of the Component field.
func (stack *SceneStack) Has(scene string) bool {
	for _, v := range stack.Scenes {
		if v.Component == scene {
			return true
		}
	}

	return false
}

// PopTo truncates the stack so that the top becomes the first occurrence
// of the provided component name (searching from the top down). If the
// component is not found, the stack remains unchanged.
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

// TruncateTo reduces the stack to at most 'num' entries. If num <= 0 the
// stack is cleared. If num >= current length, no-op.
func (stack *SceneStack) TruncateTo(num int) {
	if num <= 0 {
		stack.Scenes = nil
		return
	}

	if num >= len(stack.Scenes) {
		return
	}
    
	stack.Scenes = stack.Scenes[:num]
}

// GetTargetScene3 searches the stack from top to bottom for the first
// component whose name appears in the configured target lists in
// BasicComponentConfig (TargetScenes3 or TargetOtherScenes3 depending on
// IsOtherScene). If none match, it falls back to the current top scene
// (restoring from prs if needed) and returns that scene's GameScene.
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

// onStepStart is invoked when a new step begins; it clears the SceneStack
// to start fresh for the new step.
func (stack *SceneStack) onStepStart(_ *sgc7game.PlayResult) {
	stack.Scenes = nil
}

// NewSceneStack constructs a SceneStack. Pass isOtherScene=true to create
// a stack that operates on OtherScenes (used for auxiliary scene types).
func NewSceneStack(isOtherScene bool) *SceneStack {
	return &SceneStack{
		IsOtherScene: isOtherScene,
	}
}
