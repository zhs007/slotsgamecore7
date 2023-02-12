package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type acwData[T int | float32 | float64] struct {
	group0 []int
	val0   T
	group1 []int
	val1   T
}

func (acwd *acwData[T]) calcVal0(vm *sgc7game.ValMapping[int, T], vw *sgc7game.ValWeights) {
	var val T

	for _, v := range acwd.group0 {
		val += T(float64(vm.MapVals[v]) * float64(vw.GetWeight(v)) / float64(vw.MaxWeight))
	}

	acwd.val0 = val
}

func (acwd *acwData[T]) calcGroup1AndVal1(vm *sgc7game.ValMapping[int, T], vw *sgc7game.ValWeights) {
	acwd.group1 = nil

	for k := range vm.MapVals {
		acwd.group1 = append(acwd.group1, k)
	}

	var val T

	for _, v := range acwd.group1 {
		val += T(float64(vm.MapVals[v]) * float64(vw.GetWeight(v)) / float64(vw.MaxWeight))
	}

	acwd.val1 = val
}

type FuncRunnerWithValWeights[T int | float32 | float64] func(nvw *sgc7game.ValWeights, isfastmode bool) T

type funcFEAWL func([]int)

func forEachArrWithLength(dest []int, src []int, length int, onforeach funcFEAWL) {
	if length == 1 {
		for i := 0; i < len(src); i++ {
			dest = append(dest, src[i])

			onforeach(dest)

			dest = dest[0 : len(dest)-1]
		}

		return
	}

	for i := 0; i < len(src); i++ {
		dest = append(dest, src[i])

		forEachArrWithLength(dest, append(src[0:i], src[i+1:]...), length-1, onforeach)

		dest = dest[0 : len(dest)-1]
	}
}

type funcFEACWD[T int | float32 | float64] func(*acwData[T])

func forEachACWData[T int | float32 | float64](vm *sgc7game.ValMapping[int, T], vw *sgc7game.ValWeights, foreach funcFEACWD[T]) {
	arr := vm.Vals()

	for i := 1; i <= len(vm.MapVals)/2; i++ {
		forEachArrWithLength(nil, arr, i, func(group0 []int) {
			acwd := &acwData[T]{
				group0: group0,
			}

			acwd.calcVal0(vm, vw)
			acwd.calcGroup1AndVal1(vm, vw)

			foreach(acwd)
		})
	}
}

func AnalyzeWeights[T int | float32 | float64](vw *sgc7game.ValWeights,
	runner FuncRunnerWithValWeights[T]) (*sgc7game.ValMapping[int, T], error) {

	mappingVals := sgc7game.NewValMappingEx[int, T]()
	for _, v := range vw.Vals {
		nvw := vw.Clone()

		nvw.ClearExcludeVal(v)

		mappingVals.MapVals[v] = runner(nvw, true)
	}

	return mappingVals, nil
}

func AutoChgWeights[T int | float32 | float64](vw *sgc7game.ValWeights, target T,
	runner FuncRunnerWithValWeights[T]) (*sgc7game.ValWeights, error) {

	if len(vw.Vals) <= 1 {
		goutils.Error("AutoChgWeights",
			zap.Error(ErrValidParamInAutoChgWeights))

		return nil, ErrValidParamInAutoChgWeights
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
			zap.Error(ErrValidParamInAutoChgWeights))

		return nil, ErrValidParamInAutoChgWeights
	}

	forEachACWData(mappingVals, vw, func(acwd *acwData[T]) {
		if acwd.val0 > target && acwd.val1 > target {
			return
		} else if acwd.val0 < target && acwd.val1 < target {
			return
		}

	})

	return nil, nil
}
