package lowcode

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/zhs007/goutils"
)

type ScriptCore struct {
	Cel *cel.Env
}

func newScriptVariables(*GameProperty) []cel.EnvOption {
	return []cel.EnvOption{}
}

func newScriptBasicFuncs(gameProp *GameProperty) []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("setVal",
			cel.Overload("setVal_string_int",
				[]*cel.Type{cel.StringType, cel.IntType},
				cel.NullType,
				cel.BinaryBinding(func(param0, param1 ref.Val) ref.Val {
					prop, err := String2Property(param0.Value().(string))
					if err != nil {
						goutils.Error("newScriptBasicFuncs:setVal:String2Property",
							goutils.Err(err))

						return types.NullType
					}

					err = gameProp.SetVal(prop, param1.Value().(int))
					if err != nil {
						goutils.Error("newScriptBasicFuncs:setVal",
							goutils.Err(err))

						return types.NullType
					}

					return types.NullType
				},
				),
			),
		),
		cel.Function("setStrVal",
			cel.Overload("setStrVal_string_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.NullType,
				cel.BinaryBinding(func(param0, param1 ref.Val) ref.Val {
					prop, err := String2Property(param0.Value().(string))
					if err != nil {
						goutils.Error("newScriptBasicFuncs:setStrVal:String2Property",
							goutils.Err(err))

						return types.NullType
					}

					err = gameProp.SetStrVal(prop, param1.Value().(string))
					if err != nil {
						goutils.Error("newScriptBasicFuncs:SetStrVal",
							goutils.Err(err))

						return types.NullType
					}

					return types.NullType
				},
				),
			),
		),
	}
}

func NewScriptCore(gameProp *GameProperty) (*ScriptCore, error) {
	options := []cel.EnvOption{}
	options = append(options, newScriptVariables(gameProp)...)
	options = append(options, newScriptBasicFuncs(gameProp)...)

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
