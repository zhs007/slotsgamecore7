package lowcode

import (
	"strings"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type ScriptParser struct {
	Text       string
	ScriptName string
	Params     []string
}

func (sp *ScriptParser) Parse() error {
	arr0 := strings.Split(sp.Text, "(")
	if len(arr0) != 2 {
		return ErrNoFunctionInScript
	}

	sp.ScriptName = strings.TrimSpace(arr0[0])

	arr1 := strings.Split(arr0[1], ")")
	if len(arr1) != 2 {
		return ErrWrongFunctionInScript
	}

	arr2 := strings.Split(arr1[0], ",")
	if len(arr2) > 0 {
		for _, v := range arr2 {
			sp.Params = append(sp.Params, strings.TrimSpace(v))
		}
	}

	return nil
}

func NewScriptParser(str string) (*ScriptParser, error) {
	sp := &ScriptParser{
		Text: str,
	}

	err := sp.Parse()
	if err != nil {
		goutils.Error("NewScriptParser:Parse",
			zap.String("text", str),
			zap.Error(err))

		return nil, err
	}

	return sp, nil
}
