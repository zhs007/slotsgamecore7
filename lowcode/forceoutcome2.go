package lowcode

import (
	"log/slog"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// ForceOutcome2 - 通过 results 来分析数据，所以只能做为临时变量用
type ForceOutcome2 struct {
	cel     *cel.Env
	results []*sgc7game.PlayResult
	program cel.Program
}

func (fo2 *ForceOutcome2) SetScript(code string) error {
	ast, issues := fo2.cel.Compile(code)
	if issues != nil {
		goutils.Error("ForceOutcome2.SetScript:Compile",
			slog.String("code", code),
			slog.Any("issues", issues),
			goutils.Err(ErrInvalidForceOutcome2Code))

		return ErrInvalidForceOutcome2Code
	}

	prg, err := fo2.cel.Program(ast)
	if err != nil {
		goutils.Error("ForceOutcome2.SetScript:Program",
			slog.String("code", code),
			goutils.Err(err))

		return err
	}

	fo2.program = prg

	return nil
}

func (fo2 *ForceOutcome2) IsValid(results []*sgc7game.PlayResult) bool {
	fo2.results = results

	totalwins := 0
	for _, v := range results {
		totalwins += v.CoinWin
	}

	// 必须返回一个 bool
	out, _, err := fo2.program.Eval(map[string]any{
		"totalWins": totalwins,
	})
	if err != nil {
		goutils.Error("ForceOutcome2.IsValid:Eval",
			goutils.Err(err))

		return false
	}

	ret, isok := out.Value().(bool)
	if !isok {
		goutils.Error("ForceOutcome2.IsValid:ret",
			goutils.Err(ErrInvalidForceOutcome2ReturnVal))

		return false
	}

	return ret
}

func (fo2 *ForceOutcome2) CalcVal(results []*sgc7game.PlayResult) int {
	fo2.results = results

	totalwins := 0
	for _, v := range results {
		totalwins += v.CoinWin
	}

	// 必须返回一个 int
	out, _, err := fo2.program.Eval(map[string]any{
		"totalWins": totalwins,
	})
	if err != nil {
		goutils.Error("ForceOutcome2.CalcVal:Eval",
			goutils.Err(err))

		return -1
	}

	ret, isok := out.Value().(int64)
	if !isok {
		goutils.Error("ForceOutcome2.CalcVal:ret",
			goutils.Err(ErrInvalidForceOutcome2ReturnVal))

		return -1
	}

	return int(ret)
}

func (fo2 *ForceOutcome2) hasComponent(component string) bool {
	for _, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for _, v := range gp.HistoryComponents {
				if isComponent(v, component) {
					return true
				}
			}
		}
	}

	return false
}

func (fo2 *ForceOutcome2) getComponentValAt(hasComponent string, component string, val string) int {
	for i, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k := range gp.MapComponentData {
				if isComponent(k, hasComponent) {
					v, hasv := fo2.getComponentValEx(i, component, val)
					if hasv {
						return v
					}
				}
			}
		}
	}

	return 0
}

func (fo2 *ForceOutcome2) getComponentValEx(iStep int, component string, val string) (int, bool) {
	if iStep >= 0 && iStep < len(fo2.results) {
		ret := fo2.results[iStep]
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k, v := range gp.MapComponentData {
				if isComponent(k, component) {
					curval, isok2 := v.GetVal(val)
					if isok2 {
						return curval, true
					}
				}
			}
		}
	}

	return 0, false
}

func (fo2 *ForceOutcome2) getComponentData(iStep int, component string) IComponentData {
	if iStep >= 0 && iStep < len(fo2.results) {
		ret := fo2.results[iStep]
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k, v := range gp.MapComponentData {
				if isComponent(k, component) {
					return v
				}
			}
		}
	}

	return nil
}

