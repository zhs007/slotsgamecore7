package gatiserv

import (
	jsoniter "github.com/json-iterator/go"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// BasicService - basic service
type BasicService struct {
	Game     sgc7game.IGame
	GameInfo *GATIGameInfo
}

// NewBasicService - new a BasicService
func NewBasicService(game sgc7game.IGame, gifn string) (*BasicService, error) {

	gi, err := LoadGATIGameInfo(gifn)
	if err != nil {
		sgc7utils.Error("NewBasicService:LoadGATIGameInfo",
			zap.String("gifn", gifn),
			zap.Error(err))

		return nil, err
	}

	return &BasicService{
		Game:     game,
		GameInfo: gi,
	}, nil
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

	plugin := sv.Game.NewPlugin()
	defer sv.Game.FreePlugin(plugin)

	sv.ProcCheat(plugin, params.Cheat)

	stake := BuildStake(params.Stake)
	err = sv.Game.CheckStake(stake)
	if err != nil {
		sgc7utils.Error("BasicService.Play:CheckStake",
			sgc7utils.JSON("stake", stake),
			zap.Error(err))

		return nil, err
	}

	results := []*sgc7game.PlayResult{}

	cmd := params.Cmd

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := sv.Game.Play(plugin, cmd, params.Params, ips, stake, results)
		if err != nil {
			sgc7utils.Error("BasicService.Play:Play",
				zap.Int("results", len(results)),
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

		if len(pr.NextCmds) > 0 {
			cmd = pr.NextCmds[0]
		} else {
			cmd = ""
		}
	}

	pr := &PlayResult{
		RandomNumbers: plugin.GetUsedRngs(),
		// JackpotData   []string             `json:"jackpotData"`
		// AnalyticsData AnalyticsData        `json:"analyticsData"`
		// BoostData     string               `json:"boostData"`
	}

	ps, err := BuildPlayerState(ips)
	if err != nil {
		sgc7utils.Error("BasicService.Play:BuildPlayerState",
			zap.Error(err))

		return nil, err
	}

	pr.PlayerState = ps

	if len(results) > 0 {
		AddPlayResult(pr, params.Stake, results)

		lastr := results[len(results)-1]

		pr.Finished = lastr.IsFinish
		pr.NextCommands = lastr.NextCmds
	}

	return pr, nil
}

// ProcCheat - process cheat
func (sv *BasicService) ProcCheat(plugin sgc7plugin.IPlugin, cheat string) error {
	if cheat != "" {
		str := sgc7utils.AppendString("[", cheat, "]")

		json := jsoniter.ConfigCompatibleWithStandardLibrary

		rngs := []int{}
		err := json.Unmarshal([]byte(str), &rngs)
		if err != nil {
			return err
		}

		plugin.SetCache(rngs)
	}

	return nil
}

// Checksum - checksum
func (sv *BasicService) Checksum(lst []*CriticalComponent) ([]*ComponentChecksum, error) {
	lstret := []*ComponentChecksum{}

	for _, v := range lst {
		cc, isok := sv.GameInfo.Components[v.ID]
		if !isok {
			return nil, ErrInvalidCriticalComponentID
		}

		lstret = append(lstret, &ComponentChecksum{
			ID:       v.ID,
			Checksum: cc.Checksum,
		})
	}

	return lstret, nil
}

// Version - version
func (sv *BasicService) Version() *VersionInfo {
	return &sv.GameInfo.Info
}
