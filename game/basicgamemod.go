package sgc7game

// BasicGameMod - basic gameMod
type BasicGameMod struct {
	Name  string
	Scene *GameScene
}

// NewBasicGameMod - new a BasicGameMod
func NewBasicGameMod(name string, width int, height int) BasicGameMod {
	return BasicGameMod{
		Name:  name,
		Scene: NewGameScene(width, height),
	}
}

// GetName - get mode name
func (mod *BasicGameMod) GetName() string {
	return mod.Name
}

// GetGameScene - get GameScene
func (mod *BasicGameMod) GetGameScene() *GameScene {
	return mod.Scene
}

// OnPlay - on play
func (mod *BasicGameMod) OnPlay(cmd string, params interface{}) error {
	return nil
}
