package lowcode

import (
	"log/slog"

	"github.com/google/cel-go/cel"
	"github.com/zhs007/goutils"
)

// CalculatorCore - 基础的计算器核心
type CalculatorCore struct {
	cel     *cel.Env
	program cel.Program
}

func (cc *CalculatorCore) newScriptVariables() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("input1", cel.IntType),
		cel.Variable("input2", cel.IntType),
	}
}

func (cc *CalculatorCore) SetScript(code string) error {
	ast, issues := cc.cel.Compile(code)
	if issues != nil {
		goutils.Error("CalculatorCore.SetScript:Compile",
			slog.String("code", code),
			slog.Any("issues", issues),
			goutils.Err(ErrInvalidForceOutcome2Code))

		return ErrInvalidForceOutcome2Code
	}

	prg, err := cc.cel.Program(ast)
	if err != nil {
		goutils.Error("CalculatorCore.SetScript:Program",
			slog.String("code", code),
			goutils.Err(err))

		return err
	}

	cc.program = prg

	return nil
}

func (cc *CalculatorCore) CalcVal(inputs []int) (int, error) {
	// 必须返回一个 int
	out, _, err := cc.program.Eval(map[string]any{
		"input1": inputs[0],
		"input2": inputs[1],
	})
	if err != nil {
		goutils.Error("CalculatorCore.CalcVal:Eval",
			goutils.Err(err))

		return 0, err
	}

	ret, isok := out.Value().(int64)
	if !isok {
		goutils.Error("CalculatorCore.CalcVal:ret",
			goutils.Err(ErrInvalidForceOutcome2ReturnVal))

		return 0, ErrInvalidForceOutcome2ReturnVal
	}

	return int(ret), nil
}

func NewCalculatorCore(code string) (*CalculatorCore, error) {
	cc := &CalculatorCore{}

	options := []cel.EnvOption{}
	options = append(options, cc.newScriptVariables()...)

	cel, err := cel.NewEnv(options...)
	if err != nil {
		goutils.Error("NewCalculatorCore:NewEnv",
			goutils.Err(err))

		return nil, err
	}

	cc.cel = cel

	err = cc.SetScript(code)
	if err != nil {
		goutils.Error("NewCalculatorCore:SetScript",
			goutils.Err(err))

		return nil, err
	}

	return cc, nil
}
