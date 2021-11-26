package sgc7pbutils

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// BuildPBLineData - *sgc7game.LineData -> *sgc7pb.LinesData
func BuildPBLineData(ld *sgc7game.LineData) *sgc7pb.LinesData {
	pbld := &sgc7pb.LinesData{}

	for _, l := range ld.Lines {
		pbl := &sgc7pb.Row{}

		for _, v := range l {
			pbl.Values = append(pbl.Values, int32(v))
		}

		pbld.Lines = append(pbld.Lines, pbl)
	}

	return pbld
}

// BuildPBReelsData - *sgc7game.ReelsData -> *sgc7pb.ReelsData
func BuildPBReelsData(rd *sgc7game.ReelsData) *sgc7pb.ReelsData {
	pbrd := &sgc7pb.ReelsData{}

	for _, l := range rd.Reels {
		pbl := &sgc7pb.Column{}

		for _, v := range l {
			pbl.Values = append(pbl.Values, int32(v))
		}

		pbrd.Reels = append(pbrd.Reels, pbl)
	}

	return pbrd
}

// BuildPBPayTables - *sgc7game.PayTables -> map[int32]*Row
func BuildPBPayTables(pt *sgc7game.PayTables) map[int32]*sgc7pb.Row {
	pbpt := make(map[int32]*sgc7pb.Row)

	for k, l := range pt.MapPay {
		pbl := &sgc7pb.Row{}

		for _, v := range l {
			pbl.Values = append(pbl.Values, int32(v))
		}

		pbpt[int32(k)] = pbl
	}

	return pbpt
}

// BuildPBGameScene - *sgc7game.GameScene -> *sgc7pb.GameScene
func BuildPBGameScene(gs *sgc7game.GameScene) *sgc7pb.GameScene {
	pbgs := &sgc7pb.GameScene{}

	for _, l := range gs.Arr {
		pbl := &sgc7pb.Column{}

		for _, v := range l {
			pbl.Values = append(pbl.Values, int32(v))
		}

		pbgs.Values = append(pbgs.Values, pbl)
	}

	if len(gs.Indexes) > 0 {
		for _, v := range gs.Indexes {
			pbgs.Indexes = append(pbgs.Indexes, int32(v))
		}
	}

	if len(gs.ValidRow) > 0 {
		for _, v := range gs.ValidRow {
			pbgs.ValidRow = append(pbgs.ValidRow, int32(v))
		}
	}

	return pbgs
}

// BuildPBGameConfig - *sgc7game.Config -> *sgc7pb.GameConfig
func BuildPBGameConfig(cfg *sgc7game.Config) *sgc7pb.GameConfig {
	pbcfg := &sgc7pb.GameConfig{
		Width:   int32(cfg.Width),
		Height:  int32(cfg.Height),
		Ver:     cfg.Ver,
		CoreVer: cfg.CoreVer,
		BetMuls: cfg.BetMuls,
	}

	if cfg.Lines != nil {
		pbcfg.Lines = BuildPBLineData(cfg.Lines)
	}

	if cfg.Reels != nil {
		pbcfg.Reels = make(map[string]*sgc7pb.ReelsData)

		for k, rd := range cfg.Reels {
			pbcfg.Reels[k] = BuildPBReelsData(rd)
		}
	}

	if cfg.PayTables != nil {
		pbcfg.PayTables = BuildPBPayTables(cfg.PayTables)
	}

	if cfg.DefaultScene != nil {
		pbcfg.DefaultScene = BuildPBGameScene(cfg.DefaultScene)
	}

	if len(cfg.DefaultScene2) > 0 {
		for _, v := range cfg.DefaultScene2 {
			pbcfg.DefaultScene2 = append(pbcfg.DefaultScene2, BuildPBGameScene(v))
		}
	}

	return pbcfg
}

// BuildStake - PlayerState => sgc7game.IPlayerState
func BuildStake(stake *sgc7pb.Stake) *sgc7game.Stake {
	return &sgc7game.Stake{
		CoinBet:  int64(stake.CoinBet),
		CashBet:  int64(stake.CashBet),
		Currency: stake.Currency,
	}
}

// BuildPBRngs - []*RngInfo => []*sgc7pb.RngInfo
func BuildPBRngs(rngs []*sgc7utils.RngInfo) []*sgc7pb.RngInfo {
	pbrngs := []*sgc7pb.RngInfo{}

	for _, v := range rngs {
		pbrngs = append(pbrngs, &sgc7pb.RngInfo{
			Bits:  int32(v.Bits),
			Range: int32(v.Range),
			Value: int32(v.Value),
		})
	}

	return pbrngs
}

// BuildPBGameScenePlayResult - *sgc7game.Result -> *sgc7pb.GameScenePlayResult
func BuildPBGameScenePlayResult(r *sgc7game.Result) *sgc7pb.GameScenePlayResult {
	pr := &sgc7pb.GameScenePlayResult{
		Type:       int32(r.Type),
		Symbol:     int32(r.Symbol),
		Mul:        int32(r.Mul),
		CoinWin:    int32(r.CoinWin),
		CashWin:    int32(r.CashWin),
		OtherMul:   int32(r.OtherMul),
		Wilds:      int32(r.Wilds),
		SymbolNums: int32(r.SymbolNums),
		LineIndex:  int32(r.LineIndex),
	}

	for _, v := range r.Pos {
		pr.Pos = append(pr.Pos, int32(v))
	}

	return pr
}

// MergeReplyPlay - merge ReplyPlay
func MergeReplyPlay(dst *sgc7pb.ReplyPlay, src *sgc7pb.ReplyPlay) {
	if len(src.RandomNumbers) > 0 {
		dst.RandomNumbers = append(dst.RandomNumbers, src.RandomNumbers...)
	}

	if src.PlayerState != nil {
		dst.PlayerState = src.PlayerState
	}

	if src.Finished {
		dst.Finished = true
	}

	if len(src.Results) > 0 {
		dst.Results = append(dst.Results, src.Results...)
	}

	if src.NextCommands != nil {
		dst.NextCommands = src.NextCommands
	}
}
