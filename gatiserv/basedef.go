package gatiserv

import (
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// PlayerState - player state
type PlayerState struct {
	Public  string `json:"playerStatePublic"`
	Private string `json:"playerStatePrivate"`
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
	CoinWin    int     `json:"coinWin"`
	CashWin    float64 `json:"cashWin"`
	ClientData string  `json:"clientData"`
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
	PlayerState       PlayerState        `json:"playerState"`
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
	RandomNumbers []*sgc7utils.RngInfo `json:"randomNumbers"`
	PlayerState   *PlayerState         `json:"playerState"`
	JackpotData   []string             `json:"jackpotData"`
	Finished      bool                 `json:"finished"`
	Results       []*Result            `json:"results"`
	NextCommands  []string             `json:"nextCommands"`
	AnalyticsData AnalyticsData        `json:"analyticsData"`
	BoostData     string               `json:"boostData"`
}
