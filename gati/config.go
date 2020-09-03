package gati

import "os"

// Config - configuration
type Config struct {
	GameID  string
	RNGURL  string
	RngNums int
}

// NewConfig - new a config
func NewConfig(gameid string, rngnums int) *Config {
	rngservaddr := os.Getenv("RNG_SERVICE_ADDRESS")

	return &Config{
		GameID:  gameid,
		RNGURL:  rngservaddr + "/numbers",
		RngNums: rngnums,
	}
}
