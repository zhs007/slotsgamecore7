package stats2

var gWinRange []int

func SetWinRange(winRange []int) {
	gWinRange = winRange
}

func init() {
	SetWinRange([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200, 500, 1000, 2000, 3000, 5000})
}