func (fo2 *ForceOutcome2) getComponentValNext(hasComponent string, component string, val string) int {
	for i, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k := range gp.MapComponentData {
				if isComponent(k, hasComponent) {
					v, hasv := fo2.getComponentValEx(i+1, component, val)
					if hasv {
						return v
					}
				}
			}
		}
	}

	return 0
}

func (fo2 *ForceOutcome2) getComponentVal(component string, val string) int {
	for _, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k, v := range gp.MapComponentData {
				if isComponent(k, component) {
					curval, isok2 := v.GetVal(val)
					if isok2 {
						return curval
					}
				}
			}
		}
	}

	return 0
}

func (fo2 *ForceOutcome2) getMaxComponentVal(component string, val string) int {
	hasval := false
	maxval := 0
	for _, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k, v := range gp.MapComponentData {
				if isComponent(k, component) {
					curval, isok2 := v.GetVal(val)
					if isok2 {
						if !hasval {
							maxval = curval
							hasval = true
						} else if maxval < curval {
							maxval = curval
						}
					}
				}
			}
		}
	}

	if hasval {
		return maxval
	}

	return 0
}

func (fo2 *ForceOutcome2) getLatestSymbolVal() *sgc7game.GameScene {
	if len(fo2.results) == 0 {
		return nil
	}

	ret := fo2.results[len(fo2.results)-1]
	if len(ret.OtherScenes) <= 0 {
		return nil
	}

	return ret.OtherScenes[len(ret.OtherScenes)-1]
}

func (fo2 *ForceOutcome2) getSymbolVal(x, y int, defval int) int {
	os := fo2.getLatestSymbolVal()
	if os == nil {
		return defval
	}

	return os.Arr[x][y]
}

func (fo2 *ForceOutcome2) countSymbolVal(op string, target int) int {
	os := fo2.getLatestSymbolVal()
	if os == nil {
		return 0
	}

	num := 0
	for _, arr := range os.Arr {
		for _, v := range arr {
			if CmpVal(v, op, target) {
				num++
			}
		}
	}

	return num
}

