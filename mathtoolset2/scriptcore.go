package mathtoolset2

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type ScriptCore struct {
	Cel      *cel.Env
	JsonData string
	ErrInRun []error
}

func (sc *ScriptCore) pushError(err error) {
	sc.ErrInRun = append(sc.ErrInRun, err)
}

func (sc *ScriptCore) Run(code string) error {
	ast, issues := sc.Cel.Compile(code)
	if issues != nil {
		goutils.Error("ScriptCore.Run:Compile",
			zap.String("code", code),
			goutils.JSON("issues", issues),
			zap.Error(ErrInvalidCode))

		return ErrInvalidCode
	}

	prg, err := sc.Cel.Program(ast)
	if err != nil {
		goutils.Error("ScriptCore.Run:Program",
			zap.String("code", code),
			zap.Error(err))

		return err
	}

	// 必须返回一个 bool
	out, _, err := prg.Eval(map[string]any{})
	if err != nil {
		goutils.Error("ScriptCore.Run:Eval",
			zap.Error(err))

		return err
	}

	if !out.Value().(bool) {
		goutils.Error("ScriptCore.Run:result",
			zap.Error(ErrReturnNotOK))

		return ErrReturnNotOK
	}

	if len(sc.ErrInRun) > 0 {
		for _, v := range sc.ErrInRun {
			goutils.Error("ScriptCore.Run:check errors",
				zap.Error(v))
		}

		return ErrRunError
	}

	return nil
}

func (sc *ScriptCore) newGenStackReels() cel.EnvOption {
	return cel.Function("genStackReels",
		cel.Overload("genStackReels_string_string_string_list_int",
			[]*cel.Type{cel.StringType, cel.StringType, cel.StringType, cel.ListType(cel.StringType), cel.IntType},
			cel.BoolType,
			cel.FunctionBinding(func(params ...ref.Val) ref.Val {

				// if len(params) != 5 {
				// 	goutils.Error("genReelsMainSymbolsDistance",
				// 		zap.Error(ErrInvalidFunctionParams))

				// 	return types.Double(0)
				// }

				// targetfn := params[0].Value().(string)
				// paytablefn := params[1].Value().(string)
				// rssfn := params[2].Value().(string)

				// err := mgrGenMath.LoadPaytables(paytablefn)
				// if err != nil {
				// 	goutils.Error("calcWaysRTP2:LoadPaytables",
				// 		zap.Error(err))

				// 	return types.Double(0)
				// }

				// rss, err := LoadReelsStats(rssfn)
				// if err != nil {
				// 	goutils.Error("genReelsMainSymbolsDistance:LoadReelsStats",
				// 		zap.Error(err))

				// 	return types.Double(0)
				// }

				// mainSymbolsWithStr := array2StrSlice(params[3])
				// offset := int(params[4].Value().(int64))

				// mainSymbols := GetSymbols(mainSymbolsWithStr, mgrGenMath.Paytables)

				// reels, err := GenReelsMainSymbolsDistance(rss, mainSymbols, offset, 100)
				// if err != nil {
				// 	goutils.Error("genReelsMainSymbolsDistance:GenReelsMainSymbolsDistance",
				// 		zap.Error(err))

				// 	return types.Double(0)
				// }

				// err = reels.SaveExcelEx(targetfn, mgrGenMath.Paytables)
				// if err != nil {
				// 	goutils.Error("genReelsMainSymbolsDistance:SaveExcelEx",
				// 		zap.Error(err))

				// 	return types.Double(0)
				// }

				return types.Bool(true)
			},
			),
		),
	)
}

func (sc *ScriptCore) newBasicScriptFuncs() []cel.EnvOption {
	return []cel.EnvOption{
		sc.newGenStackReels(),
	}
}

func NewScriptCore(jsonData string) (*ScriptCore, error) {
	scriptCore := &ScriptCore{
		JsonData: jsonData,
	}

	options := []cel.EnvOption{}
	options = append(options, scriptCore.newBasicScriptFuncs()...)

	cel, err := cel.NewEnv(options...)
	if err != nil {
		goutils.Error("NewScriptCore:cel.NewEnv",
			zap.Error(err))

		return nil, err
	}

	scriptCore.Cel = cel

	return scriptCore, nil
}
