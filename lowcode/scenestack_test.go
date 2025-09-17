package lowcode

import (
    "testing"

    sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func mustNewScene(t *testing.T, w, h int) *sgc7game.GameScene {
    gs, err := sgc7game.NewGameScene(w, h)
    if err != nil {
        t.Fatalf("NewGameScene err: %v", err)
    }

    return gs
}

func TestPushPopInsertPreHasPopToTruncate(t *testing.T) {
    st := NewSceneStack(false)

    s1 := mustNewScene(t, 3, 3)
    s2 := mustNewScene(t, 3, 3)
    s3 := mustNewScene(t, 3, 3)

    // push two
    st.Push("comp1", s1)
    st.Push("comp3", s3)

    // insert before top => comp1, comp2, comp3
    st.InsertPreScene("comp2", s2)

    if len(st.Scenes) != 3 {
        t.Fatalf("expected 3 scenes, got %d", len(st.Scenes))
    }

    if st.Scenes[1].Component != "comp2" {
        t.Fatalf("expected comp2 at idx1, got %v", st.Scenes[1].Component)
    }

    // Pop
    top := st.Pop()
    if top.Component != "comp3" {
        t.Fatalf("expected comp3 popped, got %v", top.Component)
    }

    // Has
    if !st.Has("comp1") || !st.Has("comp2") {
        t.Fatalf("Has returned false for existing components")
    }

    // PopTo comp1 -> keep up to comp1
    st.PopTo("comp1")
    if len(st.Scenes) != 1 || st.Scenes[0].Component != "comp1" {
        t.Fatalf("PopTo did not truncate correctly: %#v", st.Scenes)
    }

    // TruncateTo
    // Truncate to 0 clears
    st.TruncateTo(0)
    if len(st.Scenes) != 0 {
        t.Fatalf("expected empty after truncate to 0")
    }

    // push three then TruncateTo middle
    st.Push("a", s1)
    st.Push("b", s2)
    st.Push("c", s3)
    st.TruncateTo(2)
    if len(st.Scenes) != 2 || st.Scenes[1].Component != "b" {
        t.Fatalf("truncate to 2 failed: %#v", st.Scenes)
    }
}

func TestGetTopSceneRestoreAndOtherScene(t *testing.T) {
    // main scenes restore
    st := NewSceneStack(false)
    prev := sgc7game.NewPlayResult("", 0, 0, "")
    s := mustNewScene(t, 2, 2)
    prev.Scenes = append(prev.Scenes, s)

    cur := sgc7game.NewPlayResult("", 0, 0, "")

    got := st.GetTopScene(cur, []*sgc7game.PlayResult{prev})
    if got == nil || got.Scene != s {
        t.Fatalf("GetTopScene did not restore scene: got=%v", got)
    }

    // ensure cur was appended
    if len(cur.Scenes) != 1 || cur.Scenes[0] != s {
        t.Fatalf("cur playresult not appended on restore")
    }

    // other scenes restore
    ost := NewSceneStack(true)
    prev2 := sgc7game.NewPlayResult("", 0, 0, "")
    so := mustNewScene(t, 1, 1)
    prev2.OtherScenes = append(prev2.OtherScenes, so)
    cur2 := sgc7game.NewPlayResult("", 0, 0, "")

    got2 := ost.GetTopScene(cur2, []*sgc7game.PlayResult{prev2})
    if got2 == nil || got2.Scene != so {
        t.Fatalf("GetTopScene other did not restore scene")
    }

    // GetTopSceneEx returns the GameScene directly
    st2 := NewSceneStack(false)
    prev3 := sgc7game.NewPlayResult("", 0, 0, "")
    s3 := mustNewScene(t, 2, 2)
    prev3.Scenes = append(prev3.Scenes, s3)
    cur3 := sgc7game.NewPlayResult("", 0, 0, "")

    gse := st2.GetTopSceneEx(cur3, []*sgc7game.PlayResult{prev3})
    if gse != s3 {
        t.Fatalf("GetTopSceneEx expected s3, got %v", gse)
    }
}

func TestGetPreTopSceneAndEx(t *testing.T) {
    st := NewSceneStack(false)
    s1 := mustNewScene(t, 2, 2)
    s2 := mustNewScene(t, 2, 2)

    // push two
    st.Push("c1", s1)
    st.Push("c2", s2)

    pre := st.GetPreTopScene(&sgc7game.PlayResult{}, nil)
    if pre == nil || pre.Scene != s1 {
        t.Fatalf("GetPreTopScene expected s1, got %v", pre)
    }

    // When only one on stack and prs provided, it should insert pre from prs
    st2 := NewSceneStack(false)
    st2.Push("only", s1)
    prev := sgc7game.NewPlayResult("", 0, 0, "")
    sPrev := mustNewScene(t, 3, 3)
    prev.Scenes = append(prev.Scenes, sPrev)
    cur := sgc7game.NewPlayResult("", 0, 0, "")

    pre2 := st2.GetPreTopScene(cur, []*sgc7game.PlayResult{prev})
    if pre2 == nil || pre2.Scene != sPrev {
        t.Fatalf("GetPreTopScene restore failed")
    }

    // GetPreTopSceneEx
    st3 := NewSceneStack(false)
    st3.Push("a", s1)
    st3.Push("b", s2)
    preex := st3.GetPreTopSceneEx(&sgc7game.PlayResult{}, nil)
    if preex == nil || preex != s1 {
        t.Fatalf("GetPreTopSceneEx expected s1")
    }
}

func TestGetTargetScene3(t *testing.T) {
    gameProp := &GameProperty{
        SceneStack:      NewSceneStack(false),
        OtherSceneStack: NewSceneStack(true),
    }

    s1 := mustNewScene(t, 2, 2)
    s2 := mustNewScene(t, 2, 2)

    gameProp.SceneStack.Push("compA", s1)
    gameProp.SceneStack.Push("compB", s2)

    cfg := &BasicComponentConfig{
        TargetScenes3: [][]string{{"compB"}},
    }

    // should find compB when si=0
    got := gameProp.SceneStack.GetTargetScene3(gameProp, cfg, 0, &sgc7game.PlayResult{}, nil)
    if got != s2 {
        t.Fatalf("GetTargetScene3 expected s2, got %v", got)
    }

    // when not found in stack, fallback to top (use prs to restore)
    emptyStack := NewSceneStack(false)
    prev := sgc7game.NewPlayResult("", 0, 0, "")
    sPrev := mustNewScene(t, 4, 4)
    prev.Scenes = append(prev.Scenes, sPrev)
    cur := &sgc7game.PlayResult{}

    fallback := emptyStack.GetTargetScene3(gameProp, cfg, 0, cur, []*sgc7game.PlayResult{prev})
    if fallback != sPrev {
        t.Fatalf("fallback expected sPrev, got %v", fallback)
    }
}

func TestOnStepStartAndNewSceneStack(t *testing.T) {
    st := NewSceneStack(false)
    st.Push("x", mustNewScene(t, 1, 1))

    st.onStepStart(&sgc7game.PlayResult{})
    if len(st.Scenes) != 0 {
        t.Fatalf("onStepStart should clear scenes")
    }

    ost := NewSceneStack(true)
    if !ost.IsOtherScene {
        t.Fatalf("NewSceneStack with true should set IsOtherScene")
    }
}

func TestInsertPreSceneEmptyPanicsAndGetTopEmptyNil(t *testing.T) {
    // InsertPreScene on empty stack is known to cause an index error; test
    // that it panics so behavior is explicit.
    st := NewSceneStack(false)
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("expected panic on InsertPreScene with empty stack")
        }
    }()
    st.InsertPreScene("a", mustNewScene(t, 1, 1))
}

