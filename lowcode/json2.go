package lowcode

import (
	"log/slog"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
)

func getConfigInCell(cell *ast.Node) (*ast.Node, string, *ast.Node, error) {
	var componentValues *ast.Node
	if cell.Get("componentValues") != nil {
		componentValues = cell.Get("componentValues")
	} else {
		componentValues = cell.Get("data")
	}

	if componentValues == nil {
		goutils.Error("getConfigInCell:componentValues|data",
			goutils.Err(ErrNoComponentValues))

		return nil, "", nil, ErrNoComponentValues
	}

	cfg := componentValues.Get("configuration")
	// if cfg == nil {
	// 	goutils.Error("getConfigInCell:configuration",
	// 		goutils.Err(ErrInvalidJsonNode))

	// 	return nil, "", nil, ErrInvalidJsonNode
	// }

	label := componentValues.Get("label")
	if label == nil {
		goutils.Error("getConfigInCell:label",
			goutils.Err(ErrInvalidJsonNode))

		return nil, "", nil, ErrInvalidJsonNode
	}

	str, err := label.String()
	if err != nil {
		goutils.Error("getConfigInCell:label.String",
			goutils.Err(err))

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
	Type            string   `json:"type"`
	StrParams       string   `json:"strParams"`
	Vals            int      `json:"vals"`
	TriggerNum      string   `json:"triggerNum"`
	Target          string   `json:"target"`
	TargetArr       []string `json:"targetArr"`
	Value           string   `json:"value"`
	Times           int      `json:"times"`
	ValueNum        int      `json:"valueNum"`
	Source          []string `json:"source"`
	StringVal       string   `json:"stringVal"`
	OnTriggerRespin string   `json:"onTriggerRespin"`
}

func (jcd *jsonControllerData) build() *Award {
	if jcd.Type == "addRespinTimes" {
		return &Award{
			AwardType:       "respinTimes",
			Vals:            []int{jcd.Times},
			StrParams:       []string{jcd.Target},
			OnTriggerRespin: jcd.OnTriggerRespin,
		}
	} else if jcd.Type == "chgComponentConfigIntVal" {
		if len(jcd.Source) == 0 {
			return &Award{
				AwardType:       "chgComponentConfigIntVal",
				Vals:            []int{jcd.ValueNum},
				StrParams:       []string{strings.Join(jcd.TargetArr, ".")},
				OnTriggerRespin: jcd.OnTriggerRespin,
			}
		}

		return &Award{
			AwardType:       "chgComponentConfigIntVal",
			StrParams:       []string{strings.Join(jcd.TargetArr, ".")},
			ComponentVals:   []string{strings.Join(jcd.Source, ".")},
			OnTriggerRespin: jcd.OnTriggerRespin,
		}
	} else if jcd.Type == "setComponentConfigIntVal" {
		if len(jcd.Source) == 0 {
			return &Award{
				AwardType:       "setComponentConfigIntVal",
				Vals:            []int{jcd.ValueNum},
				StrParams:       []string{strings.Join(jcd.TargetArr, ".")},
				OnTriggerRespin: jcd.OnTriggerRespin,
			}
		}

		return &Award{
			AwardType:       "setComponentConfigIntVal",
			StrParams:       []string{strings.Join(jcd.TargetArr, ".")},
			ComponentVals:   []string{strings.Join(jcd.Source, ".")},
			OnTriggerRespin: jcd.OnTriggerRespin,
		}
	} else if jcd.Type == "setComponentConfigVal" {
		if len(jcd.Source) == 0 {
			return &Award{
				AwardType:       "setComponentConfigVal",
				StrParams:       []string{strings.Join(jcd.TargetArr, "."), jcd.Value},
				OnTriggerRespin: jcd.OnTriggerRespin,
			}
		}

		return &Award{
			AwardType:       "setComponentConfigVal",
			StrParams:       []string{strings.Join(jcd.TargetArr, ".")},
			ComponentVals:   []string{strings.Join(jcd.Source, ".")},
			OnTriggerRespin: jcd.OnTriggerRespin,
		}
	}

	goutils.Error("jsonControllerData.build",
		slog.Any("controller", jcd),
		goutils.Err(ErrUnsupportedControllerType))

	return nil
}

func (jcd *jsonControllerData) buildWithTriggerNum() (string, *Award) {
	if jcd.TriggerNum == "" {
		goutils.Error("jsonControllerData.buildWithTriggerNum",
			slog.Any("triggerNum", jcd.TriggerNum),
			goutils.Err(ErrUnsupportedControllerType))

		return "", nil
	}

	return jcd.TriggerNum, jcd.build()
}

func (jcd *jsonControllerData) buildWithStringVal() (string, *Award) {
	if jcd.StringVal == "" {
		goutils.Error("jsonControllerData.buildWithStringVal",
			slog.Any("stringVal", jcd.StringVal),
			goutils.Err(ErrUnsupportedControllerType))

		return "", nil
	}

	return jcd.StringVal, jcd.build()
}

func (jcd *jsonControllerData) build4Map() (string, *Award) {
	if jcd.StringVal == "" {
		goutils.Error("jsonControllerData.build4Map",
			slog.Any("stringVal", jcd.StringVal),
			goutils.Err(ErrUnsupportedControllerType))

		return "", nil
	}

	return jcd.StringVal, jcd.build()
}

func parseControllers(controller *ast.Node) ([]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseControllers:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseControllers:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	awards := []*Award{}

	for i, v := range lst {
		a := v.build()
		if a != nil {
			awards = append(awards, a)
		} else {
			goutils.Error("parseControllers:build",
				slog.Int("i", i),
				goutils.Err(ErrUnsupportedControllerType))

			return nil, ErrUnsupportedControllerType
		}
	}

	return awards, nil
}

func parseCollectorControllers(controller *ast.Node) ([]*Award, map[int][]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseControllers:MarshalJSON",
			goutils.Err(err))

		return nil, nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseControllers:Unmarshal",
			goutils.Err(err))

		return nil, nil, err
	}

	awards := []*Award{}
	mapawards := make(map[int][]*Award)

	for i, v := range lst {
		str, a := v.buildWithTriggerNum()
		if a != nil {
			if str == "per" {
				awards = append(awards, a)
			} else if str == "all" {
				mapawards[-1] = append(mapawards[-1], a)
			} else {
				i64, err := goutils.String2Int64(str)
				if err != nil {
					goutils.Error("parseControllers:String2Int64",
						slog.Int("i", i),
						slog.String("str", str),
						goutils.Err(err))

					return nil, nil, err
				}

				mapawards[int(i64)] = append(mapawards[int(i64)], a)
			}
		} else {
			goutils.Error("parseControllers:buildWithTriggerNum",
				slog.Int("i", i),
				goutils.Err(ErrUnsupportedControllerType))

			return nil, nil, ErrUnsupportedControllerType
		}
	}

	return awards, mapawards, nil
}

