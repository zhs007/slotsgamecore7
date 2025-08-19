package sgc7plugin

import (
	"context"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// MockPlugin - mock plugin
// 这个 plugin 是专门给测试使用的
// SetCache 一个 int 数组，后面 random 就会顺序返回这个数组作为随机数
// 如果 cache 为空，则返回 0
type MockPlugin struct {
	PluginBase
}

// NewMockPlugin - new a MockPlugin
func NewMockPlugin() *MockPlugin {
	fp := &MockPlugin{
		PluginBase: NewPluginBase(),
	}

	return fp
}

// Random - return [0, r)
func (fp *MockPlugin) Random(ctx context.Context, r int) (int, error) {
	if IsNoRNGCache {
		return 0, nil
	}

	var ci int
	if len(fp.Cache) > 0 {
		ci = fp.Cache[0]
		fp.Cache = fp.Cache[1:]
	} else {
		ci = 0
	}

	cr := ci % r

	fp.AddRngUsed(&sgc7utils.RngInfo{
		Bits:  cr,
		Range: r,
		Value: cr,
	})

	return cr, nil
}

// Init - initial
func (fp *MockPlugin) Init() {
}

// SetSeed - set a seed
func (fp *MockPlugin) SetSeed(seed int) {
}
