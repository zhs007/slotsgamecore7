package lowcode

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_NewGenSymbolValsInReels(t *testing.T) {
    mgr := NewComponentMgr()

    cfg := &ComponentConfig{
        Name: "test1",
        Type: GenSymbolValsInReelsTypeName,
    }

    comp := mgr.NewComponent(cfg)
    assert.NotNil(t, comp)
}
