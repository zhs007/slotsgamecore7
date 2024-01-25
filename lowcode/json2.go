package lowcode

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

func getConfigInCell(cell *ast.Node) (*ast.Node, string, *ast.Node, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("getConfigInCell:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, "", nil, ErrNoComponentValues
	}

	cfg := componentValues.Get("configuration")
	if cfg == nil {
		goutils.Error("getConfigInCell:configuration",
			zap.Error(ErrInvalidJsonNode))

		return nil, "", nil, ErrInvalidJsonNode
	}

	label := componentValues.Get("label")
	if cfg == nil {
		goutils.Error("getConfigInCell:label",
			zap.Error(ErrInvalidJsonNode))

		return nil, "", nil, ErrInvalidJsonNode
	}

	str, err := label.String()
	if cfg == nil {
		goutils.Error("getConfigInCell:label.String",
			zap.Error(err))

		return nil, "", nil, err
	}

	controller := componentValues.Get("controller")

	return cfg, str, controller, nil
}

// "controller": [
//
//	{
//		"type": "AwardRespinTimes",
//		"strParams": "fg-start",
//		"vals": 15
//	}
//
// ]
type jsonControllerData struct {
	Type      string `json:"type"`
	StrParams string `json:"strParams"`
	Vals      int    `json:"vals"`
}

func (jcd *jsonControllerData) build() *Award {
	if jcd.Type == "addRespinTimes" {
		return &Award{
			AwardType: "respinTimes",
			Vals:      []int{jcd.Vals},
			StrParams: []string{jcd.StrParams},
		}
	}

	goutils.Error("jsonControllerData.build",
		goutils.JSON("controller", jcd),
		zap.Error(ErrUnsupportedControllerType))

	return nil
}

func parseControllers(gamecfg *Config, controller *ast.Node) ([]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseControllers:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseControllers:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	awards := []*Award{}

	for i, v := range lst {
		a := v.build()
		if a != nil {
			awards = append(awards, a)
		} else {
			goutils.Error("parseControllers:build",
				zap.Int("i", i),
				zap.Error(ErrUnsupportedControllerType))

			return nil, ErrUnsupportedControllerType
		}
	}

	return awards, nil
}