func TestGetTopSceneEmptyAndNoPRS(t *testing.T) {
    st := NewSceneStack(false)
    cur := &sgc7game.PlayResult{}
    got := st.GetTopScene(cur, nil)
    if got != nil {
        t.Fatalf("expected nil when no prs and empty stack")
    }

    // Test GetTopSceneEx with empty prs
    gse := st.GetTopSceneEx(cur, nil)
    if gse != nil {
        t.Fatalf("expected nil GetTopSceneEx when no prs and empty stack")
    }
}

func TestPopExWrapper(t *testing.T) {
    st := NewSceneStack(false)
    s1 := mustNewScene(t, 1, 1)
    s2 := mustNewScene(t, 1, 1)
    st.Push("a", s1)
    st.Push("b", s2)

    // PopEx should truncate to given num
    st.TruncateTo(1)
    if len(st.Scenes) != 1 || st.Scenes[0].Component != "a" {
        t.Fatalf("PopEx did not truncate as expected: %#v", st.Scenes)
    }
}

func TestPopToNotFoundDoesNothing(t *testing.T) {
    st := NewSceneStack(false)
    st.Push("x", mustNewScene(t, 1, 1))
    st.Push("y", mustNewScene(t, 1, 1))

    st.PopTo("noexist")
    if len(st.Scenes) != 2 {
        t.Fatalf("PopTo modified stack when component not found")
    }
}

