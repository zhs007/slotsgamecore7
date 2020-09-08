package sgc7game

// BasicGameMod - basic gameMod
type BasicGameMod struct {
	Name   string
	Width  int
	Height int
}

// NewBasicGameMod - new a BasicGameMod
func NewBasicGameMod(name string, w int, h int) *BasicGameMod {
	return &BasicGameMod{
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
	return nil, ErrInvalidCommand
}

// RandomScene - on random scene
func (mod *BasicGameMod) RandomScene(game IGame, param string, prs []*PlayResult, pr *PlayResult, reelsName string) error {
	if mod.Width > 0 && mod.Height > 0 {
		cs := NewGameScene(mod.Width, mod.Height)

		err := cs.RandReels(game, reelsName)
		if err != nil {
			return err
		}

		pr.Scenes = append(pr.Scenes, cs)

		return nil
	}

	return ErrInvalidWHGameMod
}
