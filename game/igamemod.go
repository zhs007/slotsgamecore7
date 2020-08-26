package sgc7game

// IGameMod - game
type IGameMod interface {
	// GetName - get mode name
	GetName() string

	// OnPlay - on play
	OnPlay(game IGame, cmd string, param string, prs []*PlayResult) (*PlayResult, error)
}
