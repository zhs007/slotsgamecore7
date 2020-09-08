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

	var iservice IService
	iservice = bs
	assert.NotNil(t, iservice, "Test_BasicService IService")

	t.Logf("Test_BasicService OK")
}
