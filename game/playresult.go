package sgc7game

// PlayResult - result for play
type PlayResult struct {
	CurGameMod    string       `json:"curgamemod"`
	NextGameMod   string       `json:"nextgamemod"`
	Scenes        []*GameScene `json:"scenes"`
	Results       []*Result    `json:"results"`
	NextCmds      []string     `json:"-"`
	NextCmdParams []string     `json:"-"`
	CoinWin       int          `json:"-"`
	CashWin       int64        `json:"-"`
	IsFinish      bool         `json:"-"`
	IsWait        bool         `json:"-"`
}
