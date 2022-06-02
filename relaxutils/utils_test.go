package relaxutils

import (
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
		RTP:            0.9607,
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
		cfg,
		&Reels{
			Tables: []*Table{
				{
					TableComment: "table 0",
					Reel: []*StringList{
						{
							Vals: []string{"A", "B", "A"},
						},
						{
							Vals: []string{"B", "B", "A"},
						},
					},
				},
				{
					TableComment: "table 1",
					Reel: []*StringList{
						{
							Vals: []string{"A", "B", "A"},
						},
						{
							Vals: []string{"B", "B", "A"},
						},
					},
				},
			},
		},
	}

	err = SaveConfig("../unittestdata/relaxconfig2.xml", tcfg)
	assert.NoError(t, err)

	t.Logf("Test_CalcScatter OK")
}
