package mathtoolset

import (
	"log/slog"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type ScriptCore struct {
	Cel *cel.Env
	Ast *cel.Ast
	Prg *cel.Program
}

func (sc *ScriptCore) Compile(code string) error {
	ast, issues := sc.Cel.Compile(code)
	if issues != nil {
		goutils.Error("ScriptCore.Compile:Compile",
			slog.String("code", code),
			slog.Any("issues", issues),
			goutils.Err(ErrInvalidCode))

		return ErrInvalidCode
	}

	sc.Ast = ast

	prg, err := sc.Cel.Program(ast)
	if err != nil {
		goutils.Error("ScriptCore.Compile:Program",
			slog.String("code", code),
			goutils.Err(err))

		return err
	}

	sc.Prg = &prg

	return nil
}

func (sc *ScriptCore) Eval(mgr *GenMathMgr) (ref.Val, error) {
	if sc.Prg != nil {
		out, _, err := (*sc.Prg).Eval(map[string]any{
			"rets":    float64s2list(mgr.Rets),
			"mapRets": mapSF2mapSF(mgr.MapRets),
		})
		if err != nil {
			goutils.Error("ScriptCore.Eval:Eval",
				goutils.Err(err))

			return types.Null(0), err
		}

		return out, nil
	}

	return types.Null(0), nil
}

func array2StrSlice(val ref.Val) []string {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []string{}

		for _, v := range lst0 {
			v1, isok := v.Value().(string)
			if isok {
				lst = append(lst, v1)
			}
		}

		return lst
	}

	return nil
}

func array2IntSlice(val ref.Val) []int {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []int{}

		for _, v := range lst0 {
			v1, isok := v.Value().(int64)
			if isok {
				lst = append(lst, int(v1))
			}
		}

		return lst
	}

	return nil
}

func list2listmapintfloat(val ref.Val) []map[int]float64 {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []map[int]float64{}

		for _, n := range lst0 {
			curmap := make(map[int]float64)

			cm, isok := n.Value().(map[ref.Val]ref.Val)
			if isok {
				for k, v := range cm {
					curmap[int(k.Value().(int64))] = v.Value().(float64)
				}
			}

			lst = append(lst, curmap)
		}

		return lst
	}

	return nil
}

func array2SymbolTypeSlice(val ref.Val) []SymbolType {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []SymbolType{}

		for _, v := range lst0 {
			v1, isok := v.Value().(int64)
			if isok {
				lst = append(lst, SymbolType(v1))
			}
		}

		return lst
	}

	return nil
}

func array2SymbolMapping(val ref.Val) *SymbolMapping {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		sm := NewSymbolMapping()
		var k int

		for i, v := range lst0 {
			v1, isok := v.Value().(int64)
			if isok {
				if i%2 == 1 {
					sm.MapSymbols[SymbolType(k)] = SymbolType(v1)
				} else {
					k = int(v1)
				}
			}
		}

		return sm
	}

	return nil
}

func float64s2list(arr []float64) ref.Val {
	return types.NewDynamicList(types.DefaultTypeAdapter, arr)
}

func mapSF2mapSF(mapSF map[string]float64) ref.Val {
	return types.NewDynamicMap(types.DefaultTypeAdapter, mapSF)
}

func array2SymbolMulti(val0 ref.Val, val1 ref.Val) *sgc7game.ValMapping2 {
	lst0, isok0 := val0.Value().([]ref.Val)
	lst1, isok1 := val1.Value().([]ref.Val)

	if isok0 && isok1 && len(lst0) == len(lst1) {
		sm := sgc7game.NewValMappingEx2()
		for i, v := range lst0 {
			k, isok0 := v.Value().(int64)
			v, isok1 := lst1[i].Value().(float64)
			if isok0 && isok1 {
				sm.MapVals[int(k)] = sgc7game.NewFloatValEx(v)
			}
		}

		return sm
	}

	return nil
}

func array2OverlaySyms(val ref.Val) *sgc7game.ValMapping2 {
	lst, isok := val.Value().([]ref.Val)

	if isok && len(lst)%3 == 0 {
		sm := sgc7game.NewValMappingEx2()
		for i := 0; i < len(lst)/3; i++ {
			x, isok0 := lst[i*3].Value().(int64)
			y, isok1 := lst[i*3+1].Value().(int64)
			s, isok2 := lst[i*3+2].Value().(int64)

			if isok0 && isok1 && isok2 {
				sm.MapVals[pos2int(int(x), int(y))] = sgc7game.NewIntValEx(int(s))
			}
		}

		return sm
	}

	return nil
}

