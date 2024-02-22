package mathtoolset2

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type ScriptCore struct {
	Cel            *cel.Env
	FileData       string
	MapFiles       *FileDataMap
	ErrInRun       []error
	MapOutputFiles *FileDataMap
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

// newGenStackReels - genStackReels(targetfn string, sourcefn string, stack []int, excludeSymbol []string)
func (sc *ScriptCore) newGenStackReels() cel.EnvOption {
	return cel.Function("genStackReels",
		cel.Overload("genStackReels_string_string_list_list",
			[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.StringType)},
			cel.BoolType,
			cel.FunctionBinding(func(params ...ref.Val) ref.Val {

				if len(params) != 4 {
					goutils.Error("genStackReels",
						zap.Error(ErrInvalidFunctionParams))

					sc.pushError(ErrInvalidFunctionParams)

					return types.Bool(false)
				}

				targetfn := params[0].Value().(string)
				srcfn := params[1].Value().(string)
				stack := List2IntSlice(params[2])
				excludeSymbol := List2StrSlice(params[3])

				fd := sc.MapFiles.GetReader(srcfn)
				if fd == nil {
					goutils.Error("genStackReels:GetReader",
						zap.Error(ErrInvalidFileData))

					sc.pushError(ErrInvalidFileData)

					return types.Bool(false)
				}

				rd, err := GenStackReels(fd, stack, excludeSymbol)
				if err != nil {
					goutils.Error("genStackReels:GenStackReels",
						zap.Error(err))

					sc.pushError(err)

					return types.Bool(false)
				}

				f := NewExcelFile(rd)

				buf, err := f.WriteToBuffer()
				if err != nil {
					goutils.Error("genStackReels:WriteToBuffer",
						zap.Error(err))

					sc.pushError(err)

					return types.Bool(false)
				}

				sc.MapOutputFiles.AddBuffer(targetfn, buf)

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

func NewScriptCore(fileData string) (*ScriptCore, error) {
	mapfd, err := NewFileDataMap(fileData)
	if err != nil {
		goutils.Error("NewScriptCore:NewFileDataMap",
			zap.Error(err))

		return nil, err
	}

	out, err := NewFileDataMap("")
	if err != nil {
		goutils.Error("NewScriptCore:NewFileDataMap:output",
			zap.Error(err))

		return nil, err
	}

	scriptCore := &ScriptCore{
		FileData:       fileData,
		MapFiles:       mapfd,
		MapOutputFiles: out,
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
