package relaxutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Config
	BaseReels *Reels `xml:"basereels"`
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
		CreditsPerBet:      25,
		FSAwarded:          15,
		FSAwardedRetrigger: 5,
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
	}

	tcfg := &TestConfig{
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

	err := SaveConfig("../unittestdata/relaxconfig.xml", tcfg)
	assert.NoError(t, err)

	t.Logf("Test_CalcScatter OK")
}
