package sgc7rtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

func Test_RTP(t *testing.T) {
	rtp := NewRTP()
	assert.NotNil(t, rtp)

	bg := NewRTPGameMod("bg")
	InitGameMod(bg, nil, nil, []int{0, 1, 2, 3, 4, 5, 6, 7}, []int{1, 2, 3, 4, 5})
	rtp.Root.AddChild("bg", bg)

	fg := NewRTPGameMod("fg")
	InitGameMod(fg, nil, nil, []int{0, 1, 2, 3, 4, 5, 6, 7}, []int{1, 2, 3, 4, 5})
	rtp.Root.AddChild("fg", fg)

	rtp.CalcRTP()

	pr0 := &sgc7game.PlayResult{
		CurGameMod: "bg",
		CashWin:    300,
		IsFinish:   true,
	}

	pr0.Results = append(pr0.Results, &sgc7game.Result{
		Symbol:     0,
		CashWin:    200,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr0.Results = append(pr0.Results, &sgc7game.Result{
		Symbol:     1,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1, 3, 1},
		SymbolNums: 4,
	})

	pr0.Results = append(pr0.Results, &sgc7game.Result{
		Symbol:     1,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	rtp.Bet(100)
	rtp.OnResult(pr0)

	assert.Equal(t, rtp.BetNums, int64(1))
	assert.Equal(t, rtp.TotalBet, int64(100))
	assert.Equal(t, rtp.Root.TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].MapChildren["3"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["3"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["4"].TotalWin, int64(50))

	pr1 := &sgc7game.PlayResult{
		CurGameMod: "fg",
		CashWin:    0,
		IsFinish:   true,
	}

	pr1.Results = append(pr1.Results, &sgc7game.Result{
		Symbol:     2,
		CashWin:    200,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr1.Results = append(pr1.Results, &sgc7game.Result{
		Symbol:     3,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1, 3, 1},
		SymbolNums: 4,
	})

	pr1.Results = append(pr1.Results, &sgc7game.Result{
		Symbol:     4,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	rtp.Bet(100)
	rtp.OnResult(pr1)

	assert.Equal(t, rtp.BetNums, int64(2))
	assert.Equal(t, rtp.TotalBet, int64(200))
	assert.Equal(t, rtp.Root.TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].MapChildren["3"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["3"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["4"].TotalWin, int64(50))

	pr2 := &sgc7game.PlayResult{
		CurGameMod: "fg",
		CashWin:    300,
		IsFinish:   true,
	}

	pr2.Results = append(pr2.Results, &sgc7game.Result{
		Symbol:     2,
		CashWin:    200,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr2.Results = append(pr2.Results, &sgc7game.Result{
		Symbol:     3,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1, 3, 1},
		SymbolNums: 4,
	})

	pr2.Results = append(pr2.Results, &sgc7game.Result{
		Symbol:     4,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr0.IsFinish = false

	rtp.Bet(100)
	rtp.OnResult(pr0)
	rtp.OnResult(pr2)

	assert.Equal(t, rtp.BetNums, int64(3))
	assert.Equal(t, rtp.TotalBet, int64(300))
	assert.Equal(t, rtp.Root.TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.TotalWin, int64(900))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TotalWin, int64(600))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].TotalWin, int64(400))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].MapChildren["3"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["0"].MapChildren["3"].TotalWin, int64(400))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].TriggerNums, int64(4))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["3"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["3"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["4"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["1"].MapChildren["4"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["fg"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["2"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["2"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["2"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["2"].MapChildren["3"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["3"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["3"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["3"].MapChildren["4"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["4"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["4"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["4"].MapChildren["3"].TotalWin, int64(50))

	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["2"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["2"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["3"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["3"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["4"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["4"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["5"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["5"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["0"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["0"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["1"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["1"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["5"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["5"].TotalWin, int64(0))

	rtp.Save2CSV("../unittestdata/rtptest.csv")

	issame := sgc7utils.IsSameFile("../unittestdata/rtptest.csv", "../unittestdata/rtptestok.csv")
	assert.Equal(t, issame, true)

	t.Logf("Test_RTP OK")
}

func Test_RTPTags(t *testing.T) {
	type gameparams struct {
		isSpecial bool
	}

	rtp := NewRTP()
	assert.NotNil(t, rtp)

	bg := NewRTPGameMod("bg")
	InitGameMod(bg, []string{"normal", "special"}, []FuncOnResult{
		func(node *RTPNode, pr *sgc7game.PlayResult) bool {
			gp, isok := pr.CurGameModParams.(*gameparams)
			if isok && gp != nil {
				if !gp.isSpecial && node.TagName == "normal" {
					node.TriggerNums++
					node.TotalWin += pr.CashWin

					return true
				}
			}

			return false
		},
		func(node *RTPNode, pr *sgc7game.PlayResult) bool {
			gp, isok := pr.CurGameModParams.(*gameparams)
			if isok && gp != nil {
				if gp.isSpecial && node.TagName == "special" {
					node.TriggerNums++
					node.TotalWin += pr.CashWin

					return true
				}
			}

			return false
		},
	}, []int{0, 1, 2, 3, 4, 5, 6, 7}, []int{1, 2, 3, 4, 5})
	rtp.Root.AddChild("bg", bg)

	fg := NewRTPGameMod("fg")
	InitGameMod(fg, []string{"normal", "special"}, []FuncOnResult{
		func(node *RTPNode, pr *sgc7game.PlayResult) bool {
			gp, isok := pr.CurGameModParams.(*gameparams)
			if isok && gp != nil {
				if !gp.isSpecial && node.TagName == "normal" {
					node.TriggerNums++
					node.TotalWin += pr.CashWin

					return true
				}
			}

			return false
		},
		func(node *RTPNode, pr *sgc7game.PlayResult) bool {
			gp, isok := pr.CurGameModParams.(*gameparams)
			if isok && gp != nil {
				if gp.isSpecial && node.TagName == "special" {
					node.TriggerNums++
					node.TotalWin += pr.CashWin

					return true
				}
			}

			return false
		},
	}, []int{0, 1, 2, 3, 4, 5, 6, 7}, []int{1, 2, 3, 4, 5})
	rtp.Root.AddChild("fg", fg)

	rtp.AddHitRateNode("special", func(node *HitRateNode, pr *sgc7game.PlayResult) bool {
		gp, isok := pr.CurGameModParams.(*gameparams)
		if isok && gp != nil {
			if gp.isSpecial && node.TagName == "special" {
				node.TriggerNums++

				for range pr.Results {
					node.TotalNums++
				}

				return true
			}
		}

		return false
	})

	rtp.CalcRTP()

	pr0 := &sgc7game.PlayResult{
		CurGameMod: "bg",
		CashWin:    300,
		IsFinish:   true,
		CurGameModParams: &gameparams{
			isSpecial: false,
		},
	}

	pr0.Results = append(pr0.Results, &sgc7game.Result{
		Symbol:     0,
		CashWin:    200,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr0.Results = append(pr0.Results, &sgc7game.Result{
		Symbol:     1,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1, 3, 1},
		SymbolNums: 4,
	})

	pr0.Results = append(pr0.Results, &sgc7game.Result{
		Symbol:     1,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	rtp.Bet(100)
	rtp.OnResult(pr0)

	assert.Equal(t, rtp.BetNums, int64(1))
	assert.Equal(t, rtp.TotalBet, int64(100))
	assert.Equal(t, rtp.Root.TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].MapChildren["3"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["3"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["4"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["special"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["special"].TotalWin, int64(0))

	pr1 := &sgc7game.PlayResult{
		CurGameMod: "fg",
		CashWin:    0,
		IsFinish:   true,
		CurGameModParams: &gameparams{
			isSpecial: false,
		},
	}

	pr1.Results = append(pr1.Results, &sgc7game.Result{
		Symbol:     2,
		CashWin:    200,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr1.Results = append(pr1.Results, &sgc7game.Result{
		Symbol:     3,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1, 3, 1},
		SymbolNums: 4,
	})

	pr1.Results = append(pr1.Results, &sgc7game.Result{
		Symbol:     4,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	rtp.Bet(100)
	rtp.OnResult(pr1)

	assert.Equal(t, rtp.BetNums, int64(2))
	assert.Equal(t, rtp.TotalBet, int64(200))
	assert.Equal(t, rtp.Root.TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].MapChildren["3"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["3"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["4"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["special"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["special"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].TotalWin, int64(0))

	pr2 := &sgc7game.PlayResult{
		CurGameMod: "fg",
		CashWin:    300,
		IsFinish:   true,
		CurGameModParams: &gameparams{
			isSpecial: true,
		},
	}

	pr2.Results = append(pr2.Results, &sgc7game.Result{
		Symbol:     2,
		CashWin:    200,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr2.Results = append(pr2.Results, &sgc7game.Result{
		Symbol:     3,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1, 3, 1},
		SymbolNums: 4,
	})

	pr2.Results = append(pr2.Results, &sgc7game.Result{
		Symbol:     4,
		CashWin:    50,
		Pos:        []int{0, 1, 1, 1, 2, 1},
		SymbolNums: 3,
	})

	pr0.IsFinish = false

	rtp.Bet(100)
	rtp.OnResult(pr0)
	rtp.OnResult(pr2)

	assert.Equal(t, rtp.BetNums, int64(3))
	assert.Equal(t, rtp.TotalBet, int64(300))
	assert.Equal(t, rtp.Root.TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.TotalWin, int64(900))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].TotalWin, int64(600))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].TotalWin, int64(400))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].MapChildren["3"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["0"].MapChildren["3"].TotalWin, int64(400))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].TriggerNums, int64(4))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["3"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["3"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["4"].TriggerNums, int64(2))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["normal"].MapChildren["1"].MapChildren["4"].TotalWin, int64(100))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["special"].TriggerNums, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["special"].TotalWin, int64(0))
	assert.Equal(t, rtp.Root.MapChildren["fg"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].TotalWin, int64(300))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["2"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["2"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["2"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["2"].MapChildren["3"].TotalWin, int64(200))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["3"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["3"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["3"].MapChildren["4"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["4"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["4"].TotalWin, int64(50))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["4"].MapChildren["3"].TriggerNums, int64(1))
	assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["special"].MapChildren["4"].MapChildren["3"].TotalWin, int64(50))

	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["2"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["2"].TotalWin, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["3"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["3"].TotalWin, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["4"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["4"].TotalWin, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["5"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["bg"].MapChildren["5"].TotalWin, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["0"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["0"].TotalWin, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["1"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["1"].TotalWin, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["5"].TriggerNums, int64(0))
	// assert.Equal(t, rtp.Root.MapChildren["fg"].MapChildren["5"].TotalWin, int64(0))

	rtp.Save2CSV("../unittestdata/rtptesttag.csv")

	issame := sgc7utils.IsSameFile("../unittestdata/rtptesttag.csv", "../unittestdata/rtptesttagok.csv")
	assert.Equal(t, issame, true)

	rtp1 := rtp.Clone()
	assert.NotNil(t, rtp1)

	rtp.Add(rtp1)

	rtp.Save2CSV("../unittestdata/rtptesttag1.csv")

	issame = sgc7utils.IsSameFile("../unittestdata/rtptesttag1.csv", "../unittestdata/rtptesttag1ok.csv")
	assert.Equal(t, issame, true)

	t.Logf("Test_RTPTags OK")
}
