package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type FuncRunner[T int | float32 | float64] func(nvw *sgc7game.ValWeights, isfastmode bool) T

func AutoChgWeights[T int | float32 | float64](vw *sgc7game.ValWeights, target T, runner FuncRunner[T]) (*sgc7game.ValWeights, error) {
	if len(vw.Vals) <= 1 {
		goutils.Error("AutoChgWeights",
			zap.Error(ErrNoValidParamInAutoChgWeights))

		return nil, ErrNoValidParamInAutoChgWeights
	}

	curval := runner(vw, false)
	if curval == target {
		return vw, nil
	}

	hasbigger := false
	hassmaller := false
	mappingVals := sgc7game.NewValMappingEx[int, T]()
	for _, v := range vw.Vals {
		nvw := vw.Clone()

		nvw.ClearExcludeVal(v)

		mappingVals.MapVals[v] = runner(nvw, true)

		goutils.Info("AutoChgWeights:runner",
			goutils.JSON("ValWeights", nvw),
			zap.Any("return", mappingVals.MapVals[v]))

		if mappingVals.MapVals[v] > target {
			hasbigger = true
		}

		if mappingVals.MapVals[v] < target {
			hassmaller = true
		}
	}

	if !hasbigger || !hassmaller {
		goutils.Error("AutoChgWeights",
			zap.Error(ErrNoValidParamInAutoChgWeights))

		return nil, ErrNoValidParamInAutoChgWeights
	}

	return nil, nil
}
