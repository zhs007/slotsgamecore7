package mathtoolset2

import (
	"io"
	"log/slog"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
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
			slog.String("code", code),
			slog.Any("issues", issues),
			goutils.Err(ErrInvalidCode))

		return ErrInvalidCode
	}

	prg, err := sc.Cel.Program(ast)
	if err != nil {
		goutils.Error("ScriptCore.Run:Program",
			slog.String("code", code),
			goutils.Err(err))

		return err
	}

	// 必须返回一个 bool
	out, _, err := prg.Eval(map[string]any{})
	if err != nil {
		goutils.Error("ScriptCore.Run:Eval",
			goutils.Err(err))

		return err
	}

	if !out.Value().(bool) {
		goutils.Error("ScriptCore.Run:result",
			goutils.Err(ErrReturnNotOK))

		return ErrReturnNotOK
	}

	if len(sc.ErrInRun) > 0 {
		for _, v := range sc.ErrInRun {
			goutils.Error("ScriptCore.Run:check errors",
				goutils.Err(v))
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
						goutils.Err(ErrInvalidFunctionParams))

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
						goutils.Err(ErrInvalidFileData))

					sc.pushError(ErrInvalidFileData)

					return types.Bool(false)
				}

				rd, err := GenStackReels(fd, stack, excludeSymbol)
				if err != nil {
					goutils.Error("genStackReels:GenStackReels",
						goutils.Err(err))

					sc.pushError(err)

					return types.Bool(false)
				}

				f := NewExcelFile(rd)

				buf, err := f.WriteToBuffer()
				if err != nil {
					goutils.Error("genStackReels:WriteToBuffer",
						goutils.Err(err))

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

// newGenStackReelsStrict - genStackReelsStrict(targetfn string, sourcefn string, stack []int, excludeSymbol []string)
func (sc *ScriptCore) newGenStackReelsStrict() cel.EnvOption {
	return cel.Function("genStackReelsStrict",
		cel.Overload("genStackReelsStrict_string_string_list_list",
			[]*cel.Type{cel.StringType, cel.StringType, cel.ListType(cel.IntType), cel.ListType(cel.StringType)},
			cel.BoolType,
			cel.FunctionBinding(func(params ...ref.Val) ref.Val {

				if len(params) != 4 {
					goutils.Error("genStackReelsStrict",
						goutils.Err(ErrInvalidFunctionParams))

					sc.pushError(ErrInvalidFunctionParams)

					return types.Bool(false)
				}

				targetfn := params[0].Value().(string)
				srcfn := params[1].Value().(string)
				stack := List2IntSlice(params[2])
				excludeSymbol := List2StrSlice(params[3])

				fd := sc.MapFiles.GetReader(srcfn)
				if fd == nil {
					goutils.Error("genStackReelsStrict:GetReader",
						goutils.Err(ErrInvalidFileData))

					sc.pushError(ErrInvalidFileData)

					return types.Bool(false)
				}

				rd, err := genStackReelsStrict(fd, stack, excludeSymbol)
				if err != nil {
					goutils.Error("genStackReelsStrict:GenStackReels",
						goutils.Err(err))

					sc.pushError(err)

					return types.Bool(false)
				}

				f := NewExcelFile(rd)

				buf, err := f.WriteToBuffer()
				if err != nil {
					goutils.Error("genStackReelsStrict:WriteToBuffer",
						goutils.Err(err))

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

// newGenStackReels - mergeReels(targetfn string, files []string)
func (sc *ScriptCore) newMergeReels() cel.EnvOption {
	return cel.Function("mergeReels",
		cel.Overload("mergeReels_string_list",
			[]*cel.Type{cel.StringType, cel.ListType(cel.StringType)},
			cel.BoolType,
			cel.FunctionBinding(func(params ...ref.Val) ref.Val {

				if len(params) != 2 {
					goutils.Error("mergeReels",
						goutils.Err(ErrInvalidFunctionParams))

					sc.pushError(ErrInvalidFunctionParams)

					return types.Bool(false)
				}

				targetfn := params[0].Value().(string)
				files := List2StrSlice(params[1])
				readers := []io.Reader{}

				for _, v := range files {
					readers = append(readers, sc.MapFiles.GetReader(v))
				}

				rd, err := MergeReels(readers)
				if err != nil {
					goutils.Error("mergeReels:MergeReels",
						goutils.Err(err))

					sc.pushError(err)

					return types.Bool(false)
				}

				f := NewExcelFile(rd)

				buf, err := f.WriteToBuffer()
				if err != nil {
					goutils.Error("mergeReels:WriteToBuffer",
						goutils.Err(err))

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
		sc.newGenStackReelsStrict(),
		sc.newMergeReels(),
	}
}

func NewScriptCore(fileData string) (*ScriptCore, error) {
	mapfd, err := NewFileDataMap(fileData)
	if err != nil {
		goutils.Error("NewScriptCore:NewFileDataMap",
			goutils.Err(err))

		return nil, err
	}

	out, err := NewFileDataMap("")
	if err != nil {
		goutils.Error("NewScriptCore:NewFileDataMap:output",
			goutils.Err(err))

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
			goutils.Err(err))

		return nil, err
	}

	scriptCore.Cel = cel

	return scriptCore, nil
}
