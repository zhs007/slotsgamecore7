package gatiserv

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// IService - service
type IService interface {
	// Config - get configuration
	Config() *sgc7game.Config
	// Initialize - initialize a player
	Initialize() *PlayerState
	// Validate - validate game
	Validate(params *ValidateParams) []ValidationError
	// Play - play game
	Play(params *PlayParams) (*PlayResult, error)
	// Checksum - checksum
	Checksum(lst []*CriticalComponent) ([]*ComponentChecksum, error)
	// Version - version
	Version() *VersionInfo
	// OnPlayBoostData - after call Play
	OnPlayBoostData(params *PlayParams, result *PlayResult) error
	// GetGameConfig - get GATIGameConfig
	GetGameConfig() *GATIGameConfig
	// Evaluate -
	Evaluate(params *EvaluateParams, id string) (*EvaluateResult, error)
}
