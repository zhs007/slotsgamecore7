package mathtoolset

import (
	"fmt"
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type FuncCmpTarget[T float32 | float64] func(v0 T, v1 T) int

type acwData[T float32 | float64] struct {
	group0  []int
	val0    T
	group1  []int
	val1    T
	weight0 float64
}

func (acwd *acwData[T]) outputString() string {
	return fmt.Sprintf("group0 - %v val0 - %v group1 - %v val1 - %v weight0 - %v",
		acwd.group0, acwd.val0, acwd.group1, acwd.val1, acwd.weight0)
}

func (acwd *acwData[T]) calcVal0(vm *sgc7game.FloatValMapping[int, T], vw *sgc7game.ValWeights) {
	var val T
	maxweight := 0

	for _, v := range acwd.group0 {
		maxweight += vw.GetWeight(v)
	}

	for _, v := range acwd.group0 {
		val += T(float64(vm.MapVals[v]) * float64(vw.GetWeight(v)) / float64(maxweight))
	}

	acwd.val0 = val
}

func (acwd *acwData[T]) calcGroup1AndVal1(vm *sgc7game.FloatValMapping[int, T], vw *sgc7game.ValWeights) {
	acwd.group1 = nil

	for k := range vm.MapVals {
		if goutils.IndexOfIntSlice(acwd.group0, k, 0) == -1 {
			acwd.group1 = append(acwd.group1, k)
		}
	}

	var val T
	maxweight := 0

	for _, v := range acwd.group1 {
		maxweight += vw.GetWeight(v)
	}

	for _, v := range acwd.group1 {
		val += T(float64(vm.MapVals[v]) * float64(vw.GetWeight(v)) / float64(maxweight))
	}

	acwd.val1 = val
}

func (acwd *acwData[T]) calcTarget(target T) bool {
	if acwd.val0 > target && acwd.val1 > target {
		return false
	} else if acwd.val0 < target && acwd.val1 < target {
		return false
	}

	acwd.weight0 = float64(target-acwd.val1) / float64(acwd.val0-acwd.val1)

	return true
}

func (acwd *acwData[T]) calcOff(vw *sgc7game.ValWeights) float64 {
	var off float64

	maxweight := 0

	for _, v := range acwd.group0 {
		maxweight += vw.GetWeight(v)
	}

	for _, v := range acwd.group0 {
		co := (float64(vw.GetWeight(v))/float64(vw.MaxWeight) -
			float64(vw.GetWeight(v))/float64(maxweight)*acwd.weight0) /
			float64(vw.GetWeight(v)) / float64(vw.MaxWeight)
		if co < 0 {
			off += -co
		} else {
			off += co
		}
	}

	maxweight = 0

	for _, v := range acwd.group1 {
		maxweight += vw.GetWeight(v)
	}

	for _, v := range acwd.group1 {
		co := (float64(vw.GetWeight(v))/float64(vw.MaxWeight) -
			float64(vw.GetWeight(v))/float64(maxweight)*(1-acwd.weight0)) /
			float64(vw.GetWeight(v)) / float64(vw.MaxWeight)
		if co < 0 {
			off += -co
		} else {
			off += co
		}
	}

	return off
}

func (acwd *acwData[T]) calcNewValWeights(vw *sgc7game.ValWeights, precision int) *sgc7game.ValWeights {
	nvw := sgc7game.NewValWeightsEx()

	maxweight := 0

	for _, v := range acwd.group0 {
		maxweight += vw.GetWeight(v)
	}

	for _, v := range acwd.group0 {
		nvw.Add(v, int(float64(vw.GetWeight(v))/float64(maxweight)*acwd.weight0*float64(precision)))
	}

	maxweight = 0

	for _, v := range acwd.group1 {
		maxweight += vw.GetWeight(v)
	}

	for _, v := range acwd.group1 {
		nvw.Add(v, int(float64(vw.GetWeight(v))/float64(maxweight)*(1-acwd.weight0)*float64(precision)))
	}

	return nvw
}

type FuncRunnerWithValWeights[T float32 | float64] func(nvw *sgc7game.ValWeights, isfastmode bool) T

type funcFEAWL func([]int)

func forEachArrWithLength(dest []int, src []int, length int, onforeach funcFEAWL) {
	if length > len(src) {
		return
	}

	if length == len(src) {
		dest = append(dest, src...)

		onforeach(dest)

		// dest = dest[0 : len(dest)-len(src)]

		return
	}

	if length == 1 {
		for i := 0; i < len(src); i++ {
			dest = append(dest, src[i])

			onforeach(dest)

			dest = dest[0 : len(dest)-1]
		}

		return
	}

	nsrc := make([]int, 0, len(src))

	for i := 0; i <= len(src)-length; i++ {
		dest = append(dest, src[i])

		nsrc = append(nsrc, src[i+1:]...)
		forEachArrWithLength(dest, nsrc, length-1, onforeach)

		dest = dest[0 : len(dest)-1]
		nsrc = nsrc[:0]
	}
}

type funcFEACWD[T float32 | float64] func(*acwData[T])

func forEachACWData[T float32 | float64](vm *sgc7game.FloatValMapping[int, T], vw *sgc7game.ValWeights, foreach funcFEACWD[T]) {
	arr := vm.Keys()

	// num := len(vm.MapVals) / 2
	// if len(vm.MapVals)%2 > 1 {
	// 	num++
	// }

	for i := 1; i <= len(vm.MapVals)/2; i++ {
		forEachArrWithLength(nil, arr, i, func(group0 []int) {
			acwd := &acwData[T]{
				group0: make([]int, len(group0)),
			}

			copy(acwd.group0, group0)

			acwd.calcVal0(vm, vw)
			acwd.calcGroup1AndVal1(vm, vw)

			foreach(acwd)
		})
	}
}

func AnalyzeWeights[T float32 | float64](vw *sgc7game.ValWeights,
	runner FuncRunnerWithValWeights[T]) (*sgc7game.FloatValMapping[int, T], error) {

	mappingVals := sgc7game.NewFloatValMappingEx[int, T]()
	for _, v := range vw.Vals {
		nvw := vw.Clone()

		nvw.ClearExcludeVal(v)

		mappingVals.MapVals[v] = runner(nvw, true)
	}

	return mappingVals, nil
}

func AutoChgWeights[T float32 | float64](vw *sgc7game.ValWeights, target T,
	runner FuncRunnerWithValWeights[T], precision int, cmpTarget FuncCmpTarget[T]) (*sgc7game.ValWeights, error) {

	if len(vw.Vals) <= 1 {
		goutils.Error("AutoChgWeights",
			goutils.Err(ErrValidParamInAutoChgWeights))

		return nil, ErrValidParamInAutoChgWeights
	}

	curval := runner(vw, false)
	if cmpTarget(curval, target) == 0 {
		return vw, nil
	}

	hasbigger := false
	hassmaller := false
	mappingVals := sgc7game.NewFloatValMappingEx[int, T]()
	for _, v := range vw.Vals {
		nvw := vw.Clone()

		nvw.ClearExcludeVal(v)

		mappingVals.MapVals[v] = runner(nvw, true)

		goutils.Info("AutoChgWeights:runner",
			slog.Any("ValWeights", nvw),
			slog.Any("return", mappingVals.MapVals[v]))

		if cmpTarget(mappingVals.MapVals[v], target) > 0 {
			hasbigger = true
		}

		if cmpTarget(mappingVals.MapVals[v], target) < 0 {
			hassmaller = true
		}
	}

	if !hasbigger || !hassmaller {
		goutils.Error("AutoChgWeights",
			goutils.Err(ErrValidParamInAutoChgWeights))

		return nil, ErrValidParamInAutoChgWeights
	}

	var curacwd *acwData[T]
	var curoff float64

	forEachACWData(mappingVals, vw, func(acwd *acwData[T]) {
		if acwd.calcTarget(target) {
			if curacwd == nil {
				curacwd = acwd

				curoff = curacwd.calcOff(vw)

				goutils.Info("AutoChgWeights:result",
					slog.Any("ret", acwd.outputString()),
					slog.Any("off", curoff))
			} else {
				off := acwd.calcOff(vw)

				goutils.Info("AutoChgWeights:result",
					slog.Any("ret", acwd.outputString()),
					slog.Any("off", off))

				if off < curoff {
					curacwd = acwd
					curoff = off
				}
			}
		} else {
			goutils.Info("AutoChgWeights:result",
				slog.Any("ret", acwd.outputString()))
		}
	})

	if curacwd != nil {
		nvw := curacwd.calcNewValWeights(vw, precision)

		nvw.SortBy(vw)

		return nvw, nil
	}

	goutils.Error("AutoChgWeights",
		goutils.Err(ErrNoResultInAutoChgWeights))

	return nil, ErrNoResultInAutoChgWeights
}

func AutoChgWeightsEx[T float32 | float64](vm *sgc7game.FloatValMapping[int, T],
	vw *sgc7game.ValWeights, target T,
	runner FuncRunnerWithValWeights[T], precision int, cmpTarget FuncCmpTarget[T]) (*sgc7game.ValWeights, error) {

	if len(vw.Vals) <= 1 {
		goutils.Error("AutoChgWeightsEx",
			goutils.Err(ErrValidParamInAutoChgWeights))

		return nil, ErrValidParamInAutoChgWeights
	}

	curval := runner(vw, false)
	if cmpTarget(curval, target) == 0 {
		return vw, nil
	}

	hasbigger := false
	hassmaller := false
	mappingVals := vm
	for _, v := range vw.Vals {
		if cmpTarget(mappingVals.MapVals[v], target) > 0 {
			hasbigger = true
		}

		if cmpTarget(mappingVals.MapVals[v], target) < 0 {
			hassmaller = true
		}
	}

	if !hasbigger || !hassmaller {
		goutils.Error("AutoChgWeightsEx",
			goutils.Err(ErrValidParamInAutoChgWeights))

		return nil, ErrValidParamInAutoChgWeights
	}

	var curacwd *acwData[T]
	var curoff float64

	forEachACWData(mappingVals, vw, func(acwd *acwData[T]) {
		if acwd.calcTarget(target) {
			if curacwd == nil {
				curacwd = acwd

				curoff = curacwd.calcOff(vw)

				goutils.Info("AutoChgWeightsEx:result",
					slog.Any("ret", acwd.outputString()),
					slog.Any("off", curoff))
			} else {
				off := acwd.calcOff(vw)

				goutils.Info("AutoChgWeightsEx:result",
					slog.Any("ret", acwd.outputString()),
					slog.Any("off", off))

				if off < curoff {
					curacwd = acwd
					curoff = off
				}
			}
		} else {
			goutils.Info("AutoChgWeightsEx:result",
				slog.Any("ret", acwd.outputString()))
		}
	})

	if curacwd != nil {
		nvw := curacwd.calcNewValWeights(vw, precision)

		nvw.SortBy(vw)

		return nvw, nil
	}

	goutils.Error("AutoChgWeightsEx",
		goutils.Err(ErrNoResultInAutoChgWeights))

	return nil, ErrNoResultInAutoChgWeights
}
