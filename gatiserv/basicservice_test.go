package gatiserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_BasicService(t *testing.T) {
	bg := sgc7game.NewBasicGame(func() sgc7plugin.IPlugin {
		return sgc7plugin.NewBasicPlugin()
	})
	bs := NewBasicService(bg)

	bs.Config()

	bs.Initialize()

	bs.Validate(&ValidateParams{})

	bs.Play(&PlayParams{})

	plugin := bg.NewPlugin()
	err := bs.ProcCheat(plugin, "1,2,3")
	assert.NoError(t, err)

	cr0, err := plugin.Random(100)
	assert.NoError(t, err)
	assert.Equal(t, cr0, 1)

	cr1, err := plugin.Random(100)
	assert.NoError(t, err)
	assert.Equal(t, cr1, 2)

	cr2, err := plugin.Random(100)
	assert.NoError(t, err)
	assert.Equal(t, cr2, 3)

	cr3, err := plugin.Random(100)
	assert.NoError(t, err)
	assert.Equal(t, func() bool { return cr3 >= 0 }(), true)

	var iservice IService
	iservice = bs
	assert.NotNil(t, iservice, "Test_BasicService IService")

	t.Logf("Test_BasicService OK")
}
