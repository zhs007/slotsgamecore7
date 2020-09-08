package sgc7plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PluginsMgr(t *testing.T) {
	mgr := NewPluginsMgr(func() IPlugin {
		return NewBasicPlugin()
	})

	p0 := mgr.NewPlugin()
	assert.Equal(t, len(mgr.plugins), 0)

	mgr.FreePlugin(p0)
	assert.Equal(t, len(mgr.plugins), 1)

	p1 := mgr.NewPlugin()
	assert.Equal(t, len(mgr.plugins), 0)
	assert.Equal(t, p0, p1)

	p2 := mgr.NewPlugin()
	assert.Equal(t, len(mgr.plugins), 0)

	mgr.FreePlugin(p2)
	mgr.FreePlugin(p1)
	assert.Equal(t, len(mgr.plugins), 2)

	p3 := mgr.NewPlugin()
	assert.Equal(t, len(mgr.plugins), 1)
	assert.Equal(t, p3, p2)

	p4 := mgr.NewPlugin()
	assert.Equal(t, len(mgr.plugins), 0)
	assert.Equal(t, p4, p0)

	p5 := mgr.NewPlugin()
	assert.Equal(t, len(mgr.plugins), 0)
	assert.NotNil(t, p5)

	t.Logf("Test_PluginsMgr OK")
}
