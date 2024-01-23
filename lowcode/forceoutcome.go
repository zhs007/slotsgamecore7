package lowcode

import (
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type FOData struct {
	Component   string
	Value       string
	Operator    string
	TargetValue int
}

func (fod *FOData) IsValid(lst []*sgc7game.PlayResult) bool {
	return true
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
			zap.String("str", str))

		return nil
	}

	arr1 := strings.Split(arr[0], ".")
	if len(arr1) != 2 {
		goutils.Error("ParseFOData:Split0",
			zap.String("str", str))

		return nil
	}

	i64, err := goutils.String2Int64(arr[2])
	if err != nil {
		goutils.Error("ParseFOData:String2Int64",
			zap.String("str", str),
			zap.Error(err))

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
