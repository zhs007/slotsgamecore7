package sgc7game

// IndexOfIntSlice - indexof for []int
func IndexOfIntSlice(arr []int, v int, start int) int {
	if start < 0 {
		start = 0
	}

	for i := start; i < len(arr); i++ {
		if arr[i] == v {
			return i
		}
	}

	return -1
}