func (fo2 *ForceOutcome2) hasSamePosNext(src string, target string) bool {
	for i, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k, v := range gp.MapComponentData {
				if isComponent(k, src) {
					for ti := i + 1; ti < len(fo2.results); ti++ {
						tcd := fo2.getComponentData(ti, target)
						if tcd != nil {
							if HasSamePos(v.GetPos(), tcd.GetPos()) {
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

func (fo2 *ForceOutcome2) newScriptVariables() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("totalWins", cel.IntType),
	}
}

func (fo2 *ForceOutcome2) newScriptBasicFuncs() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("get",
			cel.Overload("get_string_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.IntType,
				cel.BinaryBinding(func(param0 ref.Val, param1 ref.Val) ref.Val {
					val := fo2.getComponentVal(param0.Value().(string), param1.Value().(string))

					return types.Int(val)
				},
				),
			),
		),
		cel.Function("getMax",
			cel.Overload("getMax_string_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.IntType,
				cel.BinaryBinding(func(param0 ref.Val, param1 ref.Val) ref.Val {
					val := fo2.getMaxComponentVal(param0.Value().(string), param1.Value().(string))

					return types.Int(val)
				},
				),
			),
		),
		cel.Function("getSymbolVal",
			cel.Overload("getSymbolVal_int_int_int",
				[]*cel.Type{cel.IntType, cel.IntType, cel.IntType},
				cel.IntType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 3 {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getSymbolVal",
							goutils.Err(ErrInvalidScriptParamsNumber))

						return types.Int(0)
					}

					x, isok := params[0].Value().(int64)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getSymbolVal",
							slog.Int("i", 0),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					y, isok := params[1].Value().(int64)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getSymbolVal",
							slog.Int("i", 1),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					defval, isok := params[2].Value().(int64)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getSymbolVal",
							slog.Int("i", 2),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					val := fo2.getSymbolVal(int(x), int(y), int(defval))

					return types.Int(val)
				},
				),
			),
		),
		cel.Function("has",
			cel.Overload("has_string",
				[]*cel.Type{cel.StringType},
				cel.BoolType,
				cel.UnaryBinding(func(param ref.Val) ref.Val {
					if fo2.hasComponent(param.Value().(string)) {
						return types.Bool(true)
					}

					return types.Bool(false)
				},
				),
			),
		),
		cel.Function("getValAt",
			cel.Overload("getValAt_string_string_string",
				[]*cel.Type{cel.StringType, cel.StringType, cel.StringType},
				cel.IntType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 3 {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValAt",
							goutils.Err(ErrInvalidScriptParamsNumber))

						return types.Int(0)
					}

					hasComponent, isok := params[0].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValAt",
							slog.Int("i", 0),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					component, isok := params[1].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValAt",
							slog.Int("i", 1),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					componentVal, isok := params[2].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValAt",
							slog.Int("i", 2),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					val := fo2.getComponentValAt(hasComponent, component, componentVal)

					return types.Int(val)
				},
				),
			),
		),
		cel.Function("getValNext",
			cel.Overload("getValNext_string_string_string",
				[]*cel.Type{cel.StringType, cel.StringType, cel.StringType},
				cel.IntType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 3 {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValNext",
							goutils.Err(ErrInvalidScriptParamsNumber))

						return types.Int(0)
					}

					hasComponent, isok := params[0].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValNext",
							slog.Int("i", 0),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					component, isok := params[1].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValNext",
							slog.Int("i", 1),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					componentVal, isok := params[2].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:getValNext",
							slog.Int("i", 2),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					val := fo2.getComponentValNext(hasComponent, component, componentVal)

					return types.Int(val)
				},
				),
			),
		),
		cel.Function("countSymbolVal",
			cel.Overload("countSymbolVal_string_int",
				[]*cel.Type{cel.StringType, cel.IntType},
				cel.IntType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 2 {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:countSymbolVal",
							goutils.Err(ErrInvalidScriptParamsNumber))

						return types.Int(0)
					}

					op, isok := params[0].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:countSymbolVal",
							slog.Int("i", 0),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					target, isok := params[1].Value().(int64)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:countSymbolVal",
							slog.Int("i", 1),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Int(0)
					}

					val := fo2.countSymbolVal(op, int(target))

					return types.Int(val)
				},
				),
			),
		),
		cel.Function("hasSamePosNext",
			cel.Overload("hasSamePosNext_string_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.BoolType,
				cel.FunctionBinding(func(params ...ref.Val) ref.Val {
					if len(params) != 2 {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:hasSamePosNext",
							goutils.Err(ErrInvalidScriptParamsNumber))

						return types.Bool(false)
					}

					src, isok := params[0].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:hasSamePosNext",
							slog.Int("i", 0),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Bool(false)
					}

					target, isok := params[1].Value().(string)
					if !isok {
						goutils.Error("ForceOutcome2.newScriptBasicFuncs:hasSamePosNext",
							slog.Int("i", 1),
							goutils.Err(ErrInvalidScriptParamType))

						return types.Bool(false)
					}

					val := fo2.hasSamePosNext(src, target)

					return types.Bool(val)
				},
				),
			),
		),
	}
}

func NewForceOutcome2(code string) (*ForceOutcome2, error) {
	fo2 := &ForceOutcome2{}

	options := []cel.EnvOption{}
	options = append(options, fo2.newScriptVariables()...)
	options = append(options, fo2.newScriptBasicFuncs()...)

	cel, err := cel.NewEnv(options...)
	if err != nil {
		goutils.Error("NewForceOutcome2:NewEnv",
			goutils.Err(err))

		return nil, err
	}

	fo2.cel = cel

	err = fo2.SetScript(code)
	if err != nil {
		goutils.Error("NewForceOutcome2:SetScript",
			goutils.Err(err))

		return nil, err
	}

	return fo2, nil
}
