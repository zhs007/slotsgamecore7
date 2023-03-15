package mathtoolset

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

func GenMath(fn string) error {
	cfg, err := LoadConfig(fn)
	if err != nil {
		goutils.Error("GenMath:LoadConfig",
			zap.Error(err))

		return err
	}

	if cfg.Type == "genReelsState" {
		mgrGenMath := NewGamMathMgr()

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

		out, err := script.Eval()
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
	}

	return nil
}
