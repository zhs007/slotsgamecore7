package dtserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// BuildLineData - *sgc7game.LineData -> *sgc7pb.LinesData
func BuildLineData(ld *sgc7game.LineData) *sgc7pb.LinesData {
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

// BuildReelsData - *sgc7game.ReelsData -> *sgc7pb.ReelsData
func BuildReelsData(rd *sgc7game.ReelsData) *sgc7pb.ReelsData {
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

// BuildPayTables - *sgc7game.PayTables -> map[int32]*Row
func BuildPayTables(pt *sgc7game.PayTables) map[int32]*sgc7pb.Row {
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

// BuildGameScene - *sgc7game.GameScene -> *sgc7pb.GameScene
func BuildGameScene(gs *sgc7game.GameScene) *sgc7pb.GameScene {
	pbgs := &sgc7pb.GameScene{}

	for _, l := range gs.Arr {
		pbl := &sgc7pb.Column{}

		for _, v := range l {
			pbl.Values = append(pbl.Values, int32(v))
		}

		pbgs.Values = append(pbgs.Values, pbl)
	}

	return pbgs
}

// BuildGameConfig - *sgc7game.Config -> *sgc7pb.GameConfig
func BuildGameConfig(cfg *sgc7game.Config) *sgc7pb.GameConfig {
	pbcfg := &sgc7pb.GameConfig{
		Width:   int32(cfg.Width),
		Height:  int32(cfg.Height),
		Ver:     cfg.Ver,
		CoreVer: cfg.CoreVer,
	}

	if cfg.Lines != nil {
		pbcfg.Lines = BuildLineData(cfg.Lines)
	}

	if cfg.Reels != nil {
		pbcfg.Reels = make(map[string]*sgc7pb.ReelsData)

		for k, rd := range cfg.Reels {
			pbcfg.Reels[k] = BuildReelsData(rd)
		}
	}

	if cfg.PayTables != nil {
		pbcfg.PayTables = BuildPayTables(cfg.PayTables)
	}

	if cfg.DefaultScene != nil {
		pbcfg.DefaultScene = BuildGameScene(cfg.DefaultScene)
	}

	return pbcfg
}
