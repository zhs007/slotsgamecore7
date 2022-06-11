package relaxutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Config
	CreditsPerBet      int      `xml:"CREDITS_PER_BET"`
	FSAwarded          int      `xml:"FS_AWARDED"`
	FSAwardedRetrigger int      `xml:"FS_AWARDED_RETRIGGER"`
	BaseReels          *Reels   `xml:"basereels"`
	BaseReelWeights    *Weights `xml:"baseReelWeights"`
}

func Test_Config(t *testing.T) {
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
		Symbols: &Symbols{
			Symbols: []*Symbol{
				{
					Name: "A",
					Val:  1,
				},
				{
					Name: "B",
					Val:  2,
				},
			},
		},
		Wilds: &StringList{
			Vals: []string{"WL"},
		},
		Payouts: &Payouts{
			Payouts: []*Payout{
				{
					Name: "A",
					Wins: []*PayoutWin{
						{
							Count:  3,
							Payout: 100,
						},
						{
							Count:  4,
							Payout: 200,
						},
						{
							Count:  5,
							Payout: 300,
						},
					},
				},
				{
					Name: "B",
					Wins: []*PayoutWin{
						{
							Count:  3,
							Payout: 150,
						},
						{
							Count:  4,
							Payout: 250,
						},
						{
							Count:  5,
							Payout: 350,
						},
					},
				},
			},
		},
		PayingSymbols: &StringList{
			Vals: []string{"A", "B", "C", "D"},
		},
	}

	tcfg := &TestConfig{
		cfg,
		25,
		15,
		5,
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
		&Weights{
			[]*IntList{
				{
					Vals: []int{0, 15},
				},
				{
					Vals: []int{1, 25},
				},
			},
		},
	}

	err := SaveConfig("../unittestdata/relaxconfig.xml", tcfg, func(str string) string {
		str = strings.ReplaceAll(str, "></symbol>", "/>")
		str = strings.ReplaceAll(str, "></win>", "/>")

		return str
	})
	assert.NoError(t, err)

	t.Logf("Test_CalcScatter OK")
}
