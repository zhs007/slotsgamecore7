package lowcode

import (
    "testing"

    sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// TestSceneStackCoverage executes scenestack functions to increase file
// coverage for scenestack.go. Run this test specifically to measure the
// file's coverage.
func TestSceneStackCoverage(t *testing.T) {
    gs1 := mustNewScene(t, 2, 2)
    gs2 := mustNewScene(t, 2, 2)
    gs3 := mustNewScene(t, 2, 2)

    st := NewSceneStack(false)

    // push and pop
    st.Push("c1", gs1)
    st.Push("c2", gs2)
    if st.Pop().Scene != gs2 {
        t.Fatalf("pop mismatch")
    }

    // insert pre when stack has elements
    st.Push("top", gs3)
    st.InsertPreScene("mid", gs2)

    // access top and pretop
    cur := &sgc7game.PlayResult{}
    if st.GetTopScene(cur, nil).Scene != gs3 {
        t.Fatalf("GetTopScene mismatch")
    }
    if st.GetTopSceneEx(cur, nil) != gs3 {
        t.Fatalf("GetTopSceneEx mismatch")
    }
    if st.GetPreTopScene(cur, nil) == nil {
        t.Fatalf("GetPreTopScene nil")
    }
    if st.GetPreTopSceneEx(cur, nil) == nil {
        t.Fatalf("GetPreTopSceneEx nil")
    }

    // Has
    if !st.Has("c1") {
        t.Fatalf("Has failed")
    }

    // PopTo
    st.PopTo("c1")
    if len(st.Scenes) == 0 || st.Scenes[0].Component != "c1" {
        t.Fatalf("PopTo failed")
    }

    // TruncateTo and PopEx
    st.TruncateTo(5)
    st.TruncateTo(1)
    st.TruncateTo(0)
    st.Push("a", gs1)
    st.Push("b", gs2)
    st.TruncateTo(1)

    // GetTargetScene3 match
    gp := &GameProperty{SceneStack: NewSceneStack(false), OtherSceneStack: NewSceneStack(true)}
    gp.SceneStack.Push("a", gs1)
    gp.SceneStack.Push("b", gs2)
    cfg := &BasicComponentConfig{TargetScenes3: [][]string{{"b"}}}
    if gp.SceneStack.GetTargetScene3(gp, cfg, 0, &sgc7game.PlayResult{}, nil) != gs2 {
        t.Fatalf("GetTargetScene3 match failed")
    }

    // GetTargetScene3 fallback
    empty := NewSceneStack(false)
    pr := sgc7game.NewPlayResult("", 0, 0, "")
    ps := mustNewScene(t, 1, 1)
    pr.Scenes = append(pr.Scenes, ps)
    if empty.GetTargetScene3(gp, cfg, 0, &sgc7game.PlayResult{}, []*sgc7game.PlayResult{pr}) != ps {
        t.Fatalf("GetTargetScene3 fallback failed")
    }

    // onStepStart and NewSceneStack
    st.onStepStart(&sgc7game.PlayResult{})
    if len(st.Scenes) != 0 {
        t.Fatalf("onStepStart failed to clear")
    }

    if !NewSceneStack(true).IsOtherScene {
        t.Fatalf("NewSceneStack flag failed")
    }
}
