package mathtoolset2

import "github.com/google/cel-go/common/types/ref"

func List2IntSlice(val ref.Val) []int {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []int{}

		for _, v := range lst0 {
			v1, isok := v.Value().(int64)
			if isok {
				lst = append(lst, int(v1))
			}
		}

		return lst
	}

	return nil
}

func List2StrSlice(val ref.Val) []string {
	lst0, isok := val.Value().([]ref.Val)
	if isok {
		lst := []string{}

		for _, v := range lst0 {
			v1, isok := v.Value().(string)
			if isok {
				lst = append(lst, v1)
			}
		}

		return lst
	}

	return nil
}