func parseIntValAndAllControllers(controller *ast.Node) ([]*Award, map[int][]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValAndAllControllers:MarshalJSON",
			goutils.Err(err))

		return nil, nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseIntValAndAllControllers:Unmarshal",
			goutils.Err(err))

		return nil, nil, err
	}

	awards := []*Award{}
	mapawards := make(map[int][]*Award)

	for i, v := range lst {
		str, a := v.buildWithStringVal()
		if a != nil {
			if str == "all" || str == "" {
				awards = append(awards, a)
			} else {
				i64, err := goutils.String2Int64(str)
				if err != nil {
					goutils.Error("parseIntValAndAllControllers:String2Int64",
						slog.Int("i", i),
						slog.String("str", str),
						goutils.Err(err))

					return nil, nil, err
				}

				mapawards[int(i64)] = append(mapawards[int(i64)], a)
			}
		} else {
			goutils.Error("parseIntValAndAllControllers:buildWithTriggerNum",
				slog.Int("i", i),
				goutils.Err(ErrUnsupportedControllerType))

			return nil, nil, ErrUnsupportedControllerType
		}
	}

	return awards, mapawards, nil
}

func parseMapControllers(controller *ast.Node) (map[string][]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseMapControllers:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseMapControllers:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	mapawards := make(map[string][]*Award)

	for i, v := range lst {
		str, a := v.build4Map()
		if a != nil {
			mapawards[str] = append(mapawards[str], a)
		} else {
			goutils.Error("parseMapControllers:build4Map",
				slog.Int("i", i),
				goutils.Err(ErrUnsupportedControllerType))

			return nil, ErrUnsupportedControllerType
		}
	}

	return mapawards, nil
}

func parseReelTriggerControllers(controller *ast.Node) (map[int][]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseReelTriggerControllers:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseReelTriggerControllers:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	mapawards := make(map[int][]*Award)

	for i, v := range lst {
		str, a := v.buildWithTriggerNum()
		if a != nil {
			if strings.HasPrefix(str, "row") {
				arr := strings.Split(str, "row")
				if len(arr) == 2 {
					i64, err := goutils.String2Int64(arr[1])
					if err != nil {
						goutils.Error("parseReelTriggerControllers:String2Int64",
							slog.Int("i", i),
							slog.String("str", str),
							goutils.Err(err))

						return nil, err
					}

					mapawards[int(i64)] = append(mapawards[int(i64)], a)
				}
			} else if strings.HasPrefix(str, "column") {
				arr := strings.Split(str, "column")
				if len(arr) == 2 {
					i64, err := goutils.String2Int64(arr[1])
					if err != nil {
						goutils.Error("parseReelTriggerControllers:String2Int64",
							slog.Int("i", i),
							slog.String("str", str),
							goutils.Err(err))

						return nil, err
					}

					mapawards[int(i64)] = append(mapawards[int(i64)], a)
				}
			}

		} else {
			goutils.Error("parseReelTriggerControllers:buildWithTriggerNum",
				slog.Int("i", i),
				goutils.Err(ErrUnsupportedControllerType))

			return nil, ErrUnsupportedControllerType
		}
	}

	return mapawards, nil
}

func parseMaskControllers(controller *ast.Node) ([]*Award, map[int][]*Award, error) {
	buf, err := controller.MarshalJSON()
	if err != nil {
		goutils.Error("parseMaskControllers:MarshalJSON",
			goutils.Err(err))

		return nil, nil, err
	}

	lst := []*jsonControllerData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseMaskControllers:Unmarshal",
			goutils.Err(err))

		return nil, nil, err
	}

	mapawards := make(map[int][]*Award)
	perAwards := []*Award{}

	for i, v := range lst {
		str, a := v.buildWithTriggerNum()
		if a != nil {
			if str == "per" {
				perAwards = append(perAwards, a)
			} else if str == "all" {
				mapawards[-1] = append(mapawards[-1], a)
			} else {
				i64, err := goutils.String2Int64(str)
				if err != nil {
					goutils.Error("parseMaskControllers:String2Int64",
						slog.String("str", str),
						goutils.Err(err))
				}

				mapawards[int(i64)-1] = append(mapawards[int(i64)-1], a)
			}
		} else {
			goutils.Error("parseMaskControllers:build4Map",
				slog.Int("i", i),
				goutils.Err(ErrUnsupportedControllerType))

			return nil, nil, ErrUnsupportedControllerType
		}
	}

	return perAwards, mapawards, nil
}
