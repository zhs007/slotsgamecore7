package mathtoolset

import (
	"fmt"
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func GenMath(fn string) error {
	cfg, err := LoadConfig(fn)
	if err != nil {
		goutils.Error("GenMath:LoadConfig",
			goutils.Err(err))

		return err
	}

	paytables, err := sgc7game.LoadPaytablesFromExcel(cfg.Paytables)
	if err != nil {
		goutils.Error("GenMath:LoadPaytablesFromExcel",
			goutils.Err(err))

		return err
	}

	if cfg.Type == "genReelsState" {
		mgrGenMath := NewGamMathMgr(cfg)

		script, err := NewScriptCore(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:NewScriptCore",
				goutils.Err(err))

			return err
		}

		err = script.Compile(cfg.Code)
		if err != nil {
			goutils.Error("GenMath:Compile",
				goutils.Err(err))

			return err
		}

		out, err := script.Eval(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:Eval",
				goutils.Err(err))

			return err
		}

		if out.Value().(float64) < cfg.TargetRTP {
			goutils.Error("GenMath:Eval",
				slog.Float64("cur-rtp", out.Value().(float64)),
				goutils.Err(ErrInvalidTargetRTP))

			return err
		}

		mgrGenMath.Save()
	} else if cfg.Type == "calcRTPWithReelsState" {
		mgrGenMath := NewGamMathMgr(cfg)

		script, err := NewScriptCore(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:NewScriptCore",
				goutils.Err(err))

			return err
		}

		err = script.Compile(cfg.Code)
		if err != nil {
			goutils.Error("GenMath:Compile",
				goutils.Err(err))

			return err
		}

		out, err := script.Eval(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:Eval",
				goutils.Err(err))

			return err
		}

		mgrGenMath.Save()

		fmt.Printf("The RTP is %v\n", out.Value().(float64))
	} else if cfg.Type == "runCodes" {
		mgrGenMath := NewGamMathMgr(cfg)

		for _, v := range cfg.Codes {
			if !v.DisableAutoRun {
				mgrGenMath.RunCodeEx(v.Name)
			}
		}

		mgrGenMath.Save()
	} else if cfg.Type == "genReels" {
		rss, err := LoadReelsStats(cfg.GenReelsConfig.ReelsStatsFilename)
		if err != nil {
			goutils.Error("GenMath:LoadReelsStats",
				goutils.Err(err))

			return err
		}

		mainSymbols := GetSymbols(cfg.GenReelsConfig.MainSymbols, paytables)

		reels, err := GenReelsMainSymbolsDistance(rss, mainSymbols, cfg.GenReelsConfig.Offset, 100)
		if err != nil {
			goutils.Error("GenMath:GenReelsMainSymbolsDistance",
				goutils.Err(err))

			return err
		}

		err = reels.SaveExcel(cfg.GenReelsConfig.ReelsFilename)
		if err != nil {
			goutils.Error("GenMath:SaveExcel",
				goutils.Err(err))

			return err
		}
	}

	return nil
}
