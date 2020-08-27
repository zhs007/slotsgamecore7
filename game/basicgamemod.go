package sgc7game

// BasicGameMod - basic gameMod
type BasicGameMod struct {
	Name   string
	Width  int
	Height int
}

// NewBasicGameMod - new a BasicGameMod
func NewBasicGameMod(name string, w int, h int) BasicGameMod {
	return BasicGameMod{
		Name:   name,
		Width:  w,
		Height: h,
	}
}

// GetName - get mode name
func (mod *BasicGameMod) GetName() string {
	return mod.Name
}

// OnPlay - on play
func (mod *BasicGameMod) OnPlay(game IGame, cmd string, param string, stake *Stake, prs []*PlayResult) (*PlayResult, error) {
	if cmd == "SPIN" {
		return mod.OnSpin(game, param, stake, prs)
	}

	return nil, ErrInvalidCommand
}

// OnSpin - on spin
func (mod *BasicGameMod) OnSpin(game IGame, param string, stake *Stake, prs []*PlayResult) (*PlayResult, error) {
	pr := &PlayResult{}

	err := mod.OnRandomScene(game, param, prs, pr, mod.Name)
	if err != nil {
		return nil, err
	}

	err = mod.OnCalcScene(game, param, prs, pr)
	if err != nil {
		return nil, err
	}

	err = mod.OnPayout(game, param, prs, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// OnRandomScene - on random scene
func (mod *BasicGameMod) OnRandomScene(game IGame, param string, prs []*PlayResult, pr *PlayResult, reelsName string) error {
	if mod.Width > 0 && mod.Height > 0 {
		pr.Scene = NewGameScene(mod.Width, mod.Height)

		err := pr.Scene.RandReels(game, reelsName)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrInvalidWHGameMod
}

// OnCalcScene - on calc scene
func (mod *BasicGameMod) OnCalcScene(game IGame, param string, prs []*PlayResult, pr *PlayResult) error {
	return nil
}

// OnPayout - on payout
func (mod *BasicGameMod) OnPayout(game IGame, param string, prs []*PlayResult, pr *PlayResult) error {
	return nil
}
