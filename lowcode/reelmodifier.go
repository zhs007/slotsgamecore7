package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"github.com/zhs007/slotsgamecore7/stats2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const ReelModifierTypeName = "reelModifier"

const (
	maxRetryNum int = 100
)

// ReelModifierConfig - configuration for ReelModifier feature
type ReelModifierConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Reel                 string              `yaml:"reel" json:"reel"`               // 用这个轮子roll
	ReelData             *sgc7game.ReelsData `yaml:"-" json:"-"`                     // 用这个轮子roll
	Mask                 string              `yaml:"mask" json:"mask"`               // 如果mask不为空，则用这个mask的1来roll，可以配置 isReverse 来roll 0
	IsReverse            bool                `yaml:"isReverse" json:"isReverse"`     // 如果isReverse，表示roll 0
	HoldSymbols          []string            `yaml:"holdSymbols" json:"holdSymbols"` // 这些符号保留
	HoldSymbolCodes      []int               `yaml:"-" json:"-"`
	Triggers             []string            `yaml:"triggers" json:"triggers"` // 替换完轮子后需要保证所有trigger返回true
}

type ReelModifier struct {
	*BasicComponent `json:"-"`
	Config          *ReelModifierConfig `json:"config"`
}

// Init -
func (reelModifier *ReelModifier) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReelModifier.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReelModifierConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReelModifier.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return reelModifier.InitEx(cfg, pool)
}

// InitEx -
func (reelModifier *ReelModifier) InitEx(cfg any, pool *GamePropertyPool) error {
	reelModifier.Config = cfg.(*ReelModifierConfig)
	reelModifier.Config.ComponentType = ReelModifierTypeName

	for _, s := range reelModifier.Config.HoldSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ReelModifier.InitEx:HoldSymbols",
				zap.String("symbol", s),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		reelModifier.Config.HoldSymbolCodes = append(reelModifier.Config.HoldSymbolCodes, sc)
	}

	rd, isok := pool.Config.MapReels[reelModifier.Config.Reel]
	if !isok {
		goutils.Error("ReelModifier.InitEx:Reels",
			zap.String("reels", reelModifier.Config.Reel),
			zap.Error(ErrInvalidReels))

		return ErrInvalidReels
	}

	reelModifier.Config.ReelData = rd

	reelModifier.onInit(&reelModifier.Config.BasicComponentConfig)

	return nil
}

// procSymbolsRandPos
func (reelModifier *ReelModifier) canModify(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	for _, st := range reelModifier.Config.Triggers {
		if !gameProp.CanTrigger(st, gs, curpr, stake) {
			return false
		}
	}

	return true
}

// procSymbolsRandPos
func (reelModifier *ReelModifier) holdSymbol(src *sgc7game.GameScene, gs *sgc7game.GameScene) {
	if len(reelModifier.Config.HoldSymbolCodes) > 0 {
		for x, arr := range src.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(reelModifier.Config.HoldSymbolCodes, s, 0) >= 0 {
					gs.Arr[x][y] = s
				}
			}
		}
	}
}

// chgReel
func (reelModifier *ReelModifier) chgReel(gameProp *GameProperty, plugin sgc7plugin.IPlugin, src *sgc7game.GameScene, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	trynum := 0
	for {
		err := gs.RandReelsWithReelData(reelModifier.Config.ReelData, plugin)
		if err != nil {
			goutils.Error("ReelModifier.chgReel:RandReelsWithReelData",
				zap.Error(err))

			break
		}

		reelModifier.holdSymbol(src, gs)

		if reelModifier.canModify(gameProp, gs, curpr, stake) {
			return true
		}

		trynum++

		if trynum >= maxRetryNum {
			break
		}
	}

	return false
}

// chgReel
func (reelModifier *ReelModifier) chgReelWithMask(gameProp *GameProperty, plugin sgc7plugin.IPlugin, src *sgc7game.GameScene, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, mask string) bool {
	maskval, err := gameProp.Pool.GetMask(mask, gameProp)
	if err != nil {
		goutils.Error("ReelModifier.chgReelWithMask:GetMask",
			zap.Error(err))

		return false
	}

	trynum := 0
	for {
		err := gs.RandMaskReelsWithReelData(reelModifier.Config.ReelData, plugin, maskval, reelModifier.Config.IsReverse)
		if err != nil {
			goutils.Error("ReelModifier.chgReelWithMask:RandMaskReelsWithReelData",
				zap.Error(err))

			break
		}

		reelModifier.holdSymbol(src, gs)

		if reelModifier.canModify(gameProp, gs, curpr, stake) {
			return true
		}

		trynum++

		if trynum >= maxRetryNum {
			break
		}
	}

	return false
}

// playgame
func (reelModifier *ReelModifier) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) error {

	reelModifier.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	if reelModifier.Config.Mask != "" {
		gs := reelModifier.GetTargetScene3(gameProp, curpr, prs, bcd, reelModifier.Name, "", 0)
		gs1 := gs.CloneEx(gameProp.PoolScene)

		if reelModifier.chgReelWithMask(gameProp, plugin, gs, gs1, curpr, stake, reelModifier.Config.Mask) {
			reelModifier.AddScene(gameProp, curpr, gs1, bcd)
		}
	} else {
		gs := reelModifier.GetTargetScene3(gameProp, curpr, prs, bcd, reelModifier.Name, "", 0)
		gs1 := gs.CloneEx(gameProp.PoolScene)

		if reelModifier.chgReel(gameProp, plugin, gs, gs1, curpr, stake) {
			reelModifier.AddScene(gameProp, curpr, gs1, bcd)
		}
	}

	reelModifier.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (reelModifier *ReelModifier) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("reelModifier symbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (reelModifier *ReelModifier) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewStats2 -
func (reelModifier *ReelModifier) NewStats2() *stats2.Feature {
	return stats2.NewFeature(stats2.Options{stats2.OptStepTrigger})
}

// OnStats2
func (reelModifier *ReelModifier) OnStats2(icd IComponentData, s2 *stats2.Stats) {
	s2.PushStepTrigger(reelModifier.Name, true)
}

// // OnStats2Trigger
// func (reelModifier *ReelModifier) OnStats2Trigger(s2 *Stats2) {
// 	s2.pushTriggerStats(reelModifier.Name, true)
// }

func NewReelModifier(name string) IComponent {
	return &ReelModifier{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
