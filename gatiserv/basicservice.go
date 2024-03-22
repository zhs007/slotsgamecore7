package gatiserv

import (
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// BasicService - basic service
type BasicService struct {
	Game       sgc7game.IGame
	GameInfo   *GATIGameInfo
	GameConfig *GATIGameConfig
}

// NewBasicService - new a BasicService
func NewBasicService(game sgc7game.IGame, gifn string, gcfn string) (*BasicService, error) {

	gi, err := LoadGATIGameInfo(gifn)
	if err != nil {
		goutils.Error("NewBasicService:LoadGATIGameInfo",
			slog.String("gifn", gifn),
			goutils.Err(err))

		return nil, err
	}

	gc, err := LoadGATIGameConfig(gcfn)
	if err != nil {
		curpwd, _ := os.Getwd()

		goutils.Error("NewBasicService:LoadGATIGameConfig",
			slog.String("gcfn", gcfn),
			slog.String("pwd", curpwd),
			goutils.Err(err))

		return nil, err
	}

	return &BasicService{
		Game:       game,
		GameInfo:   gi,
		GameConfig: gc,
	}, nil
}

// Config - get configuration
func (sv *BasicService) Config() *sgc7game.Config {
	return sv.Game.GetConfig()
}

// Initialize - initialize a player
func (sv *BasicService) Initialize() *PlayerState {
	ips := sv.Game.Initialize()

	return &PlayerState{
		Public:  ips.GetPublic(),
		Private: ips.GetPrivate(),
	}
}

// Validate - validate game
func (sv *BasicService) Validate(params *ValidateParams) []ValidationError {
	return []ValidationError{}
}

// Play - play game
func (sv *BasicService) Play(params *PlayParams) (*PlayResult, error) {
	ips := sv.Game.NewPlayerState()
	if params.PlayerState != nil {
		err := BuildIPlayerState(ips, params.PlayerState)
		if err != nil {
			goutils.Error("BasicService.Play:BuildIPlayerState",
				slog.Any("PlayerState", params.PlayerState),
				goutils.Err(err))

			return nil, err
		}
	}

	plugin := sv.Game.NewPlugin()
	defer sv.Game.FreePlugin(plugin)

	sv.ProcCheat(plugin, params.Cheat)

	stake := BuildStake(params.Stake)
	err := sv.Game.CheckStake(stake)
	if err != nil {
		goutils.Error("BasicService.Play:CheckStake",
			slog.Any("stake", stake),
			goutils.Err(err))

		return nil, err
	}

	results := []*sgc7game.PlayResult{}
	gameData := sv.Game.NewGameData(stake)
	defer sv.Game.DeleteGameData(gameData)

	cmd := params.Cmd

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := sv.Game.Play(plugin, cmd, params.Params, ips, stake, results, gameData)
		if err != nil {
			goutils.Error("BasicService.Play:Play",
				slog.Int("results", len(results)),
				goutils.Err(err))

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
		goutils.Error("BasicService.Play:BuildPlayerState",
			goutils.Err(err))

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
		str := goutils.AppendString("[", cheat, "]")

		rngs := []int{}
		err := sonic.Unmarshal([]byte(str), &rngs)
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
		cc := sv.GameInfo.FindComponentChecksum(v)
		if cc != nil {
			lstret = append(lstret, cc)
		}
		// cc, isok := sv.GameInfo.Components[v.ID]
		// if !isok {
		// 	return nil, ErrInvalidCriticalComponentID
		// }
	}

	return lstret, nil
}

// Version - version
func (sv *BasicService) Version() *VersionInfo {
	return &sv.GameInfo.Info
}

// // NewBoostData - new a BoostData
// func (sv *BasicService) NewBoostData() any {
// 	return nil
// }

// // NewBoostDataList - new a list for BoostData
// func (sv *BasicService) NewBoostDataList() []any {
// 	return []*BasicMissionBoostDataMap{}
// }

// // NewPlayerBoostData - new a PlayerBoostData
// func (sv *BasicService) NewPlayerBoostData() any {
// 	return nil
// }

// OnPlayBoostData - after call Play
func (sv *BasicService) OnPlayBoostData(params *PlayParams, result *PlayResult) error {
	return nil
}

// GetGameConfig - get GATIGameConfig
func (sv *BasicService) GetGameConfig() *GATIGameConfig {
	return sv.GameConfig
}

// Evaluate -
func (sv *BasicService) Evaluate(params *EvaluateParams, id string) (*EvaluateResult, error) {
	result := &EvaluateResult{}

	var mc *MissionObject
	for _, v := range sv.GameConfig.GameObjectives {
		if v.ObjectiveID == id {
			mc = v

			break
		}
	}

	if mc == nil {
		return nil, ErrInvalidObjectiveID
	}

	var cs *BasicMissionState
	if params.State == nil {
		result.State = &BasicMissionStateMap{
			MapState: make(map[string]*BasicMissionState),
		}

		cs = &BasicMissionState{
			ObjectiveID: id,
			Goal:        mc.Goal,
		}
	} else {
		cs1, isok := params.State.MapState[id]
		if !isok {
			cs = &BasicMissionState{
				ObjectiveID: id,
				Goal:        mc.Goal,
			}
		} else {
			cs = cs1

			cs.Goal = mc.Goal
		}
	}

	result.State = &BasicMissionStateMap{
		MapState: make(map[string]*BasicMissionState),
	}

	result.State.MapState[id] = cs

	for _, v := range params.BoostData {
		cbd, isok := v.MapBoostData[id]
		if isok {
			if cbd.Type == 0 {
				cs.Current += cbd.Counter

				result.Progress += cbd.Counter
			} else if cbd.Type == 1 {
				ln := len(cs.Arr)
				for _, n := range cbd.Arr {
					hasn := false
					for _, sn := range cs.Arr {
						if n == sn {
							hasn = true

							break
						}
					}

					if !hasn {
						cs.Arr = append(cs.Arr, n)
					}
				}

				result.Progress += len(cs.Arr) - ln
			}
		}
	}

	// if result.Progress <= 0 {
	// 	if cs.Current > 0 {
	// 		result.Progress = cs.Current
	// 	} else if len(cs.Arr) > 0 {
	// 		result.Progress = len(cs.Arr)
	// 	}
	// }

	return result, nil
}
