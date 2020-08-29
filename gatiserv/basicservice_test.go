package gatiserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_BasicService(t *testing.T) {
	bg := sgc7game.NewBasicGame()
	bs := NewBasicService(&bg)

	bs.Config()

	var iservice IService
	iservice = bs
	assert.NotNil(t, iservice, "Test_BasicService IService")

	t.Logf("Test_BasicService OK")
}