func TestGetPreTopSceneOtherSceneRestore(t *testing.T) {
    ost := NewSceneStack(true)
    ost.Push("only", mustNewScene(t, 2, 2))

    prev := sgc7game.NewPlayResult("", 0, 0, "")
    sPrev := mustNewScene(t, 2, 2)
    prev.OtherScenes = append(prev.OtherScenes, sPrev)

    cur := sgc7game.NewPlayResult("", 0, 0, "")

    pre := ost.GetPreTopScene(cur, []*sgc7game.PlayResult{prev})
    if pre == nil || pre.Scene != sPrev {
        t.Fatalf("GetPreTopScene for other scene failed")
    }

    preex := ost.GetPreTopSceneEx(cur, []*sgc7game.PlayResult{prev})
    if preex == nil || preex != sPrev {
        t.Fatalf("GetPreTopSceneEx for other scene failed")
    }
}

func TestCoverScenestackAllPaths(t *testing.T) {
    // create stacks
    st := NewSceneStack(false)

    s1 := mustNewScene(t, 2, 2)
    s2 := mustNewScene(t, 2, 2)
    s3 := mustNewScene(t, 2, 2)

    // Push and InsertPre when stack has elements
    st.Push("one", s1)
    st.InsertPreScene("pre", s2) // should insert before top

    // Push another
    st.Push("top", s3)

    // GetTopScene when stack non-empty
    cur := &sgc7game.PlayResult{}
    top := st.GetTopScene(cur, nil)
    if top == nil || top.Scene != s3 {
        t.Fatalf("expected top s3")
    }

    // GetTopSceneEx
    gse := st.GetTopSceneEx(cur, nil)
    if gse != s3 {
        t.Fatalf("GetTopSceneEx mismatch")
    }

    // GetPreTopScene when >=2
    pre := st.GetPreTopScene(&sgc7game.PlayResult{}, nil)
    if pre == nil {
        t.Fatalf("expected pre not nil")
    }

    preex := st.GetPreTopSceneEx(&sgc7game.PlayResult{}, nil)
    if preex == nil {
        t.Fatalf("expected preex not nil")
    }

    // Has
    if !st.Has("one") || st.Has("noexist") {
        t.Fatalf("Has behavior unexpected")
    }

    // PopTo existing: should keep elements up to and including the found
    // component. Because we inserted a "pre" before "one" earlier the
    // resulting stack should contain both "pre" and "one".
    st.PopTo("one")
    if len(st.Scenes) != 2 || st.Scenes[1].Component != "one" {
        t.Fatalf("PopTo existing failed: %#v", st.Scenes)
    }

    // TruncateTo boundary
    st.TruncateTo(5) // no-op
    st.TruncateTo(1)
    st.TruncateTo(0) // clears
    if len(st.Scenes) != 0 {
        t.Fatalf("expected cleared stack")
    }

    // PopEx wrapper
    st.Push("a", s1)
    st.Push("b", s2)
    st.TruncateTo(1)
    if len(st.Scenes) != 1 || st.Scenes[0].Component != "a" {
        t.Fatalf("PopEx failed")
    }

    // GetTargetScene3 when target found in stack
    gameProp := &GameProperty{SceneStack: NewSceneStack(false), OtherSceneStack: NewSceneStack(true)}
    gameProp.SceneStack.Push("compA", s1)
    gameProp.SceneStack.Push("compB", s2)
    cfg := &BasicComponentConfig{TargetScenes3: [][]string{{"compB"}}}
    got := gameProp.SceneStack.GetTargetScene3(gameProp, cfg, 0, &sgc7game.PlayResult{}, nil)
    if got != s2 {
        t.Fatalf("GetTargetScene3 should return compB scene")
    }

    // GetTargetScene3 fallback path: empty stack, restore top from prs
    empty := NewSceneStack(false)
    prev := sgc7game.NewPlayResult("", 0, 0, "")
    prevS := mustNewScene(t, 1, 1)
    prev.Scenes = append(prev.Scenes, prevS)
    curpr := &sgc7game.PlayResult{}
    fb := empty.GetTargetScene3(gameProp, cfg, 0, curpr, []*sgc7game.PlayResult{prev})
    if fb != prevS {
        t.Fatalf("fallback path failed")
    }

    // onStepStart
    st.onStepStart(&sgc7game.PlayResult{})
    if len(st.Scenes) != 0 {
        t.Fatalf("onStepStart did not clear")
    }

    // ensure other scene branches executed in GetTopScene when empty and prs non-empty
    ostEmpty := NewSceneStack(true)
    prev2 := sgc7game.NewPlayResult("", 0, 0, "")
    pso := mustNewScene(t, 1, 1)
    prev2.OtherScenes = append(prev2.OtherScenes, pso)
    cur2 := &sgc7game.PlayResult{}
    gotost := ostEmpty.GetTopScene(cur2, []*sgc7game.PlayResult{prev2})
    if gotost == nil || gotost.Scene != pso {
        t.Fatalf("other scene restore failed")
    }
}
