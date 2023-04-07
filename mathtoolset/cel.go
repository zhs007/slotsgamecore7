package mathtoolset

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
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
			zap.String("code", code),
			goutils.JSON("issues", issues),
			zap.Error(ErrInvalidCode))

		return ErrInvalidCode
	}

	sc.Ast = ast

	prg, err := sc.Cel.Program(ast)
	if err != nil {
		goutils.Error("ScriptCore.Compile:Program",
			zap.String("code", code),
			zap.Error(err))

		return err
	}

	sc.Prg = &prg

	return nil
}

func (sc *ScriptCore) Eval() (ref.Val, error) {
	if sc.Prg != nil {
		out, _, err := (*sc.Prg).Eval(map[string]any{})
		if err != nil {
			goutils.Error("ScriptCore.Eval:Eval",
				zap.Error(err))

			return types.Null(0), err
		}

		return out, nil
	}

	return types.Null(0), nil
}

func array2IntSlice(val ref.Val) []int {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []int{}

		for _, v := range lst0 {
			v1, isok := v.Value().(int)
			if isok {
				lst = append(lst, v1)
			}
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
		cel.Function("calcLineRTP",
			cel.Overload("calcLineRTP_string_string_list_list_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 6 {
						goutils.Error("calcLineRTP",
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					err := mgrGenMath.LoadPaytables(params[0].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(params[1].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadReelsState",
							zap.Error(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])
					wilds := array2SymbolTypeSlice(params[3])

					ssws, err := AnalyzeReelsWithLineEx(mgrGenMath.Paytables, mgrGenMath.RSS, syms, wilds, int(params[4].Value().(int64)), int(params[5].Value().(int64)))
					if err != nil {
						goutils.Error("calcLineRTP:AnalyzeReelsWithLineEx",
							zap.Error(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					return types.Double(float64(ssws.TotalWins) / float64(ssws.TotalBet))
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
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					err := mgrGenMath.LoadPaytables(params[0].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(params[1].Value().(string))
					if err != nil {
						goutils.Error("calcLineRTP:LoadReelsState",
							zap.Error(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])

					ssws, err := AnalyzeReelsScatterEx(mgrGenMath.Paytables, mgrGenMath.RSS, syms, int(params[3].Value().(int64)))
					if err != nil {
						goutils.Error("calcLineRTP:AnalyzeReelsScatterEx",
							zap.Error(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					return types.Double(float64(ssws.TotalWins) / float64(ssws.TotalBet))
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
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					ptfn := params[0].Value().(string)
					rssfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(ptfn)
					if err != nil {
						goutils.Error("calcWaysRTP:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(rssfn)
					if err != nil {
						goutils.Error("calcWaysRTP:LoadReelsState",
							zap.Error(err))

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
							zap.Error(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					return types.Double(float64(ssws.TotalWins) / float64(ssws.TotalBet))
				},
				),
			),
		),
		cel.Function("calcWaysRTPWithSymbolMulyi",
			cel.Overload("calcWaysRTPWithSymbolMulyi_string_string_list_list_list_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.ListType(cel.IntType),
					cel.ListType(cel.IntType), cel.ListType(cel.DoubleType), cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 10 {
						goutils.Error("calcWaysRTPWithSymbolMulyi",
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					ptfn := params[0].Value().(string)
					rssfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(ptfn)
					if err != nil {
						goutils.Error("calcWaysRTPWithSymbolMulyi:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(rssfn)
					if err != nil {
						goutils.Error("calcWaysRTPWithSymbolMulyi:LoadReelsState",
							zap.Error(err))

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
						goutils.Error("calcWaysRTPWithSymbolMulyi:AnalyzeReelsWaysSymbolMulti",
							zap.Error(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					return types.Double(float64(ssws.TotalWins) / float64(ssws.TotalBet))
				},
				),
			),
		),
		cel.Function("calcWaysRTP2",
			cel.Overload("calcWaysRTP2_string_string_bool_list_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.BoolType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.ListType(cel.IntType),
					cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 9 {
						goutils.Error("calcWaysRTP2",
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					paytablefn := params[0].Value().(string)
					reelfn := params[1].Value().(string)

					err := mgrGenMath.LoadPaytables(paytablefn)
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					isStrReel := params[2].Value().(bool)

					rd, err := mgrGenMath.LoadReelsData(paytablefn, reelfn, isStrReel)
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadReelsData",
							zap.Error(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[3])
					wilds := array2SymbolTypeSlice(params[4])
					sm := array2SymbolMapping(params[5])
					height := int(params[6].Value().(int64))
					bet := int(params[7].Value().(int64))
					mul := int(params[8].Value().(int64))

					wrss := BuildWaysReelsStatsEx(rd, height, syms, wilds, sm)

					ssws, err := AnalyzeReelsWaysEx2(mgrGenMath.Paytables, wrss, syms, height, bet, mul)
					if err != nil {
						goutils.Error("calcWaysRTP2:AnalyzeReelsWaysEx2",
							zap.Error(err))

						return types.Double(0)
					}

					mgrGenMath.RetStats = append(mgrGenMath.RetStats, ssws)

					return types.Double(float64(ssws.TotalWins) / float64(ssws.TotalBet))
				},
				),
			),
		),
	}
}

func NewScriptCore(mgrGenMath *GenMathMgr) (*ScriptCore, error) {
	options := []cel.EnvOption{}
	options = appendEnvOptions(options, newScriptVariables(mgrGenMath))
	options = appendEnvOptions(options, newBasicScriptFuncs(mgrGenMath))

	cel, err := cel.NewEnv(options...)
	if err != nil {
		goutils.Error("NewScriptCore:cel.NewEnv",
			zap.Error(err))

		return nil, err
	}

	return &ScriptCore{
		Cel: cel,
	}, nil
}