func getFloat64Slice(val ref.Val) []float64 {
	lst, isok := val.Value().([]ref.Val)

	if isok {
		arr := []float64{}

		for _, v := range lst {
			v1, isok1 := v.Value().(float64)
			if isok1 {
				arr = append(arr, v1)
			}
		}

		return arr
	}

	return nil
}

func appendEnvOptions(dst []cel.EnvOption, src []cel.EnvOption) []cel.EnvOption {
	if len(src) > 0 {
		dst = append(dst, src...)
	}

	return dst
}

func newScriptVariables(mgrGenMath *GenMathMgr) []cel.EnvOption {
	return []cel.EnvOption{}
}

func newBasicScriptFuncs(mgrGenMath *GenMathMgr) []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("rets", cel.ListType(cel.DoubleType)),
		cel.Variable("mapRets", cel.MapType(cel.StringType, cel.DoubleType)),
		cel.Function("calcLineRTP",
			cel.Overload("calcLineRTP_string_string_list_list_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 6 {
						goutils.Error("calcLineRTP",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					err := mgrGenMath.LoadPaytables(params[0].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(params[1].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadReelsState",
							goutils.Err(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])
					wilds := array2SymbolTypeSlice(params[3])

					ssws, err := AnalyzeReelsWithLineEx(mgrGenMath.Paytables, mgrGenMath.RSS, syms, wilds, int(params[4].Value().(int64)), int(params[5].Value().(int64)))
					if err != nil {
						goutils.Error("calcLineRTP:AnalyzeReelsWithLineEx",
							goutils.Err(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					ret := float64(ssws.TotalWins) / float64(ssws.TotalBet)
					// mgrGenMath.pushRet(ret)

					return types.Double(ret)
				},
				),
			),
		),
		cel.Function("calcScatterRTP",
			cel.Overload("calcScatterRTP_string_string_list_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 4 {
						goutils.Error("calcLineRTP",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					err := mgrGenMath.LoadPaytables(params[0].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(params[1].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadReelsState",
							goutils.Err(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])

					ssws, err := AnalyzeReelsScatterEx(mgrGenMath.Paytables, mgrGenMath.RSS, syms, int(params[3].Value().(int64)))
					if err != nil {
						goutils.Error("calcLineRTP:AnalyzeReelsScatterEx",
							goutils.Err(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					ret := float64(ssws.TotalWins) / float64(ssws.TotalBet)
					// mgrGenMath.pushRet(ret)

					return types.Double(ret)
				},
				),
			),
		),
		cel.Function("calcWaysRTP",
			cel.Overload("calcWaysRTP_string_string_list_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 8 {
						goutils.Error("calcWaysRTP",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					ptfn := params[0].Value().(string)
					rssfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(ptfn)
					if err != nil {
						goutils.Error("calcWaysRTP:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(rssfn)
					if err != nil {
						goutils.Error("calcWaysRTP:LoadReelsState",
							goutils.Err(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])
					wilds := array2SymbolTypeSlice(params[3])
					sm := array2SymbolMapping(params[4])
					height := int(params[5].Value().(int64))
					bet := int(params[6].Value().(int64))
					mul := int(params[7].Value().(int64))

					ssws, err := AnalyzeReelsWaysEx(mgrGenMath.Paytables, mgrGenMath.RSS, syms, wilds, sm, height, bet, mul)
					if err != nil {
						goutils.Error("calcWaysRTP:AnalyzeReelsWaysEx",
							goutils.Err(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					ret := float64(ssws.TotalWins) / float64(ssws.TotalBet)
					// mgrGenMath.pushRet(ret)

					return types.Double(ret)
				},
				),
			),
		),
		cel.Function("calcWaysRTPWithSymbolMulti",
			cel.Overload("calcWaysRTPWithSymbolMulti_string_string_list_list_list_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.ListType(cel.IntType),
					cel.ListType(cel.IntType), cel.ListType(cel.DoubleType), cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 10 {
						goutils.Error("calcWaysRTPWithSymbolMulti",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					ptfn := params[0].Value().(string)
					rssfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(ptfn)
					if err != nil {
						goutils.Error("calcWaysRTPWithSymbolMulti:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(rssfn)
					if err != nil {
						goutils.Error("calcWaysRTPWithSymbolMulti:LoadReelsState",
							goutils.Err(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])
					wilds := array2SymbolTypeSlice(params[3])
					sm := array2SymbolMapping(params[4])
					symMul := array2SymbolMulti(params[5], params[6])
					height := int(params[7].Value().(int64))
					bet := int(params[8].Value().(int64))
					mul := int(params[9].Value().(int64))

					ssws, err := AnalyzeReelsWaysSymbolMulti(mgrGenMath.Paytables, mgrGenMath.RSS, syms, wilds, sm, symMul, height, bet, mul)
					if err != nil {
						goutils.Error("calcWaysRTPWithSymbolMulti:AnalyzeReelsWaysSymbolMulti",
							goutils.Err(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					ret := float64(ssws.TotalWins) / float64(ssws.TotalBet)
					// mgrGenMath.pushRet(ret)

					return types.Double(ret)
				},
				),
			),
		),
		cel.Function("calcWaysRTP2",
			cel.Overload("calcWaysRTP2_string_string_bool_list_list_list_list_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.BoolType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.ListType(cel.IntType),
					cel.ListType(cel.IntType), cel.ListType(cel.DoubleType), cel.ListType(cel.IntType), cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 12 {
						goutils.Error("calcWaysRTP2",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					paytablefn := params[0].Value().(string)
					reelfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(paytablefn)
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					isStrReel := params[2].Value().(bool)

					rd, err := mgrGenMath.LoadReelsData(paytablefn, reelfn, isStrReel)
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadReelsData",
							goutils.Err(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[3])
					wilds := array2SymbolTypeSlice(params[4])
					sm := array2SymbolMapping(params[5])
					symMul := array2SymbolMulti(params[6], params[7])
					overlaySyms := array2OverlaySyms(params[8])
					height := int(params[9].Value().(int64))
					bet := int(params[10].Value().(int64))
					mul := int(params[11].Value().(int64))

					ssws, err := AnalyzeReelsWaysEx3(mgrGenMath.Paytables, rd, syms, wilds, sm, symMul, overlaySyms, height, bet, mul)
					if err != nil {
						goutils.Error("calcWaysRTP2:AnalyzeReelsWaysEx3",
							goutils.Err(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					ret := float64(ssws.TotalWins) / float64(ssws.TotalBet)
					// mgrGenMath.pushRet(ret)

					return types.Double(ret)
				},
				),
			),
		),
		cel.Function("calcScatterProbability",
			cel.Overload("calcScatterProbability_string_int_int_int",
				[]*cel.Type{cel.StringType, cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 4 {
						goutils.Error("calcScatterProbability",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					rssfn := params[0].Value().(string)

					err := mgrGenMath.LoadReelsState(rssfn)
					if err != nil {
						goutils.Error("calcScatterProbability:LoadReelsState",
							goutils.Err(err))

						return types.Double(0)
					}

					sym := params[1].Value().(int64)
					num := params[2].Value().(int64)
					height := params[3].Value().(int64)

					prob := CalcScatterProbability(mgrGenMath.RSS, SymbolType(sym), int(num), int(height))
					// mgrGenMath.pushRet(prob)

					return types.Double(prob)
				},
				),
			),
		),
		cel.Function("calcProbWithWeights",
			cel.Overload("calcProbWithWeights_string_list",
				[]*cel.Type{cel.StringType, cel.ListType(cel.DoubleType)},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 2 {
						goutils.Error("calcProbWithWeights",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					vwfn := params[0].Value().(string)

					vw, err := sgc7game.LoadValWeights2FromExcel(vwfn, "val", "weight", sgc7game.NewStrVal)
					if err != nil {
						goutils.Error("calcProbWithWeights:LoadValMapping2FromExcel",
							goutils.Err(err))

						return types.Double(0)
					}

					arr := getFloat64Slice(params[1])

					prob := CalcProbWithWeights(vw, arr)
					// mgrGenMath.pushRet(prob)

					return types.Double(prob)
				},
				),
			),
		),
		cel.Function("genReelsMainSymbolsDistance",
			cel.Overload("genReelsMainSymbolsDistance_string_string_string_list_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.StringType, cel.ListType(cel.StringType), cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {

					if len(params) != 5 {
						goutils.Error("genReelsMainSymbolsDistance",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					targetfn := params[0].Value().(string)
					paytablefn := params[1].Value().(string)
					rssfn := params[2].Value().(string)

					err := mgrGenMath.LoadPaytables(paytablefn)
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					rss, err := LoadReelsStats(rssfn)
					if err != nil {
						goutils.Error("genReelsMainSymbolsDistance:LoadReelsStats",
							goutils.Err(err))

						return types.Double(0)
					}

					mainSymbolsWithStr := array2StrSlice(params[3])
					offset := int(params[4].Value().(int64))

					mainSymbols := GetSymbols(mainSymbolsWithStr, mgrGenMath.Paytables)

					reels, err := GenReelsMainSymbolsDistance(rss, mainSymbols, offset, 100)
					if err != nil {
						goutils.Error("genReelsMainSymbolsDistance:GenReelsMainSymbolsDistance",
							goutils.Err(err))

						return types.Double(0)
					}

					err = reels.SaveExcelEx(targetfn, mgrGenMath.Paytables)
					if err != nil {
						goutils.Error("genReelsMainSymbolsDistance:SaveExcelEx",
							goutils.Err(err))

						return types.Double(0)
					}

					return types.Double(1)
				},
				),
			),
		),
		cel.Function("runCode",
			cel.Overload("runCode_string",
				[]*cel.Type{cel.StringType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {

					if len(params) != 1 {
						goutils.Error("runCode",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					codeName := params[0].Value().(string)

					ret, err := mgrGenMath.RunCodeEx(codeName)
					if err != nil {
						goutils.Error("runCode:RunCodeEx",
							goutils.Err(err))

						return types.Double(0)
					}

					return ret
				},
				),
			),
		),
		cel.Function("calcScatterProbabilitWithReels",
			cel.Overload("calcScatterProbabilitWithReels_string_string_bool_string_list_list_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.BoolType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType),
					cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 8 {
						goutils.Error("calcScatterProbabilitWithReels",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					paytablefn := params[0].Value().(string)
					reelfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(paytablefn)
					if err != nil {
						goutils.Error("calcScatterProbabilitWithReels:LoadPaytables",
							goutils.Err(err))

						return types.Double(0)
					}

					isStrReel := params[2].Value().(bool)

					rd, err := mgrGenMath.LoadReelsData(paytablefn, reelfn, isStrReel)
					if err != nil {
						goutils.Error("calcScatterProbabilitWithReels:LoadReelsData",
							goutils.Err(err))

						return types.Double(0)
					}

					sym := mgrGenMath.Paytables.MapSymbols[params[3].Value().(string)]
					sm := array2SymbolMapping(params[4])
					overlaySyms := array2OverlaySyms(params[5])
					num := int(params[6].Value().(int64))
					height := int(params[7].Value().(int64))

					ret := CalcScatterProbabilitWithReels(rd, SymbolType(sym), sm, overlaySyms, num, height)
					if err != nil {
						goutils.Error("CalcScatterProbabilitWithReels:CalcScatterProbabilitWithReels",
							goutils.Err(err))

						return types.Double(0)
					}

					return types.Double(ret)
				},
				),
			),
		),
		cel.Function("calcMulLevelRTP",
			cel.Overload("calcMulLevelRTP_list_list_int_list",
				[]*cel.Type{cel.ListType(cel.DoubleType), cel.ListType(cel.MapType(cel.IntType, cel.DoubleType)), cel.IntType, cel.ListType(cel.IntType)},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 4 {
						goutils.Error("calcMulLevelRTP",
							goutils.Err(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					levelRTPs := getFloat64Slice(params[0])
					levelUpProbs := list2listmapintfloat(params[1])
					spinNum := int(params[2].Value().(int64))
					levelUpAddSpinNum := array2IntSlice(params[3])

					ret := CalcMulLevelRTP2(levelRTPs, levelUpProbs, spinNum, levelUpAddSpinNum)

					return types.Double(ret)
				},
				),
			),
		)}
}

func NewScriptCore(mgrGenMath *GenMathMgr) (*ScriptCore, error) {
	options := []cel.EnvOption{}
	options = appendEnvOptions(options, newScriptVariables(mgrGenMath))
	options = appendEnvOptions(options, newBasicScriptFuncs(mgrGenMath))

	cel, err := cel.NewEnv(options...)
	if err != nil {
		goutils.Error("NewScriptCore:cel.NewEnv",
			goutils.Err(err))

		return nil, err
	}

	return &ScriptCore{
		Cel: cel,
	}, nil
}
