package sgc7game

// PlayResult - result for play
type PlayResult struct {
	Scene         *GameScene
	Results       []*Result
	NextCmd       string
	NextCmdParam  string
	NextCmdParams []string
	TotalWin      int
	RealWin       int
}
