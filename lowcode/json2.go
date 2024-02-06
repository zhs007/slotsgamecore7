package lowcode

import (
	"strings"

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
//	},
//
//	{
//		"type": "addRespinTimes",
//		"target": "fg-start",
//		"times": 15
//	},
//
//	{
//		"triggerNum": "all",
//		"type": "chgComponentConfigIntVal",
//		"target": [
//			"bg-blueeffect",
//			"queue"
//		],
//		"value": 1
//	},
//
//	{
//		"type": "chgComponentConfigIntVal",
//		"targetArr": [
//			"bg-blue",
//			"valueNum"
//		],
//		"valueNum": 0,
//		"source": [
//			"bg-payblue",
//			"symbolNum"
//		]
//	}
//
// ]
type jsonControllerData struct {
	Type       string   `json:"type"`
	StrParams  string   `json:"strParams"`
	Vals       int      `json:"vals"`
	TriggerNum string   `json:"triggerNum"`
	Target     string   `json:"target"`
	TargetArr  []string `json:"targetArr"`
	Value      int      `json:"value"`
	Times      int      `json:"times"`
	ValueNum   int      `json:"valueNum"`
	Source     []string `json:"source"`
}

func (jcd *jsonControllerData) build() *Award {
	if jcd.Type == "addRespinTimes" {
		return &Award{
			AwardType: "respinTimes",
			Vals:      []int{jcd.Times},
			StrParams: []string{jcd.Target},
		}
	} else if jcd.Type == "chgComponentConfigIntVal" {
		if len(jcd.Source) == 0 {
			return &Award{
				AwardType: "chgComponentConfigIntVal",
				Vals:      []int{jcd.ValueNum},
				StrParams: []string{strings.Join(jcd.TargetArr, ".")},
			}
		}

		return &Award{
			AwardType:     "chgComponentConfigIntVal",
			StrParams:     []string{strings.Join(jcd.TargetArr, ".")},
			ComponentVals: []string{strings.Join(jcd.Source, ".")},
		}
	}

	goutils.Error("jsonControllerData.build",
		goutils.JSON("controller", jcd),
		zap.Error(ErrUnsupportedControllerType))

	return nil
}

func (jcd *jsonControllerData) build4Collector() (string, *Award) {
	if jcd.TriggerNum == "" {
		goutils.Error("jsonControllerData.build4Collector",
			goutils.JSON("triggerNum", jcd.TriggerNum),
			zap.Error(ErrUnsupportedControllerType))

		return "", nil
	}

	return jcd.TriggerNum, jcd.build()
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

func parseCollectorControllers(gamecfg *Config, controller *ast.Node) ([]*Award, map[int][]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseControllers:MarshalJSON",
			zap.Error(err))

		return nil, nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseControllers:Unmarshal",
			zap.Error(err))

		return nil, nil, err
	}

	awards := []*Award{}
	mapawards := make(map[int][]*Award)

	for i, v := range lst {
		str, a := v.build4Collector()
		if a != nil {
			if str == "per" {
				awards = append(awards, a)
			} else if str == "all" {
				mapawards[-1] = append(mapawards[-1], a)
			} else {
				i64, err := goutils.String2Int64(str)
				if err != nil {
					goutils.Error("parseControllers:String2Int64",
						zap.Int("i", i),
						zap.String("str", str),
						zap.Error(err))

					return nil, nil, err
				}

				mapawards[int(i64)] = append(mapawards[int(i64)], a)
			}
		} else {
			goutils.Error("parseControllers:build4Collector",
				zap.Int("i", i),
				zap.Error(ErrUnsupportedControllerType))

			return nil, nil, ErrUnsupportedControllerType
		}
	}

	return awards, mapawards, nil
}
