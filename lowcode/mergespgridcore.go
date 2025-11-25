package lowcode

import (
	"log/slog"

	"github.com/google/cel-go/cel"
	"github.com/zhs007/goutils"
)

// MergeSPGridCore - MergeSPGrid核心
type MergeSPGridCore struct {
	cel     *cel.Env
	program cel.Program
}

func (cc *MergeSPGridCore) newScriptVariables() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("source1", cel.IntType),
		cel.Variable("source2", cel.IntType),
	}
}

func (cc *MergeSPGridCore) SetScript(code string) error {
	ast, issues := cc.cel.Compile(code)
	if issues != nil {
		goutils.Error("MergeSPGridCore.SetScript:Compile",
			slog.String("code", code),
			slog.Any("issues", issues),
			goutils.Err(ErrInvalidForceOutcome2Code))

		return ErrInvalidForceOutcome2Code
	}

	prg, err := cc.cel.Program(ast)
	if err != nil {
		goutils.Error("MergeSPGridCore.SetScript:Program",
			slog.String("code", code),
			goutils.Err(err))

		return err
	}

	cc.program = prg

	return nil
}

func (cc *MergeSPGridCore) CalcVal(source []int) (int, error) {
	// 必须返回一个 int
	out, _, err := cc.program.Eval(map[string]any{
		"source1": source[0],
		"source2": source[1],
	})
	if err != nil {
		goutils.Error("MergeSPGridCore.CalcVal:Eval",
			goutils.Err(err))

		return 0, err
	}

	ret, isok := out.Value().(int64)
	if !isok {
		goutils.Error("MergeSPGridCore.CalcVal:ret",
			goutils.Err(ErrInvalidForceOutcome2ReturnVal))

		return 0, ErrInvalidForceOutcome2ReturnVal
	}

	return int(ret), nil
}

func NewMergeSPGridCore(code string) (*MergeSPGridCore, error) {
	cc := &MergeSPGridCore{}

	options := []cel.EnvOption{}
	options = append(options, cc.newScriptVariables()...)

	cel, err := cel.NewEnv(options...)
	if err != nil {
		goutils.Error("NewMergeSPGridCore:NewEnv",
			goutils.Err(err))

		return nil, err
	}

	cc.cel = cel

	err = cc.SetScript(code)
	if err != nil {
		goutils.Error("NewMergeSPGridCore:SetScript",
			goutils.Err(err))

		return nil, err
	}

	return cc, nil
}
