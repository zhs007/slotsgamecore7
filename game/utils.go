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

// IndexOfInt2Slice - indexof for []int2, []int2 is like [x0, y0, x1, y1, ...]
//		start * 2 <--> len([]int)
func IndexOfInt2Slice(arr []int, x, y int, start int) int {
	if start < 0 {
		start = 0
	}

	for i := start * 2; i < len(arr); i += 2 {
		if arr[i] == x && arr[i+1] == y {
			return i / 2
		}
	}

	return -1
}
