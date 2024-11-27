package sgc7plugin

import "fmt"

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
