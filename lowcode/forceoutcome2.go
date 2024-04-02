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

	// 必须返回一个 bool
	out, _, err := fo2.program.Eval(map[string]any{})
	if err != nil {
		goutils.Error("ForceOutcome2.Check:Eval",
			goutils.Err(err))

		return false
	}

	ret, isok := out.Value().(bool)
	if !isok {
		goutils.Error("ForceOutcome2.Check:ret",
			goutils.Err(ErrInvalidForceOutcome2CheckErr))

		return false
	}

	return ret
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

func (fo2 *ForceOutcome2) getComponentVal(component string, val string) int {
	for _, ret := range fo2.results {
		gp, isok := ret.CurGameModParams.(*GameParams)
		if isok {
			for k, v := range gp.MapComponents {
				if isComponent(k, component) {
					curval, isok2 := GetComponentDataVal(v, val)
					if isok2 {
						return curval
					}
				}
			}
		}
	}

	return 0
}

func (fo2 *ForceOutcome2) newScriptVariables() []cel.EnvOption {
	return []cel.EnvOption{}
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
