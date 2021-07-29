package sgc7utils

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

// IndexOfStringSlice - indexof for []string
func IndexOfStringSlice(arr []string, v string, start int) int {
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

// InsUniqueIntSlice - Insert unique int array
func InsUniqueIntSlice(arr []int, v int) []int {
	if IndexOfIntSlice(arr, v, 0) >= 0 {
		return arr
	}

	return append(arr, v)
}

// IntArr2ToInt32Arr - [][]int -> []int32
func IntArr2ToInt32Arr(arr [][]int) []int32 {
	arr2 := []int32{}

	for _, arr1 := range arr {
		for _, v := range arr1 {
			arr2 = append(arr2, int32(v))
		}
	}

	return arr2
}
