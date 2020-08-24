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
	if cmd == "SPIN" {
		return mod.OnSpin(params)
	}

	return ErrInvalidCommand
}

// OnSpin - on spin
func (mod *BasicGameMod) OnSpin(params interface{}) error {
	err := mod.OnRandomScene(params)
	if err != nil {
		return err
	}

	err = mod.OnCalcScene(params)
	if err != nil {
		return err
	}

	err = mod.OnPayout(params)
	if err != nil {
		return err
	}

	return nil
}

// OnRandomScene - on random scene
func (mod *BasicGameMod) OnRandomScene(params interface{}) error {
	return nil
}

// OnCalcScene - on calc scene
func (mod *BasicGameMod) OnCalcScene(params interface{}) error {
	return nil
}

// OnPayout - on payout
func (mod *BasicGameMod) OnPayout(params interface{}) error {
	return nil
}
