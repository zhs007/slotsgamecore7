package gatiserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// BasicService - basic service
type BasicService struct {
	Game sgc7game.IGame
}

// NewBasicService - new a BasicService
func NewBasicService(game sgc7game.IGame) *BasicService {
	return &BasicService{
		Game: game,
	}
}

// Config - get configuration
func (sv *BasicService) Config() *sgc7game.Config {
	return sv.Game.GetConfig()
}

// Initialize - initialize a player
func (sv *BasicService) Initialize() sgc7game.IPlayerState {
	return sv.Game.Initialize()
}

// Validate - validate game
func (sv *BasicService) Validate(params *ValidateParams) []ValidationError {
	return nil
}

// Play - play game
func (sv *BasicService) Play(params *PlayParams) (*PlayResult, error) {
	ips := sv.Game.NewPlayerState()
	err := BuildIPlayerState(ips, params.PlayerState)
	if err != nil {
		sgc7utils.Error("BasicService.Play:BuildIPlayerState",
			zap.Error(err))

		return nil, err
	}

	stake := BuildStake(params.Stake)

	results := []*sgc7game.PlayResult{}

	for {
		pr, err := sv.Game.Play(params.Cmd, params.Params, ips, stake, results)
		if err != nil {
			sgc7utils.Error("BasicService.Play:Play",
				zap.Error(err))

			return nil, err
		}

		if pr == nil {
			break
		}

		results = append(results, pr)
		if pr.IsFinish {
			break
		}

		if pr.IsWait {
			break
		}
	}

	pr := &PlayResult{}

	AddPlayResult(pr, params.Stake, results)

	return pr, nil
}
