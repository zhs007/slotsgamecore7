package relaxutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type TestConfig2 struct {
	Config
	BaseReels *Reels `xml:"baseReels"`
}

func Test_Config2(t *testing.T) {
	cfg := Config{
		GeneralComment: "General game information",
		SD:             1.56,
		ID:             123,
		Name:           "game001",
		ConfigVersion:  "1.0",
		RTPsComment:    "RTP infomation",
		RTPs: &FloatList{
			Vals: []float32{0.9607, 0.9608, 0.9609},
		},
		ConfigComment: "Configuration",
		Wilds: &StringList{
			Vals: []string{"WL"},
		},
	}

	pt, err := sgc7game.LoadPaytablesFromExcel("../unittestdata/paytables.xlsx")
	assert.NoError(t, err)

	cfg.Payouts = BuildPayouts(pt)

	tcfg := &TestConfig2{
		Config: cfg,
	}

	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)

	tcfg.BaseReels = BuildReels([]*sgc7game.ReelsData{reels}, pt)

	err = SaveConfig("../unittestdata/relaxconfig2.xml", tcfg, func(str string) string {
		str = strings.ReplaceAll(str, "></symbol>", "/>")
		str = strings.ReplaceAll(str, "></win>", "/>")

		return str
	})
	assert.NoError(t, err)

	t.Logf("Test_CalcScatter OK")
}
