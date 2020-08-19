package sgc7game

// GameI - game
type GameI interface {
	// GetConfig - get config
	GetConfig() *Config
	// Initialize - initialize PlayerState
	Initialize() *PlayerState
}
