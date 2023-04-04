package mathtoolset

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
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
			cel.Overload("calcWaysRTP_string_string_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 7 {
						goutils.Error("calcWaysRTP",
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					err := mgrGenMath.LoadPaytables(params[0].Value().(string))
					if err != nil {
						goutils.Error("calcWaysRTP:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					err = mgrGenMath.LoadReelsState(params[1].Value().(string))
					if err != nil {
						goutils.Error("calcWaysRTP:LoadReelsState",
							zap.Error(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])
					wilds := array2SymbolTypeSlice(params[3])

					ssws, err := AnalyzeReelsWaysEx(mgrGenMath.Paytables, mgrGenMath.RSS, syms, wilds, int(params[4].Value().(int64)), int(params[5].Value().(int64)), int(params[6].Value().(int64)))
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
		cel.Function("calcWaysRTP2",
			cel.Overload("calcWaysRTP2_string_string_list_list_int_int_int",
				[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.IntType), cel.IntType, cel.IntType, cel.IntType},
				cel.DoubleType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 7 {
						goutils.Error("calcWaysRTP2",
							zap.Error(ErrInvalidFunctionParams))

						return types.Double(0)
					}

					err := mgrGenMath.LoadPaytables(params[0].Value().(string))
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadPaytables",
							zap.Error(err))

						return types.Double(0)
					}

					rd, err := mgrGenMath.LoadReelsData2(params[0].Value().(string), params[1].Value().(string))
					if err != nil {
						goutils.Error("calcWaysRTP2:LoadReelsData2",
							zap.Error(err))

						return types.Double(0)
					}

					syms := array2SymbolTypeSlice(params[2])
					wilds := array2SymbolTypeSlice(params[3])
					height := int(params[4].Value().(int64))
					bet := int(params[5].Value().(int64))
					mul := int(params[6].Value().(int64))

					wrss := BuildWaysReelsStatsEx(rd, height, syms, wilds)

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
