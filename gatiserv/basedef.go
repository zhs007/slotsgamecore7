package gatiserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// PlayerState - player state
type PlayerState struct {
	Public  any `json:"playerStatePublic"`
	Private any `json:"playerStatePrivate"`
}

// Stake - stake
type Stake struct {
	CoinBet  float64 `json:"coinBet"`
	CashBet  float64 `json:"cashBet"`
	Currency string  `json:"currency"`
}

// Jackpot - jackpot
type Jackpot struct {
	Value float64 `json:"value"`
}

// Result - game result
type Result struct {
	CoinWin    int                  `json:"coinWin"`
	CashWin    float64              `json:"cashWin"`
	ClientData *sgc7game.PlayResult `json:"clientData"`
}

// AnalyticEvent - analytic event
type AnalyticEvent struct {
	EventID string `json:"eventId"`
	Data    string `json:"data"`
}

// ResultEvent - result event
type ResultEvent struct {
	ResultIndex int             `json:"resultIndex"`
	Events      []AnalyticEvent `json:"events"`
}

// AnalyticsData - analytics data
type AnalyticsData struct {
	GameEvents   []AnalyticEvent `json:"gameEvents"`
	ResultEvents []ResultEvent   `json:"resultEvents"`
}

// ValidateParams - validate input parameters for the game
type ValidateParams struct {
	PlayerState PlayerState `json:"playerState"`
	Stake       Stake       `json:"stakeValue"`
	Params      string      `json:"clientParams"`
	Cmd         string      `json:"command"`
}

// ValidationError - validate error
type ValidationError struct {
	ErrorCode int    `json:"errorCode"`
	Reason    string `json:"reason"`
	Data      string `json:"data"`
}

// PlayParams - play input parameters for the game
type PlayParams struct {
	PlayerState       *PlayerState       `json:"playerState"`
	Cheat             string             `json:"cheat"`
	Stake             Stake              `json:"stakeValue"`
	Params            string             `json:"clientParams"`
	Cmd               string             `json:"command"`
	JackpotStakeValue float64            `json:"jackpotStakeValue"`
	FreespinsActive   bool               `json:"freespinsActive"`
	JackpotValues     map[string]Jackpot `json:"jackpotValues"`
}

// PlayResult - play output parameters for the game
type PlayResult struct {
	RandomNumbers []*sgc7utils.RngInfo      `json:"randomNumbers"`
	PlayerState   *PlayerState              `json:"playerState"`
	JackpotData   []string                  `json:"jackpotData"`
	Finished      bool                      `json:"finished"`
	Results       []*Result                 `json:"results"`
	NextCommands  []string                  `json:"nextCommands"`
	AnalyticsData *AnalyticsData            `json:"analyticsData"`
	BoostData     *BasicMissionBoostDataMap `json:"boostData"`
}

// CriticalComponent -
type CriticalComponent struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// ComponentChecksum -
type ComponentChecksum struct {
	ID       int    `json:"id"`
	Checksum string `json:"checksum"`
}

// GATICriticalComponent -
type GATICriticalComponent struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
}

// GATICriticalComponents -
type GATICriticalComponents struct {
	Components map[int]*GATICriticalComponent `json:"components"`
}

// VersionInfo -
type VersionInfo struct {
	GameTitle     string `json:"gameTitle"`
	GameVersion   string `json:"gameVersion"`
	VCSVersion    string `json:"vcsVersion"`
	BuildChecksum string `json:"buildChecksum"`
	BuildTime     string `json:"buildTime"`
	Vendor        string `json:"vendor"`
}

// GATIGameInfo - GATIGameInfo
type GATIGameInfo struct {
	Components map[int]*GATICriticalComponent `json:"components"`
	Info       VersionInfo                    `json:"info"`
}

// Checksum - checksum
func (gatiGI *GATIGameInfo) FindComponentChecksum(cc *CriticalComponent) *ComponentChecksum {
	for _, v := range gatiGI.Components {
		if cc.Name == v.Name && cc.Location == v.Location {
			return &ComponentChecksum{
				ID:       v.ID,
				Checksum: v.Checksum,
			}
		}
	}

	return nil
}

// MissionObject -
type MissionObject struct {
	ObjectiveID string `json:"objectiveId"`
	Description string `json:"description"`
	Goal        int    `json:"goal"`
	Period      int    `json:"period"`
}

// GATIGameConfig - game_configuration.json
type GATIGameConfig struct {
	GameObjectives []*MissionObject `json:"gameObjectives"`
}

// EvaluateParams -
type EvaluateParams struct {
	BoostData []*BasicMissionBoostDataMap `json:"boostData"`
	State     *BasicMissionStateMap       `json:"state"`
}

// EvaluateResult -
type EvaluateResult struct {
	Progress int                   `json:"progress"`
	State    *BasicMissionStateMap `json:"state"`
}

// BasicMissionState -
type BasicMissionState struct {
	ObjectiveID string `json:"objectiveId"`
	Goal        int    `json:"goal"`
	Current     int    `json:"current"`
	Arr         []int  `json:"arr"`
}

// BasicMissionBoostData -
type BasicMissionBoostData struct {
	ObjectiveID string `json:"objectiveId"`
	Counter     int    `json:"counter"`
	Arr         []int  `json:"arr"`
	Type        int    `json:"type"` // 0 - counter, 1 - arr
}

// BasicMissionStateMap -
type BasicMissionStateMap struct {
	MapState map[string]*BasicMissionState `json:"mapstate"`
}

// BasicMissionBoostDataMap -
type BasicMissionBoostDataMap struct {
	MapBoostData map[string]*BasicMissionBoostData `json:"mapboostdata"`
}
