package lowcode

import (
	"log/slog"

	"github.com/zhs007/goutils"
)

type ComponentMgr struct {
	MapComponent     map[string]FuncNewComponent
	MapComponentData map[string]FuncNewComponentData // 未写完，写了一半，觉得clone效率更高，后续如果确定需要要再写
}

func (mgr *ComponentMgr) Reg(component string, funcNew FuncNewComponent) {
	mgr.MapComponent[component] = funcNew
}

func (mgr *ComponentMgr) RegComponentData(pbtype string, funcNewComponentData FuncNewComponentData) {
	mgr.MapComponentData[pbtype] = funcNewComponentData
}

func (mgr *ComponentMgr) NewComponent(cfgComponent *ComponentConfig) IComponent {
	funcNew, isok := mgr.MapComponent[cfgComponent.Type]
	if isok {
		return funcNew(cfgComponent.Name)
	}

	goutils.Error("ComponentMgr.NewComponent",
		slog.String("component", cfgComponent.Type),
		goutils.Err(ErrInvalidComponent))

	return nil
}

func NewComponentMgr() *ComponentMgr {
	mgr := &ComponentMgr{
		MapComponent:     make(map[string]FuncNewComponent),
		MapComponentData: make(map[string]FuncNewComponentData),
	}

	mgr.Reg(BasicReelsTypeName, NewBasicReels)
	mgr.Reg(CollectorTypeName, NewCollector)
	mgr.Reg(SymbolMultiTypeName, NewSymbolMulti)
	mgr.Reg(SymbolValTypeName, NewSymbolVal)
	mgr.Reg(SymbolValWinsTypeName, NewSymbolValWins)
	mgr.Reg(SymbolVal2TypeName, NewSymbolVal2)
	mgr.Reg(OverlaySymbolTypeName, NewOverlaySymbol)
	mgr.Reg(ChgSymbolTypeName, NewChgSymbol)
	mgr.Reg(RespinTypeName, NewRespin)
	mgr.Reg(MultiRespinTypeName, NewMultiRespin)
	mgr.Reg(ReplaceSymbolTypeName, NewReplaceSymbol)
	mgr.Reg(MaskTypeName, NewMask)
	mgr.Reg(FixSymbolsTypeName, NewFixSymbols)
	mgr.Reg(ReplaceReelTypeName, NewReplaceReel)
	mgr.Reg(MoveSymbolTypeName, NewMoveSymbol)
	mgr.Reg(MoveReelTypeName, NewMoveReel)
	mgr.Reg(MergeSymbolTypeName, NewMergeSymbol)
	mgr.Reg(ReRollReelTypeName, NewReRollReel)
	mgr.Reg(MultiWeightAwardsTypeName, NewMultiWeightAwards)
	mgr.Reg(MaskBranchTypeName, NewMaskBranch)
	mgr.Reg(SymbolModifierTypeName, NewSymbolModifier)
	mgr.Reg(ComponentTriggerTypeName, NewComponentTrigger)
	mgr.Reg(ComponentValTriggerTypeName, NewComponentValTrigger)
	mgr.Reg(ReelModifierTypeName, NewReelModifier)
	mgr.Reg(WeightReelsTypeName, NewWeightReels)
	mgr.Reg(ScatterTriggerTypeName, NewScatterTrigger)
	mgr.Reg(LinesTriggerTypeName, NewLinesTrigger)
	mgr.Reg(WaysTriggerTypeName, NewWaysTrigger)
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
	mgr.Reg(SymbolExpanderTypeName, NewSymbolExpander)
	mgr.Reg(PiggyBankTypeName, NewPiggyBank)
	mgr.Reg(AddSymbolsTypeName, NewAddSymbols)
	mgr.Reg(IntValMappingTypeName, NewIntValMapping)
	mgr.Reg(StringValMappingTypeName, NewStringValMapping)
	mgr.Reg(WeightBranchTypeName, NewWeightBranch)
	mgr.Reg(GenGigaSymbolTypeName, NewGenGigaSymbol)
	mgr.Reg(GenSymbolValsWithPosTypeName, NewGenSymbolValsWithPos)
	mgr.Reg(CheckSymbolValsTypeName, NewCheckSymbolVals)
	mgr.Reg(PositionCollectionTypeName, NewPositionCollection)
	mgr.Reg(ChgSymbolValsTypeName, NewChgSymbolVals)
	mgr.Reg(ChgSymbolsTypeName, NewChgSymbols)
	mgr.Reg(GenSymbolValsWithSymbolTypeName, NewGenSymbolValsWithSymbol)
	mgr.Reg(RebuildReelIndexTypeName, NewRebuildReelIndex)
	mgr.Reg(GenSymbolValsTypeName, NewGenSymbolVals)
	mgr.Reg(RebuildSymbolsTypeName, NewRebuildSymbols)
	mgr.Reg(RollNumberTypeName, NewRollNumber)
	mgr.Reg(ControllerWorkerTypeName, NewControllerWorker)
	mgr.Reg(CatchSymbolsTypeName, NewCatchSymbols)
	mgr.Reg(BurstSymbolsTypeName, NewBurstSymbols)
	mgr.Reg(WinResultModifierTypeName, NewWinResultModifier)
	mgr.Reg(ReelTriggerTypeName, NewReelTrigger)
	mgr.Reg(JackpotTypeName, NewJackpot)
	mgr.Reg(CheckValTypeName, NewCheckVal)
	mgr.Reg(AdjacentPayTriggerTypeName, NewAdjacentPayTrigger)
	mgr.Reg(WeightReels2TypeName, NewWeightReels2)
	mgr.Reg(WinResultModifierExTypeName, NewWinResultModifierEx)
	mgr.Reg(RandomMoveSymbolsTypeName, NewRandomMoveSymbols)
	mgr.Reg(GenPositionCollectionTypeName, NewGenPositionCollection)
	mgr.Reg(GenPositionCollectionWithSymbolTypeName, NewGenPositionCollectionWithSymbol)
	mgr.Reg(GenSymbolCollectionTypeName, NewGenSymbolCollection)
	mgr.Reg(DropSymbolsTypeName, NewDropSymbols)
	mgr.Reg(FeatureEmitterTypeName, NewFeatureEmitter)
	mgr.Reg(FeatureBarTypeName, NewFeatureBar)
	mgr.Reg(BombTypeName, NewBomb)
	mgr.Reg(SumSymbolValsTypeName, NewSumSymbolVals)
	mgr.Reg(TreasureChestTypeName, NewTreasureChest)
	mgr.Reg(MoveSymbols2TypeName, NewMoveSymbols2)
	mgr.Reg(GenSymbolVals2TypeName, NewGenSymbolVals2)
	mgr.Reg(MergePositionCollectionTypeName, NewMergePositionCollection)
	mgr.Reg(FeatureBar2TypeName, NewFeatureBar2)
	mgr.Reg(ChgSymbols2TypeName, NewChgSymbols2)
	mgr.Reg(FeaturePickTypeName, NewFeaturePick)
	mgr.Reg(ReelsCollectorTypeName, NewReelsCollector)
	mgr.Reg(FlowDownSymbolsTypeName, NewFlowDownSymbols)
	mgr.Reg(HoldAndWinTypeName, NewHoldAndWin)
	mgr.Reg(ChgSymbolsInReelsTypeName, NewChgSymbolsInReels)
	mgr.Reg(CascadingRegulatorTypeName, NewCascadingRegulator)
	mgr.Reg(WinResultLimiterTypeName, NewWinResultLimiter)
	mgr.Reg(SymbolValsSPTypeName, NewSymbolValsSP)
	mgr.Reg(Collector2TypeName, NewCollector2)
	mgr.Reg(DropDownSymbols2TypeName, NewDropDownSymbols2)
	mgr.Reg(HoldAndRespinReelsTypeName, NewHoldAndRespinReels)
	mgr.Reg(GenMaskTypeName, NewGenMask)
	mgr.Reg(BasicReels2TypeName, NewBasicReels2)
	mgr.Reg(RefillSymbols2TypeName, NewRefillSymbols2)
	mgr.Reg(GenSPGridTypeName, NewGenSPGrid)
	mgr.Reg(InitTropiCoolSPGridTypeName, NewInitTropiCoolSPGrid)
	mgr.Reg(DropDownTropiCoolSPGridTypeName, NewDropDownTropiCoolSPGrid)
	mgr.Reg(AlignTropiCoolSPGridTypeName, NewAlignTropiCoolSPGrid)
	mgr.Reg(RefillTropiCoolSPGridTypeName, NewRefillTropiCoolSPGrid)
	mgr.Reg(MergeSPGridTypeName, NewMergeSPGrid)
	mgr.Reg(CollectorPayTriggerTypeName, NewCollectorPayTrigger)
	mgr.Reg(CalculatorTypeName, NewCalculator)

	return mgr
}
