package lowcode

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type ComponentMgr struct {
	MapComponent map[string]FuncNewComponent
}

func (mgr *ComponentMgr) Reg(component string, funcNew FuncNewComponent) {
	mgr.MapComponent[component] = funcNew
}

func (mgr *ComponentMgr) NewComponent(cfgComponent *ComponentConfig) IComponent {
	funcNew, isok := mgr.MapComponent[cfgComponent.Type]
	if isok {
		return funcNew(cfgComponent.Name)
	}

	goutils.Error("ComponentMgr.NewComponent",
		zap.String("component", cfgComponent.Type),
		zap.Error(ErrInvalidComponent))

	return nil
}

func NewComponentMgr() *ComponentMgr {
	mgr := &ComponentMgr{
		MapComponent: make(map[string]FuncNewComponent),
	}

	mgr.Reg(BasicReelsTypeName, NewBasicReels)
	mgr.Reg(MysteryTypeName, NewMystery)
	// mgr.Reg(BasicWinsTypeName, NewBasicWins)
	// mgr.Reg(LightningTypeName, NewLightning)
	mgr.Reg(MultiLevelReelsTypeName, NewMultiLevelReels)
	mgr.Reg(CollectorTypeName, NewCollector)
	mgr.Reg(MultiLevelMysteryTypeName, NewMultiLevelMystery)
	// mgr.Reg(BookOfTypeName, NewBookOf)
	mgr.Reg(SymbolMultiTypeName, NewSymbolMulti)
	mgr.Reg(SymbolValTypeName, NewSymbolVal)
	mgr.Reg(SymbolValWinsTypeName, NewSymbolValWins)
	mgr.Reg(SymbolVal2TypeName, NewSymbolVal2)
	mgr.Reg(OverlaySymbolTypeName, NewOverlaySymbol)
	mgr.Reg(ReelSetMysteryTypeName, NewReelSetMystery)
	mgr.Reg(WeightTriggerTypeName, NewWeightTrigger)
	mgr.Reg(ChgSymbolTypeName, NewChgSymbol)
	mgr.Reg(RespinTypeName, NewRespin)
	mgr.Reg(MultiRespinTypeName, NewMultiRespin)
	mgr.Reg(ReplaceSymbolTypeName, NewReplaceSymbol)
	mgr.Reg(MaskTypeName, NewMask)
	mgr.Reg(MultiLevelReplaceReelTypeName, NewMultiLevelReplaceReel)
	mgr.Reg(FixSymbolsTypeName, NewFixSymbols)
	mgr.Reg(SymbolCollectionTypeName, NewSymbolCollection)
	mgr.Reg(WeightChgSymbolTypeName, NewWeightChgSymbol)
	// mgr.Reg(BookOf2TypeName, NewBookOf2)
	mgr.Reg(SymbolTriggerTypeName, NewSymbolTrigger)
	mgr.Reg(ReplaceReelTypeName, NewReplaceReel)
	mgr.Reg(MoveSymbolTypeName, NewMoveSymbol)
	mgr.Reg(MoveReelTypeName, NewMoveReel)
	mgr.Reg(MergeSymbolTypeName, NewMergeSymbol)
	mgr.Reg(ReRollReelTypeName, NewReRollReel)
	mgr.Reg(MultiWeightAwardsTypeName, NewMultiWeightAwards)
	mgr.Reg(MaskBranchTypeName, NewMaskBranch)
	// mgr.Reg(Respin2TypeName, NewRespin2)
	mgr.Reg(WeightTrigger2TypeName, NewWeightTrigger2)
	mgr.Reg(SymbolModifierTypeName, NewSymbolModifier)
	mgr.Reg(ComponentTriggerTypeName, NewComponentTrigger)
	mgr.Reg(ComponentValTriggerTypeName, NewComponentValTrigger)
	mgr.Reg(ReelModifierTypeName, NewReelModifier)
	mgr.Reg(WeightReelsTypeName, NewWeightReels)
	mgr.Reg(ScatterTriggerTypeName, NewScatterTrigger)
	mgr.Reg(LinesTriggerTypeName, NewLinesTrigger)
	mgr.Reg(WaysTriggerTypeName, NewWaysTrigger)
	mgr.Reg(WeightAwardsTypeName, NewWeightAwards)
	mgr.Reg(ClusterTriggerTypeName, NewClusterTrigger)
	mgr.Reg(RemoveSymbolsTypeName, NewRemoveSymbols)
	mgr.Reg(WinResultMultiTypeName, NewWinResultMulti)
	mgr.Reg(DropDownSymbolsTypeName, NewDropDownSymbols)
	mgr.Reg(RefillSymbolsTypeName, NewRefillSymbols)
	mgr.Reg(ReplaceSymbolGroupTypeName, NewReplaceSymbolGroup)
	mgr.Reg(RollSymbolTypeName, NewRollSymbol)
	mgr.Reg(QueueBranchTypeName, NewQueueBranch)
	mgr.Reg(SymbolCollection2TypeName, NewSymbolCollection2)
	mgr.Reg(ReplaceReelWithMaskTypeName, NewReplaceReelWithMask)

	return mgr
}
