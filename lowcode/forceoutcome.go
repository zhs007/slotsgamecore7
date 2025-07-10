package lowcode

import (
	"log/slog"
	"strings"

	any1 "github.com/golang/protobuf/ptypes/any"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func cmpVal(srcVal int, op string, targetVal int) bool {
	switch op {
	case "==":
		return srcVal == targetVal
	case ">=":
		return srcVal >= targetVal
	case "<=":
		return srcVal <= targetVal
	case ">":
		return srcVal > targetVal
	case "<":
		return srcVal < targetVal
	}

	return false
}

func isComponent(target string, src string) bool {
	if src == target {
		return true
	}

	arr0 := strings.Split(target, ":")
	if len(arr0) > 1 {
		if arr0[0] == src {
			return true
		}
	}

	arr1 := strings.Split(target, "/")
	if len(arr1) > 1 {
		if arr1[len(arr1)-1] == src {
			return true
		}
	}

	return false
}

func hasComponentInHistory(lst []string, component string) bool {
	for _, v := range lst {
		if isComponent(v, component) {
			return true
		}
	}

	return false
}

func checkComponentVal(mapComponent map[string]*any1.Any, component string, val string, op string, targetVal int) bool {
	for k, v := range mapComponent {
		if isComponent(k, component) {
			curval, isok2 := GetComponentDataVal(v, val)
			if isok2 {
				if cmpVal(curval, op, targetVal) {
					return true
				}
			}
		}
	}

	return false
}

type FOData struct {
	Component   string
	Value       string
	Operator    string
	TargetValue int
}

func (fod *FOData) IsValid(lst []*sgc7game.PlayResult) bool {
	if fod.Value == "" {
		for _, pr := range lst {
			gp := pr.CurGameModParams.(*GameParams)
			if hasComponentInHistory(gp.HistoryComponents, fod.Component) {
				return true
			}
		}
	} else {
		for _, pr := range lst {
			gp := pr.CurGameModParams.(*GameParams)
			if checkComponentVal(gp.MapComponents, fod.Component, fod.Value, fod.Operator, fod.TargetValue) {
				return true
			}

			// cdpb, isok := gp.MapComponents[fod.Component]
			// if isok {
			// 	val, isok2 := GetComponentDataVal(cdpb, fod.Value)
			// 	if isok2 {
			// 		if cmpVal(val, fod.Operator, fod.TargetValue) {
			// 			return true
			// 		}
			// 	}
			// }
		}
	}

	return false
}

// parse a
// parse a.b >= 1
func ParseFOData(str string) *FOData {
	arr := strings.Split(str, " ")
	if len(arr) == 1 {
		return &FOData{
			Component: arr[0],
		}
	}

	if len(arr) != 3 {
		goutils.Error("ParseFOData:Split",
			slog.String("str", str))

		return nil
	}

	arr1 := strings.Split(arr[0], ".")
	if len(arr1) != 2 {
		goutils.Error("ParseFOData:Split0",
			slog.String("str", str))

		return nil
	}

	i64, err := goutils.String2Int64(arr[2])
	if err != nil {
		goutils.Error("ParseFOData:String2Int64",
			slog.String("str", str),
			goutils.Err(err))

		return nil
	}

	return &FOData{
		Component:   arr1[0],
		Value:       arr1[1],
		Operator:    arr[1],
		TargetValue: int(i64),
	}
}

type ForceOutcome struct {
	Data []*FOData
}

func (fo *ForceOutcome) IsValid(lst []*sgc7game.PlayResult) bool {
	for _, v := range fo.Data {
		if !v.IsValid(lst) {
			return false
		}
	}

	return true
}

func ParseForceOutcome(str string) *ForceOutcome {
	fod := ParseFOData(str)

	if fod == nil {
		return nil
	}

	return &ForceOutcome{
		Data: []*FOData{fod},
	}
}
