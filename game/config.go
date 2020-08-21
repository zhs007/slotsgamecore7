package sgc7game

// Config - config
type Config struct {
	Line      *LineData             `json:"line"`
	Reels     map[string]*ReelsData `json:"reels"`
	PayTables *PayTables            `json:"paytables"`
	Width     int                   `json:"width"`
	Height    int                   `json:"height"`
}
