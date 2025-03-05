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

// // LoadPB
// func (mgr *ComponentMgr) LoadPB(pb *anypb.Any) (IComponentData, error) {
// 	funcCD, isok := mgr.MapComponentData[pb.TypeUrl]
// 	if isok {
// 		icd := funcCD()
// 		if icd == nil {
// 			goutils.Error("ComponentMgr.LoadPB",
// 				goutils.Err(ErrInvalidFuncNewComponentData))

// 			return nil, ErrInvalidFuncNewComponentData
// 		}

// 		err := icd.LoadPB(pb)
// 		if icd == nil {
// 			goutils.Error("ComponentMgr.LoadPB:LoadPB",
// 				goutils.Err(err))

// 			return nil, err
// 		}

// 		return icd, nil
// 	}

// 	goutils.Error("ComponentMgr.LoadPB",
// 		goutils.Err(ErrInvalidAnypbTypeURL))

// 	return nil, ErrInvalidAnypbTypeURL
// }

func NewComponentMgr() *ComponentMgr {
	mgr := &ComponentMgr{
		MapComponent:     make(map[string]FuncNewComponent),
		MapComponentData: make(map[string]FuncNewComponentData),
	}

	mgr.Reg(BasicReelsTypeName, NewBasicReels)
	// mgr.Reg(MysteryTypeName, NewMystery)
	// mgr.Reg(BasicWinsTypeName, NewBasicWins)
	// mgr.Reg(LightningTypeName, NewLightning)
	// mgr.Reg(MultiLevelReelsTypeName, NewMultiLevelReels)
	mgr.Reg(CollectorTypeName, NewCollector)
	// mgr.Reg(MultiLevelMysteryTypeName, NewMultiLevelMystery)
	// mgr.Reg(BookOfTypeName, NewBookOf)
	mgr.Reg(SymbolMultiTypeName, NewSymbolMulti)
	mgr.Reg(SymbolValTypeName, NewSymbolVal)
	mgr.Reg(SymbolValWinsTypeName, NewSymbolValWins)
	mgr.Reg(SymbolVal2TypeName, NewSymbolVal2)
	mgr.Reg(OverlaySymbolTypeName, NewOverlaySymbol)
	// mgr.Reg(ReelSetMysteryTypeName, NewReelSetMystery)
	mgr.Reg(WeightTriggerTypeName, NewWeightTrigger)
	mgr.Reg(ChgSymbolTypeName, NewChgSymbol)
	mgr.Reg(RespinTypeName, NewRespin)
	mgr.Reg(MultiRespinTypeName, NewMultiRespin)
	mgr.Reg(ReplaceSymbolTypeName, NewReplaceSymbol)
	mgr.Reg(MaskTypeName, NewMask)
	// mgr.Reg(MultiLevelReplaceReelTypeName, NewMultiLevelReplaceReel)
	mgr.Reg(FixSymbolsTypeName, NewFixSymbols)
	// mgr.Reg(SymbolCollectionTypeName, NewSymbolCollection)
	mgr.Reg(WeightChgSymbolTypeName, NewWeightChgSymbol)
	// mgr.Reg(BookOf2TypeName, NewBookOf2)
	// mgr.Reg(SymbolTriggerTypeName, NewSymbolTrigger)
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
	mgr.Reg(PiggyBankTypeName, NewPiggyBank)
	mgr.Reg(AddSymbolsTypeName, NewAddSymbols)
	mgr.Reg(IntValMappingTypeName, NewIntValMapping)
	mgr.Reg(WeightBranchTypeName, NewWeightBranch)
	mgr.Reg(GenGigaSymbolTypeName, NewGenGigaSymbol)
	mgr.Reg(WinResultCacheTypeName, NewWinResultCache)
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
	mgr.Reg(FeatureBarTypeName, NewFeatureBar)
	mgr.Reg(BombTypeName, NewBomb)
	mgr.Reg(SumSymbolValsTypeName, NewSumSymbolVals)
	mgr.Reg(TreasureChestTypeName, NewTreasureChest)
	mgr.Reg(MoveSymbols2TypeName, NewMoveSymbols2)
	mgr.Reg(GenSymbolVals2TypeName, NewGenSymbolVals2)
	mgr.Reg(MergePositionCollectionTypeName, NewMergePositionCollection)

	return mgr
}
