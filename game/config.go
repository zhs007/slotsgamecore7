package sgc7game

// Config - config
type Config struct {
	Line      LineData  `json:"line"`
	Reels     ReelsData `json:"reels"`
	PayTables PayTables `json:"paytables"`
}
