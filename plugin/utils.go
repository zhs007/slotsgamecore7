package sgc7plugin

import (
	"fmt"
	"strings"

	"github.com/zhs007/goutils"
)

func GetRngs(plugin IPlugin) []int {
	rngs := []int{}

	lst := plugin.GetUsedRngs()

	for _, v := range lst {
		rngs = append(rngs, v.Value)
	}

	return rngs
}

func GenRngsString(rngs []int) string {
	str := ""

	for _, v := range rngs {
		str += fmt.Sprintf("%v,", v)
	}

	return str
}

func String2Rngs(str string) []int {
	arr := strings.Split(str, ",")
	rngs := []int{}

	for _, v := range arr {
		if v != "" {
			i64, err := goutils.String2Int64(v)
			if err != nil {
				// goutils.Error("String2Rngs:String2Int64",
				// 	slog.String("v", v),
				// 	goutils.Err(err))

				continue
			}

			rngs = append(rngs, int(i64))
		}
	}

	return rngs
}