package mathtoolset

import (
	"fmt"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

func GenMath(fn string) error {
	cfg, err := LoadConfig(fn)
	if err != nil {
		goutils.Error("GenMath:LoadConfig",
			zap.Error(err))

		return err
	}

	paytables, err := sgc7game.LoadPaytablesFromExcel(cfg.Paytables)
	if err != nil {
		goutils.Error("GenMath:LoadPaytablesFromExcel",
			zap.Error(err))

		return err
	}

	if cfg.Type == "genReelsState" {
		mgrGenMath := NewGamMathMgr(cfg)

		script, err := NewScriptCore(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:NewScriptCore",
				zap.Error(err))

			return err
		}

		err = script.Compile(cfg.Code)
		if err != nil {
			goutils.Error("GenMath:Compile",
				zap.Error(err))

			return err
		}

		out, err := script.Eval(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:Eval",
				zap.Error(err))

			return err
		}

		if out.Value().(float64) < cfg.TargetRTP {
			goutils.Error("GenMath:Eval",
				zap.Float64("cur-rtp", out.Value().(float64)),
				zap.Error(ErrInvalidTargetRTP))

			return err
		}

		mgrGenMath.Save()
	} else if cfg.Type == "calcRTPWithReelsState" {
		mgrGenMath := NewGamMathMgr(cfg)

		script, err := NewScriptCore(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:NewScriptCore",
				zap.Error(err))

			return err
		}

		err = script.Compile(cfg.Code)
		if err != nil {
			goutils.Error("GenMath:Compile",
				zap.Error(err))

			return err
		}

		out, err := script.Eval(mgrGenMath)
		if err != nil {
			goutils.Error("GenMath:Eval",
				zap.Error(err))

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
				zap.Error(err))

			return err
		}

		mainSymbols := GetSymbols(cfg.GenReelsConfig.MainSymbols, paytables)

		reels, err := GenReelsMainSymbolsDistance(rss, mainSymbols, cfg.GenReelsConfig.Offset, 100)
		if err != nil {
			goutils.Error("GenMath:GenReelsMainSymbolsDistance",
				zap.Error(err))

			return err
		}

		err = reels.SaveExcel(cfg.GenReelsConfig.ReelsFilename)
		if err != nil {
			goutils.Error("GenMath:SaveExcel",
				zap.Error(err))

			return err
		}
	}

	return nil
}
