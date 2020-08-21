package sgc7game

// IGameMod - game
type IGameMod interface {
	// GetName - get mode name
	GetName() string

	// GetGameScene - get GameScene
	GetGameScene() *GameScene

	// OnPlay - on play
	OnPlay(cmd string, params interface{}) error
}
